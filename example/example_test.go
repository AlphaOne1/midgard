// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestExampleMain(t *testing.T) {
	t.Parallel()

	go main()

	time.Sleep(500 * time.Millisecond)
	req, reqErr := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)

	if reqErr != nil {
		t.Errorf("unexpected request error: %v", reqErr)
	}

	req.Header.Add("Origin", "localhost")

	res, resErr := http.DefaultClient.Do(req)

	if resErr != nil {
		t.Errorf("got error for hello test page: %v", resErr)
	}

	if resErr == nil && res.StatusCode != http.StatusOK {
		body := make([]byte, res.ContentLength)
		_, _ = io.ReadFull(res.Body, body)
		t.Errorf("expected an OK HTTP status but got %v: %v", res.StatusCode, string(body))
	}

	time.Sleep(700 * time.Millisecond) // Give server time to process the shutdown
}
