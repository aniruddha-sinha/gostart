package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	embed "github.com/aniruddha-sinha/gostart/config"
	"github.com/aniruddha-sinha/gostart/config/userconfig"
)

type OSOpAbstractions interface {
	Stat(dir string) (os.FileInfo, error)
	Mkdir(filePath string, permissions os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type OSExecAbstraction interface {
	CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd
	Look(file string) (string, error)
}

type (
	OSOperations struct{}
	OSExecutions struct{}
)

type UmbrellaConfig struct {
	BaseDir              string
	ProjectName          string
	SkipMise             bool
	FileSystemOperations OSOperations
	OSExecutions         OSExecutions
	GoModuleName         string
}

func (ops OSOperations) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (ops OSOperations) Mkdir(filePath string, permissions os.FileMode) error {
	return os.Mkdir(filePath, permissions)
}

func (ops OSOperations) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (oex OSExecutions) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	//nolint:gosec // This is a generic wrapper for testing; inputs are trusted CLI arguments
	return exec.CommandContext(ctx, name, arg...)
}

func (oex OSExecutions) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (u UmbrellaConfig) OrchestrateGoProjectCreation() error {
	slog.Info("Base Dev Dir ", "dirpath ", u.BaseDir)
	slog.Info("Project Name ", "name ", u.ProjectName)
	slog.Info("Use Mise ", "= ", u.SkipMise)

	slog.Info("validating Dependencies")
	slog.Info("Running pre-flight dependency checks...")
	if err := u.validateDependencies(); err != nil {
		return err
	}

	slog.Info("Validating the base Directory...")
	if err := u.validateBaseDir(); err != nil {
		return fmt.Errorf("directory validation encountered an error %v", err)
	}

	slog.Info("Creating Project Directory")
	if err := u.createProjectDirectory(); err != nil {
		return fmt.Errorf("problems encountered when initialising go project %v", err)
	}

	slog.Info("Mise")
	if !u.SkipMise {
		slog.Info("Configuring mise env")
		if err := u.configureMise(); err != nil {
			return fmt.Errorf("problems encountered while configuring mise %v", err)
		}

		slog.Info("Run mise trust")
		output, err := u.execMiseTrust()
		if err != nil {
			return fmt.Errorf("problems encountered while running mise trust %v", err)
		}

		slog.Info("mise trust ", "output ", string(output))

		slog.Info("initialise Go Module")
		outputGoInit, err1 := u.initialiseGoModule()
		if err != nil {
			return fmt.Errorf("failed to initialise Go module %v", err1)
		} else {
			slog.Info("Go mod init ", "output ", string(outputGoInit))
		}

		slog.Info("Run cobra-cli init")
		cobraCliOut, err2 := u.cobraCLIInitialise()
		if err1 != nil {
			return fmt.Errorf("problems encountered while running cobra-cli init %v", err2)
		}

		slog.Info("cobra-cli init ", "output ", string(cobraCliOut))
	} else {
		slog.Info("Skipping Mise Configuration; make sure you have golang installed; for now skipping go mod init; just exiting with dir creation")
	}

	slog.Info("Go Project Created")

	return nil
}

func (u UmbrellaConfig) validateDependencies() error {
	// exec.LookPath searches the system $PATH for the executable
	if _, err := u.OSExecutions.LookPath("go"); err != nil {
		return fmt.Errorf("the 'go' command is required but was not found in your $PATH. Please ensure Go is installed")
	}

	// If you rely on mise to be installed globally for the user, check for it too!
	if !u.SkipMise {
		if _, err := u.OSExecutions.LookPath("mise"); err != nil {
			return fmt.Errorf("the 'mise' command was not found. Install it, or run with --skip-mise")
		}
	}

	return nil
}

func (u UmbrellaConfig) validateBaseDir() error {
	// check if base dir exists
	info, err := u.FileSystemOperations.Stat(u.BaseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("base dir does not exist %s, %v", u.BaseDir, err)
		} else {
			return fmt.Errorf("error stat %s, %v", u.BaseDir, err)
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory ", u.BaseDir)
	}

	return nil
}

func (u UmbrellaConfig) createProjectDirectory() error {
	config, err := userconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("error reading config %v", err)
	}

	defaultFilePermission, err := stringToFileMode(config.DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("error conversion from DEFAULT_FILE_PERMISSION string to os.FileMode %v", err)
	}

	projectPath := filepath.Join(u.BaseDir, u.ProjectName)
	if err := u.FileSystemOperations.Mkdir(projectPath, defaultFilePermission); err != nil {
		return fmt.Errorf("error creating directory %v", err)
	}
	return nil
}

func (u *UmbrellaConfig) initialiseGoModule() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := u.OSExecutions.CommandContext(ctx, "go", "mod", "init", u.GoModuleName)

	projectPath := filepath.Join(u.BaseDir, u.ProjectName)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// If the error was caused by the timeout, we can return a very specific error message
		if ctx.Err() == context.DeadlineExceeded {
			return []byte{}, fmt.Errorf("go mod init timed out after 30 seconds: %w", err)
		}
		return []byte{}, fmt.Errorf("failed to initialize go module: %w\nDetails: %s", err, string(output))
	}

	return output, nil
}

func (u UmbrellaConfig) configureMise() error {
	config, err := userconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("error reading config %v", err)
	}

	readyMadeFileBytes, err := embed.PreconfiguredMise.ReadFile(config.MiseGoPath)
	if err != nil {
		return fmt.Errorf("failed to read ready made mise file: %w", err)
	}

	projectPath := filepath.Join(u.BaseDir, u.ProjectName)
	targetPath := filepath.Join(projectPath, ".mise.toml")

	defaultFilePermission, err := stringToFileMode(config.DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("error conversion from DEFAULT_FILE_PERMISSION string to os.FileMode %v", err)
	}

	if err := u.FileSystemOperations.WriteFile(targetPath, readyMadeFileBytes, defaultFilePermission); err != nil {
		return fmt.Errorf("failed to write .mise.toml to target directory: %w", err)
	}

	return nil
}

func (u UmbrellaConfig) execMiseTrust() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := u.OSExecutions.CommandContext(ctx, "mise", "trust")

	projectPath := filepath.Join(u.BaseDir, u.ProjectName)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// If the error was caused by the timeout
		if ctx.Err() == context.DeadlineExceeded {
			return []byte{}, fmt.Errorf("go mod init timed out after 30 seconds: %w", err)
		}
		return []byte{}, fmt.Errorf("failed to run \"mise trust\": %w\nDetails: %s", err, string(output))
	}

	return output, nil
}

func (u UmbrellaConfig) cobraCLIInitialise() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := u.OSExecutions.CommandContext(ctx, "cobra-cli", "init")

	projectPath := filepath.Join(u.BaseDir, u.ProjectName)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// If the error was caused by the timeout
		if ctx.Err() == context.DeadlineExceeded {
			return []byte{}, fmt.Errorf("go mod init timed out after 30 seconds: %w", err)
		}
		return []byte{}, fmt.Errorf("failed to run \"cobra-cli init\": %w\nDetails: %s", err, string(output))
	}

	return output, nil
}
