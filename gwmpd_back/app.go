package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

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

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:         "Restricted zone",
		Key:           []byte(viper.GetString("ginserver.jwtSecretKey")),
		Timeout:       time.Minute,
		MaxRefresh:    time.Minute,
		Authenticator: authenticator,
		Authorizator:  refreshToken,
		Unauthorized:  unauthorized,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}

	g.POST("/login", authMiddleware.LoginHandler)

	auth := g.Group("/v1")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/addSongToPlaylist", com.addSongToPlaylist)
		auth.GET("/allPlaylists", com.getAllPlaylists)
		auth.POST("/clearPlaylist", com.clearPlaylist)
		auth.GET("/clearCurrentPlaylist", com.getClearCurrentPlaylist)
		auth.GET("/currentPlaylist", com.getCurrentPlaylist)
		auth.GET("/currentSong", com.getCurrentSong)
		auth.GET("/filesList", com.getFilesList)
		auth.GET("/loadPlaylist", com.getLoadPlaylist)
		auth.GET("/pauseSong", com.getPauseSong)
		auth.GET("/playlistSongsList", com.getPlaylistSongsList)
		auth.GET("/playSong", com.getPlaySong)
		auth.GET("/previousSong", com.getPreviousSong)
		auth.POST("moveSong", com.moveSong)
		auth.GET("/nextSong", com.getNextSong)
		auth.GET("/refresh", authMiddleware.RefreshHandler)
		auth.POST("/removePlaylist", com.removePlaylist)
		auth.POST("/removeSong", com.removeSong)
		auth.POST("/renamePlaylist", com.renamePlaylist)
		auth.POST("/savePlaylist", com.savePlaylist)
		auth.POST("/setPositionTimeInCurrentSong", com.setPositionTimeInCurrentSong)
		auth.POST("/setVolume", com.setVolume)
		auth.GET("/shuffle", com.shuffle)
		auth.GET("/statusMPD", com.getStatusMPD)
		auth.GET("/stopSong", com.getStopSong)
		auth.PUT("/toggleConsume", com.toggleConsume)
		auth.PUT("/toggleRandom", com.toggleRandom)
		auth.PUT("/toggleSingle", com.toggleSingle)
		auth.PUT("/toggleRepeat", com.toggleRepeat)
		auth.PUT("/toggleMuteVolume", com.toggleMuteVolume)
		auth.GET("/updateDB", com.updateDB)
	}

	log.Debugf("Port: %d", viper.GetInt("ginserver.port"))
	if err := g.Run(":" + strconv.Itoa(viper.GetInt("ginserver.port"))); err != nil {
		log.Criticalf("Unable to start server: %s", err)
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

func main() {
	confPath := "/etc/gwmpd"
	confFilename := "gwmpd"
	logFilename := "/var/log/gwmpd/error.log"

	// confPath := "cfg/"
	// confFilename := "gwmpd_sample"
	// logFilename := "error.log"

	fd := initLogging(&logFilename)
	defer fd.Close()

	loadConfig(&confPath, &confFilename)
	startApp()
}
