package repo

import (
	"context"
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPowerRepository_Create(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	// 迁移表结构
	err := power.AutoMigrate(db)
	require.NoError(t, err)

	repo := NewPowerRepository(db)
	ctx := context.Background()

	t.Run("成功创建电源", func(t *testing.T) {
		ps := &power.PowerSupply{
			Name:       "Test Power Supply",
			Brand:      "Test Brand",
			Model:      "TP-1000",
			Power:      1000,
			Efficiency: "80Plus Gold",
			Modular:    true,
			Price:      299.99,
			Stock:      10,
			Status:     1,
		}

		err := repo.Create(ctx, ps)
		assert.NoError(t, err)
		assert.NotZero(t, ps.ID)
	})
}

func TestPowerRepository_FindByID(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	repo := NewPowerRepository(db)
	ctx := context.Background()

	// 创建测试电源
	ps := &power.PowerSupply{
		Name:   "Find Power Supply",
		Brand:  "Test Brand",
		Power:  850,
		Price:  199.99,
		Status: 1,
	}
	err = repo.Create(ctx, ps)
	require.NoError(t, err)

	t.Run("成功查找电源", func(t *testing.T) {
		found, err := repo.FindByID(ctx, ps.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, ps.Name, found.Name)
		assert.Equal(t, ps.Power, found.Power)
	})

	t.Run("查找不存在的电源", func(t *testing.T) {
		found, err := repo.FindByID(ctx, 99999)
		assert.Error(t, err)
		assert.Nil(t, found)
		assert.True(t, common.IsAppError(err))
	})
}

func TestPowerRepository_Update(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	repo := NewPowerRepository(db)
	ctx := context.Background()

	// 创建测试电源
	ps := &power.PowerSupply{
		Name:   "Update Power Supply",
		Brand:  "Old Brand",
		Power:  500,
		Price:  99.99,
		Stock:  5,
		Status: 1,
	}
	err = repo.Create(ctx, ps)
	require.NoError(t, err)

	t.Run("成功更新电源", func(t *testing.T) {
		updates := map[string]any{
			"brand": "New Brand",
			"power": 750,
			"price": 149.99,
		}

		err := repo.Update(ctx, ps, updates)
		assert.NoError(t, err)

		// 验证更新
		updated, err := repo.FindByID(ctx, ps.ID)
		assert.NoError(t, err)
		assert.Equal(t, "New Brand", updated.Brand)
		assert.Equal(t, 750, updated.Power)
		assert.Equal(t, 149.99, updated.Price)
	})
}

func TestPowerRepository_Delete(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	repo := NewPowerRepository(db)
	ctx := context.Background()

	// 创建测试电源
	ps := &power.PowerSupply{
		Name:   "Delete Power Supply",
		Power:  600,
		Price:  89.99,
		Status: 1,
	}
	err = repo.Create(ctx, ps)
	require.NoError(t, err)

	t.Run("成功删除电源", func(t *testing.T) {
		err := repo.Delete(ctx, ps.ID)
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.FindByID(ctx, ps.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestPowerRepository_List(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	repo := NewPowerRepository(db)
	ctx := context.Background()

	// 创建多个测试电源
	powerSupplies := []*power.PowerSupply{
		{Name: "PSU 500W", Brand: "Brand A", Power: 500, Price: 99.99, Status: 1},
		{Name: "PSU 750W", Brand: "Brand B", Power: 750, Price: 149.99, Status: 1},
		{Name: "PSU 1000W", Brand: "Brand A", Power: 1000, Price: 199.99, Status: 0},
	}

	for _, ps := range powerSupplies {
		err := repo.Create(ctx, ps)
		require.NoError(t, err)
	}

	t.Run("查询所有电源", func(t *testing.T) {
		query := &power.QueryOptions{}
		psList, err := repo.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(psList), 3)
	})

	t.Run("按品牌筛选", func(t *testing.T) {
		query := &power.QueryOptions{
			Brand: "Brand A",
		}
		psList, err := repo.List(ctx, query)
		assert.NoError(t, err)
		for _, ps := range psList {
			assert.Equal(t, "Brand A", ps.Brand)
		}
	})

	t.Run("按功率范围筛选", func(t *testing.T) {
		minPower := 600
		maxPower := 900
		query := &power.QueryOptions{
			MinPower: &minPower,
			MaxPower: &maxPower,
		}
		psList, err := repo.List(ctx, query)
		assert.NoError(t, err)
		for _, ps := range psList {
			assert.GreaterOrEqual(t, ps.Power, 600)
			assert.LessOrEqual(t, ps.Power, 900)
		}
	})

	t.Run("按价格范围筛选", func(t *testing.T) {
		minPrice := 100.0
		maxPrice := 200.0
		query := &power.QueryOptions{
			MinPrice: &minPrice,
			MaxPrice: &maxPrice,
		}
		psList, err := repo.List(ctx, query)
		assert.NoError(t, err)
		for _, ps := range psList {
			assert.GreaterOrEqual(t, ps.Price, 100.0)
			assert.LessOrEqual(t, ps.Price, 200.0)
		}
	})

	t.Run("分页查询", func(t *testing.T) {
		query := &power.QueryOptions{
			Page:     1,
			PageSize: 2,
		}
		psList, err := repo.List(ctx, query)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(psList), 2)
	})
}

func TestPowerRepository_Count(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	repo := NewPowerRepository(db)
	ctx := context.Background()

	// 创建测试电源
	powerSupplies := []*power.PowerSupply{
		{Name: "Count PSU 1", Power: 500, Status: 1},
		{Name: "Count PSU 2", Power: 750, Status: 1},
		{Name: "Count PSU 3", Power: 1000, Status: 0},
	}

	for _, ps := range powerSupplies {
		err := repo.Create(ctx, ps)
		require.NoError(t, err)
	}

	t.Run("统计所有电源", func(t *testing.T) {
		count, err := repo.Count(ctx, nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})

	t.Run("按状态统计", func(t *testing.T) {
		status := 1
		query := &power.QueryOptions{
			Status: &status,
		}
		count, err := repo.Count(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2))
	})
}
