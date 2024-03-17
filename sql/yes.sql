-- Вставка данных о фильмах с указанием массива идентификаторов актеров
INSERT INTO movies (title, description, release_date, rating, actors_id)
VALUES ('Фильм 1', 'Описание фильма 1', '2022-01-01', 8.5, ARRAY[1, 2]), -- Фильм 1 с участием актеров с идентификаторами 1 и 2
       ('Фильм 3', 'Описание фильма 2', '2023-05-15', 0.0, ARRAY[2]);   -- Фильм 2 с участием только актера с идентификатором 2


INSERT INTO users (email, password) VALUES ('zhiborkin@ya.ru', 'kekus123'), ('kekus@mail.ru', 'lolip12344')

DROP TABLE users;