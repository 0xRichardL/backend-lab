package main

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type CreateTaskRequest struct {
	Title string `json:"title,required" vd:"len($)>0"`
}

func newServer() *server.Hertz {
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8888"),
		// The shutdown budget starts when Spin receives a termination signal.
		server.WithExitWaitTime(30*time.Second),
	)

	// Server middleware wraps every matched route.
	h.Use(addRequestTiming())

	// Minimal health route and JSON response.
	h.GET("/health", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
	})

	// Route groups share a path prefix and can own middleware.
	tasks := h.Group("/tasks")
	tasks.Use(requireRequestID())
	tasks.GET("/recent", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"tasks": []any{}})
	})
	tasks.GET("/:id", func(_ context.Context, c *app.RequestContext) {
		id := c.Param("id")
		c.JSON(consts.StatusOK, utils.H{"id": id})
	})

	// BindAndValidate combines JSON binding and tag validation.
	tasks.POST("", createTask)

	// Live-process probes for an in-flight request during SIGTERM.
	// These sleeps simulate long-running work; they are not background workers.
	grateful := h.Group("/grateful")
	// This request can finish within the 30-second shutdown budget.
	grateful.GET("/slow", func(_ context.Context, c *app.RequestContext) {
		println("Entering /grateful/slow handler, sleeping for 10 seconds...")
		time.Sleep(10 * time.Second)
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
		println("Exiting /grateful/slow handler after 10 seconds.")
	})
	// Signal within the first five seconds so the remaining work exceeds the budget.
	grateful.GET("/exceed", func(_ context.Context, c *app.RequestContext) {
		println("Entering /grateful/exceed handler, sleeping for 35 seconds...")
		time.Sleep(35 * time.Second)
		c.JSON(consts.StatusOK, utils.H{"status": "ok"})
		println("Exiting /grateful/exceed handler after 35 seconds.")
	})

	return h
}

func main() {
	newServer().Spin()
}
