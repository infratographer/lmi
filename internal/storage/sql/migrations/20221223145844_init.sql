-- +goose Up
-- +goose StatementBegin

-- tracked_subjects table
-- It stores the subjects that are tracked by the
-- authorization service
CREATE TABLE IF NOT EXISTS tracked_subjects (
    subject_id TEXT NOT NULL PRIMARY KEY
);

-- Permissions table
-- It stores the permissions that can be allowed
-- if a subject is assigned to a role
CREATE TABLE IF NOT EXISTS permissions (
    target TEXT NOT NULL PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Roles table
-- It stores the roles that can be assigned to a subject
-- NOTE/TODO: Should we add a column for the directory
--           where the role is defined? This way we can
--           have roles with the same name but different
--           permissions depending on the directory
--           as well as global roles people can use.
CREATE TABLE IF NOT EXISTS roles (
    role_id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- role_permissions table
-- It stores the permissions that are assigned to a role
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL,
    target TEXT NOT NULL,
    PRIMARY KEY (role_id, target),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (target) REFERENCES permissions(target) ON DELETE CASCADE
);

-- role_assignments tab;
-- It stores the subjects that are assigned to a rol;
CREATE TABLE IF NOT EXISTS role_assignments (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    role_id UUID NOT NULL,
    subject_id TEXT NOT NULL,
    scope UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (subject_id) REFERENCES tracked_subjects(subject_id) ON DELETE CASCADE,
    FOREIGN KEY (scope) REFERENCES tracked_directories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS effective_permissions (
    subject_id TEXT NOT NULL,
    target TEXT NOT NULL,
    scope UUID NOT NULL,
    from_role UUID NOT NULL,
    PRIMARY KEY (subject_id, target, scope),
    FOREIGN KEY (subject_id) REFERENCES tracked_subjects(subject_id) ON DELETE CASCADE,
    FOREIGN KEY (target) REFERENCES permissions(target) ON DELETE CASCADE,
    FOREIGN KEY (scope) REFERENCES tracked_directories(id) ON DELETE CASCADE,
    FOREIGN KEY (from_role) REFERENCES roles(role_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracked_subjects;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS role_assignments;
-- +goose StatementEnd
