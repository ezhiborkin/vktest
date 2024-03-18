package handler

import (
	"bytes"
	"filmlibrary/internal/handler/mocks"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler_loginUser(t *testing.T) {
	type fields struct {
		log           *slog.Logger
		userProvider  UserProvider
		actorProvider ActorProvider
		movieProvider MovieProvider
		authProvider  AuthProvider
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test login user",
			fields: fields{
				log:          slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				authProvider: mocks.NewAuthProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(`{"email":"33kjdkjj123kk@al.ru","password":"opopop111"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMock := &mocks.AuthProvider{}
			authMock.On("LoginUser", "33kjdkjj123kk@al.ru", "opopop111").Return("token", nil)

			h := &Handler{
				log:          tt.fields.log,
				authProvider: authMock,
			}

			h.loginUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Header().Get("Authorization"), "Bearer token")
			assert.Contains(t, resp.Body.String(), "Successfully logged in, your token: token")
		})
	}
}

func TestHandler_loginUser_EmptyBody(t *testing.T) {
	type fields struct {
		log           *slog.Logger
		userProvider  UserProvider
		actorProvider ActorProvider
		movieProvider MovieProvider
		authProvider  AuthProvider
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test login user empty body",
			fields: fields{
				log:          slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				authProvider: mocks.NewAuthProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMock := &mocks.AuthProvider{}
			authMock.On("LoginUser", "33kjdkjj123kk@al.ru", "opopop111").Return("token", nil)

			h := &Handler{
				log:          tt.fields.log,
				authProvider: authMock,
			}

			h.loginUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			//assert.Contains(t, resp.Header().Get("Authorization"), "Bearer token")
			assert.Contains(t, resp.Body.String(), "failed to decode request\n")
		})
	}
}

func TestHandler_loginUser_EmailMissing(t *testing.T) {
	type fields struct {
		log           *slog.Logger
		userProvider  UserProvider
		actorProvider ActorProvider
		movieProvider MovieProvider
		authProvider  AuthProvider
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test login user empty body",
			fields: fields{
				log:          slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				authProvider: mocks.NewAuthProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authMock := &mocks.AuthProvider{}
			authMock.On("LoginUser", "33kjdkjj123kk@al.ru", "opopop111").Return("token", []byte(`{"body":"33kjdkjj123kk@al.ru",password":"opopop111"}`))

			h := &Handler{
				log:          tt.fields.log,
				authProvider: authMock,
			}

			h.loginUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), "failed to decode request\n")
		})
	}
}
