package magpie

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
)

func NewFileSystem(nest Nest, prefix string) http.FileSystem {
	return &fileSystem{
		prefix: prefix,
		nest:   nest,
	}
}

type fileSystem struct {
	prefix string
	nest   Nest
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
	a, err := fs.nest(name)
	if err != nil {
		return nil, err
	}
	return newFile(a), nil
}
