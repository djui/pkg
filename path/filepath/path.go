package filepath

import "path/filepath"

// SplitBase returns a file path's basename and extension.
func SplitBase(path string) (string, string) {
	ext := filepath.Ext(path)
	return path[0 : len(path)-len(ext)], ext
}
