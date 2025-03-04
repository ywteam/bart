package spinners_test

import (
	"os"
	"testing"
	spinners "yellowteam/lib/spinners"
)

func TestInit(t *testing.T) {
	// Ensure the spinners are downloaded and initialized
	spinners.Download()

	// Check if the default spinner is initialized
	defaultSpinner := spinners.Spin("default")
	if defaultSpinner == nil {
		t.Errorf("Expected default spinner to be initialized")
	}

	// Check if the spinners are loaded from the cache file
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("Failed to get user cache directory: %v", err)
	}
	cacheFile := cacheDir + "/spinners.json"
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Errorf("Expected spinners cache file to exist")
	}

	// Check if the spinners map is populated
	if len(spinners.Spinners()) == 0 {
		t.Errorf("Expected spinners map to be populated")
	}
}
func TestAllSpinners(t *testing.T) {
	for _, spinner := range spinners.Spinners() {
		if spinners.Spin(spinner) == nil {
			t.Errorf("Expected spinner %s to be initialized", spinner)
		}
		spinners.Spin(spinner).Start(
			"Testing spinner: "+spinner,
			func() bool {
				return len(spinner)%2 == 0
			},
		)
	}
}
