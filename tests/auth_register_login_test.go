package test

import (
	"sso/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/emp2ty0/sso-protos/gen/go/sso"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func Test_HappyPath(t *testing.T) {
	ctx, suite := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := suite.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLog, err := suite.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})

	require.NoError(t, err)

	token := respLog.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

}

func TestRegister_DuplicatedRegistration(t *testing.T) {
	ctx, suite := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	respReg, err := suite.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = suite.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailedCase(t *testing.T) {
	ctx, suite := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		wantCode    codes.Code
		expectedErr string
	}{
		{
			name:        "Register with Empty password",
			email:       gofakeit.Email(),
			password:    "",
			wantCode:    codes.InvalidArgument,
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			wantCode:    codes.InvalidArgument,
			expectedErr: "email is required",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			wantCode:    codes.InvalidArgument,
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := suite.AuthClient.Register(ctx, &sso.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, tt.wantCode, st.Code())
		})
	}
}

func TestLogin_FailedCase(t *testing.T) {
	ctx, suite := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int
		wantCode    codes.Code
		expectedErr string
	}{
		{
			name:        "Login with Empty password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			wantCode:    codes.InvalidArgument,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			appID:       appID,
			wantCode:    codes.InvalidArgument,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty",
			email:       "",
			password:    "",
			appID:       appID,
			wantCode:    codes.InvalidArgument,
			expectedErr: "email is required",
		},
		{
			name:        "App ID Empty",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			appID:       0,
			wantCode:    codes.InvalidArgument,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := suite.AuthClient.Login(ctx, &sso.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    int32(tt.appID),
			})

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, tt.wantCode, st.Code())
		})
	}
}
