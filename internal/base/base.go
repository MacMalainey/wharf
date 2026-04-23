// Licensed Materials - Property of IBM
// Copyright IBM Corp. 2023.

package base

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/zosopentools/wharf/internal/util"
)

// Config holds all configuration and environment information for Wharf
type Config struct {
	goenv     map[string]string
	BuildTags map[string]bool
	ImportDir string
	Cache     string
}

// NewConfig creates a new Config instance with default values from the Go environment.
// If workspacePath is empty, it runs 'go env' in the current directory.
// Otherwise, it runs 'go env' in the provided workspace directory.
func NewConfig(workspacePath string) (*Config, error) {
	// Run go env in the specified directory (or current directory if empty)
	goenv, err := util.GoEnv(workspacePath)
	if err != nil {
		return nil, fmt.Errorf("unable to inspect Go environment (cannot execute 'go env'): %v", err)
	}

	cfg := &Config{
		goenv:     goenv,
		BuildTags: make(map[string]bool),
	}

	// Set tags that Go figures out from the environment, such as GOARCH, CGO, and GOVERSION
	cfg.BuildTags[goenv["GOARCH"]] = true
	cfg.BuildTags[build.Default.Compiler] = true
	if goenv["CGO_ENABLED"] == "1" {
		cfg.BuildTags["cgo"] = true
	}

	var vnum int
	if match := regexp.MustCompile(`go1\.(\d+)(?:(?:\.|-).+)?$`).FindStringSubmatch(goenv["GOVERSION"]); match != nil {
		var err error
		vnum, err = strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("go version minor number unable to parse to int")
		}
	} else {
		vnum = 18
		fmt.Fprintf(os.Stderr, "unknown go version number (%v) - assuming go1.18\n", goenv["GOVERSION"])
	}

	for vnum >= 0 {
		cfg.BuildTags[fmt.Sprintf("go1.%v", vnum)] = true
		vnum -= 1
	}

	// Use GOWORK from the environment (which was obtained from the specified directory)
	goWorkDir := filepath.Dir(cfg.GOWORK())

	// Initialize some variables here to default values (can be overwritten)
	cfg.Cache = filepath.Join(goWorkDir, ".wharf_cache") // TODO: move this to TMPDIR

	// TODO: make this relative to the position of the GOWORK folder
	// so that `go work use` uses a relative position instead of absolute
	cfg.ImportDir = filepath.Join(goWorkDir, "wharf_port")

	return cfg, nil
}

func (c *Config) GOOS() string {
	return c.goenv["GOOS"]
}

func (c *Config) GOARCH() string {
	return c.goenv["GOARCH"]
}

func (c *Config) GOWORK() string {
	return c.goenv["GOWORK"]
}

func (c *Config) GoEnv(key string) string {
	return c.goenv[key]
}
