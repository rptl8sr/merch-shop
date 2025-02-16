package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"merch-shop/internal/api"
)

type mockInfoRepository struct {
	mock.Mock
}

func (m *mockInfoRepository) GetInfo(ctx context.Context, userID int) (*api.InfoResponse, error) {
	args := m.Called(ctx, userID)
	if res, ok := args.Get(0).(*api.InfoResponse); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func PtrInt(value int) *int {
	return &value
}

func PtrString(value string) *string {
	return &value
}

func TestInfoService_GetUserInfo(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockInfoRepository)
	infoService := NewInfoService(mockRepo)

	testCases := []struct {
		name         string
		userID       uint
		mockReturn   *api.InfoResponse
		mockError    error
		expectedResp *api.InfoResponse
		expectedErr  string
	}{
		{
			name:   "Success: valid user with coins and inventory",
			userID: 1,
			mockReturn: &api.InfoResponse{
				Coins: PtrInt(120),
				Inventory: &[]struct {
					Quantity *int    `json:"quantity,omitempty"`
					Type     *string `json:"type,omitempty"`
				}{
					{Quantity: PtrInt(2), Type: PtrString("sword")},
					{Quantity: PtrInt(1), Type: PtrString("shield")},
				},
			},
			mockError: nil,
			expectedResp: &api.InfoResponse{
				Coins: PtrInt(120),
				Inventory: &[]struct {
					Quantity *int    `json:"quantity,omitempty"`
					Type     *string `json:"type,omitempty"`
				}{
					{Quantity: PtrInt(2), Type: PtrString("sword")},
					{Quantity: PtrInt(1), Type: PtrString("shield")},
				},
			},
			expectedErr: "",
		},
		{
			name:         "Error: repository failure",
			userID:       2,
			mockReturn:   nil,
			mockError:    errors.New("repository error"),
			expectedResp: nil,
			expectedErr:  "repository error",
		},
		{
			name:   "Success: user with empty data",
			userID: 3,
			mockReturn: &api.InfoResponse{
				Coins: PtrInt(0),
				Inventory: &[]struct {
					Quantity *int    `json:"quantity,omitempty"`
					Type     *string `json:"type,omitempty"`
				}{},
			},
			mockError: nil,
			expectedResp: &api.InfoResponse{
				Coins: PtrInt(0),
				Inventory: &[]struct {
					Quantity *int    `json:"quantity,omitempty"`
					Type     *string `json:"type,omitempty"`
				}{},
			},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetInfo", ctx, int(tc.userID)).Return(tc.mockReturn, tc.mockError)
			result, err := infoService.GetUserInfo(ctx, tc.userID)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
			}
			assert.Equal(t, tc.expectedResp, result)
			mockRepo.AssertCalled(t, "GetInfo", ctx, int(tc.userID))
		})
	}
}
