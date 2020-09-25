## Workspace Repository Manager

This is a basic utility used for managing related repositories in a workspace, specifically for Go projects. Using
modules and a simple configuration file you can create a workspace and manage `go.mod` `replace` directives, easing
development with related projects.

### Workspace Configuration

A configuration file for a single workspace is a TOML file with a list of repositories. The default configuration file
name is `workspace.toml`, and a different name can be set by using the `--config` flag.

```toml
repos = [
    "git@github.com:org/repo-one.git",
    "git@github.com:org/repo-two.git",
]
```

From the directory containing that file, run `rman init`, and it will clone the repositories. It will then look through
the `go.mod` files for any `require` dependencies that mention a sibling repository _(defined by the module name from
`go.mod`)_, then add `replace` directives to point to the sibling directory.

### Global Configuration

Configurations can be stored in `~/.config/rman/workspaces.toml` as a map, like so:

```toml
[foo]
repos = [
    # repositories here    
]

[bar]
repos = [
    # different repositories
]
```

When `rman init --config foo` is run it will first look for a file named `foo`, and if it does not exist it will try to
load the global configuration and return the list of repositories named `foo`.

**NB:** I was debating between a single global configuration file or using `~/.config/rman/{name}.toml`, and I'm _still_
undecided. /shrug

### Standalone usage

This utility can also be used standalone. Pointing it at a directory containing several projects it will look for
dependencies that live in sibling directories, then create a `replace` directive for sibling projects. Usage is:
`rman wire /path/to/workspace`.
