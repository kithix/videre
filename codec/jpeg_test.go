package codec

import (
	"bytes"
	"image/jpeg"
	"testing"
)

var subject = NewJPEGCodec(jpeg.Options{})

func TestJPEGEncode(t *testing.T) {

}

func TestJPEGDecode(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	subject.Decode(buf)
}
