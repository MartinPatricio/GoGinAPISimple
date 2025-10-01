-- internal/repository/queries.sql

-- name: CreateUser :one
INSERT INTO tblusers (
    "idRol", "NameUser", "Email", "LastName", "Password", "LastActivitie", "DateCreated"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING "idUser", "idRol", "NameUser", "Email", "LastName", "DateCreated", "LastActivitie", "Password";

-- name: GetUserByID :one
SELECT "idUser", "idRol", "NameUser", "Email", "LastName", "DateCreated", "LastActivitie", "Password" FROM tblusers
WHERE "idUser" = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT "idUser", "idRol", "NameUser", "Email", "LastName", "DateCreated", "LastActivitie", "Password" FROM tblusers
WHERE "Email" = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM tblusers
WHERE "idUser" = $1;

-- name: GetAllUsers :many
SELECT "idUser", "idRol", "NameUser", "Email", "LastName", "DateCreated", "LastActivitie", "Password" FROM tblusers
ORDER BY "idUser"
LIMIT $1
OFFSET $2;

-- name: GetUsersWithFilters :many
SELECT "idUser", "idRol", "NameUser", "Email", "LastName", "DateCreated", "LastActivitie", "Password" FROM tblusers
WHERE
    (sqlc.narg(name_filter)::VARCHAR IS NULL OR "NameUser" ILIKE sqlc.narg(name_filter))
AND
    (sqlc.narg(email_filter)::VARCHAR IS NULL OR "Email" ILIKE sqlc.narg(email_filter))
ORDER BY "idUser"
LIMIT $1
OFFSET $2;