package utils

import (
	"bytes"
	"github.com/yargevad/filepathx"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"text/template"
)

func IsYamlValid(fileName, searchPath string) bool {
	Logger.Debug().Str("name", fileName).Msg("Check if YAML exists")
	fullPath := CheckExistence(fileName, searchPath)
	if fullPath == "" {
		Logger.Error().Str("name", fileName).Msg("YAML not found")
		return false
	}

	// check if workflow template is valid
	Logger.Debug().Str("name", fileName).Msg("Check YAML syntax")
	if !CheckSyntax(fullPath) {
		return false
	}

	return true
}

func CheckExistence(flowName, flowsPath string) string {
	for _, flow := range ListYAML(flowsPath) {
		if flowName == FileWithoutExtension(flow) {
			return flow
		}
	}
	return ""
}

func CheckSyntax(filename string) bool {
	tpl, err := template.ParseFiles(filename)
	if err != nil {
		Logger.Error().Str("name", filename).Msg("YAML template syntax is invalid")
		return false
	}

	var result bytes.Buffer
	_ = tpl.Execute(&result, nil)

	buffer := make(map[interface{}]interface{})
	err = yaml.Unmarshal(result.Bytes(), buffer)
	if err != nil {
		Logger.Error().Str("name", filename).Msg("YAML syntax is invalid")
		return false
	}
	return true
}

func ListYAML(filesPath string) []string {
	if _, err := os.Stat(filesPath); os.IsNotExist(err) {
		if err = os.Mkdir(filesPath, 0700); err != nil {
			Logger.Error().Msg(err.Error())
		}
	}

	searchPath := path.Join(filesPath, "/**/*.yaml")
	result, err := filepathx.Glob(searchPath)
	if err != nil {
		Logger.Error().Msg(err.Error())
		return result
	}
	return result
}
