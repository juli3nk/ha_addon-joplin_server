package main

import (
	"dagger/ha-addon-joplin-server/internal/dagger"
)

const (
	addonName      = "ha_addon-joplin_server"
	addonSourceUrl = "github.com/juli3nk/ha_addon-joplin_server"
	appSourceUrl   = "https://github.com/laurent22/joplin"
)

type RegistryAuth struct {
	Address  string
	Username string
	Secret   *dagger.Secret
}

type HaAddonJoplinServer struct {
	Worktree     *dagger.Directory
	RegistryAuth *RegistryAuth
	Containers   []*dagger.Container
}

func New(
	source *dagger.Directory,
	// +optional
	registryAddress string,
	// +optional
	registryUsername string,
	// +optional
	registrySecret *dagger.Secret,
) *HaAddonJoplinServer {
	haAddon := HaAddonJoplinServer{Worktree: source}

	if len(registryAddress) > 0 {
		registryAuth := RegistryAuth{
			Address:  registryAddress,
			Username: registryUsername,
			Secret:   registrySecret,
		}

		haAddon.RegistryAuth = &registryAuth
	}

	return &haAddon
}
