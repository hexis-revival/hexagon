CREATE TABLE users (
    id serial NOT NULL PRIMARY KEY,
    name character varying(32) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    country character varying(2) NOT NULL DEFAULT 'XX' ,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    latest_activity timestamp with time zone NOT NULL DEFAULT now(),
    restricted boolean DEFAULT false NOT NULL,
    activated boolean DEFAULT false NOT NULL
);

CREATE TABLE stats (
    user_id int NOT NULL PRIMARY KEY REFERENCES users (id),
    rank bigint NOT NULL DEFAULT 0,
    total_score bigint NOT NULL DEFAULT 0,
    ranked_score bigint NOT NULL DEFAULT 0,
    playcount bigint NOT NULL DEFAULT 0,
    playtime bigint NOT NULL DEFAULT 0,
    accuracy numeric NOT NULL DEFAULT 0,
    max_combo bigint NOT NULL DEFAULT 0,
    total_hits bigint NOT NULL DEFAULT 0,
    xh_count bigint NOT NULL DEFAULT 0,
    x_count bigint NOT NULL DEFAULT 0,
    sh_count bigint NOT NULL DEFAULT 0,
    s_count bigint NOT NULL DEFAULT 0,
    a_count bigint NOT NULL DEFAULT 0,
    b_count bigint NOT NULL DEFAULT 0,
    c_count bigint NOT NULL DEFAULT 0,
    d_count bigint NOT NULL DEFAULT 0
);

CREATE TYPE relationship_status AS ENUM (
    'friend',
    'blocked'
);

CREATE TABLE relationships
(
    user_id int NOT NULL REFERENCES users (id),
    target_id int NOT NULL REFERENCES users (id),
    status relationship_status NOT NULL,
    PRIMARY KEY (user_id, target_id)
);

CREATE INDEX idx_users_id ON users (id);
CREATE INDEX idx_users_name ON users (name);
CREATE INDEX idx_stats_user_id ON stats (user_id);

CREATE INDEX idx_relationships_user_id ON relationships (user_id);
CREATE INDEX idx_relationships_target_id ON relationships (target_id);