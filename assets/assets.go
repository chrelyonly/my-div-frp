package assets

import (
	"io/fs"
	"net/http"
)

var (
	// read-only filesystem created by "embed" for embedded files
	content fs.FS

	FileSystem http.FileSystem

	// if prefix is not empty, we get file content from disk
	prefixPath string
)

// if path is empty, load assets in memory
// or set FileSystem using disk files
func Load(path string) {
	prefixPath = path
	if prefixPath != "" {
		FileSystem = http.Dir(prefixPath)
	} else {
		FileSystem = http.FS(content)
	}
}

func Register(fileSystem fs.FS) {
	subFs, err := fs.Sub(fileSystem, "static")
	if err == nil {
		content = subFs
	}
}
