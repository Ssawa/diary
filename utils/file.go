package utils

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Entry represents an entry in the diary
type Entry struct {
	timestamp    time.Time
	relativePath string
	basePath     string
	fullPath     string
}

// NewEntry forms a new entry in the diary
func NewEntry() Entry {
	now := time.Now()
	relativePath := now.Format(viper.GetString("file.filename"))
	basePath := viper.GetString("file.base")
	fullPath := path.Join(basePath, relativePath)
	return Entry{
		timestamp:    now,
		relativePath: relativePath,
		basePath:     basePath,
		fullPath:     fullPath,
	}
}

// StartEntry is the entry point for creating a new entry file and allowing the user to edit it
func StartEntry(entry Entry) error {
	Verbose.Println("Today's File: ", entry.fullPath)
	if err := os.MkdirAll(path.Dir(entry.fullPath), os.ModePerm); err != nil {
		return err
	}
	if err := templateEntry(entry); err != nil {
		return err
	}
	if err := spawnEditor(entry.fullPath); err != nil {
		return nil
	}
	Verbose.Println("File closed")
	return nil
}

// spawnEditor spawns an editor and waits for it to close
func spawnEditor(filename string) error {
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

// templateEntry checks to see if this is a new file and, if it is,
// creates it with the templated data rendered. If it already exists,
// and the "append_template" config has been set then we append it to the
// file
func templateEntry(entry Entry) error {
	var file *os.File
	if _, err := os.Stat(entry.fullPath); os.IsNotExist(err) {
		Verbose.Println("File does not exist, creating template...")
		newTemplate := entry.timestamp.Format(viper.GetString("file.template.new"))
		file, err = os.OpenFile(entry.fullPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		if _, err := file.WriteString(newTemplate); err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(entry.fullPath, os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil
		}
	}
	return appendEntry(file, entry.timestamp)
}

// appendEntry checks if the append_template has been set and, if it has,
func appendEntry(file *os.File, now time.Time) error {
	appendTemplate := now.Format(viper.GetString("file.template.append"))
	_, err := file.WriteString(appendTemplate)
	return err
}
