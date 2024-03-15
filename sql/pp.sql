SELECT
    a.id AS actor_id,
    a.name AS actor_name,
    a.sex AS actor_sex,
    a.birthday AS actor_birthday,
    json_agg(m.title) AS movies
FROM actors a
         LEFT JOIN movies m ON a.id = ANY(m.actors_id) AND m.deleted_at IS NULL
WHERE a.deleted_at IS NULL -- Exclude actors with non-null deleted_at
GROUP BY a.id, a.name, a.sex, a.birthday;