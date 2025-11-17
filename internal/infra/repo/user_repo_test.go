package repo

import (
	"context"
	"power-supply-sys/internal/domain/user"
	dbpkg "power-supply-sys/internal/infra/db"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	// 迁移表结构
	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("成功创建用户", func(t *testing.T) {
		u := &user.User{
			Username: "testuser",
			Password: "hashed_password",
			Email:    "test@example.com",
			Status:   1,
		}

		err := repo.Create(ctx, u)
		assert.NoError(t, err)
		assert.NotZero(t, u.ID)
	})

	t.Run("创建重复用户名失败", func(t *testing.T) {
		u1 := &user.User{
			Username: "duplicate",
			Password: "password1",
			Status:   1,
		}
		err := repo.Create(ctx, u1)
		require.NoError(t, err)

		u2 := &user.User{
			Username: "duplicate",
			Password: "password2",
			Status:   1,
		}
		err = repo.Create(ctx, u2)
		assert.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	u := &user.User{
		Username: "findbyid",
		Password: "password",
		Email:    "findbyid@example.com",
		Status:   1,
	}
	err = repo.Create(ctx, u)
	require.NoError(t, err)

	t.Run("成功查找用户", func(t *testing.T) {
		found, err := repo.FindByID(ctx, u.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, u.Username, found.Username)
		assert.Equal(t, u.Email, found.Email)
	})

	t.Run("查找不存在的用户", func(t *testing.T) {
		found, err := repo.FindByID(ctx, 99999)
		assert.Error(t, err)
		assert.Nil(t, found)
		assert.True(t, common.IsAppError(err))
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	u := &user.User{
		Username: "findbyusername",
		Password: "password",
		Email:    "findbyusername@example.com",
		Status:   1,
	}
	err = repo.Create(ctx, u)
	require.NoError(t, err)

	t.Run("成功通过用户名查找", func(t *testing.T) {
		found, err := repo.FindByUsername(ctx, "findbyusername")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, u.Username, found.Username)
	})

	t.Run("查找不存在的用户名", func(t *testing.T) {
		found, err := repo.FindByUsername(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	u := &user.User{
		Username: "updateuser",
		Password: "password",
		Email:    "old@example.com",
		Nickname: "OldNick",
		Status:   1,
	}
	err = repo.Create(ctx, u)
	require.NoError(t, err)

	t.Run("成功更新用户", func(t *testing.T) {
		updates := map[string]any{
			"email":    "new@example.com",
			"nickname": "NewNick",
		}

		err := repo.Update(ctx, u, updates)
		assert.NoError(t, err)

		// 验证更新
		updated, err := repo.FindByID(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", updated.Email)
		assert.Equal(t, "NewNick", updated.Nickname)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	u := &user.User{
		Username: "deleteuser",
		Password: "password",
		Status:   1,
	}
	err = repo.Create(ctx, u)
	require.NoError(t, err)

	t.Run("成功删除用户", func(t *testing.T) {
		err := repo.Delete(ctx, u.ID)
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.FindByID(ctx, u.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("删除不存在的用户", func(t *testing.T) {
		err := repo.Delete(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestUserRepository_List(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建多个测试用户
	users := []*user.User{
		{Username: "user1", Password: "pwd", Email: "user1@test.com", Status: 1},
		{Username: "user2", Password: "pwd", Email: "user2@test.com", Status: 1},
		{Username: "user3", Password: "pwd", Email: "user3@test.com", Status: 0},
	}

	for _, u := range users {
		err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	t.Run("查询所有用户", func(t *testing.T) {
		query := &user.QueryOptions{}
		users, err := repo.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3)
	})

	t.Run("按状态筛选", func(t *testing.T) {
		status := 1
		query := &user.QueryOptions{
			Status: &status,
		}
		users, err := repo.List(ctx, query)
		assert.NoError(t, err)
		for _, u := range users {
			assert.Equal(t, 1, u.Status)
		}
	})

	t.Run("按用户名模糊查询", func(t *testing.T) {
		query := &user.QueryOptions{
			Username: "user",
		}
		users, err := repo.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3)
	})

	t.Run("分页查询", func(t *testing.T) {
		query := &user.QueryOptions{
			Page:     1,
			PageSize: 2,
		}
		users, err := repo.List(ctx, query)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(users), 2)
	})
}

func TestUserRepository_Count(t *testing.T) {
	db := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, db)

	err := dbpkg.Migrate(db)
	require.NoError(t, err)

	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	users := []*user.User{
		{Username: "count1", Password: "pwd", Status: 1, Email: "count1@test.com"},
		{Username: "count2", Password: "pwd", Status: 1, Email: "count2@test.com"},
		{Username: "count3", Password: "pwd", Status: 0, Email: "count3@test.com"},
	}

	for _, u := range users {
		err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	t.Run("统计所有用户", func(t *testing.T) {
		count, err := repo.Count(ctx, nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})

	t.Run("按状态统计", func(t *testing.T) {
		status := 1
		query := &user.QueryOptions{
			Status: &status,
		}
		count, err := repo.Count(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2))
	})
}
