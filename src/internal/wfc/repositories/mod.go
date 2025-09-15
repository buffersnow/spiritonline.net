package repositories

import (
	"buffersnow.com/spiritonline/pkg/db"
	"buffersnow.com/spiritonline/pkg/log"
)

type WFCRepo struct {
	Account    *WFCAccountRepo
	Suspension *WFCSuspensionRepo
}

func New(logger *log.Logger, sql *db.SQL) (*WFCRepo, error) {

	repo := &WFCRepo{
		Account: &WFCAccountRepo{
			sql:    sql,
			logger: logger.Factory("WFC Account Repo"),
		},
		Suspension: &WFCSuspensionRepo{
			sql:    sql,
			logger: logger.Factory("WFC Suspension Repo"),
		},
	}

	return repo, nil
}
