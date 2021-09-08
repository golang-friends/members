package config

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetConfigName("members")
	viper.SetConfigType("yaml")
	// internal/config/../.. => project root.
	viper.AddConfigPath("../..")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed while reading config via viper: %v", err))
	}
}

func TestFromViper(t *testing.T) {
	cfg := FromViper()
	assert.NotEmpty(t, cfg.Orgname)
	assert.NotEmpty(t, cfg.Admins)
	assert.NotEmpty(t, cfg.Members)
}
