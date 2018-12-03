package storage

import (
	"bufio"
	"fmt"
	"github.com/spf13/afero"
	"net/http"
	"strconv"
	"strings"
)

type respNode struct {
	folder     string
	statusCode int
}

type reqSeeker struct {
	root     *afero.Afero
	requests map[string]respNode
}

func NewSeekerWithFs(folder string, fs afero.Fs) (ReqSeeker, error) {
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

			key := fields[2] + ":" + fields[3]
			// TODO: log
			requests[key] = node
		}
	}

	defer file.Close()

	seeker := &reqSeeker{
		requests: requests,
		root:     root,
	}

	return seeker, nil
}

func (seeker *reqSeeker) Look(method string, url string) *http.Response {
	return nil
}
