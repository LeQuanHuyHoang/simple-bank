package gapi

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/token"
	"simple-bank/utils"
	"testing"
	"time"
)

func TestUpdateUser(t *testing.T) {
	user, _ := randomUser(t)

	newName := utils.RandomOwner()
	newEmail := utils.RandomEmail()
	invalidEmail := "invalid_email"

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateUserParams{
					Username: user.Username,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: newEmail,
						Valid:  true,
					},
				}

				updatedUser := db.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(updatedUser, nil)

			},

			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},

			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, updatedUser.Username, user.Username)
				require.Equal(t, updatedUser.FullName, newName)
				require.Equal(t, updatedUser.Email, newEmail)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},

			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},

			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.NotFound)
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},

			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, -time.Minute)
			},

			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.Unauthenticated)
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},

			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},

			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.Unauthenticated)
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &invalidEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},

			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, user.Role, time.Minute)
			},

			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.InvalidArgument)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)

			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)

			tc.checkResponse(t, res, err)
		})
	}
}
