package cmd

import (
	"log"
	"os"
	"path"

	"github.com/cj-dimaggio/diary/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var home string

var rootCmd = &cobra.Command{
	Use:   "diary",
	Short: "Maintain diary entries",
	Long:  `Create, maintain, and backup your daily text entries`,
	Run: func(cmd *cobra.Command, args []string) {
		entry := utils.NewEntry(utils.Entry{})
		if err := utils.StartEntry(entry); err != nil {
			utils.Fail(err)
		}

		if viper.GetBool("s3.enabled") {
			if err := utils.BackupFile(entry); err != nil {
				utils.Fail(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Fail(err)
	}
}

func init() {
	// Find home directory.
	var err error
	home, err = homedir.Dir()
	if err != nil {
		utils.Fail(err)
	}

	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.diary.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print verbose output")

	// Viper has support for environment variable bindings but we want a different precedence order
	// here than what they provide
	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		defaultEditor = "vim"
	}

	rootCmd.Flags().StringP("editor", "e", defaultEditor, "The editor to spawn")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())

	viper.SetDefault("file.base", path.Join(home, ".diary"))
	viper.SetDefault("file.template.path", "2006/1/2-Mon-Jan-2006.md")
	viper.SetDefault("file.template.new", "# Monday January 2, 2006\n")
	viper.SetDefault("file.template.append", "\n## At 3:04pm...\n")

	viper.SetDefault("s3.enabled", false)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".diary" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".diary")
	}

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

// initLogger initializes the verbose logger
func initLogger() {
	if viper.GetBool("verbose") {
		utils.Verbose = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}
}
