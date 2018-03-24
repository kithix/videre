package sources

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"math/rand"
	"time"
)

func makeWhiteNoise(x, y int) image.Image {
	rect := image.Rect(0, 0, x, y)

	rand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	image := image.NewRGBA(rect.Bounds())
	for y := 0; y < rect.Bounds().Max.Y-1; y++ {
		for x := 0; x < rect.Bounds().Max.X-1; x++ {
			alpha := uint8(rand.Float32() * 255)
			rgbc := uint8(rand.Float32() * 255)
			rgbc = rgbc * alpha
			image.SetRGBA(x, y, color.RGBA{rgbc, rgbc, rgbc, alpha})
		}
	}
	return image
}

func NewWhiteNoiseAsJPEG(x, y int) (io.Reader, error) {
	image := makeWhiteNoise(x, y)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

type WhiteNoiseGenerator struct {
	X int
	Y int
}

func (s WhiteNoiseGenerator) Fetch() (io.Reader, error) {
	return NewWhiteNoiseAsJPEG(s.X, s.Y)
}

func NewWhiteNoiseGenerator(x, y int) *WhiteNoiseGenerator {
	return &WhiteNoiseGenerator{X: x, Y: y}
}
