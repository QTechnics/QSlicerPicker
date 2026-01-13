//go:build linux
// +build linux

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RegisterFileAssociations registers file associations on Linux
func RegisterFileAssociations() error {
	extensions := []string{"3mf", "step", "stl", "svg", "obj", "amf", "usd", "usda", "usdc", "abc", "ply", "sla"}

	appPath, err := getAppPath()
	if err != nil {
		return fmt.Errorf("failed to get app path: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	desktopDir := filepath.Join(homeDir, ".local", "share", "applications")
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		return fmt.Errorf("failed to create desktop directory: %w", err)
	}

	// Create .desktop file
	desktopFile := filepath.Join(desktopDir, "qslicerpicker.desktop")
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Name=3D Slicer Picker
Exec=%s %%f
Type=Application
MimeType=%s
`, appPath, strings.Join(getMimeTypes(extensions), ";"))

	if err := os.WriteFile(desktopFile, []byte(desktopContent), 0644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	// Update MIME database
	cmd := exec.Command("update-desktop-database", desktopDir)
	cmd.Run() // Ignore errors

	// Update MIME types
	for _, ext := range extensions {
		mimeType := getMimeType(ext)
		cmd := exec.Command("xdg-mime", "default", "qslicerpicker.desktop", mimeType)
		cmd.Run() // Ignore errors
	}

	return nil
}

func getMimeTypes(extensions []string) []string {
	mimeTypes := make([]string, 0, len(extensions))
	for _, ext := range extensions {
		mimeTypes = append(mimeTypes, getMimeType(ext))
	}
	return mimeTypes
}

func getMimeType(ext string) string {
	mimeMap := map[string]string{
		"3mf":  "model/3mf",
		"step": "application/step",
		"stl":  "model/stl",
		"svg":  "image/svg+xml",
		"obj":  "model/obj",
		"amf":  "application/x-amf",
		"usd":  "model/vnd.usd",
		"usda": "model/vnd.usd",
		"usdc": "model/vnd.usd",
		"abc":  "application/x-abc",
		"ply":  "model/ply",
		"sla":  "application/x-sla",
	}

	if mime, ok := mimeMap[ext]; ok {
		return mime
	}

	return "application/octet-stream"
}

func getAppPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
