package tutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type SlackHttpClientForTokenExpiresTest struct {
	OtherError bool

	TokenExpired       bool
	GetOAuthToken      bool
	GetOAuthTokenError bool
	SendMessageOK      bool

	invokeNums int // use to check Do(*http.Request) function invoked count
}

func (c *SlackHttpClientForTokenExpiresTest) GetInvokeNums() int {
	return c.invokeNums
}

// Do will return different result for test.
// If token expired set to true, it will return "token_expired" error,
// and will return normal response next time.
func (c *SlackHttpClientForTokenExpiresTest) Do(*http.Request) (*http.Response, error) {
	defer func() {
		c.invokeNums++
	}()

	var bodySource map[string]string

	switch {
	case c.TokenExpired:
		bodySource = map[string]string{
			"error": "token_expired",
		}

		c.TokenExpired = false
		c.GetOAuthToken = true
	case c.GetOAuthToken:
		if c.GetOAuthTokenError {
			return &http.Response{
				Body:       io.NopCloser(bytes.NewBufferString("")),
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		bodySource = map[string]string{
			"access_token":  "ACCESS_TOKEN",
			"refresh_token": "REFRESH_TOKEN",
		}

		c.GetOAuthToken = false
		c.SendMessageOK = true
	case c.SendMessageOK:
		bodySource = map[string]string{
			"text":       "OK",
			"channel":    "EXAMPLE",
			"ts":         strconv.Itoa(int(time.Now().UnixMilli())),
			"message_ts": strconv.Itoa(int(time.Now().UnixMilli())),
		}
	case c.OtherError:
		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString("")),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	bodyJson, _ := json.Marshal(bodySource)
	body := io.NopCloser(bytes.NewBuffer(bodyJson))

	return &http.Response{
		Body:       body,
		StatusCode: http.StatusOK,
	}, nil
}
