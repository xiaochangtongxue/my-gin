// Package integration 集成测试
//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiaochangtongxue/my-gin/internal/model"
)

// TestUserCRUD 测试用户增删改查
func TestUserCRUD(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TeardownSuite()

	// 清空表
	suite.TruncateUsers(t)
	defer suite.TruncateUsers(t)

	t.Run("CreateUser", func(t *testing.T) {
		user := suite.CreateTestUser(t)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.UID)
		assert.NotEmpty(t, user.Username)
		assert.NotEmpty(t, user.Mobile)
	})

	t.Run("GetUser", func(t *testing.T) {
		createdUser := suite.CreateTestUser(t)

		var user model.User
		err := suite.DB.Where("id = ?", createdUser.ID).First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, createdUser.Username, user.Username)
		assert.Equal(t, createdUser.Mobile, user.Mobile)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user := suite.CreateTestUser(t)

		updates := map[string]interface{}{
			"username": "updated_username",
		}
		err := suite.DB.Model(user).Updates(updates).Error
		assert.NoError(t, err)

		var updatedUser model.User
		suite.DB.First(&updatedUser, user.ID)
		assert.Equal(t, "updated_username", updatedUser.Username)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		user := suite.CreateTestUser(t)

		err := suite.DB.Delete(user).Error
		assert.NoError(t, err)

		var deletedUser model.User
		err = suite.DB.First(&deletedUser, user.ID).Error
		assert.Error(t, err)
	})
}

// TestUserQuery 测试用户查询
func TestUserQuery(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TeardownSuite()

	suite.TruncateUsers(t)
	defer suite.TruncateUsers(t)

	// 创建测试数据
	for i := 0; i < 10; i++ {
		suite.CreateTestUser(t)
	}

	t.Run("ListUsers", func(t *testing.T) {
		var users []model.User
		err := suite.DB.Limit(5).Find(&users).Error
		assert.NoError(t, err)
		assert.Len(t, users, 5)
	})

	t.Run("CountUsers", func(t *testing.T) {
		var count int64
		err := suite.DB.Model(&model.User{}).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(10), count)
	})
}

// TestUserCacheIntegration 测试用户缓存集成
func TestUserCacheIntegration(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TeardownSuite()

	suite.TruncateUsers(t)
	defer suite.TruncateUsers(t)

	t.Run("CacheUser", func(t *testing.T) {
		user := suite.CreateTestUser(t)

		// 缓存用户信息
		key := "user:id"
		ctx := context.Background()
		err := suite.Cache.Set(ctx, key, user, 300)
		assert.NoError(t, err)

		// 从缓存获取
		var cachedUser model.User
		err = suite.Cache.Get(ctx, key, &cachedUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, cachedUser.Username)
	})

	t.Run("DeleteCache", func(t *testing.T) {
		user := suite.CreateTestUser(t)

		key := "user:delete_test"
		ctx := context.Background()
		suite.Cache.Set(ctx, key, user, 300)

		// 删除缓存
		err := suite.Cache.Del(ctx, key)
		assert.NoError(t, err)

		// 验证缓存已删除
		var cachedUser model.User
		err = suite.Cache.Get(ctx, key, &cachedUser)
		assert.Error(t, err)
	})
}