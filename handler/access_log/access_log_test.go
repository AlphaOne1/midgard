// SPDX-FileCopyrightText: 2025 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

package access_log

import (
	"bytes"
	"encoding/base64"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

//nolint:paralleltest // testing output, manipulating global log behaviour
func TestAccessLogging(t *testing.T) {
	oldLog := slog.Default()
	defer slog.SetDefault(oldLog)

	logBuf := bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{})))

	handler := util.Must(New())(http.HandlerFunc(util.DummyHandler))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	slog.SetDefault(oldLog)

	messageMatch := regexp.MustCompile("msg=access")
	correlationMatch := regexp.MustCompile("correlation_id=")
	clientMatch := regexp.MustCompile("client=[^ ]+")
	methodMatch := regexp.MustCompile("method=GET")
	targetMatch := regexp.MustCompile("target=/")
	userMatch := regexp.MustCompile("user=")

	if !messageMatch.Match(logBuf.Bytes()) {
		t.Errorf("message not logged correctly: %v", logBuf.String())
	}
	if correlationMatch.Match(logBuf.Bytes()) {
		t.Errorf("correlation_id logged but not set: %v", logBuf.String())
	}
	if !clientMatch.Match(logBuf.Bytes()) {
		t.Errorf("client not logged correctly: %v", logBuf.String())
	}
	if !methodMatch.Match(logBuf.Bytes()) {
		t.Errorf("method not logged correctly: %v", logBuf.String())
	}
	if !targetMatch.Match(logBuf.Bytes()) {
		t.Errorf("target not logged correctly: %v", logBuf.String())
	}
	if userMatch.Match(logBuf.Bytes()) {
		t.Errorf("user logged but not set: %v", logBuf.String())
	}
}

//nolint:paralleltest // testing output, manipulating global log behaviour
func TestAccessLoggingCorrelationID(t *testing.T) {
	oldLog := slog.Default()
	defer slog.SetDefault(oldLog)

	logBuf := bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{})))

	handler := util.Must(New())(http.HandlerFunc(util.DummyHandler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Correlation-ID", "setOutside")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	slog.SetDefault(oldLog)

	correlationSetOutsideMatch := regexp.MustCompile("correlation_id=setOutside")

	if !correlationSetOutsideMatch.Match(logBuf.Bytes()) {
		t.Errorf("correlation_id from outside not logged correctly: %v", logBuf.String())
	}
}

//nolint:paralleltest // testing output, manipulating global log behaviour
func TestAccessLoggingUser(t *testing.T) {
	oldLog := slog.Default()
	defer slog.SetDefault(oldLog)

	logBuf := bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{})))

	handler := util.Must(New())(http.HandlerFunc(util.DummyHandler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:testpass")))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	slog.SetDefault(oldLog)

	userMatch := regexp.MustCompile("user=testuser")

	if !userMatch.Match(logBuf.Bytes()) {
		t.Errorf("user not logged correctly: %v", logBuf.String())
	}
}
