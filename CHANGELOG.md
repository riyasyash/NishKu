# Changelog

All notable changes to nishku will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.1] - 2026-02-03

### Added
- Initial development release of nishku window position manager
- Window position capture using CoreGraphics API
- Window restoration using Accessibility API
- Multi-display support with automatic adaptation
- Profile management (save, load, list, delete)
- Hybrid window matching for multiple windows per app
- Silent display adaptation when monitor configuration changes
- Built-in diagnostics tool (`nishku doctor`)
- Version command (`nishku version`)
- Comprehensive documentation (README, TESTING, IMPLEMENTATION, SECURITY)
- Secure file permissions (0700/0600)
- Support for macOS (darwin/amd64 and darwin/arm64)

### Commands
- `nishku save <profile>` - Save current window layout
- `nishku load <profile>` - Restore window layout
- `nishku list` - List all saved profiles
- `nishku delete <profile>` - Delete a profile
- `nishku doctor` - Check setup and permissions
- `nishku version` - Show version information

### Features
- Position-based window indexing for reliable multi-window matching
- Relative positioning fallback when displays are missing
- Bounds checking to keep windows on screen
- Zero network access, local-only storage
- No telemetry or tracking
- Single binary with no runtime dependencies

### Platforms
- macOS 10.14+ (Mojave or later)
- Apple Silicon (arm64) and Intel (amd64) support

### Known Limitations
- Requires accessibility permissions for window restoration
- Apps must be running before restore (doesn't launch apps)
- Minimized windows may not restore correctly
- Full-screen windows may exit full-screen when restored
- Some system apps don't support Accessibility API

## Roadmap to v1.0.0

### Planned for v0.1.0
- [ ] Additional testing with various apps
- [ ] Bug fixes from initial user feedback
- [ ] Shell completions (bash/zsh/fish)
- [ ] Improved error messages

### Planned for v0.2.0
- [ ] Profile templates
- [ ] Auto-detect profile based on display configuration
- [ ] Verbose/debug mode
- [ ] Performance optimizations

### Planned for v1.0.0
- [ ] Stable API and CLI interface
- [ ] Comprehensive test coverage
- [ ] Production-ready documentation
- [ ] Homebrew formula
- [ ] Linux support (X11)

[0.0.1]: https://github.com/riyasyash/nishku/releases/tag/v0.0.1
