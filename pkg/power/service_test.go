package power

import (
	"context"
	"errors"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository 是 Repository 的 mock 实现
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, powerSupply *PowerSupply) error {
	args := m.Called(ctx, powerSupply)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id uint) (*PowerSupply, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PowerSupply), args.Error(1)
}

func (m *MockRepository) FindOne(ctx context.Context, opts ...common.QueryOption) (*PowerSupply, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PowerSupply), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, powerSupply *PowerSupply, updates map[string]interface{}) error {
	args := m.Called(ctx, powerSupply, updates)
	return args.Error(0)
}

func (m *MockRepository) UpdateByID(ctx context.Context, id uint, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Exists(ctx context.Context, opts ...common.QueryOption) (bool, error) {
	args := m.Called(ctx, opts)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, query *QueryOptions) ([]*PowerSupply, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*PowerSupply), args.Error(1)
}

func (m *MockRepository) Count(ctx context.Context, query *QueryOptions) (int64, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(int64), args.Error(1)
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		req       *PowerSupplyCreateRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *PowerSupply, error)
	}{
		{
			name: "成功创建电源",
			req: &PowerSupplyCreateRequest{
				Name:        "Corsair RM850x",
				Brand:       "Corsair",
				Model:       "RM850x",
				Power:       850,
				Efficiency:  "80Plus Gold",
				Modular:     true,
				Price:       899.99,
				Stock:       100,
				Description: "高品质全模组电源",
			},
			mockSetup: func(m *MockRepository) {
				m.On("Create", ctx, mock.MatchedBy(func(ps *PowerSupply) bool {
					return ps.Name == "Corsair RM850x" && ps.Power == 850
				})).Return(nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, ps)
				assert.Equal(t, "Corsair RM850x", ps.Name)
				assert.Equal(t, "Corsair", ps.Brand)
				assert.Equal(t, 850, ps.Power)
				assert.Equal(t, 899.99, ps.Price)
				assert.Equal(t, true, ps.Modular)
				assert.Equal(t, 1, ps.Status)
			},
		},
		{
			name: "最小化信息创建电源",
			req: &PowerSupplyCreateRequest{
				Name:  "Basic PSU",
				Power: 500,
				Price: 299.99,
			},
			mockSetup: func(m *MockRepository) {
				m.On("Create", ctx, mock.Anything).Return(nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, ps)
				assert.Equal(t, "Basic PSU", ps.Name)
				assert.Equal(t, 500, ps.Power)
			},
		},
		{
			name: "数据库错误",
			req: &PowerSupplyCreateRequest{
				Name:  "Test PSU",
				Power: 600,
				Price: 399.99,
			},
			mockSetup: func(m *MockRepository) {
				m.On("Create", ctx, mock.Anything).Return(errors.New("database error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.Error(t, err)
				assert.Nil(t, ps)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			ps, err := svc.Create(ctx, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, ps, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		psID      uint
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *PowerSupply, error)
	}{
		{
			name: "成功获取电源",
			psID: 1,
			mockSetup: func(m *MockRepository) {
				expectedPS := &PowerSupply{
					ID:    1,
					Name:  "Corsair RM850x",
					Power: 850,
					Price: 899.99,
				}
				m.On("FindByID", ctx, uint(1)).Return(expectedPS, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, ps)
				assert.Equal(t, uint(1), ps.ID)
				assert.Equal(t, "Corsair RM850x", ps.Name)
			},
		},
		{
			name: "电源不存在",
			psID: 999,
			mockSetup: func(m *MockRepository) {
				m.On("FindByID", ctx, uint(999)).Return(nil, common.ErrNotFound("电源"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.Error(t, err)
				assert.Nil(t, ps)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			ps, err := svc.GetByID(ctx, tt.psID)

			if tt.checkFunc != nil {
				tt.checkFunc(t, ps, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	newPower := 1000
	newPrice := 1099.99
	newStock := 50
	newModular := false
	newStatus := 0

	tests := []struct {
		name      string
		psID      uint
		req       *PowerSupplyUpdateRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *PowerSupply, error)
	}{
		{
			name: "成功更新电源",
			psID: 1,
			req: &PowerSupplyUpdateRequest{
				Name:  "Updated PSU",
				Brand: "Updated Brand",
				Power: &newPower,
				Price: &newPrice,
			},
			mockSetup: func(m *MockRepository) {
				existingPS := &PowerSupply{
					ID:    1,
					Name:  "Old PSU",
					Power: 850,
				}
				updatedPS := &PowerSupply{
					ID:    1,
					Name:  "Updated PSU",
					Brand: "Updated Brand",
					Power: 1000,
					Price: 1099.99,
				}
				m.On("FindByID", ctx, uint(1)).Return(existingPS, nil).Once()
				m.On("Update", ctx, existingPS, map[string]interface{}{
					"name":  "Updated PSU",
					"brand": "Updated Brand",
					"power": 1000,
					"price": 1099.99,
				}).Return(nil)
				m.On("FindByID", ctx, uint(1)).Return(updatedPS, nil).Once()
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, ps)
				assert.Equal(t, "Updated PSU", ps.Name)
				assert.Equal(t, 1000, ps.Power)
			},
		},
		{
			name: "更新所有字段",
			psID: 1,
			req: &PowerSupplyUpdateRequest{
				Name:        "Complete Update",
				Brand:       "New Brand",
				Model:       "New Model",
				Power:       &newPower,
				Efficiency:  "80Plus Platinum",
				Modular:     &newModular,
				Price:       &newPrice,
				Stock:       &newStock,
				Description: "New Description",
				Status:      &newStatus,
			},
			mockSetup: func(m *MockRepository) {
				existingPS := &PowerSupply{ID: 1}
				updatedPS := &PowerSupply{
					ID:         1,
					Name:       "Complete Update",
					Power:      1000,
					Price:      1099.99,
					Stock:      50,
					Modular:    false,
					Status:     0,
					Efficiency: "80Plus Platinum",
				}
				m.On("FindByID", ctx, uint(1)).Return(existingPS, nil).Once()
				m.On("Update", ctx, existingPS, mock.MatchedBy(func(updates map[string]interface{}) bool {
					return len(updates) == 10 // 应该更新10个字段
				})).Return(nil)
				m.On("FindByID", ctx, uint(1)).Return(updatedPS, nil).Once()
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, ps)
			},
		},
		{
			name: "电源不存在",
			psID: 999,
			req:  &PowerSupplyUpdateRequest{Name: "Test"},
			mockSetup: func(m *MockRepository) {
				m.On("FindByID", ctx, uint(999)).Return(nil, common.ErrNotFound("电源"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.Error(t, err)
				assert.Nil(t, ps)
			},
		},
		{
			name: "没有更新字段",
			psID: 1,
			req:  &PowerSupplyUpdateRequest{},
			mockSetup: func(m *MockRepository) {
				existingPS := &PowerSupply{
					ID:   1,
					Name: "Test PSU",
				}
				m.On("FindByID", ctx, uint(1)).Return(existingPS, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, ps)
			},
		},
		{
			name: "更新失败",
			psID: 1,
			req:  &PowerSupplyUpdateRequest{Name: "Test"},
			mockSetup: func(m *MockRepository) {
				existingPS := &PowerSupply{ID: 1}
				m.On("FindByID", ctx, uint(1)).Return(existingPS, nil).Once()
				m.On("Update", ctx, existingPS, mock.Anything).Return(errors.New("update error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, ps *PowerSupply, err error) {
				assert.Error(t, err)
				assert.Nil(t, ps)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			ps, err := svc.Update(ctx, tt.psID, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, ps, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		psID      uint
		mockSetup func(*MockRepository)
		wantErr   bool
	}{
		{
			name: "成功删除电源",
			psID: 1,
			mockSetup: func(m *MockRepository) {
				m.On("Delete", ctx, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "删除不存在的电源",
			psID: 999,
			mockSetup: func(m *MockRepository) {
				m.On("Delete", ctx, uint(999)).Return(common.ErrNotFound("电源"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			err := svc.Delete(ctx, tt.psID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	ctx := context.Background()
	minPower := 500
	maxPower := 1000
	minPrice := 300.0
	maxPrice := 1000.0
	status := 1

	tests := []struct {
		name      string
		req       *PowerSupplyQueryRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, []*PowerSupply, int64, error)
	}{
		{
			name: "成功获取电源列表",
			req: &PowerSupplyQueryRequest{
				Page:     1,
				PageSize: 10,
				Name:     "Corsair",
			},
			mockSetup: func(m *MockRepository) {
				powerSupplies := []*PowerSupply{
					{ID: 1, Name: "Corsair RM850x", Power: 850},
					{ID: 2, Name: "Corsair RM750x", Power: 750},
				}
				m.On("Count", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Name == "Corsair" && opts.Page == 1 && opts.PageSize == 10
				})).Return(int64(2), nil)
				m.On("List", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Name == "Corsair"
				})).Return(powerSupplies, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps []*PowerSupply, total int64, err error) {
				assert.NoError(t, err)
				assert.Len(t, ps, 2)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "使用所有过滤条件",
			req: &PowerSupplyQueryRequest{
				Page:       1,
				PageSize:   20,
				Name:       "RM",
				Brand:      "Corsair",
				MinPower:   &minPower,
				MaxPower:   &maxPower,
				MinPrice:   &minPrice,
				MaxPrice:   &maxPrice,
				Efficiency: "80Plus Gold",
				Status:     &status,
			},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Name == "RM" &&
						opts.Brand == "Corsair" &&
						opts.MinPower != nil && *opts.MinPower == 500 &&
						opts.MaxPower != nil && *opts.MaxPower == 1000 &&
						opts.Efficiency == "80Plus Gold" &&
						opts.Status != nil && *opts.Status == 1
				})).Return(int64(5), nil)
				m.On("List", ctx, mock.Anything).Return([]*PowerSupply{}, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps []*PowerSupply, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(5), total)
			},
		},
		{
			name: "使用默认分页参数",
			req:  &PowerSupplyQueryRequest{},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Page == 1 && opts.PageSize == 10
				})).Return(int64(0), nil)
				m.On("List", ctx, mock.Anything).Return([]*PowerSupply{}, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps []*PowerSupply, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "Count 失败",
			req:  &PowerSupplyQueryRequest{},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.Anything).Return(int64(0), errors.New("database error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, ps []*PowerSupply, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, ps)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "List 失败",
			req:  &PowerSupplyQueryRequest{},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.Anything).Return(int64(10), nil)
				m.On("List", ctx, mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, ps []*PowerSupply, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, ps)
			},
		},
		{
			name: "按品牌查询",
			req: &PowerSupplyQueryRequest{
				Brand: "EVGA",
			},
			mockSetup: func(m *MockRepository) {
				powerSupplies := []*PowerSupply{
					{ID: 1, Name: "EVGA SuperNOVA", Brand: "EVGA"},
				}
				m.On("Count", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Brand == "EVGA"
				})).Return(int64(1), nil)
				m.On("List", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Brand == "EVGA"
				})).Return(powerSupplies, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, ps []*PowerSupply, total int64, err error) {
				assert.NoError(t, err)
				assert.Len(t, ps, 1)
				assert.Equal(t, "EVGA", ps[0].Brand)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			ps, total, err := svc.List(ctx, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, ps, total, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
