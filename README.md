<div align="center">

# 🪟 Nishku

**Save and restore window positions across displays**

[![Version](https://img.shields.io/badge/version-0.0.1-blue.svg)](https://github.com/riyasyash/nishku/releases)
[![Platform](https://img.shields.io/badge/platform-macOS-lightgrey.svg)](https://github.com/riyasyash/nishku)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)

Perfect for switching between different monitor setups (home vs office, docked vs undocked)

[Features](#features) •
[Installation](#installation) •
[Usage](#usage) •
[Documentation](#documentation) •
[Contributing](#contributing)

</div>

---

## Overview

Nishku is a command-line tool that captures and restores window positions and sizes on macOS. It's designed for people who frequently switch between different monitor configurations and want their window layouts to adapt automatically.

### Why Nishku?

- 🏠 **Home/Office Setup**: Different monitor configs? No problem.
- 💻 **Docked/Undocked**: Seamlessly switch between laptop-only and multi-monitor.
- 🎯 **Smart Adaptation**: Windows automatically adapt when displays change.
- ⚡ **Fast & Light**: Single binary, no runtime dependencies.

## Features

### Core Capabilities

- 💾 **Save window layouts** as named profiles
- 🔄 **Restore layouts instantly** with one command
- 🖥️ **Multi-display support** with automatic adaptation
- 🎯 **Smart window matching** - handles multiple windows per app
- 📐 **Position + Size** - captures both location and dimensions
- 🔒 **Secure & Private** - local storage only, no network access

### Smart Features

- **Hybrid Matching**: Works with apps that have multiple windows (Chrome, iTerm2, VS Code)
- **Display Adaptation**: Automatically adjusts when monitors are disconnected
- **Bounds Checking**: Ensures windows stay visible on screen
- **Built-in Diagnostics**: `nishku doctor` checks your setup

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/riyasyash/nishku
cd nishku

# Build and install
make build
sudo make install

# Verify installation
nishku version
```

### Using Go Install

```bash
go install github.com/riyasyash/nishku@latest
```

### Requirements

- macOS 10.14 (Mojave) or later
- ~5MB disk space
- Accessibility permissions (one-time setup)

## Quick Start

### 1. Check Your Setup

```bash
nishku doctor
```

This will show your current status and guide you through any required setup.

### 2. Grant Accessibility Permissions

To move windows, nishku needs accessibility permissions:

1. Open **System Settings** → **Privacy & Security** → **Accessibility**
2. Click the 🔒 lock icon and authenticate
3. Click the **+** button
4. Add your terminal app (Terminal.app, iTerm.app, etc.)
5. Enable the checkbox ✓
6. **Restart your terminal**

### 3. Save Your First Profile

```bash
# Save current window positions
nishku save work-setup
```

### 4. Restore Anytime

```bash
# Restore to saved layout
nishku load work-setup
```

## Usage

### Basic Commands

```bash
# Save current window layout
nishku save <profile-name>

# Restore window layout
nishku load <profile-name>

# List all profiles
nishku list

# Delete a profile
nishku delete <profile-name>

# Check setup and permissions
nishku doctor

# Show version
nishku version
```

### Example Workflow

```bash
# At the office with dual monitors
nishku save office-dual-monitor

# At home with laptop only
nishku save home-laptop

# Switching between setups
nishku load office-dual-monitor  # When at office
nishku load home-laptop          # When at home
```

## How It Works

### Window Matching

Nishku uses a hybrid approach to match windows:

1. **Groups windows by application** (e.g., all Chrome windows)
2. **Sorts by screen position** (top-to-bottom, left-to-right)
3. **Assigns indices** for reliable matching
4. **Matches on restore** by app + window index

This means:
- ✓ Works with multiple windows per app
- ✓ Handles dynamic window titles (browser tabs, document names)
- ✓ Gracefully skips windows if app isn't running

### Display Adaptation

When your display configuration changes:

- **Stores both absolute and relative positions**
- **Adapts automatically** if displays are missing
- **Ensures visibility** - windows never go off-screen
- **Silent operation** - no prompts or warnings needed

**Example**: Profile saved with 2 monitors, restored with 1 monitor:
- Windows from Display 1 → Restored to same position
- Windows from Display 2 → Moved to primary display (equivalent relative position)

## Documentation

- [TESTING.md](TESTING.md) - Testing guide and troubleshooting
- [SECURITY.md](SECURITY.md) - Security considerations
- [CHANGELOG.md](CHANGELOG.md) - Version history

## Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| macOS (Apple Silicon) | ✅ Supported | arm64 |
| macOS (Intel) | ✅ Supported | amd64 |
| Linux | 🚧 Planned | X11/Wayland support coming |
| Windows | 🚧 Planned | Win32 API support coming |

## Known Limitations

⚠️ **Alpha Release** - This is v0.0.1, expect some rough edges

- Requires accessibility permissions (macOS security requirement)
- Apps must be running before restore (doesn't launch apps)
- Minimized windows may not restore correctly
- Full-screen windows may exit full-screen when restored
- Some system apps don't support Accessibility API

## Troubleshooting

### "Cannot move windows"

Run `nishku doctor` and follow the permission instructions. Make sure to **restart your terminal** after granting permissions.

### "Skipped N windows"

This is normal when:
- Profile has more windows than currently open
- App isn't running
- Window index doesn't exist

### Windows not moving

1. Check permissions: `nishku doctor`
2. Verify app is running
3. Try with a single-window app first
4. See [TESTING.md](TESTING.md) for detailed troubleshooting

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt
```

### Version Management

```bash
# Build with version
VERSION=v0.0.1 make build

# Check version
./nishku version --verbose
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Areas for Contribution

- 🐧 Linux support (X11/Wayland)
- 🪟 Windows support (Win32 API)
- 🧪 Testing and bug reports
- 📚 Documentation improvements
- ✨ Feature suggestions

### Development Setup

```bash
git clone https://github.com/riyasyash/nishku
cd nishku
make deps
make build
./nishku doctor
```

## Roadmap

### v0.1.0
- [ ] Bug fixes from user feedback
- [ ] Shell completions (bash/zsh/fish)
- [ ] Improved error messages

### v0.2.0
- [ ] Profile templates
- [ ] Auto-detect profile based on displays
- [ ] Verbose/debug mode

### v1.0.0
- [ ] Stable API and CLI
- [ ] Comprehensive tests
- [ ] Homebrew formula
- [ ] Linux support

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://go.dev/)
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- Uses macOS CoreGraphics and Accessibility APIs

## Author

Created by [Riyas Yash](https://github.com/riyasyash)

## Support

- 📫 [Issues](https://github.com/riyasyash/nishku/issues)
- 💬 [Discussions](https://github.com/riyasyash/nishku/discussions)
- ⭐ Star this repo if you find it useful!

---

<div align="center">

Made with ❤️ for productivity enthusiasts

</div>
