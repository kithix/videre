package main

import (
	"bytes"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/managers"
	"github.com/kithix/videre/servers"
	"github.com/kithix/videre/sources"
	"github.com/kithix/videre/stores"

	"github.com/kithix/stoppable"
)

func main() {
	// source := sources.NewWhiteNoiseGenerator(1024, 760)
	source, err := sources.HTTPBodySource(sources.HttpRequestDetails{
		URL: "http://thecatapi.com/api/images/get?format=src&type=jpg",
	})
	if err != nil {
		log.Fatal(err)
	}
	store := stores.NewSingleStore()

	jpegCodec := codec.NewJPEGCodec(jpeg.Options{
		Quality: 100,
	})

	buf := bytes.NewBuffer([]byte{})
	// Our lifecycle as a startable and stoppable thing
	proxy := managers.NewProxyManager(func() io.Closer {
		closer, err := stoppable.Open(
			func() error {
				log.Println("Starting")
				return nil
			},
			func() error {
				reader, err := source.Fetch()
				time.Sleep(1 * time.Second)
				if err != nil {
					return err
				}
				img, err := jpegCodec.Decode(reader)
				if err != nil {
					return err
				}
				err = jpegCodec.Encode(buf, img)
				if err != nil {
					return err
				}
				defer buf.Reset()
				_, err = store.Write(buf.Bytes())
				if err != nil {
					return err
				}
				return nil
			},
			func() error {
				log.Println("Stopping")
				return nil
			})
		if err != nil {
			return nil
		}
		return closer
	}, servers.MJPEGHTTPHandler(store, codec.NewPassthrough()))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<ul>
			<li><a href=/mjpeg> MJPEG </a></li>
			<li><a href=/jpeg> JPEG </a></li>
			</ul>
			</html>
			`))
	})
	http.Handle("/mjpeg", proxy)
	http.Handle("/jpeg", servers.JPEGHTTPHandler(store, codec.NewPassthrough()))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
