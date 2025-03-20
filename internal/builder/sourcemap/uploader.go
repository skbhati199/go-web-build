package sourcemap

import (
	"bytes"
	"fmt"
	"net/http"
)

type ErrorTracker interface {
	UploadSourceMap(filename string, sourceMap []byte) error
	ValidateSourceMap(sourceMap []byte) error
}

type SentryUploader struct {
	organization string
	project      string
	authToken    string
	release      string
}

func NewSentryUploader(org, project, token, release string) *SentryUploader {
	return &SentryUploader{
		organization: org,
		project:      project,
		authToken:    token,
		release:      release,
	}
}

func (s *SentryUploader) UploadSourceMap(filename string, sourceMap []byte) error {
	url := fmt.Sprintf("https://sentry.io/api/0/organizations/%s/releases/%s/files/",
		s.organization, s.release)

	req, err := http.NewRequest("POST", url, bytes.NewReader(sourceMap))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.authToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to upload source map: %s", resp.Status)
	}

	return nil
}

// Add Rollbar implementation
type RollbarUploader struct {
	projectID string
	authToken string
	release   string
}

func NewRollbarUploader(projectID, token, release string) *RollbarUploader {
	return &RollbarUploader{
		projectID: projectID,
		authToken: token,
		release:   release,
	}
}

func (r *RollbarUploader) UploadSourceMap(filename string, sourceMap []byte) error {
	url := fmt.Sprintf("https://api.rollbar.com/api/1/sourcemap")

	req, err := http.NewRequest("POST", url, bytes.NewReader(sourceMap))
	if err != nil {
		return err
	}

	req.Header.Set("X-Rollbar-Access-Token", r.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload source map to Rollbar: %s", resp.Status)
	}

	return nil
}

func (r *RollbarUploader) ValidateSourceMap(sourceMap []byte) error {
	// Implement source map validation for Rollbar
	return nil
}

func (s *SentryUploader) ValidateSourceMap(sourceMap []byte) error {
	url := fmt.Sprintf("https://sentry.io/api/0/organizations/%s/releases/%s/files/validate/",
		s.organization, s.release)

	req, err := http.NewRequest("POST", url, bytes.NewReader(sourceMap))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.authToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to validate source map: %s", resp.Status)
	}

	return nil
}
