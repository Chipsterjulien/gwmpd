package main

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	logging "github.com/op/go-logging"
)

func getSongInfos(e *com, song *mpdCurrentSong, location *string, songFile *string) error {
	log := logging.MustGetLogger("log")

	e.sendCmdToMPDChan <- []byte(fmt.Sprintf("lsinfo \"%s\"", path.Join(*location, *songFile)))

	for {
		line := <-e.cmdToConsumeChan
		if bytes.Equal(line, []byte("OK")) {
			return nil
		} else if bytes.Contains(line, []byte("ACK [50@0]")) {
			return errors.New(string(line))
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

func splitLine(line *[]byte) (string, string) {
	lineSplitted := strings.Split(string(*line), ":")
	end := strings.TrimLeft(strings.Join(lineSplitted[1:], ":"), " ")

	return lineSplitted[0], end
}
