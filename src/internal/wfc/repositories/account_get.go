package repositories

import (
	sq "github.com/Masterminds/squirrel"
)

func (w *WFCAccountRepo) Get(query WFCAccountQuery) (*WFCAccount, error) {

	acc, err := w.GetByCID(query.ConsoleID)
	if err == nil {
		return acc, nil
	}
	w.logger.Warning("Query", "Failed to get WFCAccount by ConsoleID (query.ConsoleID: %s)", query.ConsoleID)

	acc, err = w.GetByNandID(query.NandID)
	if err == nil {
		return acc, nil
	}
	w.logger.Warning("Query", "Failed to get WFCAccount by NandID (query.NandID: %d)", query.NandID)

	acc, err = w.GetByIP(query.IP)
	if err == nil {
		return acc, nil
	}
	w.logger.Warning("Query", "Failed to get WFCAccount by IP (query.IP: %v)", query.IP)

	acc, err = w.GetByMAC(query.MAC)
	if err == nil {
		return acc, nil
	}
	w.logger.Warning("Query", "Failed to get WFCAccount by NandID (query.MAC: %s)", query.MAC)

	return nil, err
}

func (w *WFCAccountRepo) GetByWFCID(wfcid int64) (*WFCAccount, error) {
	account := new(WFCAccount)
	err := w.sql.Get(account,
		sq.Select("*").From("wfc_accounts").Where(sq.Eq{"wfc_id": wfcid}),
	)

	return account, err
}

func (w *WFCAccountRepo) GetByCID(cid string) (*WFCAccount, error) {
	account := new(WFCAccount)
	err := w.sql.Get(account,
		sq.Select("*").From("wfc_accounts").Where(
			"JSON_CONTAINS(console_ids, JSON_QUOTE(?))", cid,
		),
	)

	return account, err
}

func (w *WFCAccountRepo) GetByNandID(cfc int64) (*WFCAccount, error) {
	account := new(WFCAccount)
	err := w.sql.Get(account,
		sq.Select("*").From("wfc_accounts").Where(
			"JSON_CONTAINS(nand_ids, JSON_QUOTE(?))", cfc,
		),
	)

	return account, err
}

func (w *WFCAccountRepo) GetByIP(ip string) (*WFCAccount, error) {
	account := new(WFCAccount)
	err := w.sql.Get(account,
		sq.Select("*").From("wfc_accounts").Where(
			"JSON_CONTAINS(ip_addrs, JSON_QUOTE(?))", ip,
		),
	)

	return account, err
}

func (w *WFCAccountRepo) GetByMAC(mac string) (*WFCAccount, error) {
	account := new(WFCAccount)
	err := w.sql.Get(account,
		sq.Select("*").From("wfc_accounts").Where(
			"JSON_CONTAINS(mac_addrs, JSON_QUOTE(?))", mac,
		),
	)

	return account, err
}
