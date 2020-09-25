package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
	Replace []Replace
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

func getProjectsInDir(workDir string) ([]string, error) {
	dir, err := os.Open(workDir)
	if err != nil {
		fmt.Println("unable to open work dir", err)
		return nil, err
	}

	names, err := dir.Readdirnames(-1)
	if err != nil {
		fmt.Println("unable to read directory names")
		return nil, err
	}

	projects := make([]string, 0)

	for _, name := range names {
		projectDir := path.Join(workDir, name)

		if _, err := os.Stat(path.Join(projectDir, ".git")); err != nil {
			continue
		}

		projects = append(projects, projectDir)
	}

	return projects, nil
}

func wireSiblingsInPath(workDir string) error {
	projects, err := getProjectsInDir(workDir)
	if err != nil {
		return err
	}

	modInfo := make(modList, 0)

	for _, projectDir := range projects {
		goMod, err := goModFromRepoPath(projectDir)
		if err != nil {
			// This happens when we try to read from directories that don't contain a go module, or from a file
			continue
		}

		modInfo = append(modInfo, &ModDir{
			gm:  goMod,
			dir: projectDir,
		})
	}

	return wireSiblings(modInfo)
}

func wireSiblings(modInfo modList) error {
	for _, mod := range modInfo {
		for _, req := range mod.gm.Require {
			// We require a sibling project, so lets add a replace
			reqMod := modInfo.getRequired(&req)
			if reqMod != nil {
				fmt.Printf("Adding go.mod replace for %s: %s => ../%s\n",
					mod.gm.Module.Path,
					req.Path,
					path.Base(reqMod.dir),
				)
				cmd := exec.Command(
					"go",
					"mod",
					"edit",
					"-replace",
					fmt.Sprintf("%s=../%s", req.Path, path.Base(reqMod.dir)),
				)
				cmd.Dir = mod.dir

				if err := cmd.Run(); err != nil {
					return fmt.Errorf("unable to add mod replace: %v", err)
				}
			}
		}
	}

	return nil
}

func goModFromRepoPath(repoDir string) (*GoMod, error) {
	var output bytes.Buffer

	cmd := exec.Command("go", "mod", "edit", "-json")
	cmd.Dir = repoDir
	cmd.Stdout = &output

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	goMod := &GoMod{}

	if err := json.Unmarshal(output.Bytes(), goMod); err != nil {
		return nil, err
	}

	return goMod, nil
}

type ModDir struct {
	gm  *GoMod
	dir string
}

type modList []*ModDir

func (l modList) getRequired(req *Require) *ModDir {
	for _, mod := range l {
		if mod.gm.Module.Path == req.Path {
			return mod
		}
	}

	return nil
}

func (l modList) containsRequired(req *Require) bool {
	for _, mod := range l {
		if mod.gm.Module.Path == req.Path {
			return true
		}
	}

	return false
}

func dirExists(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func fileExists(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}

	return stat.Mode().IsRegular()
}
