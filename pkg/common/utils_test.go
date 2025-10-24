package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseUintParam(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		paramName string
		paramVal  string
		wantID    uint
		wantErr   bool
	}{
		{
			name:      "有效的ID",
			paramName: "id",
			paramVal:  "123",
			wantID:    123,
			wantErr:   false,
		},
		{
			name:      "ID为0",
			paramName: "id",
			paramVal:  "0",
			wantID:    0,
			wantErr:   false,
		},
		{
			name:      "ID为1",
			paramName: "id",
			paramVal:  "1",
			wantID:    1,
			wantErr:   false,
		},
		{
			name:      "最大uint32值",
			paramName: "id",
			paramVal:  "4294967295",
			wantID:    4294967295,
			wantErr:   false,
		},
		{
			name:      "无效的ID - 负数",
			paramName: "id",
			paramVal:  "-1",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "无效的ID - 字符串",
			paramName: "id",
			paramVal:  "abc",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "无效的ID - 浮点数",
			paramName: "id",
			paramVal:  "12.34",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "无效的ID - 空字符串",
			paramName: "id",
			paramVal:  "",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "无效的ID - 超过uint32最大值",
			paramName: "id",
			paramVal:  "4294967296",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "不同的参数名",
			paramName: "user_id",
			paramVal:  "456",
			wantID:    456,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置URL参数
			c.Params = []gin.Param{
				{Key: tt.paramName, Value: tt.paramVal},
			}

			// 调用函数
			id, err := ParseUintParam(c, tt.paramName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.IsType(t, &AppError{}, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

func TestGetPageInfo(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{
			name:         "正常的分页参数",
			page:         2,
			pageSize:     20,
			wantPage:     2,
			wantPageSize: 20,
		},
		{
			name:         "页码小于1，使用默认值1",
			page:         0,
			pageSize:     20,
			wantPage:     1,
			wantPageSize: 20,
		},
		{
			name:         "页码为负数，使用默认值1",
			page:         -1,
			pageSize:     20,
			wantPage:     1,
			wantPageSize: 20,
		},
		{
			name:         "页大小小于1，使用默认值10",
			page:         1,
			pageSize:     0,
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "页大小为负数，使用默认值10",
			page:         1,
			pageSize:     -5,
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "页大小超过100，限制为100",
			page:         1,
			pageSize:     150,
			wantPage:     1,
			wantPageSize: 100,
		},
		{
			name:         "页大小等于100",
			page:         1,
			pageSize:     100,
			wantPage:     1,
			wantPageSize: 100,
		},
		{
			name:         "页大小等于1",
			page:         1,
			pageSize:     1,
			wantPage:     1,
			wantPageSize: 1,
		},
		{
			name:         "所有参数都无效",
			page:         -1,
			pageSize:     -1,
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "所有参数都为0",
			page:         0,
			pageSize:     0,
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "页码很大",
			page:         10000,
			pageSize:     50,
			wantPage:     10000,
			wantPageSize: 50,
		},
		{
			name:         "边界值测试 - 页大小为99",
			page:         1,
			pageSize:     99,
			wantPage:     1,
			wantPageSize: 99,
		},
		{
			name:         "边界值测试 - 页大小为101",
			page:         1,
			pageSize:     101,
			wantPage:     1,
			wantPageSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize := GetPageInfo(tt.page, tt.pageSize)
			assert.Equal(t, tt.wantPage, gotPage, "page should match")
			assert.Equal(t, tt.wantPageSize, gotPageSize, "pageSize should match")
		})
	}
}

func TestParseUintParam_WithRealRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/users/:id", func(c *gin.Context) {
		id, err := ParseUintParam(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedID     uint
	}{
		{
			name:           "有效请求",
			url:            "/users/123",
			expectedStatus: http.StatusOK,
			expectedID:     123,
		},
		{
			name:           "无效请求",
			url:            "/users/abc",
			expectedStatus: http.StatusBadRequest,
			expectedID:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
