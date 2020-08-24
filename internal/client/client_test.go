package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFataWhenNoTokenIsProvided(t *testing.T) {
	assert.PanicsWithError(t, "GitHub access token was not provided", func() { NewOAuthClient("") })
}
