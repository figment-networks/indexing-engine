package pipeline

var (
	logger Logger
)

// logInfo log with info message
func logInfo(msg string) {
	if logger != nil {
		logger.Info(msg)
	}
}

// logDebug log with debug
func logDebug(msg string) {
	if logger != nil {
		logger.Debug(msg)
	}
}
