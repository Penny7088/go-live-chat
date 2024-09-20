package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newLanguagesCache() *gotest.Cache {
	record1 := &model.Languages{}
	record1.ID = 1
	record2 := &model.Languages{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLanguagesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_languagesCache_Set(t *testing.T) {
	c := newLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Languages)
	err := c.ICache.(LanguagesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LanguagesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_languagesCache_Get(t *testing.T) {
	c := newLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Languages)
	err := c.ICache.(LanguagesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LanguagesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LanguagesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_languagesCache_MultiGet(t *testing.T) {
	c := newLanguagesCache()
	defer c.Close()

	var testData []*model.Languages
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Languages))
	}

	err := c.ICache.(LanguagesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LanguagesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Languages))
	}
}

func Test_languagesCache_MultiSet(t *testing.T) {
	c := newLanguagesCache()
	defer c.Close()

	var testData []*model.Languages
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Languages))
	}

	err := c.ICache.(LanguagesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_languagesCache_Del(t *testing.T) {
	c := newLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Languages)
	err := c.ICache.(LanguagesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_languagesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Languages)
	err := c.ICache.(LanguagesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLanguagesCache(t *testing.T) {
	c := NewLanguagesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLanguagesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLanguagesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
