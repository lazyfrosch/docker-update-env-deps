package repos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadGitHub(t *testing.T) {
	repo, err := LoadGitHub("https://github.com/lazyfrosch/docker-icingaweb2.git")
	assert.NoError(t, err)
	assert.Equal(t, "lazyfrosch", repo.Owner)
	assert.Equal(t, "docker-icingaweb2", repo.Repo)

	_, err = LoadGitHub("https://github.com/lazyfrosch/nonexisting.git")
	assert.Error(t, err)
}

func TestGitHub_LoadReleases(t *testing.T) {
	repo, err := LoadGitHub("https://github.com/Icinga/icinga2.git")
	assert.NoError(t, err)

	list, err := repo.LoadReleases()
	assert.NoError(t, err)
	assert.Greater(t, 1, len(list))
}
