package main

import (
	"crypto/sha1"
	"hash"
	"io"
)

type hashWriter struct {
	io.WriteCloser
	hash.Hash
}

func newHashWriter(w io.WriteCloser) *hashWriter {
	return &hashWriter{
		WriteCloser: w,
		Hash:        sha1.New(),
	}
}

func (hw *hashWriter) Write(p []byte) (n int, err error) {
	hw.Hash.Write(p)
	return hw.WriteCloser.Write(p)
}
