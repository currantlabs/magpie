package magpie

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type asset struct {
	path     string
	name     string
	constant string
	info     os.FileInfo
}

type fileCollection struct {
	visited   map[string]bool
	constants map[string]int
	assets    []asset
}

func findFiles(c *config) (*fileCollection, error) {
	a := &fileCollection{
		visited:   make(map[string]bool),
		constants: make(map[string]int),
	}
	for _, dir := range c.inputs {
		err := find(dir.path, dir.recursive, c.prefix, a, c)
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func find(dir string, recursive bool, prefix string, fc *fileCollection, c *config) error {
	path := dir
	if prefix != "" {
		path, _ = filepath.Abs(path)
		prefix, _ = filepath.Abs(prefix)
		prefix = filepath.ToSlash(prefix)
	}

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	var list []os.FileInfo

	if fi.IsDir() {
		fc.visited[path] = true
		d, err := os.Open(path)
		if err != nil {
			return err
		}
		defer d.Close()
		list, err = d.Readdir(0)
		if err != nil {
			return err
		}
		// Sort to make output stable between invocations
		sort.Sort(ByName(list))
	} else {
		path = filepath.Dir(path)
		list = []os.FileInfo{fi}
	}

	for _, file := range list {
		p := filepath.Join(path, file.Name())
		var ignored bool
		for _, pattern := range c.ignore {
			if pattern.MatchString(p) {
				ignored = true
				break
			}
		}
		if ignored {
			continue
		}

		if file.IsDir() {
			if recursive {
				fc.visited[p] = true
				find(filepath.Join(path, file.Name()), recursive, prefix, fc, c)
			}
			continue
		}
		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			var linkPath string
			if linkPath, err = os.Readlink(p); err != nil {
				return err
			}
			if !filepath.IsAbs(linkPath) {
				if linkPath, err = filepath.Abs(p + "/" + linkPath); err != nil {
					return err
				}
			}
			if _, ok := fc.visited[linkPath]; !ok {
				fc.visited[linkPath] = true
				find(p, recursive, prefix, fc, c)
			}
			continue
		}
		var a asset
		a.name = filepath.ToSlash(p)
		if strings.HasPrefix(a.name, prefix) {
			a.name = a.name[len(prefix):]
		}
		if len(a.name) > 0 && a.name[0] == '/' {
			a.name = a.name[1:]
		}
		if len(a.name) == 0 {
			return errors.New("Invalid file: " + p)
		}
		a.path, _ = filepath.Abs(p)
		a.constant = safeConstantName(a.name, fc)
		a.info = file
		fc.assets = append(fc.assets, a)
	}
	return nil
}

var regFuncName = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func safeConstantName(name string, fc *fileCollection) string {
	var inBytes, outBytes []byte
	var toUpper bool

	name = strings.ToLower(name)
	inBytes = []byte(name)

	for i := 0; i < len(inBytes); i++ {
		if regFuncName.Match([]byte{inBytes[i]}) {
			toUpper = true
		} else if toUpper {
			outBytes = append(outBytes, []byte(strings.ToUpper(string(inBytes[i])))...)
			toUpper = false
		} else {
			outBytes = append(outBytes, inBytes[i])
		}
	}

	name = string(outBytes)

	// Identifier can't start with a digit.
	if unicode.IsDigit(rune(name[0])) {
		name = "_" + name
	}

	if num, ok := fc.constants[name]; ok {
		fc.constants[name] = num + 1
		name = fmt.Sprintf("%s%d", name, num)
	} else {
		fc.constants[name] = 2
	}

	return name
}

type ByName []os.FileInfo

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }
