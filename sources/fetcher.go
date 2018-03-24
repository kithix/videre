package sources

import "io"

type Fetcher interface {
	Fetch() (io.Reader, error)
}

type FetcherFunc func() (io.Reader, error)

func (f FetcherFunc) Fetch() (io.Reader, error) {
	return f()
}

func MakeFetcher(f func() (io.Reader, error)) FetcherFunc {
	return FetcherFunc(f)
}
