-- Usuarios por defecto (contraseñas hasheadas con bcrypt, costo 12):
--   captador / captador123  -> Captador de Terreno
--   admin    / admin123     -> Administrativo
INSERT INTO usuarios (username, password_hash, nombre, rol) VALUES
    ('captador', '$2a$12$OBBo9qlBRxms5DGoFRzUM.2hBd8rmnTIJiF6Wfw9khdcI4/s0LTK.', 'Captador de Terreno', 'CAPTADOR'),
    ('admin',    '$2a$12$8cvzj7i84DRBBYSVPuDTwOZq1cemPzyrVdPsh03Vbbbi4ZHN3Ansy', 'Administrativo',      'ADMINISTRATIVO');
