package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"inet.af/tcpproxy"
	"main/handle"
	"time"
)


func main() {

	host := flag.String("host", "0.0.0.0:8080", "http server listen to host")
	flag.Parse()

	tcpSrv := &handle.TcpServer{
		ProxyMap: make(map[string]*tcpproxy.Proxy),
	}

	httpSrv := &handle.HttpServer{
		ProxyMap: make(map[string]string),
	}

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(
		func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage)
	}))
	router.Use(gin.Recovery())

	apiTcpGroup := router.Group("/api/tcp")
	apiTcpGroup.GET("/query", tcpSrv.HandleProxyQuery)
	apiTcpGroup.POST("/create", tcpSrv.HandleProxyCreate)
	apiTcpGroup.POST("/delete", tcpSrv.HandleProxyDelete)

	apiHttpGroup := router.Group("/api/http")
	apiHttpGroup.GET("/query", httpSrv.HandleProxyQuery)
	apiHttpGroup.POST("/create", httpSrv.HandleProxyCreate)
	apiHttpGroup.POST("/delete", httpSrv.HandleProxyDelete)

	err := router.Run(*host)
	if err != nil {
		return
	}
}


