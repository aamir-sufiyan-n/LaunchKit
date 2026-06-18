-- name: AddProjectMember :one
INSERT INTO project_members (project_id, user_id, role, status)
VALUES ($1, $2, $3, $4)
ON CONFLICT (project_id, user_id) DO UPDATE
    SET role = $3, status = $4, updated_at = NOW()
RETURNING *;

-- name: GetProjectMember :one
SELECT * FROM project_members
WHERE project_id = $1 AND user_id = $2;

-- name: ListProjectMembers :many
SELECT pm.*, u.wallet_address, u.display_name, u.avatar_url
FROM project_members pm
JOIN users u ON u.id = pm.user_id
WHERE pm.project_id = $1
ORDER BY pm.created_at ASC;

-- name: ListProjectsForUser :many
SELECT p.*, pm.role, pm.status
FROM project_members pm
JOIN projects p ON p.id = pm.project_id
WHERE pm.user_id = $1 AND pm.status = 'active';

-- name: UpdateMemberRole :one
UPDATE project_members
SET role = $3, updated_at = NOW()
WHERE project_id = $1 AND user_id = $2
RETURNING *;

-- name: UpdateMemberStatus :one
UPDATE project_members
SET status = $3, updated_at = NOW()
WHERE project_id = $1 AND user_id = $2
RETURNING *;

-- name: GetProjectMembership :one
SELECT project_id, role
FROM project_members
WHERE user_id = $1 AND status = 'active'
LIMIT 1;

-- name: RemoveProjectMember :exec
DELETE FROM project_members
WHERE project_id = $1 AND user_id = $2;
