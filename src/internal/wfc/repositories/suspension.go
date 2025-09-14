package repositories

import (
	"database/sql"
	"time"
)

type WFCSuspension struct {
	AuditID      int64          `db:"audit_id"`
	WFCID        int64          `db:"wfc_id"`
	ModeratorID  int64          `db:"moderator_id"`
	BanMessage   sql.NullString `db:"ban_message"`
	BanReason    sql.NullString `db:"ban_reason"`
	LastBannedOn time.Time      `db:"last_banned_on"`
	BanExpiresOn sql.NullTime   `db:"ban_expires_on"`
}
