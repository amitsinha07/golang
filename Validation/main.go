// main.go

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

func validateHandler(w http.ResponseWriter, r *http.Request) {
	// Load the YAML schema file
	schemaFile, err := os.Open("schema.yml")
	if err != nil {
		http.Error(w, "Error opening schema file", http.StatusInternalServerError)
		return
	}
	defer schemaFile.Close()

	schemaData, err := io.ReadAll(schemaFile)
	if err != nil {
		http.Error(w, "Error reading schema file", http.StatusInternalServerError)
		return
	}

	// Convert the YAML schema to JSON
	schemaJSON, err := yaml.YAMLToJSON(schemaData)
	if err != nil {
		http.Error(w, "Error converting YAML to JSON", http.StatusInternalServerError)
		return
	}

	// Load the JSON schema
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Load the JSON data to be validated
	dataLoader := gojsonschema.NewBytesLoader(body)

	// Validate the data against the schema
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		http.Error(w, "Error validating request body", http.StatusInternalServerError)
		return
	}

	// Check if the data is valid
	if result.Valid() {
		fmt.Fprint(w, "The request body is valid")
	} else {
		fmt.Fprint(w, "The request body is not valid. See errors :\n")
		for _, desc := range result.Errors() {
			fmt.Fprintf(w, "- %s\n", desc)
		}
	}
}
