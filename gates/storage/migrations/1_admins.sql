-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admins(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL
);

-- Добавляем первого админа
INSERT INTO admins (user_id, role) 
SELECT '104186268', 'creator'
WHERE NOT EXISTS (
    SELECT 1 FROM admins WHERE user_id = '104186268'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS admins;
-- +goose StatementEnd 
