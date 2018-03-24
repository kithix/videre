package servers

import (
	"log"
	"net/http"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/sources"
)

func JPEGHTTPHandler(fetcher sources.Fetcher, codec codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reader, err := fetcher.Fetch()
		if err != nil {
			log.Println(err)
			return
		}
		img, err := codec.Decode(reader)
		if err != nil {
			log.Println(err)
			return
		}

		err = codec.Encode(w, img)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
