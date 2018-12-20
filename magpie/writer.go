package main

import (
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const hex = "0123456789abcdef"

type assetWriter interface {
	io.ByteWriter
	io.Writer
}

type writer struct {
	assetWriter
}

func (w *writer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		_, err = w.assetWriter.Write([]byte("\\x"))
		if err != nil {
			return
		}
		err = w.assetWriter.WriteByte(hex[b/16])
		if err != nil {
			return
		}
		err = w.assetWriter.WriteByte(hex[b%16])
		if err != nil {
			return
		}
		n++
	}
	return
}

func (w *writer) Close() error {
	return nil
}

func writeNest(fc *fileCollection, c *config, whitelistPath string) (err error) {
	err = os.MkdirAll(filepath.Dir(c.output), os.ModePerm)
	if err != nil {
		return
	}

	fd, err := os.Create(c.output)
	if err != nil {
		return
	}

	defer fd.Close()

	w := bufio.NewWriter(fd)
	defer w.Flush()

	err = writeHeader(w, fc, c)
	if err != nil {
		return
	}

	err = writeInit(w, fc, c)
	if err != nil {
		return
	}

	var assetKeys []string
	for assetKey := range fc.assets {
		assetKeys = append(assetKeys, assetKey)
	}
	sort.Strings(assetKeys)

	for _, assetKey := range assetKeys {
		asset := fc.assets[assetKey]
		if whitelistPath != "" {
			//writeLog("whitelistPath %s asset.path %s", whitelistPath, asset.path)
			if strings.HasSuffix(asset.path, whitelistPath) {
				writeLog("writing %s to %s", asset.path, "/assets_" + asset.constant + ".go")
				writeAssetFile(asset, c)
			}
		} else {
			writeAssetFile(asset, c)
		}
	}

	if c.fileSystem {
		err = writeFileSystem(w, c)
		if err != nil {
			return
		}
	}

	return nil
}

func writeHeader(w assetWriter, fc *fileCollection, c *config) (err error) {
	if _, err = fmt.Fprint(w, "// Code generated by magpie\n"); err != nil {
		return
	}
	if _, err = fmt.Fprint(w, "// sources:\n"); err != nil {
		return
	}

	writePackage(w, fc, c)

	imports := []string{"os", "time"}

	if c.fileSystem {
		imports = append(imports, "net/http")
	}
	if c.unsafe {
		imports = append(imports, "reflect", "unsafe")
	}

	// Write imports
	_, err = fmt.Fprint(w, "import (\n")
	if err != nil {
		return
	}

	sort.Strings(imports)

	for _, imp := range imports {
		_, err = fmt.Fprintf(w, "\t%q\n", imp)
		if err != nil {
			return
		}

	}

	_, err = fmt.Fprintf(w, `
	"github.com/currantlabs/magpie"
)

var %s magpie.Nest

`, c.nest)
	if err != nil {
		return
	}
	return
}

func writePackage(w assetWriter, fc *fileCollection, c *config) (err error) {
	if _, err = fmt.Fprint(w, "// Code generated by magpie\n"); err != nil {
		return
	}
	if _, err = fmt.Fprint(w, "// sources:\n"); err != nil {
		return
	}

	if fc != nil {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return
		}

		var paths []string
		for _, asset := range fc.assets {
			relative, _ := filepath.Rel(wd, asset.path)
			paths = append(paths, relative)
		}
		sort.Strings(paths)

		for _, path := range paths {
			//println("WD", wd)
			//println("path", asset.path)
			if _, err = fmt.Fprintf(w, "// %s\n", path); err != nil {
				return err
			}
		}
	}

	if _, err = fmt.Fprint(w, "// DO NOT EDIT!\n\n"); err != nil {
		return
	}
	// Write build tags, if applicable.
	if len(c.tags) > 0 {
		if _, err = fmt.Fprintf(w, "// +build %s\n\n", strings.Join(c.tags, " ")); err != nil {
			return
		}
	}

	packageName := c.packageName
	if packageName == "" {
		outpath, _ := filepath.Abs(c.output)
		outpath = filepath.Dir(outpath)
		path := filepath.Base(outpath)
		if path != "" {
			packageName = path
		}
	}
	if packageName == "" {
		packageName = "magpie"
	}

	// Write package declaration.
	_, err = fmt.Fprintf(w, "package %s\n\n", packageName)
	return
}

func writeInit(w assetWriter, fc *fileCollection, c *config) (err error) {
	_, err = fmt.Fprintf(w, "func init() {\n")
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "\n")
	if err != nil {
		return
	}

	writeRead(w, c)

	err = insertAssets(w, fc, c)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "\n\t%s = magpie.NewNest(assets)\n", c.nest)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "}\n")
	if err != nil {
		return
	}
	return
}

func writeRead(w io.Writer, c *config) (err error) {
	_, err = fmt.Fprintf(w, `	read := func(s string) []byte {
`)
	if err != nil {
		return
	}
	if c.unsafe {
		_, err = fmt.Fprintf(w, `		var empty [0]byte
		sx := (*reflect.StringHeader)(unsafe.Pointer(&s))
		b := empty[:]
		bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		bx.Data = sx.Data
		bx.Len = len(s)
		bx.Cap = bx.Len
`)
		if err != nil {
			return
		}
	} else {
		_, err = fmt.Fprintf(w, `		b := []byte(s)
`)
	}
	if c.compress {
		_, err = fmt.Fprint(w, `		b, err := magpie.Decompress(b)
		if err != nil {
			panic("Error decompressing asset" + err.Error())
		}
`)
	}
	_, err = fmt.Fprint(w, `		return b
	}
`)
	return
}

func insertAssets(w io.Writer, fc *fileCollection, c *config) (err error) {
	// NewAsset(name string, content []byte, size int64, mode os.FileMode, modTime time.Time) *Asset

	_, err = fmt.Fprint(w, "\n\tvar assets = map[string]magpie.Asset{\n")
	if err != nil {
		return
	}
	for _, a := range fc.assets {
		l := base64.StdEncoding.EncodeToString(a.hash)
		_, err = fmt.Fprintf(w, "\t\t%q: magpie.NewAsset(%q, read(_%s), %d, os.FileMode(%d), time.Unix(%d, 0), %q),\n", a.name, a.name, a.constant, a.info.Size(), a.info.Mode(), a.info.ModTime().Unix(), l)
		if err != nil {
			return
		}

	}
	_, err = fmt.Fprint(w, "\n\t}\n")
	return
}

func writeAssetFile(a *asset, c *config) (err error) {
	dir := filepath.Dir(c.output)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}

	fd, err := os.Create(dir + "/assets_" + a.constant + ".go")
	if err != nil {
		return
	}

	defer fd.Close()

	w := bufio.NewWriter(fd)
	defer w.Flush()

	err = writePackage(w, nil, c)
	if err != nil {
		return err
	}

	return writeAsset(w, a, c)
}

func writeAsset(w assetWriter, a *asset, c *config) error {
	fd, err := os.Open(a.path)
	if err != nil {
		return err
	}

	defer fd.Close()
	// TODO: Compress

	_, err = fmt.Fprintf(w, "const _%s = \"", a.constant)
	if err != nil {
		return err
	}

	var wr io.WriteCloser
	hr := newHashWriter(&writer{assetWriter: w})
	if c.compress {
		wr = gzip.NewWriter(hr)
	} else {
		wr = hr
	}
	_, err = io.Copy(wr, fd)
	if c.compress {
		wr.Close()
	}
	a.hash = hr.Sum(nil)
	_, err = fmt.Fprintf(w, "\"\n")
	if err != nil {
		return err
	}
	return nil
}

func writeFileSystem(w io.Writer, c *config) (err error) {
	fmt.Fprintf(w, `
func FileSystem(prefix string) http.FileSystem {
	return magpie.NewFileSystem(%s, prefix)
}
`, c.nest)
	return
}