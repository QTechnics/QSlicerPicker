package i18n

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"qslicerpicker/internal/config"
	"runtime"
)

//go:embed locales/tr.json
var trJSON []byte

//go:embed locales/en.json
var enJSON []byte

//go:embed locales/de.json
var deJSON []byte

//go:embed locales/fr.json
var frJSON []byte

var translations map[string]map[string]string
var currentLang string

func init() {
	translations = make(map[string]map[string]string)
	loadTranslations()
}

func loadTranslations() {
	// Load embedded translations first
	langData := map[string][]byte{
		"tr": trJSON,
		"en": enJSON,
		"de": deJSON,
		"fr": frJSON,
	}

	for lang, data := range langData {
		if len(data) == 0 {
			continue
		}
		var trans map[string]string
		if err := json.Unmarshal(data, &trans); err == nil {
			translations[lang] = trans
		}
	}

	// Try to load from file system (for development/overrides)
	// Get the directory where locales are stored
	var localesDir string

	// Try to get executable path first (works better on Windows)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// Check if locales directory exists relative to executable
		possibleLocalesDir := filepath.Join(exeDir, "locales")
		if _, err := os.Stat(possibleLocalesDir); err == nil {
			localesDir = possibleLocalesDir
		} else {
			// Try parent directory (for development)
			possibleLocalesDir = filepath.Join(filepath.Dir(exeDir), "internal", "i18n", "locales")
			if _, err := os.Stat(possibleLocalesDir); err == nil {
				localesDir = possibleLocalesDir
			}
		}
	}

	// Fallback to runtime.Caller (for development)
	if localesDir == "" {
		_, filename, _, _ := runtime.Caller(0)
		localesDir = filepath.Join(filepath.Dir(filename), "locales")
	}

	// Load from file system if available (overrides embedded)
	languages := []string{"tr", "en", "de", "fr"}
	for _, lang := range languages {
		filePath := filepath.Join(localesDir, lang+".json")
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var trans map[string]string
		if err := json.Unmarshal(data, &trans); err != nil {
			continue
		}

		// Override embedded translations with file-based ones
		translations[lang] = trans
	}

	// Set current language from config or system
	cfg := config.GetConfig()
	if cfg.Language != "" {
		currentLang = cfg.Language
	} else {
		currentLang = getSystemLanguage()
	}

	// Fallback to English if current language not available
	if _, ok := translations[currentLang]; !ok {
		currentLang = "en"
	}
}

func getSystemLanguage() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang == "" {
		return "en"
	}

	// Extract language code (e.g., "tr_TR.UTF-8" -> "tr")
	if len(lang) >= 2 {
		code := lang[:2]
		if code == "tr" || code == "en" || code == "de" || code == "fr" {
			return code
		}
	}

	return "en"
}

// T translates a key to the current language
func T(key string) string {
	if trans, ok := translations[currentLang]; ok {
		if val, ok := trans[key]; ok {
			return val
		}
	}

	// Fallback to English
	if trans, ok := translations["en"]; ok {
		if val, ok := trans[key]; ok {
			return val
		}
	}

	// Return key if translation not found
	return key
}

// SetLanguage sets the current language
func SetLanguage(lang string) {
	if _, ok := translations[lang]; ok {
		currentLang = lang
		cfg := config.GetConfig()
		cfg.Language = lang
		config.SaveConfig()
	}
}

// GetLanguage returns the current language
func GetLanguage() string {
	return currentLang
}

// GetAvailableLanguages returns all available language codes
func GetAvailableLanguages() []string {
	return []string{"tr", "en", "de", "fr"}
}
