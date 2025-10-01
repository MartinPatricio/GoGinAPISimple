-- sql/schema.sql

CREATE TABLE tblusers (
    "idUser" SERIAL PRIMARY KEY,
    "idRol" SMALLINT NOT NULL,
    "NameUser" VARCHAR(30) NOT NULL,
    "Email" VARCHAR(80) UNIQUE NOT NULL,
    "LastName" VARCHAR(40) NOT NULL,
    "DateCreated" DATE NOT NULL DEFAULT CURRENT_DATE,
    "LastActivitie" DATE NOT NULL DEFAULT CURRENT_DATE,
    "Password" VARCHAR(255) NOT NULL
);