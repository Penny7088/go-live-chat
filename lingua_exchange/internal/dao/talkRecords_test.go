package dao

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
)

func newTalkRecordsDao() *gotest.Dao {
	testData := &model.TalkRecords{}
	testData.ID = 1
	// you can set the other fields of testData here, such as:
	// testData.CreatedAt = time.Now()
	// testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	// c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // to test mysql, disable caching
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewTalkRecordsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = NewTalkRecordsDao(d.DB, c.ICache.(cache.TalkRecordsCache))

	return d
}

func Test_talkRecordsDao_Create(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(TalkRecordsDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsDao_DeleteByID(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)
	expectedSQLForDeletion := "UPDATE .*"

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(TalkRecordsDao).DeleteByID(d.Ctx, testData.MsgID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(TalkRecordsDao).DeleteByID(d.Ctx, "0")
	assert.Error(t, err)
}

func Test_talkRecordsDao_UpdateByID(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(TalkRecordsDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(TalkRecordsDao).UpdateByID(d.Ctx, &model.TalkRecords{})
	assert.Error(t, err)

}

func Test_talkRecordsDao_GetByID(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(TalkRecordsDao).GetByID(d.Ctx, testData.MsgID)
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// notfound error
	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(2).
		WillReturnRows(rows)
	_, err = d.IDao.(TalkRecordsDao).GetByID(d.Ctx, "2")
	assert.Error(t, err)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(3, 4).
		WillReturnRows(rows)
	_, err = d.IDao.(TalkRecordsDao).GetByID(d.Ctx, "4")
	assert.Error(t, err)
}

func Test_talkRecordsDao_DeleteByIDs(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(TalkRecordsDao).DeleteByID(d.Ctx, testData.MsgID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(TalkRecordsDao).DeleteByIDs(d.Ctx, []string{"0"})
	assert.Error(t, err)
}

func Test_talkRecordsDao_GetByCondition(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(TalkRecordsDao).GetByCondition(d.Ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: testData.ID,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// notfound error
	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(2).
		WillReturnRows(rows)
	_, err = d.IDao.(TalkRecordsDao).GetByCondition(d.Ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: 2,
			},
		},
	})
	assert.Error(t, err)
}

func Test_talkRecordsDao_GetByIDs(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(TalkRecordsDao).GetByIDs(d.Ctx, []string{testData.MsgID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.IDao.(TalkRecordsDao).GetByIDs(d.Ctx, []string{"111"})
	assert.Error(t, err)

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsDao_GetByLastID(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, err := d.IDao.(TalkRecordsDao).GetByLastID(d.Ctx, "0", 10, "")
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, err = d.IDao.(TalkRecordsDao).GetByLastID(d.Ctx, "0", 10, "unknown-column")
	assert.Error(t, err)
}

func Test_talkRecordsDao_DeleteByTx(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)
	expectedSQLForDeletion := "UPDATE .*"

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(TalkRecordsDao).DeleteByTx(d.Ctx, d.DB, testData.MsgID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_talkRecordsDao_UpdateByTx(t *testing.T) {
	d := newTalkRecordsDao()
	defer d.Close()
	testData := d.TestData.(*model.TalkRecords)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(TalkRecordsDao).UpdateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}
