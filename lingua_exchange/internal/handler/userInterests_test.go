package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/httpcli"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
)

func newUserInterestsHandler() *gotest.Handler {
	testData := &model.UserInterests{}
	testData.ID = 1
	// you can set the other fields of testData here, such as:
	// testData.CreatedAt = time.Now()
	// testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewUserInterestsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = dao.NewUserInterestsDao(d.DB, c.ICache.(cache.UserInterestsCache))

	// init mock handler
	h := gotest.NewHandler(d, testData)
	h.IHandler = &userInterestsHandler{iDao: d.IDao.(dao.UserInterestsDao)}
	iHandler := h.IHandler.(UserInterestsHandler)

	testFns := []gotest.RouterInfo{
		{
			FuncName:    "Create",
			Method:      http.MethodPost,
			Path:        "/userInterests",
			HandlerFunc: iHandler.Create,
		},
		{
			FuncName:    "DeleteByID",
			Method:      http.MethodDelete,
			Path:        "/userInterests/:id",
			HandlerFunc: iHandler.DeleteByID,
		},
		// {
		// 	FuncName:    "UpdateByID",
		// 	Method:      http.MethodPut,
		// 	Path:        "/userInterests/:id",
		// 	HandlerFunc: iHandler.UpdateByID,
		// },
		{
			FuncName:    "GetByID",
			Method:      http.MethodGet,
			Path:        "/userInterests/:id",
			HandlerFunc: iHandler.GetByID,
		},
		{
			FuncName:    "List",
			Method:      http.MethodPost,
			Path:        "/userInterests/list",
			HandlerFunc: iHandler.List,
		},
		{
			FuncName:    "DeleteByIDs",
			Method:      http.MethodPost,
			Path:        "/userInterests/delete/ids",
			HandlerFunc: iHandler.DeleteByIDs,
		},
		{
			FuncName:    "GetByCondition",
			Method:      http.MethodPost,
			Path:        "/userInterests/condition",
			HandlerFunc: iHandler.GetByCondition,
		},
		{
			FuncName:    "ListByIDs",
			Method:      http.MethodPost,
			Path:        "/userInterests/list/ids",
			HandlerFunc: iHandler.ListByIDs,
		},
		{
			FuncName:    "ListByLastID",
			Method:      http.MethodGet,
			Path:        "/userInterests/list",
			HandlerFunc: iHandler.ListByLastID,
		},
	}

	h.GoRunHTTPServer(testFns)

	time.Sleep(time.Millisecond * 200)
	return h
}

func Test_userInterestsHandler_Create(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := &types.CreateUserInterestsRequest{}
	_ = copier.Copy(testData, h.TestData.(*model.UserInterests))

	h.MockDao.SQLMock.ExpectBegin()
	args := h.MockDao.GetAnyArgs(h.TestData)
	h.MockDao.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(args[:len(args)-1]...). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(1, 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Post(result, h.GetRequestURL("Create"), testData)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", result)

}

func Test_userInterestsHandler_DeleteByID(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)
	expectedSQLForDeletion := "UPDATE .*"
	expectedArgsForDeletionTime := h.MockDao.AnyTime

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec(expectedSQLForDeletion).
		WithArgs(expectedArgsForDeletionTime, testData.ID). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Delete(result, h.GetRequestURL("DeleteByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = httpcli.Delete(result, h.GetRequestURL("DeleteByID", 0))
	assert.NoError(t, err)

	// delete error test
	err = httpcli.Delete(result, h.GetRequestURL("DeleteByID", 111))
	assert.Error(t, err)
}

func Test_userInterestsHandler_UpdateByID(t *testing.T) {
	// h := newUserInterestsHandler()
	// defer h.Close()
	// testData := &types.UpdateUserInterestsByIDRequest{}
	// _ = copier.Copy(testData, h.TestData.(*model.UserInterests))
	//
	// h.MockDao.SQLMock.ExpectBegin()
	// h.MockDao.SQLMock.ExpectExec("UPDATE .*").
	// 	WithArgs(h.MockDao.AnyTime, testData.ID). // adjusted for the amount of test data
	// 	WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	// h.MockDao.SQLMock.ExpectCommit()
	//
	// result := &httpcli.StdResult{}
	// err := httpcli.Put(result, h.GetRequestURL("UpdateByID", testData.ID), testData)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if result.Code != 0 {
	// 	t.Fatalf("%+v", result)
	// }
	//
	// // zero id error test
	// err = httpcli.Put(result, h.GetRequestURL("UpdateByID", 0), testData)
	// assert.NoError(t, err)
	//
	// // update error test
	// err = httpcli.Put(result, h.GetRequestURL("UpdateByID", 111), testData)
	// assert.Error(t, err)
}

func Test_userInterestsHandler_GetByID(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Get(result, h.GetRequestURL("GetByID", testData.ID))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = httpcli.Get(result, h.GetRequestURL("GetByID", 0))
	assert.NoError(t, err)

	// get error test
	err = httpcli.Get(result, h.GetRequestURL("GetByID", 111))
	assert.Error(t, err)
}

func Test_userInterestsHandler_List(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Post(result, h.GetRequestURL("List"), &types.ListUserInterestssRequest{query.Params{
		Page:  0,
		Limit: 10,
		Sort:  "ignore count", // ignore test count
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// nil params error test
	err = httpcli.Post(result, h.GetRequestURL("List"), nil)
	assert.NoError(t, err)

	// get error test
	err = httpcli.Post(result, h.GetRequestURL("List"), &types.ListUserInterestssRequest{query.Params{
		Page:  0,
		Limit: 10,
		Sort:  "unknown-column",
	}})
	assert.Error(t, err)
}

func Test_userInterestsHandler_DeleteByIDs(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)

	h.MockDao.SQLMock.ExpectBegin()
	h.MockDao.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(h.MockDao.AnyTime, testData.ID). // adjusted for the amount of test data
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	h.MockDao.SQLMock.ExpectCommit()

	result := &httpcli.StdResult{}
	err := httpcli.Post(result, h.GetRequestURL("DeleteByIDs"), &types.DeleteUserInterestssByIDsRequest{IDs: []uint64{testData.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	err = httpcli.Post(result, h.GetRequestURL("DeleteByIDs"), nil)
	assert.NoError(t, err)

	// get error test
	err = httpcli.Post(result, h.GetRequestURL("DeleteByIDs"), &types.DeleteUserInterestssByIDsRequest{IDs: []uint64{111}})
	assert.Error(t, err)
}

func Test_userInterestsHandler_GetByCondition(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Post(result, h.GetRequestURL("GetByCondition"), &types.GetUserInterestsByConditionRequest{
		query.Conditions{
			Columns: []query.Column{
				{
					Name:  "id",
					Value: testData.ID,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero error test
	err = httpcli.Post(result, h.GetRequestURL("GetByCondition"), nil)
	assert.NoError(t, err)

	// get error test
	err = httpcli.Post(result, h.GetRequestURL("GetByCondition"), &types.GetUserInterestsByConditionRequest{
		query.Conditions{
			Columns: []query.Column{
				{
					Name:  "id",
					Value: 2,
				},
			},
		},
	})
	assert.Error(t, err)
}

func Test_userInterestsHandler_ListByIDs(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Post(result, h.GetRequestURL("ListByIDs"), &types.ListUserInterestssByIDsRequest{IDs: []uint64{testData.ID}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// zero id error test
	_ = httpcli.Post(result, h.GetRequestURL("ListByIDs"), nil)

	// get error test
	err = httpcli.Post(result, h.GetRequestURL("ListByIDs"), &types.ListUserInterestssByIDsRequest{IDs: []uint64{111}})
	assert.Error(t, err)
}

func Test_userInterestsHandler_ListByLastID(t *testing.T) {
	h := newUserInterestsHandler()
	defer h.Close()
	testData := h.TestData.(*model.UserInterests)

	// column names and corresponding data
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(testData.ID)

	h.MockDao.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	result := &httpcli.StdResult{}
	err := httpcli.Get(result, h.GetRequestURL("ListByLastID"), httpcli.WithParams(map[string]interface{}{"id": 10}))
	if err != nil {
		t.Fatal(err)
	}
	if result.Code != 0 {
		t.Fatalf("%+v", result)
	}

	// error test
	err = httpcli.Get(result, h.GetRequestURL("ListByLastID"), httpcli.WithParams(map[string]interface{}{"lastID": 0, "limit": 10, "sort": "unknown-column"}))
	assert.Error(t, err)
}

func TestNewUserInterestsHandler(t *testing.T) {
	defer func() {
		recover()
	}()
	_ = NewUserInterestsHandler()
}
