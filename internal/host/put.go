package host

import (
	"database/sql"
	"encoding/json"
	"os"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func PutHost(db *sql.DB, filepath string) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	host := Host{Id: id}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &host); err != nil {
		return err
	}

	return host.Store(db)

}

func PutTagMapping(db *sql.DB, hostName, tagName string) error {
	host, nTag, err := getHostAndTag(db, hostName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.NewTagMapping(nTag.Id, host.Id)
	if err != nil {
		log.Debugf("error occurred while creating new tag mapping object; %v", err)
		return err
	}

	return tm.Store(db)
}
