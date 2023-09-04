package internal

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Workspace interface {
	WorkDir(path string) string
	CacheDir(path string) string
}

func InitWorkspace() (Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	w := &workspace{cwd: cwd}

	if root, ok := resolvePkgRoot(cwd); ok {
		w.root = root
	} else {
		return nil, errors.New("must under some dir with go.mod")
	}

	return w, nil
}

type workspace struct {
	cwd  string
	root string
}

func (w *workspace) WorkDir(path string) string {
	return filepath.Join(w.cwd, path)
}

func (w *workspace) CacheDir(path string) string {
	dir := filepath.Join(w.root, ".gio-compose", path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	return dir
}

func resolvePkgRoot(from string) (string, bool) {
	if from == "" || from == "/" {
		return "", false
	}
	_, err := os.Stat(filepath.Join(from, "go.mod"))
	if os.IsNotExist(err) {
		return resolvePkgRoot(filepath.Dir(from))
	}
	return from, true
}

func WriteFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	if err := f.Truncate(0); err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(data))
	return err
}
