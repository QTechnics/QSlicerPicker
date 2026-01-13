package slicer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"qslicerpicker/internal/config"
	"runtime"
)

type Slicer struct {
	ID          string
	Name        string
	DefaultPath string
	Path        string
	Enabled     bool
	Order       int
	Arguments   []string
	WorkingDir  string
	IsCustom    bool
}

var defaultSlicers = []struct {
	ID          string
	Name        string
	DefaultPath map[string]string // platform -> path
}{
	{"cura", "Cura", map[string]string{
		"darwin":  "/Applications/Ultimaker Cura.app/Contents/MacOS/Ultimaker Cura",
		"windows": "C:\\Program Files\\Ultimaker Cura\\Ultimaker-Cura.exe",
		"linux":   "/usr/bin/cura",
	}},
	{"prusaslicer", "PrusaSlicer", map[string]string{
		"darwin":  "/Applications/PrusaSlicer.app/Contents/MacOS/PrusaSlicer",
		"windows": "C:\\Program Files\\PrusaSlicer\\prusa-slicer.exe",
		"linux":   "/usr/bin/prusa-slicer",
	}},
	{"superslicer", "SuperSlicer", map[string]string{
		"darwin":  "/Applications/SuperSlicer.app/Contents/MacOS/SuperSlicer",
		"windows": "C:\\Program Files\\SuperSlicer\\super-slicer.exe",
		"linux":   "/usr/bin/super-slicer",
	}},
	{"orcaslicer", "OrcaSlicer", map[string]string{
		"darwin":  "/Applications/OrcaSlicer.app/Contents/MacOS/OrcaSlicer",
		"windows": "C:\\Program Files\\OrcaSlicer\\OrcaSlicer.exe",
		"linux":   "/usr/bin/orca-slicer",
	}},
	{"bambustudio", "Bambu Studio", map[string]string{
		"darwin":  "/Applications/BambuStudio.app/Contents/MacOS/BambuStudio",
		"windows": "C:\\Program Files\\BambuStudio\\BambuStudio.exe",
		"linux":   "/usr/bin/bambu-studio",
	}},
	{"slic3r", "Slic3r", map[string]string{
		"darwin":  "/Applications/Slic3r.app/Contents/MacOS/Slic3r",
		"windows": "C:\\Program Files\\Slic3r\\slic3r.exe",
		"linux":   "/usr/bin/slic3r",
	}},
	{"ideamaker", "IdeaMaker", map[string]string{
		"darwin":  "/Applications/IdeaMaker.app/Contents/MacOS/IdeaMaker",
		"windows": "C:\\Program Files\\Raise3D\\IdeaMaker\\IdeaMaker.exe",
		"linux":   "/usr/bin/ideamaker",
	}},
	{"simplify3d", "Simplify3D", map[string]string{
		"darwin":  "/Applications/Simplify3D.app/Contents/MacOS/Simplify3D",
		"windows": "C:\\Program Files\\Simplify3D\\Simplify3D.exe",
		"linux":   "/usr/bin/simplify3d",
	}},
	{"kisslicer", "KISSlicer", map[string]string{
		"darwin":  "/Applications/KISSlicer.app/Contents/MacOS/KISSlicer",
		"windows": "C:\\Program Files\\KISSlicer\\KISSlicer.exe",
		"linux":   "/usr/bin/kisslicer",
	}},
	{"slic3rpe", "Slic3r PE", map[string]string{
		"darwin":  "/Applications/Slic3r PE.app/Contents/MacOS/Slic3r PE",
		"windows": "C:\\Program Files\\Slic3r PE\\slic3r-pe.exe",
		"linux":   "/usr/bin/slic3r-pe",
	}},
}

// GetDefaultSlicers returns all default slicer definitions
func GetDefaultSlicers() []Slicer {
	platform := runtime.GOOS
	slicers := make([]Slicer, 0, len(defaultSlicers))

	for i, ds := range defaultSlicers {
		defaultPath := ds.DefaultPath[platform]
		slicers = append(slicers, Slicer{
			ID:          ds.ID,
			Name:        ds.Name,
			DefaultPath: defaultPath,
			Path:        defaultPath,
			Enabled:     true,
			Order:       i * 10, // Default order with spacing for reordering
			IsCustom:    false,
		})
	}

	return slicers
}

// LoadSlicers loads slicers from config and merges with defaults
func LoadSlicers() []Slicer {
	cfg := config.GetConfig()
	defaultSlicers := GetDefaultSlicers()

	// If config is empty, initialize with defaults
	if len(cfg.Slicers) == 0 {
		cfg.Slicers = make([]config.SlicerConfig, 0, len(defaultSlicers))
		for i, ds := range defaultSlicers {
			cfg.Slicers = append(cfg.Slicers, config.SlicerConfig{
				ID:      ds.ID,
				Enabled: true,
				Order:   i * 10,
			})
		}
		config.SaveConfig()
	}

	// Create a map of configured slicers
	configMap := make(map[string]config.SlicerConfig)
	for _, sc := range cfg.Slicers {
		configMap[sc.ID] = sc
	}

	// Merge defaults with config
	slicers := make([]Slicer, 0)
	for _, ds := range defaultSlicers {
		slicer := ds
		if sc, ok := configMap[ds.ID]; ok {
			slicer.Enabled = sc.Enabled
			slicer.Order = sc.Order
			if sc.CustomPath != "" {
				slicer.Path = sc.CustomPath
			}
			slicer.Arguments = sc.Arguments
			slicer.WorkingDir = sc.WorkingDir
		}
		slicers = append(slicers, slicer)
	}

	// Add custom slicers
	for i, cs := range cfg.CustomSlicers {
		slicers = append(slicers, Slicer{
			ID:         fmt.Sprintf("custom_%d", i),
			Name:       cs.Name,
			Path:       cs.Path,
			Enabled:    cs.Enabled,
			Order:      cs.Order,
			Arguments:  cs.Arguments,
			WorkingDir: cs.WorkingDir,
			IsCustom:   true,
		})
	}

	// Sort all slicers by order
	for i := 0; i < len(slicers)-1; i++ {
		for j := i + 1; j < len(slicers); j++ {
			if slicers[i].Order > slicers[j].Order {
				slicers[i], slicers[j] = slicers[j], slicers[i]
			}
		}
	}

	return slicers
}

// GetEnabledSlicers returns only enabled slicers, sorted by order
func GetEnabledSlicers() []Slicer {
	allSlicers := LoadSlicers()
	enabled := make([]Slicer, 0)

	for _, s := range allSlicers {
		if s.Enabled && s.Path != "" {
			// Check if path exists
			if _, err := os.Stat(s.Path); err == nil {
				enabled = append(enabled, s)
			}
		}
	}

	// Sort by order
	for i := 0; i < len(enabled)-1; i++ {
		for j := i + 1; j < len(enabled); j++ {
			if enabled[i].Order > enabled[j].Order {
				enabled[i], enabled[j] = enabled[j], enabled[i]
			}
		}
	}

	return enabled
}

// LaunchSlicer launches a slicer with the given file
func LaunchSlicer(slicer Slicer, filePath string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "darwin" {
		// Check if path is a .app bundle
		if filepath.Ext(slicer.Path) == ".app" || (filepath.Dir(slicer.Path) != "" && filepath.Ext(filepath.Dir(slicer.Path)) == ".app") {
			// macOS .app bundle - find the .app path
			appPath := slicer.Path
			if filepath.Ext(slicer.Path) != ".app" {
				// Path is inside .app bundle, find the .app
				dir := filepath.Dir(slicer.Path)
				for dir != "/" && dir != "." {
					if filepath.Ext(dir) == ".app" {
						appPath = dir
						break
					}
					dir = filepath.Dir(dir)
				}
			}
			// Use open command for .app bundles
			args := []string{"-a", appPath, filePath}
			cmd = exec.Command("open", args...)
		} else {
			// Regular executable
			args := append(slicer.Arguments, filePath)
			cmd = exec.Command(slicer.Path, args...)
		}
	} else {
		// Windows and Linux
		args := append(slicer.Arguments, filePath)
		cmd = exec.Command(slicer.Path, args...)
	}

	if slicer.WorkingDir != "" {
		cmd.Dir = slicer.WorkingDir
	}

	return cmd.Start()
}

// FindSlicerByID finds a slicer by its ID
func FindSlicerByID(id string) *Slicer {
	slicers := LoadSlicers()
	for _, s := range slicers {
		if s.ID == id {
			return &s
		}
	}
	return nil
}
