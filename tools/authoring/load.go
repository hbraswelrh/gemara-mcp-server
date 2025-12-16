package authoring

import (
	"github.com/ossf/gemara/layer1"
	"github.com/ossf/gemara/layer2"
	"github.com/ossf/gemara/layer3"
)

// LoadArtifactsDir loads Gemara artifacts from the artifacts directory
// This synchronizes the in-memory cache with what's stored on disk
func (g *GemaraAuthoringTools) LoadArtifactsDir() {
	// If storage is available, use its index to know what to load
	// Otherwise, scan directories directly
	if g.storage != nil {
		// Use storage index to efficiently load artifacts
		// Layer 1
		for _, entry := range g.storage.List(1) {
			if _, exists := g.layer1Guidance[entry.ID]; !exists {
				if retrieved, err := g.storage.Retrieve(1, entry.ID); err == nil {
					if guidance, ok := retrieved.(*layer1.GuidanceDocument); ok {
						g.layer1Guidance[entry.ID] = guidance
					}
				}
			}
		}
		// Layer 2
		for _, entry := range g.storage.List(2) {
			if _, exists := g.layer2Catalogs[entry.ID]; !exists {
				if retrieved, err := g.storage.Retrieve(2, entry.ID); err == nil {
					if catalog, ok := retrieved.(*layer2.Catalog); ok {
						g.layer2Catalogs[entry.ID] = catalog
					}
				}
			}
		}
		// Layer 3
		for _, entry := range g.storage.List(3) {
			if _, exists := g.layer3Policies[entry.ID]; !exists {
				if retrieved, err := g.storage.Retrieve(3, entry.ID); err == nil {
					if policy, ok := retrieved.(*layer3.PolicyDocument); ok {
						g.layer3Policies[entry.ID] = policy
					}
				}
			}
		}
		return
	}
}
