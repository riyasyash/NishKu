package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/riyasyash/nishku/internal/window"
)

type Profile struct {
	Name      string           `json:"name"`
	Windows   []window.Window  `json:"windows"`
	Displays  []window.Display `json:"displays"`  // Display configuration when saved
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// getProfileDir returns the directory where profiles are stored
func getProfileDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	profileDir := filepath.Join(home, ".nishku", "profiles")
	
	// Create directory if it doesn't exist with restrictive permissions (0700 = owner only)
	if err := os.MkdirAll(profileDir, 0700); err != nil {
		return "", err
	}

	return profileDir, nil
}

// getProfilePath returns the full path to a profile file
func getProfilePath(name string) (string, error) {
	dir, err := getProfileDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}

// SaveProfile saves a profile to disk
func SaveProfile(prof Profile) error {
	path, err := getProfilePath(prof.Name)
	if err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	prof.UpdatedAt = now
	if prof.CreatedAt.IsZero() {
		prof.CreatedAt = now
	}

	data, err := json.MarshalIndent(prof, "", "  ")
	if err != nil {
		return err
	}

	// Write with restrictive permissions (0600 = owner read/write only)
	return os.WriteFile(path, data, 0600)
}

// LoadProfile loads a profile from disk
func LoadProfile(name string) (*Profile, error) {
	path, err := getProfilePath(name)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile '%s' not found", name)
		}
		return nil, err
	}

	var prof Profile
	if err := json.Unmarshal(data, &prof); err != nil {
		return nil, err
	}

	return &prof, nil
}

// ListProfiles returns all saved profiles
func ListProfiles() ([]Profile, error) {
	dir, err := getProfileDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()[:len(entry.Name())-5] // Remove .json extension
		prof, err := LoadProfile(name)
		if err != nil {
			continue // Skip invalid profiles
		}

		profiles = append(profiles, *prof)
	}

	return profiles, nil
}

// DeleteProfile deletes a profile from disk
func DeleteProfile(name string) error {
	path, err := getProfilePath(name)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("profile '%s' not found", name)
		}
		return err
	}

	return nil
}
