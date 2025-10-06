package dirs

import (
	"os"
	"path/filepath"
	"strings"
	"vic3-data-reader/internal/env"
	"vic3-data-reader/internal/read/files"
)

func dataRootPath() (string, error) {
	var fp string
	dir, err := env.Vic3Dir.GetValue()
	if err == nil {
		fp = filepath.Join(dir, "game/common")
	}
	return fp, err
}

type DataDir string

const (
	BuildingGroups         DataDir = "building_groups"
	Buildings              DataDir = "buildings"
	Goods                  DataDir = "goods"
	ProductionMethodGroups DataDir = "production_method_groups"
	ProductionMethods      DataDir = "production_methods"
	Technologies           DataDir = "technology/technologies"
)

func (d DataDir) DirPath() (string, error) {
	var fp string
	rt, err := dataRootPath()
	if err == nil {
		fp = filepath.Join(rt, string(d))
	}
	return fp, err
}

func (d DataDir) Files() ([]files.DataFile, error) {
	var dfs []files.DataFile
	var err error

	// load dir and read contents
	dir, err := d.DirPath()
	if err != nil {
		return dfs, err
	}
	fps, err := os.ReadDir(dir)
	if err != nil {
		return dfs, err
	}

	// filter files in dir
	isTxt := func(entry os.DirEntry) bool {
		return strings.HasSuffix(entry.Name(), ".txt")
	}
	isDummy := func(entry os.DirEntry) bool {
		return strings.Contains(entry.Name(), "dummy")
	}
	isValid := func(entry os.DirEntry) bool {
		return !entry.IsDir() && isTxt(entry) && !isDummy(entry)
	}

	for _, entry := range fps {
		if isValid(entry) {
			fp := filepath.Join(dir, entry.Name())
			dfs = append(dfs, files.DataFile(fp))
		}
	}

	return dfs, err
}
