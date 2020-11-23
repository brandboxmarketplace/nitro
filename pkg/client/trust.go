package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/craftcms/nitro/pkg/sudo"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Trust is used to find the proxy container and get the certificate from the container
// and place into the host machine.
func (cli *Client) Trust(ctx context.Context, env string, args []string) error {
	// find the nitro proxy for the environment
	filter := filters.NewArgs()
	filter.Add("label", EnvironmentLabel+"="+env)
	filter.Add("label", "com.craftcms.nitro.proxy="+env)

	// find the container, should only be one
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get the list of containers, %w", err)
	}

	// make sure there is at least one container
	if len(containers) == 0 {
		return fmt.Errorf("unable to find the container for the proxy")
	}

	// get the contents of the certificate from the container
	cli.Info(fmt.Sprintf("Trusting CA certificate for %s...", env))

	content, err := cli.Exec(ctx, containers[0].ID, []string{"less", "/data/caddy/pki/authorities/local/root.crt"})
	if err != nil {
		return fmt.Errorf("unable to retreive the certificate from the proxy, %w", err)
	}

	// remove special characters from the output
	var stop int
	for i, s := range content {
		if s != 0 && s != 1 {
			stop = i + 1
			break
		}
	}

	// TODO(jasonmccallister) move this to the cmd pkg

	// create a temp file
	f, err := ioutil.TempFile(os.TempDir(), "nitro-cert")
	if err != nil {
		return fmt.Errorf("unable to create a temporary file, %w", err)
	}

	// write the certificate to the temporary file
	if _, err := f.Write(content[stop+1:]); err != nil {
		return fmt.Errorf("unable to write the certificate to the temporary file, %w", err)
	}
	defer f.Close()

	cli.InfoPending("trusting certificate (you will be prompted for a password):\n")

	if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", f.Name()); err != nil {
		cli.Info("Unable to automatically add the certificate\n")
		cli.Info("To install the certificate, run the following command:")

		// TODO show os specific commands
		switch runtime.GOOS {
		default:
			cli.Info(fmt.Sprintf("  sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", f.Name()))
		}

		return nil
	}

	// clean up
	cli.InfoPending("cleaning up")

	if err := os.Remove(f.Name()); err != nil {
		cli.Info("unable to remove temporary file, it will be automatically removed on reboot")
	}

	cli.InfoDone()

	return nil
}