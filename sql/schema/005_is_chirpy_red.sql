-- +goose Up
ALTER TABLE chirps 
ADD is_chirpy_red BOOLEAN
DEFAULT FALSE;

-- +goose Down
ALTER TABLE chirps
DROP COLUMN is_chirpy_red;

