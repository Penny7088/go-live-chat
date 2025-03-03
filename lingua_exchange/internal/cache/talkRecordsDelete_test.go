package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newTalkRecordsDeleteCache() *gotest.Cache {
	record1 := &model.TalkRecordsDelete{}
	record1.ID = 1
	record2 := &model.TalkRecordsDelete{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewTalkRecordsDeleteCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_talkRecordsDeleteCache_Set(t *testing.T) {
	c := newTalkRecordsDeleteCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsDelete)
	err := c.ICache.(TalkRecordsDeleteCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(TalkRecordsDeleteCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_talkRecordsDeleteCache_Get(t *testing.T) {
	c := newTalkRecordsDeleteCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsDelete)
	err := c.ICache.(TalkRecordsDeleteCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TalkRecordsDeleteCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(TalkRecordsDeleteCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_talkRecordsDeleteCache_MultiGet(t *testing.T) {
	c := newTalkRecordsDeleteCache()
	defer c.Close()

	var testData []*model.TalkRecordsDelete
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.TalkRecordsDelete))
	}

	err := c.ICache.(TalkRecordsDeleteCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TalkRecordsDeleteCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.TalkRecordsDelete))
	}
}

func Test_talkRecordsDeleteCache_MultiSet(t *testing.T) {
	c := newTalkRecordsDeleteCache()
	defer c.Close()

	var testData []*model.TalkRecordsDelete
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.TalkRecordsDelete))
	}

	err := c.ICache.(TalkRecordsDeleteCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsDeleteCache_Del(t *testing.T) {
	c := newTalkRecordsDeleteCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsDelete)
	err := c.ICache.(TalkRecordsDeleteCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsDeleteCache_SetCacheWithNotFound(t *testing.T) {
	c := newTalkRecordsDeleteCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsDelete)
	err := c.ICache.(TalkRecordsDeleteCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTalkRecordsDeleteCache(t *testing.T) {
	c := NewTalkRecordsDeleteCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewTalkRecordsDeleteCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewTalkRecordsDeleteCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
