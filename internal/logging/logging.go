package logging

import "log"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func Info(format string, args ...any) {
	log.Printf(colorCyan+"INFO "+format+colorReset, args...)
}

func Warn(format string, args ...any) {
	log.Printf(colorYellow+"WARN "+format+colorReset, args...)
}

func Error(format string, args ...any) {
	log.Printf(colorRed+"ERROR "+format+colorReset, args...)
}

func Fatal(format string, args ...any) {
	log.Fatalf(colorRed+"FATAL "+format+colorReset, args...)
}
