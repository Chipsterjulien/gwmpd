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

//************************************
// Exemple de playlist renvoyé par mpd
//************************************

// file: Mp3 varies/Dessins animés/Attaquantes_G�n�rique.mp3
// Last-Modified: 2008-05-26T21:12:13Z
// Artist: ?
// Title: Les attaquantes
// Album: Attaquantes
// Date: 1988
// Genre: Oldies
// Time: 155
// duration: 154.706
// Pos: 0
// Id: 1

// file: Toy-Box_-_Eenie_Meenie_Miney_Mo_(1999).mp3
// Last-Modified: 2015-11-26T13:17:22Z
// Time: 198
// duration: 198.217
// Pos: 1
// Id: 2

// file: Skyfall - Adele.mp3
// Last-Modified: 2014-09-03T13:04:02Z
// Artist: Adele
// Title: Skyfall - Adele
// Track: 4294967295
// Date: -1
// Time: 289
// duration: 289.336
// Pos: 2
// Id: 3

// file: Krewella-Alive.mp3
// Last-Modified: 2013-10-13T12:48:19Z
// Title: Krewella - Alive
// Track: 1
// Time: 207
// duration: 206.544
// Pos: 3
// Id: 4

// Regarder ici pour l'authentification: https://github.com/appleboy/gin-jwt

// Pour les commandes mpd via le socket:
// https://www.musicpd.org/doc/protocol/command_reference.html

// Com ma chin chose à mettre correctement
// Pour plus d'info, voir ici: https://github.com/gin-gonic/gin/issues/932
type com struct {
	sendCmdToMPDChan chan []byte
	mpdResponseChan  chan []byte
	cmdToConsume     chan []byte
	info             *mpdInfos
	// waitPermissionToSendJSONAtVue chan bool
}

type mpdInfos struct {
	currentSong *mpdCurrentSong
	status      *mpdStatus
	stat        *mpdStat
}

type mpdCurrentSong struct {
	album        string
	artist       string
	date         string
	duration     float64
	file         string
	genre        string
	id           int
	lastModified string
	pos          int
	time         int
	title        string
	track        int
}

type mpdStatus struct {
	audio          string
	bitrate        string
	consume        bool
	duration       float64
	elapsed        float64
	mixrampDB      float64
	nextSong       int
	nextSongID     int
	playlist       int
	playlistLength int
	random         bool
	repeat         bool
	single         bool
	song           int
	songID         int
	state          string
	time           string
	volume         int
	volumeSav      int
}

type mpdStat struct {
	albums     int
	artists    int
	dbPlaytime string
	dbUpdate   string
	playtime   string
	songs      int
	uptime     string
}

type volumeForm struct {
	Volume int `form:"volume" binding:"required"`
}

func (e *com) getCurrentSong(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("currentsong")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, end := splitLine(&line)
		switch first {
		case "file":
			e.info.currentSong.file = end
		case "Last-Modified":
			e.info.currentSong.lastModified = end
		case "Time":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			e.info.currentSong.time = i
		case "duration":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"duration\" %s", end)
				continue
			}
			e.info.currentSong.duration = f
		case "Pos":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Pos\" %s", end)
				continue
			}
			e.info.currentSong.pos = i
		case "Id":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Id\" %s", end)
				continue
			}
			e.info.currentSong.id = i
		default:
			log.Errorf("In getCurrentSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{
		"file":          e.info.currentSong.file,
		"Last-Modified": e.info.currentSong.lastModified,
		"Time":          e.info.currentSong.time,
		"duration":      e.info.currentSong.duration,
		"Pos":           e.info.currentSong.pos,
		"Id":            e.info.currentSong.id,
	})
}

func (e *com) getPreviousSong(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("previous")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Errorf("In getPreviousSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"previousSong": "ok"})
}

func (e *com) getNextSong(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("next")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Errorf("In getNextSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"nextSong": "ok"})
}

func (e *com) getStopSong(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("stop")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Errorf("In getStopSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"stopSong": "ok"})
}

func (e *com) getPlaySong(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("play")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Errorf("In getPlaySong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"playSong": "ok"})
}

func (e *com) getPauseSong(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("pause")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Errorf("In getPauseSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"pauseSong": "ok"})
}

func (e *com) getStatusMPD(c *gin.Context) {
	log := logging.MustGetLogger("log")
	e.sendCmdToMPDChan <- []byte("status")

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, end := splitLine(&line)
		switch first {
		case "consume":
			if end == "1" {
				e.info.status.consume = true
			} else {
				e.info.status.consume = false
			}
		case "mixrampdb":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"mixrampdb\" %s", end)
				continue
			}
			e.info.status.mixrampDB = f
		case "playlist":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"playlist\" %s", end)
				continue
			}
			e.info.status.playlist = i
		case "playlistlength":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"playlistlength\" %s", end)
				continue
			}
			e.info.status.playlistLength = i
		case "random":
			if end == "1" {
				e.info.status.random = true
			} else {
				e.info.status.random = false
			}
		case "repeat":
			if end == "1" {
				e.info.status.repeat = true
			} else {
				e.info.status.repeat = false
			}
		case "single":
			if end == "1" {
				e.info.status.single = true
			} else {
				e.info.status.single = false
			}
		case "song":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"song\" %s", end)
				continue
			}
			e.info.status.song = i
		case "songid":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Debug("coin")
				log.Warningf("Unable to convert \"songid\" %s", end)
				continue
			}
			e.info.status.songID = i
		case "nextsong":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"nextsong\" %s", end)
				continue
			}
			e.info.status.nextSong = i
		case "nextsongid":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"nextsongid\" %s", end)
				continue
			}
			e.info.status.nextSongID = i
		case "state":
			e.info.status.state = end
		case "volume":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			e.info.status.volume = i
		default:
			log.Errorf("Unknown: \"%s\"\n", first)
		}
	}
	c.JSON(200, gin.H{
		"consume":        e.info.status.consume,
		"mixrampdb":      e.info.status.mixrampDB,
		"playlist":       e.info.status.playlist,
		"playlistlength": e.info.status.playlistLength,
		"random":         e.info.status.random,
		"repeat":         e.info.status.repeat,
		"single":         e.info.status.single,
		"song":           e.info.status.song,
		"songid":         e.info.status.songID,
		"nextsong":       e.info.status.nextSong,
		"nextsongid":     e.info.status.nextSongID,
		"state":          e.info.status.state,
		"volume":         e.info.status.volume,
	})
}

func (e *com) getPlaylist(c *gin.Context) {
	// e.sendCmdToMPDChan <- []byte("playlistinfo")
	// <-e.waitPermissionToSendJSONAtVue
}

func (e *com) setChangeVolume(c *gin.Context) {
	log := logging.MustGetLogger("log")
	var vol volumeForm

	if err := c.ShouldBind(&vol); err == nil {
		e.info.status.volume = vol.Volume
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", vol.Volume))
		for {
			line := <-e.cmdToConsume
			if bytes.Equal(line, []byte("OK")) {
				return
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Errorf("In setChangeVolume, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"setVolume": "ok", "volume": e.info.status.volume})
	} else {
		log.Warningf("Unable to set volume to \"%v\": %s\n", vol.Volume, err)
	}
}

func (e *com) toggleMuteVolume(c *gin.Context) {
	log := logging.MustGetLogger("log")
	if e.info.status.volume == 0 {
		e.info.status.volume = e.info.status.volumeSav
		e.info.status.volumeSav = 0
	} else {
		e.info.status.volumeSav = e.info.status.volume
		e.info.status.volume = 0
	}
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", e.info.status.volume))

	for {
		line := <-e.cmdToConsume
		if bytes.Equal(line, []byte("OK")) {
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Errorf("In toggleMuteVolume, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleMute": "ok", "volume": e.info.status.volume})
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
		v1.POST("/changeVolume", com.setChangeVolume)
		v1.GET("/currentSong", com.getCurrentSong)
		v1.GET("/pauseSong", com.getPauseSong)
		v1.GET("/playSong", com.getPlaySong)
		v1.GET("/previousSong", com.getPreviousSong)
		v1.GET("/nextSong", com.getNextSong)
		v1.GET("/statusMPD", com.getStatusMPD)
		v1.GET("/stopSong", com.getStopSong)
		v1.PUT("/toggleMuteVolume", com.toggleMuteVolume)
		// v1.GET("/getPlaylist", com.getPlaylist)
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

func splitLine(line *[]byte) (string, string) {
	lineSplitted := strings.Split(string(*line), ":")
	end := strings.TrimLeft(strings.Join(lineSplitted[1:], ":"), " ")

	return lineSplitted[0], end
}

func startApp() {
	// Réponse venant de mpd
	mpdResponseChan := make(chan []byte, 100)
	// Commande à envoyer à mpd
	sendCmdToMPDChan := make(chan []byte, 100)
	// Commande à traiter en interne
	cmdToConsume := make(chan []byte, 100)
	// Autorisation pour envoyer les infos à la vue
	// waitPermissionToSendJSONAtVue := make(chan bool)

	// currentSong := &mpdCurrentSong{}
	// status := &mpdStatus{}
	// stat := &mpdStat{}

	// info := &mpdInfos{
	// 	currentSong: &mpdCurrentSong{},
	// 	status:      &mpdStatus{},
	// 	stat:        &mpdStat{},
	// }

	com := &com{
		sendCmdToMPDChan: sendCmdToMPDChan,
		mpdResponseChan:  mpdResponseChan,
		cmdToConsume:     cmdToConsume,

		info: &mpdInfos{
			currentSong: &mpdCurrentSong{},
			status:      &mpdStatus{},
			stat:        &mpdStat{},
		},
		// waitPermissionToSendJSONAtVue: waitPermissionToSendJSONAtVue,
	}

	socket := initMPDSocket()
	go readLineProcess(mpdResponseChan, socket)
	go writeProcess(sendCmdToMPDChan, socket)
	go eventManagement(com)

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

func eventManagement(com *com) {
	log := logging.MustGetLogger("log")

	// Consume that connect is ok
	if log.IsEnabledFor(5) {
		log.Debugf("> %s\n", <-com.mpdResponseChan)
	} else {
		<-com.mpdResponseChan
	}

	ticker := time.NewTicker(5 * time.Second)
	past := time.Now()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			delta := now.Sub(past)
			if delta > time.Duration(50)*time.Second {
				past = now
				com.sendCmdToMPDChan <- []byte("ping")

				if log.IsEnabledFor(5) {
					log.Debugf("> %s\n", <-com.mpdResponseChan)
				} else {
					<-com.mpdResponseChan
				}
			}
		case line := <-com.mpdResponseChan:
			log.Debugf("> %s\n", string(line))

			past = time.Now()
			if bytes.Contains(line, []byte("ACK")) {
				return
			}
			com.cmdToConsume <- line
		}
	}
}

// func event(com *com, waitPermissionToSendJSONAtVue chan<- bool, line *[]byte) {
// 	log := logging.MustGetLogger("log")
//
// 	if log.IsEnabledFor(5) {
// 		log.Debugf("> %s", string(*line))
// 	}
//
// 	if bytes.Contains(*line, []byte("OK MPD")) {
// 		return
// 	} else if bytes.Equal(*line, []byte("OK")) {
// 		waitPermissionToSendJSONAtVue <- true
// 		return
// 	} else if bytes.Contains(*line, []byte("ACK")) {
// 		waitPermissionToSendJSONAtVue <- true
// 		return
// 	}
//
// 	// lineSplitted := strings.Split(string(*line), ":")
// 	// end := strings.TrimLeft(strings.Join(lineSplitted[1:], ":"), " ")
// 	// switch lineSplitted[0] {
// 	// case "Album":
// 	// 	com.info.album = end
// 	// case "albums":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"albums\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.albums = i
// 	// case "Artist":
// 	// 	com.info.artist = end
// 	// case "artists":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"artists\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.artists = i
// 	// case "audio":
// 	// 	com.info.audio = end
// 	// case "bitrate":
// 	// 	com.info.bitrate = end
// 	// case "consume":
// 	// 	if end == "1" {
// 	// 		com.info.consume = true
// 	// 	} else {
// 	// 		com.info.consume = false
// 	// 	}
// 	// case "Date":
// 	// 	com.info.date = end
// 	// case "db_playtime":
// 	// 	com.info.dbPlaytime = end
// 	// case "db_update":
// 	// 	com.info.dbUpdate = end
// 	// case "duration":
// 	// 	f, err := strconv.ParseFloat(end, 64)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"duration\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.duration = f
// 	// case "elapsed":
// 	// 	f, err := strconv.ParseFloat(end, 64)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"elapsed\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.elapsed = f
// 	// case "file":
// 	// 	com.info.file = end
// 	// case "Genre":
// 	// 	com.info.genre = end
// 	// case "Id":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"Id\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.id = i
// 	// case "Last-Modified":
// 	// 	com.info.lastModified = end
// 	// case "mixrampdb":
// 	// 	f, err := strconv.ParseFloat(end, 64)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"mixrampdb\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.mixrampDB = f
// 	// case "Name":
// 	// 	com.info.name = end
// 	// case "nextsong":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"nextsong\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.nextSong = i
// 	// case "nextsongid":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"nextsongid\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.nextSongID = i
// 	// case "playlist":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"playlist\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.playlist = i
// 	// case "playlistlength":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"playlistlength\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.playlistLength = i
// 	// case "playtime":
// 	// 	com.info.playtime = end
// 	// case "Pos":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"Pos\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.pos = i
// 	// case "random":
// 	// 	if end == "1" {
// 	// 		com.info.random = true
// 	// 	} else {
// 	// 		com.info.random = false
// 	// 	}
// 	// case "repeat":
// 	// 	if end == "1" {
// 	// 		com.info.repeat = true
// 	// 	} else {
// 	// 		com.info.repeat = false
// 	// 	}
// 	// case "single":
// 	// 	if end == "1" {
// 	// 		com.info.single = true
// 	// 	} else {
// 	// 		com.info.single = false
// 	// 	}
// 	// case "song":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"song\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.song = i
// 	// case "songs":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"songs\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.songs = i
// 	// case "songid":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Debug("coin")
// 	// 		log.Warningf("Unable to convert \"songid\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.songID = i
// 	// case "state":
// 	// 	com.info.state = end
// 	// case "Time":
// 	// 	com.info.timeSong = end
// 	// case "time":
// 	// 	com.info.timeElapsed = end
// 	// case "Title":
// 	// 	com.info.title = end
// 	// case "Track":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"Track\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.track = i
// 	// case "uptime":
// 	// 	com.info.uptime = end
// 	// case "volume":
// 	// 	i, err := strconv.Atoi(end)
// 	// 	if err != nil {
// 	// 		log.Warningf("Unable to convert \"volume\" %s", end)
// 	// 		return
// 	// 	}
// 	// 	com.info.volume = i
// 	// default:
// 	// 	log.Errorf("Unknown: \"%s\"\n", lineSplitted[0])
// 	// }
// }
