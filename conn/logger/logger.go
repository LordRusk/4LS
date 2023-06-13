/*
github.com/lordrusk/dctl/logger ported and reworked for 4LS's purposes

log's to an io.Writer and a file
*/
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	errs "github.com/pkg/errors"
)

/* allows for default log.Logger fallback if
 * logger fails to create */
type Logger interface {
	Close() error

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Flags() int
	Output(calldepth int, s string) error
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	SetFlags(flag int)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
	Writer() io.Writer
}

type logger struct {
	*log.Logger
	close func() error // closes the file
}

func (l *logger) Close() error { return l.close() }

type default_logger struct {
	*log.Logger
}

func (_ default_logger) Close() error { return nil }

func new_default_logger() *default_logger {
	return &default_logger{
		Logger: log.Default(),
	}
}

/* returns the parent directories of a path */
func Parent(path string) string {
	strs := strings.Split(path, "/")
	return strings.Join(strs[:len(strs)-1], "/")
}

/* Allows for logging to an io.Writer and a file easily
 *
 * if it fails to create the file, it fall's back on
 * log.Default() with a dummy .Close() */
func New(out io.Writer, prefix string, flag int, path string) Logger {
	path = strings.ReplaceAll(path, "~", "$HOME")
	path = os.ExpandEnv(path)
	spath := strings.Split(path, "/")
	var f *os.File
	var err error
	var mwr io.Writer

	if err = os.MkdirAll(strings.Join(spath[0:len(spath)-1], "/"), os.ModePerm); err != nil {
		fmt.Printf("%s: Cannot create '%s'", err, path)
		goto DEFAULT_LOGGER
	}
	if _, err = os.Create(path); err != nil {
		fmt.Printf("%s: Failed to create %s", err, path)
		goto DEFAULT_LOGGER
	}
	if f, err = os.OpenFile(path, os.O_RDWR, 0666); err != nil {
		fmt.Printf("%s Failed to open %s", err, path)
		goto DEFAULT_LOGGER
	}
	mwr = io.MultiWriter(out, f)
	return &logger{
		Logger: log.New(mwr, prefix, flag),
		close: func() error {
			if err := f.Close(); err != nil {
				return errs.Wrap(err, "unable to close file")
			}
			return nil
		},
	}

DEFAULT_LOGGER:
	return new_default_logger()
}
