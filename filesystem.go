package magpie

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func NewFileSystem(assets map[string]*Asset, prefix string) http.FileSystem {
	return &fileSystem{
		prefix: prefix,
		assets: assets,
	}
}

type fileSystem struct {
	prefix string
	assets map[string]*Asset
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 || strings.Contains(name, "\x00") {
		return nil, errors.New("magpie: invalid character in file path")
	}
	if fs.prefix != "" && strings.HasPrefix(name, fs.prefix) {
		name = name[len(fs.prefix):]
	}
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	if asset, ok := fs.assets[name]; ok {
		return newFile(asset), nil
	}
	return nil, fmt.Errorf("Unknown asset: %s", name)
}
