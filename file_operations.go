package mktools

import "os"

// FileExist returns whether the given file or directory exists
func FileExist(path string) bool {
	fi, err := os.Lstat(path)

	if err != nil {
		return !fi.IsDir()
	}

	return !os.IsNotExist(err)
}

// DelFile removes path and any children it contains
func DelFile(path string) error {
	return os.RemoveAll(path)
}
