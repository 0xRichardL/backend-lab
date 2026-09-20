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

	h.GET("/health", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
	})

	tasks := h.Group("/tasks")
	tasks.GET("/recent", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"tasks": []any{}})
	})
	tasks.GET("/:id", func(_ context.Context, c *app.RequestContext) {
		id := c.Param("id")
		c.JSON(consts.StatusOK, utils.H{"id": id})
	})
	tasks.POST("", func(_ context.Context, c *app.RequestContext) {
		var request CreateTaskRequest
		if err := c.BindAndValidate(&request); err != nil {
			c.JSON(consts.StatusBadRequest, utils.H{"error": "invalid request"})
			return
		}

		c.JSON(consts.StatusCreated, utils.H{"title": request.Title})
	})

	return h
}

func main() {
	newServer().Spin()
}
