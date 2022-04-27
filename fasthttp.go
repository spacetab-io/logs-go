package log

type FHLogger struct {
	Logger
}

func (fhl FHLogger) Printf(format string, v ...interface{}) {
	fhl.Logger.Printf(format, v...)
}
