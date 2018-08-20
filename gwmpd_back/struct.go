package main

import "sync"

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
	ID           int
	LastModified string
	Name         string
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

type addSongForm struct {
	SongFilename string `form:"songFilename" binding:"required"`
	PlaylistName string `form:"playlistName" binding:"required"`
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

type seekForm struct {
	Position float64 `form:"position" json:"position" binding:"required"`
}

type songForm struct {
	PlaylistName string `form:"playlistName" json:"playlistName" binding:"required"`
	OldPos       int    `form:"oldPos" json:"oldPos" binding:"exists"`
	NewPos       int    `form:"newPos" binding:"exists"`
}

type volumeForm struct {
	Volume int `form:"volume" binding:"required"`
}
