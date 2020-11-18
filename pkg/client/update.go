package client

import (
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func (cli *Client) Update(ctx context.Context, images []string) error {
	cli.Info("Updating images from docker hub")

	for _, image := range images {
		rdr, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
		if err != nil {
			return fmt.Errorf("unable to pull image %s, %w", image, err)
		}

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			return fmt.Errorf("unable to read the output from pulling the image, %w", err)
		}

		cli.SubInfo(image, "update complete")
	}

	cli.Info("All images updated successfully")

	return nil
}
