package main

import (
	"bytes"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	h := newServer()
	response := ut.PerformRequest(h.Engine, consts.MethodGet, "/health", nil).Result()

	assert.Equal(t, consts.StatusOK, response.StatusCode())
	assert.Equal(t, `{"status":"ok"}`, string(response.Body()))
}

func TestTasksGroup(t *testing.T) {
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
			response := ut.PerformRequest(h.Engine, test.Method, test.Path, nil).Result()
			assert.Equal(t, test.ExpectedStatus, response.StatusCode())
			assert.Equal(t, test.ExpectedBody, string(response.Body()))
		})
	}
}

func TestCreateTask(t *testing.T) {
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
			).Result()

			assert.Equal(t, test.ExpectedStatus, response.StatusCode())
			assert.Equal(t, consts.MIMEApplicationJSONUTF8, string(response.Header.ContentType()))
			assert.Equal(t, test.ExpectedBody, string(response.Body()))
		})
	}
}
