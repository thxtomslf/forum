ALTER
USER postgres WITH ENCRYPTED PASSWORD 'admin';

DROP SCHEMA IF EXISTS dbforum CASCADE;

CREATE
EXTENSION IF NOT EXISTS citext;
CREATE SCHEMA dbforum;

CREATE
UNLOGGED TABLE dbforum.users
(
    id       BIGSERIAL PRIMARY KEY NOT NULL,

    nickname CITEXT UNIQUE         NOT NULL,
    fullname TEXT                  NOT NULL,
    about    TEXT                  NOT NULL,
    email    CITEXT UNIQUE         NOT NULL
);


CREATE
UNLOGGED TABLE dbforum.forum
(
    id            BIGSERIAL PRIMARY KEY NOT NULL,
    user_nickname CITEXT                NOT NULL,

    title         TEXT                  NOT NULL,
    slug          CITEXT UNIQUE         NOT NULL,
    posts         BIGINT DEFAULT 0      NOT NULL,
    threads       INT    DEFAULT 0      NOT NULL,

    FOREIGN KEY (user_nickname)
        REFERENCES dbforum.users (nickname)
);


CREATE
UNLOGGED TABLE dbforum.thread
(
    id              BIGSERIAL PRIMARY KEY    NOT NULL,
    forum_slug      CITEXT                   NOT NULL,
    author_nickname CITEXT                   NOT NULL,

    title           TEXT                     NOT NULL,
    message         TEXT                     NOT NULL,
    votes           INT DEFAULT 0            NOT NULL,
    slug            citext UNIQUE,
    created         TIMESTAMP WITH TIME ZONE NOT NULL,

    FOREIGN KEY (forum_slug)
        REFERENCES dbforum.forum (slug),
    FOREIGN KEY (author_nickname)
        REFERENCES dbforum.users (nickname)
);

CREATE
UNLOGGED TABLE dbforum.votes
(
    nickname  CITEXT        NOT NULL,
    voice     INT DEFAULT 0 NOT NULL,
    thread_id BIGINT        NOT NULL,
    PRIMARY KEY (nickname, thread_id),

    FOREIGN KEY (nickname)
        REFERENCES dbforum.users (nickname),
    FOREIGN KEY (thread_id)
        REFERENCES dbforum.thread (id)
);


CREATE
UNLOGGED TABLE dbforum.post
(
    id              BIGSERIAL PRIMARY KEY               NOT NULL,
    author_nickname CITEXT                              NOT NULL,
    forum_slug      CITEXT                              NOT NULL,
    thread_id       BIGINT                              NOT NULL,
    message         TEXT                                NOT NULL,

    parent          BIGINT   DEFAULT 0                  NOT NULL,
    is_edited       BOOLEAN  DEFAULT false              NOT NULL,
    created         TIMESTAMP WITH TIME ZONE            NOT NULL,
    tree            BIGINT[] DEFAULT ARRAY []::BIGINT[] NOT NULL,

    FOREIGN KEY (author_nickname)
        REFERENCES dbforum.users (nickname),
    FOREIGN KEY (forum_slug)
        REFERENCES dbforum.forum (slug),
    FOREIGN KEY (thread_id)
        REFERENCES dbforum.thread (id)
);

CREATE
UNLOGGED TABLE dbforum.forum_users
(
    forum_slug CITEXT NOT NULL,
    nickname   CITEXT NOT NULL,
    fullname   TEXT   NOT NULL,
    about      TEXT   NOT NULL,
    email      TEXT   NOT NULL,

    FOREIGN KEY (nickname)
        REFERENCES dbforum.users (nickname),
    FOREIGN KEY (forum_slug)
        REFERENCES dbforum.forum (slug),

    PRIMARY KEY (nickname, forum_slug)
);
