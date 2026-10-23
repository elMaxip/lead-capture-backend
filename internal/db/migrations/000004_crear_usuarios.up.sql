CREATE TABLE usuarios (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL,               -- bcrypt
    nombre        VARCHAR(100) NOT NULL,
    rol           VARCHAR(20)  NOT NULL CHECK (rol IN ('CAPTADOR','ADMINISTRATIVO'))
);

CREATE TABLE sesiones (
    token_hash CHAR(64) PRIMARY KEY,                    -- SHA-256 del token entregado a la app
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    creada_en  TIMESTAMPTZ NOT NULL DEFAULT now()
);
