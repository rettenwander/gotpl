package gotpl

import (
	"embed"
	"fmt"
	"path"
	"path/filepath"
)

type file struct {
	name string
	path string
}

// readDir reads the embed.FS and returns all .html files inside the given directory.
func readDir(fs embed.FS, dir ...string) ([]file, error) {
	var files []file

	fullDir := path.Join(dir...)

	if ok := exists(fs, fullDir); !ok {
		return nil, fmt.Errorf("directory not found: %s", fullDir)
	}

	allFiles, err := fs.ReadDir(fullDir)
	if err != nil {
		return nil, err
	}

	for _, f := range allFiles {
		if f.IsDir() || filepath.Ext(f.Name()) != ".html" {
			continue
		}

		files = append(files, file{name: f.Name(), path: path.Join(fullDir, f.Name())})
	}

	return files, nil
}

// getPaths returns the paths from a slice of files.
func getPaths(files []file) []string {
	var p []string
	for _, f := range files {
		p = append(p, f.path)
	}
	return p
}

// exists reports whether the given file or directory exists in the embed.FS.
func exists(fs embed.FS, path string) bool {
	f, err := fs.Open(path)
	if err != nil {
		return false
	}
	f.Close()
	return true
}
