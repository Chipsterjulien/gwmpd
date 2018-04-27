package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

// Regarder ici pour l'authentification: https://github.com/appleboy/gin-jwt

// Pour les commandes mpd via le socket:
// https://www.musicpd.org/doc/protocol/command_reference.html

// Com ma chin chose à mettre correctement
// Pour plus d'info, voir ici: https://github.com/gin-gonic/gin/issues/932
type com struct {
	sendCmdToMPDChan              chan []byte
	mpdResponseChan               chan []byte
	waitPermissionToSendJSONAtVue chan bool
	info                          *mpdInfos
}

type mpdInfos struct {
	album          string
	albums         int
	artist         string
	artists        int
	audio          string
	bitrate        string
	consume        bool
	date           string
	dbPlaytime     string
	dbUpdate       string
	duration       float64
	elapsed        float64
	file           string
	genre          string
	id             int
	lastModified   string
	name           string
	nextSong       int
	nextSongID     int
	mixrampDB      float64
	playlist       int
	playlistLength int
	playtime       string
	pos            int
	random         bool
	repeat         bool
	single         bool
	state          string
	song           int
	songs          int
	songID         int
	timeSong       string
	timeElapsed    string
	title          string
	track          int
	uptime         string
	volume         int
	volumeSav      int
}

type volumeForm struct {
	Volume int `form:"volume" binding:"required"`
}

func (e *com) getStatMPD(c *gin.Context) {
	e.sendCmdToMPDChan <- []byte("currentsong")
	<-e.waitPermissionToSendJSONAtVue
	e.sendCmdToMPDChan <- []byte("status")
	<-e.waitPermissionToSendJSONAtVue
	e.sendCmdToMPDChan <- []byte("stats")
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{
		"album":          e.info.album,
		"albums":         e.info.albums,
		"artist":         e.info.artist,
		"artists":        e.info.artists,
		"audio":          e.info.audio,
		"bitrate":        e.info.bitrate,
		"consume":        e.info.consume,
		"date":           e.info.date,
		"dbplaytime":     e.info.dbPlaytime,
		"dbupdate":       e.info.dbUpdate,
		"duration":       e.info.duration,
		"elapsed":        e.info.elapsed,
		"file":           e.info.file,
		"genre":          e.info.genre,
		"id":             e.info.id,
		"Last-Modified":  e.info.lastModified,
		"mixrampdb":      e.info.mixrampDB,
		"Name":           e.info.name,
		"nextsong":       e.info.nextSong,
		"nextsongid":     e.info.nextSongID,
		"playlist":       e.info.playlist,
		"playlistlength": e.info.playlistLength,
		"playtime":       e.info.playtime,
		"pos":            e.info.pos,
		"random":         e.info.random,
		"repeat":         e.info.repeat,
		"single":         e.info.single,
		"state":          e.info.state,
		"song":           e.info.song,
		"songs":          e.info.songs,
		"songid":         e.info.songID,
		"TimeSong":       e.info.timeSong,
		"timeElapsed":    e.info.timeElapsed,
		"title":          e.info.title,
		"track":          e.info.track,
		"uptime":         e.info.uptime,
		"volume":         e.info.volume,
	})
}

func (e *com) getPreviousSong(c *gin.Context) {
	e.sendCmdToMPDChan <- []byte("previous")
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{"previousSong": "ok"})
}

func (e *com) getNextSong(c *gin.Context) {
	e.sendCmdToMPDChan <- []byte("next")
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{"nextSong": "ok"})
}

func (e *com) getStopSong(c *gin.Context) {
	e.sendCmdToMPDChan <- []byte("stop")
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{"stopSong": "ok"})
}

func (e *com) getPlaySong(c *gin.Context) {
	e.sendCmdToMPDChan <- []byte("play")
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{"playSong": "ok"})
}

func (e *com) getPauseSong(c *gin.Context) {
	e.sendCmdToMPDChan <- []byte("pause")
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{"pauseSong": "ok"})
}

func (e *com) setChangeVolume(c *gin.Context) {
	log := logging.MustGetLogger("log")
	var vol volumeForm

	if err := c.ShouldBind(&vol); err == nil {
		e.info.volume = vol.Volume
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", vol.Volume))
		<-e.waitPermissionToSendJSONAtVue
		c.JSON(200, gin.H{"setVolume": "ok", "volume": e.info.volume})
	} else {
		log.Warningf("Unable to set volume to \"%v\": %s", vol.Volume, err)
	}
}

func (e *com) toggleMuteVolume(c *gin.Context) {
	if e.info.volume == 0 {
		e.info.volume = e.info.volumeSav
		e.info.volumeSav = 0
	} else {
		e.info.volumeSav = e.info.volume
		e.info.volume = 0
	}
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", e.info.volume))
	<-e.waitPermissionToSendJSONAtVue
	c.JSON(200, gin.H{"toggleMute": "ok", "volume": e.info.volume})
}

func initGin(com *com) {
	log := logging.MustGetLogger("log")

	if viper.GetString("logtype") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.Default()
	g.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	v1 := g.Group("/v1")
	{
		v1.GET("/stateMPD", com.getStatMPD)
		v1.GET("/previousSong", com.getPreviousSong)
		v1.GET("/nextSong", com.getNextSong)
		v1.GET("/stopSong", com.getStopSong)
		v1.GET("/playSong", com.getPlaySong)
		v1.GET("/pauseSong", com.getPauseSong)
		v1.POST("/changeVolume", com.setChangeVolume)
		v1.PUT("/toggleMuteVolume", com.toggleMuteVolume)
	}

	log.Debugf("Port: %d", viper.GetInt("ginserver.port"))
	if err := g.Run(":" + strconv.Itoa(viper.GetInt("ginserver.port"))); err != nil {
		log.Criticalf("Unable to start serveur: %s", err)
		os.Exit(1)
	}
}

func initMPDSocket() net.Conn {
	log := logging.MustGetLogger("log")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("mpdserver.ip"), viper.GetInt("mpdserver.port")))
	if err != nil {
		log.Criticalf("Unable to connect to mpd server: %s", err)
		os.Exit(1)
	}

	return conn
}

func main() {
	confPath := "cfg/"
	confFilename := "gwmpd_sample"
	logFilename := "error.log"

	fd := initLogging(&logFilename)
	defer fd.Close()

	loadConfig(&confPath, &confFilename)
	startApp()
}

func startApp() {
	// Réponse venant de mpd
	mpdResponseChan := make(chan []byte, 100)
	// Commande à envoyer à mpd
	sendCmdToMPDChan := make(chan []byte, 100)
	// Autorisation pour envoyer les infos à la vue
	waitPermissionToSendJSONAtVue := make(chan bool)

	info := &mpdInfos{}
	com := &com{sendCmdToMPDChan: sendCmdToMPDChan,
		mpdResponseChan:               mpdResponseChan,
		waitPermissionToSendJSONAtVue: waitPermissionToSendJSONAtVue,
		info: info,
	}

	socket := initMPDSocket()
	go readLineProcess(mpdResponseChan, socket)
	go writeProcess(sendCmdToMPDChan, socket)
	go eventManagement(com, waitPermissionToSendJSONAtVue)

	initGin(com)
}

func readLineProcess(mpdResponseChan chan<- []byte, socket net.Conn) {
	log := logging.MustGetLogger("log")

	reader := bufio.NewReader(socket)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLineBytes()
		if err != nil {
			log.Criticalf("Unable to read line from %s: %s", viper.GetString("mpdserver.ip"), err)
			os.Exit(1)
		}

		mpdResponseChan <- line
	}
}

func writeProcess(sendCmdToMPDChan <-chan []byte, socket net.Conn) {
	log := logging.MustGetLogger("log")

	for {
		line := <-sendCmdToMPDChan
		if log.IsEnabledFor(5) {
			log.Debugf("< %s", line)
		}

		fmt.Fprintf(socket, fmt.Sprintf("%s\n", line))
	}
}

func eventManagement(com *com, waitPermissionToSendJSONAtVue chan<- bool) {
	ticker := time.NewTicker(55 * time.Second)
	sendPing := false

	for {
		select {
		case <-ticker.C:
			// Send ping to socket every 55s
			com.sendCmdToMPDChan <- []byte("ping")
			sendPing = true
		case line := <-com.mpdResponseChan:
			event(com, waitPermissionToSendJSONAtVue, &sendPing, &line)
		}
	}
}

func event(com *com, waitPermissionToSendJSONAtVue chan<- bool, sendPing *bool, line *[]byte) {
	log := logging.MustGetLogger("log")

	if log.IsEnabledFor(5) {
		log.Debugf("> %s", string(*line))
	}

	if bytes.Contains(*line, []byte("OK MPD")) {
		return
	} else if bytes.Equal(*line, []byte("OK")) {
		if *sendPing {
			*sendPing = false
		} else {
			waitPermissionToSendJSONAtVue <- true
		}
		return
	} else if bytes.Contains(*line, []byte("ACK")) {
		waitPermissionToSendJSONAtVue <- true
		return
	}

	lineSplitted := strings.Split(string(*line), ":")
	end := strings.TrimLeft(strings.Join(lineSplitted[1:], ":"), " ")
	switch lineSplitted[0] {
	case "Album":
		com.info.album = end
	case "albums":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"albums\" %s", end)
			return
		}
		com.info.albums = i
	case "Artist":
		com.info.artist = end
	case "artists":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"artists\" %s", end)
			return
		}
		com.info.artists = i
	case "audio":
		com.info.audio = end
	case "bitrate":
		com.info.bitrate = end
	case "consume":
		if end == "1" {
			com.info.consume = true
		} else {
			com.info.consume = false
		}
	case "Date":
		com.info.date = end
	case "db_playtime":
		com.info.dbPlaytime = end
	case "db_update":
		com.info.dbUpdate = end
	case "duration":
		f, err := strconv.ParseFloat(end, 64)
		if err != nil {
			log.Warningf("Unable to convert \"duration\" %s", end)
			return
		}
		com.info.duration = f
	case "elapsed":
		f, err := strconv.ParseFloat(end, 64)
		if err != nil {
			log.Warningf("Unable to convert \"elapsed\" %s", end)
			return
		}
		com.info.elapsed = f
	case "file":
		com.info.file = end
	case "Genre":
		com.info.genre = end
	case "Id":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"Id\" %s", end)
			return
		}
		com.info.id = i
	case "Last-Modified":
		com.info.lastModified = end
	case "mixrampdb":
		f, err := strconv.ParseFloat(end, 64)
		if err != nil {
			log.Warningf("Unable to convert \"mixrampdb\" %s", end)
			return
		}
		com.info.mixrampDB = f
	case "Name":
		com.info.name = end
	case "nextsong":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"nextsong\" %s", end)
			return
		}
		com.info.nextSong = i
	case "nextsongid":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"nextsongid\" %s", end)
			return
		}
		com.info.nextSongID = i
	case "playlist":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"playlist\" %s", end)
			return
		}
		com.info.playlist = i
	case "playlistlength":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"playlistlength\" %s", end)
			return
		}
		com.info.playlistLength = i
	case "playtime":
		com.info.playtime = end
	case "Pos":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"Pos\" %s", end)
			return
		}
		com.info.pos = i
	case "random":
		if end == "1" {
			com.info.random = true
		} else {
			com.info.random = false
		}
	case "repeat":
		if end == "1" {
			com.info.repeat = true
		} else {
			com.info.repeat = false
		}
	case "single":
		if end == "1" {
			com.info.single = true
		} else {
			com.info.single = false
		}
	case "song":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"song\" %s", end)
			return
		}
		com.info.song = i
	case "songs":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"songs\" %s", end)
			return
		}
		com.info.songs = i
	case "songid":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Debug("coin")
			log.Warningf("Unable to convert \"songid\" %s", end)
			return
		}
		com.info.songID = i
	case "state":
		com.info.state = end
	case "Time":
		com.info.timeSong = end
	case "time":
		com.info.timeElapsed = end
	case "Title":
		com.info.title = end
	case "Track":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"Track\" %s", end)
			return
		}
		com.info.track = i
	case "uptime":
		com.info.uptime = end
	case "volume":
		i, err := strconv.Atoi(end)
		if err != nil {
			log.Warningf("Unable to convert \"volume\" %s", end)
			return
		}
		com.info.volume = i
	default:
		log.Errorf("Unknown: \"%s\"\n", lineSplitted[0])
	}
}
