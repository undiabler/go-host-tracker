package htracker

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testLogin = os.Getenv("TEST_LOGIN")
	testPass  = os.Getenv("TEST_PASS")
)

func TestConnect(t *testing.T) {

	htClient := &HostTrackerClient{
		Login:    testLogin,
		Password: testPass,
	}

	err := htClient.auth()
	// assert.NoError(t, err, "auth failed")
	if err != nil {
		log.Fatalf("auth failed: %s", err)
	}

	assert.NotNil(t, htClient.token, "token must not be empty")
	assert.NotEqual(t, "", htClient.token.Ticket, "ticket must not be empty")

	t.Logf("token: %+v", htClient.token)

	id, err := htClient.NewHttpTask(map[string]interface{}{
		"url":                   "https://google.com",
		"checkDomainExpiration": true,
	})
	assert.NoError(t, err, "new task failed")

	t.Logf("task id: %+v", id)
}
