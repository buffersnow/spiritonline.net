package repositories

import (
	"database/sql"
	"time"

	"buffersnow.com/spiritonline/pkg/db"
	"buffersnow.com/spiritonline/pkg/log"

	sq "github.com/Masterminds/squirrel"
)

type IWMasterBlockedRepo struct {
	sql    *db.SQL
	logger log.LoggingFactory
}

type IWMasterBlocked struct {
	BlockID      int64          `db:"block_id"`
	IPAddress    string         `db:"ip_addr"`
	BanMessage   sql.NullString `db:"ban_message"`
	IsActive     bool           `db:"is_active"`
	LastBannedOn time.Time      `db:"last_banned_on"`
}

func (w *IWMasterBlockedRepo) GetByIP(ip string) (*IWMasterBlocked, error) {
	suspension := new(IWMasterBlocked)
	err := w.sql.Get(suspension,
		sq.Select("*").From("iwmaster_blocked").Where(sq.Eq{"ip_addr": ip}),
	)

	return suspension, err
}
