SELECT
    m.id,
    m.title,
    m.description,
    m.release_date,
    m.rating,
    array_agg(a.name) AS actors
FROM
    movies m
        JOIN
    actors a ON a.id = ANY(m.actors_id)
WHERE
    m.title LIKE '%' || vibor || '%'
   OR EXISTS (
    SELECT 1
    FROM
        actors a2
            JOIN
        movies m2 ON a2.id = ANY(m2.actors_id) AND m2.id = m.id
    WHERE
        a2.name LIKE '%' || 'vibor' || '%'
)
GROUP BY
    m.id, m.title, m.description, m.release_date, m.rating