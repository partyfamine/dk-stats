package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dk-stats",
	Short: "Distrokid stats",
	Long:  "Distrokid stats",
}

func init() {
	rootCmd.AddCommand(monthlyCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
