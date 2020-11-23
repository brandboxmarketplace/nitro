package ls

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// LSCommand is the command for creating new development environments
var LSCommand = &cobra.Command{
	Use:   "ls",
	Short: "List environment containers",
	RunE:  lsMain,
	Example: `  # list all containers for the environment
  nitro ls`,
}

func lsMain(cmd *cobra.Command, args []string) error {
	env := cmd.Flag("environment").Value.String()

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	containers, err := nitro.LS(cmd.Context(), env, args)
	if err != nil {
		return err
	}

	for name, id := range containers {
		fmt.Println(name, id)
	}

	return nil
}