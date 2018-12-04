package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
)

type respNode struct {
	folder     string
	statusCode int
}

type seekerImpl struct {
	root     *afero.Afero
	requests map[string]respNode
}

func NewSeekerWithFs(folder string, fs afero.Fs) (Seeker, error) {
	folder = strings.Trim(folder, " \\/")
	logDirExists, _ := afero.DirExists(fs, folder)
	if !logDirExists {
		return nil, fmt.Errorf("Folder %s doesn't exists", folder)
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, folder)}

	file, err := root.Open("index.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	requests := make(map[string]respNode)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			// TODO: Error log
			continue
		}
		if fields[0] == "F" { // focused
			statusCode, err := strconv.Atoi(fields[4])
			if err != nil {
				// TODO: Error log
				continue
			}
			node := respNode{
				folder:     fields[1],
				statusCode: statusCode,
			}

			key := createKey(fields[2], fields[3])
			// TODO: log
			requests[key] = node
		}
	}

	seeker := &seekerImpl{
		requests: requests,
		root:     root,
	}

	return seeker, nil
}

func (seeker *seekerImpl) Look(method string, url string) *http.Response {
	key := createKey(method, url)
	req, ok := seeker.requests[key]
	if !ok {
		return nil
	}

	// read headers
	header, err := seeker.readHeader(req.folder + "/resp_header.json")
	if err != nil {
		fmt.Printf("Read header error for %s: %v", key, err)
		// TODO: Log here
		return nil
	}

	body, err := seeker.readBody(req.folder + "/resp_body.json")
	if err != nil {
		fmt.Printf("Read body error for %s: %v", key, err)
		// TODO: Log here
		return nil
	}

	response := &http.Response{
		StatusCode: req.statusCode,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
		Body:       body,
	}

	return response
}

/*
 * Private
 */

func (seeker *seekerImpl) readHeader(fname string) (http.Header, error) {
	exists, err := seeker.root.Exists(fname)
	if err != nil {
		return nil, err
	}

	if !exists {
		return make(http.Header), nil
	}

	fp, err := seeker.root.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	header, err := utils.DecodeHeaders(buf)
	return header, err
}

func (seeker *seekerImpl) readBody(fname string) (io.ReadCloser, error) {
	exists, err := seeker.root.Exists(fname)
	if err != nil {
		return nil, err
	}

	if !exists {
		return ioutil.NopCloser(strings.NewReader("")), nil
	}

	fp, err := seeker.root.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(buf)), nil
}

func createKey(method string, url string) string {
	return method + ":" + url
}
