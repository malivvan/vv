package editor

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed syntax themes
var assets embed.FS

var Assets = struct {
	Themes []Asset
	Syntax []Asset
}{
	Themes: loadAssets("themes", ".micro"),
	Syntax: loadAssets("syntax", ".yaml"),
}

type Asset struct {
	Name string
	Data []byte
}

func (a Asset) String() string {
	return a.Name
}

func loadAssets(directory, extension string) []Asset {
	dir, err := assets.Open(directory)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	var files []Asset
	entries, err := fs.ReadDir(assets, directory)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if ext == extension && !entry.IsDir() {
			file, err := assets.Open(filepath.Join(directory, entry.Name()))
			if err != nil {
				continue
			}
			defer file.Close()
			data, err := io.ReadAll(file)
			if err == nil {
				files = append(files, Asset{strings.TrimSuffix(entry.Name(), ext), data})
			}
		}
	}
	return files
}
