package route

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
)

func TestNegotiateEncoding(t *testing.T) {
	tests := []struct {
		header string
		want   contentEncoding
	}{
		{"", identityEncoding},
		{"br", brotliEncoding},
		{"gzip", gzipEncoding},
		{"gzip, br", brotliEncoding},
		{"gzip;q=1, br;q=0.5", gzipEncoding},
		{"br;q=0, gzip;q=1", gzipEncoding},
		{"br;q=0, gzip;q=0", identityEncoding},
		{"*", brotliEncoding},
		{"identity", identityEncoding},
	}

	for _, test := range tests {
		if got := negotiateEncoding(test.header); got != test.want {
			t.Errorf("negotiateEncoding(%q) = %v, want %v", test.header, got, test.want)
		}
	}
}

func TestConditionalRequestsAreNotCached(t *testing.T) {
	handler := Compression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-Modified-Since") != "" {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		io.WriteString(w, strings.Repeat("fresh content ", 100))
	}))

	conditional := httptest.NewRequest(http.MethodGet, "/conditional", nil)
	conditional.Header.Set("If-Modified-Since", "Wed, 21 Oct 2015 07:28:00 GMT")
	conditionalResponse := httptest.NewRecorder()
	handler.ServeHTTP(conditionalResponse, conditional)
	if conditionalResponse.Code != http.StatusNotModified {
		t.Fatalf("conditional status = %d, want %d", conditionalResponse.Code, http.StatusNotModified)
	}

	fresh := httptest.NewRequest(http.MethodGet, "/conditional", nil)
	freshResponse := httptest.NewRecorder()
	handler.ServeHTTP(freshResponse, fresh)
	if freshResponse.Code != http.StatusOK || freshResponse.Body.Len() == 0 {
		t.Fatalf("fresh response = %d with %d bytes, want 200 with a body", freshResponse.Code, freshResponse.Body.Len())
	}
}

func TestCompression(t *testing.T) {
	body := strings.Repeat("the quick brown fox ", 100)
	handler := Compression(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, body)
	}))

	tests := []struct {
		name     string
		path     string
		header   string
		encoding string
	}{
		{"brotli", "/brotli", "br", "br"},
		{"gzip", "/gzip", "gzip", "gzip"},
		{"negotiated gzip", "/negotiated", "gzip;q=1, br;q=0.5", "gzip"},
		{"identity", "/identity", "", ""},
		{"image passthrough", "/image.png", "br", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			if test.header != "" {
				request.Header.Set("Accept-Encoding", test.header)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if got := response.Header().Get("Content-Encoding"); got != test.encoding {
				t.Fatalf("Content-Encoding = %q, want %q", got, test.encoding)
			}

			var reader io.Reader = response.Body
			switch test.encoding {
			case "br":
				reader = brotli.NewReader(response.Body)
			case "gzip":
				gzipReader, err := gzip.NewReader(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				defer gzipReader.Close()
				reader = gzipReader
			}
			decoded, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(decoded, []byte(body)) {
				t.Fatalf("decoded body = %q, want %q", decoded, body)
			}
		})
	}
}
