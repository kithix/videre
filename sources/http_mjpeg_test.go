package sources

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
)

func mockMultiPartHTTPServer(contents ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mpw := multipart.NewWriter(w)
		w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+mpw.Boundary())
		w.Header().Add("Cache-Control", "no-cache")
		for _, content := range contents {
			headers := textproto.MIMEHeader{}
			headers.Add("Content-Type", "image/jpeg")
			headers.Add("Content-Length", strconv.Itoa(len(content)))
			part, err := mpw.CreatePart(headers)
			if err != nil {
				log.Println(err)
				return
			}
			part.Write([]byte(content))
		}
		mpw.Close()
	}
}

func TestHTTPMJPEGReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(mockMultiPartHTTPServer("test1", "test2", "test3")))
	r, err := HTTPMJPEGReader(HttpRequestDetails{
		URL: ts.URL,
	}, "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	buf := bytes.NewBuffer([]byte{})
	for {
		_, err = io.Copy(buf, r)
		if err != nil {
			if strings.Contains(err.Error(), io.EOF.Error()) {
				break
			} else {
				t.Error(err)
				break
			}
		}
	}
	if string(buf.Bytes()) != "test1test2test3" {
		t.Error("Buffer did not contain test1test2test3, instead contained " + string(buf.Bytes()))
	}
	if err = r.Close(); err != nil {
		t.Error(err)
	}
}
