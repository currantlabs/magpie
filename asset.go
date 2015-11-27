package magpie

import (
	"os"
	"time"
)

type Asset struct {
	Content []byte
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func NewAsset(name string, content []byte, size int64, mode os.FileMode, modTime time.Time) *Asset {
	return &Asset{
		name:    name,
		Content: content,
		size:    size,
		mode:    mode,
		modTime: modTime,
	}
}

func (a Asset) Name() string {
	return a.name
}

func (a Asset) Size() int64 {
	return a.size
}

func (a Asset) Mode() os.FileMode {
	return a.mode
}

func (a Asset) ModTime() time.Time {
	return a.modTime
}

func (a Asset) IsDir() bool {
	return false
}

func (a Asset) Sys() interface{} {
	return nil
}
