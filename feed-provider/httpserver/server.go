package httpserver

import (
	"github.com/KyberNetwork/reserve-data/feed-provider/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//HTTPServer Httpserver that serve exporter
type HTTPServer struct {
	sugar *zap.SugaredLogger
	s     storage.Storage
	port  string
}

//NewHTTPServer new instance of HTTPServer
func NewHTTPServer(sugar *zap.SugaredLogger, s storage.Storage, port string) (*HTTPServer, error) {
	return &HTTPServer{
		sugar: sugar,
		s:     s,
		port:  port,
	}, nil
}

//Run start HTTPServer
func (h *HTTPServer) Run() error {
	r := gin.Default()
	r.GET("/feed/:name", func(c *gin.Context) {
		name := c.Param("name")
		data := h.s.Load(name)
		c.JSON(200, data)
	})
	return r.Run(h.port)
}
