package main

import (
	"io"
	"log"
	"net/http"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/servers"
	"github.com/kithix/videre/sources"
	"github.com/kithix/videre/stores"

	_ "net/http/pprof"
)

type closeable chan struct{}

func (c closeable) Close() error {
	c <- struct{}{}
	return nil
}

func main() {
	source := sources.NewScreenshotter() //{X: 1920, Y: 1080}
	store := stores.NewSingleStore()
	server := servers.MJPEGHTTPHandler(store, codec.NewPassthrough())

	closerMaker := func() io.Closer {
		closer := make(closeable)
		go func() {
			log.Println("Starting")
			for {
				select {
				case <-closer:
					log.Println("Closing")
					return
				default:
					reader, err := source.Fetch()
					if err != nil {
						log.Fatal(err)
					}
					io.Copy(store, reader)
				}
			}
		}()
		return closer
	}
	closerMaker()

	// proxy := managers.NewProxyManager(closerMaker, server)

	http.Handle("/jpeg", servers.JPEGHTTPHandler(store, codec.NewPassthrough()))
	http.Handle("/", server)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
