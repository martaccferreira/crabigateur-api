
CREATE TYPE word_types AS ENUM ('regular', 'irregular', 'verb');
CREATE TYPE genders AS ENUM ('m', 'f');
CREATE TYPE grammatical_number AS ENUM ('singular', 'plural');

DROP TABLE IF EXISTS Users; 
CREATE TABLE Users (
    user_id VARCHAR(255),
    username VARCHAR (50),
    email VARCHAR (255) not null,
    level INT default 1,
    date_joined DATE,
    PRIMARY KEY (user_id)
);

DROP TABLE IF EXISTS Cards; 
CREATE TABLE Cards (
    card_id SERIAL,
    word VARCHAR(50) not null,
    translation JSONB not null,
    word_type word_types default 'regular',
    gender genders,
    level INT not null,
    PRIMARY KEY (card_id)
); 

DROP TABLE IF EXISTS Conjugations; 
CREATE TABLE Conjugations (
    conjugation_id SERIAL,
    card_id INT not null,
    tense VARCHAR(50) not null,
    forms JSONB not null,
    irregular BOOLEAN default FALSE,
    PRIMARY KEY (conjugation_id),
    FOREIGN KEY (card_id) REFERENCES Cards(card_id)
);

DROP TABLE IF EXISTS Forms; 
CREATE TABLE Forms (
    form_id SERIAL,
    card_id INT not null,
    gender genders not null,
    number grammatical_number not null,
    form VARCHAR(50) not null,
    PRIMARY KEY (form_id),
    FOREIGN KEY (card_id) REFERENCES Cards(card_id)
);

DROP TABLE IF EXISTS SRSStages; 
CREATE TABLE SRSStages (
    stage_id INT CHECK (stage_id BETWEEN 1 AND 9),
    stage_name VARCHAR(12) not null,
    stage_interval INTERVAL,
    stage_penalty INT,
    PRIMARY KEY (stage_id)
);

DROP TABLE IF EXISTS UserCardStatus; 
CREATE TABLE UserCardStatus (
    user_id VARCHAR(255),
    card_id INT,
    stage_id INT not null,
    next_review_date TIMESTAMPTZ,
    PRIMARY KEY (user_id, card_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (card_id) REFERENCES Cards(card_id),
    FOREIGN KEY (stage_id) REFERENCES SRSStages(stage_id)
);

DROP TABLE IF EXISTS Reviews; 
CREATE TABLE Reviews (
    review_id SERIAL,
    user_id VARCHAR(255) not null,
    card_id INT not null,
    review_date DATE,
    success BOOLEAN not null,
    previous_stage INT CHECK (previous_stage BETWEEN 1 AND 9),
    PRIMARY KEY (review_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (card_id) REFERENCES Cards(card_id)
);


