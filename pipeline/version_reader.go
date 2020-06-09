package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/figment-networks/indexing-engine/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// NewVersionReader VersionReader constructor
func NewVersionReader(dir string) *versionReader {
	return &versionReader{VersionsDir: dir}
}

type versionReader struct {
	VersionsDir string
}

type Mode byte

const (
	ModeAll Mode = iota
	ModeUp
	ModeVersion
)

// VersionTasks slice of task names
type VersionTasks []string

// All combines all files
func (p *versionReader) All() (*int64, VersionTasks, error) {
	dir := p.VersionsDir

	files, err := p.getDirFiles(dir)
	if err != nil {
		return nil, nil, err
	}

	var lastVersionNumber int64
	var combinedTasks VersionTasks
	for _, f := range files {
		filename := f.Name()

		fileVersionNumber, err := p.getVersionNumber(filename)
		if err != nil {
			return nil, nil, err
		}

		fileTasks, err := p.getTasks(dir, filename)
		if err != nil {
			return nil, nil, err
		}

		lastVersionNumber = *fileVersionNumber
		combinedTasks = append(combinedTasks, fileTasks...)
	}
	return &lastVersionNumber, utils.UniqueStr(combinedTasks), nil
}

// Up combines only files that are greater than currentVersion
func (p *versionReader) Up(currentVersion *int64) (*int64, VersionTasks, error) {
	dir := p.VersionsDir

	files, err := p.getDirFiles(dir)
	if err != nil {
		return nil, nil, err
	}

	var lastVersionNumber int64
	var combinedTasks VersionTasks
	for _, f := range files {
		filename := f.Name()

		fileVersionNumber, err:= p.getVersionNumber(filename)
		if err != nil {
			return nil, nil, err
		}

		if *fileVersionNumber > *currentVersion {
			fileTasks, err := p.getTasks(dir, filename)
			if err != nil {
				return nil, nil, err
			}

			lastVersionNumber = *fileVersionNumber
			combinedTasks = append(combinedTasks, fileTasks...)
		}
	}
	return &lastVersionNumber, utils.UniqueStr(combinedTasks), nil
}

// Version uses version file for provided desiredVersion
func (p *versionReader) Version(desiredVersion *int64) (*int64, VersionTasks, error) {
	dir := p.VersionsDir

	files, err := p.getDirFiles(dir)
	if err != nil {
		return nil, nil, err
	}

	var fileTasks VersionTasks
	found := false
	for _, f := range files {
		filename := f.Name()

		fileVersionNumber, err:= p.getVersionNumber(filename)
		if err != nil {
			return nil, nil, err
		}

		if *fileVersionNumber == *desiredVersion {
			fileTasks, err = p.getTasks(dir, filename)
			if err != nil {
				return nil, nil, err
			}
			found = true
			break
		}
	}

	if !found {
		return nil, nil, errors.New(fmt.Sprintf("desired version %d not found", desiredVersion))
	}

	return desiredVersion, fileTasks, nil
}

// getTasks gets tasks from json files from given directory
func (p *versionReader) getTasks(dir string, filename string) (VersionTasks, error) {
	file := fmt.Sprintf("%s/%s", dir, filename)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var fileTasks VersionTasks
	err = json.Unmarshal(data, &fileTasks)
	if err != nil {
		return nil, err
	}
	return fileTasks, nil
}

// getVersionNumber gets version number from filename
func (p *versionReader) getVersionNumber(filename string) (*int64, error) {
	filenameParts := strings.Split(filename, "_")
	fileVersion, err := strconv.ParseInt(filenameParts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	return &fileVersion, nil
}

// getDirFiles gets list of files in directory
func (p *versionReader) getDirFiles(dir string) ([]os.FileInfo, error) {
	if dir == "" {
		return nil, ErrVersionsDirNotSet
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return files, nil
}