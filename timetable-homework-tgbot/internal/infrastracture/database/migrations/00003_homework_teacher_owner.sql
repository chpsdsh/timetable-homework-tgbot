-- +goose Up
ALTER TABLE homeworks
    ADD COLUMN IF NOT EXISTS teacher_id BIGINT;

ALTER TABLE homeworks
    ADD CONSTRAINT fk_homeworks_teacher
        FOREIGN KEY (teacher_id) REFERENCES teachers (id)
            ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_homeworks_teacher_id
    ON homeworks (teacher_id);

-- +goose Down
DROP INDEX IF EXISTS idx_homeworks_teacher_id;

ALTER TABLE homeworks
    DROP CONSTRAINT IF EXISTS fk_homeworks_teacher;

ALTER TABLE homeworks
    DROP COLUMN IF EXISTS teacher_id;
