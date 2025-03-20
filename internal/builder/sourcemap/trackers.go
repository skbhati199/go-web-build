package sourcemap

func createErrorTracker(config ErrorTrackingConfig) ErrorTracker {
	switch config.Provider {
	case "sentry":
		return NewSentryUploader(
			config.ProjectID,
			"", // project name
			config.AuthToken,
			"latest", // release version
		)
	case "rollbar":
		// Implement Rollbar uploader
		return nil
	default:
		return nil
	}
}
