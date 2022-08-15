package avro_test

import "github.com/odenio/avro/v2"

func ConfigTeardown() {
	// Reset the caches
	avro.DefaultConfig = avro.Config{}.Freeze()
}
