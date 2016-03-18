package Log

type LogContext struct {
	Tag string
}

// todo: log library leveled log
func (l *LogContext) LogMessagef(f string, v ...interface{}) {

}
