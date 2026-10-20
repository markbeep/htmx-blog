package route

import (
	"bytes"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
)

type contentEncoding uint8

const (
	identityEncoding contentEncoding = iota
	gzipEncoding
	brotliEncoding
)

type encodedResponse struct {
	status int
	header http.Header
	raw    []byte
	gzip   []byte
	brotli []byte
}

var compressedResponses sync.Map

// Compression serves brotli/gzip/identity responses. Responses are encoded on
// first request and kept in memory; the blog is fully static, so they never
// change. Already-compressed asset types are passed through untouched.
func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !compressiblePath(r.URL.Path) || conditionalRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		encoding := negotiateEncoding(r.Header.Get("Accept-Encoding"))
		if cached, ok := compressedResponses.Load(r.URL.Path); ok {
			writeEncoded(w, cached.(*encodedResponse), encoding)
			return
		}

		recorder := &bufferedResponse{header: make(http.Header)}
		next.ServeHTTP(recorder, r)

		response := &encodedResponse{
			status: recorder.Status(),
			header: recorder.header,
			raw:    recorder.body.Bytes(),
		}
		if len(response.raw) >= 512 {
			response.gzip = encodeGzip(response.raw)
			response.brotli = encodeBrotli(response.raw)
		}

		compressedResponses.Store(r.URL.Path, response)
		writeEncoded(w, response, encoding)
	})
}

func writeEncoded(w http.ResponseWriter, response *encodedResponse, encoding contentEncoding) {
	for key, values := range response.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Del("Content-Length")
	w.Header().Del("Accept-Ranges")
	w.Header().Del("Content-Encoding")

	body := response.raw
	if encoding == gzipEncoding && response.gzip != nil {
		w.Header().Set("Content-Encoding", "gzip")
		body = response.gzip
	} else if encoding == brotliEncoding && response.brotli != nil {
		w.Header().Set("Content-Encoding", "br")
		body = response.brotli
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(response.status)
	w.Write(body)
}

func encodeGzip(body []byte) []byte {
	var buffer bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	if err != nil {
		return nil
	}
	if _, err := writer.Write(body); err != nil {
		return nil
	}
	if err := writer.Close(); err != nil {
		return nil
	}
	return buffer.Bytes()
}

func encodeBrotli(body []byte) []byte {
	var buffer bytes.Buffer
	writer := brotli.NewWriterLevel(&buffer, brotli.BestCompression)
	if _, err := writer.Write(body); err != nil {
		return nil
	}
	if err := writer.Close(); err != nil {
		return nil
	}
	return buffer.Bytes()
}

func negotiateEncoding(header string) contentEncoding {
	brotliQuality, gzipQuality := -1.0, -1.0
	for part := range strings.SplitSeq(header, ",") {
		params := strings.Split(strings.TrimSpace(part), ";")
		name := strings.ToLower(strings.TrimSpace(params[0]))
		quality := 1.0
		for _, param := range params[1:] {
			param = strings.TrimSpace(param)
			if strings.HasPrefix(param, "q=") {
				if parsed, err := strconv.ParseFloat(strings.TrimSpace(param[2:]), 64); err == nil {
					quality = parsed
				}
			}
		}
		switch name {
		case "br":
			brotliQuality = quality
		case "gzip":
			gzipQuality = quality
		case "*":
			if brotliQuality < 0 {
				brotliQuality = quality
			}
			if gzipQuality < 0 {
				gzipQuality = quality
			}
		}
	}

	switch {
	case brotliQuality > 0 && brotliQuality >= gzipQuality:
		return brotliEncoding
	case gzipQuality > 0:
		return gzipEncoding
	default:
		return identityEncoding
	}
}

var precompressedExtensions = map[string]bool{
	".avif":     true,
	".br":       true,
	".gif":      true,
	".gz":       true,
	".ico":      true,
	".jpeg":     true,
	".jpg":      true,
	".mp3":      true,
	".mp4":      true,
	".otf":      true,
	".png":      true,
	".shortcut": true,
	".ttf":      true,
	".webm":     true,
	".webp":     true,
	".woff":     true,
	".woff2":    true,
	".zip":      true,
}

func compressiblePath(path string) bool {
	return !precompressedExtensions[strings.ToLower(filepath.Ext(path))]
}

// Conditional and range responses must not be cached: a 304 or 206 stored for
// a path would be replayed for later requests that expect the full response.
func conditionalRequest(r *http.Request) bool {
	return r.Header.Get("Range") != "" ||
		r.Header.Get("If-Range") != "" ||
		r.Header.Get("If-Modified-Since") != "" ||
		r.Header.Get("If-None-Match") != ""
}

type bufferedResponse struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func (b *bufferedResponse) Header() http.Header {
	return b.header
}

func (b *bufferedResponse) WriteHeader(status int) {
	if b.status != 0 {
		return
	}
	b.status = status
}

func (b *bufferedResponse) Write(body []byte) (int, error) {
	if b.status == 0 {
		b.status = http.StatusOK
	}
	return b.body.Write(body)
}

func (b *bufferedResponse) Status() int {
	if b.status == 0 {
		return http.StatusOK
	}
	return b.status
}
