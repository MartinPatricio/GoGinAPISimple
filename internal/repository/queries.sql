-- name: CreateUser :one
INSERT INTO tblUsers (
    IdRol, NameUser, Email, LastName, Password, LastActivitie
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM tblUsers
WHERE idUser = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM tblUsers
WHERE Email = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM tblUsers
WHERE idUser = $1;

-- name: GetAllUsers :many
SELECT * FROM tblUsers
ORDER BY idUser
LIMIT $1
OFFSET $2;

-- name: GetUsersWithFilters :many
SELECT * FROM tblUsers
WHERE
    (sqlc.narg(name_filter)::VARCHAR IS NULL OR NameUser ILIKE sqlc.narg(name_filter))
AND
    (sqlc.narg(email_filter)::VARCHAR IS NULL OR Email ILIKE sqlc.narg(email_filter))
ORDER BY idUser
LIMIT $1
OFFSET $2;