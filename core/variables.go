package core

import (
	"bytes"
	"path"
	"text/template"

	"github.com/Ether-Security/leviathan/utils"
	"golang.org/x/exp/slices"
)

var ParamsBlacklist = []string{
	"source",
	"source_id",
	"source_safe",
}

func (r *Runner) initParams() {
	if r.Params == nil {
		r.Params = make(map[string]string)
	}

	r.Params["Target"] = r.Target
	r.Params["Workspace"] = r.Workspace
	r.Params["ReportsDir"] = path.Join(r.Workspace, "reports")
	r.Params["TempDir"] = path.Join(r.Workspace, ".tmp")

	// Check if there is workflow params
	for _, param := range r.Workflow.Params {
		for k, v := range param {
			utils.Logger.Info().Str("param", k).Str("value", v).Msg("Set workflow param")
			r.setParam(k, v)
		}
	}

	// Check if there is user params (override module and workflow params)
	for k, v := range r.Options.Scan.Params {
		r.setParam(k, v)
	}
}

func (r *Runner) setParam(key, value string) {
	// if in blacklist
	if slices.Contains(ParamsBlacklist, key) {
		return
	}
	r.Params[key] = value
}

func (r *Runner) deleteParam(key string) {
	if slices.Contains(ParamsBlacklist, key) {
		return
	}
	delete(r.Params, key)
}

func (r *Runner) ParseTemplate(filePath string) (output []byte, err error) {
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		return output, err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, r.Params)
	output = result.Bytes()
	return output, err
}

func (r *Runner) ParseString(value string) (output string) {
	tpl, err := template.New("value").Parse(value)
	if err != nil {
		utils.Logger.Error().Str("string", value).Msg("Unable to parse string")
		return output
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, r.Params)
	if err != nil {
		utils.Logger.Error().Str("string", value).Msg("Unable to parse string")
		return output
	}

	output = result.String()
	return output
}
