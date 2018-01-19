package magpie

import (
	"path/filepath"
	"regexp"
	"strings"
)

type input struct {
	path      string
	recursive bool
}

type config struct {
	nest        string
	inputs      []*input
	packageName string
	tags        []string
	output      string
	prefix      string
	unsafe      bool
	compress    bool
	fileMode    uint
	fileModTime int64
	ignore      []*regexp.Regexp
	fileSystem  bool
	live        bool
}

func newConfig() *config {
	return &config{
		nest:     "Nest",
		output:   "magpie.go",
		compress: true,
	}
}

func buildConfig(configPath string, cliInputPaths []string, cliOptions ...Option) ([]*config, error) {
	baseConfig := newConfig()
	var configs []*config
	configJSON, err := getJSONConfig(configPath)
	if err != nil {
		return nil, err
	}
	if configJSON != nil {
		for _, option := range getCommonOptionsJSON(configJSON.commonOptionsJSON) {
			option(baseConfig)
		}
		for _, option := range cliOptions {
			option(baseConfig)
		}
		for _, asset := range configJSON.Assets {
			pc := *baseConfig
			for _, option := range getAssetJSONOptions(asset) {
				option(&pc)
			}
			pc.inputs = make([]*input, len(asset.Paths))
			for i, path := range asset.Paths {
				pc.inputs[i] = parseInputPath(path)
			}
			configs = append(configs, &pc)
		}
	} else {
		for _, option := range cliOptions {
			option(baseConfig)
		}
	}
	if len(cliInputPaths) > 0 {
		pc := *baseConfig
		pc.inputs = make([]*input, len(cliInputPaths))
		for i, path := range cliInputPaths {
			pc.inputs[i] = parseInputPath(path)
		}
		configs = append(configs, &pc)
	}
	return configs, nil
}

func parseInputPath(path string) *input {
	if strings.HasSuffix(path, "/...") {
		return &input{
			path:      filepath.Clean(path[:len(path)-4]),
			recursive: true,
		}
	} else {
		return &input{
			path:      filepath.Clean(path),
			recursive: false,
		}
	}
}
