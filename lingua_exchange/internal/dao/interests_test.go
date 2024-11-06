package dao

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/model"
)

func newInterestsDao() *gotest.Dao {
	testData := &model.Interests{}
	testData.ID = 1
	// you can set the other fields of testData here, such as:
	//testData.CreatedAt = time.Now()
	//testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	//c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // to test mysql, disable caching
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewInterestsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = NewInterestsDao(d.DB, c.ICache.(cache.InterestsCache))

	return d
}

func Test_interestsDao_Create(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(InterestsDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsDao_DeleteByID(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)
	expectedSQLForDeletion := "UPDATE .*"

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(InterestsDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(InterestsDao).DeleteByID(d.Ctx, 0)
	assert.Error(t, err)
}

func Test_interestsDao_UpdateByID(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(InterestsDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(InterestsDao).UpdateByID(d.Ctx, &model.Interests{})
	assert.Error(t, err)

}

func Test_interestsDao_GetByID(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(InterestsDao).GetByID(d.Ctx, testData.ID)
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
	_, err = d.IDao.(InterestsDao).GetByID(d.Ctx, 2)
	assert.Error(t, err)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(3, 4).
		WillReturnRows(rows)
	_, err = d.IDao.(InterestsDao).GetByID(d.Ctx, 4)
	assert.Error(t, err)
}

func Test_interestsDao_GetByColumns(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(InterestsDao).GetByColumns(d.Ctx, &query.Params{
		Page:  0,
		Limit: 10,
		Sort:  "ignore count", // ignore test count(*)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, _, err = d.IDao.(InterestsDao).GetByColumns(d.Ctx, &query.Params{
		Page:  0,
		Limit: 10,
		Columns: []query.Column{
			{
				Name:  "id",
				Exp:   "<",
				Value: 0,
			},
		},
	})
	assert.Error(t, err)

	// error test
	dao := &interestsDao{}
	_, _, err = dao.GetByColumns(context.Background(), &query.Params{Columns: []query.Column{{}}})
	t.Log(err)
}

func Test_interestsDao_DeleteByIDs(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(InterestsDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(InterestsDao).DeleteByIDs(d.Ctx, []uint64{0})
	assert.Error(t, err)
}

func Test_interestsDao_GetByCondition(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(InterestsDao).GetByCondition(d.Ctx, &query.Conditions{
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
	_, err = d.IDao.(InterestsDao).GetByCondition(d.Ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: 2,
			},
		},
	})
	assert.Error(t, err)
}

func Test_interestsDao_GetByIDs(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(InterestsDao).GetByIDs(d.Ctx, []uint64{testData.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.IDao.(InterestsDao).GetByIDs(d.Ctx, []uint64{111})
	assert.Error(t, err)

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsDao_GetByLastID(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, err := d.IDao.(InterestsDao).GetByLastID(d.Ctx, 0, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, err = d.IDao.(InterestsDao).GetByLastID(d.Ctx, 0, 10, "unknown-column")
	assert.Error(t, err)
}

func Test_interestsDao_CreateByTx(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	_, err := d.IDao.(InterestsDao).CreateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsDao_DeleteByTx(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)
	expectedSQLForDeletion := "UPDATE .*"

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(InterestsDao).DeleteByTx(d.Ctx, d.DB, testData.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_interestsDao_UpdateByTx(t *testing.T) {
	d := newInterestsDao()
	defer d.Close()
	testData := d.TestData.(*model.Interests)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(InterestsDao).UpdateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}
