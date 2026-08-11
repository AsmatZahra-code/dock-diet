package scanner

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ScanRemoteImage fetches image metadata from the registry and analyzes it
func ScanRemoteImage(imageName string) error {
	fmt.Printf("🔍 Fetching metadata for image: %s...\n", imageName)

	// Parse the image name (e.g., ubuntu:latest)
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return fmt.Errorf("invalid image name: %v", err)
	}

	// Fetch image details from the remote registry (Docker Hub)
	img, err := remote.Image(ref)
	if err != nil {
		return fmt.Errorf("could not fetch image from registry: %v", err)
	}

	// Get total size of the image
	layers, err := img.Layers()
	if err != nil {
		return fmt.Errorf("could not read layers: %v", err)
	}

	var totalSize int64 = 0
	for _, layer := range layers {
		size, _ := layer.Size()
		totalSize += size
	}

	// Convert bytes to Megabytes (MB)
	sizeMB := float64(totalSize) / (1024 * 1024)
	layerCount := len(layers)

	fmt.Println("-------------------------------------------------")
	fmt.Printf("📦 IMAGE: %s\n", imageName)
	fmt.Printf("📏 SIZE: %.2f MB\n", sizeMB)
	fmt.Printf("🥞 LAYERS: %d\n", layerCount)
	fmt.Println("-------------------------------------------------")

	// Analyze the results
	issues := 0
	if sizeMB > 500 {
		fmt.Println("⚠️  ISSUE: Image size is extremely large (> 500MB). Consider using an Alpine or Slim base.")
		issues++
	}
	if layerCount > 15 {
		fmt.Println("⚠️  ISSUE: Too many layers (> 15). Try consolidating RUN instructions.")
		issues++
	}

	if issues == 0 {
		fmt.Println("✅ This image is well-optimized in terms of size and layers!")
	} else {
		fmt.Printf("🚨 Found %d optimization warning(s) for this image.\n", issues)
	}

	return nil
}