/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/aniruddha-sinha/gostart/internal"
	"github.com/spf13/cobra"
)

var (
	baseDevDir  string
	projectName string
	skipMise    bool
	moduleName  string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a project",
}

var goCmd = &cobra.Command{
	Use:          "golang",
	Aliases:      []string{"go"},
	Short:        "add golang project",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		umbrellaCfg := internal.UmbrellaConfig{
			BaseDir:              baseDevDir,
			ProjectName:          projectName,
			SkipMise:             skipMise,
			FileSystemOperations: internal.OSOperations{},
			OSExecutions:         internal.OSExecutions{},
			GoModuleName:         moduleName,
		}

		err := umbrellaCfg.OrchestrateGoProjectCreation()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(goCmd)

	goCmd.Flags().StringVarP(&baseDevDir, "base-dir", "d", "", "The base dev directory where you want to create the go project")
	goCmd.Flags().StringVarP(&projectName, "proj", "p", "", "The name of the Golang project")
	goCmd.Flags().BoolVarP(&skipMise, "skip-mise", "s", false, "Skip mise configuration for the current project; do it only if you have system levels tools installed")
	goCmd.Flags().StringVarP(&moduleName, "module-name", "m", "", "The name of the go module")
}
