package cmd

import (
	"fmt"

	"github.com/Ether-Security/leviathan/core"
	"github.com/Ether-Security/leviathan/libs"
	"github.com/Ether-Security/leviathan/utils"
	"github.com/spf13/cobra"
)

var options = libs.Options{}

var rootCmd = &cobra.Command{
	Use:   libs.NAME,
	Short: fmt.Sprintf("%s - %s", libs.NAME, libs.DESC),
}

func init() {
	rootCmd.PersistentFlags().StringVar(&options.ConfigFile, "config", libs.CONFIGFILE, "File from where config is loaded")

	rootCmd.PersistentFlags().StringVar(&options.Log.Directory, "log", libs.LOGDIR, "Directory where logs will be saved")
	rootCmd.PersistentFlags().BoolVar(&options.Log.JSON, "log-json", false, "Output log as JSON")
	rootCmd.PersistentFlags().BoolVarP(&options.Log.Debug, "debug", "d", false, "Set log to debug level")
	rootCmd.PersistentFlags().BoolVarP(&options.Log.Quiet, "quiet", "q", false, "Remove all logs from output")

	cobra.OnInitialize(initCore)
}

func initCore() {
	// Init config and create config file
	if err := core.InitConfig(&options); err != nil {
		utils.Logger.Fatal().Msgf(err.Error())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Logger.Fatal().Msg(err.Error())
	}
}
