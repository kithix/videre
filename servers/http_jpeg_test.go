package servers

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/sources"
	. "github.com/kithix/videre/test_helpers"
)

func TestJPEGHTTPHander(t *testing.T) {

	fetch := sources.MakeFetcher(func() (io.Reader, error) {
		return bytes.NewReader(SmallestJPEG), nil
	})
	w := httptest.NewRecorder()
	JPEGHTTPHandler(fetch, codec.NewPassthrough())(w, nil)

	if !strings.Contains(string(w.Body.Bytes()), string(SmallestJPEG)) {
		t.Error("Expected " + strconv.Itoa(len(w.Body.Bytes())) + " to be " + strconv.Itoa(len(SmallestJPEG)))
	}
}
