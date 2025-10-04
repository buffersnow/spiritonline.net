package repositories

import (
	"buffersnow.com/spiritonline/pkg/db"
	"buffersnow.com/spiritonline/pkg/log"
)

type IWMasterRepo struct {
	Blocked *IWMasterBlockedRepo
}

func New(logger *log.Logger, sql *db.SQL) (*IWMasterRepo, error) {

	repo := &IWMasterRepo{
		Blocked: &IWMasterBlockedRepo{
			sql:    sql,
			logger: logger.Factory("IWMaster Blocked Repo"),
		},
	}

	return repo, nil
}
