package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newCountriesCache() *gotest.Cache {
	record1 := &model.Countries{}
	record1.ID = 1
	record2 := &model.Countries{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewCountriesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_countriesCache_Set(t *testing.T) {
	c := newCountriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Countries)
	err := c.ICache.(CountriesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(CountriesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_countriesCache_Get(t *testing.T) {
	c := newCountriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Countries)
	err := c.ICache.(CountriesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CountriesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(CountriesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_countriesCache_MultiGet(t *testing.T) {
	c := newCountriesCache()
	defer c.Close()

	var testData []*model.Countries
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Countries))
	}

	err := c.ICache.(CountriesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CountriesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Countries))
	}
}

func Test_countriesCache_MultiSet(t *testing.T) {
	c := newCountriesCache()
	defer c.Close()

	var testData []*model.Countries
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Countries))
	}

	err := c.ICache.(CountriesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_countriesCache_Del(t *testing.T) {
	c := newCountriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Countries)
	err := c.ICache.(CountriesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_countriesCache_SetCacheWithNotFound(t *testing.T) {
	c := newCountriesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Countries)
	err := c.ICache.(CountriesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCountriesCache(t *testing.T) {
	c := NewCountriesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewCountriesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewCountriesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
