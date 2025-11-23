package week1

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestServerWithTestdata tests the server using commands from the testdata file
func TestServerWithTestdata(t *testing.T) {
	// Parse testdata file
	commands, err := parseTestdata("testdata/week1.txt")
	if err != nil {
		t.Fatalf("Failed to parse testdata: %v", err)
	}

	t.Logf("Loaded %d commands from testdata", len(commands))

	// Get the server handler (no TCP server needed)
	handler := Server()

	// Execute commands
	for i, cmd := range commands {
		lineNum := i + 1

		switch cmd.Operation {
		case "PUT":
			err := executePUT(handler, cmd.Key, cmd.Value)
			if err != nil {
				t.Errorf("Line %d: PUT %s %s failed: %v", lineNum, cmd.Key, cmd.Value, err)
			}

		case "GET":
			statusCode, body, err := executeGET(handler, cmd.Key)
			if err != nil {
				t.Errorf("Line %d: GET %s failed: %v", lineNum, cmd.Key, err)
				continue
			}

			if cmd.Value == "NOT_FOUND" {
				// Expect 404 status
				if statusCode != http.StatusNotFound {
					t.Errorf("Line %d: GET %s expected NOT_FOUND (404) but got status %d with body: %s",
						lineNum, cmd.Key, statusCode, body)
				}
			} else {
				// Expect 200 status and matching value
				if statusCode != http.StatusOK {
					t.Errorf("Line %d: GET %s expected status 200 but got %d with body: %s",
						lineNum, cmd.Key, statusCode, body)
				} else if body != cmd.Value {
					t.Errorf("Line %d: GET %s expected value %q but got %q",
						lineNum, cmd.Key, cmd.Value, body)
				}
			}

		default:
			t.Errorf("Line %d: Unknown operation %q", lineNum, cmd.Operation)
		}
	}
}

// Command represents a parsed testdata command
type Command struct {
	Operation string // "PUT" or "GET"
	Key       string
	Value     string // For PUT: value to store; For GET: expected value or "NOT_FOUND"
}

// parseTestdata reads and parses the testdata file
func parseTestdata(filepath string) ([]Command, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open testdata file: %w", err)
	}
	defer file.Close()

	var commands []Command
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid command at line %d: %s", lineNum, line)
		}

		cmd := Command{
			Operation: parts[0],
			Key:       parts[1],
		}

		if len(parts) > 2 {
			cmd.Value = parts[2]
		}

		commands = append(commands, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading testdata file: %w", err)
	}

	return commands, nil
}

// executePUT sends a PUT request to store a key-value pair using the handler directly
func executePUT(handler http.Handler, key, value string) error {
	req := httptest.NewRequest("POST", "/"+key, strings.NewReader(value))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		return fmt.Errorf("PUT failed with status %d: %s", rec.Code, rec.Body.String())
	}

	return nil
}

// executeGET sends a GET request and returns the response status and body using the handler directly
func executeGET(handler http.Handler, key string) (int, string, error) {
	req := httptest.NewRequest("GET", "/"+key, nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	return rec.Code, rec.Body.String(), nil
}
