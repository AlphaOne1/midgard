package access_log

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func TestAccessLogging(t *testing.T) {
	oldLog := slog.Default()
	defer slog.SetDefault(oldLog)

	logBuf := bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{})))

	handler := AccessLogging(http.HandlerFunc(DummyHandler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Correlation-ID", "setOutside")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	req.Header.Del("X-Correlation-ID")

	handler.ServeHTTP(rec, req)

	slog.SetDefault(oldLog)

	messageMatch := regexp.MustCompile("msg=access")
	correlationSetOutsideMatch := regexp.MustCompile("correlation_id=setOutside")
	correlationUnknownMatch := regexp.MustCompile("correlation_id=unknown")
	clientMatch := regexp.MustCompile("client=[^ ]+")
	methodMatch := regexp.MustCompile("method=GET")
	targetMatch := regexp.MustCompile("target=/")

	if !messageMatch.Match(logBuf.Bytes()) {
		t.Errorf("message not logged correctly: %v", logBuf.String())
	}
	if !correlationSetOutsideMatch.Match(logBuf.Bytes()) {
		t.Errorf("correlation_id from outside not logged correctly: %v", logBuf.String())
	}
	if !correlationUnknownMatch.Match(logBuf.Bytes()) {
		t.Errorf("correlation_id unset not logged correctly: %v", logBuf.String())
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
}
