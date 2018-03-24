package sources

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	FetchTimeout = time.Second * 5
	TimeBetween  = time.Millisecond * 500
)

type _HTTPSource struct {
	request *http.Request
	client  *http.Client

	lastBody io.ReadCloser
	closed   error
}

func (s *_HTTPSource) Close() error {
	// Ensure our closed state is always EOF.
	defer func() { s.closed = io.EOF }()
	if s.lastBody != nil {
		// If we have a last body
		defer func() { s.lastBody = nil }()
		return s.lastBody.Close()
	}
	return s.closed
}

func (s *_HTTPSource) Read(p []byte) (int, error) {
	if s.closed != nil {
		return 0, s.closed
	}
	// If we have a last body to read from.
	if s.lastBody != nil {
		i, err := s.lastBody.Read(p)
		if err != nil {
			s.lastBody.Close()
			s.lastBody = nil
			return i, err
		}
		return i, err
	}

	log.Println("Doing request")
	response, err := s.client.Do(s.request)
	if err != nil {
		return 0, err
	}
	log.Println("Retrieved response")
	// Check statuscode
	s.lastBody = response.Body

	return s.lastBody.Read(p)
}

func (s *_HTTPSource) Fetch() (io.Reader, error) {
	log.Println("Doing request")
	response, err := s.client.Do(s.request)
	if err != nil {
		return nil, err
	}
	log.Println("Retrieved response")
	// Check statuscode
	if response.StatusCode > 400 {
		return nil, errors.New("Failed to retrieve your mother")
	}

	return response.Body, nil
}

func HTTPBodySource(requestDetails HttpRequestDetails) (*_HTTPSource, error) {
	request, client, err := buildRequest(requestDetails)
	if err != nil {
		return nil, err
	}

	return &_HTTPSource{
		request: request,
		client:  client,
	}, nil
}

func HTTPBodyReader(requestDetails HttpRequestDetails) (io.ReadCloser, error) {
	request, client, err := buildRequest(requestDetails)
	if err != nil {
		return nil, err
	}
	client.Timeout = FetchTimeout

	return &_HTTPSource{
		request: request,
		client:  client,

		lastBody: nil,
	}, nil
}
