package import_xsh

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"xsh/internal/host"
	"xsh/internal/identity"

	"github.com/charmbracelet/log"
	"github.com/google/shlex"
	"github.com/google/uuid"
)

var (
	// Regex to match ssh-related commands
	re             = regexp.MustCompile(`^\s*ssh(?:\s|$)`)
	aliasRe        = regexp.MustCompile(`^\s*alias\s+\w+=['"]ssh(?:\s|$).*['"]`)
	aliasExtractRe = regexp.MustCompile(`['"](ssh(?:\s|$).*)['"]`)
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
		} else if aliasRe.MatchString(line) {
			match := aliasExtractRe.FindStringSubmatch(line)
			if len(match) > 1 {
				if !slices.Contains(results, line) {
					results = append(results, match[1]) // actual ssh command
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, i := range results {
		host, err := commandToHost(i, db)
		if err != nil {
			log.Warn("[import] error occurred while parsing ssh command", "error", err, "command", i)
			continue
		}
		if err := host.Store(db); err != nil {
			log.Warnf("[import] error occurred while storing host to database: %v", err)
		}

	}

	return nil
}

func commandToHost(command string, db *sql.DB) (*host.Host, error) {
	log.Debugf("Parshing ssh string to create host: %s", command)
	if strings.Contains(command, "$") {
		return nil, fmt.Errorf("found `$` in the command, seems like there is bash variable replacement present, not importing this one")
	}

	h, err := host.GetDefaultHost()
	if err != nil {
		return nil, err
	}
	tokens, err := shlex.Split(command)
	if err != nil {
		return nil, err
	}

	// Skipping the first index, as that will be someething that won't add much of the value in host insert
	index := 1

	for index < len(tokens) {
		data := tokens[index]
		if index == len(tokens)-1 {
			// Reached the last token
			if strings.HasPrefix(data, "-") {
				h.UpdateExtraFlags(data)
			} else {
				h.UpdateUserAddress(data)
			}
			break
		}

		switch data {
		case "-i":
			h.IdentityFile = tokens[index+1]
			h.IdentityID, err = identity.CheckOrCreateIdentity(h.IdentityFile, db)
			if err != nil {
				return nil, err
			}
			index += 2
			continue
		case "-p":
			port, err := strconv.Atoi(tokens[index+1])
			if err != nil {
				log.Warnf("[import] error occurred while conerting port string to integer, falling back to default 22: %v", err)
				port = 22
			}
			h.Port = port

			index += 2
			continue
		case "-o":
			log.Debugf("SSH Option flag found in use: %s", tokens[index+1])
			if strings.Contains(tokens[index+1], "ProxyCommand") {
				log.Debug("import jumphost from ProxyCommand option of ssh")
				jumpHost, err := commandToHost(tokens[index+1], db)
				if err != nil {
					return nil, err
				}
				h.JumphostID = uuid.NullUUID{
					UUID:  jumpHost.Id,
					Valid: true,
				}

				if err := jumpHost.Store(db); err != nil {
					log.Warnf("[import] error occurred while storing jumphost: %v", err)
					return nil, err
				}
			}
			index += 2
			continue
		}

		if strings.HasPrefix(data, "-") {
			// If data is the last token in the list or if next token in the list also starts with hypen
			// in both the cases we can conclude that its an extra flag candidate
			if strings.HasPrefix(tokens[index+1], "-") {
				h.UpdateExtraFlags(data)
				index++
				continue
			}
		}

		h.UpdateUserAddress(data)
		index++
	}

	h.Name = h.Address

	return h, nil
}
