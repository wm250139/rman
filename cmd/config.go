package cmd

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config represents a global/user config, can be used rather than individual files
type Config map[string]*Workspace

// Workspace represents a list of repositories that make up a single workspace
type Workspace struct {
	Repositories []string `toml:"repos"`
}

func getGlobalConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "unable to find user home directory")
		return nil, err
	}

	gf := filepath.Join(home, ".config", "rman", "workspaces.toml")
	if !fileExists(gf) {
		return nil, fmt.Errorf("no global configuration found at %s", gf)
	}

	gc := &Config{}
	if err := parseConfigFile(gf, gc); err != nil {
		return nil, err
	}
	return *gc, nil
}

func getWorkspaceConfig(name string) (*Workspace, error) {
	if fileExists(name) {
		fileConfig := &Workspace{}
		if err := parseConfigFile(name, fileConfig); err != nil {
			return nil, err
		}
		return fileConfig, nil
	}

	gc, err := getGlobalConfig()
	if err != nil {
		return nil, err
	}

	config, ok := gc[name]
	if !ok {
		return nil, fmt.Errorf("unable to find valid config for '%s'", name)
	}

	return config, nil
}

func parseConfigFile(name string, target interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return toml.Unmarshal(bytes, target)
}
