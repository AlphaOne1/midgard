// Copyright the midgard contributors.
// SPDX-License-Identifier: MPL-2.0
package access_log

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/AlphaOne1/midgard/util"
)

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

func TestOptionError(t *testing.T) {
	errOpt := func(h *Handler) error {
		return errors.New("testerror")
	}

	_, err := New(errOpt)

	if err == nil {
		t.Errorf("expected middleware creation to fail")
	}
}

func TestOptionWithLevel(t *testing.T) {
	h := util.Must(New(WithLogLevel(slog.LevelDebug)))(nil)

	if h.(*Handler).level != slog.LevelDebug {
		t.Errorf("wanted loglevel debug not set")
	}
}

func TestOptionWithLogger(t *testing.T) {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	h := util.Must(New(WithLogger(l)))(nil)

	if h.(*Handler).log != l {
		t.Errorf("logger not set correctly")
	}
}

func TestOptionWithNilLogger(t *testing.T) {
	var l *slog.Logger = nil
	_, hErr := New(WithLogger(l))

	if hErr == nil {
		t.Errorf("expected error on configuration with nil logger")
	}
}
