package handle

import "github.com/gin-gonic/gin"

type HttpServer struct {
	ProxyMap map[string]string
}

func (srv *HttpServer) HandleProxyQuery(c *gin.Context)  {
	// 请求参数结构体
	var req struct {
		ServerName string `validate:"required"`
		ServerUrl string `validate:"required"`
		DstAddr string `validate:"required"`
	}
}