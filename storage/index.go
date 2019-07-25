package storage

import (
	"bufio"
	"fmt"
	"github.com/gavrilaf/grouter"
	"github.com/spf13/afero"

	"chuck/utils"
)

type Index interface {
	Add(item IndexItem) error
	Find(method string, url string) *IndexItem
	Size() int
	Get(index int) IndexItem
}

func NewIndex() Index {
	return &indexImpl{
		items:  make([]IndexItem, 0),
		router: grouter.NewRouter(),
	}
}

// Index creation

func LoadIndex(fp afero.File, focused bool, log utils.Logger) (Index, error) {
	index := NewIndex()

	scanner := bufio.NewScanner(fp)
	lineIndex := 0
	for scanner.Scan() {
		line := scanner.Text()
		item := ParseIndexItem(line)
		if item == nil {
			return nil, fmt.Errorf("Couldn't parse line %d", lineIndex)
		}
		lineIndex += 1
		if !focused || item.Focused {
			err := index.Add(*item)
			if err != nil {
				if err == grouter.ErrAlreadyAdded {
					if log != nil {
						log.Error("File %s, line %d. Route already added", fp.Name(), lineIndex)
					}
					continue
				}
				return nil, err
			}
		}
	}

	return index, nil
}

func LoadIndex2(fs afero.Fs, file string, focused bool, log utils.Logger) (Index, error) {
	fp, err := fs.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return LoadIndex(fp, focused, log)
}

// indexImpl

type indexImpl struct {
	items  []IndexItem
	router grouter.Router
}

func (self *indexImpl) Add(item IndexItem) error {
	self.items = append(self.items, item)

	err := self.router.AddRoute(item.Method, item.Url, len(self.items)-1)
	return err
}

func (self *indexImpl) Find(method string, url string) *IndexItem {
	p, _ := self.router.Lookup(method, url)
	if p == nil { // TODO: add errors handling
		return nil
	}

	index := p.Value.(int)
	if index >= 0 && index < len(self.items) {
		return &self.items[index]
	}

	return nil
}

func (self *indexImpl) Size() int {
	return len(self.items)
}

func (self *indexImpl) Get(index int) IndexItem {
	return self.items[index]
}
