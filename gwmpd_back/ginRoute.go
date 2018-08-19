package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	logging "github.com/op/go-logging"
)

func (e *com) addSongToPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var songForm addSongForm

	if err := c.ShouldBind(&songForm); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("playlistadd \"%s\" \"%s\"", songForm.PlaylistName, songForm.SongFilename))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"addSongToPlaylist": "failed"})

				return
			}

			first, _ := splitLine(&line)
			switch first {
			default:
				log.Infof("In clearPlaylist, unknown: \"%s\"\n", first)
			}
		}

		c.JSON(200, gin.H{"addSongToPlaylist": "ok", "playlistName": songForm.PlaylistName, "songFilename": songForm.SongFilename})
	} else {
		log.Warningf("Unable to add song in playlist: %s\n", err)
	}
}

func (e *com) clearPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var playlist playlistNameForm

	if err := c.ShouldBind(&playlist); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("playlistclear \"%s\"", playlist.PlaylistName))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"clearPlaylist": "failed"})

				return
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
}

func (e *com) getAllPlaylists(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("listplaylists")
	playlists := []string{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"allPlaylists": "failed"})

			return
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
}

func (e *com) getClearCurrentPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("clear")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"clearCurrentPlaylist": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getClearCurrentPlaylist, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"clearCurrentPlaylist": "ok"})
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
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"currentPlaylist": "failed"})

			return
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
}

func (e *com) getCurrentSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("currentsong")

	mySong := mpdCurrentSong{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"currentSong": "failed"})

			return
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
		case "Id":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Id\" %s", end)
				continue
			}
			mySong.Id = i
		case "Last-Modified":
			mySong.LastModified = end
		case "Name":
			mySong.Name = end
		case "Pos":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"Pos\" %s", end)
				continue
			}
			mySong.Pos = i
		case "Title":
			mySong.Title = end
		case "Time":
			i, err := strconv.Atoi(end)
			if err != nil {
				log.Warningf("Unable to convert \"volume\" %s", end)
				continue
			}
			mySong.Time = i
		case "Track":
		default:
			log.Infof("In getCurrentSong, unknown: \"%s\"\n", first)
		}
	}

	e.info.currentSong = &mySong

	c.JSON(200, gin.H{
		"Album":         e.info.currentSong.Album,
		"Artist":        e.info.currentSong.Artist,
		"Date":          e.info.currentSong.Date,
		"duration":      e.info.currentSong.Duration,
		"file":          e.info.currentSong.File,
		"Id":            e.info.currentSong.Id,
		"Last-Modified": e.info.currentSong.LastModified,
		"Name":          e.info.currentSong.Name,
		"Pos":           e.info.currentSong.Pos,
		"Title":         e.info.currentSong.Title,
		"Time":          e.info.currentSong.Time,
	})
}

func (e *com) getFilesList(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	location := c.DefaultQuery("location", "")
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("listfiles \"%s\"", location))

	directories := []string{}
	songs := []mpdCurrentSong{}

	mySong := mpdCurrentSong{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"filesList": "failed"})

			return
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

	newSongsList := []mpdCurrentSong{}
	for pos, song := range songs {
		if err := getSongInfos(e, &songs[pos], &location, &song.File); err == nil {
			newSongsList = append(newSongsList, songs[pos])
		}
	}

	e.mutex.Unlock()

	infos := make(map[string]interface{})
	infos["directories"] = directories
	infos["songs"] = newSongsList

	c.JSON(200, infos)
}

func (e *com) getLoadPlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	name := c.DefaultQuery("name", "")
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("load \"%s\"", name))

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"loadPlaylist": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getloadPlaylist, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"loadPlaylist": name})
}

func (e *com) getNextSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("next")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"nextSong": "ok"})
			return
		} else if bytes.Contains(line, []byte("ACK [55@0]")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"nextSong": "failed"})
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getNextSong, unknown: \"%s\"\n", first)
		}
	}
}

func (e *com) getPauseSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("pause")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"pauseSong": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPauseSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"pauseSong": "ok"})
}

func (e *com) getPlaylistSongsList(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	name := c.DefaultQuery("playlistName", "")
	if name == "" {
		e.mutex.Unlock()
		return
	}

	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("listplaylistinfo \"%s\"", name))

	playlist := []mpdCurrentSong{}
	mySong := mpdCurrentSong{}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"playlistSongsList": "failed"})

			return
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
			for _, protocol := range strings.Split("http:// https:// mms:// mmsh:// mmst:// mmsu:// gopher:// rtp:// rtsp:// rtmp:// rtmpt:// rtmps:// cdda:// alsa:// qobuz:// tidal://", " ") {
				if strings.HasPrefix(end, protocol) {
					playlist = append(playlist, mySong)
					mySong = mpdCurrentSong{}
					break
				}
			}
			// // file:// http:// https:// mms:// mmsh:// mmst:// mmsu:// gopher:// rtp:// rtsp:// rtmp:// rtmpt:// rtmps:// smb:// nfs:// cdda:// alsa:// qobuz:// tidal://
			// if strings.HasPrefix(end, "http://") || strings.HasPrefix(end, "https://") || strings.HasPrefix(end, "rtsp://") {
			// }
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

	// fmt.Println(playlist)
	// os.Exit(0)

	c.JSON(200, playlist)
}

func (e *com) getPlaySong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	pos := c.DefaultQuery("pos", "")
	if pos == "" {
		e.sendCmdToMPDChan <- []byte("play")
	} else {
		e.sendCmdToMPDChan <- append([]byte("play "), []byte(pos)...)
	}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"playSong": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPlaySong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"playSong": "ok"})
}

func (e *com) getPreviousSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("previous")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"previousSong": "ok"})
			return
		} else if bytes.Contains(line, []byte("ACK [55@0]")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"previousSong": "failed"})
			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getPreviousSong, unknown: \"%s\"\n", first)
		}
	}
}

func (e *com) getStopSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("stop")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"stopSong": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In getStopSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"stopSong": "ok"})
}

func (e *com) getStatusMPD(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("status")
	errorIsNotDefine := true

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"statusMPD": "failed"})

			return
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
}

func (e *com) moveSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var song songForm

	if err := c.ShouldBind(&song); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("playlistmove \"%s\" %d %d",
			song.PlaylistName,
			song.OldPos,
			song.NewPos,
		))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"moveSong": "failed"})

				return
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
}

func (e *com) removePlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var playlistName playlistNameForm

	if err := c.ShouldBind(&playlistName); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("rm \"%s\"", playlistName.PlaylistName))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"removePlaylist": "failed"})

				return
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
}

func (e *com) removeSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var song removeSongForm

	if err := c.ShouldBind(&song); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("playlistdelete \"%s\" \"%d\"", song.PlaylistName, song.Pos))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"removeSong": "failed"})

				return
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
}

func (e *com) renamePlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var name renamePlaylistForm

	if err := c.ShouldBind(&name); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("rename \"%s\" \"%s\"", name.OldName, name.NewName))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"renamePlaylist": "failed"})

				return
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
}

func (e *com) savePlaylist(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var playlistName playlistNameForm

	if err := c.ShouldBind(&playlistName); err == nil {
		e.mutex.Lock()
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("save \"%s\"", playlistName.PlaylistName))
		for {
			line := <-e.cmdToConsumeChan
			if bytes.Equal(line, []byte("OK")) {
				e.mutex.Unlock()
				break
			} else if bytes.Contains(line, []byte("ACK")) {
				e.mutex.Unlock()
				c.JSON(200, gin.H{"savePlaylist": "failed"})

				return
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
}

func (e *com) setPositionTimeInCurrentSong(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var seek seekForm

	c.ShouldBind(&seek)
	e.mutex.Lock()

	if seek.Position <= 0 {
		e.sendCmdToMPDChan <- []byte("seekcur 0")
		e.info.status.Elapsed = 0
	} else if seek.Position > e.info.status.Duration {
		e.sendCmdToMPDChan <- []byte("next")
		e.info.status.Elapsed = e.info.status.Duration
	} else {
		e.sendCmdToMPDChan <- []byte(fmt.Sprintf("seekcur %f", seek.Position))
		e.info.status.Elapsed = seek.Position
	}

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"setPosition": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In setPositionTimeInCurrentSong, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"setPosition": "ok", "position": e.info.status.Elapsed})
}

func (e *com) setVolume(c *gin.Context) {
	log := logging.MustGetLogger("log")

	var vol volumeForm

	c.ShouldBind(&vol)
	if vol.Volume < 0 || vol.Volume > 100 {
		c.JSON(200, gin.H{"setVolume": "ok", "volume": e.info.status.Volume})
		return
	}

	e.mutex.Lock()
	e.info.status.Volume = vol.Volume
	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("setvol %d", vol.Volume))
	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"setVolume": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In setVolume, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"setVolume": "ok", "volume": e.info.status.Volume})
}

func (e *com) shuffle(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("shuffle")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"shuffle": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In shuffle, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"shuffle": "ok"})
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
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"toggleConsume": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleConsume, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleConsume": "ok", "consume": e.info.status.Consume})
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
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"toggleRandom": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleRandom, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleRandom": "ok", "random": e.info.status.Random})
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
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"toggleRepeat": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleRepeat, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleRepeat": "ok", "repeat": e.info.status.Repeat})
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
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"toggleSingle": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleSingle, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleSingle": "ok", "single": e.info.status.Single})
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
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"toggleMuteVolume": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		default:
			log.Infof("In toggleMuteVolume, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"toggleMute": "ok", "volume": e.info.status.Volume})
}

func (e *com) updateDB(c *gin.Context) {
	log := logging.MustGetLogger("log")

	e.mutex.Lock()
	e.sendCmdToMPDChan <- []byte("update")

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			e.mutex.Unlock()
			break
		} else if bytes.Contains(line, []byte("ACK")) {
			e.mutex.Unlock()
			c.JSON(200, gin.H{"updateDB": "failed"})

			return
		}

		first, _ := splitLine(&line)
		switch first {
		case "updating_db":
		default:
			log.Infof("In updateDB, unknown: \"%s\"\n", first)
		}
	}

	c.JSON(200, gin.H{"updateDB": "ok"})
}