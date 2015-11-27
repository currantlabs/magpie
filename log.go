package magpie

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

func writeLog(format string, args ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, args...))
}
