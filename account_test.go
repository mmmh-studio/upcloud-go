package upcloud

import (
	"os"
	"testing"
)

func TestGetAccount(t *testing.T) {
	var (
		username = os.Getenv("UPCLOUD_USERNAME")
		svc      = newTestService()
	)

	acc, err := svc.GetAccount()
	if err != nil {
		t.Fatal(err)
	}

	if acc.Username != username {
		t.Fatalf("have %q, want %q", acc.Username, username)
	}
}
