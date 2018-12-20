package main

import (
	"fmt"
)

type Magpie struct {
	configs []*config
}

func New(configPath string, paths []string, options ...Option) (*Magpie, error) {
	configs, err := buildConfig(configPath, paths, options...)
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		fmt.Printf("Config: %+v\n", config)
	}
	return &Magpie{
		configs: configs,
	}, nil
}
