package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newCountryLanguagesCache() *gotest.Cache {
	record1 := &model.CountryLanguages{}
	record1.ID = 1
	record2 := &model.CountryLanguages{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewCountryLanguagesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_countryLanguagesCache_Set(t *testing.T) {
	c := newCountryLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CountryLanguages)
	err := c.ICache.(CountryLanguagesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(CountryLanguagesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_countryLanguagesCache_Get(t *testing.T) {
	c := newCountryLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CountryLanguages)
	err := c.ICache.(CountryLanguagesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CountryLanguagesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(CountryLanguagesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_countryLanguagesCache_MultiGet(t *testing.T) {
	c := newCountryLanguagesCache()
	defer c.Close()

	var testData []*model.CountryLanguages
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.CountryLanguages))
	}

	err := c.ICache.(CountryLanguagesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CountryLanguagesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.CountryLanguages))
	}
}

func Test_countryLanguagesCache_MultiSet(t *testing.T) {
	c := newCountryLanguagesCache()
	defer c.Close()

	var testData []*model.CountryLanguages
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.CountryLanguages))
	}

	err := c.ICache.(CountryLanguagesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_countryLanguagesCache_Del(t *testing.T) {
	c := newCountryLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CountryLanguages)
	err := c.ICache.(CountryLanguagesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_countryLanguagesCache_SetCacheWithNotFound(t *testing.T) {
	c := newCountryLanguagesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CountryLanguages)
	err := c.ICache.(CountryLanguagesCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCountryLanguagesCache(t *testing.T) {
	c := NewCountryLanguagesCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewCountryLanguagesCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewCountryLanguagesCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
