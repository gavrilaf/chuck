package utils

import (
	"fmt"
	"github.com/spf13/afero"
	. "path"
	"strings"
	"time"
)

/*
 * Create folder for the traffic log.
 * Check if 'folder' exists in the filesystem 'fs'. Create it if not.
 * If 'createNewFolder' is true, create unique named subfloder inside folder.
 * Return unique folder name (or empty string) and full path to the folder.
 */
func PrepareStorageFolder(fs afero.Fs, folder string, createNewFolder bool) (string, string, error) {
	folder = strings.Trim(folder, " \\/")
	logDirExists, err := afero.DirExists(fs, folder)
	if err != nil {
		return "", "", err
	}

	if !logDirExists {
		err := fs.Mkdir(folder, 0777)
		if err != nil {
			return "", "", err
		}
	}

	name := ""
	path := folder
	if createNewFolder {
		tm := time.Now()
		name = fmt.Sprintf("%d_%d_%d_%d_%d_%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
		path = Join(folder, name)

		err = fs.Mkdir(path, 0777)
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}
