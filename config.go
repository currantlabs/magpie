package magpie

import (
	"path/filepath"
	"regexp"
	"strings"
)

type input struct {
	Path      string
	Recursive bool
}

type config struct {
	Inputs     []*input
	Package    string
	Tags       []string
	Output     string
	Prefix     string
	Unsafe     bool
	Compress   bool
	Mode       uint
	ModTime    int64
	Ignore     []*regexp.Regexp
	FileSystem bool
}

func newConfig() *config {
	return &config{
		Output:   "magpie.go",
		Compress: true,
	}
}

func buildConfig(configPath string, paths []string, options ...Option) ([]*config, error) {
	baseConfig := newConfig()
	var configs []*config
	rootJSON, err := getJSONConfig(configPath)
	if err != nil {
		return nil, err
	}
	if rootJSON != nil {
		mergeConfig(baseConfig, rootJSON.configJSON)
		for _, option := range options {
			option(baseConfig)
		}
		for _, asset := range rootJSON.Assets {
			pc := *baseConfig
			mergeConfig(&pc, asset.configJSON)
			pc.Inputs = make([]*input, len(asset.Paths))
			for i, path := range asset.Paths {
				pc.Inputs[i] = parsePath(path)
			}
			configs = append(configs, &pc)
		}
	} else {
		for _, option := range options {
			option(baseConfig)
		}
	}
	if len(paths) > 0 {
		pc := *baseConfig
		pc.Inputs = make([]*input, len(paths))
		for i, path := range paths {
			pc.Inputs[i] = parsePath(path)
		}
		configs = append(configs, &pc)
	}
	return configs, nil
}

func mergeConfig(c *config, json configJSON) {
	for _, option := range getJSONOptions(json) {
		option(c)
	}
}

func parsePath(path string) *input {
	if strings.HasSuffix(path, "/...") {
		return &input{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	} else {
		return &input{
			Path:      filepath.Clean(path),
			Recursive: false,
		}
	}

}
