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
	key        string
	folder     string
	statusCode int
}

type seekerImpl struct {
	root     *afero.Afero
	requests []respNode
	log      utils.Logger
}

func NewSeekerWithFs(folder string, fs afero.Fs, log utils.Logger) (Seeker, error) {
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

	requests := make([]respNode, 0)
	scanner := bufio.NewScanner(file)
	linesCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			log.Error("Invalid fields count in %d", linesCount)
			continue
		}
		if fields[0] == "F" { // focused
			statusCode, err := strconv.Atoi(fields[4])
			if err != nil {
				log.Error("Invalid status code in %d, %v", linesCount, err)
				continue
			}
			node := respNode{
				key:        createKey(fields[2], fields[3]),
				folder:     fields[1],
				statusCode: statusCode,
			}
			requests = append(requests, node)
		}
		linesCount += 1
	}

	log.Info("Loaded index in %s, lines %d, focused %d", folder, linesCount, len(requests))

	seeker := &seekerImpl{
		requests: requests,
		root:     root,
		log:      log,
	}

	return seeker, nil
}

func (seeker *seekerImpl) Look(method string, url string) *http.Response {
	key := createKey(method, url)
	var req respNode
	ok := false
	for _, node := range seeker.requests {
		if strings.HasPrefix(key, node.key) {
			req = node
			ok = true
			break
		}
	}

	if !ok {
		return nil
	}

	header, err := seeker.readHeader(req.folder + "/resp_header.json")
	if err != nil {
		seeker.log.Error("Read header error for %s: %v", key, err)
		return nil
	}

	body, err := seeker.readBody(req.folder + "/resp_body.json")
	if err != nil {
		seeker.log.Error("Read header body for %s: %v", key, err)
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
