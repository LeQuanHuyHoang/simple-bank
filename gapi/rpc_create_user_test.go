package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/utils"
	"simple-bank/worker"
	mockwk "simple-bank/worker/mock"
	"testing"
)

type eqCreateUserTxParamsMatcher struct {
	args     db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	//Casts the underlying VALUE of interface x to a specific type CreateUserParams
	//db.CreateUserParams store the data send to server
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}
	err := utils.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	expected.args.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.args.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err = actualArg.AfterCreate(expected.user)
	return err == nil
}

func (expected eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password  %v", expected.args, expected.password)
}

func EqCreateUserTxParams(args db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{args, password, user}
}

func randomUser(t *testing.T) (db.User, string) {
	password := utils.RandomString(6)
	hastPassword, err := utils.HasdPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       utils.RandomString(6),
		FullName:       utils.RandomString(6),
		Role:           utils.DepositorRole,
		HashedPassword: hastPassword,
		Email:          utils.RandomEmail(),
	}
	return user, password
}

func TestCreateUserTx(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Password: password,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				args := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(args, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createUser := res.GetUser()
				require.Equal(t, createUser.Username, user.Username)
				require.Equal(t, createUser.FullName, user.FullName)
				require.Equal(t, createUser.Email, user.Email)
			},
		},
		{
			name: "InternalServerError",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Password: password,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.Internal)
			},
		},
		{
			name: "DuplicateUserName",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Password: password,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, db.ErrUniqueViolation)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.AlreadyExists)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)
			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)

			res, err := server.CreateUser(context.Background(), tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}
