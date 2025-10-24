package user

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
	err := db.AutoMigrate(&User{})
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
		user    *User
		wantErr bool
	}{
		{
			name: "成功创建用户",
			user: &User{
				Username: "testuser",
				Password: "hashedpassword",
				Email:    "test@example.com",
				Nickname: "Test User",
				Status:   1,
			},
			wantErr: false,
		},
		{
			name: "创建重复用户名",
			user: &User{
				Username: "testuser",
				Password: "hashedpassword",
				Email:    "test2@example.com",
				Status:   1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.user.ID)
			}
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试用户
	user := &User{
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "test@example.com",
		Status:   1,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "成功查找用户",
			userID:  user.ID,
			wantErr: false,
		},
		{
			name:    "用户不存在",
			userID:  9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := repo.FindByID(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
				assert.Equal(t, tt.userID, foundUser.ID)
			}
		})
	}
}

func TestRepository_FindByUsername(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试用户
	user := &User{
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "test@example.com",
		Status:   1,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "成功查找用户",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "用户不存在",
			username: "notexist",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := repo.FindByUsername(ctx, tt.username)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
				assert.Equal(t, tt.username, foundUser.Username)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试用户
	user := &User{
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "test@example.com",
		Status:   1,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	tests := []struct {
		name      string
		updates   map[string]interface{}
		wantErr   bool
		checkFunc func(*testing.T)
	}{
		{
			name: "更新邮箱",
			updates: map[string]interface{}{
				"email": "newemail@example.com",
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updatedUser, err := repo.FindByID(ctx, user.ID)
				require.NoError(t, err)
				assert.Equal(t, "newemail@example.com", updatedUser.Email)
			},
		},
		{
			name: "更新多个字段",
			updates: map[string]interface{}{
				"nickname": "New Nickname",
				"phone":    "1234567890",
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updatedUser, err := repo.FindByID(ctx, user.ID)
				require.NoError(t, err)
				assert.Equal(t, "New Nickname", updatedUser.Nickname)
				assert.Equal(t, "1234567890", updatedUser.Phone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, user, tt.updates)

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
			name: "成功删除用户",
			setup: func() uint {
				user := &User{
					Username: "deleteuser",
					Password: "password",
					Email:    "delete@example.com",
					Status:   1,
				}
				err := repo.Create(ctx, user)
				require.NoError(t, err)
				return user.ID
			},
			wantErr: false,
		},
		{
			name: "删除不存在的用户",
			setup: func() uint {
				return 9999
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setup()
			err := repo.Delete(ctx, userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 验证用户已被删除
				_, err := repo.FindByID(ctx, userID)
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
	users := []*User{
		{Username: "user1", Password: "pass1", Email: "user1@example.com", Status: 1},
		{Username: "user2", Password: "pass2", Email: "user2@example.com", Status: 1},
		{Username: "user3", Password: "pass3", Email: "user3@example.com", Status: 0},
		{Username: "testuser", Password: "pass4", Email: "test@example.com", Status: 1},
	}

	for _, u := range users {
		err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	t.Run("查询所有用户", func(t *testing.T) {
		result, err := repo.List(ctx, nil)
		assert.NoError(t, err)
		assert.Len(t, result, 4)
	})
}

func TestRepository_Count(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试数据
	users := []*User{
		{Username: "user1", Password: "pass1", Email: "user1@example.com", Status: 1},
		{Username: "user2", Password: "pass2", Email: "user2@example.com", Status: 1},
		{Username: "admin1", Password: "pass3", Email: "admin1@example.com", Status: 0},
	}

	for _, u := range users {
		err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	t.Run("统计所有用户", func(t *testing.T) {
		count, err := repo.Count(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestRepository_Exists(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// 创建测试用户
	user := &User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Status:   1,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	tests := []struct {
		name       string
		opts       []common.QueryOption
		wantExists bool
	}{
		{
			name:       "用户存在",
			opts:       []common.QueryOption{common.Where("username", "testuser")},
			wantExists: true,
		},
		{
			name:       "用户不存在",
			opts:       []common.QueryOption{common.Where("username", "notexist")},
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
