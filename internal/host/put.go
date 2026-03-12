package host

import (
	"database/sql"
	"encoding/json"
	"os"

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

	json.Unmarshal(data, &host)

	return host.Store(db)

}
