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
	"strings"
	"testing"
	"time"
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
		{
			name: "Test add actor with empty request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/actor/add/movies", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("AddMoviesToActor", mock.AnythingOfType("int64"), mock.AnythingOfType("[]int64")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.addMoviesToActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			assert.Equal(t, http.StatusBadRequest, resp.Code)

			// Assert response body
			expectedBody := "failed to decode request\n"
			assert.Equal(t, expectedBody, resp.Body.String())
		})
	}
}

func TestHandler_addMoviesToActor_InvalidRequestBody(t *testing.T) {
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
			name: "Test add movies to actor with invalid request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/actor/add/movies", strings.NewReader("invalid JSON")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("AddMoviesToActor", mock.AnythingOfType("int64"), mock.AnythingOfType("[]int64")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.addMoviesToActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			assert.Equal(t, http.StatusBadRequest, resp.Code)

			// Assert response body
			expectedBody := "failed to decode request\n"
			assert.Equal(t, expectedBody, resp.Body.String())
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
		{
			name: "Test delete actor with valid actor ID",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/delete/actor?id=123", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("DeleteActor", mock.AnythingOfType("int64")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.deleteActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					resp.Code, http.StatusOK)
			}

			// Assert response body
			expectedBody := "Successfully deleted an actor"
			if resp.Body.String() != expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					resp.Body.String(), expectedBody)
			}
		})
	}
}

func TestHandler_deleteActor_Invalid(t *testing.T) {
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
			name: "Test delete actor with invalid actor ID",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/delete/actor?id=invalid", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("DeleteActor", mock.AnythingOfType("int64")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.deleteActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v",
					resp.Code, http.StatusBadRequest)
			}

			// Assert response body
			expectedBody := "invalid actor ID\n"
			if resp.Body.String() != expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					resp.Body.String(), expectedBody)
			}
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
		{
			name: "Test edit actor with valid request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/edit/actor", bytes.NewBuffer([]byte(`{
					"id": 123,
					"name": "Updated Name",
					"sex": "male",
					"birthday": "2023-03-16T12:34:56Z"
				}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("EditActor", mock.AnythingOfType("*models.Actor")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.editActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					resp.Code, http.StatusOK)
			}

			// Assert response body
			expectedBody := "Successfully edited an actor"
			if resp.Body.String() != expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					resp.Body.String(), expectedBody)
			}
		})
	}
}

func TestHandler_editActor_NotFullBody(t *testing.T) {
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
			name: "Test edit actor with invalid request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/edit/actor", bytes.NewBuffer([]byte(`{
					"sex": "male",
					"birthday": "2023-03-16T12:34:56Z"
				}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorMock := &mocks.ActorProvider{}
			actorMock.On("EditActor", mock.AnythingOfType("*models.Actor")).Return(nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorMock,
			}

			// Execute the handler
			h.editActor(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					resp.Code, http.StatusOK)
			}

			// Assert response body
			expectedBody := "Successfully edited an actor"
			if resp.Body.String() != expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					resp.Body.String(), expectedBody)
			}
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
		{
			name: "Test get actors",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				actorProvider: mocks.NewActorProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/get/actors", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actorProviderMock := &mocks.ActorProvider{}
			timing := time.Now()
			actorProviderMock.On("GetActors").Return([]*models.ActorListing{
				{ID: 1, Name: "Actor 1", Sex: "Male", Birthday: timing, Movies: nil},
				{ID: 2, Name: "Actor 2", Sex: "Female", Birthday: timing, Movies: nil},
			}, nil) // Mock expectation

			h := &Handler{
				log:           tt.fields.log,
				actorProvider: actorProviderMock,
			}

			// Execute the handler
			h.getActors(tt.args.w, tt.args.r)

			// Assert HTTP response
			resp := tt.args.w.(*httptest.ResponseRecorder)
			if resp.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					resp.Code, http.StatusOK)
			}

			// Assert response body
			expectedBody := `[{"id":1,"name":"Actor 1","sex":"Male","birthday":"` + timing.Format("2006-01-02T15:04:05.999999999-07:00") + `"},{"id":2,"name":"Actor 2","sex":"Female","birthday":"` + timing.Format("2006-01-02T15:04:05.999999999-07:00") + `"}]`
			assert.Equal(t, expectedBody, resp.Body.String())

		})
	}
}
