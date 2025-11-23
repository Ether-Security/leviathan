package core

import (
	"os"
	"path"
	"path/filepath"

	"github.com/Ether-Security/leviathan/libs"
	"github.com/Ether-Security/leviathan/utils"
	"github.com/spf13/viper"
)

func InitConfig(options *libs.Options) error {
	// Check if config file directory exists or create it
	RootFolder := filepath.Dir(utils.NormalizePath(options.ConfigFile))
	if _, err := os.Stat(RootFolder); os.IsNotExist(err) {
		if err = os.Mkdir(RootFolder, 0700); err != nil {
			return err
		}
	}

	v := viper.New()
	v.AddConfigPath(RootFolder)
	v.SetConfigType("yaml")

	v.SetDefault("Environment", map[string]string{
		"workspaces": path.Join(RootFolder, "workspaces"),
		"workflows":  path.Join(RootFolder, "workflows"),
		"modules":    path.Join(RootFolder, "modules"),
		"binaries":   path.Join(RootFolder, "binaries"),
	})

	_ = v.SafeWriteConfigAs(utils.NormalizePath(options.ConfigFile))

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return v.Unmarshal(options)
}
