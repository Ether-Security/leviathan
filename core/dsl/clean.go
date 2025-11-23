package dsl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ether-Security/leviathan/utils"
	"github.com/yargevad/filepathx"
)

func CleanWorkspace(workspace string, whitelist []string) bool {
	items, err := filepathx.Glob(fmt.Sprintf("%s/**", workspace))
	if err != nil {
		return false
	}

WhitelistLoop:
	for _, item := range items {
		for _, wl := range whitelist {
			rel, _ := filepath.Rel(item, wl)

			// Check if the file or directory is under a whitelisted directory
			if strings.HasSuffix(rel, "..") {
				continue WhitelistLoop
			}

			// Check if the item is above a whitelisted file or directory
			if !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
				continue WhitelistLoop
			}
		}

		fi, err := os.Stat(item)
		if err != nil {
			continue
		}

		if fi.Mode().IsDir() {
			utils.Logger.Debug().Str("folder", item).Msg("Delete folder")
			os.RemoveAll(item)
		} else if fi.Mode().IsRegular() {
			utils.Logger.Debug().Str("file", item).Msg("Delete file")
			os.Remove(item)
		}
	}
	return true
}
