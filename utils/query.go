package utils

import "github.com/google/uuid"

const FindSecretsByEnvIdQuery = `
SELECT * 
FROM (
	SELECT 
		s.id,
		s.user_id,
		s.key,
		s.value,
		s.created_at,
		s.updated_at,
		jsonb_agg(envs) as environments
	FROM secrets s
	JOIN environment_secrets es ON s.id = es.secret_id
	JOIN environments envs on es.environment_id = envs.id
	WHERE s.user_id = ?
	GROUP BY s.id
) r
WHERE r.environments @> ?;
`

func GenerateJSONIDString(id uuid.UUID) string {
	return `[{"id":"` + id.String() + `"}]`
}

func GenerateFindSecretByEnvIdsQuery(ids []uuid.UUID) string {
	var queryEnvironments string
	for _, value := range ids {
		if len(queryEnvironments) == 0 {
			queryEnvironments += `r.environments @> '[{"id": "` + value.String() + `"}]'`
		} else {
			queryEnvironments += `OR r.environments @> '[{"id": "` + value.String() + `"}]'`
		}
	}

	RAWSQL := `
	       SELECT *
	       FROM (
	        SELECT
		        s.id,
		        s.user_id,
		        s.key,
		        s.value,
		        s.created_at,
		        s.updated_at,
		        jsonb_agg(envs) as environments
	        FROM secrets s
	        JOIN environment_secrets es ON s.id = es.secret_id
	        JOIN environments envs on es.environment_id = envs.id
	        WHERE s.user_id = ?
	        GROUP BY s.id
	       ) r
	       WHERE `

	RAWSQL += "(" + queryEnvironments + ") AND r.key = ?"

	return RAWSQL
}
