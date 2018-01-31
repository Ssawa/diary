package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/Ssawa/diary/utils"
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
		now := time.Now()
		relativePath := now.Format(viper.GetString("filename_format"))
		basePath := viper.GetString("base")
		fullPath := path.Join(relativePath, basePath)
		os.MkdirAll(path.Dir(fullPath), 0600)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Find home directory.
	var err error
	home, err = homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.diary.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print verbose output")
	rootCmd.Flags().StringP("editor", "e", "vim", "The editor to spawn")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindPFlags(rootCmd.Flags())

	viper.SetDefault("base", path.Join(home, "diary"))
	viper.SetDefault("filename_format", "/2006/1/Mon-Jan-2-2006.md")

	viper.BindEnv("editor", "EDITOR")
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

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}

// initLogger initializes the verbose logger
func initLogger() {
	if viper.GetBool("verbose") {
		utils.Verbose = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}
}
