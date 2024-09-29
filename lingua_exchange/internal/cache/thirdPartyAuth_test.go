package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newThirdPartyAuthCache() *gotest.Cache {
	record1 := &model.ThirdPartyAuth{}
	record1.ID = 1
	record2 := &model.ThirdPartyAuth{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewThirdPartyAuthCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_thirdPartyAuthCache_Set(t *testing.T) {
	c := newThirdPartyAuthCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ThirdPartyAuth)
	err := c.ICache.(ThirdPartyAuthCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ThirdPartyAuthCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_thirdPartyAuthCache_Get(t *testing.T) {
	c := newThirdPartyAuthCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ThirdPartyAuth)
	err := c.ICache.(ThirdPartyAuthCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ThirdPartyAuthCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ThirdPartyAuthCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_thirdPartyAuthCache_MultiGet(t *testing.T) {
	c := newThirdPartyAuthCache()
	defer c.Close()

	var testData []*model.ThirdPartyAuth
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.ThirdPartyAuth))
	}

	err := c.ICache.(ThirdPartyAuthCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ThirdPartyAuthCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.ThirdPartyAuth))
	}
}

func Test_thirdPartyAuthCache_MultiSet(t *testing.T) {
	c := newThirdPartyAuthCache()
	defer c.Close()

	var testData []*model.ThirdPartyAuth
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.ThirdPartyAuth))
	}

	err := c.ICache.(ThirdPartyAuthCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_thirdPartyAuthCache_Del(t *testing.T) {
	c := newThirdPartyAuthCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ThirdPartyAuth)
	err := c.ICache.(ThirdPartyAuthCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_thirdPartyAuthCache_SetCacheWithNotFound(t *testing.T) {
	c := newThirdPartyAuthCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ThirdPartyAuth)
	err := c.ICache.(ThirdPartyAuthCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewThirdPartyAuthCache(t *testing.T) {
	c := NewThirdPartyAuthCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewThirdPartyAuthCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewThirdPartyAuthCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
