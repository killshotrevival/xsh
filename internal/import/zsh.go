package import_xsh

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"slices"
)

func cleanHistoryLine(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == ';' {
			return line[i+1:]
		}
	}
	return line
}

func Import(path string, db *sql.DB) error {

	// Regex to match ssh-related commands
	re := regexp.MustCompile(`^\s*ssh(?:\s|$)`)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		line = cleanHistoryLine(line)
		if re.MatchString(line) {
			if !slices.Contains(results, line) {
				results = append(results, line)

			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, i := range results {
		fmt.Println(i)
	}

	return nil
}
