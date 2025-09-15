package repositories

import (
	"database/sql"
	"net"
	"time"

	"buffersnow.com/spiritonline/pkg/db"
	"buffersnow.com/spiritonline/pkg/log"

	sq "github.com/Masterminds/squirrel"
)

type WFCAccountRepo struct {
	sql    *db.SQL
	logger log.LoggingFactory
}

type WFCAccount struct {
	WFCID         int64          `db:"wfc_id"`
	LinkedID      sql.NullInt64  `db:"linked_id"`
	NandIDs       db.IntegerList `db:"nand_ids"`
	ConsoleIDs    db.StringList  `db:"console_ids"`
	IPs           db.IPList      `db:"ip_addrs"`
	MacAddresses  db.StringList  `db:"mac_addrs"`
	LastUpdatedOn time.Time      `db:"last_updated_on"`
}

type WFCAccountQuery struct {
	ConsoleID string
	NandID    int64
	IP        net.IP
	MAC       string
}

func (w *WFCAccountRepo) Insert(query WFCAccountQuery) (int64 /*wfc_id*/, error) {
	_, err := w.sql.Insert(
		sq.Insert("wfc_account").
			Columns("console_ids", "nand_ids", "ip_addrs", "mac_addrs", "last_update_on").
			Values(
				db.StringList{query.ConsoleID}, db.IntegerList{query.NandID},
				db.IPList{query.IP}, db.StringList{query.MAC}, time.Now(),
			),
	)

	if err != nil {
		return -1, nil
	}

	wfcid := int64(0)
	err = w.sql.Get(&wfcid, sq.Select("wfc_id").From("wfc_accounts").Where(
		"JSON_CONTAINS(console_ids, JSON_QUOTE(?))", query.ConsoleID,
	))

	return wfcid, err
}
