package main

import (
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// https://github.com/go-yaml/yaml/issues/100
type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

type Project struct {
	Label string
	Path  StringArray
	Skip  StringArray
	Env   map[string]string
}

func (p *Project) getMainPath() string {
	return p.Path[0]
}

func (p *Project) checkProjectRules(step map[interface{}]interface{}) bool {
	for _, pattern := range p.Skip {
		label := step["label"].(string)
		if matched, _ := filepath.Match(pattern, label); matched {
			return false
		}
	}
	return true
}

func (p *Project) checkAffected(changedFiles []string) bool {
	for _, filePath := range p.Path {
		if filePath == "." {
			return true
		}
		normalizedPath := path.Clean(filePath)
		projectDirs := strings.Split(normalizedPath, "/")
		for _, changedFile := range changedFiles {
			changedDirs := strings.Split(changedFile, "/")
			if reflect.DeepEqual(changedDirs[:Min(len(projectDirs), len(changedDirs))], projectDirs) {
				return true
			}
		}
	}
	return false
}
