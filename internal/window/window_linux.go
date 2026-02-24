//go:build linux

package window

func getAllWindowsImpl() ([]Window, error) {
	return getWindowsLinux()
}

func getDisplaysImpl() ([]Display, error) {
	return getDisplaysLinux()
}

func restoreWindowsImpl(windows []Window) (restored, failed int) {
	return restoreWindowsLinux(windows)
}

func checkPermissionsImpl() (canRead, canWrite bool, err error) {
	return checkPermissionsLinux()
}

// getWindowsLinux retrieves all visible windows on Linux
func getWindowsLinux() ([]Window, error) {
	// TODO: Implement using X11 or Wayland
	// X11: github.com/BurntSushi/xgb
	// Wayland: More complex, compositor-specific
	return nil, nil
}

// getDisplaysLinux retrieves display information on Linux
func getDisplaysLinux() ([]Display, error) {
	// TODO: Implement using X11 or Wayland
	return nil, nil
}

// restoreWindowsLinux restores window positions on Linux
func restoreWindowsLinux(windows []Window) (restored, failed int) {
	// TODO: Implement using X11 or Wayland
	return 0, len(windows)
}

func checkPermissionsLinux() (canRead, canWrite bool, err error) {
	// TODO: Implement permission check for Linux
	return false, false, nil
}
