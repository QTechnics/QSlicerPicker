# Contributing to 3D Slicer Picker

Thank you for your interest in contributing to 3D Slicer Picker! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check the [existing issues](https://github.com/QTechnics/QSlicerPicker/issues) to ensure the bug hasn't already been reported.

When creating a bug report, please include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior**
- **Actual behavior**
- **Screenshots** (if applicable)
- **Environment details**:
  - Operating system and version
  - Application version
  - Go version (if building from source)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- **Clear title and description**
- **Use case**: Why is this feature useful?
- **Proposed solution** (if you have one)
- **Alternatives considered** (if any)

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the coding standards
3. **Test your changes** thoroughly
4. **Update documentation** if needed
5. **Commit your changes** with clear, descriptive messages
6. **Push to your fork** and submit a pull request

#### Pull Request Guidelines

- **Keep PRs focused**: One feature or bug fix per PR
- **Write clear commit messages**: Follow [Conventional Commits](https://www.conventionalcommits.org/)
- **Update tests**: Add tests for new features
- **Update documentation**: Update README or other docs as needed
- **Link related issues**: Reference any related issues in your PR description

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- Platform-specific build tools (see [README.md](README.md#building-from-source))

### Getting Started

1. **Fork and clone**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/QSlicerPicker.git
   cd QSlicerPicker
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Build the project**:
   ```bash
   go build -o qslicerpicker .
   ```

4. **Run tests**:
   ```bash
   go test ./...
   ```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format code
- Run `golint` or `golangci-lint` before committing
- Keep functions small and focused
- Use meaningful variable and function names

### Code Organization

- **Package structure**: Follow the existing package structure
- **Error handling**: Always handle errors explicitly
- **Comments**: Add comments for exported functions and types
- **Imports**: Organize imports (stdlib, third-party, local)

### Example

```go
// ProcessFile handles file processing with error handling
func ProcessFile(path string) error {
    if path == "" {
        return fmt.Errorf("path cannot be empty")
    }
    
    // Process the file
    // ...
    
    return nil
}
```

## Adding Translations

The application uses Go's `embed` feature to include translation files directly in the binary, ensuring they work on all platforms without requiring external files.

### Adding a New Language

1. **Create translation file**: Create a new JSON file in `internal/i18n/locales/` (e.g., `es.json` for Spanish)
2. **Copy structure**: Copy the structure from `en.json` and translate all strings
3. **Add embed directive**: In `internal/i18n/i18n.go`, add:
   ```go
   //go:embed locales/es.json
   var esJSON []byte
   ```
4. **Update loadTranslations()**: Add the new language to the `langData` map:
   ```go
   langData := map[string][]byte{
       "tr": trJSON,
       "en": enJSON,
       "de": deJSON,
       "fr": frJSON,
       "es": esJSON,  // Add new language
   }
   ```
5. **Update language list**: Add the language code to `GetAvailableLanguages()`:
   ```go
   return []string{"tr", "en", "de", "fr", "es"}
   ```
6. **Update UI**: Add the language to the language selection in `internal/ui/settings.go`:
   - Add to `langRadio` options
   - Add to `langMap` for selection handling
   - Add to `langNames` for display

### Translation File Structure

Example (`es.json`):
```json
{
  "app_title": "3D Slicer Picker",
  "open_in": "Abrir en...",
  "choose_slicer": "Elegir un slicer",
  "cancel": "Cancelar",
  "open": "Abrir",
  "settings": "Configuración",
  "slicers": "Slicers",
  "enabled": "Habilitado",
  "disabled": "Deshabilitado",
  "path": "Ruta",
  "custom_path": "Ruta personalizada",
  "arguments": "Argumentos",
  "working_directory": "Directorio de trabajo",
  "add_custom_slicer": "Agregar Slicer Personalizado",
  "name": "Nombre",
  "save": "Guardar",
  "delete": "Eliminar",
  "move_up": "Mover arriba",
  "move_down": "Mover abajo",
  "language": "Idioma",
  "turkish": "Turco",
  "english": "Inglés",
  "german": "Alemán",
  "french": "Francés",
  "file_not_found": "Archivo no encontrado",
  "slicer_not_found": "Slicer no encontrado",
  "error_launching": "Error al iniciar el slicer",
  "no_slicers_available": "No hay slicers disponibles",
  "about": "Acerca de",
  "version": "Versión",
  "license": "Licencia",
  "author": "Autor",
  "source_code": "Código fuente"
}
```

### Important Notes

- **All keys must be present**: Every translation file must contain all keys from `en.json`
- **Embed directive path**: The `//go:embed` path is relative to the file location
- **Testing**: After adding a language, test it thoroughly in the UI
- **Fallback**: If a translation is missing, the system falls back to English, then to the key name

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/config
```

### Writing Tests

- Write tests for new features
- Aim for good test coverage
- Use table-driven tests when appropriate
- Test both success and error cases

Example:
```go
func TestProcessFile(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        wantErr bool
    }{
        {"valid path", "/path/to/file.stl", false},
        {"empty path", "", true},
        {"invalid path", "/nonexistent/file.stl", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ProcessFile(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

Examples:
```
feat: Add support for new slicer type
fix: Resolve crash when path is empty
docs: Update installation instructions
refactor: Simplify slicer detection logic
```

## Review Process

1. **Automated checks**: All PRs must pass CI checks
2. **Code review**: At least one maintainer will review your PR
3. **Feedback**: Address any feedback or requested changes
4. **Merge**: Once approved, your PR will be merged

## Questions?

- Open a [discussion](https://github.com/QTechnics/QSlicerPicker/discussions)
- Check existing [issues](https://github.com/QTechnics/QSlicerPicker/issues)
- Contact maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to 3D Slicer Picker! 🎉
