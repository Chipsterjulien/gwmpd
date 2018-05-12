package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
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
	cmdToConsumeChan chan []byte
	mpdResponseChan  chan []byte
	sendCmdToMPDChan chan []byte
	mutex            *sync.Mutex
	info             *mpdInfos
}

type mpdInfos struct {
	currentPlaylist []mpdCurrentSong
	currentSong     *mpdCurrentSong
	stat            *mpdStat
	status          *mpdStatus
}

type mpdCurrentSong struct {
	Album        string
	Artist       string
	Date         string
	Duration     float64
	File         string
	Genre        string
	Id           int
	LastModified string
	Pos          int
	Time         int
	Title        string
	Track        int
}

type mpdStatus struct {
	Audio          string
	Bitrate        string
	Consume        bool
	Duration       float64
	Elapsed        float64
	Error          string
	MixrampDB      float64
	NextSong       int
	NextSongID     int
	Playlist       int
	PlaylistLength int
	Random         bool
	Repeat         bool
	Single         bool
	Song           int
	SongID         int
	State          string
	Time           string
	Volume         int
	VolumeSav      int
}

type mpdStat struct {
	Albums     int
	Artists    int
	DBPlaytime string
	DBUpdate   string
	Playtime   string
	Songs      int
	Uptime     string
}

type locationForm struct {
	Location string `form:"location" binding:"required"`
}

type playlistNameForm struct {
	PlaylistName string `form:"playlistName" binding:"required"`
}

type removeSongForm struct {
	PlaylistName string `form:"playlistName" json:"playlistName" binding:"required"`
	Pos          int    `form:"pos" json:"pos" binding:"exists"`
}

type renamePlaylistForm struct {
	OldName string `form:"oldName" json:"oldName" binding:"required"`
	NewName string `form:"newName" json:"newName" binding:"required"`
}

type songForm struct {
	PlaylistName string `form:"playlistName" json:"playlistName" binding:"required"`
	OldPos       int    `form:"oldPos" json:"oldPos" binding:"exists"`
	NewPos       int    `form:"newPos" binding:"exists"`
}

type volumeForm struct {
	Volume int `form:"volume" binding:"required"`
}

func (e *com) getAllPlaylists(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("listplaylists")
	playlists := []string{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "playlist":
			playlists = append(playlists, end)
		case "Last-Modified":
		default:
			log.Infof("In getAllFiles, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, playlists)
	e.mutex.Unlock()
}

func (e *com) clearPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var playlist playlistNameForm

	if err := c.ShouldBind(&playlist); err == nil {
		e.sendCmdToMPDChan <- append([]byte("playlistclear "), []byte(playlist.PlaylistName)...)
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In clearPlaylist, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"clearPlaylist": "ok"})
	} else {
		log.Warningf("Unable to clear playlist \"%v\": %s\n", playlist.PlaylistName, err)
	}
	e.mutex.Unlock()
}

func (e *com) getClearCurrentPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("clear")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getClearCurrentPlaylist, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"clearCurrentPlaylist": "ok"})
	e.mutex.Unlock()
}

func (e *com) getCurrentSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("currentsong")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "Album":
			e.info.currentSong.Album = end
		case "Artist":
			e.info.currentSong.Artist = end
		case "Composer":
		case "Date":
			e.info.currentSong.Date = end
		case "duration":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"duration\" %s", end)
				continue
			}
			e.info.currentSong.Duration = f
		case "file":
			e.info.currentSong.File = end
		case "Genre":
		case "Id":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Id\" %s", end)
				continue
			}
			e.info.currentSong.Id = i
		case "Last-Modified":
			e.info.currentSong.LastModified = end
		case "Pos":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Pos\" %s", end)
				continue
			}
			e.info.currentSong.Pos = i
		case "Title":
			e.info.currentSong.Title = end
		case "Time":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			e.info.currentSong.Time = i
		case "Track":
		default:
			log.Infof("In getCurrentSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{
		"Album":         e.info.currentSong.Album,
		"Artist":        e.info.currentSong.Artist,
		"Date":          e.info.currentSong.Date,
		"duration":      e.info.currentSong.Duration,
		"file":          e.info.currentSong.File,
		"Id":            e.info.currentSong.Id,
		"Last-Modified": e.info.currentSong.LastModified,
		"Pos":           e.info.currentSong.Pos,
		"Title":         e.info.currentSong.Title,
		"Time":          e.info.currentSong.Time,
	})
	e.mutex.Unlock()
}

func (e *com) getFilesList(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	location := c.Param("location")
	if len(location) > 0 {
		location = location[1:]
	}
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("listfiles \"%s\"", location))

	directories := []string{}
	songs := []mpdCurrentSong{}

	mySong := mpdCurrentSong{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			// e.mutex.Unlock()
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "directory":
			directories = append(directories, end)
		case "file":
			mySong.File = end
			songs = append(songs, mySong)
			mySong.File = ""
		case "Last-Modified":
		case "size":
		default:
			log.Infof("In getFilesList, unknown: \"%s\"\n", first)
		}
	}

	for pos, song := range songs {
		// e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("lsinfo \"%s\"", path.Join(location, song.File)))
		getSongInfos(e, &songs[pos])
		// e.mutex.Unlock()
	}

	infos := make(map[string]interface{})
	infos["directories"] = directories
	infos["songs"] = songs
	// os.Exit(0)

	c.JSON(200, infos)
	e.mutex.Unlock()
}

func getSongInfos(e *com, song *mpdCurrentSong) {
	log := logging.MustGetLogger("log")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			// e.mutex.Unlock()
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "Album":
			(*song).Album = end
		case "Artist":
			(*song).Artist = end
		case "Composer":
		case "Date":
		case "duration":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"duration\" %s", end)
				continue
			}
			(*song).Duration = f
		case "file":
		case "Genre":
		case "Last-Modified":
		case "Title":
			(*song).Title = end
		case "Time":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			(*song).Time = i
		case "Track":
		default:
			log.Infof("In getSongInfos, unknown: \"%s\"\n", first)
		}
	}
}

func (e *com) getLoadPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	name := c.Param("name")
	e.sendCmdToMPDChan <- append([]byte("load "), []byte(name)...)

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getloadPlaylist, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"loadPlaylist": name})
	e.mutex.Unlock()
}

func (e *com) getPreviousSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("previous")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPreviousSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"previousSong": "ok"})
	e.mutex.Unlock()
}

func (e *com) getNextSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("next")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getNextSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"nextSong": "ok"})
	e.mutex.Unlock()
}

func (e *com) getStopSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("stop")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getStopSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"stopSong": "ok"})
	e.mutex.Unlock()
}

func (e *com) getPlaylistSongsList(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	name := c.Param("name")
	e.sendCmdToMPDChan <- append([]byte("listplaylistinfo "), []byte(name)...)

	playlist := []mpdCurrentSong{}
	mySong := mpdCurrentSong{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "Album":
			mySong.Album = end
		case "Artist":
			mySong.Artist = end
		case "Composer":
		case "Date":
			mySong.Date = end
		case "duration":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"duration\" %s", end)
				continue
			}
			mySong.Duration = f
			playlist = append(playlist, mySong)
			mySong = mpdCurrentSong{}
		case "file":
			mySong.File = end
		case "Genre":
			mySong.Genre = end
		case "Last-Modified":
			mySong.LastModified = end
		case "Time":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			mySong.Time = i
		case "Title":
			mySong.Title = end
		case "Track":
		default:
			log.Infof("In getPlaylistSongsList, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, playlist)
	e.mutex.Unlock()
}

func (e *com) getPlaySong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("play")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPlaySong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"playSong": "ok"})
	e.mutex.Unlock()
}

func (e *com) getPlaySongWithID(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	pos := c.Param("pos")
	e.sendCmdToMPDChan <- append([]byte("play "), []byte(pos)...)

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPlaySong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"playSong": "ok"})
	e.mutex.Unlock()
}

func (e *com) getPauseSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("pause")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPauseSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"pauseSong": "ok"})
	e.mutex.Unlock()
}

func (e *com) getCurrentPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("playlistinfo")
	e.info.currentPlaylist = []mpdCurrentSong{}
	mySong := mpdCurrentSong{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "Album":
			mySong.Album = end
		case "Artist":
			mySong.Artist = end
		case "Composer":
		case "Date":
			mySong.Date = end
		case "duration":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"duration\" %s", end)
				continue
			}
			mySong.Duration = f
		case "file":
			mySong.File = end
		case "Genre":
			mySong.Genre = end
		case "Id":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Id\" %s", end)
				continue
			}
			mySong.Id = i
			e.info.currentPlaylist = append(e.info.currentPlaylist, mySong)
			mySong = mpdCurrentSong{}
		case "Last-Modified":
			mySong.LastModified = end
		case "Pos":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Pos\" %s", end)
				continue
			}
			mySong.Pos = i
		case "Time":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			mySong.Time = i
		case "Title":
			mySong.Title = end
		case "Track":
		default:
			log.Infof("In getCurrentPlaylist, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, e.info.currentPlaylist)
	e.mutex.Unlock()
}

func (e *com) getStatusMPD(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("status")
	errorIsNotDefine := true

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, end := splitLine(&line)
		switch first {
		case "audio":
			e.info.status.Audio = end
		case "bitrate":
			e.info.status.Bitrate = end
		case "consume":
			if end == "1" {
				e.info.status.Consume = true
			} else {
				e.info.status.Consume = false
			}
		case "duration":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"duration\" %s", end)
				continue
			}
			e.info.status.Duration = f
		case "elapsed":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"elapsed\" %s", end)
				continue
			}
			e.info.status.Elapsed = f
		case "error":
			errorIsNotDefine = false
			e.info.status.Error = end
		case "mixrampdb":
			f, err := strconv.ParseFloat(end, 64)
			if err != nil {
				log.Warningf("Unable to convert \"mixrampdb\" %s", end)
				continue
			}
			e.info.status.MixrampDB = f
		case "playlist":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"playlist\" %s", end)
				continue
			}
			e.info.status.Playlist = i
		case "playlistlength":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"playlistlength\" %s", end)
				continue
			}
			e.info.status.PlaylistLength = i
		case "random":
			if end == "1" {
				e.info.status.Random = true
			} else {
				e.info.status.Random = false
			}
		case "repeat":
			if end == "1" {
				e.info.status.Repeat = true
			} else {
				e.info.status.Repeat = false
			}
		case "single":
			if end == "1" {
				e.info.status.Single = true
			} else {
				e.info.status.Single = false
			}
		case "song":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"song\" %s", end)
				continue
			}
			e.info.status.Song = i
		case "songid":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Debug("coin")
				log.Warningf("Unable to convert \"songid\" %s", end)
				continue
			}
			e.info.status.SongID = i
		case "nextsong":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"nextsong\" %s", end)
				continue
			}
			e.info.status.NextSong = i
		case "nextsongid":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"nextsongid\" %s", end)
				continue
			}
			e.info.status.NextSongID = i
		case "state":
			e.info.status.State = end
		case "time":
			e.info.status.Time = end
		case "volume":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			e.info.status.Volume = i
		default:
			log.Infof("In getStatusMPD, unknown: \"%s\"\n", first)
		}
	}

	if errorIsNotDefine {
		e.info.status.Error = ""
	}

	c.JSON(200, gin.H{
		"audio":          e.info.status.Audio,
		"bitrate":        e.info.status.Bitrate,
		"consume":        e.info.status.Consume,
		"duration":       e.info.status.Duration,
		"elapsed":        e.info.status.Elapsed,
		"error":          e.info.status.Error,
		"mixrampdb":      e.info.status.MixrampDB,
		"playlist":       e.info.status.Playlist,
		"playlistlength": e.info.status.PlaylistLength,
		"random":         e.info.status.Random,
		"repeat":         e.info.status.Repeat,
		"single":         e.info.status.Single,
		"song":           e.info.status.Song,
		"songid":         e.info.status.SongID,
		"nextsong":       e.info.status.NextSong,
		"nextsongid":     e.info.status.NextSongID,
		"state":          e.info.status.State,
		"volume":         e.info.status.Volume,
	})
	e.mutex.Unlock()
}

func (e *com) moveSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var song songForm

	if err := c.ShouldBind(&song); err == nil {
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("playlistmove %s %d %d",
			song.PlaylistName,
			song.OldPos,
			song.NewPos,
		))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In moveSong, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"moveSong": "ok"})
	} else {
		log.Warningf("Unable to move song \"%v\" in playlist: %s\n", song.PlaylistName, err)
	}
	e.mutex.Unlock()
}

func (e *com) removePlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var playlistName playlistNameForm

	if err := c.ShouldBind(&playlistName); err == nil {
		e.sendCmdToMPDChan <- append([]byte("rm "), []byte(playlistName.PlaylistName)...)
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In removePlaylist, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"removePlaylist": "ok"})
	} else {
		log.Warningf("Unable to remove playlist \"%v\": %s\n", playlistName.PlaylistName, err)
	}
	e.mutex.Unlock()
}

func (e *com) removeSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var song removeSongForm

	if err := c.ShouldBind(&song); err == nil {
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("playlistdelete %s %d", song.PlaylistName, song.Pos))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In removeSong, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"removeSong": "ok"})
	} else {
		log.Warningf("Unable to remove song in \"%v\" at \"%d\": %s\n", song.PlaylistName, song.Pos, err)
	}
	e.mutex.Unlock()
}

func (e *com) renamePlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var name renamePlaylistForm

	if err := c.ShouldBind(&name); err == nil {
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("rename %s %s", name.OldName, name.NewName))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In renamePlaylist, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"renamePlaylist": "ok", "newName": name.NewName})
	} else {
		log.Warningf("Unable to rename playlist \"%ss\" to \"%s\": %s\n", name.OldName, name.NewName, err)
	}
	e.mutex.Unlock()
}

func (e *com) savePlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var playlistName playlistNameForm

	if err := c.ShouldBind(&playlistName); err == nil {
		// append([]byte("load "), []byte(name)...)
		// e.sendCmdToMPDChan <- []byte(fmt.Sprintf("save %s", playlistName.PlaylistName))
		e.sendCmdToMPDChan <- append([]byte("save "), []byte(playlistName.PlaylistName)...)
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In savePlaylist, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"savePlaylist": "ok"})
	} else {
		log.Warningf("Unable to save playlist \"%v\": %s\n", playlistName.PlaylistName, err)
	}
	e.mutex.Unlock()
}

func (e *com) setVolume(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	var vol volumeForm

	if err := c.ShouldBind(&vol); err == nil {
		e.info.status.Volume = vol.Volume
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", vol.Volume))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				break
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In setVolume, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"setVolume": "ok", "volume": e.info.status.Volume})
	} else {
		log.Warningf("Unable to set volume to \"%v\": %s\n", vol.Volume, err)
	}
	e.mutex.Unlock()
}

func (e *com) shuffle(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("shuffle")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In shuffle, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"shuffle": "ok"})
	e.mutex.Unlock()
}

func (e *com) toggleConsume(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.info.status.Consume = !e.info.status.Consume

	if e.info.status.Consume {
		e.sendCmdToMPDChan <- []byte("consume 1")
	} else {
		e.sendCmdToMPDChan <- []byte("consume 0")
	}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleConsume, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleConsume": "ok", "consume": e.info.status.Consume})
	e.mutex.Unlock()
}

func (e *com) toggleRandom(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.info.status.Random = !e.info.status.Random

	if e.info.status.Random {
		e.sendCmdToMPDChan <- []byte("random 1")
	} else {
		e.sendCmdToMPDChan <- []byte("random 0")
	}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleRandom, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleRandom": "ok", "random": e.info.status.Random})
	e.mutex.Unlock()
}

func (e *com) toggleRepeat(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.info.status.Repeat = !e.info.status.Repeat

	if e.info.status.Repeat {
		e.sendCmdToMPDChan <- []byte("repeat 1")
	} else {
		e.sendCmdToMPDChan <- []byte("repeat 0")
	}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleRepeat, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleRepeat": "ok", "repeat": e.info.status.Repeat})
	e.mutex.Unlock()
}

func (e *com) toggleSingle(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.info.status.Single = !e.info.status.Single

	if e.info.status.Single {
		e.sendCmdToMPDChan <- []byte("single 1")
	} else {
		e.sendCmdToMPDChan <- []byte("single 0")
	}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleSingle, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleSingle": "ok", "single": e.info.status.Single})
	e.mutex.Unlock()
}

func (e *com) toggleMuteVolume(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	if e.info.status.Volume == 0 {
		e.info.status.Volume = e.info.status.VolumeSav
		e.info.status.VolumeSav = 0
	} else {
		e.info.status.VolumeSav = e.info.status.Volume
		e.info.status.Volume = 0
	}
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", e.info.status.Volume))

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleMuteVolume, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleMute": "ok", "volume": e.info.status.Volume})
	e.mutex.Unlock()
}

func (e *com) updateDB(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("update")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		}

		first, _ := splitLine(&line)
		switch first {
		case "updating_db":
		default:
			log.Infof("In updateDB, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"updateDB": "ok"})
	e.mutex.Unlock()
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
		v1.GET("/allPlaylists", com.getAllPlaylists)
		v1.POST("/clearPlaylist", com.clearPlaylist)
		v1.GET("/clearCurrentPlaylist", com.getClearCurrentPlaylist)
		v1.GET("/currentPlaylist", com.getCurrentPlaylist)
		v1.GET("/currentSong", com.getCurrentSong)
		// v1.GET("/filesList", com.getFilesList)
		v1.GET("/filesList/*location", com.getFilesList)
		v1.GET("/loadPlaylist/:name", com.getLoadPlaylist)
		v1.GET("/pauseSong", com.getPauseSong)
		v1.GET("/playlistSongsList/:name", com.getPlaylistSongsList)
		v1.GET("/playSong", com.getPlaySong)
		v1.GET("/playSong/:pos", com.getPlaySongWithID)
		v1.GET("/previousSong", com.getPreviousSong)
		v1.POST("moveSong", com.moveSong)
		v1.GET("/nextSong", com.getNextSong)
		v1.POST("/removePlaylist", com.removePlaylist)
		v1.POST("/removeSong", com.removeSong)
		v1.POST("/renamePlaylist", com.renamePlaylist)
		v1.POST("/savePlaylist", com.savePlaylist)
		v1.POST("/setVolume", com.setVolume)
		v1.GET("/shuffle", com.shuffle)
		v1.GET("/statusMPD", com.getStatusMPD)
		v1.GET("/stopSong", com.getStopSong)
		v1.PUT("/toggleConsume", com.toggleConsume)
		v1.PUT("/toggleRandom", com.toggleRandom)
		v1.PUT("/toggleSingle", com.toggleSingle)
		v1.PUT("/toggleRepeat", com.toggleRepeat)
		v1.PUT("/toggleMuteVolume", com.toggleMuteVolume)
		v1.GET("/updateDB", com.updateDB)
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
	mpdResponseChan := make(chan []byte)
	sendCmdToMPDChan := make(chan []byte)
	cmdToConsumeChan := make(chan []byte)
	mutex := &sync.Mutex{}

	com := &com{
		cmdToConsumeChan: cmdToConsumeChan,
		mpdResponseChan:  mpdResponseChan,
		sendCmdToMPDChan: sendCmdToMPDChan,
		mutex:            mutex,

		info: &mpdInfos{
			currentPlaylist: []mpdCurrentSong{},
			currentSong:     &mpdCurrentSong{},
			stat:            &mpdStat{},
			status:          &mpdStatus{},
		},
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
			// If no action during 50s
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
				line = []byte("OK")
			}
			com.cmdToConsumeChan <- line
		}
	}
}
