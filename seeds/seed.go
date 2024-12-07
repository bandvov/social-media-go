package seeds

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

func Seed(db *sql.DB, filePath string) {
	// Open the SQL file
	if filePath == "" {
		log.Fatal("filePath not provided")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open SQL file: %v", err)
	}

	defer file.Close()

	// Use bufio to read the file line by line
	var queryBuilder strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments and empty lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		// Accumulate the query
		queryBuilder.WriteString(line)

		// Execute the query when a semicolon is encountered
		if strings.HasSuffix(line, ";") {
			query := queryBuilder.String()
			_, err := db.Exec(query)
			if err != nil {
				log.Printf("Failed to execute query: %s\nError: %v", query, err)
			}
			queryBuilder.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading the SQL file: %v", err)
	}

	fmt.Printf("SQL file %v imported successfully!\n", filePath)
}
