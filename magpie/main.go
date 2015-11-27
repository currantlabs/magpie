package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/currantlabs/magpie"
)

const configFlag = "config"

const tagFlag = "tags"
const prefixFlag = "prefix"
const packageNameFlag = "package"
const unsafeFlag = "unsafe"
const compressFlag = "compress"
const fileModeFlag = "mode"
const fileModTimeFlag = "modtime"
const outputFlag = "output"
const ignoreFlag = "ignore"
const watchFlag = "watch"

func main() {

	app := cli.NewApp()
	app.Name = "magpie"
	app.Usage = "Bundle files into a go binary"
	app.Action = action

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  configFlag,
			Usage: "The configuration file for magpie",
		},
		cli.StringSliceFlag{
			Name:  tagFlag,
			Usage: "Optional build tag to include",
		},
		cli.StringFlag{
			Name:  prefixFlag,
			Usage: "Optional path prefix to strip off asset names",
		},
		cli.StringFlag{
			Name:  packageNameFlag,
			Usage: "Package name to use in the generated code",
		},
		cli.BoolFlag{
			Name:  unsafeFlag,
			Usage: "Use a .rodata hack to get rid of unnecessary memcopies. This requires the ability to run unsafe code",
		},
		cli.BoolTFlag{
			Name:  compressFlag,
			Usage: "Compress assets with gzip",
		},
		cli.StringFlag{
			Name:  fileModeFlag,
			Usage: "Optional file mode override for all files",
		},
		cli.StringFlag{
			Name:  fileModTimeFlag,
			Usage: "Optional modification unix timestamp override for all files",
		},
		cli.StringFlag{
			Name:  outputFlag,
			Value: "magpie.go",
			Usage: "Optional name of the output file to be generated",
		},
		cli.StringSliceFlag{
			Name:  ignoreFlag,
			Usage: "Regex pattern to ignore",
		},
		cli.BoolFlag{
			Name:  watchFlag,
			Usage: "Watch filesystem for changes",
		},
	}

	app.Run(os.Args)
}

func action(c *cli.Context) {
	configFile := c.String(configFlag)
	options, err := getOptions(c)
	if err != nil {
		fmt.Printf("Error reading options: %v", err)
		os.Exit(1)
	}
	m, err := magpie.New(configFile, c.Args(), options...)
	if err != nil {
		fmt.Printf("Error initializing: %v", err)
		os.Exit(1)
	}
	watch := c.Bool(watchFlag)
	if watch {
		m.Watch()
	} else {
		err = m.Collect()
		if err != nil {
			fmt.Printf("Error collecting: %v", err)
			os.Exit(1)
		}
	}
}

func getOptions(c *cli.Context) ([]magpie.Option, error) {
	var options []magpie.Option
	if c.IsSet(tagFlag) {
		options = append(options, magpie.Tags(c.StringSlice(tagFlag)))
	}
	if c.IsSet(prefixFlag) {
		options = append(options, magpie.Prefix(c.String(prefixFlag)))
	}
	if c.IsSet(packageNameFlag) {
		options = append(options, magpie.PackageName(c.String(packageNameFlag)))
	}
	if c.IsSet(unsafeFlag) {
		options = append(options, magpie.Unsafe(c.Bool(unsafeFlag)))
	}
	if c.IsSet(compressFlag) {
		options = append(options, magpie.Compress(c.BoolT(compressFlag)))
	}
	if c.IsSet(fileModeFlag) {
		fileMode := c.String(fileModeFlag)
		if fileMode != "" {
			n, err := strconv.ParseUint(fileMode, 10, 32)
			if err != nil {
				return nil, err
			}
			options = append(options, magpie.OverrideFileMode(uint(n)))
		}
	}
	if c.IsSet(fileModTimeFlag) {
		fileMod := c.String(fileModTimeFlag)
		if fileMod != "" {
			n, err := strconv.ParseInt(fileMod, 10, 64)
			if err != nil {
				return nil, err
			}
			options = append(options, magpie.OverrideFileModTime(n))
		}
	}
	if c.IsSet(outputFlag) {
		options = append(options, magpie.OutputFile(c.String(outputFlag)))
	}
	if c.IsSet(ignoreFlag) {
		options = append(options, magpie.Ignore(c.StringSlice(ignoreFlag)))
	}
	return options, nil
}
