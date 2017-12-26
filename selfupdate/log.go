package selfupdate

import (
	"io/ioutil"
	stdlog "log"
	"os"
)

var log = stdlog.New(ioutil.Discard, "", stdlog.Ltime)
var logEnabled = false

// EnableLog enables to output logging messages in library
func EnableLog() {
	if logEnabled {
		return
	}
	logEnabled = true
	log.SetOutput(os.Stderr)
}

// DisableLog disables to output logging messages in library
func DisableLog() {
	if !logEnabled {
		return
	}
	logEnabled = false
	log.SetOutput(ioutil.Discard)
}
