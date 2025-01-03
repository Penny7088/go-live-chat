package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newUserLanguagesCache() *gotest.Cache {
	record1 := &model.UserLanguages{}
	record1.ID = 1
	record2 := &model.UserLanguages{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewUserLanguagesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_userLanguagesCache_Set(t *testing.T) {
	c := newUserLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserLanguages)
	err := c.ICache.(UserLanguagesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(UserLanguagesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_userLanguagesCache_Get(t *testing.T) {
	c := newUserLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserLanguages)
	err := c.ICache.(UserLanguagesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserLanguagesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(UserLanguagesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_userLanguagesCache_MultiGet(t *testing.T) {
	c := newUserLanguagesCache()
	defer c.Close()

	var testData []*model.UserLanguages
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserLanguages))
	}

	err := c.ICache.(UserLanguagesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserLanguagesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.UserLanguages))
	}
}

func Test_userLanguagesCache_MultiSet(t *testing.T) {
	c := newUserLanguagesCache()
	defer c.Close()

	var testData []*model.UserLanguages
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserLanguages))
	}

	err := c.ICache.(UserLanguagesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userLanguagesCache_Del(t *testing.T) {
	c := newUserLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserLanguages)
	err := c.ICache.(UserLanguagesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userLanguagesCache_SetCacheWithNotFound(t *testing.T) {
	c := newUserLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserLanguages)
	err := c.ICache.(UserLanguagesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewUserLanguagesCache(t *testing.T) {
	c := NewUserLanguagesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewUserLanguagesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewUserLanguagesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
