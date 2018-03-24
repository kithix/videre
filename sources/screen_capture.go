package sources

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/kbinani/screenshot"
)

func takeScreenshot() (image.Image, error) {
	return screenshot.CaptureDisplay(0)
}

func NewScreenshotAsJPEG() (io.Reader, error) {
	image, err := takeScreenshot()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

type Screenshotter struct {
}

func (s Screenshotter) Fetch() (io.Reader, error) {
	return NewScreenshotAsJPEG()
}

func NewScreenshotter() *Screenshotter {
	return &Screenshotter{}
}
