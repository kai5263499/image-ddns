package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadURL(t *testing.T) {
	dir, err := ioutil.TempDir("", "image-ddns-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	previous := cfg
	defer func() { cfg = previous }()
	cfg.TmpDir = dir

	const payload = "offline image fixture"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if _, err := w.Write([]byte(payload)); err != nil {
			t.Errorf("write fixture: %v", err)
		}
	}))
	defer server.Close()

	name, err := downloadUrl(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(name) != dir {
		t.Fatalf("download escaped configured temporary directory: %s", name)
	}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != payload {
		t.Fatalf("download = %q, want %q", data, payload)
	}
}

func TestDownloadURLInvalidScheme(t *testing.T) {
	dir, err := ioutil.TempDir("", "image-ddns-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	previous := cfg
	defer func() { cfg = previous }()
	cfg.TmpDir = dir
	name, err := downloadUrl("invalid-scheme://offline-test")
	if err == nil || name != "" {
		t.Fatalf("downloadUrl returned (%q, %v), want empty path and error", name, err)
	}
}

func TestIPPortPattern(t *testing.T) {
	cases := []struct{ input, want string }{
		{"Address: 192.0.2.1:8080", "192.0.2.1"},
		{"Address: 192,0,2,1:8080", "192.0.2.1"},
	}
	for _, tc := range cases {
		matches := ipPortRegexp.FindStringSubmatch(tc.input)
		if len(matches) == 0 {
			t.Fatalf("no match for %q", tc.input)
		}
		got := strings.Split(strings.ReplaceAll(matches[0], ",", "."), ":")[0]
		if got != tc.want {
			t.Errorf("address = %q, want %q", got, tc.want)
		}
	}
	if ipPortRegexp.MatchString("no address present") {
		t.Error("unexpected address match")
	}
}

func TestCheckError(t *testing.T) {
	checkError("no failure", nil)
	want := errors.New("offline test error")
	defer func() {
		if got := recover(); got != want {
			t.Errorf("panic = %v, want original error %v", got, want)
		}
	}()
	checkError("expected failure", want)
}
