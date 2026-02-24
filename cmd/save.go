package cmd

import (
	"fmt"

	"github.com/riyasyash/nishku/internal/profile"
	"github.com/riyasyash/nishku/internal/window"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save [profile-name]",
	Short: "Save current window positions as a profile",
	Long:  `Captures all window positions, sizes, and display information and saves them as a named profile.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		// Get current window positions
		windows, err := window.GetAllWindows()
		if err != nil {
			return fmt.Errorf("failed to get window positions: %w", err)
		}

		// Get current display configuration
		displays, err := window.GetDisplays()
		if err != nil {
			return fmt.Errorf("failed to get display configuration: %w", err)
		}

		// Save as profile
		prof := profile.Profile{
			Name:     profileName,
			Windows:  windows,
			Displays: displays,
		}

		if err := profile.SaveProfile(prof); err != nil {
			return fmt.Errorf("failed to save profile: %w", err)
		}

		fmt.Printf("✓ Saved profile '%s' with %d windows across %d display(s)\n", 
			profileName, len(windows), len(displays))
		
		// Show display info
		for i, d := range displays {
			primary := ""
			if d.IsPrimary {
				primary = " (primary)"
			}
			fmt.Printf("  Display %d: %dx%d at (%d,%d)%s\n", 
				i+1, d.Width, d.Height, d.X, d.Y, primary)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
