-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Change CHARACTER SET and COLLATE to support UTF-8!
ALTER DATABASE betting CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin;

-- A competition is what we're betting on, this might be "Eurovision Song
-- Contest 2020". A competition can be locked wich means no bets for that
-- competition may be added or changed.
CREATE TABLE competition (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    image       VARCHAR(100),
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    locked      TINYINT(1) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- A competitor is a team, player or other competing in the competition. This
-- might be "Sweden".
CREATE TABLE competitor (
    id              INT PRIMARY KEY AUTO_INCREMENT,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    image           VARCHAR(100),
    name            VARCHAR(100) NOT NULL,
    description     VARCHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- A linking between a competition and a competitor. A competitor can only be
-- linked to one competition once but the same competitor can be linked to
-- multiple competitions.
CREATE TABLE competition_competitor (
    id              INT PRIMARY KEY AUTO_INCREMENT,
    id_competition  INT NOT NULL,
    id_competitor   INT NOT NULL,

    FOREIGN KEY (id_competition) REFERENCES competition(id),
    FOREIGN KEY (id_competitor) REFERENCES competitor(id),

    CONSTRAINT idx_competition_competitor UNIQUE (id_competition, id_competitor)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- A better is someone watching the competition who may cast bets and add notes
-- to the competitor. This might be "John Doe".
CREATE TABLE better (
    id          INT PRIMARY KEY AUTO_INCREMENT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    confirmed   TINYINT(1) DEFAULT 0,
    image       VARCHAR(100),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- A bet is a betters note or bet for a one competitor in a specific
-- competition. This might be "John Doe" adding a note for "Sweden" in the
-- competition "Eurovision Song Contest 2020" saying "Bad keyboard!". Only one
-- bet for each competitor in each competition may be added.
CREATE TABLE bet (
    id                          INT PRIMARY KEY AUTO_INCREMENT,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    id_better                   INT NOT NULL,
    id_competition_competitor   INT NOT NULL,
    placing                     INT,
    note                        VARCHAR(255),

    FOREIGN KEY (id_better) REFERENCES better(id),
    FOREIGN KEY (id_competition_competitor) REFERENCES competition_competitor(id),

    CONSTRAINT idx_better_competition_competitor UNIQUE (id_better, id_competition_competitor)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_bin;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE bet;
DROP TABLE better;
DROP TABLE competition_competitor;
DROP TABLE competitor;
DROP TABLE competition;
