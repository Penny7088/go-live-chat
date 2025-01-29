package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newTalkRecordsReadCache() *gotest.Cache {
	record1 := &model.TalkRecordsRead{}
	record1.ID = 1
	record2 := &model.TalkRecordsRead{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewTalkRecordsReadCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_talkRecordsReadCache_Set(t *testing.T) {
	c := newTalkRecordsReadCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsRead)
	err := c.ICache.(TalkRecordsReadCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(TalkRecordsReadCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_talkRecordsReadCache_Get(t *testing.T) {
	c := newTalkRecordsReadCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsRead)
	err := c.ICache.(TalkRecordsReadCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TalkRecordsReadCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(TalkRecordsReadCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_talkRecordsReadCache_MultiGet(t *testing.T) {
	c := newTalkRecordsReadCache()
	defer c.Close()

	var testData []*model.TalkRecordsRead
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.TalkRecordsRead))
	}

	err := c.ICache.(TalkRecordsReadCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TalkRecordsReadCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.TalkRecordsRead))
	}
}

func Test_talkRecordsReadCache_MultiSet(t *testing.T) {
	c := newTalkRecordsReadCache()
	defer c.Close()

	var testData []*model.TalkRecordsRead
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.TalkRecordsRead))
	}

	err := c.ICache.(TalkRecordsReadCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsReadCache_Del(t *testing.T) {
	c := newTalkRecordsReadCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsRead)
	err := c.ICache.(TalkRecordsReadCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsReadCache_SetCacheWithNotFound(t *testing.T) {
	c := newTalkRecordsReadCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecordsRead)
	err := c.ICache.(TalkRecordsReadCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTalkRecordsReadCache(t *testing.T) {
	c := NewTalkRecordsReadCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewTalkRecordsReadCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewTalkRecordsReadCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
