package utils

import (
	"log"
	"os"
	"strings"

	"golift.io/xtractr"
)

// Logger satisfies the xtractr.Logger interface.
type Logger struct {
	xtractr *log.Logger
	debug   *log.Logger
	info    *log.Logger
}

// Printf satisfies the xtractr.Logger interface.
func (l *Logger) Printf(msg string, v ...interface{}) {
	l.xtractr.Printf(msg, v...)
}

// Debug satisfies the xtractr.Logger interface.
func (l *Logger) Debugf(msg string, v ...interface{}) {
	l.debug.Printf(msg, v...)
}

// Infof printf an info line.
func (l *Logger) Infof(msg string, v ...interface{}) {
	l.info.Printf(msg, v...)
}

func Extract(lid string) {
	logger := &Logger{
		xtractr: log.New(os.Stdout, "[xtractr] ", 0),
		debug:   log.New(os.Stdout, "[debug] ", 0),
		info:    log.New(os.Stdout, "[info] ", 0),
	}

	queue := xtractr.NewQueue(&xtractr.Config{
		Suffix:   "_xtractr",
		Logger:   logger,
		Parallel: 1,
		FileMode: 0644,
		DirMode:  0755,
	})

	defer queue.Stop() // Stop() waits for all extractions to finish.

	response := make(chan *xtractr.Response)

	queue.Extract(&xtractr.Xtract{
		Name: "Legendas Divx " + lid, // name is not import to this library.
		Filter: xtractr.Filter{
			Path: "tmp/downloads/" + lid,
		},
		CBChannel:  response, // queue responses are sent here.
		DeleteOrig: true,
	})

	// Queue always sends two responses. 1 on start and again when finished (error or not)
	resp := <-response // wait for the response.

	if (len(resp.Archives) == 0) || (resp.Error != nil) {
		logger.Printf("Extraction failed #1: %s", resp.Error)
		return
	}

	logger.Infof("Extraction started: %+v", resp.Archives)

	resp = <-response // wait for the response.

	if resp.Error != nil {
		logger.Printf("Extraction failed #2: %s", resp.Error)
	}

	logger.Infof("Extracted Files:\n - %s", strings.Join(resp.NewFiles, "\n - "))
}
