SELECT
    m.id AS movie_id,
    m.title AS movie_title,
    m.description AS movie_description,
    m.release_date AS release_date,
    m.rating AS movie_rating,
    json_agg(a.name) AS actors
FROM movies m
         JOIN actors a ON a.id = ANY(m.actors_id)
WHERE m.deleted_at IS NULL
  AND a.deleted_at IS NULL
GROUP BY m.id, m.title, m.description, m.release_date, m.rating;


