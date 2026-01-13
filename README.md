# 3D Slicer Picker

<div align="center">

![Logo](q.png)

**Cross-platform 3D printer slicer selector application**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux-lightgrey)](https://github.com/QTechnics/QSlicerPicker)

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Contributing](#-contributing) • [License](#-license)

</div>

---

## 📖 About

**3D Slicer Picker** is an open-source, cross-platform application that allows you to choose which 3D printer slicer software to use when opening 3D model files. It provides a native dialog interface to select from multiple installed slicers.

When you double-click a 3D model file (`.3mf`, `.stl`, `.step`, `.obj`, etc.), **3D Slicer Picker** displays a beautiful, system-native dialog where you can choose which slicer to open the file with. Perfect for users who work with multiple slicers and want a unified way to manage their workflow.

### 🎯 Why 3D Slicer Picker?

- **No more default slicer conflicts**: Choose the right slicer for each project
- **Unified workflow**: One application to manage all your slicers
- **Native experience**: System-native dialogs that respect your OS theme (light/dark mode)
- **Multi-language support**: Available in Turkish, English, German, and French
- **Lightweight**: Single binary executable, minimal disk space usage
- **Customizable**: Add unlimited custom slicers with custom paths and arguments

## ✨ Features

- 🖥️ **Cross-platform**: Works on macOS, Windows, and Linux
- 🎨 **Native UI**: System-native dialogs with automatic light/dark mode support
- 🌍 **Multi-language**: Turkish, English, German, and French support
- 📦 **Single binary**: Easy to install and run
- ⚙️ **Configurable**: JSON-based configuration with easy-to-use settings UI
- 🔧 **Custom slicers**: Add unlimited custom slicers with:
  - Custom executable paths
  - Command-line arguments
  - Working directories
- 📋 **Slicer management**: 
  - Enable/disable slicers
  - Reorder slicers
  - Edit slicer configurations
- 🎯 **Smart detection**: Automatically detects common slicer installations
- 📁 **File associations**: Easy setup for supported file types

## 🎮 Supported Slicers

### Default Slicers (Auto-detected)

- **Cura** - Ultimaker Cura
- **PrusaSlicer** - Prusa Research's slicer
- **SuperSlicer** - Advanced PrusaSlicer fork
- **OrcaSlicer** - Modern, feature-rich slicer
- **Bambu Studio** - Bambu Lab's slicer
- **Slic3r** - Open-source slicer
- **IdeaMaker** - Raise3D's slicer
- **Simplify3D** - Professional slicer
- **KISSlicer** - Fast, efficient slicer
- **Slic3r PE** - Prusa Edition

### Custom Slicers

You can add any slicer application with custom paths, command-line arguments, and working directories.

## 📥 Installation

### Pre-built Releases

Download the latest release from the [Releases page](https://github.com/QTechnics/QSlicerPicker/releases):

- **macOS**: Download the `.zip` file for your architecture (Intel or Apple Silicon), extract the `.app` bundle, and move it to your Applications folder.
  
  **Note**: On first launch, macOS may show a security warning. If you encounter this:
  1. Right-click the `.app` bundle → Open (first time only)
  2. Or run in Terminal: `sudo xattr -d com.apple.quarantine /path/to/QSlicerPicker-*.app`
- **Linux**: Download the binary for your architecture, make it executable (`chmod +x qslicerpicker-linux-*`), and run it.
- **Windows**: Download the `.exe` file and run it.

### Building from Source

#### Prerequisites

- **Go** 1.21 or later ([Download](https://golang.org/dl/))
- **CGO** enabled (required for Fyne)
- Platform-specific dependencies:
  - **macOS**: Xcode Command Line Tools
  - **Linux**: `libgl1-mesa-dev`, `xorg-dev`
  - **Windows**: MinGW or MSYS2 (for CGO)

#### Build Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/QTechnics/QSlicerPicker.git
   cd QSlicerPicker
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Build for your platform**:
   ```bash
   # Native build (recommended)
   cd build
   ./build-native.sh
   
   # Or build directly
   go build -o qslicerpicker .
   
   # For macOS .app bundle
   go install fyne.io/fyne/v2/cmd/fyne@latest
   fyne package -os darwin -name "3D Slicer Picker" -appID "com.qslicerpicker.app" -icon q.png
   ```

## 📋 System Requirements

- **macOS**: 10.13 or later
- **Windows**: Windows 10 or later with OpenGL 2.1+ support
  - **Note**: Some systems (especially VMs or systems with outdated graphics drivers) may not support OpenGL. Ensure your graphics drivers are up to date.
- **Linux**: X11 with OpenGL 2.1+ support

## 🚀 Usage

### First Run

1. **Launch the application**: Run `qslicerpicker` (or double-click the `.app` bundle on macOS)
2. **Configure slicers**: The settings window will open automatically
3. **Enable slicers**: Check the slicers you want to use
4. **Set custom paths** (if needed): Edit slicer paths if they're not auto-detected
5. **Reorder slicers**: Use the up/down arrows to change the order

### Opening Files

1. **Set up file associations** (one-time setup):
   - **macOS**: Right-click a file → Open With → Choose "3D Slicer Picker" → Always Open With
   - **Windows**: Right-click a file → Open With → Choose "3D Slicer Picker" → Always use this app
   - **Linux**: Use your file manager's "Open With" option

2. **Open a 3D model file**: Double-click any supported file type
3. **Select slicer**: Choose from the list of enabled slicers
4. **Open**: Click "Open" to launch the selected slicer with your file

### Supported File Types

- `.3mf` - 3D Manufacturing Format
- `.step` / `.stp` - STEP CAD files
- `.stl` - STereoLithography
- `.svg` - Scalable Vector Graphics
- `.obj` - Wavefront OBJ
- `.amf` - Additive Manufacturing Format
- `.usd*` - Universal Scene Description (`.usd`, `.usda`, `.usdc`)
- `.abc` - Alembic
- `.ply` - Polygon File Format
- `.sla` - SLA format

### Configuration

Configuration is stored in:
- **macOS/Linux**: `~/.qslicerpicker/config.json`
- **Windows**: `%APPDATA%\.qslicerpicker\config.json`

You can edit the configuration file directly or use the settings UI.

### Language Settings

The application supports multiple languages:
- **Turkish** (tr)
- **English** (en) - Default
- **German** (de)
- **French** (fr)

To change the language:
1. Open the Settings window
2. Go to the "Language" tab
3. Select your preferred language
4. The interface will update immediately

The language preference is saved in the configuration file and persists across application restarts.

## 🛠️ Development

### Project Structure

```
QSlicerPicker/
├── internal/
│   ├── config/      # Configuration management
│   ├── filehandler/ # File handling logic
│   ├── i18n/        # Internationalization
│   ├── platform/    # Platform-specific code
│   ├── slicer/      # Slicer management
│   └── ui/          # User interface
├── build/           # Build scripts
├── main.go          # Entry point
└── README.md        # This file
```

### Running Tests

```bash
go test ./...
```

### Code Style

- Follow Go standard formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Areas for Contribution

- 🐛 Bug fixes
- ✨ New features
- 🌍 Additional language translations
- 📝 Documentation improvements
- 🎨 UI/UX enhancements
- 🧪 Test coverage
- 🚀 Performance optimizations

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Muzaffer AKYIL** ([@victorioustr](https://github.com/victorioustr)) - Initial work

## 🙏 Acknowledgments

- Built with [Fyne](https://fyne.io/) - Beautiful cross-platform GUI framework

## 📞 Support

- 🐛 **Found a bug?** [Open an issue](https://github.com/QTechnics/QSlicerPicker/issues)
- 💡 **Have a feature request?** [Open an issue](https://github.com/QTechnics/QSlicerPicker/issues)
- 📧 **Questions?** [Open a discussion](https://github.com/QTechnics/QSlicerPicker/discussions)

## ⭐ Star History

If you find this project useful, please consider giving it a star! ⭐

---

<div align="center">

**Made with ❤️ by [QTechnics](https://github.com/QTechnics)**

[Website](https://github.com/QTechnics) • [Issues](https://github.com/QTechnics/QSlicerPicker/issues) • [Releases](https://github.com/QTechnics/QSlicerPicker/releases)

</div>
