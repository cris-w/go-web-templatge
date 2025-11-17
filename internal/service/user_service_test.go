package service

import (
	"context"
	"power-supply-sys/internal/domain/user"
	"power-supply-sys/internal/infra/db"
	"power-supply-sys/internal/infra/repo"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Create(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	t.Run("成功创建用户", func(t *testing.T) {
		req := &user.UserCreateRequest{
			Username: "newuser",
			Password: "password123",
			Email:    "newuser@example.com",
			Nickname: "New User",
		}

		u, err := service.Create(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, req.Username, u.Username)
		assert.Equal(t, req.Email, u.Email)
		assert.NotEqual(t, req.Password, u.Password) // 密码应该被加密
		assert.Equal(t, 1, u.Status)

		// 验证密码已加密
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
		assert.NoError(t, err)
	})

	t.Run("创建重复用户名失败", func(t *testing.T) {
		req1 := &user.UserCreateRequest{
			Username: "duplicate_user",
			Password: "password123",
		}
		_, err := service.Create(ctx, req1)
		require.NoError(t, err)

		req2 := &user.UserCreateRequest{
			Username: "duplicate_user",
			Password: "password456",
		}
		_, err = service.Create(ctx, req2)
		assert.Error(t, err)
		assert.True(t, common.IsAppError(err))
	})
}

func TestUserService_GetByID(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	// 创建测试用户
	req := &user.UserCreateRequest{
		Username: "getbyid",
		Password: "password123",
		Email:    "getbyid@example.com",
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功获取用户", func(t *testing.T) {
		u, err := service.GetByID(ctx, created.ID)
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, created.ID, u.ID)
		assert.Equal(t, created.Username, u.Username)
	})

	t.Run("获取不存在的用户", func(t *testing.T) {
		u, err := service.GetByID(ctx, 99999)
		assert.Error(t, err)
		assert.Nil(t, u)
	})
}

func TestUserService_GetByUsername(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	// 创建测试用户
	req := &user.UserCreateRequest{
		Username: "getbyusername",
		Password: "password123",
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功通过用户名获取", func(t *testing.T) {
		u, err := service.GetByUsername(ctx, "getbyusername")
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, created.ID, u.ID)
		assert.Equal(t, created.Username, u.Username)
	})

	t.Run("获取不存在的用户名", func(t *testing.T) {
		u, err := service.GetByUsername(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, u)
	})
}

func TestUserService_Update(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	// 创建测试用户
	req := &user.UserCreateRequest{
		Username: "updateuser",
		Password: "password123",
		Email:    "old@example.com",
		Nickname: "OldNick",
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功更新用户", func(t *testing.T) {
		updateReq := &user.UserUpdateRequest{
			Email:    "new@example.com",
			Nickname: "NewNick",
		}

		updated, err := service.Update(ctx, created.ID, updateReq)
		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "new@example.com", updated.Email)
		assert.Equal(t, "NewNick", updated.Nickname)
	})

	t.Run("更新不存在的用户", func(t *testing.T) {
		updateReq := &user.UserUpdateRequest{
			Email: "test@example.com",
		}
		_, err := service.Update(ctx, 99999, updateReq)
		assert.Error(t, err)
	})
}

func TestUserService_Delete(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	// 创建测试用户
	req := &user.UserCreateRequest{
		Username: "deleteuser",
		Password: "password123",
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功删除用户", func(t *testing.T) {
		err := service.Delete(ctx, created.ID)
		assert.NoError(t, err)

		// 验证已删除
		_, err = service.GetByID(ctx, created.ID)
		assert.Error(t, err)
	})
}

func TestUserService_List(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	// 创建多个测试用户
	users := []*user.UserCreateRequest{
		{Username: "listuser1", Password: "pwd", Email: "list1@test.com"},
		{Username: "listuser2", Password: "pwd", Email: "list2@test.com"},
		{Username: "listuser3", Password: "pwd", Email: "list3@test.com"},
	}

	for _, req := range users {
		_, err := service.Create(ctx, req)
		require.NoError(t, err)
	}

	t.Run("查询所有用户", func(t *testing.T) {
		query := &user.UserQueryRequest{}
		users, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		assert.GreaterOrEqual(t, len(users), 3)
	})

	t.Run("分页查询", func(t *testing.T) {
		query := &user.UserQueryRequest{
			Page:     1,
			PageSize: 2,
		}
		users, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		assert.LessOrEqual(t, len(users), 2)
	})

	t.Run("按用户名筛选", func(t *testing.T) {
		query := &user.UserQueryRequest{
			Username: "listuser",
		}
		users, total, err := service.List(ctx, query)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(3))
		for _, u := range users {
			assert.Contains(t, u.Username, "listuser")
		}
	})
}

func TestUserService_Login(t *testing.T) {
	gormDB := common.SetupTestDB(t)
	defer common.TeardownTestDB(t, gormDB)

	err := db.Migrate(gormDB)
	require.NoError(t, err)

	userRepo := repo.NewUserRepository(gormDB)
	service := NewUserService(userRepo)
	ctx := context.Background()

	// 创建测试用户
	req := &user.UserCreateRequest{
		Username: "loginuser",
		Password: "correctpassword",
	}
	created, err := service.Create(ctx, req)
	require.NoError(t, err)

	t.Run("成功登录", func(t *testing.T) {
		loginReq := &user.LoginRequest{
			Username: "loginuser",
			Password: "correctpassword",
		}

		u, err := service.Login(ctx, loginReq)
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, created.ID, u.ID)
		assert.Equal(t, created.Username, u.Username)
	})

	t.Run("错误密码登录失败", func(t *testing.T) {
		loginReq := &user.LoginRequest{
			Username: "loginuser",
			Password: "wrongpassword",
		}

		u, err := service.Login(ctx, loginReq)
		assert.Error(t, err)
		assert.Nil(t, u)
		assert.True(t, common.IsAppError(err))
	})

	t.Run("不存在的用户登录失败", func(t *testing.T) {
		loginReq := &user.LoginRequest{
			Username: "nonexistent",
			Password: "password",
		}

		u, err := service.Login(ctx, loginReq)
		assert.Error(t, err)
		assert.Nil(t, u)
	})

	t.Run("禁用用户登录失败", func(t *testing.T) {
		// 创建禁用用户
		req := &user.UserCreateRequest{
			Username: "disableduser",
			Password: "password123",
			Email:    "disableduser@example.com",
		}
		created, err := service.Create(ctx, req)
		require.NoError(t, err)

		// 禁用用户
		status := 0
		updateReq := &user.UserUpdateRequest{
			Status: &status,
		}
		_, err = service.Update(ctx, created.ID, updateReq)
		require.NoError(t, err)

		// 尝试登录
		loginReq := &user.LoginRequest{
			Username: "disableduser",
			Password: "password123",
		}

		u, err := service.Login(ctx, loginReq)
		assert.Error(t, err)
		assert.Nil(t, u)
		assert.True(t, common.IsAppError(err))
	})
}

func TestUserService_VerifyPassword(t *testing.T) {
	service := &userService{}

	t.Run("验证正确密码", func(t *testing.T) {
		password := "testpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		err = service.VerifyPassword(string(hashedPassword), password)
		assert.NoError(t, err)
	})

	t.Run("验证错误密码", func(t *testing.T) {
		password := "testpassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		err = service.VerifyPassword(string(hashedPassword), "wrongpassword")
		assert.Error(t, err)
	})
}
