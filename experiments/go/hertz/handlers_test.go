package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
)

func TestCreateTaskHandler(t *testing.T) {
	tests := []struct {
		name                    string
		body                    string
		expectedStatus          int
		expectedBody            string
		expectedRequestIDHeader string
	}{
		{
			name:                    "valid task",
			body:                    `{"title":"Learn Hertz"}`,
			expectedStatus:          consts.StatusCreated,
			expectedBody:            `{"title":"Learn Hertz"}`,
			expectedRequestIDHeader: "ID-01",
		},
		{
			name:                    "missing title",
			body:                    `{}`,
			expectedStatus:          consts.StatusBadRequest,
			expectedBody:            `{"error":"invalid request"}`,
			expectedRequestIDHeader: "",
		},
		{
			name:                    "empty title",
			body:                    `{"title":""}`,
			expectedStatus:          consts.StatusBadRequest,
			expectedBody:            `{"error":"invalid request"}`,
			expectedRequestIDHeader: "",
		},
		{
			name:                    "malformed JSON",
			body:                    `{"title":`,
			expectedStatus:          consts.StatusBadRequest,
			expectedBody:            `{"error":"invalid request"}`,
			expectedRequestIDHeader: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ut.CreateUtRequestContext(consts.MethodPost,
				"/tasks",
				&ut.Body{Body: bytes.NewBufferString(tt.body), Len: len(tt.body)},
				ut.Header{Key: consts.HeaderContentType, Value: consts.MIMEApplicationJSON},
			)
			c.Set(REQUEST_ID_KEY, "ID-01")
			createTask(context.Background(), c)
			assert.Equal(t, tt.expectedStatus, c.Response.StatusCode())
			assert.Equal(t, tt.expectedBody, string(c.Response.Body()))
			assert.Equal(t, consts.MIMEApplicationJSONUTF8, string(c.Response.Header.ContentType()))
			assert.Equal(t, tt.expectedRequestIDHeader, c.Response.Header.Get(REQUEST_ID_HEADER))
		})
	}
}
