package user

import (
	"context"
	"errors"
	"power-supply-sys/pkg/common"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository 是 Repository 的 mock 实现
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) FindOne(ctx context.Context, opts ...common.QueryOption) (*User, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *User, updates map[string]interface{}) error {
	args := m.Called(ctx, user, updates)
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

func (m *MockRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, query *QueryOptions) ([]*User, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*User), args.Error(1)
}

func (m *MockRepository) Count(ctx context.Context, query *QueryOptions) (int64, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(int64), args.Error(1)
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		req       *UserCreateRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *User, error)
	}{
		{
			name: "成功创建用户",
			req: &UserCreateRequest{
				Username: "testuser",
				Password: "password123",
				Email:    "test@example.com",
				Nickname: "Test User",
			},
			mockSetup: func(m *MockRepository) {
				m.On("FindByUsername", ctx, "testuser").Return(nil, common.ErrNotFound("用户"))
				m.On("Create", ctx, mock.MatchedBy(func(u *User) bool {
					return u.Username == "testuser" && u.Email == "test@example.com"
				})).Return(nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "test@example.com", user.Email)
				assert.NotEmpty(t, user.Password)
				assert.Equal(t, 1, user.Status)
			},
		},
		{
			name: "用户名已存在",
			req: &UserCreateRequest{
				Username: "existuser",
				Password: "password123",
			},
			mockSetup: func(m *MockRepository) {
				existingUser := &User{
					ID:       1,
					Username: "existuser",
				}
				m.On("FindByUsername", ctx, "existuser").Return(existingUser, nil)
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				appErr, ok := err.(*common.AppError)
				assert.True(t, ok)
				assert.Equal(t, common.ErrCodeAlreadyExists, appErr.Code)
			},
		},
		{
			name: "数据库错误",
			req: &UserCreateRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(m *MockRepository) {
				m.On("FindByUsername", ctx, "testuser").Return(nil, common.ErrNotFound("用户"))
				m.On("Create", ctx, mock.Anything).Return(errors.New("database error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			user, err := svc.Create(ctx, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, user, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    uint
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *User, error)
	}{
		{
			name:   "成功获取用户",
			userID: 1,
			mockSetup: func(m *MockRepository) {
				expectedUser := &User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
				}
				m.On("FindByID", ctx, uint(1)).Return(expectedUser, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, uint(1), user.ID)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name:   "用户不存在",
			userID: 999,
			mockSetup: func(m *MockRepository) {
				m.On("FindByID", ctx, uint(999)).Return(nil, common.ErrNotFound("用户"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			user, err := svc.GetByID(ctx, tt.userID)

			if tt.checkFunc != nil {
				tt.checkFunc(t, user, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	status := 0

	tests := []struct {
		name      string
		userID    uint
		req       *UserUpdateRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *User, error)
	}{
		{
			name:   "成功更新用户",
			userID: 1,
			req: &UserUpdateRequest{
				Email:    "newemail@example.com",
				Nickname: "New Nickname",
			},
			mockSetup: func(m *MockRepository) {
				existingUser := &User{
					ID:       1,
					Username: "testuser",
					Email:    "old@example.com",
				}
				updatedUser := &User{
					ID:       1,
					Username: "testuser",
					Email:    "newemail@example.com",
					Nickname: "New Nickname",
				}
				m.On("FindByID", ctx, uint(1)).Return(existingUser, nil).Once()
				m.On("Update", ctx, existingUser, map[string]interface{}{
					"email":    "newemail@example.com",
					"nickname": "New Nickname",
				}).Return(nil)
				m.On("FindByID", ctx, uint(1)).Return(updatedUser, nil).Once()
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "newemail@example.com", user.Email)
			},
		},
		{
			name:   "更新用户状态",
			userID: 1,
			req: &UserUpdateRequest{
				Status: &status,
			},
			mockSetup: func(m *MockRepository) {
				existingUser := &User{
					ID:       1,
					Username: "testuser",
					Status:   1,
				}
				updatedUser := &User{
					ID:       1,
					Username: "testuser",
					Status:   0,
				}
				m.On("FindByID", ctx, uint(1)).Return(existingUser, nil).Once()
				m.On("Update", ctx, existingUser, map[string]interface{}{
					"status": 0,
				}).Return(nil)
				m.On("FindByID", ctx, uint(1)).Return(updatedUser, nil).Once()
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, 0, user.Status)
			},
		},
		{
			name:   "用户不存在",
			userID: 999,
			req:    &UserUpdateRequest{Email: "test@example.com"},
			mockSetup: func(m *MockRepository) {
				m.On("FindByID", ctx, uint(999)).Return(nil, common.ErrNotFound("用户"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
		{
			name:   "没有更新字段",
			userID: 1,
			req:    &UserUpdateRequest{},
			mockSetup: func(m *MockRepository) {
				existingUser := &User{
					ID:       1,
					Username: "testuser",
				}
				m.On("FindByID", ctx, uint(1)).Return(existingUser, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			user, err := svc.Update(ctx, tt.userID, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, user, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    uint
		mockSetup func(*MockRepository)
		wantErr   bool
	}{
		{
			name:   "成功删除用户",
			userID: 1,
			mockSetup: func(m *MockRepository) {
				m.On("Delete", ctx, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "删除不存在的用户",
			userID: 999,
			mockSetup: func(m *MockRepository) {
				m.On("Delete", ctx, uint(999)).Return(common.ErrNotFound("用户"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			err := svc.Delete(ctx, tt.userID)

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
	status := 1

	tests := []struct {
		name      string
		req       *UserQueryRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, []*User, int64, error)
	}{
		{
			name: "成功获取用户列表",
			req: &UserQueryRequest{
				Page:     1,
				PageSize: 10,
				Username: "test",
			},
			mockSetup: func(m *MockRepository) {
				users := []*User{
					{ID: 1, Username: "testuser1"},
					{ID: 2, Username: "testuser2"},
				}
				m.On("Count", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Username == "test" && opts.Page == 1 && opts.PageSize == 10
				})).Return(int64(2), nil)
				m.On("List", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Username == "test" && opts.Page == 1 && opts.PageSize == 10
				})).Return(users, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, users []*User, total int64, err error) {
				assert.NoError(t, err)
				assert.Len(t, users, 2)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "使用默认分页参数",
			req: &UserQueryRequest{
				Status: &status,
			},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.MatchedBy(func(opts *QueryOptions) bool {
					return opts.Page == 1 && opts.PageSize == 10 && opts.Status != nil && *opts.Status == 1
				})).Return(int64(5), nil)
				m.On("List", ctx, mock.Anything).Return([]*User{}, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, users []*User, total int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(5), total)
			},
		},
		{
			name: "Count 失败",
			req:  &UserQueryRequest{},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.Anything).Return(int64(0), errors.New("database error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, users []*User, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.Equal(t, int64(0), total)
			},
		},
		{
			name: "List 失败",
			req:  &UserQueryRequest{},
			mockSetup: func(m *MockRepository) {
				m.On("Count", ctx, mock.Anything).Return(int64(10), nil)
				m.On("List", ctx, mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, users []*User, total int64, err error) {
				assert.Error(t, err)
				assert.Nil(t, users)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			users, total, err := svc.List(ctx, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, users, total, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Login(t *testing.T) {
	ctx := context.Background()

	// 生成测试用的密码哈希
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		req       *LoginRequest
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *User, error)
	}{
		{
			name: "登录成功",
			req: &LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(m *MockRepository) {
				user := &User{
					ID:       1,
					Username: "testuser",
					Password: string(hashedPassword),
					Status:   1,
				}
				m.On("FindByUsername", ctx, "testuser").Return(user, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name: "用户不存在",
			req: &LoginRequest{
				Username: "notexist",
				Password: "password123",
			},
			mockSetup: func(m *MockRepository) {
				m.On("FindByUsername", ctx, "notexist").Return(nil, common.ErrNotFound("用户"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				appErr, ok := err.(*common.AppError)
				assert.True(t, ok)
				assert.Equal(t, common.ErrCodeUnauthorized, appErr.Code)
			},
		},
		{
			name: "密码错误",
			req: &LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockRepository) {
				user := &User{
					ID:       1,
					Username: "testuser",
					Password: string(hashedPassword),
					Status:   1,
				}
				m.On("FindByUsername", ctx, "testuser").Return(user, nil)
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				appErr, ok := err.(*common.AppError)
				assert.True(t, ok)
				assert.Equal(t, common.ErrCodeUnauthorized, appErr.Code)
			},
		},
		{
			name: "用户已被禁用",
			req: &LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(m *MockRepository) {
				user := &User{
					ID:       1,
					Username: "testuser",
					Password: string(hashedPassword),
					Status:   0, // 禁用状态
				}
				m.On("FindByUsername", ctx, "testuser").Return(user, nil)
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				appErr, ok := err.(*common.AppError)
				assert.True(t, ok)
				assert.Equal(t, common.ErrCodeForbidden, appErr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			user, err := svc.Login(ctx, tt.req)

			if tt.checkFunc != nil {
				tt.checkFunc(t, user, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_VerifyPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "密码正确",
			hashedPassword: string(hashedPassword),
			password:       "password123",
			wantErr:        false,
		},
		{
			name:           "密码错误",
			hashedPassword: string(hashedPassword),
			password:       "wrongpassword",
			wantErr:        true,
		},
		{
			name:           "空密码",
			hashedPassword: string(hashedPassword),
			password:       "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			svc := &service{repo: mockRepo}

			err := svc.VerifyPassword(tt.hashedPassword, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetByUsername(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		username  string
		mockSetup func(*MockRepository)
		wantErr   bool
		checkFunc func(*testing.T, *User, error)
	}{
		{
			name:     "成功获取用户",
			username: "testuser",
			mockSetup: func(m *MockRepository) {
				user := &User{
					ID:       1,
					Username: "testuser",
				}
				m.On("FindByUsername", ctx, "testuser").Return(user, nil)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name:     "用户不存在",
			username: "notexist",
			mockSetup: func(m *MockRepository) {
				m.On("FindByUsername", ctx, "notexist").Return(nil, common.ErrNotFound("用户"))
			},
			wantErr: true,
			checkFunc: func(t *testing.T, user *User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			svc := &service{repo: mockRepo}
			user, err := svc.GetByUsername(ctx, tt.username)

			if tt.checkFunc != nil {
				tt.checkFunc(t, user, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
