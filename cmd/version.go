package cmd

import (
	"fmt"

	"github.com/Ether-Security/leviathan/libs"
	"github.com/Ether-Security/leviathan/utils"
	"github.com/spf13/cobra"
)

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get version",
		RunE:  getVersion,
	}
	rootCmd.AddCommand(versionCmd)
}

func getVersion(_ *cobra.Command, _ []string) error {
	fmt.Printf("%v %v by %v\n", utils.Title(libs.NAME), libs.VERSION, libs.AUTHOR)
	return nil
}
