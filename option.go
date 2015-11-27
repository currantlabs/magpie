package magpie

import "regexp"

type Option func(*config)

func PackageName(packageName string) Option {
	return func(c *config) {
		c.Package = packageName
	}
}

func Tags(tags []string) Option {
	return func(c *config) {
		if len(c.Tags) == 0 {
			c.Tags = tags
			return
		}
		t := make([]string, 0, len(c.Tags))
		copy(t, c.Tags)
		t = append(t, tags...)
		c.Tags = t
	}
}

func OutputFile(outputFile string) Option {
	return func(c *config) {
		c.Output = outputFile
	}
}

func Prefix(prefix string) Option {
	return func(c *config) {
		c.Prefix = prefix
	}
}

func Unsafe(unsafe bool) Option {
	return func(c *config) {
		c.Unsafe = unsafe
	}
}

func Compress(compress bool) Option {
	return func(c *config) {
		c.Compress = compress
	}
}

func OverrideFileMode(mode uint) Option {
	return func(c *config) {
		c.Mode = mode
	}
}

func OverrideFileModTime(modTime int64) Option {
	return func(c *config) {
		c.ModTime = modTime
	}
}

func Ignore(patterns []string) Option {
	return func(c *config) {
		var ignores []*regexp.Regexp
		if c.Ignore != nil {
			ignores := make([]*regexp.Regexp, 0, len(c.Ignore))
			copy(ignores, c.Ignore)
		}
		for _, pattern := range patterns {
			ignores = append(ignores, regexp.MustCompile(pattern))
		}
		c.Ignore = ignores
	}
}

func CreateFileSystem(create bool) Option {
	return func(c *config) {
		c.FileSystem = create
	}
}
