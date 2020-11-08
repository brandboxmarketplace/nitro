package client

import (
	"context"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestInit(t *testing.T) {
	// Arrange
	mock := newMockDockerClient(nil, nil, nil)
	mock.networkCreateResponse = types.NetworkCreateResponse{ID: "test"}
	cli := Client{docker: mock}

	// Expected
	networkReq := types.NetworkCreateRequest{
		NetworkCreate: types.NetworkCreate{
			Driver:     "bridge",
			Attachable: true,
			Labels: map[string]string{
				"nitro": "testing-init",
			},
		},
		Name: "testing-init",
	}
	// TODO(jasonmccallister) set the volume create request
	// TODO(jasonmccallister) set the container create request
	// TODO(jasonmccallister) set the container start request

	// Act
	err := cli.Init(context.TODO(), "testing-init", []string{})

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil, got %v", err)
	}

	// make sure the network create matches the expected
	if !reflect.DeepEqual(mock.networkCreateRequest, networkReq) {
		t.Errorf(
			"expected network create requests to match\ngot:\n%v\nwant:\n%v",
			mock.networkCreateRequest,
			networkReq,
		)
	}
}
