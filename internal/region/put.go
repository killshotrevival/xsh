package region

import (
	"database/sql"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
)

func PutRegion(db *sql.DB, name string) error {
	region, err := NewRegion(name)
	if err != nil {
		return err
	}
	return region.Store(db)
}

func PutTagMapping(db *sql.DB, identittyName, tagName string) error {
	host, nTag, err := getRegionAndTag(db, identittyName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.NewTagMapping(nTag.Id, host.Id)
	if err != nil {
		log.Debugf("[region] failed to create tag mapping for region %q and tag %q: %v", identittyName, tagName, err)
		return err
	}

	return tm.Store(db)
}
