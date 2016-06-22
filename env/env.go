// Package env provides simple environment primitives.
//
// An environment holds the boundary for standard and logging I/O, process
// environment, and arguments.
//
// Be aware that environment variables are immutable, as opposed to os.env which
// uses syscall.
//
// Example:
//
//     func main() {
//         Main(env.Default)
//     }
//
//     func Main(e env.Env) {
//         e.Flags.Parse()
//         e.Log.Println("parsed flags")
//         h := e.Envvars["HOME"]
//         e.Log.Println(h)
//     }
package env

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

// An Env holds the command line and other values which can be injected into a
// program.
type Env struct {
	In      io.Reader
	Out     io.Writer
	Err     io.Writer
	Envvars map[string]string
	Flags   *flag.FlagSet
	Log     *log.Logger
}

// Default represents a set of expected presents for an environment.
var Default = &Env{
	In:      os.Stdin,
	Out:     os.Stdout,
	Err:     os.Stderr,
	Envvars: map[string]string{},
	Flags:   flag.CommandLine,
	Log:     log.New(os.Stderr, "", log.LstdFlags), // log.std
}

func init() {
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		k, v := parts[0], parts[1]
		Default.Envvars[k] = v
	}
}

// ExpandEnv replaces ${var} or $var in the string according to the values of
// the current environment variables. References to undefined variables are
// replaced by the empty string.
func (e *Env) ExpandEnv(s string) string {
	return os.Expand(s, e.GetEnv)
}

// GetEnv retrieves the value of the environment variable named by the key. It
// returns the value, which will be empty if the variable is not present.
func (e *Env) GetEnv(key string) string {
	v, _ := e.Envvars[key]
	return v
}

// LookupEnv retrieves the value of the environment variable named by the key.
// If the variable is present in the environment the value (which may be empty)
// is returned and the boolean is true. Otherwise the returned value will be
// empty and the boolean will be false.
func (e *Env) LookupEnv(key string) (string, bool) {
	v, ok := e.Envvars[key]
	return v, ok
}
