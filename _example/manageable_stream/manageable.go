package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/kithix/videre/codec"
	"github.com/kithix/videre/servers"
	"github.com/kithix/videre/sources"
	"github.com/kithix/videre/stores"

	"github.com/kithix/stoppable"
)

func main() {
	source := sources.NewScreenshotter()
	// source := sources.NewWhiteNoiseGenerator(1024, 760)
	//source, err := sources.HTTPBodySource(sources.HttpRequestDetails{
	//	URL: "http://thecatapi.com/api/images/get?format=src&type=jpg",
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	store := stores.NewSingleStore()

	// jpegCodec := codec.NewJPEGCodec(jpeg.Options{
	// 	Quality: 100,
	// })
	codeca := codec.NewPassthrough()

	buf := bytes.NewBuffer([]byte{})
	// Our lifecycle as a startable and stoppable thing
	watchdog := stoppable.NewWatchdog(
		func() error {
			log.Println("Starting")
			return nil
		},
		func() error {
			reader, err := source.Fetch()
			if err != nil {
				return err
			}
			img, err := codeca.Decode(reader)
			if err != nil {
				return err
			}
			err = codeca.Encode(buf, img)
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
		},
		func(watcherErr stoppable.WatcherError) bool {
			log.Println(watcherErr)
			if watcherErr.Teardown == nil && watcherErr.Setup == nil {
				log.Println("Restarting in 1 second")
				time.Sleep(1 * time.Second)
				return true
			}
			log.Println("Stopping watchdog")
			return false
		},
	)

	watchdog.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<ul>
			<li><a href=/start> START </a></li>
			<li><a href=/stop> STOP </a></li>
			<li><a href=/mjpeg> MJPEG </a></li>
			<li><a href=/jpeg> JPEG </a></li>
			</ul>
			</html>
			`))
	})
	http.Handle("/mjpeg", servers.MJPEGHTTPHandler(store, codec.NewPassthrough()))
	http.Handle("/jpeg", servers.JPEGHTTPHandler(store, codec.NewPassthrough()))
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		err := watchdog.Start()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Started!"))
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		err := watchdog.Stop()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Stopped"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
