package complete

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/sirkon/message"
)

// Log is used for debugging purposes
// since complete is running on tab completion, it is nice to
// have logs to the stderr (when writing your own completer)
// to write logs, set the COMP_DEBUG environment variable and
// use complete.Log in the complete program
var Log = getLogger()

func getLogger() func(format string, args ...interface{}) {
	var logfile io.Writer = ioutil.Discard
	if os.Getenv(envDebug) != "" {
		filePath := LogFilePath{}
		if ok, _ := filePath.Extract(os.Getenv(envDebug)); ok {
			var err error
			logfile, err = os.OpenFile(filePath.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				message.Error(err)
				logfile = os.Stderr
			}
		} else {
			logfile = os.Stderr
		}
	}
	return log.New(logfile, "complete ", log.Flags()).Printf
}
