package schema

import (
	"io/ioutil"

	"github.com/ppincak/rysen/api"
)

// TODO: remove this struct as it is now esentially useless
// Schema container
type Schema struct {
	components map[string]*ExchangeSchemaMetadata
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
	var components map[string]*ExchangeSchemaMetadata
	err := api.UnmarshallAs(jsonSchema, &components)
	if err != nil {
		return nil, err
	}
	for name, schema := range components {
		schema.Name = name
	}

	return &Schema{
		components: components,
	}, nil
}

// Return single component from schema
func (schema *Schema) Component(component string) *ExchangeSchemaMetadata {
	return schema.components[component]
}
