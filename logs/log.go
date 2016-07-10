package logs

import "os"

//Logger stores logs from a function
type Logger struct {
	Logs chan string
}

//Log a message
func (l *Logger) Log(message string) {
	l.Logs <- message
}

//WriteToFile write all the logs to a file (such as stderr)
func (l *Logger) WriteToFile(file *os.File) error {
	var err error
	for log := range l.Logs {
		_, err = file.WriteString(log)
		if err != nil {
			break
		}
	}
	return err
}

//Kill closes the logger
func (l *Logger) Kill() {
	close(l.Logs)
}
