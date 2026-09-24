-- Создание схемы приложения
CREATE SCHEMA IF NOT EXISTS todoapp;

-- =============================================================================
-- Таблица пользователей
-- =============================================================================
CREATE TABLE IF NOT EXISTS todoapp.users (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version      BIGINT NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(16) CHECK (phone_number ~ '^\+?[0-9]{10,15}$')
);

COMMENT ON TABLE todoapp.users IS 'Пользователи системы';
COMMENT ON COLUMN todoapp.users.version IS 'Версия записи для оптимистичной блокировки';
COMMENT ON COLUMN todoapp.users.phone_number IS 'Номер телефона в международном формате (с плюсом или без)';

-- =============================================================================
-- Таблица задач
-- =============================================================================
CREATE TABLE IF NOT EXISTS todoapp.tasks (
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    version        BIGINT NOT NULL DEFAULT 1,
    title          VARCHAR(100) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    description    VARCHAR(1000) CHECK (char_length(description) BETWEEN 1 AND 1000),
    completed      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at   TIMESTAMPTZ,
    author_user_id BIGINT NOT NULL,

    -- Логическая целостность статуса задачи и даты завершения
    CONSTRAINT chk_tasks_completion_state CHECK (
        (completed = FALSE AND completed_at IS NULL)
        OR
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
    ),

    -- Связь с пользователем (каскадное удаление при удалении владельца)
    CONSTRAINT fk_tasks_author FOREIGN KEY (author_user_id)
        REFERENCES todoapp.users (id)
        ON DELETE CASCADE
);

-- Индекс для ускорения выборок задач пользователя и оптимизации ON DELETE CASCADE
CREATE INDEX IF NOT EXISTS idx_tasks_author_user_id 
    ON todoapp.tasks (author_user_id);

COMMENT ON TABLE todoapp.tasks IS 'Задачи пользователей';
COMMENT ON COLUMN todoapp.tasks.version IS 'Версия записи для оптимистичной блокировки';
COMMENT ON COLUMN todoapp.tasks.completed_at IS 'Дата фактического закрытия задачи (не может быть раньше created_at)';