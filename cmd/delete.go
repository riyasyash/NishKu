package cmd

import (
	"fmt"

	"github.com/riyasyash/nishku/internal/profile"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [profile-name]",
	Short: "Delete a saved profile",
	Long:  `Permanently removes a saved window position profile.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		if err := profile.DeleteProfile(profileName); err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		fmt.Printf("✓ Deleted profile '%s'\n", profileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
