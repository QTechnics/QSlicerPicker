//go:build darwin
// +build darwin

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RegisterFileAssociations registers file associations on macOS
func RegisterFileAssociations() error {
	extensions := []string{"3mf", "step", "stl", "svg", "obj", "amf", "usd", "usda", "usdc", "abc", "ply", "sla"}

	appPath, err := getAppPath()
	if err != nil {
		return fmt.Errorf("failed to get app path: %w", err)
	}

	for _, ext := range extensions {
		// Use Launch Services to register file association
		// This requires the app to be in an .app bundle
		cmd := exec.Command("duti", "-s", appPath, "."+ext, "all")
		if err := cmd.Run(); err != nil {
			// duti might not be installed, that's okay
			continue
		}
	}

	return nil
}

func getAppPath() (string, error) {
	// Try to get the executable path
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	// If it's already an .app bundle, return it
	if filepath.Ext(execPath) == ".app" {
		return execPath, nil
	}

	// Otherwise, we'd need to create an .app bundle
	// For now, return the executable path
	return execPath, nil
}
