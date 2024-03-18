package handler

import (
	"bytes"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/handler/mocks"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHandler_addActorsToMovie(t *testing.T) {
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
			h.addActorsToMovie(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_addMovie(t *testing.T) {
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
			name: "Test add movie success",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/add/movie", bytes.NewBuffer([]byte(`{"title":"Movie Title","description":"Movie Description"}`))),
			},
		},
		//{
		//	name: "Test add movie failed to decode request body",
		//	fields: fields{
		//		log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		//		movieProvider: mocks.NewMovieProvider(t),
		//	},
		//	args: args{
		//		w: httptest.NewRecorder(),
		//		r: httptest.NewRequest(http.MethodPost, "/add/movie", nil),
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movieMock := &mocks.MovieProvider{}
			movieMock.On("AddMovie", mock.AnythingOfType("*models.Movie")).Return(nil)

			h := &Handler{
				log:           tt.fields.log,
				movieProvider: movieMock,
			}

			h.addMovie(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusCreated, resp.Code)
			assert.Equal(t, "Successfully added a movie", resp.Body.String())
		})
	}
}

func TestHandler_addMovie_EmptyBody(t *testing.T) {
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
			name: "Test add movie failed to decode request body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/add/movie", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movieMock := &mocks.MovieProvider{}
			movieMock.On("AddMovie", mock.AnythingOfType("*models.Movie")).Return(nil)

			h := &Handler{
				log:           tt.fields.log,
				movieProvider: movieMock,
			}

			h.addMovie(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Equal(t, "failed to decode request\n", resp.Body.String())
		})
	}
}

func TestHandler_deleteMovie(t *testing.T) {
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
			name: "Test delete movie success",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/delete/movie?id=1", nil),
			},
		},
		{
			name: "Test delete movie invalid ID",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/delete/movie?id=invalid", nil),
			},
		},
		{
			name: "Test delete movie failed",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/delete/movie?id=1", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movieMock := &mocks.MovieProvider{}
			var expectedErr error
			switch tt.name {
			case "Test delete movie success":
				movieMock.On("DeleteMovie", int64(1)).Return(nil)
			case "Test delete movie invalid ID":
				expectedErr = fmt.Errorf("invalid movie ID\n")
			case "Test delete movie failed":
				movieMock.On("DeleteMovie", int64(1)).Return(fmt.Errorf("failed to delete movie"))
				expectedErr = fmt.Errorf("failed to delete a movie\n")
			}

			h := &Handler{
				log:           tt.fields.log,
				movieProvider: movieMock,
			}

			h.deleteMovie(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			if expectedErr != nil {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				assert.Equal(t, expectedErr.Error(), resp.Body.String())
			} else {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, "Successfully deleted a movie", resp.Body.String())
			}
		})
	}
}

func TestHandler_editMovie(t *testing.T) {
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
		name        string
		fields      fields
		args        args
		wantStatus  int
		wantMessage string
	}{
		{
			name: "Edit Movie Success",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/edit/movie", bytes.NewBufferString(`{"id":1,"title":"Updated Movie","release_date":"2024-03-20T12:00:00Z"}`)),
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Successfully edited a movie",
		},
		{
			name: "Empty Request Body",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/edit/movie", nil),
			},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "failed to decode request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMovieProvider := &mocks.MovieProvider{}
			mockMovieProvider.On("EditMovie", mock.Anything).Return(nil) // Set up mock expectation

			h := &Handler{
				log:           tt.fields.log,
				movieProvider: mockMovieProvider,
			}

			h.editMovie(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.wantStatus, resp.Code)
			assert.Equal(t, tt.wantMessage, strings.TrimSpace(resp.Body.String()))
		})
	}
}

func TestHandler_getMovie(t *testing.T) {
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
			name: "Test get movie success",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/get/movie?input=example", nil),
			},
		},
		{
			name: "Test get movie failed",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/get/movie?input=invalid", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movieMock := &mocks.MovieProvider{}
			var expectedErr error
			var expectedResponse string
			switch tt.name {
			case "Test get movie success":
				expectedResponse = `[{"id":1,"title":"Example Movie","release_date":"0001-01-01T00:00:00Z"}]`
				movieMock.On("GetMovie", "example").Return([]*models.MovieListing{{ID: 1, Title: "Example Movie", ReleaseDate: time.Time{}}}, nil)
			case "Test get movie failed":
				expectedErr = fmt.Errorf("failed to find a movie")
				expectedResponse = `[]`
				movieMock.On("GetMovie", "invalid").Return([]*models.MovieListing{}, expectedErr)
			}

			h := &Handler{
				log:           tt.fields.log,
				movieProvider: movieMock,
			}

			h.getMovie(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			if expectedErr != nil {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				assert.Equal(t, expectedErr.Error()+"\n", resp.Body.String())
			} else {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.Equal(t, expectedResponse, resp.Body.String())
			}
		})
	}
}

func TestHandler_getMoviesSorted(t *testing.T) {
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
			name: "Test get movies sorted",
			fields: fields{
				log:           slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
				movieProvider: mocks.NewMovieProvider(t),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/movies/sorted?sortBy=title&sortDir=asc", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movieMock := &mocks.MovieProvider{}
			movies := []*models.MovieListing{
				{ID: 1, Title: "Movie 1", Description: "Description 1", ReleaseDate: time.Now(), Rating: ptrFloat64(8.5), Actors: []string{"Oleg", "Putin"}},
				{ID: 2, Title: "Movie 2", Description: "Description 2", ReleaseDate: time.Now(), Rating: ptrFloat64(7.9), Actors: []string{"OPOPOPO", "GVNO"}},
			}
			movieMock.On("GetMoviesSorted", "title", "asc").Return(movies, nil)

			h := &Handler{
				log:           tt.fields.log,
				movieProvider: movieMock,
			}

			h.getMoviesSorted(tt.args.w, tt.args.r)

			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Body.String(), "Sorted by: title asc")
			assert.Contains(t, resp.Body.String(), "\"id\":1,\"title\":\"Movie 1\",\"description\":\"Description 1\",\"release_date\":")
			assert.Contains(t, resp.Body.String(), "\"id\":2,\"title\":\"Movie 2\",\"description\":\"Description 2\",\"release_date\":")
		})
	}
}

func ptrFloat64(f float64) *float64 {
	return &f
}
