CREATE TABLE areas (
    id     INTEGER PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE cursos (
    id      INTEGER PRIMARY KEY,
    area_id INTEGER NOT NULL REFERENCES areas(id),
    nombre  VARCHAR(150) NOT NULL,
    UNIQUE (area_id, nombre)
);
