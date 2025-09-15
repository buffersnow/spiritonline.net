package repositories

import (
	"database/sql"
	"time"

	"buffersnow.com/spiritonline/pkg/db"
)

type WFCAccount struct {
	WFCID         int64          `db:"wfc_id"`
	LinkedID      sql.NullInt64  `db:"linked_id"`
	NandIDs       db.IntegerList `db:"nand_ids"`
	ConsoleIDs    db.StringList  `db:"console_ids"`
	IPs           db.IPList      `db:"ip_addrs"`
	MacAddresses  db.StringList  `db:"mac_addrs"`
	LastUpdatedOn time.Time      `db:"last_updated_on"`
}
