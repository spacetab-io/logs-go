package log

type KafkaLogger struct {
	Logger
}

func NewKafkaLogger(l Logger) KafkaLogger {
	return KafkaLogger{Logger: l}
}

func (kl KafkaLogger) Printf(format string, v ...interface{}) {
	kl.Logger.Printf(format, v...)
}
