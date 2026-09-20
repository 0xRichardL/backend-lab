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
		Name           string
		Method         string
		Path           string
		ExpectedStatus int
		ExpectedBody   string
	}

	tests := []Test{
		{
			Name:           "static route takes priority",
			Method:         consts.MethodGet,
			Path:           "/tasks/recent",
			ExpectedStatus: consts.StatusOK,
			ExpectedBody:   `{"tasks":[]}`,
		},
		{
			Name:           "parameter route captures ID",
			Method:         consts.MethodGet,
			Path:           "/tasks/42",
			ExpectedStatus: consts.StatusOK,
			ExpectedBody:   `{"id":"42"}`,
		},
		{
			Name:           "group root is not registered",
			Method:         consts.MethodGet,
			Path:           "/tasks/",
			ExpectedStatus: consts.StatusNotFound,
			ExpectedBody:   `Not Found`,
		},
		{
			Name:           "unsupported method returns not found",
			Method:         consts.MethodPost,
			Path:           "/tasks/42",
			ExpectedStatus: consts.StatusNotFound,
			ExpectedBody:   `Not Found`,
		},
		{
			Name:           "parameter matches one segment only",
			Method:         consts.MethodGet,
			Path:           "/tasks/42/comments",
			ExpectedStatus: consts.StatusNotFound,
			ExpectedBody:   `Not Found`,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			response := ut.PerformRequest(h.Engine,
				test.Method,
				test.Path,
				nil,
				ut.Header{Key: REQUEST_ID_HEADER,
					Value: "ID-01",
				}).Result()

			assert.Equal(t, test.ExpectedStatus, response.StatusCode())
			assert.Equal(t, test.ExpectedBody, string(response.Body()))
		})
	}
}

// Stage 2: JSON binding, validation, and stable public errors.

func TestRoute_CreateTask(t *testing.T) {
	h := newServer()

	type Test struct {
		Name           string
		Body           string
		ExpectedStatus int
		ExpectedBody   string
	}

	tests := []Test{
		{
			Name:           "valid task",
			Body:           `{"title":"Learn Hertz"}`,
			ExpectedStatus: consts.StatusCreated,
			ExpectedBody:   `{"title":"Learn Hertz"}`,
		},
		{
			Name:           "missing title",
			Body:           `{}`,
			ExpectedStatus: consts.StatusBadRequest,
			ExpectedBody:   `{"error":"invalid request"}`,
		},
		{
			Name:           "empty title",
			Body:           `{"title":""}`,
			ExpectedStatus: consts.StatusBadRequest,
			ExpectedBody:   `{"error":"invalid request"}`,
		},
		{
			Name:           "malformed JSON",
			Body:           `{"title":`,
			ExpectedStatus: consts.StatusBadRequest,
			ExpectedBody:   `{"error":"invalid request"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			body := &ut.Body{Body: bytes.NewBufferString(test.Body), Len: len(test.Body)}
			response := ut.PerformRequest(
				h.Engine,
				consts.MethodPost,
				"/tasks",
				body,
				ut.Header{Key: consts.HeaderContentType, Value: consts.MIMEApplicationJSON},
				ut.Header{Key: REQUEST_ID_HEADER, Value: "ID-01"},
			).Result()

			assert.Equal(t, test.ExpectedStatus, response.StatusCode())
			assert.Equal(t, consts.MIMEApplicationJSONUTF8, string(response.Header.ContentType()))
			assert.Equal(t, test.ExpectedBody, string(response.Body()))
		})
	}
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

	t.Run("request ID send back in response", func(t *testing.T) {
		bodyStr := `{"title":"Learn Hertz"}`
		body := &ut.Body{Body: bytes.NewBufferString(bodyStr), Len: len(bodyStr)}
		requestID := "ID-01"
		response := ut.PerformRequest(
			h.Engine,
			consts.MethodPost,
			"/tasks",
			body,
			ut.Header{Key: consts.HeaderContentType, Value: consts.MIMEApplicationJSON},
			ut.Header{Key: REQUEST_ID_HEADER, Value: requestID},
		).Result()
		respRequestID := response.Header.Get(REQUEST_ID_HEADER)

		assert.NotEmpty(t, respRequestID)
		assert.Equal(t, requestID, respRequestID)
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
