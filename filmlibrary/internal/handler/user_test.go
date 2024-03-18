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

func TestHandler_createUser(t *testing.T) {
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
				userProvider: mocks.NewUserProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create/user", bytes.NewBuffer([]byte(`{"email":"33kjdkjj123kk@al.ru","role":"admin","password":"opopop111"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := &mocks.UserProvider{}
			userMock.On("CreateUser", "33kjdkjj123kk@al.ru", "admin", "opopop111").Return(nil)

			h := &Handler{
				log:          tt.fields.log,
				userProvider: userMock,
			}

			h.createUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Body.String(), "Created user with email - 33kjdkjj123kk@al.ru")
		})
	}
}

func TestHandler_createUser_EmptyBody(t *testing.T) {
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
				userProvider: mocks.NewUserProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create/user", bytes.NewBuffer(nil)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := &mocks.UserProvider{}
			userMock.On("CreateUser").Return(nil)

			h := &Handler{
				log:          tt.fields.log,
				userProvider: userMock,
			}

			h.createUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), "failed to decode request\n")
		})
	}
}

func TestHandler_createUser_RoleMissing(t *testing.T) {
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
				userProvider: mocks.NewUserProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create/user", bytes.NewBuffer([]byte(`{"email":"33kjdkjj123kk@al.ru","password":"opopop111"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := &mocks.UserProvider{}
			userMock.On("CreateUser", "33kjdkjj123kk@al.ru", "admin", "opopop111").Return(nil)

			h := &Handler{
				log:          tt.fields.log,
				userProvider: userMock,
			}

			h.createUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), "role is empty\n")
		})
	}
}

func TestHandler_createUser_EmailMissing(t *testing.T) {
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
				userProvider: mocks.NewUserProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create/user", bytes.NewBuffer([]byte(`{"role":"admin","password":"opopop111"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := &mocks.UserProvider{}
			userMock.On("CreateUser", "33kjdkjj123kk@al.ru", "admin", "opopop111").Return(nil)

			h := &Handler{
				log:          tt.fields.log,
				userProvider: userMock,
			}

			h.createUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), "email is empty\n")
		})
	}
}

func TestHandler_createUser_PasswordMissing(t *testing.T) {
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
				userProvider: mocks.NewUserProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/create/user", bytes.NewBuffer([]byte(`{"email":"33kjdkjj123kk@al.ru","role":"admin"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := &mocks.UserProvider{}
			userMock.On("CreateUser", "33kjdkjj123kk@al.ru", "admin", "opopop111").Return(nil)

			h := &Handler{
				log:          tt.fields.log,
				userProvider: userMock,
			}

			h.createUser(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), "password is empty\n")
		})
	}
}
