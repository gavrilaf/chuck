package storage

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
)

type seekerImpl struct {
	root  *afero.Afero
	index Index
	log   utils.Logger
}

func NewSeekerWithFs(folder string, fs afero.Fs, log utils.Logger) (Seeker, error) {
	folder = strings.Trim(folder, " \\/")
	logDirExists, _ := afero.DirExists(fs, folder)
	if !logDirExists {
		return nil, fmt.Errorf("Folder %s doesn't exists", folder)
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, folder)}

	index, err := LoadIndex(root, "index.txt", true)
	if err != nil {
		return nil, err
	}

	seeker := &seekerImpl{
		index: index,
		root:  root,
		log:   log,
	}

	return seeker, nil
}

func (seeker *seekerImpl) Look(method string, url string) *http.Response {
	item := seeker.index.Find(method, url, SEARCH_SUBSTR)
	if item == nil {
		return nil
	}

	header, err := seeker.readHeader(item.Folder + "/resp_header.json")
	if err != nil {
		seeker.log.Error("Read header error for %s: %v", item.Folder, err)
		return nil
	}

	body, err := seeker.readBody(item.Folder + "/resp_body.json")
	if err != nil {
		seeker.log.Error("Read header body for %s: %v", item.Folder, err)
		return nil
	}

	response := &http.Response{
		StatusCode: item.Code,
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
