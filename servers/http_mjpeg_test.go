package servers

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/sources"
	. "github.com/kithix/videre/test_helpers"
)

func TestMultiPartHTTP(t *testing.T) {
	testStr := "jpegdatas"
	i := 0
	reader := sources.MakeFetcher(func() (io.Reader, error) {
		if i == 1 {
			return nil, io.EOF
		}
		i++
		return bytes.NewReader(SmallestJPEG), nil
	})
	w := httptest.NewRecorder()

	MJPEGHTTPHandler(reader, codec.NewPassthrough())(w, nil)

	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		t.Error("No content type header")
	}
	subStr := "multipart/x-mixed-replace;boundary="
	if !strings.Contains(contentType, subStr) {
		t.Error("Expected " + contentType + " to contain " + subStr)
	}

	if !strings.Contains(string(w.Body.Bytes()), testStr) {
		t.Error("Expected " + string(w.Body.Bytes()) + " to contain " + testStr)
	}
}
