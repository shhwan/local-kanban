CREATE TABLE labels (
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(20) NOT NULL UNIQUE,
    color VARCHAR(7)  NOT NULL            -- HEX形式: #RRGGBB
);

CREATE TABLE stages (
    id       SERIAL PRIMARY KEY,
    name     VARCHAR(20) NOT NULL UNIQUE,
    position INTEGER     NOT NULL          -- 表示順 (昇順、左から右)
);

CREATE TABLE tasks (
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    label_id   INTEGER      NOT NULL REFERENCES labels(id),
    stage_id   INTEGER      NOT NULL REFERENCES stages(id),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- 追記のみ・イミュータブル運用を想定
CREATE TABLE work_logs (
    id         SERIAL PRIMARY KEY,
    task_id    INTEGER   NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- work_logs と同様にイミュータブル運用
CREATE TABLE notes (
    id         SERIAL PRIMARY KEY,
    task_id    INTEGER   NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- インデックス
CREATE INDEX idx_tasks_label_id ON tasks (label_id);
CREATE INDEX idx_tasks_stage_id ON tasks (stage_id);
CREATE INDEX idx_work_logs_task_id_created_at ON work_logs (task_id, created_at);
CREATE INDEX idx_notes_task_id_created_at ON notes (task_id, created_at);

-- 初期データ
INSERT INTO labels (name, color) VALUES
    ('FRONTEND', '#3B82F6'),
    ('BACKEND',  '#10B981'),
    ('INFRA',    '#F59E0B');

INSERT INTO stages (name, position) VALUES
    ('TODO',  0),
    ('DOING', 1),
    ('DONE',  2);
