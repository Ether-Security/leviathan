package dsl

import (
	"os"

	"github.com/Ether-Security/leviathan/utils"
	cp "github.com/otiai10/copy"
)

func CopyFile(src string, dest string) (err error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, data, 0600)
	return err
}

func CopyDir(src string, dest string) (err error) {
	cp.Copy(src, dest)
	if err != nil {
		return err
	}
	return nil
}

func Append(dest string, src string) (err error) {
	if !utils.FileExists(src) {
		return
	}
	data := utils.GetFileContent(src)
	_, err = utils.AppendToContent(dest, data)
	return err
}
