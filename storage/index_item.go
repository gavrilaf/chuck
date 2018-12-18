package storage

import (
	"fmt"
	"strconv"
	"strings"
)

type IndexItem struct {
	Focused bool
	Method  string
	Url     string
	Code    int
	Folder  string
}

func (i IndexItem) Format() string {
	return FormatIndexItem(i.Method, i.Url, i.Code, i.Folder, i.Focused)
}

func FormatIndexItem(method string, url string, code int, folder string, focused bool) string {
	prefix := "N"
	if focused {
		prefix = "F"
	}

	return fmt.Sprintf("%s,\t%d,\t%s,\t%s,\t%s", prefix, code, folder, method, url)
}

func ParseIndexItem(s string) *IndexItem { // TODO: add error handling
	fields := strings.Split(s, ",\t")
	if len(fields) != 5 {
		return nil
	}

	code, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil
	}

	focused := false
	if fields[0] == "F" {
		focused = true
	}

	return &IndexItem{
		Focused: focused,
		Method:  fields[3],
		Url:     fields[4],
		Code:    code,
		Folder:  fields[2],
	}
}
