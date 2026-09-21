package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type CreateTaskRequest struct {
	Title string `json:"title,required" vd:"len($)>0"`
}

func newServer() *server.Hertz {
	h := server.Default(server.WithHostPorts("127.0.0.1:8888"))

	// Stage 3: server middleware wraps every matched route.
	h.Use(addRequestTiming())

	// Stage 0: minimum route and JSON response.
	h.GET("/health", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
	})

	// Stage 1: groups share a path prefix and can own middleware.
	tasks := h.Group("/tasks")
	tasks.Use(requireRequestID())
	tasks.GET("/recent", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"tasks": []any{}})
	})
	tasks.GET("/:id", func(_ context.Context, c *app.RequestContext) {
		id := c.Param("id")
		c.JSON(consts.StatusOK, utils.H{"id": id})
	})

	// Stage 2: BindAndValidate combines JSON binding and tag validation.
	tasks.POST("", createTask)

	return h
}

func main() {
	newServer().Spin()
}
