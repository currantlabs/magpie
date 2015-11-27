package magpie

import (
	"bytes"
	"io"
	"os"
)

type fileInfo struct {
	*Asset
	*bytes.Reader
}

func (f *fileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, io.EOF
}

func (f *fileInfo) Stat() (os.FileInfo, error) {
	return f.Asset, nil
}

func (f *fileInfo) Close() error {
	return nil
}

func newFile(asset *Asset) *fileInfo {
	return &fileInfo{
		Asset:  asset,
		Reader: bytes.NewReader(asset.Content),
	}
}
