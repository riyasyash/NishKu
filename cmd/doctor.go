package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/riyasyash/nishku/internal/window"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check nishku setup and permissions",
	Long:  `Diagnose common issues with nishku setup, permissions, and configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running nishku diagnostics...\n")

		allGood := true

		// Check profile directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("✗ Cannot find home directory")
			allGood = false
		} else {
			profileDir := filepath.Join(home, ".nishku", "profiles")
			if _, err := os.Stat(profileDir); os.IsNotExist(err) {
				fmt.Printf("⚠ Profile directory doesn't exist yet: %s\n", profileDir)
				fmt.Println("  (This is normal if you haven't saved any profiles)")
			} else {
				fmt.Printf("✓ Profile directory exists: %s\n", profileDir)
			}
		}

		// Check platform
		fmt.Printf("✓ Platform: %s\n", runtime.GOOS)

		// Check window reading capability
		canRead, canWrite, err := window.CheckPermissions()

		if err != nil {
			fmt.Printf("✗ Error checking permissions: %v\n", err)
			allGood = false
		} else {
			if canRead {
				fmt.Println("✓ Can read window positions")
			} else {
				fmt.Println("✗ Cannot read window positions")
				allGood = false
			}

			if canWrite {
				fmt.Println("✓ Can move windows (accessibility permissions granted)")
			} else {
				fmt.Println("✗ Cannot move windows - accessibility permissions needed")
				allGood = false

				// Platform-specific permission instructions
				fmt.Println()
				printPermissionInstructions()
			}
		}

		// Test getting current windows
		fmt.Println()
		windows, err := window.GetAllWindows()
		if err != nil {
			fmt.Printf("✗ Failed to get windows: %v\n", err)
			allGood = false
		} else {
			fmt.Printf("✓ Found %d windows currently open\n", len(windows))
			
			// Show app breakdown
			appCounts := make(map[string]int)
			for _, w := range windows {
				appCounts[w.AppName]++
			}
			
			if len(appCounts) > 0 {
				fmt.Println("\n  Window breakdown:")
				for app, count := range appCounts {
					fmt.Printf("    - %s: %d window(s)\n", app, count)
				}
			}
		}

		fmt.Println()
		if allGood {
			fmt.Println("✓ All checks passed! You're ready to use nishku.")
		} else {
			fmt.Println("⚠ Some issues found. See above for details.")
		}

		return nil
	},
}

func printPermissionInstructions() {
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("To grant accessibility permissions on macOS:")
		fmt.Println()
		fmt.Println("  1. Open System Settings (or System Preferences)")
		fmt.Println("  2. Go to Privacy & Security → Accessibility")
		fmt.Println("  3. Click the lock icon and authenticate")
		fmt.Println("  4. Click the + button")
		fmt.Println("  5. Navigate to Applications → Utilities")
		fmt.Println("  6. Select your terminal app (Terminal.app, iTerm.app, etc.)")
		fmt.Println("  7. Make sure the checkbox next to it is enabled")
		fmt.Println("  8. Restart your terminal application")
		fmt.Println()
		fmt.Println("For detailed instructions with screenshots, visit:")
		fmt.Println("  https://github.com/riyasyash/nishku#permissions")
	case "linux":
		fmt.Println("Linux support is not yet implemented.")
	case "windows":
		fmt.Println("Windows support is not yet implemented.")
	default:
		fmt.Println("Platform not supported yet.")
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
