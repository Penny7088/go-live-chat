package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newInterestsCache() *gotest.Cache {
	record1 := &model.Interests{}
	record1.ID = 1
	record2 := &model.Interests{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewInterestsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_interestsCache_Set(t *testing.T) {
	c := newInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Interests)
	err := c.ICache.(InterestsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(InterestsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_interestsCache_Get(t *testing.T) {
	c := newInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Interests)
	err := c.ICache.(InterestsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(InterestsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(InterestsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_interestsCache_MultiGet(t *testing.T) {
	c := newInterestsCache()
	defer c.Close()

	var testData []*model.Interests
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Interests))
	}

	err := c.ICache.(InterestsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(InterestsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Interests))
	}
}

func Test_interestsCache_MultiSet(t *testing.T) {
	c := newInterestsCache()
	defer c.Close()

	var testData []*model.Interests
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Interests))
	}

	err := c.ICache.(InterestsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsCache_Del(t *testing.T) {
	c := newInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Interests)
	err := c.ICache.(InterestsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsCache_SetCacheWithNotFound(t *testing.T) {
	c := newInterestsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Interests)
	err := c.ICache.(InterestsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInterestsCache(t *testing.T) {
	c := NewInterestsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewInterestsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewInterestsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
