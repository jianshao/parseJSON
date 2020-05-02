package parse

import (
	"fmt"
	"os"
	"strings"
)

const Delimiter = "."

type JSON struct {
	filePath string
	root Value
}

func NewParseJSON(filePath string) *JSON {
	return &JSON{
		filePath:filePath,
		root:&jsonValueObject{},
	}
}

func (p *JSON)Load() error {
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

func (p *JSON)LoadFromFile(filePath string) error {
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

func (p *JSON)LoadFromString(buf []byte) error {
	if p.root == nil {
		p.root = &jsonValueObject{}
	}
	if _, err := p.root.parse(buf, 0); err != nil {
		return err
	}
	return nil
}

func (p *JSON)getNode(path string) (Value, error) {
	var err error
	s := strings.Split(path, Delimiter)
	node := p.root
	for i := 0; i < len(s); i++ {
		if node, err = node.get(s[i]); err != nil {
			return nil, err
		}
	}

	return node, nil
}

func (p *JSON)GetStringValue(path string) (string, error) {
	if node, err := p.getNode(path); err != nil {
		return "", err
	} else {
		if node.getType() == String {
			return string(node.(*jsonValueString).v), nil
		} else {
			return "", fmt.Errorf("invalid type")
		}
	}
}

func (p *JSON)GetIntValue(path string) (int, error) {
	if node, err := p.getNode(path); err == nil {
		if node.getType() == Int {
			return int(node.(*jsonValueInt).v), nil
		} else {
			return -1, fmt.Errorf("invalid type")
		}
	} else {
		return -1, err
	}
}

