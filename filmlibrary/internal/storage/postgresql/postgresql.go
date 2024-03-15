package postgresql

import (
	"database/sql"
	"encoding/json"
	"filmlibrary/internal/domain/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"math"
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
		return fmt.Errorf("%s: %w", op, err)
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
	if movie.ActorsID != nil && len(movie.ActorsID) > 0 {
		updateBuilder = updateBuilder.Set("actors_id", movie.ActorsID)
	}

	updateBuilder = updateBuilder.PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(sqlStr, args)

	// Execute the SQL query
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

	stmt, err := s.db.Prepare("INSERT INTO actors (name, sex, birthday, movies_id) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(&actor.Name, &actor.Sex, &actor.Birthday, &actor.MoviesID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) EditActorStorage(actor *models.Actor) error {
	const op = "storage.postgresql.EditActorStorage"

	updateBuilder := sq.Update("actors").Where(sq.Eq{"ID": &actor.ID})

	if actor.Name != "" {
		updateBuilder = updateBuilder.Set("name", "kek")
	}
	if actor.Sex != "" {
		updateBuilder = updateBuilder.Set("sex", actor.Sex)
	}
	if !actor.Birthday.IsZero() {
		updateBuilder = updateBuilder.Set("birthday", actor.Birthday)
	}
	if actor.MoviesID != nil && len(actor.MoviesID) > 0 {
		updateBuilder = updateBuilder.Set("movies_id", actor.MoviesID)
	}

	updateBuilder = updateBuilder.PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(sqlStr, args)

	// Execute the SQL query
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
			//actor.Movies = append(actor.Movies, movies.String)

		}
		actor.Movies = *newMovies

		fmt.Println(actor.Movies)
		actors = append(actors, actor)
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}

func (s *Storage) GetActorStorage(actorName string) (*models.Actor, error) {
	const op = "storage.postgresql.GetActor"

	var actor models.Actor

	stmt, err := s.db.Prepare("SELECT id, name, sex, birthday, movies_id, deleted_at FROM actors WHERE name = $1 LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRow(actorName)
	err = row.Scan(&actor.ID, &actor.Name, &actor.Sex, &actor.Birthday, &actor.MoviesID, &actor.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &actor, nil
}

func (s *Storage) DeleteActorStorage(id int64) error {
	const op = "storage.postgresql.DeleteActorStorage"

	stmt, err := s.db.Prepare("UPDATE actors SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
