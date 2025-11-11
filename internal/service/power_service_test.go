package service

import (
	"context"
	"power-supply-sys/internal/domain/power"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPowerService_Create(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	service := NewPowerService(db)
	ctx := context.Background()

	t.Run("成功创建电源", func(t *testing.T) {
		req := &power.PowerSupplyCreateRequest{
			Name:        "Test PSU 850W",
			Brand:       "Test Brand",
			Model:       "TP-850",
			Power:       850,
			Efficiency:  "80Plus Gold",
			Modular:     true,
			Price:       199.99,
			Stock:       10,
			Description: "Test power supply",
		}

		ps, err := service.Create(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, ps)
		assert.Equal(t, req.Name, ps.Name)
		assert.Equal(t, req.Brand, ps.Brand)
		assert.Equal(t, req.Power, ps.Power)
		assert.Equal(t, req.Price, ps.Price)
		assert.Equal(t, 1, ps.Status) // 默认状态为1
	})
}

func TestPowerService_GetByID(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	service := NewPowerService(db)
	ctx := context.Background()

	// 创建测试电源
	req := &power.PowerSupplyCreateRequest{
		Name:  "Get PSU",
		Power: 750,
		Price: 149.99,
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功获取电源", func(t *testing.T) {
		ps, err := service.GetByID(ctx, created.ID)
		assert.NoError(t, err)
		assert.NotNil(t, ps)
		assert.Equal(t, created.ID, ps.ID)
		assert.Equal(t, created.Name, ps.Name)
	})

	t.Run("获取不存在的电源", func(t *testing.T) {
		ps, err := service.GetByID(ctx, 99999)
		assert.Error(t, err)
		assert.Nil(t, ps)
	})
}

func TestPowerService_Update(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	service := NewPowerService(db)
	ctx := context.Background()

	// 创建测试电源
	req := &power.PowerSupplyCreateRequest{
		Name:  "Update PSU",
		Brand: "Old Brand",
		Power: 500,
		Price: 99.99,
		Stock: 5,
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功更新电源", func(t *testing.T) {
		newPower := 750
		newPrice := 149.99
		updateReq := &power.PowerSupplyUpdateRequest{
			Brand: "New Brand",
			Power: &newPower,
			Price: &newPrice,
		}

		updated, err := service.Update(ctx, created.ID, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "New Brand", updated.Brand)
		assert.Equal(t, 750, updated.Power)
		assert.Equal(t, 149.99, updated.Price)
	})

	t.Run("部分更新", func(t *testing.T) {
		newStock := 20
		updateReq := &power.PowerSupplyUpdateRequest{
			Stock: &newStock,
		}

		updated, err := service.Update(ctx, created.ID, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, 20, updated.Stock)
		// 其他字段应该保持不变
		assert.Equal(t, "New Brand", updated.Brand)
	})

	t.Run("更新不存在的电源", func(t *testing.T) {
		updateReq := &power.PowerSupplyUpdateRequest{
			Name: "New Name",
		}
		_, err := service.Update(ctx, 99999, updateReq)
		assert.Error(t, err)
	})
}

func TestPowerService_Delete(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	service := NewPowerService(db)
	ctx := context.Background()

	// 创建测试电源
	req := &power.PowerSupplyCreateRequest{
		Name:  "Delete PSU",
		Power: 600,
		Price: 89.99,
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功删除电源", func(t *testing.T) {
		err := service.Delete(ctx, created.ID)
		assert.NoError(t, err)

		// 验证已删除
		_, err = service.GetByID(ctx, created.ID)
		assert.Error(t, err)
	})
}

func TestPowerService_List(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := power.AutoMigrate(db)
	require.NoError(t, err)

	service := NewPowerService(db)
	ctx := context.Background()

	// 创建多个测试电源
	powerSupplies := []*power.PowerSupplyCreateRequest{
		{Name: "PSU 500W", Brand: "Brand A", Power: 500, Price: 99.99},
		{Name: "PSU 750W", Brand: "Brand B", Power: 750, Price: 149.99},
		{Name: "PSU 1000W", Brand: "Brand A", Power: 1000, Price: 199.99},
	}

	for _, req := range powerSupplies {
		_, err := service.Create(ctx, req)
		require.NoError(t, err)
	}

	t.Run("查询所有电源", func(t *testing.T) {
		query := &power.PowerSupplyQueryRequest{}
		psList, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		assert.GreaterOrEqual(t, len(psList), 3)
	})

	t.Run("分页查询", func(t *testing.T) {
		query := &power.PowerSupplyQueryRequest{
			Page:     1,
			PageSize: 2,
		}
		psList, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		assert.LessOrEqual(t, len(psList), 2)
	})

	t.Run("按品牌筛选", func(t *testing.T) {
		query := &power.PowerSupplyQueryRequest{
			Brand: "Brand A",
		}
		psList, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))
		_ = total // 使用 total 避免未使用变量警告
		for _, ps := range psList {
			assert.Equal(t, "Brand A", ps.Brand)
		}
	})

	t.Run("按功率范围筛选", func(t *testing.T) {
		minPower := 600
		maxPower := 900
		query := &power.PowerSupplyQueryRequest{
			MinPower: &minPower,
			MaxPower: &maxPower,
		}
		psList, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		_ = total // 使用 total 避免未使用变量警告
		for _, ps := range psList {
			assert.GreaterOrEqual(t, ps.Power, 600)
			assert.LessOrEqual(t, ps.Power, 900)
		}
	})

	t.Run("按价格范围筛选", func(t *testing.T) {
		minPrice := 100.0
		maxPrice := 200.0
		query := &power.PowerSupplyQueryRequest{
			MinPrice: &minPrice,
			MaxPrice: &maxPrice,
		}
		psList, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		_ = total // 使用 total 避免未使用变量警告
		for _, ps := range psList {
			assert.GreaterOrEqual(t, ps.Price, 100.0)
			assert.LessOrEqual(t, ps.Price, 200.0)
		}
	})

	t.Run("组合筛选", func(t *testing.T) {
		minPower := 500
		brand := "Brand A"
		query := &power.PowerSupplyQueryRequest{
			Brand:    brand,
			MinPower: &minPower,
		}
		psList, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		_ = total // 使用 total 避免未使用变量警告
		for _, ps := range psList {
			assert.Equal(t, "Brand A", ps.Brand)
			assert.GreaterOrEqual(t, ps.Power, 500)
		}
	})
}
