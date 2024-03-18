package handler

import (
	"bytes"
	"encoding/json"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEditActor(t *testing.T) {
	mockActorProvider := mocks.NewActorProvider(t)
	mockUserProvider := mocks.NewUserProvider(t)
	mockMovieProvider := mocks.NewMovieProvider(t)
	mockAuthProvider := mocks.NewAuthProvider(t)

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := New(logger, mockUserProvider, mockActorProvider, mockMovieProvider, mockAuthProvider)

	actor := &models.Actor{ID: 1, Name: "John Doe"}
	actorJSON, _ := json.Marshal(actor)
	req, _ := http.NewRequest("POST", "/edit/actor", bytes.NewBuffer(actorJSON))
	rr := httptest.NewRecorder()

	mockActorProvider.On("EditActor", actor).Return(nil)

	h.editActor(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Successfully edited an actor", rr.Body.String())
}

func TestHandler_addActor(t *testing.T) {
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
			name: "Test add actor",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/add/actor", bytes.NewBuffer([]byte(`{
					"name": "vasiliy",
					"sex": "male",
					"birthday": "2023-03-16T12:34:56Z"
				}`))), // You can provide request body here if needed
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("AddActor", mock.AnythingOfType("*models.Actor")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.addActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			assert.Equal(t, http.StatusCreated, resp.Code)

			// Assert response body
			expectedBody := "Successfully added an actor"
			assert.Equal(t, expectedBody, resp.Body.String())

		})
	}
}

func TestHandler_addActor_NameMissing(t *testing.T) {
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
			name: "Test add actor with invalid request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/add/actor", bytes.NewBuffer([]byte(`{
					"sex": "male",
					"birthday": "2023-03-16T12:34:56Z"
				}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("AddActor", mock.AnythingOfType("*models.Actor")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.addActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			assert.Equal(t, http.StatusBadRequest, resp.Code)

			// Assert response body
			expectedBody := "actor name is empty\n"
			assert.Equal(t, expectedBody, resp.Body.String())

		})
	}
}

func TestHandler_addActor_BodyEmpty(t *testing.T) {
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
			name: "Test add actor with empty request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/add/actor", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("AddActor", mock.AnythingOfType("*models.Actor")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.addActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			assert.Equal(t, http.StatusBadRequest, resp.Code)

			// Assert response body
			expectedBody := "failed to decode request\n"
			assert.Equal(t, expectedBody, resp.Body.String())

		})
	}
}

func TestHandler_addMoviesToActor(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:           tt.fields.log,
				userProvider:  tt.fields.userProvider,
				actorProvider: tt.fields.actorProvider,
				movieProvider: tt.fields.movieProvider,
				authProvider:  tt.fields.authProvider,
			}
			h.addMoviesToActor(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_deleteActor(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:           tt.fields.log,
				userProvider:  tt.fields.userProvider,
				actorProvider: tt.fields.actorProvider,
				movieProvider: tt.fields.movieProvider,
				authProvider:  tt.fields.authProvider,
			}
			h.deleteActor(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_editActor(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:           tt.fields.log,
				userProvider:  tt.fields.userProvider,
				actorProvider: tt.fields.actorProvider,
				movieProvider: tt.fields.movieProvider,
				authProvider:  tt.fields.authProvider,
			}
			h.editActor(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_getActors(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:           tt.fields.log,
				userProvider:  tt.fields.userProvider,
				actorProvider: tt.fields.actorProvider,
				movieProvider: tt.fields.movieProvider,
				authProvider:  tt.fields.authProvider,
			}
			h.getActors(tt.args.w, tt.args.r)
		})
	}
}
