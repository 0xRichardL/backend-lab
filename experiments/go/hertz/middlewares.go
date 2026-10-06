package main

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Middleware, post-handler work, and request-scoped values.

const RESPONSE_TIME_HEADER = "X-Response-Time"

func addRequestTiming() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		// Next runs the remaining chain before post-handler work resumes here.
		c.Next(ctx)
		duration := time.Since(start)
		c.Header(RESPONSE_TIME_HEADER, duration.String())
	}
}

const REQUEST_ID_HEADER = "X-Request-ID"
const REQUEST_ID_KEY = "REQUEST_ID_KEY"

func requireRequestID() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		requestID := string(c.GetHeader(REQUEST_ID_HEADER))
		if requestID == "" {
			c.AbortWithStatusJSON(consts.StatusBadRequest, utils.H{"error": "missing request ID"})
			return
		}
		// RequestContext values are available only during this request lifecycle.
		c.Set(REQUEST_ID_KEY, requestID)
		c.Next(ctx)
	}
}
