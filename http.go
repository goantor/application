package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HttpServerOptions struct {
	Address      string `json:"address" toml:"address"`
	Port         int    `json:"port" toml:"port"`
	ReadTimeout  int    `json:"read_timeout" toml:"read_timeout"`
	WriteTimeout int    `json:"write_timeout" toml:"write_timeout"`
}

func (o HttpServerOptions) Addr() string {
	return fmt.Sprintf("%s:%d", o.Address, o.Port)
}

type HttpServer struct {
	server *http.Server
	entity http.Handler
	opt    *HttpServerOptions
}

func NewHttpServer(opt *HttpServerOptions) *HttpServer {
	return &HttpServer{entity: gin.New(), opt: opt}
}

func (h *HttpServer) Initialize() error {
	h.server = &http.Server{
		Addr:         h.opt.Addr(),
		Handler:      h.entity,
		ReadTimeout:  time.Duration(h.opt.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(h.opt.WriteTimeout) * time.Second,
	}

	return nil
}

func (h *HttpServer) Assign(entity interface{}) {
	entity = h.entity
}

func (h *HttpServer) Close() {
	h.Shutdown()
}

func (h *HttpServer) Boot() (err error) {
	// 需要修复
	err = h.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return
}

func (h *HttpServer) Handler() http.Handler {
	return h.entity
}

func (h *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
