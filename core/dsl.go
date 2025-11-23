package core

import (
	b64 "encoding/base64"
	"os"

	"github.com/Ether-Security/leviathan/core/dsl"
	"github.com/Ether-Security/leviathan/utils"
)

func (r *Runner) initDSL() {
	r.DSL = map[string]interface{}{
		"Append": func(dst, src string) bool {
			err := dsl.Append(dst, src)
			return err != nil
		},
		"Base64": func(value string) string {
			return b64.StdEncoding.EncodeToString([]byte(value))
		},
		"CleanWorkspace": func() bool {
			if r.Options.Scan.NoClean {
				return true
			}
			return dsl.CleanWorkspace(r.Workspace, r.Reports)
		},
		"CopyFile": func(src, dst string) bool {
			err := dsl.CopyFile(src, dst)
			return err != nil
		},
		"CopyDir": func(src, dst string) bool {
			err := dsl.CopyDir(src, dst)
			return err != nil
		},
		"CreateFolder": func(folderPath string) bool {
			err := os.MkdirAll(folderPath, 0700)
			return err != nil
		},
		"ExecCmd": func(cmd string) bool {
			parsedCmd := r.ParseString(cmd)
			_, err := utils.ExecCmd(parsedCmd, r.Workspace, r.Options.Environment.Binaries, r.Options.Log.Quiet)
			return err != nil
		},
		"ExtractNmapIP": func(src, dest string) bool {
			dsl.ExtractNmapIP(src, dest)
			return true
		},
		"FileExists": func(filename string) bool {
			return utils.FileExists(filename)
		},
	}
}
