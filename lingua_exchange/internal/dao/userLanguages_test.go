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

func newUserLanguagesDao() *gotest.Dao {
	testData := &model.UserLanguages{}
	testData.ID = 1
	// you can set the other fields of testData here, such as:
	//testData.CreatedAt = time.Now()
	//testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	//c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // to test mysql, disable caching
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewUserLanguagesCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = NewUserLanguagesDao(d.DB, c.ICache.(cache.UserLanguagesCache))

	return d
}

func Test_userLanguagesDao_Create(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(UserLanguagesDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userLanguagesDao_DeleteByID(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)
	expectedSQLForDeletion := "UPDATE .*"

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(UserLanguagesDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(UserLanguagesDao).DeleteByID(d.Ctx, 0)
	assert.Error(t, err)
}

func Test_userLanguagesDao_UpdateByID(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(UserLanguagesDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(UserLanguagesDao).UpdateByID(d.Ctx, &model.UserLanguages{})
	assert.Error(t, err)

}

func Test_userLanguagesDao_GetByID(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(UserLanguagesDao).GetByID(d.Ctx, testData.ID)
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
	_, err = d.IDao.(UserLanguagesDao).GetByID(d.Ctx, 2)
	assert.Error(t, err)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(3, 4).
		WillReturnRows(rows)
	_, err = d.IDao.(UserLanguagesDao).GetByID(d.Ctx, 4)
	assert.Error(t, err)
}

func Test_userLanguagesDao_GetByColumns(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(UserLanguagesDao).GetByColumns(d.Ctx, &query.Params{
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
	_, _, err = d.IDao.(UserLanguagesDao).GetByColumns(d.Ctx, &query.Params{
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
	dao := &userLanguagesDao{}
	_, _, err = dao.GetByColumns(context.Background(), &query.Params{Columns: []query.Column{{}}})
	t.Log(err)
}

func Test_userLanguagesDao_DeleteByIDs(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(UserLanguagesDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(UserLanguagesDao).DeleteByIDs(d.Ctx, []uint64{0})
	assert.Error(t, err)
}

func Test_userLanguagesDao_GetByCondition(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(UserLanguagesDao).GetByCondition(d.Ctx, &query.Conditions{
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
	_, err = d.IDao.(UserLanguagesDao).GetByCondition(d.Ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: 2,
			},
		},
	})
	assert.Error(t, err)
}

func Test_userLanguagesDao_GetByIDs(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(UserLanguagesDao).GetByIDs(d.Ctx, []uint64{testData.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.IDao.(UserLanguagesDao).GetByIDs(d.Ctx, []uint64{111})
	assert.Error(t, err)

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userLanguagesDao_GetByLastID(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, err := d.IDao.(UserLanguagesDao).GetByLastID(d.Ctx, 0, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, err = d.IDao.(UserLanguagesDao).GetByLastID(d.Ctx, 0, 10, "unknown-column")
	assert.Error(t, err)
}

func Test_userLanguagesDao_CreateByTx(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	_, err := d.IDao.(UserLanguagesDao).CreateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userLanguagesDao_DeleteByTx(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)
	expectedSQLForDeletion := "UPDATE .*"

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(UserLanguagesDao).DeleteByTx(d.Ctx, d.DB, testData.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userLanguagesDao_UpdateByTx(t *testing.T) {
	d := newUserLanguagesDao()
	defer d.Close()
	testData := d.TestData.(*model.UserLanguages)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(UserLanguagesDao).UpdateByTx(d.Ctx, d.DB, testData)
	if err != nil {
		t.Fatal(err)
	}
}
