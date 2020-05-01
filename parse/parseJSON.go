package parse

import (
	"fmt"
	"os"
	"strings"
)

const Delimiter = "."

type ParseJSON struct {
	filePath string
	root *jsonValue
}

func NewParseJSON(filePath string) *ParseJSON {
	return &ParseJSON{
		filePath:filePath,
		root:new(jsonValue),
	}
}

func (p *ParseJSON)Load() error {
	if 0 == len(p.filePath) {
		return fmt.Errorf("file can not be null")
	}

	file, err := os.Open(p.filePath)
	if err != nil {
		return fmt.Errorf("open file(%s) failed: %s", p.filePath, err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	if _, err := file.Read(buf); err != nil {
		return fmt.Errorf("read file(%s) failed: %s", p.filePath, err)
	}

	return p.LoadFromString(buf)
}

func (p *ParseJSON)LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file(%s) failed: %s", filePath, err)
	}
	defer file.Close()

	buf := make([]byte, 0)
	if _, err := file.Read(buf); err != nil {
		return fmt.Errorf("read file(%s) failed: %s", filePath, err)
	}

	if err := p.LoadFromString(buf); err != nil {
		return err
	}
	return nil
}

func (p *ParseJSON)LoadFromString(buf []byte) error {
	var err error
	p.root, err = parse(buf)
	return err
}


func (p* ParseJSON)GetStringValue(path string) (string, error) {
	s := strings.Split(path, Delimiter)
	node, err := p.root.Get(s)
	if err != nil {
		return "", err
	}
	if node.Type != String {
		return "", fmt.Errorf("invalid type")
	}
	return string(*node.Value.(*jsonValueString)), nil
}

func (p *ParseJSON)GetIntValue(path string) (int, error) {
	s := strings.Split(path, Delimiter)
	node, err := p.root.Get(s)
	if err != nil {
		return 0, err
	}
	if node.Type != Int {
		return 0, fmt.Errorf("invalid type")
	}
	return int(*node.Value.(*jsonValueInt)), nil
}

