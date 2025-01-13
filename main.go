package main

import (
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
	"github.com/hypermodeinc/modus/sdk/go/pkg/neo4j"
)

func SayHello(name *string) string {

	var s string
	if name == nil {
		s = "World"
	} else {
		s = *name
	}

	return fmt.Sprintf("Hello, %s!", s)
}

const modelName = "text-generator"
const connectionName = "my-neo-4j"

func GenerateQuery(prompt string) (string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", err
	}

	// Create a system message that specifically asks for a Cypher query
	systemPrompt := `Please generate just a Cypher query for Neo4j and nothing else in plain text. The schema is:
Node types (all have only name property):
:dsyn {name: STRING}
:neop {name: STRING}
:topp {name: STRING}
:gngm {name: STRING}
:imft {name: STRING}
:humn {name: STRING}
:orch {name: STRING}
:diap {name: STRING}
:aapp {name: STRING}
:carb {name: STRING}
:patf {name: STRING}
:elii {name: STRING}
:tisu {name: STRING}
:phsu {name: STRING}
:chvf {name: STRING}
:bacs {name: STRING}
:celf {name: STRING}
:mamm {name: STRING}
:inch {name: STRING}
:orgf {name: STRING}
:bpoc {name: STRING}
:cell {name: STRING}
:fndg {name: STRING}
:inpo {name: STRING}
Relationships (all have pmid and description properties):
[r:COEXISTS_WITH {pmid: STRING, description: STRING}]
[r:ISA {pmid: STRING, description: STRING}]
[r:MANIFESTATION_OF {pmid: STRING, description: STRING}]
[r:TREATS {pmid: STRING, description: STRING}]
[r:NEG_AFFECTS {pmid: STRING, description: STRING}]
[r:AFFECTS {pmid: STRING, description: STRING}]
[r:ASSOCIATED_WITH {pmid: STRING, description: STRING}]
[r:PREDISPOSES {pmid: STRING, description: STRING}]
[r:PROCESS_OF {pmid: STRING, description: STRING}]
[r:DIAGNOSES {pmid: STRING, description: STRING}]
[r:DISRUPTS {pmid: STRING, description: STRING}]
[r:PRODUCES {pmid: STRING, description: STRING}]
[r:PART_OF {pmid: STRING, description: STRING}]
[r:AUGMENTS {pmid: STRING, description: STRING}]
[r:LOCATION_OF {pmid: STRING, description: STRING}]
[r:USES {pmid: STRING, description: STRING}]
[r:CAUSES {pmid: STRING, description: STRING}]
[r:STIMULATES {pmid: STRING, description: STRING}]
[r:ADMINISTERED_TO {pmid: STRING, description: STRING}]
[r:NEG_TREATS {pmid: STRING, description: STRING}]
[r:INHIBITS {pmid: STRING, description: STRING}]
[r:INTERACTS_WITH {pmid: STRING, description: STRING}]`

	input, err := model.CreateInput(
		openai.NewSystemMessage(systemPrompt),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return "", err
	}

	input.Temperature = 0.7

	output, err := model.Invoke(input)
	if err != nil {
		return "", err
	}

	// Return the generated Cypher query
	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}

func GetDatabaseSchema() error {
	// Query to get node labels
	labelsQuery := `
	CALL db.labels() 
	YIELD label 
	RETURN collect(label) as labels`

	// Query to get relationship types
	relsQuery := `
	CALL db.relationshipTypes() 
	YIELD relationshipType 
	RETURN collect(relationshipType) as relationships`

	// Execute queries
	labelResult, err := neo4j.ExecuteQuery(connectionName, labelsQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve labels: %w", err)
	}

	relResult, err := neo4j.ExecuteQuery(connectionName, relsQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve relationships: %w", err)
	}

	// Debug print
	fmt.Println("Label Result Records:")
	for i, record := range labelResult.Records {
		labels, ok := record.Get("labels")
		fmt.Printf("Record %d: value=%v, type=%T, ok=%v\n", i, labels, labels, ok)
	}

	fmt.Println("\nRelationship Result Records:")
	for i, record := range relResult.Records {
		rels, ok := record.Get("relationships")
		fmt.Printf("Record %d: value=%v, type=%T, ok=%v\n", i, rels, rels, ok)
	}

	return nil
}

func TestNeo4jConnection() error {
	// Simple test query to count all nodes
	query := "MATCH (n) RETURN count(n) as nodeCount"

	result, err := neo4j.ExecuteQuery(connectionName, query, nil)
	if err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}

	// Print the result
	fmt.Printf("Database connection test successful!\n")
	fmt.Printf("Total number of nodes in database: %+v\n", result)

	return nil
}

func ExecuteGeneratedQuery(prompt string) (string, string, error) {
	// Generate the query
	query, err := GenerateQuery(prompt)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate query: %w", err)
	}

	// Execute the query on Neo4j
	result, err := neo4j.ExecuteQuery(connectionName, query, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute query: %w", err)
	}

	// Extract and format the data from the result
	var formattedResult strings.Builder
	for i, record := range result.Records {
		formattedResult.WriteString(fmt.Sprintf("Record %d:\n", i))
		for _, key := range record.Keys {
			value, ok := record.Get(key)
			if !ok {
				return "", "", fmt.Errorf("key '%s' not found in record %d", key, i)
			}
			formattedResult.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}

	// Return the generated query and the formatted result as strings
	return query, formattedResult.String(), nil
}
