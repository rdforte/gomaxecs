package config

import (
	"os"
	"strings"
)

const metaURIEnv = "ECS_CONTAINER_METADATA_URI_V4"

func New() Config {
	uri := metadataURI()

	pathParts := strings.Split(uri, "/")
	containerID := pathParts[len(pathParts)-1]

	return Config{
		MetadataURI: uri,
		ConainerID:  containerID,
	}
}

func metadataURI() string {
	return os.Getenv(metaURIEnv)
}

type Config struct {
	MetadataURI string
	ConainerID  string
}
