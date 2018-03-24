package codec

import (
	"image"
	"image/jpeg"
	"io"
)

type JPEGCodec struct {
	jpegOptions jpeg.Options
}

func (c *JPEGCodec) Encode(writer io.Writer, img image.Image) error {
	return jpeg.Encode(writer, img, &c.jpegOptions)
}

func (c *JPEGCodec) Decode(reader io.Reader) (image.Image, error) {
	return jpeg.Decode(reader)
}

func NewJPEGCodec(opts jpeg.Options) *JPEGCodec {
	return &JPEGCodec{
		jpegOptions: opts,
	}
}
