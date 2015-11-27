package magpie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type configJSON struct {
	Tags    []string `json:"tags"`
	Ignore  []string `json:"ignore"`
	Runtime *struct {
		Unsafe   *bool `json:"unsafe"`
		Compress *bool `json:"compress"`
	} `json:"runtime"`
	Attributes *struct {
		Mode    *uint  `json:"mode"`
		ModTime *int64 `json:"modtime"`
	} `json:"attributes"`
	Prefix     *string `json:"prefix"`
	Output     string  `json:"output"`
	Package    string  `json:"package"`
	FileSystem *bool   `json:"fileSystem"`
}

type assetJSON struct {
	configJSON
	Paths []string `json:"paths`
}

type rootConfigJSON struct {
	configJSON
	Assets []assetJSON `json:assets`
}

func getJSONConfig(configPath string) (*rootConfigJSON, error) {
	if configPath == "" {
		// If there's no config file specified, look for a local magpie.json
		root, err := getJSONConfig("magpie.json")
		if err == os.ErrNotExist {
			return nil, nil
		}
		return root, err
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil, err
	}

	var root rootConfigJSON
	err = json.Unmarshal(file, &root)
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return nil, err
	}
	return &root, nil
}

func getJSONOptions(json configJSON) []Option {
	var options []Option
	if json.Tags != nil {
		options = append(options, Tags(json.Tags))
	}
	if json.Ignore != nil {
		options = append(options, Ignore(json.Ignore))
	}
	if json.Runtime != nil {
		if json.Runtime.Compress != nil {
			options = append(options, Compress(*json.Runtime.Compress))
		}
		if json.Runtime.Unsafe != nil {
			options = append(options, Unsafe(*json.Runtime.Unsafe))
		}
	}
	if json.Attributes != nil {
		if json.Attributes.Mode != nil {
			options = append(options, OverrideFileMode(*json.Attributes.Mode))
		}
		if json.Attributes.ModTime != nil {
			options = append(options, OverrideFileModTime(*json.Attributes.ModTime))
		}
	}
	if json.Package != "" {
		options = append(options, PackageName(json.Package))
	}
	if json.Prefix != nil {
		options = append(options, Prefix(*json.Prefix))
	}
	if json.Output != "" {
		options = append(options, OutputFile(json.Output))
	}
	if json.FileSystem != nil {
		options = append(options, CreateFileSystem(*json.FileSystem))
	}
	return options
}
