package schema

import (
	"io/ioutil"

	"github.com/ppincak/rysen/api"
)

// Schema container
type Schema struct {
	// Todo: rename to components
	echangeSchemas map[string]*ExchangeSchemaMetadata
}

// Load schema.json file
func LoadAndCreateSchema(url string) (*Schema, error) {
	bytes, err := ioutil.ReadFile(url)
	if err != nil {
		return nil, err
	}
	return CreateSchema(bytes)
}

// Create schema from json
func CreateSchema(jsonSchema []byte) (*Schema, error) {
	var echangeSchemas map[string]*ExchangeSchemaMetadata
	err := api.UnmarshallAs(jsonSchema, &echangeSchemas)
	if err != nil {
		return nil, err
	}
	for name, schema := range echangeSchemas {
		schema.Name = name
	}

	return &Schema{
		echangeSchemas: echangeSchemas,
	}, nil
}

// Return single component from schema
func (schema *Schema) Component(component string) *ExchangeSchemaMetadata {
	return schema.echangeSchemas[component]
}
