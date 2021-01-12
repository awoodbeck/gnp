package ch13

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type wideResponseWriter struct {
	http.ResponseWriter
	length, status int
}

func (w *wideResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

func (w *wideResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.length += n

	if w.status == 0 {
		// 200 OK inferred on first Write if status is not yet set
		w.status = http.StatusOK
	}

	return n, err
}

func WideEventLog(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wideWriter := &wideResponseWriter{ResponseWriter: w}

			next.ServeHTTP(wideWriter, r)

			addr, _, _ := net.SplitHostPort(r.RemoteAddr)
			logger.Info("example wide event",
				zap.Int("status_code", wideWriter.status),
				zap.Int("response_length", wideWriter.length),
				zap.Int64("content_length", r.ContentLength),
				zap.String("method", r.Method),
				zap.String("proto", r.Proto),
				zap.String("remote_addr", addr),
				zap.String("uri", r.RequestURI),
				zap.String("user_agent", r.UserAgent()),
			)
		},
	)
}

func Example_wideLogEntry() {
	zl := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		),
	)
	defer func() { _ = zl.Sync() }()

	ts := httptest.NewServer(
		WideEventLog(zl, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func(r io.ReadCloser) {
					_, _ = io.Copy(ioutil.Discard, r)
					_ = r.Close()
				}(r.Body)
				_, _ = w.Write([]byte("Hello!"))
			},
		)),
	)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/test")
	if err != nil {
		zl.Fatal(err.Error())
	}
	_ = resp.Body.Close()

	// Output:
	// {"level":"info","msg":"example wide event","status_code":200,"response_length":6,"content_length":0,"method":"GET","proto":"HTTP/1.1","remote_addr":"127.0.0.1","uri":"/test","user_agent":"Go-http-client/1.1"}
}
