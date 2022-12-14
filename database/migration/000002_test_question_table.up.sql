BEGIN;
CREATE TABLE IF NOT EXISTS tests(
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(256) NOT NULL,
    description TEXT,
    is_available BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY(id)
);

CREATE TABLE  IF NOT EXISTS questions(
    id INT NOT NULL AUTO_INCREMENT,
    title TEXT NOT NULL,
    type VARCHAR(256) NOT NULL DEFAULT 'mcq',
    test_id int NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(test_id) REFERENCES tests(id)
);

CREATE TABLE  IF NOT EXISTS options(
    id INT NOT NULL AUTO_INCREMENT,
    title TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    question_id int NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(question_id) REFERENCES questions(id)
);

CREATE TABLE  IF NOT EXISTS user_test_report(
    id INT NOT NULL AUTO_INCREMENT,
    test_id int NOT NULL,
    user_id int NOT NULL,
    has_passed BOOLEAN NOT NULL,    
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(test_id) REFERENCES tests(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);

COMMIT;