package log

type FHLogger struct{}

func (fhl FHLogger) Printf(format string, v ...interface{}) {
	Printf(format, v...)
}
