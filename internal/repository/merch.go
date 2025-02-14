package repository

import (
	"context"
	"database/sql"
	"errors"
	"merch-shop/pkg/logger"

	"merch-shop/pkg/database"
)

type MerchRepository struct {
	db database.DB
}

func NewMerchRepository(db database.DB) *MerchRepository {
	return &MerchRepository{
		db: db,
	}
}

func (m *MerchRepository) GetMerchList(ctx context.Context) (map[string]int, error) {
	logger.Debug("MerchRepository.GetMerchList: ", "message", "retrieving merch list")

	query := "select item_name, price from merch_items"

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("MerchRepository.GetMerchList: ", "message", "no merch items found")
			return nil, ErrNoMerchItemsFound
		}
		logger.Error("MerchRepository.GetMerchList: ", "message", "query execution error", "error", err)
		return nil, ErrMerchItemsGettingFailed
	}
	defer rows.Close()

	merchList := make(map[string]int)
	for rows.Next() {
		var itemName string
		var price int

		e := rows.Scan(&itemName, &price)
		if e != nil {
			logger.Error("MerchRepository.GetMerchList: ", "message", "scan error", "error", e)
			return nil, ErrMerchItemScan
		}

		merchList[itemName] = price
	}

	if err = rows.Err(); err != nil {
		logger.Error("MerchRepository.GetMerchList: ", "message", "rows error", "error", err)
		return nil, ErrMerchItemsGettingFailed
	}

	logger.Debug("MerchRepository.GetMerchList: ", "message", "merch list retrieved", "count", len(merchList))
	return merchList, nil
}
