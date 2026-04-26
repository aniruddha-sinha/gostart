# gostart

A lightweight CLI tool designed to rapidly scaffold new [Cobra](https://github.com/spf13/cobra)-based Go projects. `gostart` automates the boilerplate of setting up a modern Go application, complete with a modular directory structure and out-of-the-box support for [Mise](https://mise.jdx.dev/) task automation.

## Features

- **Rapid Scaffolding**: Instantly generates a clean project structure suitable for CLI applications.
- **Cobra Integration**: Pre-configures the foundation for building command-line interfaces.
- **Mise Native**: Bundles a `.mise.toml` file by default to manage environment variables and chain development tasks (like linting and building) without repeating yourself.
- **CI/CD & Tooling Ready**: Includes pre-configured files for GitHub Actions workflows, strict linting (`.golangci.yml`), and automated binary releases (`.goreleaser.yml`).
- **Automated Initialization**: Safely executes `go mod init` inside your newly created project directory.

## Installation

Ensure you have Go installed on your system. You can install `gostart` directly via `go install`:

```bash
go install [github.com/aniruddha-sinha/gostart@latest](https://github.com/aniruddha-sinha/gostart@latest)

## Usage

Use the `gostart` command to initialize a new project by providing the base directory and your desired project name.

### Flags
```
| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--base-dir` | `-d` | `""` | The base development directory where the project will be created. |
| `--proj` | `-p` | `""` | The name of the new Golang project. |
| `--skip-mise` | `-s` | `false` | Pass this flag to explicitly opt-out of generating the `.mise.toml` configuration. |
```
### Examples

To scaffold a new project called `my-cli` inside your `~/dev/go/projects` folder (with `mise` configured by default):

```bash
gostart -d ~/dev/go/projects -p my-cli