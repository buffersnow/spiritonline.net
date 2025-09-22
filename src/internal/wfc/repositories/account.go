package repositories

import (
	"database/sql"
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
	Serials       db.StringList  `db:"console_sns"`
	FriendCodes   db.IntegerList `db:"console_fcs"`
	IPs           db.StringList  `db:"ip_addrs"`
	MacAddresses  db.StringList  `db:"mac_addrs"`
	LastUpdatedOn time.Time      `db:"last_updated_on"`
}

type WFCAccountQuery struct {
	Serial string
	FC     int64
	IP     string
	MAC    string
}

func (w *WFCAccountRepo) GetWFCID(query WFCAccountQuery) (int64, error) {

	wfcid := int64(0)
	err := w.sql.Get(&wfcid, sq.Select("wfc_id").From("wfc_accounts").Where(
		`JSON_CONTAINS(console_sns, JSON_QUOTE(?)) 
			OR JSON_CONTAINS(console_fcs, JSON_QUOTE(?)) 
			OR JSON_CONTAINS(ip_addrs, JSON_QUOTE(?))
			OR JSON_CONTAINS(mac_addrs, JSON_QUOTE(?))
		`,
		query.Serial, query.FC, query.IP, query.MAC,
	))

	if err != nil {
		w.logger.Warning("Query", "Failed to get WFCID (%+v)", query)
		return 0, err
	}

	return wfcid, nil
}

func (w *WFCAccountRepo) Insert(query WFCAccountQuery) (int64 /*wfc_id*/, error) {
	err := w.sql.Insert(
		sq.Insert("wfc_accounts").
			Columns("console_sns", "console_fcs", "ip_addrs", "mac_addrs", "last_updated_on").
			Values(
				db.StringList{query.Serial}, db.IntegerList{query.FC},
				db.StringList{query.IP}, db.StringList{query.MAC}, time.Now(),
			),
	)

	if err != nil {
		return -1, nil
	}

	wfcid := int64(0)
	err = w.sql.Get(&wfcid, sq.Select("wfc_id").From("wfc_accounts").Where(
		"JSON_CONTAINS(console_ids, JSON_QUOTE(?))", query.Serial,
	))

	return wfcid, err
}
