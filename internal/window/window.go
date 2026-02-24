package window

// Window represents a window's position and metadata
type Window struct {
	AppName     string  `json:"app_name"`
	WindowTitle string  `json:"window_title"`
	WindowIndex int     `json:"window_index"` // Position in app's window list (for matching)
	X           int     `json:"x"`
	Y           int     `json:"y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	DisplayID   string  `json:"display_id"`
	// Relative position on display (0.0-1.0) for display-independent restoration
	RelativeX   float64 `json:"relative_x"`
	RelativeY   float64 `json:"relative_y"`
	ProcessID   int     `json:"process_id,omitempty"`
}

// Display represents a display/monitor configuration
type Display struct {
	ID        string `json:"id"`
	X         int    `json:"x"`          // Position in screen space
	Y         int    `json:"y"`          // Position in screen space
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	IsPrimary bool   `json:"is_primary"`
}

// Platform-specific functions are implemented in window_*.go files with build tags
// GetAllWindows retrieves all visible windows with their positions
func GetAllWindows() ([]Window, error) {
	return getAllWindowsImpl()
}

// GetDisplays retrieves information about all connected displays
func GetDisplays() ([]Display, error) {
	return getDisplaysImpl()
}

// RestoreWindows restores window positions from a saved state
// Returns the number of successfully restored windows and the number of failures
func RestoreWindows(windows []Window) (restored, failed int) {
	return restoreWindowsImpl(windows)
}

// CheckPermissions checks if the app has necessary permissions to move windows
func CheckPermissions() (canRead, canWrite bool, err error) {
	return checkPermissionsImpl()
}
