package main

import "regexp"

type Option func(*config)

func Live(live bool) Option {
	return func(c *config) {
		c.live = live
	}
}

func NestName(nestName string) Option {
	return func(c *config) {
		c.nest = nestName
	}
}

func PackageName(packageName string) Option {
	return func(c *config) {
		c.packageName = packageName
	}
}

func Tags(tags []string) Option {
	return func(c *config) {
		if len(c.tags) == 0 {
			c.tags = tags
			return
		}
		t := make([]string, 0, len(c.tags))
		copy(t, c.tags)
		t = append(t, tags...)
		c.tags = t
	}
}

func OutputFile(outputFile string) Option {
	return func(c *config) {
		c.output = outputFile
	}
}

func Prefix(prefix string) Option {
	return func(c *config) {
		c.prefix = prefix
	}
}

func Unsafe(unsafe bool) Option {
	return func(c *config) {
		c.unsafe = unsafe
	}
}

func Compress(compress bool) Option {
	return func(c *config) {
		c.compress = compress
	}
}

func OverrideFileMode(mode uint) Option {
	return func(c *config) {
		c.fileMode = mode
	}
}

func OverrideFileModTime(modTime int64) Option {
	return func(c *config) {
		c.fileModTime = modTime
	}
}

func Ignore(patterns []string) Option {
	return func(c *config) {
		var ignores []*regexp.Regexp
		if c.ignore != nil {
			ignores := make([]*regexp.Regexp, 0, len(c.ignore))
			copy(ignores, c.ignore)
		}
		for _, pattern := range patterns {
			ignores = append(ignores, regexp.MustCompile(pattern))
		}
		c.ignore = ignores
	}
}

func CreateFileSystem(create bool) Option {
	return func(c *config) {
		c.fileSystem = create
	}
}
