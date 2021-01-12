package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodsHandler(t *testing.T) {
	handler := DefaultMethodsHandler()

	// test GET
	r := httptest.NewRequest(http.MethodGet, "http://test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %q", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if expected := "Hello, friend!"; expected != string(b) {
		t.Fatalf("expected %q; actual %q", expected, b)
	}

	// test POST
	r = httptest.NewRequest(http.MethodPost, "http://test",
		bytes.NewBufferString("<world>"))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %q", resp.Status)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if expected := "Hello, &lt;world&gt;!"; expected != string(b) {
		t.Fatalf("expected %q; actual %q", expected, b)
	}

	// test OPTIONS
	r = httptest.NewRequest(http.MethodOptions, "http://test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %q", resp.Status)
	}
	allow := resp.Header.Get("Allow")
	if allow == "" {
		t.Fatal("Allow header empty")
	}
	t.Logf("Allow: %s", allow)

	// test HEAD
	r = httptest.NewRequest(http.MethodHead, "http://test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	resp = w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status code: %q", resp.Status)
	}
}
