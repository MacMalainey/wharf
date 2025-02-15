// Licensed Materials - Property of IBM
// Copyright IBM Corp. 2023.
// US Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
package util

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Copies a module to the given path
func CopyModule(srcdir string, dstdir string, modpath string) error {
	// Copy module to workspace
	if err := copyAll(dstdir, srcdir); err != nil {
		return err
	}

	return generate(filepath.Join(dstdir, "go.mod"), goModTemplate(modpath))
}

func GitCloneModule(repo string, dest string, modpath string) error {
	// Copy module to workspace
	if err := GitClone(repo, dest); err != nil {
		return err
	}

	return generate(filepath.Join(dest, "go.mod"), goModTemplate(modpath))
}

// Util copy function
func copyAll(dst, src string) error {
	srcInfo, srcErr := os.Lstat(src)
	if srcErr != nil {
		return fmt.Errorf("util: %w", srcErr)
	}
	_, dstErr := os.Lstat(dst)
	if dstErr == nil {
		return fmt.Errorf("util: will not overwrite %q", dst)
	}
	if !errors.Is(dstErr, fs.ErrNotExist) {
		return fmt.Errorf("util: %w", dstErr)
	}
	switch mode := srcInfo.Mode(); mode & os.ModeType {
	case os.ModeSymlink:
		return fmt.Errorf("util: will not copy symbolic link")
	case os.ModeDir:
		return copyDir(dst, src)
	case 0:
		return CopyFile(dst, src)
	default:
		return fmt.Errorf("util: cannot copy file with mode %v", mode)
	}
}

// Util copy file function
func CopyFile(dst, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("util: %w", err)
	}
	defer srcf.Close()
	dstf, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("util: %w", err)
	}
	defer dstf.Close()
	if _, err := io.Copy(dstf, srcf); err != nil {
		return fmt.Errorf("util: cannot copy %q to %q: %w", src, dst, err)
	}
	return nil
}

// Util copy directory function
func copyDir(dst, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("util: %w", err)
	}
	defer srcf.Close()
	if err := os.MkdirAll(dst, 0777); err != nil {
		return fmt.Errorf("util: %w", err)
	}
	for {
		names, err := srcf.Readdirnames(100)
		for _, name := range names {
			if err := copyAll(filepath.Join(dst, name), filepath.Join(src, name)); err != nil {
				return fmt.Errorf("util: %w", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("util: error reading directory %q: %w", src, err)
		}
	}
	return nil
}

func goModTemplate(path string) string {
	// TODO: figure out a way to make sure we keep the same versions for stuff
	return "// Generated by IBM Wharf; DO NOT EDIT.\nmodule " + path + "\n"
}

func generate(file, content string) error {
	if _, err := os.Stat(file); err == nil {
		return nil
	}
	if err := ioutil.WriteFile(file, []byte(content), 0666); err != nil {
		return fmt.Errorf("util: could not generate %q: %w", file, err)
	}
	return nil
}
