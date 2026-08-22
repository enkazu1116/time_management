-- name: CreateWorkType :one
INSERT INTO work_types (
  work_type_id,
  work_type_name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateWorkType :one
UPDATE work_types
SET
  work_type_name = $2,
  updated_at = now()
WHERE work_type_id = $1
RETURNING *;

-- name: DeleteWorkType :exec
DELETE FROM work_types
WHERE work_type_id = $1;

-- name: FindWorkTypeByID :one
SELECT *
FROM work_types
WHERE work_type_id = $1;
