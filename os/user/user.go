package user

import (
	"os"
	"os/user"
	"strings"
)

// HomeDir returns a best guess of the current running user's home directory
// either by using an existing HOME environment variable or the user's database.
func HomeDir() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

// ExpandTilde expands a file path beginning with tilde to its HOME path
// equivalent.
func ExpandTilde(p string) (string, error) {
	if p == "~" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return usr.HomeDir, nil
	}

	sep := string(os.PathSeparator)
	if len(p) > 1 && p[:1+len(sep)] == "~"+sep {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		dir := usr.HomeDir

		return strings.Replace(p, "~", dir, 1), nil
	}

	return p, nil
}
