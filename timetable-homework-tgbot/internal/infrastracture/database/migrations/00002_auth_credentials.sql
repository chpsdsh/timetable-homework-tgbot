-- +goose Up
ALTER TABLE teachers
    ADD COLUMN IF NOT EXISTS login TEXT,
    ADD COLUMN IF NOT EXISTS password TEXT;

ALTER TABLE students
    ADD COLUMN IF NOT EXISTS login TEXT,
    ADD COLUMN IF NOT EXISTS password TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_teachers_login
    ON teachers (login)
    WHERE login IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_students_login
    ON students (login)
    WHERE login IS NOT NULL;

UPDATE teachers
SET login = COALESCE(login, 'teacher_' || id::TEXT),
    password = COALESCE(password, 'teacher_' || id::TEXT)
WHERE teacher_name IS NOT NULL;

INSERT INTO groups (group_name)
VALUES ('23204'), ('23203')
ON CONFLICT (group_name) DO NOTHING;

INSERT INTO users (tg_id)
VALUES (-300000000000001), (-300000000000002)
ON CONFLICT DO NOTHING;

INSERT INTO students (user_id, group_id, login, password)
SELECT -300000000000001, id, 'student1', 'student1'
FROM groups
WHERE group_name = '23204'
ON CONFLICT (user_id) DO UPDATE
SET group_id = EXCLUDED.group_id,
    login = EXCLUDED.login,
    password = EXCLUDED.password;

INSERT INTO students (user_id, group_id, login, password)
SELECT -300000000000002, id, 'student2', 'student2'
FROM groups
WHERE group_name = '23203'
ON CONFLICT (user_id) DO UPDATE
SET group_id = EXCLUDED.group_id,
    login = EXCLUDED.login,
    password = EXCLUDED.password;

-- +goose Down
DROP INDEX IF EXISTS idx_students_login;
DROP INDEX IF EXISTS idx_teachers_login;

DELETE FROM students
WHERE login IN ('student1', 'student2');

DELETE FROM users
WHERE tg_id IN (-300000000000001, -300000000000002);

ALTER TABLE students
    DROP COLUMN IF EXISTS login,
    DROP COLUMN IF EXISTS password;

ALTER TABLE teachers
    DROP COLUMN IF EXISTS login,
    DROP COLUMN IF EXISTS password;
