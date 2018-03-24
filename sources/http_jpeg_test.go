package sources

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/kithix/videre/test_helpers"
)

func TestHTTPReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(SmallestJPEG)
	}))
	r, err := HTTPBodyReader(HttpRequestDetails{
		URL: ts.URL,
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	buf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buf, r)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !reflect.DeepEqual(buf.Bytes(), SmallestJPEG) {
		t.Error("Buffer did not contain data, instead contained " + string(buf.Bytes()))
	}
	_, err = io.Copy(buf, r)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !reflect.DeepEqual(buf.Bytes(), append(SmallestJPEG, SmallestJPEG...)) {
		t.Error("Buffer is not equal to two small jpegs")
	}
	err = r.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	_, err = r.Read([]byte{})
	if err != io.EOF {
		t.Error(err)
		t.Log("Did not receive expected EOF")
		t.Fail()
	}
}
