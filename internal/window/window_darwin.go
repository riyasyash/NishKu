//go:build darwin

package window

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices

#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>

typedef struct {
    char app_name[256];
    char window_title[256];
    int x;
    int y;
    int width;
    int height;
    int display_id;
    int process_id;
} WindowInfo;

typedef struct {
    int id;
    int x;
    int y;
    int width;
    int height;
    int is_primary;
} DisplayInfo;

int get_all_windows(WindowInfo* windows, int max_windows) {
    CFArrayRef windowList = CGWindowListCopyWindowInfo(
        kCGWindowListOptionOnScreenOnly | kCGWindowListExcludeDesktopElements,
        kCGNullWindowID
    );
    
    if (!windowList) {
        return 0;
    }
    
    CFIndex count = CFArrayGetCount(windowList);
    int windowCount = 0;
    
    for (CFIndex i = 0; i < count && windowCount < max_windows; i++) {
        CFDictionaryRef window = (CFDictionaryRef)CFArrayGetValueAtIndex(windowList, i);
        
        // Get window layer (skip if not normal window)
        CFNumberRef layerRef = (CFNumberRef)CFDictionaryGetValue(window, kCGWindowLayer);
        int layer;
        if (layerRef) {
            CFNumberGetValue(layerRef, kCFNumberIntType, &layer);
            if (layer != 0) continue; // Only normal windows
        }
        
        // Get app name
        CFStringRef appName = (CFStringRef)CFDictionaryGetValue(window, kCGWindowOwnerName);
        if (!appName) continue;
        
        // Get window title
        CFStringRef windowName = (CFStringRef)CFDictionaryGetValue(window, kCGWindowName);
        
        // Get bounds
        CFDictionaryRef bounds = (CFDictionaryRef)CFDictionaryGetValue(window, kCGWindowBounds);
        if (!bounds) continue;
        
        CGRect rect;
        if (!CGRectMakeWithDictionaryRepresentation(bounds, &rect)) continue;
        
        // Get process ID
        CFNumberRef pidRef = (CFNumberRef)CFDictionaryGetValue(window, kCGWindowOwnerPID);
        int pid = 0;
        if (pidRef) {
            CFNumberGetValue(pidRef, kCFNumberIntType, &pid);
        }
        
        // Fill struct
        CFStringGetCString(appName, windows[windowCount].app_name, 256, kCFStringEncodingUTF8);
        
        if (windowName) {
            CFStringGetCString(windowName, windows[windowCount].window_title, 256, kCFStringEncodingUTF8);
        } else {
            windows[windowCount].window_title[0] = '\0';
        }
        
        windows[windowCount].x = (int)rect.origin.x;
        windows[windowCount].y = (int)rect.origin.y;
        windows[windowCount].width = (int)rect.size.width;
        windows[windowCount].height = (int)rect.size.height;
        windows[windowCount].process_id = pid;
        
        // Get display ID for this window
        CGDirectDisplayID displayID = 0;
        uint32_t matchingDisplayCount = 0;
        CGDirectDisplayID displays[32];
        
        // Get all displays that contain this window's center point
        CGPoint windowCenter = CGPointMake(
            rect.origin.x + rect.size.width / 2,
            rect.origin.y + rect.size.height / 2
        );
        
        if (CGGetDisplaysWithPoint(windowCenter, 32, displays, &matchingDisplayCount) == kCGErrorSuccess) {
            if (matchingDisplayCount > 0) {
                displayID = displays[0];
            }
        }
        
        windows[windowCount].display_id = (int)displayID;
        
        windowCount++;
    }
    
    CFRelease(windowList);
    return windowCount;
}

int get_displays(DisplayInfo* displays, int max_displays) {
    uint32_t displayCount = 0;
    CGDirectDisplayID activeDisplays[32];
    
    if (CGGetActiveDisplayList(32, activeDisplays, &displayCount) != kCGErrorSuccess) {
        return 0;
    }
    
    CGDirectDisplayID mainDisplay = CGMainDisplayID();
    
    int count = displayCount < max_displays ? displayCount : max_displays;
    
    for (int i = 0; i < count; i++) {
        CGDirectDisplayID displayID = activeDisplays[i];
        CGRect bounds = CGDisplayBounds(displayID);
        
        displays[i].id = (int)displayID;
        displays[i].x = (int)bounds.origin.x;
        displays[i].y = (int)bounds.origin.y;
        displays[i].width = (int)bounds.size.width;
        displays[i].height = (int)bounds.size.height;
        displays[i].is_primary = (displayID == mainDisplay) ? 1 : 0;
    }
    
    return count;
}

int restore_window(const char* app_name, int pid, int window_index, int x, int y, int width, int height) {
    if (pid <= 0) {
        return 0; // Invalid PID
    }
    
    // Create AXUIElement for the application
    AXUIElementRef app = AXUIElementCreateApplication(pid);
    if (!app) {
        return 0;
    }
    
    // Get all windows for this application
    CFArrayRef windows = NULL;
    AXError error = AXUIElementCopyAttributeValue(
        app,
        kAXWindowsAttribute,
        (CFTypeRef*)&windows
    );
    
    if (error != kAXErrorSuccess || !windows) {
        CFRelease(app);
        return 0;
    }
    
    // Get the count of windows
    CFIndex windowCount = CFArrayGetCount(windows);
    if (windowCount == 0 || window_index >= windowCount) {
        CFRelease(windows);
        CFRelease(app);
        return 0;
    }
    
    // Get the window at the specified index
    AXUIElementRef window = (AXUIElementRef)CFArrayGetValueAtIndex(windows, window_index);
    
    // Set the window position
    CGPoint newPosition = CGPointMake(x, y);
    CFTypeRef positionValue = (CFTypeRef)AXValueCreate(kAXValueTypeCGPoint, &newPosition);
    
    error = AXUIElementSetAttributeValue(window, kAXPositionAttribute, positionValue);
    CFRelease(positionValue);
    
    if (error != kAXErrorSuccess) {
        CFRelease(windows);
        CFRelease(app);
        return 0;
    }
    
    // Set the window size
    CGSize newSize = CGSizeMake(width, height);
    CFTypeRef sizeValue = (CFTypeRef)AXValueCreate(kAXValueTypeCGSize, &newSize);
    
    error = AXUIElementSetAttributeValue(window, kAXSizeAttribute, sizeValue);
    CFRelease(sizeValue);
    
    CFRelease(windows);
    CFRelease(app);
    
    // Return success even if size setting failed (position is more important)
    return 1;
}

int check_accessibility_permissions() {
    // Check if the process is trusted for accessibility
    NSDictionary *options = @{(__bridge NSString *)kAXTrustedCheckOptionPrompt: @NO};
    Boolean isTrusted = AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options);
    return isTrusted ? 1 : 0;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func getAllWindowsImpl() ([]Window, error) {
	return getWindowsMacOS()
}

func getDisplaysImpl() ([]Display, error) {
	return getDisplaysMacOS()
}

func restoreWindowsImpl(windows []Window) (restored, failed int) {
	return restoreWindowsMacOS(windows)
}

func checkPermissionsImpl() (canRead, canWrite bool, err error) {
	return checkPermissionsMacOS()
}

func getWindowsMacOS() ([]Window, error) {
	const maxWindows = 256
	windowInfos := make([]C.WindowInfo, maxWindows)
	
	count := C.get_all_windows(&windowInfos[0], maxWindows)
	
	// Get display information for calculating relative positions
	displays, err := getDisplaysMacOS()
	if err != nil {
		return nil, err
	}
	
	// Create a map of display ID to display for quick lookup
	displayMap := make(map[string]Display)
	for _, d := range displays {
		displayMap[d.ID] = d
	}
	
	windows := make([]Window, 0, count)
	for i := 0; i < int(count); i++ {
		info := windowInfos[i]
		displayID := fmt.Sprintf("%d", info.display_id)
		
		// Calculate relative position on display
		relX, relY := 0.0, 0.0
		if display, ok := displayMap[displayID]; ok {
			if display.Width > 0 && display.Height > 0 {
				relX = float64(int(info.x)-display.X) / float64(display.Width)
				relY = float64(int(info.y)-display.Y) / float64(display.Height)
			}
		}
		
		windows = append(windows, Window{
			AppName:     C.GoString(&info.app_name[0]),
			WindowTitle: C.GoString(&info.window_title[0]),
			X:           int(info.x),
			Y:           int(info.y),
			Width:       int(info.width),
			Height:      int(info.height),
			DisplayID:   displayID,
			RelativeX:   relX,
			RelativeY:   relY,
			ProcessID:   int(info.process_id),
		})
	}
	
	// Assign window indices using hybrid approach
	windows = assignWindowIndices(windows)
	
	return windows, nil
}

func getDisplaysMacOS() ([]Display, error) {
	const maxDisplays = 32
	displayInfos := make([]C.DisplayInfo, maxDisplays)
	
	count := C.get_displays(&displayInfos[0], maxDisplays)
	
	displays := make([]Display, 0, count)
	for i := 0; i < int(count); i++ {
		info := displayInfos[i]
		displays = append(displays, Display{
			ID:        fmt.Sprintf("%d", info.id),
			X:         int(info.x),
			Y:         int(info.y),
			Width:     int(info.width),
			Height:    int(info.height),
			IsPrimary: info.is_primary == 1,
		})
	}
	
	return displays, nil
}

// assignWindowIndices groups windows by app and assigns indices based on screen position
func assignWindowIndices(windows []Window) []Window {
	// Group windows by app name
	appWindows := make(map[string][]int)
	for i, w := range windows {
		appWindows[w.AppName] = append(appWindows[w.AppName], i)
	}
	
	// For each app, sort windows by position (left to right, top to bottom)
	for _, indices := range appWindows {
		// Sort by Y first (top to bottom), then X (left to right)
		for i := 0; i < len(indices); i++ {
			for j := i + 1; j < len(indices); j++ {
				wi := windows[indices[i]]
				wj := windows[indices[j]]
				
				// Primary sort: Y position (top to bottom)
				// Secondary sort: X position (left to right)
				if wi.Y > wj.Y || (wi.Y == wj.Y && wi.X > wj.X) {
					indices[i], indices[j] = indices[j], indices[i]
				}
			}
		}
		
		// Assign window index based on sorted order
		for idx, windowIdx := range indices {
			windows[windowIdx].WindowIndex = idx
		}
	}
	
	return windows
}

func restoreWindowsMacOS(windows []Window) (restored, failed int) {
	// Get current windows to match with saved profile
	currentWindows, err := getWindowsMacOS()
	if err != nil {
		return 0, len(windows)
	}
	
	// Get current display configuration
	currentDisplays, err := getDisplaysMacOS()
	if err != nil {
		return 0, len(windows)
	}
	
	// Create display lookup map
	displayMap := make(map[string]Display)
	for _, d := range currentDisplays {
		displayMap[d.ID] = d
	}
	
	// Find primary display as fallback
	var primaryDisplay Display
	for _, d := range currentDisplays {
		if d.IsPrimary {
			primaryDisplay = d
			break
		}
	}
	
	// Group current windows by app name
	currentByApp := make(map[string][]Window)
	for _, w := range currentWindows {
		currentByApp[w.AppName] = append(currentByApp[w.AppName], w)
	}
	
	// Match and restore windows using adaptive approach
	for _, savedWindow := range windows {
		appWindows, exists := currentByApp[savedWindow.AppName]
		if !exists {
			failed++
			continue
		}
		
		// Match by window_index within the app
		if savedWindow.WindowIndex >= len(appWindows) {
			// Not enough windows of this app currently open
			failed++
			continue
		}
		
		currentWindow := appWindows[savedWindow.WindowIndex]
		
		// Determine target position with display adaptation
		targetX, targetY := adaptWindowPosition(savedWindow, displayMap, primaryDisplay)
		
		// Restore the window position using the actual window index
		appName := C.CString(savedWindow.AppName)
		result := C.restore_window(
			appName,
			C.int(currentWindow.ProcessID),
			C.int(currentWindow.WindowIndex), // Use the current window's index
			C.int(targetX),
			C.int(targetY),
			C.int(savedWindow.Width),
			C.int(savedWindow.Height),
		)
		C.free(unsafe.Pointer(appName))
		
		if result == 1 {
			restored++
		} else {
			failed++
		}
	}
	
	return restored, failed
}

// adaptWindowPosition adapts saved window position to current display configuration
func adaptWindowPosition(window Window, displayMap map[string]Display, primaryDisplay Display) (x, y int) {
	// Try to find the saved display
	targetDisplay, displayExists := displayMap[window.DisplayID]
	
	if displayExists {
		// Display still exists, use absolute position
		x = window.X
		y = window.Y
	} else {
		// Display doesn't exist, use relative position on primary display
		x = primaryDisplay.X + int(float64(primaryDisplay.Width)*window.RelativeX)
		y = primaryDisplay.Y + int(float64(primaryDisplay.Height)*window.RelativeY)
		targetDisplay = primaryDisplay
	}
	
	// Ensure window is within display bounds
	x, y = ensureWindowVisible(x, y, window.Width, window.Height, targetDisplay)
	
	return x, y
}

// ensureWindowVisible ensures a window is visible within a display's bounds
func ensureWindowVisible(x, y, width, height int, display Display) (int, int) {
	// Ensure window isn't off the right edge
	if x+width > display.X+display.Width {
		x = display.X + display.Width - width
	}
	
	// Ensure window isn't off the bottom edge
	if y+height > display.Y+display.Height {
		y = display.Y + display.Height - height
	}
	
	// Ensure window isn't off the left edge
	if x < display.X {
		x = display.X
	}
	
	// Ensure window isn't off the top edge
	if y < display.Y {
		y = display.Y
	}
	
	return x, y
}

func checkPermissionsMacOS() (canRead, canWrite bool, err error) {
	// Try to read windows - this should always work
	windows, err := getWindowsMacOS()
	if err != nil || len(windows) == 0 {
		return false, false, fmt.Errorf("cannot read window information")
	}
	canRead = true
	
	// Check accessibility permissions using the C function
	hasPermission := C.check_accessibility_permissions()
	canWrite = (hasPermission == 1)
	
	return canRead, canWrite, nil
}
