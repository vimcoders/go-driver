# go-driver
日志


	type Logger interface {
	
		Debug(format string, v ...interface{})
	
		Info(format string, v ...interface{})
	
		Warning(format string, v ...interface{})
	
		Error(format string, v ...interface{})
	
		Close() (err error)
	}
