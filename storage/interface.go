package storage

// Storage defines the interface for storing and retrieving Gemara artifacts.
// Implementations can provide local file-based storage or remote storage clients.
type Storage interface {
	// StoreRawYAML stores raw YAML content for a given layer and returns the artifact ID.
	// The YAML content must include metadata.id. Returns the artifact ID on success.
	StoreRawYAML(layer int, yamlContent string) (string, error)

	// Retrieve loads an artifact by layer and ID.
	// Returns the artifact as an interface{} which should be cast to the appropriate type:
	Retrieve(layer int, artifactID string) (interface{}, error)

	// List returns all artifacts for a given layer.
	// If layer is 0, returns artifacts from all layers.
	List(layer int) []*ArtifactIndexEntry

	// Rescan rescans the storage and rebuilds the index.
	// This is useful to discover new artifacts that may have been added externally.
	Rescan() error

	// GetBaseDir returns the base directory path for local storage.
	// For remote storage implementations, this may return an empty string or a logical identifier.
	GetBaseDir() string
}
