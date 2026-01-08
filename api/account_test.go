package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/DakshChawla/simplebank/db/mock"
	db "github.com/DakshChawla/simplebank/db/sqlc"
	"github.com/DakshChawla/simplebank/token"
	"github.com/DakshChawla/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t,
					req,
					tokenMaker,
					authorizationTypeBearer,
					account.Owner,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t,
					req,
					tokenMaker,
					authorizationTypeBearer,
					account.Owner,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t,
					req,
					tokenMaker,
					authorizationTypeBearer,
					account.Owner,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t,
					req,
					tokenMaker,
					authorizationTypeBearer,
					account.Owner,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).
					Return(account, nil)
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(
					t,
					req,
					tokenMaker,
					authorizationTypeBearer,
					"wrong_user",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func addAuthorization(
	t *testing.T,
	req *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var got db.Account
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, account, got)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
