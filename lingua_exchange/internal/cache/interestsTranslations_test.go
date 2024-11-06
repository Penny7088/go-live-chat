package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newInterestsTranslationsCache() *gotest.Cache {
	record1 := &model.InterestsTranslations{}
	record1.ID = 1
	record2 := &model.InterestsTranslations{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewInterestsTranslationsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_interestsTranslationsCache_Set(t *testing.T) {
	c := newInterestsTranslationsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.InterestsTranslations)
	err := c.ICache.(InterestsTranslationsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(InterestsTranslationsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_interestsTranslationsCache_Get(t *testing.T) {
	c := newInterestsTranslationsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.InterestsTranslations)
	err := c.ICache.(InterestsTranslationsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(InterestsTranslationsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(InterestsTranslationsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_interestsTranslationsCache_MultiGet(t *testing.T) {
	c := newInterestsTranslationsCache()
	defer c.Close()

	var testData []*model.InterestsTranslations
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.InterestsTranslations))
	}

	err := c.ICache.(InterestsTranslationsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(InterestsTranslationsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.InterestsTranslations))
	}
}

func Test_interestsTranslationsCache_MultiSet(t *testing.T) {
	c := newInterestsTranslationsCache()
	defer c.Close()

	var testData []*model.InterestsTranslations
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.InterestsTranslations))
	}

	err := c.ICache.(InterestsTranslationsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsTranslationsCache_Del(t *testing.T) {
	c := newInterestsTranslationsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.InterestsTranslations)
	err := c.ICache.(InterestsTranslationsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsTranslationsCache_SetCacheWithNotFound(t *testing.T) {
	c := newInterestsTranslationsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.InterestsTranslations)
	err := c.ICache.(InterestsTranslationsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInterestsTranslationsCache(t *testing.T) {
	c := NewInterestsTranslationsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewInterestsTranslationsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewInterestsTranslationsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
