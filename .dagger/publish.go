package main

import (
	"context"
	"fmt"

	"dagger/ha-addon-joplin-server/internal/dagger"
)

func (m *HaAddonJoplinServer) Publish(
	ctx context.Context,
	// +optional
	registryNamespace string,
) error {
	if len(m.Containers) == 0 {
		return fmt.Errorf("error: build containers first")
	}

	addonVersion, err := m.Containers[0].Label(ctx, "io.hass.version")
	if err != nil {
		return err
	}

	ctr := dag.Container()

	imageName := fmt.Sprintf("%s:%s", addonName, addonVersion)

	if len(m.RegistryAuth.Address) > 0 {
		imageName = fmt.Sprintf("%s/%s/%s", m.RegistryAuth.Address, registryNamespace, imageName)

		ctr = ctr.WithRegistryAuth(m.RegistryAuth.Address, m.RegistryAuth.Username, m.RegistryAuth.Secret)
	}

	_, err = ctr.Publish(ctx, imageName, dagger.ContainerPublishOpts{
		PlatformVariants: m.Containers,
	})

	return err
}
