CREATE TABLE prospectos (
    id              UUID PRIMARY KEY,                  -- generado en el dispositivo
    nombre          VARCHAR(100) NOT NULL,
    email           VARCHAR(150) NOT NULL,
    telefono        VARCHAR(12)  NOT NULL,             -- formato +569XXXXXXXX
    nivel_educativo VARCHAR(30)  NOT NULL CHECK (nivel_educativo IN
                      ('BASICA','MEDIA','TECNICO_PROFESIONAL','UNIVERSITARIA','POSTGRADO')),
    curso_id        INTEGER NOT NULL REFERENCES cursos(id),
    modalidad       VARCHAR(20) NOT NULL CHECK (modalidad IN ('ONLINE','PRESENCIAL')),
    canal_contacto  VARCHAR(20) NOT NULL CHECK (canal_contacto IN ('EMAIL','TELEFONO','WHATSAPP')),
    observaciones   VARCHAR(500),
    capturado_en    TIMESTAMPTZ NOT NULL,              -- fecha/hora de captura en terreno
    recibido_en     TIMESTAMPTZ NOT NULL DEFAULT now() -- fecha/hora de llegada al servidor
);

CREATE INDEX idx_prospectos_capturado_en ON prospectos (capturado_en DESC);
