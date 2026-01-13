//go:build windows
// +build windows

package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

var (
	shell32           = syscall.NewLazyDLL("shell32.dll")
	procSHChangeNotify = shell32.NewProc("SHChangeNotify")
)

// RegisterFileAssociations registers file associations on Windows
func RegisterFileAssociations() error {
	extensions := []string{"3mf", "step", "stl", "svg", "obj", "amf", "usd", "usda", "usdc", "abc", "ply", "sla"}

	appPath, err := getAppPath()
	if err != nil {
		return fmt.Errorf("failed to get app path: %w", err)
	}

	appName := "QSlicerPicker"

	for _, ext := range extensions {
		// Create registry entries for file association
		keyPath := fmt.Sprintf(`Software\Classes\.%s`, ext)
		key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
		if err != nil {
			continue
		}
		defer key.Close()

		progID := fmt.Sprintf("%s.%s", appName, ext)
		key.SetStringValue("", progID)

		// Create ProgID
		progIDPath := fmt.Sprintf(`Software\Classes\%s`, progID)
		progIDKey, _, err := registry.CreateKey(registry.CURRENT_USER, progIDPath, registry.ALL_ACCESS)
		if err != nil {
			continue
		}
		defer progIDKey.Close()

		progIDKey.SetStringValue("", fmt.Sprintf("%s File", ext))

		// Create shell\open\command
		commandPath := fmt.Sprintf(`Software\Classes\%s\shell\open\command`, progID)
		commandKey, _, err := registry.CreateKey(registry.CURRENT_USER, commandPath, registry.ALL_ACCESS)
		if err != nil {
			continue
		}
		defer commandKey.Close()

		commandKey.SetStringValue("", fmt.Sprintf(`"%s" "%%1"`, appPath))
	}

	// Notify shell of changes
	procSHChangeNotify.Call(
		uintptr(0x8000000), // SHCNE_ASSOCCHANGED
		uintptr(0x00000001), // SHCNF_IDLIST
		0,
		0,
	)

	return nil
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
