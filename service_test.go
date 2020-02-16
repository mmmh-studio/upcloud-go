package upcloud

import (
	"os"

	"github.com/mmmh-studio/upcloud-go/client"
)

func newTestService() *Service {
	var (
		user     = os.Getenv("UPCLOUD_USERNAME")
		password = os.Getenv("UPCLOUD_PASSWORD")
		client   = client.New(user, password)
	)

	return NewService(client)
}
