package registry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/containrrr/watchtower/pkg/registry/auth"
	"github.com/containrrr/watchtower/pkg/registry/helpers"
	"github.com/containrrr/watchtower/pkg/registry/manifest"
	watchtowerTypes "github.com/containrrr/watchtower/pkg/types"
	ref "github.com/distribution/reference"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

// GetPullOptions creates a struct with all options needed for pulling images from a registry
func GetPullOptions(imageName string) (types.ImagePullOptions, error) {
	auth, err := EncodedAuth(imageName)
	log.Debugf("Got image name: %s", imageName)
	if err != nil {
		return types.ImagePullOptions{}, err
	}

	if auth == "" {
		return types.ImagePullOptions{}, nil
	}

	// CREDENTIAL: Uncomment to log docker config auth
	// log.Tracef("Got auth value: %s", auth)

	return types.ImagePullOptions{
		RegistryAuth:  auth,
		PrivilegeFunc: DefaultAuthHandler,
	}, nil
}

// DefaultAuthHandler will be invoked if an AuthConfig is rejected
// It could be used to return a new value for the "X-Registry-Auth" authentication header,
// but there's no point trying again with the same value as used in AuthConfig
func DefaultAuthHandler() (string, error) {
	log.Debug("Authentication request was rejected. Trying again without authentication")
	return "", nil
}

// WarnOnAPIConsumption will return true if the registry is known-expected
// to respond well to HTTP HEAD in checking the container digest -- or if there
// are problems parsing the container hostname.
// Will return false if behavior for container is unknown.
func WarnOnAPIConsumption(container watchtowerTypes.Container) bool {

	normalizedRef, err := ref.ParseNormalizedNamed(container.ImageName())
	if err != nil {
		return true
	}

	containerHost, err := helpers.GetRegistryAddress(normalizedRef.Name())
	if err != nil {
		return true
	}

	if containerHost == helpers.DefaultRegistryHost || containerHost == "ghcr.io" {
		return true
	}

	return false
}

type tagsListResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func ListTags(container watchtowerTypes.Container, registryAuth string) ([]string, error) {
	token, err := auth.GetToken(container, registryAuth)

	if err != nil {
		return nil, err
	}

	tagsUrl, err := manifest.BuildTagURL(container)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", tagsUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("registry responded to tags list request with %q", res.Status)
	}

	var response tagsListResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Tags, nil
}
