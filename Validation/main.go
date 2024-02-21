// main.go

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {
	http.HandleFunc("/validate", validateHandler)
	http.ListenAndServe(":8080", nil);
}

// validateRequestBody validates the request body against the provided JSON schema.
func validateRequestBody(body []byte, schemaPath string) error {
	// Load the YAML schema file
	schemaFile, err := os.Open(schemaPath)
	if err != nil {
		return fmt.Errorf("error opening schema file: %w", err)
	}
	defer schemaFile.Close()

	schemaData, err := io.ReadAll(schemaFile)
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}

	// Convert the YAML schema to JSON
	schemaJSON, err := yaml.YAMLToJSON(schemaData)
	if err != nil {
		return fmt.Errorf("error converting YAML to JSON: %w", err)
	}

	// Load the JSON schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)

	// Load the JSON data to be validated
	dataLoader := gojsonschema.NewBytesLoader(body)

	// Validate the data against the schema
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return fmt.Errorf("error validating request body: %w", err)
	}

	// Check if the data is valid
	if !result.Valid() {
		var errStrings []string
		for _, desc := range result.Errors() {
			errStrings = append(errStrings, desc.String())
		}
		return fmt.Errorf("the request body is not valid: %s", strings.Join(errStrings, ", "))
	}

	return nil
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Validate the request body using the separate function
	schemaPath := "schema.yml" // Replace with the actual path to your YAML schema file
	if err := validateRequestBody(body, schemaPath); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "The request body is valid")
}


