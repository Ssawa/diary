package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Manually upload entries to S3",
	Long: `Manually triggers an upload of a file to S3
	
Can be called with a specific path to upload that one file or without any
path to upload all`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upload called")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
