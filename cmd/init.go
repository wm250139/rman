package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var (
	wsConfig string
)

func init() {
	initCommand.Flags().StringVarP(&wsConfig, "config", "c", "workspace.toml", "Path to configuration file")

	rootCommand.AddCommand(initCommand)
}

var initCommand = &cobra.Command{
	Use:   "init [flags] [path]",
	Short: "Initialize a workspace, cloning and configuring repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Fail early if git isn't available
		if exec.Command("git", "--version").Run() != nil {
			_, _ = fmt.Fprintln(os.Stderr, "[error] git is required to use this utility")
			os.Exit(1)
		}

		workDir := "."
		if len(args) == 1 {
			workDir = args[0]
		}

		// Create the directory if it does not exist
		if !dirExists(workDir) {
			if err := os.MkdirAll(workDir, 0755); err != nil {
				fmt.Printf("Unable to create workdir '%s': %s\n", workDir, err)
				return err
			}
		}

		// Change working directory to the workdir
		if err := os.Chdir(workDir); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "unable to change directory to workdir", err)
			os.Exit(1)
		}

		cn, _ := cmd.Flags().GetString("config")
		config, err := getWorkspaceConfig(cn)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "unable to load configuration:", err)
			os.Exit(1)
		}

		modInfo := make(modList, 0)

		for _, repo := range config.Repositories {
			if err := cloneRepo(repo); err != nil {
				return fmt.Errorf("unable to download repo %s: %s", repo, err)
			}

			repoName := strings.TrimSuffix(repo, ".git")
			repoName = repoName[strings.LastIndex(repoName, "/")+1:]

			mod, err := goModFromRepoPath(repoName)
			if err != nil {
				return err
			}

			modDir := &ModDir{
				gm:  mod,
				dir: repoName,
			}

			modInfo = append(modInfo, modDir)
		}

		return wireSiblings(modInfo)
	},
}

func cloneRepo(repo string) error {
	repoName := strings.TrimSuffix(repo, ".git")
	repoName = repoName[strings.LastIndex(repoName, "/")+1:]

	// Check to see if repo exists before downloading
	if dirExists(repoName) {
		fmt.Printf("%s already exists\n", repoName)
		return nil
	}

	fmt.Printf("running 'git clone %s %s'\n", repo, repoName)
	cmd := exec.Command("git", "clone", repo, repoName)

	fmt.Printf("Cloning repo: %s\n", repo)

	return cmd.Run()
}
