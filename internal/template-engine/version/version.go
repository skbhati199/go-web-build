package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
)

type TemplateVersion struct {
	Version  *semver.Version
	Path     string
	Metadata map[string]interface{}
}

type VersionManager struct {
	templatesDir string
}

func NewTemplateVersion(version string) (*TemplateVersion, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}
	return &TemplateVersion{Version: v}, nil
}

func (tv *TemplateVersion) String() string {
	if tv.Version != nil {
		return tv.Version.String()
	}
	return ""
}

func NewVersionManager(templatesDir string) *VersionManager {
	return &VersionManager{
		templatesDir: templatesDir,
	}
}

func (vm *VersionManager) GetAllVersions(templateName string) ([]*TemplateVersion, error) {
	versionsPath := filepath.Join(vm.templatesDir, templateName, "versions.json")
	data, err := os.ReadFile(versionsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions file: %w", err)
	}

	var versionData struct {
		Versions []struct {
			Version     string                 `json:"version"`
			Path        string                 `json:"path"`
			ReleaseDate string                 `json:"releaseDate"`
			IsLatest    bool                   `json:"isLatest"`
			Features    []string               `json:"features"`
			Metadata    map[string]interface{} `json:"metadata"`
		} `json:"versions"`
		Latest string `json:"latest"`
		Stable string `json:"stable"`
	}

	if err := json.Unmarshal(data, &versionData); err != nil {
		return nil, fmt.Errorf("failed to parse versions: %w", err)
	}

	var versions []*TemplateVersion
	for _, v := range versionData.Versions {
		version, err := NewTemplateVersion(v.Version)
		if err != nil {
			return nil, err
		}
		version.Path = v.Path
		version.Metadata = v.Metadata
		versions = append(versions, version)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version.LessThan(versions[j].Version)
	})

	return versions, nil
}

func (vm *VersionManager) GetVersion(templateName, version string) (*TemplateVersion, error) {
	versions, err := vm.GetAllVersions(templateName)
	if err != nil {
		return nil, err
	}

	constraint, err := semver.NewConstraint(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint: %w", err)
	}

	for _, v := range versions {
		if constraint.Check(v.Version) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("version %s not found for template %s", version, templateName)
}

func (vm *VersionManager) GetLatestVersion(templateName string) (*TemplateVersion, error) {
	versions, err := vm.GetAllVersions(templateName)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for template %s", templateName)
	}

	return versions[len(versions)-1], nil
}

func (vm *VersionManager) IsValidVersion(templateName, version string) bool {
	versions, err := vm.GetAllVersions(templateName)
	if err != nil {
		return false
	}

	constraint, err := semver.NewConstraint(version)
	if err != nil {
		return false
	}

	for _, v := range versions {
		if constraint.Check(v.Version) {
			return true
		}
	}
	return false
}
