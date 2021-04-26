package tool

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 非常严重的问题
)

func init() {
	errorFile, err := os.OpenFile("/var/log/proxy/errors.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	proxyFile, err := os.OpenFile("/var/log/proxy/proxy.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Info = log.New(io.MultiWriter(proxyFile, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(proxyFile, os.Stdout), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(errorFile, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
