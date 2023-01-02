-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tracked_subjects (
    ID TEXT PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS base_subject_permissions (
    subject_id TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    action TEXT NOT NULL,
    directory UUID NOT NULL,
    FOREIGN KEY (subject_id) REFERENCES tracked_subjects(ID) ON DELETE CASCADE,
    FOREIGN KEY (directory) REFERENCES tracked_directories(id) ON DELETE CASCADE,
    PRIMARY KEY (subject_id, endpoint, action, directory)
);

CREATE TABLE IF NOT EXISTS effective_subject_permissions (
    subject_id TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    action TEXT NOT NULL,
    directory UUID NOT NULL,
    FOREIGN KEY (subject_id) REFERENCES tracked_subjects(ID) ON DELETE CASCADE,
    FOREIGN KEY (directory) REFERENCES tracked_directories(id) ON DELETE CASCADE,
    PRIMARY KEY (subject_id, endpoint, action, directory)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracked_subjects;
DROP TABLE IF EXISTS base_subject_permissions;
-- +goose StatementEnd
