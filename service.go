package upcloud

import (
	"github.com/mmmh-studio/upcloud-go/client"
)

const apiURL = "https://api.upcloud.com/1.3"

type Service struct {
	client *client.Client
}

func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
	}
}
