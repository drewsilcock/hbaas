package cmd

import (
	"fmt"

	"github.com/drewsilcock/hbaas-server/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get server version.",
	Long:  "Gets the current server version.",
	Run:   getVersion,
}

func getVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version.Version, "built", version.BuildTime)
}
