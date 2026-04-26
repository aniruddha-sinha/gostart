package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aniruddha-sinha/gostart/config/userconfig"
)

type OSOpAbstractions interface {
	Stat(dir string) (os.FileInfo, error)
	Mkdir(filePath string, permissions os.FileMode) error
}

type OSOperations struct{}

type UmbrellaConfig struct {
	BaseDir              string
	ProjectName          string
	SkipMise             bool
	FileSystemOperations OSOperations
}

func (ops OSOperations) Stat(dir string) (os.FileInfo, error) {
	return os.Stat(dir)
}

func (ops OSOperations) Mkdir(filePath string, permissions os.FileMode) error {
	return os.Mkdir(filePath, permissions)
}

func (u UmbrellaConfig) OrchestrateGoProjectCreation() error {
	slog.Info("Base Dev Dir ", "dirpath ", u.BaseDir)
	slog.Info("Project Name ", "name ", u.ProjectName)
	slog.Info("Use Mise ", "= ", u.SkipMise)

	slog.Info("Validating the base Directory...")
	if err := u.validateBaseDir(); err != nil {
		return fmt.Errorf("directory validation encountered an error %v", err)
	}

	slog.Info("initializing go mod")
	if err := u.createProjectDirectory(); err != nil {
		return fmt.Errorf("problems encountered when initialising go project %v", err)
	}

	slog.Info("Mise")
	if !u.SkipMise {
		slog.Info("Configuring mise env")
		if err := u.configureMise(); err != nil {
			return fmt.Errorf("problems encountered while configuring mise %v", err)
		}
	} else {
		slog.Info("Skipping Mise Configuration")
	}

	slog.Info("Go Project Created")

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
	// here i will create a new project
	// then run go mod init
	projectPath := filepath.Join(u.BaseDir, u.ProjectName)
	if err := u.FileSystemOperations.Mkdir(projectPath, defaultFilePermission); err != nil {
		return fmt.Errorf("error creating directory %v", err)
	}
	return nil
}

func (u UmbrellaConfig) configureMise() error {
	return nil
}
