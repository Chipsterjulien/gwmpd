package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

// Regarder ici pour l'authentification: https://github.com/appleboy/gin-jwt

func initGin() {
	log := logging.MustGetLogger("log")

	if viper.GetString("logtype") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.Default()

	v1 := g.Group("api/v1")
	{
		v1.GET("/stateMpd", getStateMpd)
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
	jobsChan := make(chan []byte, 100)
	cmdToWriteChan := make(chan []byte, 100)

	socket := initMPDSocket()
	go readLineProcess(jobsChan, socket)
	go writeProcess(cmdToWriteChan, socket)
	go eventManagement(jobsChan, cmdToWriteChan)

	initGin()
}

func readLineProcess(jobsChan chan<- []byte, socket net.Conn) {
	log := logging.MustGetLogger("log")

	reader := bufio.NewReader(socket)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLineBytes()
		if err != nil {
			log.Criticalf("Unable to read line from %s: %s", viper.GetString("mpdserver.ip"), err)
			os.Exit(1)
		}

		jobsChan <- line
	}
}

func writeProcess(cmdToWriteChan <-chan []byte, socket net.Conn) {
	log := logging.MustGetLogger("log")

	for {
		line := <-cmdToWriteChan
		if log.IsEnabledFor(5) {
			log.Debugf("< %s", line)
		}

		fmt.Fprintf(socket, fmt.Sprintf("%s\n", line))
	}
}

func eventManagement(jobsChan <-chan []byte, cmdToWriteChan chan<- []byte) {

	ticker := time.NewTicker(55 * time.Second)

	for {
		select {
		case <-ticker.C:
			// Send ping to socket every 55s
			cmdToWriteChan <- []byte("ping")
		case line := <-jobsChan:
			event(jobsChan, cmdToWriteChan, &line)
		}
	}
}

func getStateMpd(c *gin.Context) {
}

func event(jobsChan <-chan []byte, cmdToWriteChan chan<- []byte, line *[]byte) {
	log := logging.MustGetLogger("log")

	if log.IsEnabledFor(5) {
		log.Debugf("> %s", string(*line))
	}

	if bytes.Contains(*line, []byte("OK")) {
		return
	}
}
