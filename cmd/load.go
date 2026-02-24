package cmd

import (
	"fmt"

	"github.com/riyasyash/nishku/internal/profile"
	"github.com/riyasyash/nishku/internal/window"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [profile-name]",
	Short: "Restore window positions from a profile",
	Long:  `Restores all window positions and sizes from a previously saved profile.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		// Load profile
		prof, err := profile.LoadProfile(profileName)
		if err != nil {
			return fmt.Errorf("failed to load profile: %w", err)
		}

		// Restore window positions
		restored, failed := window.RestoreWindows(prof.Windows)

		if restored == 0 && failed > 0 {
			fmt.Printf("✗ Failed to restore any windows from profile '%s'\n", profileName)
			fmt.Println()
			fmt.Println("Common causes:")
			fmt.Println("  • Apps from the profile are not currently running")
			fmt.Println("  • Accessibility permissions not granted (run 'nishku doctor' to check)")
			fmt.Println()
			return fmt.Errorf("restoration failed")
		}

		fmt.Printf("✓ Restored %d windows from profile '%s'\n", restored, profileName)
		if failed > 0 {
			fmt.Printf("⚠ Skipped %d windows (apps may not be running or have fewer windows open)\n", failed)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
