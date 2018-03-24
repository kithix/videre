package codec

import (
	"image"
	"io"
	"sync"
)

type Codec interface {
	Encode(io.Writer, image.Image) error
	Decode(io.Reader) (image.Image, error)
}

// Passthrough should be used as a codec when both sides are the same.
type Passthrough struct {
	reader io.Reader
	lock   sync.Mutex
}

func (c *Passthrough) Encode(writer io.Writer, img image.Image) error {
	c.lock.Lock()
	_, err := io.Copy(writer, c.reader)
	c.lock.Unlock()
	return err
}

func (c *Passthrough) Decode(reader io.Reader) (image.Image, error) {
	c.lock.Lock()
	c.reader = reader
	c.lock.Unlock()
	return nil, nil
}

func NewPassthrough() *Passthrough {
	return &Passthrough{}
}
