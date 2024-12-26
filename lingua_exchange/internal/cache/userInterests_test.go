package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newUserInterestsCache() *gotest.Cache {
	record1 := &model.UserInterests{}
	record1.ID = 1
	record2 := &model.UserInterests{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewUserInterestsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_userInterestsCache_Set(t *testing.T) {
	c := newUserInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserInterests)
	err := c.ICache.(UserInterestsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(UserInterestsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_userInterestsCache_Get(t *testing.T) {
	c := newUserInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserInterests)
	err := c.ICache.(UserInterestsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserInterestsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(UserInterestsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_userInterestsCache_MultiGet(t *testing.T) {
	c := newUserInterestsCache()
	defer c.Close()

	var testData []*model.UserInterests
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserInterests))
	}

	err := c.ICache.(UserInterestsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserInterestsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.UserInterests))
	}
}

func Test_userInterestsCache_MultiSet(t *testing.T) {
	c := newUserInterestsCache()
	defer c.Close()

	var testData []*model.UserInterests
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserInterests))
	}

	err := c.ICache.(UserInterestsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userInterestsCache_Del(t *testing.T) {
	c := newUserInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserInterests)
	err := c.ICache.(UserInterestsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userInterestsCache_SetCacheWithNotFound(t *testing.T) {
	c := newUserInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserInterests)
	err := c.ICache.(UserInterestsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewUserInterestsCache(t *testing.T) {
	c := NewUserInterestsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewUserInterestsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewUserInterestsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
