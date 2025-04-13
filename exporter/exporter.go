package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const fileExtensionsToProcess = ".gpg"

type PassExporter struct {
	baseDir string
	data    map[string]interface{}
}

func (p *PassExporter) Export() error {
	fileCount := 0
	successCount := 0

	return filepath.Walk(p.baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only skip hidden subdirectories, not the base directory itself
		if info.IsDir() && path != p.baseDir && strings.HasPrefix(info.Name(), ".") {
			logrus.Warnf("skipping directory %s/%s", p.baseDir, path)
			return filepath.SkipDir
		}

		// Process only .gpg files
		if !info.IsDir() && strings.HasSuffix(info.Name(), fileExtensionsToProcess) {
			fileCount++
			logrus.Infof("processing file %s/%s", path, info.Name())

			// get the relative path
			relPath, err := filepath.Rel(p.baseDir, path)
			if err != nil {
				return err
			}
			cmdStr := convertRelativePathToPassCommand(relPath)
			logrus.Info(fmt.Sprintf("executing command: %s", cmdStr))
			cmd := exec.Command("sh", "-c", cmdStr)

			// capture stdout & stderr

			output, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}

			pathComponents := strings.Split(relPath, string(os.PathSeparator))

			// key name (gpg file without extension)
			key := strings.TrimSuffix(pathComponents[len(pathComponents)-1], ".gpg")
			// Build the nested structure
			current := p.data
			for i := 0; i < len(pathComponents)-1; i++ {
				component := pathComponents[i]

				// Create new map if this component doesn't exist
				if _, exists := current[component]; !exists {
					current[component] = make(map[string]interface{})
				}

				// Navigate deeper
				current = current[component].(map[string]interface{})
			}

			// Add the value (trim newlines)
			current[key] = strings.TrimSpace(string(output))
			successCount++
			logrus.Infof("successfully processed %s", relPath)
		}
		return nil
	})

}

func convertRelativePathToPassCommand(relPath string) string {
	return fmt.Sprintf("pass show %s", strings.TrimSuffix(relPath, ".gpg"))
}
