package cmd

import (
	"fmt"
	"time"

	"github.com/riyasyash/nishku/internal/profile"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved profiles",
	Long:  `Display all saved window position profiles with details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := profile.ListProfiles()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles saved yet. Use 'nishku save <name>' to create one.")
			return nil
		}

		fmt.Println("Saved Profiles:")
		fmt.Println()
		for _, p := range profiles {
			fmt.Printf("  %s\n", p.Name)
			fmt.Printf("    Windows: %d\n", len(p.Windows))
			fmt.Printf("    Displays: %d\n", len(p.Displays))
			if len(p.Displays) > 0 {
				for i, d := range p.Displays {
					primary := ""
					if d.IsPrimary {
						primary = " (primary)"
					}
					fmt.Printf("      Display %d: %dx%d%s\n", 
						i+1, d.Width, d.Height, primary)
				}
			}
			fmt.Printf("    Created: %s\n", p.CreatedAt.Format(time.RFC822))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
