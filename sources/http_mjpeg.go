package sources

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

type _HTTPMJPEGSource struct {
	response        *http.Response
	multiPartReader *multipart.Reader
}

func (s *_HTTPMJPEGSource) Read(p []byte) (int, error) {
	part, err := s.multiPartReader.NextPart()
	if err != nil {
		return 0, err
	}
	defer part.Close()

	return part.Read(p)
}

func (s *_HTTPMJPEGSource) Close() error {
	return s.response.Body.Close()
}

func HTTPMJPEGReader(requestDetails HttpRequestDetails, boundary string) (io.ReadCloser, error) {
	request, client, err := buildRequest(requestDetails)
	if err != nil {
		return nil, err
	}
	response, err := getResponse(request, client)
	if err != nil {
		return nil, err
	}
	if boundary == "" {
		boundary, err = getBoundary(response)
		if err != nil {
			return nil, err
		}
	}

	return &_HTTPMJPEGSource{
		response:        response,
		multiPartReader: multipart.NewReader(response.Body, boundary),
	}, nil
}

func getResponse(request *http.Request, client *http.Client) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		response.Body.Close()
		errs := "Got invalid response status: " + response.Status
		return nil, errors.New(errs)
	}

	return response, nil
}

// Looks through HTTP response headers to get a boundary string for reading the multipart data
func getBoundary(response *http.Response) (string, error) {
	header := response.Header.Get("Content-Type")
	if header == "" {
		return "", errors.New("Content-Type isn't specified!")
	}
	contentType, params, err := mime.ParseMediaType(header)
	if err != nil {
		return "", err
	}
	if contentType != "multipart/x-mixed-replace" {
		errs := "Wrong Content-Type: expected multipart/x-mixed-replace, got " + contentType
		return "", errors.New(errs)
	}
	boundary, ok := params["boundary"]
	if !ok {
		return "", errors.New("No multipart boundary param in Content-Type!")
	}
	// TODO fact check this
	// A boundary string with "--" is a terminator. Maybe this case was intended.
	// Some IP-cameras screw up boundary strings so we
	// have to remove excessive "--" characters manually.
	boundary = strings.Replace(boundary, "--", "", -1)
	return boundary, nil
}
