package storage

import (
	"bufio"
	"fmt"
	"github.com/spf13/afero"
	"strings"
)

type SearchType int

const (
	SEARCH_EQ SearchType = iota + 1
	SEARCH_SUBSTR
)

type Index interface {
	Add(item IndexItem)
	Find(method string, url string, searchType SearchType) *IndexItem
	Size() int
	Get(index int) IndexItem
}

func NewIndex() Index {
	return &indexImpl{
		nodes: make([]indexNode, 0),
	}
}

func LoadIndex(fs afero.Fs, file string) (Index, error) {
	fp, err := fs.Open(file)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

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
		index.Add(*item)
	}

	return index, nil
}

////////////////////////////////////////////////////////////////////////////////////////

type indexNode struct {
	key  string
	item IndexItem
}

type indexImpl struct {
	nodes []indexNode
}

func (p *indexImpl) Add(item IndexItem) {
	key := item.Method + ":" + item.Url
	p.nodes = append(p.nodes, indexNode{
		key:  key,
		item: item,
	})
}

func (p *indexImpl) Find(method string, url string, searchType SearchType) *IndexItem {
	var found indexNode

	key := method + ":" + url
	ok := false

	// should be replaced by more effective implementation
	for _, node := range p.nodes {
		if (searchType == SEARCH_EQ && key == node.key) || (searchType == SEARCH_SUBSTR && strings.HasPrefix(node.key, key)) {
			found = node
			ok = true
			break
		}
	}

	if !ok {
		return nil
	}

	return &found.item
}

func (p *indexImpl) Size() int {
	return len(p.nodes)
}

func (p *indexImpl) Get(index int) IndexItem {
	return p.nodes[index].item
}
