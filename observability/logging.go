package observability

import (
	"fmt"
	"os"
	"time"
)

var appName string
var loggingOn, loggingLevel string

// SetAppName -
func SetAppName(s string) {
	appName = s
}

// getLogHdr function
func getLogHdr() string {
	if appName != "" {
		return "[" + appName + "] "
	}
	return ""
}

// LogMemory prints memory usage to the trace
func LogMemory(errorType string) {
	logN(errorType, fmt.Sprintf("%s", getMemUsageStr()), 4)
}

// Logger is externalised for first level caller that doesnt care
func Logger(errorType string, logString string) {
	log(errorType, logString, 3)
}

// logN is not externalised to cater for internal callers from observability
func logN(errorType string, logString string, n int) {
	log(errorType, logString, n)
}

// Logger is externalised for first level caller that doesnt care
func Debug(logString string) {
	log("Debug", logString, 3)
}
func Info(logString string) {
	log("Info", logString, 3)
}
func Warn(logString string) {
	log("Warn", logString, 3)
}
func Error(logString string) {
	log("Error", logString, 3)
}
func Fatal(logString string) {
	log("Fatal", logString, 3)
}

func LogEnvVars() {
	for _, pair := range os.Environ() {
		log("Info", pair, 3)
	}
}

// Log wraps glog
func log(errorType string, logString string, n int) {

	if loggingOn == "" {
		loggingOn = os.Getenv("LOG_ENABLED")
		loggingLevel = os.Getenv("LOG_LEVEL")

		if loggingOn == "" {
			loggingOn = "true"
		}
		if loggingLevel == "" {
			loggingLevel = "DEBUG"
		}
	}

	if loggingOn == "true" {

		caller := Caller{}
		t := time.Now()

		// In a single-threaded process, the thread ID is equal to the process ID
		// p := fmt.Sprintf("%d:%d", syscall.Getpid(), syscall.Gettid())
		cor := GetCorrId()
		if cor == "" {
			cor = "NoCorrId"
		}

		caus := GetCausationId()
		if caus == "" {
			caus = "NoCausationId"
		}
		str := fmt.Sprintf("%s => %s", caus, cor)

		// format message
		msg := fmt.Sprintf("%s [%s] %s %s\t%s\n", t.Format("2006-01-02 15:04:05.0000"), str, getLogHdr(), caller.get(n), logString)

		if errorType == "Exit" {
			fmt.Fprintf(os.Stdout, "Q %s", msg)
			os.Exit(0)
		} else if errorType == "Fatal" {
			fmt.Fprintf(os.Stdout, "F %s", msg)
			os.Exit(3)
		} else if errorType == "Debug" && loggingLevel == "DEBUG" {
			fmt.Fprintf(os.Stdout, "D %s", msg)
		} else if errorType == "Info" && (loggingLevel == "INFO" || loggingLevel == "DEBUG") {
			fmt.Fprintf(os.Stdout, "I %s", msg)
		} else if errorType == "Warn" && (loggingLevel == "WARN" || loggingLevel == "INFO" || loggingLevel == "DEBUG") {
			fmt.Fprintf(os.Stdout, "W %s", msg)
		} else if errorType == "Error" && (loggingLevel == "ERROR" || loggingLevel == "WARN" || loggingLevel == "INFO" || loggingLevel == "DEBUG") {
			fmt.Fprintf(os.Stdout, "E %s", msg)
		} else if errorType == "" {
			fmt.Fprintf(os.Stdout, "? %s", msg)
		}
	}

}
