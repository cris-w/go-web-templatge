package power

import (
	"context"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (Repository, func()) {
	db := common.SetupTestDB(t)

	// 迁移表结构
	err := db.AutoMigrate(&PowerSupply{})
	require.NoError(t, err)

	repo := NewRepository(db)

	cleanup := func() {
		common.TeardownTestDB(t, db)
	}

	return repo, cleanup
}

func TestRepository_Create(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		ps      *PowerSupply
		wantErr bool
	}{
		{
			name: "成功创建电源",
			ps: &PowerSupply{
				Name:        "Corsair RM850x",
				Brand:       "Corsair",
				Model:       "RM850x",
				Power:       850,
				Efficiency:  "80Plus Gold",
				Modular:     true,
				Price:       899.99,
				Stock:       100,
				Description: "高品质全模组电源",
				Status:      1,
			},
			wantErr: false,
		},
		{
			name: "创建最小信息电源",
			ps: &PowerSupply{
				Name:   "Basic PSU",
				Power:  500,
				Price:  299.99,
				Status: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.ps)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.ps.ID)
			}
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试数据
	ps := &PowerSupply{
		Name:   "Test PSU",
		Power:  750,
		Price:  699.99,
		Status: 1,
	}
	err := repo.Create(ctx, ps)
	require.NoError(t, err)

	tests := []struct {
		name    string
		psID    uint
		wantErr bool
	}{
		{
			name:    "成功查找电源",
			psID:    ps.ID,
			wantErr: false,
		},
		{
			name:    "电源不存在",
			psID:    9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundPS, err := repo.FindByID(ctx, tt.psID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, foundPS)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundPS)
				assert.Equal(t, tt.psID, foundPS.ID)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试数据
	ps := &PowerSupply{
		Name:   "Test PSU",
		Brand:  "TestBrand",
		Power:  750,
		Price:  699.99,
		Stock:  50,
		Status: 1,
	}
	err := repo.Create(ctx, ps)
	require.NoError(t, err)

	tests := []struct {
		name      string
		updates   map[string]interface{}
		wantErr   bool
		checkFunc func(*testing.T)
	}{
		{
			name: "更新价格",
			updates: map[string]interface{}{
				"price": 799.99,
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updatedPS, err := repo.FindByID(ctx, ps.ID)
				require.NoError(t, err)
				assert.Equal(t, 799.99, updatedPS.Price)
			},
		},
		{
			name: "更新多个字段",
			updates: map[string]interface{}{
				"stock":  100,
				"status": 0,
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updatedPS, err := repo.FindByID(ctx, ps.ID)
				require.NoError(t, err)
				assert.Equal(t, 100, updatedPS.Stock)
				assert.Equal(t, 0, updatedPS.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, ps, tt.updates)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t)
				}
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() uint
		wantErr bool
	}{
		{
			name: "成功删除电源",
			setup: func() uint {
				ps := &PowerSupply{
					Name:   "Delete PSU",
					Power:  600,
					Price:  499.99,
					Status: 1,
				}
				err := repo.Create(ctx, ps)
				require.NoError(t, err)
				return ps.ID
			},
			wantErr: false,
		},
		{
			name: "删除不存在的电源",
			setup: func() uint {
				return 9999
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			psID := tt.setup()
			err := repo.Delete(ctx, psID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证电源已被删除
				_, err := repo.FindByID(ctx, psID)
				assert.Error(t, err)
			}
		})
	}
}

func TestRepository_List(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试数据
	powerSupplies := []*PowerSupply{
		{Name: "Corsair RM850x", Brand: "Corsair", Power: 850, Price: 899.99, Efficiency: "80Plus Gold", Status: 1},
		{Name: "Corsair RM750x", Brand: "Corsair", Power: 750, Price: 799.99, Efficiency: "80Plus Gold", Status: 1},
		{Name: "EVGA SuperNOVA", Brand: "EVGA", Power: 650, Price: 599.99, Efficiency: "80Plus Bronze", Status: 1},
		{Name: "Seasonic Prime", Brand: "Seasonic", Power: 1000, Price: 1299.99, Efficiency: "80Plus Platinum", Status: 0},
	}

	for _, ps := range powerSupplies {
		err := repo.Create(ctx, ps)
		require.NoError(t, err)
	}

	t.Run("查询所有电源", func(t *testing.T) {
		ps, err := repo.List(ctx, nil)
		assert.NoError(t, err)
		assert.Len(t, ps, 4)
	})
}

func TestRepository_Count(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试数据
	powerSupplies := []*PowerSupply{
		{Name: "Corsair RM850x", Brand: "Corsair", Power: 850, Price: 899.99, Status: 1},
		{Name: "Corsair RM750x", Brand: "Corsair", Power: 750, Price: 799.99, Status: 1},
		{Name: "EVGA SuperNOVA", Brand: "EVGA", Power: 650, Price: 599.99, Status: 0},
	}

	for _, ps := range powerSupplies {
		err := repo.Create(ctx, ps)
		require.NoError(t, err)
	}

	t.Run("统计所有电源", func(t *testing.T) {
		count, err := repo.Count(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestRepository_Exists(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试数据
	ps := &PowerSupply{
		Name:   "Test PSU",
		Power:  750,
		Price:  699.99,
		Status: 1,
	}
	err := repo.Create(ctx, ps)
	require.NoError(t, err)

	tests := []struct {
		name       string
		opts       []common.QueryOption
		wantExists bool
	}{
		{
			name:       "电源存在",
			opts:       []common.QueryOption{common.Where("name", "Test PSU")},
			wantExists: true,
		},
		{
			name:       "电源不存在",
			opts:       []common.QueryOption{common.Where("name", "Not Exist")},
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.Exists(ctx, tt.opts...)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantExists, exists)
		})
	}
}
