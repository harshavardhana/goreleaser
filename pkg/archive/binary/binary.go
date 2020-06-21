// Package binary implements the Archive interface providing binary archiving
// and compression.
package binary

import (
	"io"
	"os"
)

// Archive binary struct.
type Archive struct {
	io.WriteCloser
}

// Close all closeables.
func (a Archive) Close() error {
	return a.WriteCloser.Close()
}

// New binary archive.
func New(target io.WriteCloser) Archive {
	return Archive{target}
}

// Add a file to the binary archive.
func (a Archive) Add(name, path string) (err error) {
	file, err := os.Open(path) // #nosec
	if err != nil {
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return
	}
	if info.IsDir() {
		return
	}
	_, err = io.Copy(a.WriteCloser, file)
	return err
}
