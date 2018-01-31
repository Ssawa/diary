package utils

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// SpawnEditor spawns an editor and waits for it to close
func SpawnEditor(filename string) error {
	editor := strings.Split(viper.GetString("editor"), " ")
	Verbose.Printf("Spawning editor '%s' for file '%s'", editor, filename)
	cmd := exec.Command(editor[0], append(editor[1:], filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return err
}
