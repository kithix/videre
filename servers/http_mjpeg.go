package servers

import (
	"bytes"
	"errors"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/sources"
)

func MJPEGHTTPHandler(fetcher sources.Fetcher, codec codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mpw := multipart.NewWriter(w)
		w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+mpw.Boundary())
		w.Header().Add("Cache-Control", "no-cache")
		w.WriteHeader(200)

		buf := bytes.NewBuffer([]byte{})
		var err error
		var reader io.Reader
		var part io.Writer
		var img image.Image
		var header textproto.MIMEHeader
		for {
			reader, err = fetcher.Fetch()
			if err != nil {
				log.Println(err)
				break
			}
			img, err = codec.Decode(reader)
			if err != nil {
				log.Println(err)
				return
			}
			err = codec.Encode(buf, img)
			if err != nil {
				return
			}
			if err != nil {
				// check EOF
				break
			}
			if buf.Len() == 0 {
				err = errors.New("Nothing to write")
				break
			}
			header = textproto.MIMEHeader{}
			header.Add("Content-Type", "image/jpeg")
			header.Add("Content-Length", strconv.Itoa(buf.Len()))
			header.Add("X-Timestamp", "0.000000")
			part, err = mpw.CreatePart(header)
			if err != nil {
				break
			}
			_, err = part.Write(buf.Bytes())
			if err != nil {
				break
			}
			buf.Reset()
		}
		// Error handler
		if err != nil {
			log.Println(err)
		}
	}
}
