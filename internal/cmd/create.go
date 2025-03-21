package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	templateengine "github.com/skbhati199/go-web-build/internal/template-engine"
	"github.com/spf13/cobra"
)

var (
	framework    string
	templateName string
)

func init() {
	createCmd := &cobra.Command{
		Use:     "create [name]",
		Aliases: []string{"new", "init", "n"},
		Short:   "Create a new project",
		Long:    `Create a new web application project with the specified framework and template.`,
		Example: `  # Create a React TypeScript project
  gobuild create myapp --framework react --template typescript

  # Create a Vue.js project
  gobuild create myapp --framework vue --template composition`,
		RunE: runCreate,
	}

	// Add flags
	createCmd.Flags().StringVarP(&framework, "framework", "f", "", "Framework to use (react, vue, next)")
	createCmd.Flags().StringVarP(&templateName, "template", "t", "", "Template to use (javascript, typescript)")

	// Mark required flags
	createCmd.MarkFlagRequired("framework")

	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("project name is required")
	}

	projectName := args[0]

	// Enable debug output
	fmt.Printf("Creating project: %s\nFramework: %s\nTemplate: %s\n", projectName, framework, templateName)

	// Create template manager with proper path handling
	templatesDir := filepath.Join("/Users/sonukumar/go-web-build", "templates")

	// Try both directory structures
	templatePath := filepath.Join(templatesDir, fmt.Sprintf("%s-%s", framework, templateName))
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// Try nested structure
		templatePath = filepath.Join(templatesDir, framework, fmt.Sprintf("react-%s", templateName))
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("template not found: %s-%s", framework, templateName)
		}
	}

	// Create template manager
	manager := templateengine.NewManager(templatesDir)

	// Create project with specified template
	err := manager.CreateProject(projectName, framework, templateName, "", nil)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	fmt.Printf("Successfully created project %s\n", projectName)
	return nil
}
