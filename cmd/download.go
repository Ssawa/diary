package cmd

import (
	"github.com/cj-dimaggio/diary/utils"
	"github.com/spf13/cobra"
)

var list bool

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download entries from S3",
	Long: `Download entries that have been backed up to S3
	
Can be used with a specific path to download or without a specific
path to download all or use "--list" or "-l" to list the existing files

Warning that this process might not work well if you've changed your key prefix`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if list {
			if err := utils.ListExistingKeys(); err != nil {
				utils.Fail(err)
			}
			return
		}

		if len(args) == 1 {
			if err := utils.DownloadEntry(utils.NewEntry(utils.Entry{
				RelativePath: args[0],
			})); err != nil {
				utils.Fail(err)
			}
			return
		}

		if err := utils.DownloadAll(); err != nil {
			utils.Fail(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolVarP(&list, "list", "l", false, "list files on remote")
}
