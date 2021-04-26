package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inet.af/tcpproxy"
	"main/tool"
	"strings"
)

type TcpServer struct {
	ProxyMap map[string]*tcpproxy.Proxy
}

func (srv *TcpServer) HandleProxyCreate(c *gin.Context) {

	// 请求参数结构体
	var req struct {
		SrcAddr string `validate:"required"`
		DstAddr string `validate:"required"`
	}

	// 提取请求参数
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	// 判断是否已经添加过 proxy
	_, ok := srv.ProxyMap[req.SrcAddr + "/" +  req.DstAddr]
	if ok {
		tool.Info.Printf("SrcAddr:%s DstAddr:%s is already proxy\n", req.SrcAddr, req.DstAddr)
		c.JSON(200, gin.H{
			"body": "SrcAddr: " + req.SrcAddr +  " DstAddr: " + req.DstAddr +  " already proxy",
			"code": 200,
			"success": false,
		})
		return
	}

	// create tcp proxy
	go func() {
		var p tcpproxy.Proxy
		p.AddRoute(req.SrcAddr, tcpproxy.To(req.DstAddr))
		err := p.Start()
		if err != nil {
			return
		} else {
			go func() {
				err := p.Wait()
				if err != nil {
					return
				}
			}()
			srv.ProxyMap[req.SrcAddr + "/" +  req.DstAddr] = &p
		}
	}()

	tool.Info.Printf("Add proxy, src:%s to dst:%s\n", req.SrcAddr, req.DstAddr)
	c.JSON(200, gin.H{
		"body": "Add proxy " + "SrcAddr: " + req.SrcAddr +  " DstAddr: " + req.DstAddr,
		"code": 200,
		"success": true,
	})
}

func (srv *TcpServer) HandleProxyQuery(c *gin.Context) {

	// create response body data structure
	type ProxyItem struct {
		k string
		p *tcpproxy.Proxy
	}
	res := make([]ProxyItem, 0)

	// format response body step 1
	for k, v := range srv.ProxyMap {
		pi := ProxyItem{
			k: k,
			p: v,
		}
		res = append(res, pi)
	}

	if len(res) != 0 {
		// format response body step 2
		fmt.Println(res[0].k)
		response:= make(map[string]string)
		for i := 0; i < len(res); i++ {
			response[strings.Split(res[i].k, "/")[0]] = strings.Split(res[i].k, "/")[1]
		}

		// return response
		c.JSON(200, gin.H{
			"body": response,
			"code": 200,
			"success": true,
		})
		return
	}
	// return response
	c.JSON(200, gin.H{
		"body": "proxy list is nil",
		"code": 200,
		"success": true,
	})
	return
}

func (srv *TcpServer) HandleProxyDelete(c *gin.Context) {
	// 请求参数结构体
	var req struct {
		SrcAddr string `validate:"required"`
		DstAddr string `validate:"required"`
	}

	// 提取请求参数
	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	// 判断是否是已经添加过 proxy
	proxy, ok := srv.ProxyMap[req.SrcAddr + "/" +  req.DstAddr]
	if ok {
		err := proxy.Close()
		if err != nil {
			c.JSON(200, gin.H{
				"body": "Delete proxy " + "SrcAddr: " + req.SrcAddr +  " DstAddr: " + req.DstAddr,
				"code": 400,
				"success": false,
			})
			return
		}
		delete(srv.ProxyMap, req.SrcAddr + "/" +  req.DstAddr)
		tool.Info.Printf("Delete proxy, src:%s to dst:%s", req.SrcAddr, req.DstAddr)
		c.JSON(200, gin.H{
			"body": "Delete proxy " + "SrcAddr: " + req.SrcAddr +  " DstAddr: " + req.DstAddr,
			"code": 200,
			"success": true,
		})
		return
	}
	c.JSON(200, gin.H{
		"body": "Delete proxy " + "SrcAddr: " + req.SrcAddr +  " DstAddr: " + req.DstAddr,
		"code": 400,
		"success": false,
	})
	return
}


