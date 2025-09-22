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
	err := w.sql.Get(&wfcid, sq.Select("wfc_id").From("wfc_accounts").Where(sq.Or{
		sq.Expr("JSON_SEARCH(console_sns, 'one', ?) IS NOT NULL", query.Serial),
		sq.Expr("JSON_SEARCH(console_fcs, 'one', ?) IS NOT NULL", query.FC),
		sq.Expr("JSON_SEARCH(ip_addrs, 'one', ?) IS NOT NULL", query.IP),
		sq.Expr("JSON_SEARCH(mac_addrs, 'one', ?) IS NOT NULL", query.MAC),
	}))

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

func (w *WFCAccountRepo) Update(query WFCAccountQuery) error {

	queries := []func(sq.UpdateBuilder) sq.UpdateBuilder{
		func(ub sq.UpdateBuilder) sq.UpdateBuilder {
			return ub.Set("console_sns", sq.Expr("JSON_ARRAY_APPEND(console_sns, '$', ?)", query.Serial)).Where(sq.Or{
				sq.Expr("NOT JSON_CONTAINS(console_sns, JSON_QUOTE(?), '$')", query.Serial),
			})
		},

		func(ub sq.UpdateBuilder) sq.UpdateBuilder {
			return ub.Set("console_fcs", sq.Expr("JSON_ARRAY_APPEND(console_fcs, '$', ?)", query.FC)).Where(sq.Or{
				sq.Expr("NOT JSON_CONTAINS(console_fcs, CAST(? AS JSON), '$')", query.FC),
			})
		},

		func(ub sq.UpdateBuilder) sq.UpdateBuilder {
			return ub.Set("ip_addrs", sq.Expr("JSON_ARRAY_APPEND(ip_addrs, '$', ?)", query.IP)).Where(sq.Or{
				sq.Expr("NOT JSON_CONTAINS(ip_addrs, JSON_QUOTE(?), '$')", query.IP),
			})
		},

		func(ub sq.UpdateBuilder) sq.UpdateBuilder {
			return ub.Set("mac_addrs", sq.Expr("JSON_ARRAY_APPEND(mac_addrs, '$', ?)", query.MAC)).Where(sq.Or{
				sq.Expr("NOT JSON_CONTAINS(mac_addrs, JSON_QUOTE(?), '$')", query.MAC),
			})
		},
	}

	for _, query := range queries {
		err := w.sql.Update(query(sq.Update("wfc_accounts")))
		if err != nil {
			return err
		}
	}

	return nil
}
