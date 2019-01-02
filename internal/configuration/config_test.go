package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanGetDefaultConfig(t *testing.T) {
	assert := assert.New(t)

	config := DefaultConfig()
	assert.Equal(":80", config.Server.Addr)
	assert.True(config.Log.Debug)
	assert.Equal("secret", config.Github.SecretKey)
	assert.Equal("bot-user", config.Github.Username)
	assert.Equal("password", config.Github.Password)
}

func TestCanLoadFromFile(t *testing.T) {
	assert := assert.New(t)

	cwd, err := os.Getwd()
	assert.NoError(err)

	config := MustLoadFromFile(cwd + "/testdata/config.yaml")
	assert.Equal(":8080", config.Server.Addr)
	assert.Equal("secret", config.Github.SecretKey)
	assert.Equal(2, len(config.Reviewers))

	phpReviewers := config.Reviewers[0]
	assert.Equal("PHP", phpReviewers.Language)
	assert.Equal([]string{"username1", "username2"}, phpReviewers.Users)
}

func TestFromEnvironment(t *testing.T) {
	assert := assert.New(t)

	envs := map[string]string{
		"GITHUB_PASSWORD": "bar",
		"GO_REVIEWERS":    "foo,bar",
	}

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			t.Fatal(err)
		}
	}

	config := FromEnvironment()
	assert.Equal("bar", config.Github.Password)

	for _, rc := range config.Reviewers {
		if rc.Language != "GO" {
			continue
		}

		assert.Subset(rc.Users, []string{"foo", "bar"})
	}
}
