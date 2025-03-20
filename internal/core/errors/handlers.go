package errors

import (
	"fmt"
	"os"
	"time"
)

func DefaultBuildRecoveryHandler(err interface{}) error {
	logError := fmt.Sprintf("[%s] Build process panic: %v\n", time.Now().Format(time.RFC3339), err)

	// Log the error
	if f, err := os.OpenFile("build.error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		defer f.Close()
		f.WriteString(logError)
	}

	// Cleanup temporary files
	cleanupBuildArtifacts()

	return nil
}

func DefaultTemplateRecoveryHandler(err interface{}) error {
	logError := fmt.Sprintf("[%s] Template generation panic: %v\n", time.Now().Format(time.RFC3339), err)

	// Log the error
	if f, err := os.OpenFile("template.error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		defer f.Close()
		f.WriteString(logError)
	}

	// Cleanup template cache
	cleanupTemplateCache()

	return nil
}

func cleanupBuildArtifacts() {
	// Implementation for cleaning up temporary build files
}

func cleanupTemplateCache() {
	// Implementation for cleaning up template cache
}
