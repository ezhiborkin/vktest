package postgresql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/storage"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"math"
	"strings"
)

type Storage struct {
	db *sql.DB
}

func New(dataSourceName string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgresql.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddMovieStorage(movie *models.Movie) error {
	const op = "storage.postgresql.AddMovie"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	movieInsert := sq.Insert("movies").
		Columns("title", "description", "release_date", "rating", "actors_id").
		Values(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, movie.ActorsID).
		Suffix("RETURNING id")

	var movieID int64
	err = movieInsert.RunWith(tx).PlaceholderFormat(sq.Dollar).Scan(&movieID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrMovieExists)
	}

	actorUpdate := sq.Update("actors").
		Set("movies_id", sq.Expr("array_append(movies_id, ?)", movieID)).
		Where(sq.Expr("id = ANY(?)", movie.ActorsID))

	_, err = actorUpdate.RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetMoviesSortedStorage(sortBy string, sortDirection string) ([]*models.MovieListing, error) {
	const op = "storage.postgresql.GetMoviesSorted"

	sortColumn := "rating"

	switch sortBy {
	case "title":
		sortColumn = "title"
	case "release_date":
		sortColumn = "release_date"
	}

	if sortDirection != "ASC" && sortDirection != "DESC" {
		sortDirection = "DESC"
	}

	query, args, err := sq.
		Select("m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, m.release_date AS release_date, m.rating AS movie_rating, json_agg(a.name) AS actors").
		From("movies m").
		Join("actors a ON a.id = ANY(m.actors_id)").
		Where("m.deleted_at IS NULL AND a.deleted_at IS NULL").
		GroupBy("m.id, m.title, m.description, m.release_date, m.rating").
		OrderBy(sortColumn + " " + sortDirection).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var movies []*models.MovieListing
	for rows.Next() {
		movie := &models.MovieListing{}
		var actors sql.NullString
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &actors)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		var newActors *[]string
		if actors.Valid {
			if err := json.Unmarshal([]byte(actors.String), &newActors); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		}

		movie.Actors = *newActors

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *Storage) EditMovieStorage(movie *models.Movie) error {
	const op = "storage.postgresql.EditMovieStorage"

	updateBuilder := sq.Update("movies").Where(sq.Eq{"ID": &movie.ID})

	if movie.Title != "" {
		updateBuilder = updateBuilder.Set("title", movie.Title)
	}
	if movie.Description != "" {
		updateBuilder = updateBuilder.Set("description", movie.Description)
	}
	if !movie.ReleaseDate.IsZero() {
		updateBuilder = updateBuilder.Set("release_date", movie.ReleaseDate)
	}
	if movie.Rating != nil && *movie.Rating >= 0 && *movie.Rating <= 10 {
		roundedRating := math.Round(*movie.Rating*10) / 10
		updateBuilder = updateBuilder.Set("rating", roundedRating)
	}

	updateBuilder = updateBuilder.PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(sqlStr, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteMovieStorage(id int64) error {
	const op = "storage.postgresql.DeleteMovieStorage"

	query, args, err := sq.Update("movies").
		Set("deleted_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddActorStorage(actor *models.Actor) error {
	const op = "storage.postgresql.AddActor"

	query, args, err := sq.Insert("actors").
		Columns("name", "sex", "birthday", "movies_id").
		Values(&actor.Name, &actor.Sex, &actor.Birthday, &actor.MoviesID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) EditActorStorage(actor *models.Actor) error {
	const op = "storage.postgresql.EditActorStorage"

	updateBuilder := sq.Update("actors").Where(sq.Eq{"ID": &actor.ID})

	if actor.Name != "" {
		updateBuilder = updateBuilder.Set("name", actor.Name)
	}
	if actor.Sex != "" {
		updateBuilder = updateBuilder.Set("sex", actor.Sex)
	}
	if !actor.Birthday.IsZero() {
		updateBuilder = updateBuilder.Set("birthday", actor.Birthday)
	}

	updateBuilder = updateBuilder.PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(sqlStr, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetActorsStorage() ([]*models.ActorListing, error) {
	const op = "storage.postgresql.GetActorsStorage"

	query, args, err := sq.
		Select("a.id AS actor_id, a.name AS actor_name, a.sex AS actor_sex, a.birthday AS actor_birthday, json_agg(m.title) AS movies").
		From("actors a").
		LeftJoin("movies m ON a.id = ANY(m.actors_id) AND m.deleted_at IS NULL").
		Where("a.deleted_at IS NULL").
		GroupBy("a.id, a.name, a.sex, a.birthday").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var actors []*models.ActorListing
	for rows.Next() {
		actor := &models.ActorListing{}
		var movies sql.NullString
		err := rows.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday, &movies)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		var newMovies *[]string
		if movies.Valid {
			if err := json.Unmarshal([]byte(movies.String), &newMovies); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

		}
		actor.Movies = *newMovies

		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}

func (s *Storage) GetActorStorage(actorName string) (*models.Actor, error) {
	const op = "storage.postgresql.GetActor"

	var actor models.Actor

	query, args, err := sq.Select("id", "name", "sex", "birthday", "movies_id", "deleted_at").
		From("actors").
		Where(sq.Eq{"name": actorName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//stmt, err := s.db.Prepare("SELECT id, name, sex, birthday, movies_id, deleted_at FROM actors WHERE name = $1 LIMIT 1")

	row := s.db.QueryRow(query, args...)
	err = row.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday, &actor.MoviesID, &actor.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &actor, nil
}

func (s *Storage) DeleteActorStorage(id int64) error {
	const op = "storage.postgresql.DeleteActorStorage"

	query, args, err := sq.Update("actors").
		Set("deleted_at", sq.Expr("CURRENT_TIMESTAMP")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddActorsToMovieStorage(movieID int64, actors []int64) error {
	const op = "storage.postgresql.AddActorsToMovieStorage"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	moviessUpdate := sq.Update("movies").
		Set("actors_id", sq.Expr("array_cat(actors_id, ?)", actors)).
		Where(sq.Eq{"id": movieID})

	_, err = moviessUpdate.RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, actorID := range actors {
		movieUpdate := sq.Update("actors").
			Set("movies_id", sq.Expr("array_append(movies_id, ?)", movieID)).
			Where(sq.Eq{"id": actorID})

		_, err = movieUpdate.RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) AddMoviesToActorStorage(actorID int64, movies []int64) error {
	const op = "storage.postgresql.AddMoviesToActorStorage"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	actorsUpdate := sq.Update("actors").
		Set("movies_id", sq.Expr("array_cat(movies_id, ?)", movies)).
		Where(sq.Eq{"id": actorID})

	_, err = actorsUpdate.RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, movieID := range movies {
		actorUpdate := sq.Update("movies").
			Set("actors_id", sq.Expr("array_append(actors_id, ?)", actorID)).
			Where(sq.Eq{"id": movieID})

		_, err = actorUpdate.RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) GetMovieStorage(input string) ([]*models.MovieListing, error) {
	const op = "storage.postgresql.GetMovieStorage"

	inputLower := strings.ToLower(input)

	query, args, err := sq.Select("m.id AS movie_id").
		From("movies m").
		Join("actors a ON a.id = ANY(m.actors_id)").
		Where(sq.Or{sq.Expr("LOWER(m.title) LIKE ?", "%"+inputLower+"%"), sq.Expr("LOWER(a.name) LIKE ?", "%"+inputLower+"%")}).
		Where("m.deleted_at IS NULL AND a.deleted_at IS NULL").
		GroupBy("m.id").
		OrderBy("m.id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var moviesID []int64
	for rows.Next() {
		var movie int64
		err := rows.Scan(&movie)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		moviesID = append(moviesID, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var movies []*models.MovieListing
	for _, tempMovie := range moviesID {
		movie, err := s.GetMovieStorageByID(tempMovie)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		movies = append(movies, movie)
	}

	return movies, nil

}

func (s *Storage) GetMovieStorageByID(id int64) (*models.MovieListing, error) {
	const op = "storage.postgresql.GetMovieStorage"

	query, args, err := sq.
		Select("m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, m.release_date AS release_date, m.rating AS movie_rating, json_agg(a.name) AS actors").
		From("movies m").
		Join("actors a ON a.id = ANY(m.actors_id)").
		Where(sq.And{sq.Eq{"m.id": id}, sq.Expr("m.deleted_at IS NULL"), sq.Expr("a.deleted_at IS NULL")}).
		GroupBy("m.id, m.title, m.description, m.release_date, m.rating").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := s.db.QueryRow(query, args...)

	movie := &models.MovieListing{}
	var actors sql.NullString
	err = row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &actors)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrMovieNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var newActors []string
	if actors.Valid {
		if err := json.Unmarshal([]byte(actors.String), &newActors); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	movie.Actors = newActors

	return movie, nil
}

func (s *Storage) CreateUserStorage(email, role string, passHash []byte) error {
	const op = "storage.postgresql.CreateUserStorage"

	query, args, err := sq.Insert("users").
		Columns("email", "role", "password_hash").
		Values(email, role, passHash).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserStorage(email string) (*models.User, error) {
	const op = "storage.postgresql.GetUserStorage"

	query, args, err := sq.Select("id", "email", "role", "password_hash").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := s.db.QueryRow(query, args...)

	user := &models.User{}
	err = row.Scan(&user.ID, &user.Email, &user.Role, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
