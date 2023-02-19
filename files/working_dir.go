package files

import (
	"os"
	"path/filepath"
)

func WorkingDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
