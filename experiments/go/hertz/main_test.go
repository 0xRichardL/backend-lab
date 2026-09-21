package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
)

// Stage 0: baseline route behavior.

func TestRoute_Health(t *testing.T) {
	h := newServer()
	response := ut.PerformRequest(h.Engine, consts.MethodGet, "/health", nil).Result()

	assert.Equal(t, consts.StatusOK, response.StatusCode())
	assert.Equal(t, `{"status":"ok"}`, string(response.Body()))
}

// Stage 1: route matching, precedence, parameters, and missing routes.

func TestRoute_TasksGroup(t *testing.T) {
	h := newServer()

	type Test struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}

	tests := []Test{
		{
			name:           "static route takes priority",
			method:         consts.MethodGet,
			path:           "/tasks/recent",
			expectedStatus: consts.StatusOK,
			expectedBody:   `{"tasks":[]}`,
		},
		{
			name:           "parameter route captures ID",
			method:         consts.MethodGet,
			path:           "/tasks/42",
			expectedStatus: consts.StatusOK,
			expectedBody:   `{"id":"42"}`,
		},
		{
			name:           "group root is not registered",
			method:         consts.MethodGet,
			path:           "/tasks/",
			expectedStatus: consts.StatusNotFound,
			expectedBody:   `Not Found`,
		},
		{
			name:           "unsupported method returns not found",
			method:         consts.MethodPost,
			path:           "/tasks/42",
			expectedStatus: consts.StatusNotFound,
			expectedBody:   `Not Found`,
		},
		{
			name:           "parameter matches one segment only",
			method:         consts.MethodGet,
			path:           "/tasks/42/comments",
			expectedStatus: consts.StatusNotFound,
			expectedBody:   `Not Found`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := ut.PerformRequest(h.Engine,
				test.method,
				test.path,
				nil,
				ut.Header{Key: REQUEST_ID_HEADER,
					Value: "ID-01",
				}).Result()

			assert.Equal(t, test.expectedStatus, response.StatusCode())
			assert.Equal(t, test.expectedBody, string(response.Body()))
		})
	}
}

// Stage 2: JSON binding, validation, and stable public errors.

func TestRoute_CreateTask(t *testing.T) {
	h := newServer()

	t.Run("valid task", func(t *testing.T) {
		body := `{"title":"Learn Hertz"}`
		response := ut.PerformRequest(
			h.Engine,
			consts.MethodPost,
			"/tasks",
			&ut.Body{Body: bytes.NewBufferString(body), Len: len(body)},
			ut.Header{Key: consts.HeaderContentType, Value: consts.MIMEApplicationJSON},
			ut.Header{Key: REQUEST_ID_HEADER, Value: "ID-01"},
		).Result()

		assert.Equal(t, consts.StatusCreated, response.StatusCode())
		assert.Equal(t, consts.MIMEApplicationJSONUTF8, string(response.Header.ContentType()))
		assert.Equal(t, `{"title":"Learn Hertz"}`, string(response.Body()))
		assert.Equal(t, "ID-01", response.Header.Get(REQUEST_ID_HEADER))
	})
}

// Stage 3: middleware behavior and request-scoped values.

func TestMiddleware_RequestTiming(t *testing.T) {
	h := newServer()

	t.Run("valid request", func(t *testing.T) {
		response := ut.PerformRequest(h.Engine, consts.MethodGet, "/health", nil).Result()
		xResponseTimeHeader := response.Header.Get(RESPONSE_TIME_HEADER)
		assert.NotEmpty(t, xResponseTimeHeader)
		_, err := time.ParseDuration(xResponseTimeHeader)
		assert.NoError(t, err)
	})

	t.Run("response time remains in middleware rejection", func(t *testing.T) {
		response := ut.PerformRequest(h.Engine, consts.MethodGet, "/tasks/42", nil).Result()

		assert.Equal(t, consts.StatusBadRequest, response.StatusCode())
		xResponseTimeHeader := response.Header.Get(RESPONSE_TIME_HEADER)
		assert.NotEmpty(t, xResponseTimeHeader)
		_, err := time.ParseDuration(xResponseTimeHeader)
		assert.NoError(t, err)
	})
}

func TestMiddleware_RequestID(t *testing.T) {
	h := newServer()

	t.Run("valid request", func(t *testing.T) {
		response := ut.PerformRequest(
			h.Engine, consts.MethodGet,
			"/tasks/recent",
			nil,
			ut.Header{Key: REQUEST_ID_HEADER, Value: "ID-01"},
		).Result()
		assert.Equal(t, consts.StatusOK, response.StatusCode())
	})
	t.Run("missing request ID", func(t *testing.T) {
		response := ut.PerformRequest(h.Engine, consts.MethodGet, "/tasks/recent", nil).Result()

		assert.Equal(t, consts.StatusBadRequest, response.StatusCode())
		assert.JSONEq(t, `{"error":"missing request ID"}`, string(response.Body()))
	})

	t.Run("Health route doesn't require request ID", func(t *testing.T) {
		h := newServer()
		response := ut.PerformRequest(h.Engine, consts.MethodGet, "/health", nil).Result()

		assert.Equal(t, consts.StatusOK, response.StatusCode())
		assert.Equal(t, `{"status":"ok"}`, string(response.Body()))
	})

	t.Run("successful chain executes nested middleware in order", func(t *testing.T) {
		order := []string{}
		h := server.New()
		h.Use(func(ctx context.Context, c *app.RequestContext) {
			order = append(order, "outer before")
			c.Next(ctx)
			order = append(order, "outer after")
		})
		h.Use(requireRequestID())
		h.Use(func(ctx context.Context, c *app.RequestContext) {
			order = append(order, "inner before")
			c.Next(ctx)
			order = append(order, "inner after")
		})
		h.GET(
			"/probe",
			func(_ context.Context, c *app.RequestContext) {
				order = append(order, "handler")
				c.Status(consts.StatusNoContent)
			},
		)
		ut.PerformRequest(
			h.Engine,
			consts.MethodGet,
			"/probe",
			nil,
			ut.Header{Key: REQUEST_ID_HEADER, Value: "ID-01"},
		).Result()

		assert.Equal(t, []string{
			"outer before",
			"inner before",
			"handler",
			"inner after",
			"outer after",
		}, order)
	})

	t.Run("abort skips pending handlers but resumes outer middleware", func(t *testing.T) {
		order := []string{}
		h := server.New()
		h.Use(func(ctx context.Context, c *app.RequestContext) {
			order = append(order, "outer before")
			c.Next(ctx)
			order = append(order, "outer after")
		})
		h.Use(requireRequestID())
		h.Use(func(ctx context.Context, c *app.RequestContext) {
			order = append(order, "inner before")
			c.Next(ctx)
			order = append(order, "inner after")
		})
		h.GET(
			"/probe",
			func(_ context.Context, c *app.RequestContext) {
				order = append(order, "handler")
				c.Status(consts.StatusNoContent)
			},
		)
		ut.PerformRequest(h.Engine, consts.MethodGet, "/probe", nil).Result()

		assert.Equal(t, []string{
			"outer before",
			"outer after",
		}, order)
	})
}
