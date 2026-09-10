-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL 
    CONSTRAINT red_default DEFAULT false;