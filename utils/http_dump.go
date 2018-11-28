package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func DumpRespBody(resp *http.Response) ([]byte, error) {
	var err error
	save := resp.Body
	savecl := resp.ContentLength

	if resp.Body == nil {
		resp.Body = emptyBody
	} else {
		save, resp.Body, err = drainBody(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)

	resp.Body = save
	resp.ContentLength = savecl

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

///////////////////////////////////////////////////////////////////////////////////
// Private

// emptyBody is an instance of empty reader.
var emptyBody = ioutil.NopCloser(strings.NewReader(""))

// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}

	if err = b.Close(); err != nil {
		return nil, b, err
	}

	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
