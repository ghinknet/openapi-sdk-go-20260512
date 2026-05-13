package shortlink

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.gh.ink/openapi/sdk/20260512/v3/client"
	"go.gh.ink/openapi/sdk/20260512/v3/errors"
	"go.gh.ink/toolbox/expr"
	"go.gh.ink/toolbox/pointer"
)

// Add a short link
func Add(c *client.Client, link string, validity *time.Time) (ok string, err error) {
	// Build payload
	payload := struct {
		Link     string `json:"link"`
		Validity *int64 `json:"validity"`
	}{
		Link: link,
		Validity: expr.Ternary(
			validity != nil,
			new(pointer.SafeDeref(validity).Unix()),
			nil,
		),
	}

	// Send request
	result := c.Send(
		strings.Join([]string{c.GetEndpoint(), Endpoint, "/add"}, ""),
		http.MethodPost,
		payload,
	).WithToken()
	if result.Err != nil {
		c.Logger.Error(nil, fmt.Sprintf(
			"failed to add short link, sender error: %s", result.Err.Error(),
		))
		return "", errors.ErrRequestSendFailed.
			WithApiCode(result.Code).
			WithApiMessage(result.Msg).
			WithRequestID(result.RequestID).
			WithResponse(result)
	}

	// Check status code
	if !result.OK() {
		c.Logger.Error(nil, fmt.Sprintf(
			"failed to add short link, upstream failed: code: %d, msg: %s", result.Code, result.Msg,
		))
		return "", errors.ErrStatusError.
			WithApiCode(result.Code).
			WithApiMessage(result.Msg).
			WithRequestID(result.RequestID).
			WithResponse(result)
	}

	// Build verify result struct
	var Link struct {
		LinkID string `json:"link_id"`
	}

	// Unmarshal token data
	if err = result.Unmarshal(&Link); err != nil {
		c.Logger.Error(nil, fmt.Sprintf(
			"failed to add short link, unmarshal error: %s", err.Error(),
		))
		return "", errors.ErrResponseUnmarshalFailed.
			WithApiCode(result.Code).
			WithApiMessage(result.Msg).
			WithRequestID(result.RequestID).
			WithResponse(result)
	}

	return Link.LinkID, nil
}
