package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

func newTalkRecordsCache() *gotest.Cache {
	record1 := &model.TalkRecords{}
	record1.ID = 1
	record2 := &model.TalkRecords{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewTalkRecordsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_talkRecordsCache_Set(t *testing.T) {
	c := newTalkRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecords)
	err := c.ICache.(TalkRecordsCache).Set(c.Ctx, record.MsgID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(TalkRecordsCache).Set(c.Ctx, "0", nil, time.Hour)
	assert.NoError(t, err)
}

func Test_talkRecordsCache_Get(t *testing.T) {
	c := newTalkRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecords)
	err := c.ICache.(TalkRecordsCache).Set(c.Ctx, record.MsgID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TalkRecordsCache).Get(c.Ctx, record.MsgID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(TalkRecordsCache).Get(c.Ctx, "0")
	assert.Error(t, err)
}

func Test_talkRecordsCache_MultiSet(t *testing.T) {
	c := newTalkRecordsCache()
	defer c.Close()

	var testData []*model.TalkRecords
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.TalkRecords))
	}

	err := c.ICache.(TalkRecordsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsCache_Del(t *testing.T) {
	c := newTalkRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecords)
	err := c.ICache.(TalkRecordsCache).Del(c.Ctx, record.MsgID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsCache_SetCacheWithNotFound(t *testing.T) {
	c := newTalkRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.TalkRecords)
	err := c.ICache.(TalkRecordsCache).SetCacheWithNotFound(c.Ctx, record.MsgID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTalkRecordsCache(t *testing.T) {
	c := NewTalkRecordsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewTalkRecordsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewTalkRecordsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
