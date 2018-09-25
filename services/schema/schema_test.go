package schema

import (
	"testing"
)

// Test creation of the schema structure
func TestLoadAndCreateSchema(t *testing.T) {
	schema, err := LoadAndCreateSchema("./schema_test.json")

	if err != nil {
		t.Error(err)
		t.Fatalf("Schema not loaded")
	}
	if len(schema.components) == 0 {
		t.Fatalf("Map of schemas is empty")
	}

	testSchema := schema.components["testSchema"]
	if testSchema == nil {
		t.Fatalf("Test schema is empty")
	}

	if len(testSchema.Scrapers) != 2 {
		t.Fatalf("Number of scrapers should be 2")
	}
	if len(testSchema.Aggregators) != 1 {
		t.Fatalf("Number of aggregators should be 1")
	}
	if len(testSchema.Feeds) != 1 {
		t.Fatalf("Number of feeds should be 1")
	}
	if testSchema.Name != "testSchema" {
		t.Fatalf("Name of the schema is not set")
	}
	if testSchema.Exchange != "binance" {
		t.Fatalf("Exchange should be set to \"binance\"")
	}
}
