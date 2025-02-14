package repository

import (
	"context"
	"database/sql"
	"errors"
	internalErrors "merch-shop/internal/errors"

	"merch-shop/pkg/database"
	"merch-shop/pkg/logger"
)

type MerchDTO struct {
	ID       uint
	ItemName string
	Price    int
}

type MerchRepository struct {
	db database.DB
}

func NewMerchRepository(db database.DB) *MerchRepository {
	return &MerchRepository{
		db: db,
	}
}

func (m *MerchRepository) GetMerchList(ctx context.Context) ([]MerchDTO, error) {
	logger.Debug("MerchRepository.GetMerchList: ", "message", "retrieving merch list")

	query := "select id, item_name, price from merch_items"

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("MerchRepository.GetMerchList: ", "message", "no merch items found")
			return nil, internalErrors.ErrNoMerchItemsFound
		}
		logger.Error("MerchRepository.GetMerchList: ", "message", "query execution error", "error", err)
		return nil, internalErrors.ErrMerchItemsGettingFailed
	}
	defer rows.Close()

	var merchList []MerchDTO
	for rows.Next() {
		var dto MerchDTO

		e := rows.Scan(&dto.ID, &dto.ItemName, &dto.Price)
		if e != nil {
			logger.Error("MerchRepository.GetMerchList: ", "message", "scan error", "error", e)
			return nil, internalErrors.ErrMerchItemScan
		}

		merchList = append(merchList, dto)
	}

	if err = rows.Err(); err != nil {
		logger.Error("MerchRepository.GetMerchList: ", "message", "rows error", "error", err)
		return nil, internalErrors.ErrMerchItemsGettingFailed
	}

	logger.Debug("MerchRepository.GetMerchList: ", "message", "merch list retrieved", "count", len(merchList))
	return merchList, nil
}
