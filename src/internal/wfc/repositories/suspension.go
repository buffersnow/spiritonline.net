package repositories

import (
	"database/sql"
	"time"

	"buffersnow.com/spiritonline/pkg/db"
	"buffersnow.com/spiritonline/pkg/log"

	sq "github.com/Masterminds/squirrel"
)

type WFCSuspensionRepo struct {
	sql    *db.SQL
	logger log.LoggingFactory
}

type WFCSuspension struct {
	AuditID      int64          `db:"audit_id"`
	WFCID        int64          `db:"wfc_id"`
	ModeratorID  int64          `db:"moderator_id"`
	BanMessage   sql.NullString `db:"ban_message"`
	BanReason    sql.NullString `db:"ban_reason"`
	LastBannedOn time.Time      `db:"last_banned_on"`
	BanExpiresOn sql.NullTime   `db:"ban_expires_on"`
}

func (w *WFCSuspensionRepo) Get(wfcid int64) (*WFCSuspension, error) {
	suspension := new(WFCSuspension)
	err := w.sql.Get(suspension,
		sq.Select("*").From("wfc_suspensions").Where(sq.Eq{"wfc_id": wfcid}),
	)

	return suspension, err
}
