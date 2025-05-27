CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS webpage_requests (
    request_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NULL,
    status_code int NOT NULL,
    html_version TEXT NULL,
    title TEXT NULL,
    headings JSONB,
    internal_links_number INT NULL,
    external_links_number INT NULL,
    inaccessible_links_number INT NULL,
    contains_login_form BOOLEAN NOT NULL,
    error_description TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);