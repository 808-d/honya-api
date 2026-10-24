-- name: UserHasFeature :one
SELECT EXISTS (
    SELECT 1
    FROM users u
    JOIN user_roles ur    ON ur.user_id = u.id
    JOIN roles r          ON r.id = ur.role_id AND r.is_delete = false
    JOIN role_features rf ON rf.role_id = r.id
    JOIN features f       ON f.id = rf.feature_id AND f.is_delete = false
    WHERE u.public_id::text = sqlc.arg(public_id)::text
      AND u.is_delete = false
      AND f.code = sqlc.arg(code)::text
) AS has_feature;
