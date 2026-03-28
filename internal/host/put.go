package host

import (
	"database/sql"
	"encoding/json"
	"os"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func PutHost(db *sql.DB, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	hosts := []Host{}
	host := Host{}
	if err := json.Unmarshal(data, &host); err != nil {
		log.Debugf("[host] error occurred while parsing single host file, check for multipl host in the file: %v", err)
		if err := json.Unmarshal(data, &hosts); err != nil {
			log.Debugf("[host] error occurred while reading multiplke hosts from the file: %v", err)
			return err
		}
	} else {
		hosts = append(hosts, host)
	}

	for _, h := range hosts {
		if h.Id == uuid.Nil {
			id, err := uuid.NewUUID()
			if err != nil {
				log.Warnf("[host] error occurred while trying to generate the id for host: %v", err)
				continue
			}
			h.Id = id
		}

		// validate jumphost id
		if h.JumphostID.Valid {
			_, err := GetHostByID(db, h.JumphostID.UUID.String())
			if err != nil {
				log.Warnf("error occurred while verifying the jumphost ID provided: %v", err)
				return err
			}
		}

		// validate region id
		if _, err = region.GetRegionByID(db, h.RegionID.String()); err != nil {
			log.Warnf("error occurred while verifying the region ID provided: %v", err)
			return err
		}

		// validate identity id
		if _, err = identity.GetIdentityByID(db, h.IdentityID); err != nil {
			log.Warnf("error occurred while verifying the identity ID provided: %v", err)
			return err
		}

		if err := h.Store(db); err != nil {
			log.Warnf("[host] error occurred while writing hosts: %v", err)
		}

	}

	return nil

}

func PutTagMapping(db *sql.DB, hostName, tagName string) error {
	host, nTag, err := getHostAndTag(db, hostName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.NewTagMapping(nTag.Id, host.Id)
	if err != nil {
		log.Debugf("[host] failed to create tag mapping for host %q and tag %q: %v", hostName, tagName, err)
		return err
	}

	return tm.Store(db)
}
