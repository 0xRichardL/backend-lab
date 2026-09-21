package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func createTask(ctx context.Context, c *app.RequestContext) {
	var request CreateTaskRequest
	if err := c.BindAndValidate(&request); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}
	requestID := c.GetString(REQUEST_ID_KEY)
	c.Header(REQUEST_ID_HEADER, requestID)
	c.JSON(consts.StatusCreated, utils.H{"title": request.Title})
}
