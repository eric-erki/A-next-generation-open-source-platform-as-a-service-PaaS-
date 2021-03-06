package data

import (
	"github.com/flynn/que-go"
	"github.com/jackc/pgx"
)

var preparedStatements = map[string]string{
	"ping":                                  pingQuery,
	"app_list":                              appListQuery,
	"app_select_by_name":                    appSelectByNameQuery,
	"app_select_by_name_for_update":         appSelectByNameForUpdateQuery,
	"app_select_by_name_or_id":              appSelectByNameOrIDQuery,
	"app_select_by_name_or_id_for_update":   appSelectByNameOrIDForUpdateQuery,
	"app_insert":                            appInsertQuery,
	"app_update_strategy":                   appUpdateStrategyQuery,
	"app_update_meta":                       appUpdateMetaQuery,
	"app_update_release":                    appUpdateReleaseQuery,
	"app_update_deploy_timeout":             appUpdateDeployTimeoutQuery,
	"app_delete":                            appDeleteQuery,
	"app_next_name_id":                      appNextNameIDQuery,
	"app_get_release":                       appGetReleaseQuery,
	"release_list":                          releaseListQuery,
	"release_select":                        releaseSelectQuery,
	"release_insert":                        releaseInsertQuery,
	"release_app_list":                      releaseAppListQuery,
	"release_artifacts_insert":              releaseArtifactsInsertQuery,
	"release_artifacts_delete":              releaseArtifactsDeleteQuery,
	"release_delete":                        releaseDeleteQuery,
	"artifact_list":                         artifactListQuery,
	"artifact_list_ids":                     artifactListIDsQuery,
	"artifact_select":                       artifactSelectQuery,
	"artifact_select_by_type_and_uri":       artifactSelectByTypeAndURIQuery,
	"artifact_insert":                       artifactInsertQuery,
	"artifact_delete":                       artifactDeleteQuery,
	"artifact_release_count":                artifactReleaseCountQuery,
	"artifact_layer_count":                  artifactLayerCountQuery,
	"deployment_list":                       deploymentListQuery,
	"deployment_select":                     deploymentSelectQuery,
	"deployment_insert":                     deploymentInsertQuery,
	"deployment_update_finished_at":         deploymentUpdateFinishedAtQuery,
	"deployment_update_finished_at_now":     deploymentUpdateFinishedAtNowQuery,
	"deployment_delete":                     deploymentDeleteQuery,
	"event_select":                          eventSelectQuery,
	"event_insert":                          eventInsertQuery,
	"event_insert_op":                       eventInsertOpQuery,
	"event_insert_unique":                   eventInsertUniqueQuery,
	"formation_list_by_app":                 formationListByAppQuery,
	"formation_list_by_release":             formationListByReleaseQuery,
	"formation_list_active":                 formationListActiveQuery,
	"formation_list_since":                  formationListSinceQuery,
	"formation_select":                      formationSelectQuery,
	"formation_select_expanded":             formationSelectExpandedQuery,
	"formation_insert":                      formationInsertQuery,
	"formation_delete":                      formationDeleteQuery,
	"formation_delete_by_app":               formationDeleteByAppQuery,
	"scale_request_insert":                  scaleRequestInsertQuery,
	"scale_request_cancel":                  scaleRequestCancelQuery,
	"scale_request_update":                  scaleRequestUpdateQuery,
	"job_list":                              jobListQuery,
	"job_list_active":                       jobListActiveQuery,
	"job_select":                            jobSelectQuery,
	"job_insert":                            jobInsertQuery,
	"job_volume_insert":                     jobVolumeInsertQuery,
	"provider_list":                         providerListQuery,
	"provider_select_by_name":               providerSelectByNameQuery,
	"provider_select_by_name_or_id":         providerSelectByNameOrIDQuery,
	"provider_insert":                       providerInsertQuery,
	"resource_list":                         resourceListQuery,
	"resource_list_by_provider":             resourceListByProviderQuery,
	"resource_list_by_app":                  resourceListByAppQuery,
	"resource_select":                       resourceSelectQuery,
	"resource_insert":                       resourceInsertQuery,
	"resource_delete":                       resourceDeleteQuery,
	"app_resource_insert_app_by_name":       appResourceInsertAppByNameQuery,
	"app_resource_insert_app_by_name_or_id": appResourceInsertAppByNameOrIDQuery,
	"app_resource_delete_by_app":            appResourceDeleteByAppQuery,
	"app_resource_delete_by_resource":       appResourceDeleteByResourceQuery,
	"domain_migration_insert":               domainMigrationInsert,
	"backup_insert":                         backupInsert,
	"backup_update":                         backupUpdate,
	"backup_select_latest":                  backupSelectLatest,
	"sink_list":                             sinkListQuery,
	"sink_list_since":                       sinkListSinceQuery,
	"sink_select":                           sinkSelectQuery,
	"sink_insert":                           sinkInsertQuery,
	"sink_delete":                           sinkDeleteQuery,
	"volume_list":                           volumeListQuery,
	"volume_app_list":                       volumeAppListQuery,
	"volume_list_since":                     volumeListSinceQuery,
	"volume_select":                         volumeSelectQuery,
	"volume_insert":                         volumeInsertQuery,
	"volume_decommission":                   volumeDecommissionQuery,
}

func PrepareStatements(conn *pgx.Conn) error {
	for name, sql := range preparedStatements {
		if _, err := conn.Prepare(name, sql); err != nil {
			return err
		}
	}
	return que.PrepareStatements(conn)
}

const (
	// misc
	pingQuery = `SELECT 1`
	// apps
	appListQuery = `
SELECT app_id, name, meta, strategy, release_id, deploy_timeout, created_at, updated_at
FROM apps WHERE deleted_at IS NULL ORDER BY created_at DESC`
	appSelectByNameQuery = `
SELECT app_id, name, meta, strategy, release_id, deploy_timeout, created_at, updated_at
FROM apps WHERE deleted_at IS NULL AND name = $1`
	appSelectByNameForUpdateQuery = `
SELECT app_id, name, meta, strategy, release_id, deploy_timeout, created_at, updated_at
FROM apps WHERE deleted_at IS NULL AND name = $1 FOR UPDATE`
	appSelectByNameOrIDQuery = `
SELECT app_id, name, meta, strategy, release_id, deploy_timeout, created_at, updated_at
FROM apps WHERE deleted_at IS NULL AND (app_id = $1 OR name = $2) LIMIT 1`
	appSelectByNameOrIDForUpdateQuery = `
SELECT app_id, name, meta, strategy, release_id, deploy_timeout, created_at, updated_at
FROM apps WHERE deleted_at IS NULL AND (app_id = $1 OR name = $2) LIMIT 1 FOR UPDATE`
	appInsertQuery = `
INSERT INTO apps (app_id, name, meta, strategy, deploy_timeout) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at`
	appUpdateStrategyQuery = `
UPDATE apps SET strategy = $2, updated_at = now() WHERE app_id = $1`
	appUpdateMetaQuery = `
UPDATE apps SET meta = $2, updated_at = now() WHERE app_id = $1`
	appUpdateReleaseQuery = `
UPDATE apps SET release_id = $2, updated_at = now() WHERE app_id = $1
RETURNING updated_at`
	appUpdateDeployTimeoutQuery = `
UPDATE apps SET deploy_timeout = $2, updated_at = now() WHERE app_id = $1`
	appDeleteQuery = `
UPDATE apps SET deleted_at = now() WHERE app_id = $1 AND deleted_at IS NULL`
	appNextNameIDQuery = `
SELECT nextval('name_ids')`
	appGetReleaseQuery = `
SELECT r.release_id, r.app_id,
  ARRAY(
	SELECT a.artifact_id
	FROM release_artifacts a
	WHERE a.release_id = r.release_id AND a.deleted_at IS NULL
	ORDER BY a.index
  ), r.env, r.processes, r.meta, r.created_at
FROM apps a JOIN releases r USING (release_id) WHERE a.app_id = $1 AND r.deleted_at IS NULL`

	releaseListQuery = `
SELECT r.release_id, r.app_id,
  ARRAY(
	SELECT a.artifact_id
	FROM release_artifacts a
	WHERE a.release_id = r.release_id AND a.deleted_at IS NULL
	ORDER BY a.index
  ), r.env, r.processes, r.meta, r.created_at
FROM releases r WHERE r.deleted_at IS NULL ORDER BY r.created_at DESC`
	releaseSelectQuery = `
SELECT r.release_id, r.app_id,
  ARRAY(
	SELECT a.artifact_id
	FROM release_artifacts a
	WHERE a.release_id = r.release_id AND a.deleted_at IS NULL
	ORDER BY a.index
  ), r.env, r.processes, r.meta, r.created_at
FROM releases r WHERE r.release_id = $1 AND r.deleted_at IS NULL`
	releaseInsertQuery = `
INSERT INTO releases (release_id, app_id, env, processes, meta)
VALUES ($1, $2, $3, $4, $5) RETURNING created_at`
	releaseAppListQuery = `
SELECT r.release_id, r.app_id,
  ARRAY(
	SELECT a.artifact_id
	FROM release_artifacts a
	WHERE a.release_id = r.release_id AND a.deleted_at IS NULL
	ORDER BY a.index
  ), r.env, r.processes, r.meta, r.created_at
FROM releases r WHERE r.app_id = $1 AND r.deleted_at IS NULL ORDER BY r.created_at DESC`
	releaseArtifactsInsertQuery = `
INSERT INTO release_artifacts (release_id, artifact_id, index) VALUES ($1, $2, $3)`
	releaseArtifactsDeleteQuery = `
UPDATE release_artifacts SET deleted_at = now() WHERE release_id = $1 AND artifact_id = $2 AND deleted_at IS NULL`
	releaseDeleteQuery = `
UPDATE releases SET deleted_at = now() WHERE release_id = $1 AND deleted_at IS NULL`
	artifactListQuery = `
SELECT artifact_id, type, uri, meta, manifest, hashes, size, layer_url_template, created_at FROM artifacts
WHERE deleted_at IS NULL ORDER BY created_at DESC`
	artifactListIDsQuery = `
SELECT artifact_id, type, uri, meta, manifest, hashes, size, layer_url_template, created_at FROM artifacts
WHERE deleted_at IS NULL AND artifact_id = ANY($1)`
	artifactSelectQuery = `
SELECT artifact_id, type, uri, meta, manifest, hashes, size, layer_url_template, created_at FROM artifacts
WHERE artifact_id = $1 AND deleted_at IS NULL`
	artifactSelectByTypeAndURIQuery = `
SELECT artifact_id, meta, manifest, hashes, size, layer_url_template, created_at FROM artifacts WHERE type = $1 AND uri = $2 AND deleted_at IS NULL`
	artifactInsertQuery = `
INSERT INTO artifacts (artifact_id, type, uri, meta, manifest, hashes, size, layer_url_template) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING created_at`
	artifactDeleteQuery = `
UPDATE artifacts SET deleted_at = now() WHERE artifact_id = $1 AND deleted_at IS NULL`
	artifactReleaseCountQuery = `
SELECT COUNT(*) FROM release_artifacts WHERE artifact_id = $1 AND deleted_at IS NULL`
	artifactLayerCountQuery = `
SELECT COUNT(*) FROM (
  SELECT jsonb_array_elements(jsonb_array_elements(manifest->'rootfs')->'layers')->'id' AS layer_id
  FROM artifacts
  WHERE deleted_at IS NULL
) AS l WHERE l.layer_id = $1`
	deploymentInsertQuery = `
INSERT INTO deployments (deployment_id, app_id, old_release_id, new_release_id, strategy, processes, tags, deploy_timeout, deploy_batch_size)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING created_at`
	deploymentUpdateFinishedAtQuery = `
UPDATE deployments SET finished_at = $2 WHERE deployment_id = $1`
	deploymentUpdateFinishedAtNowQuery = `
UPDATE deployments SET finished_at = now() WHERE deployment_id = $1`
	deploymentDeleteQuery = `
DELETE FROM deployments WHERE deployment_id = $1`
	deploymentSelectQuery = `
WITH deployment_events AS (SELECT * FROM events WHERE object_type = 'deployment')
SELECT d.deployment_id, d.app_id, d.old_release_id, d.new_release_id,
  strategy, e1.data->>'status' AS status,
  processes, tags, deploy_timeout, deploy_batch_size, d.created_at, d.finished_at
FROM deployments d
LEFT JOIN deployment_events e1
  ON d.deployment_id = e1.object_id::uuid
LEFT OUTER JOIN deployment_events e2
  ON (d.deployment_id = e2.object_id::uuid AND e1.created_at < e2.created_at)
WHERE e2.created_at IS NULL AND d.deployment_id = $1`
	deploymentListQuery = `
WITH deployment_events AS (SELECT * FROM events WHERE object_type = 'deployment')
SELECT d.deployment_id, d.app_id, d.old_release_id, d.new_release_id,
  strategy, e1.data->>'status' AS status,
  processes, tags, deploy_timeout, deploy_batch_size, d.created_at, d.finished_at
FROM deployments d
LEFT JOIN deployment_events e1
  ON d.deployment_id = e1.object_id::uuid
LEFT OUTER JOIN deployment_events e2
  ON (d.deployment_id = e2.object_id::uuid AND e1.created_at < e2.created_at)
WHERE e2.created_at IS NULL AND d.app_id = $1 ORDER BY d.created_at DESC`
	eventSelectQuery = `
SELECT event_id, app_id, object_id, object_type, data, op, created_at
FROM events WHERE event_id = $1`
	eventInsertQuery = `
INSERT INTO events (app_id, object_id, object_type, data)
VALUES ($1, $2, $3, $4)`
	eventInsertOpQuery = `
INSERT INTO events (app_id, object_id, object_type, data, op)
VALUES ($1, $2, $3, $4, $5)`
	eventInsertUniqueQuery = `
INSERT INTO events (app_id, object_id, unique_id, object_type, data)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (unique_id) DO NOTHING`
	formationListByAppQuery = `
SELECT app_id, release_id, processes, tags, created_at, updated_at
FROM formations WHERE app_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	formationListByReleaseQuery = `
SELECT app_id, release_id, processes, tags, created_at, updated_at
FROM formations WHERE release_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	formationListActiveQuery = `
SELECT
  apps.app_id, apps.name, apps.meta, apps.strategy, apps.release_id,
  apps.deploy_timeout, apps.created_at, apps.updated_at,
  releases.release_id,
  ARRAY(
	SELECT r.artifact_id
	FROM release_artifacts r
	WHERE r.release_id = releases.release_id AND r.deleted_at IS NULL
	ORDER BY r.index
  ),
  releases.meta, releases.env, releases.processes, releases.created_at,
  scale_requests.scale_request_id, scale_requests.old_processes, scale_requests.new_processes,
  scale_requests.old_tags, scale_requests.new_tags, scale_requests.created_at,
  formations.processes, formations.tags, formations.updated_at, formations.deleted_at IS NOT NULL
FROM formations
JOIN apps USING (app_id)
JOIN releases ON releases.release_id = formations.release_id
LEFT OUTER JOIN scale_requests
  ON scale_requests.app_id = formations.app_id
  AND scale_requests.release_id = formations.release_id
  AND scale_requests.state = 'pending'
WHERE (formations.app_id, formations.release_id) IN (
  SELECT app_id, release_id
  FROM formations, json_each_text(formations.processes::json)
  WHERE processes != 'null'
  GROUP BY app_id, release_id
  HAVING SUM(value::int) > 0
)
AND formations.deleted_at IS NULL
ORDER BY formations.updated_at DESC`
	formationListSinceQuery = `
SELECT
  apps.app_id, apps.name, apps.meta, apps.strategy, apps.release_id,
  apps.deploy_timeout, apps.created_at, apps.updated_at,
  releases.release_id,
  ARRAY(
	SELECT r.artifact_id
	FROM release_artifacts r
	WHERE r.release_id = releases.release_id AND r.deleted_at IS NULL
	ORDER BY r.index
  ),
  releases.meta, releases.env, releases.processes, releases.created_at,
  scale_requests.scale_request_id, scale_requests.old_processes, scale_requests.new_processes,
  scale_requests.old_tags, scale_requests.new_tags, scale_requests.created_at,
  formations.processes, formations.tags, formations.updated_at, formations.deleted_at IS NOT NULL
FROM formations
JOIN apps USING (app_id)
JOIN releases ON releases.release_id = formations.release_id
LEFT OUTER JOIN scale_requests
  ON scale_requests.app_id = formations.app_id
  AND scale_requests.release_id = formations.release_id
  AND scale_requests.state = 'pending'
WHERE formations.updated_at >= $1 AND formations.deleted_at IS NULL
ORDER BY formations.updated_at DESC`
	formationSelectQuery = `
SELECT app_id, release_id, processes, tags, created_at, updated_at
FROM formations WHERE app_id = $1 AND release_id = $2 AND deleted_at IS NULL`
	formationSelectExpandedQuery = `
SELECT
  apps.app_id, apps.name, apps.meta, apps.strategy, apps.release_id,
  apps.deploy_timeout, apps.created_at, apps.updated_at,
  releases.release_id,
  ARRAY(
	SELECT a.artifact_id
	FROM release_artifacts a
	WHERE a.release_id = releases.release_id AND a.deleted_at IS NULL
	ORDER BY a.index
  ),
  releases.meta, releases.env, releases.processes, releases.created_at,
  scale_requests.scale_request_id, scale_requests.old_processes, scale_requests.new_processes,
  scale_requests.old_tags, scale_requests.new_tags, scale_requests.created_at,
  formations.processes, formations.tags, formations.updated_at, formations.deleted_at IS NOT NULL
FROM formations
JOIN apps USING (app_id)
JOIN releases ON releases.release_id = formations.release_id
LEFT OUTER JOIN scale_requests
  ON scale_requests.app_id = formations.app_id
  AND scale_requests.release_id = formations.release_id
  AND scale_requests.state = 'pending'
WHERE formations.app_id = $1 AND formations.release_id = $2`
	formationInsertQuery = `
INSERT INTO formations (app_id, release_id, processes, tags)
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT formations_pkey DO UPDATE
SET processes = $3, tags = $4, updated_at = now(), deleted_at = NULL
RETURNING created_at, updated_at`
	formationDeleteQuery = `
UPDATE formations SET deleted_at = now(), processes = NULL, updated_at = now()
WHERE app_id = $1 AND release_id = $2 AND deleted_at IS NULL`
	formationDeleteByAppQuery = `
UPDATE formations SET deleted_at = now(), processes = NULL, updated_at = now()
WHERE app_id = $1 AND deleted_at IS NULL`
	scaleRequestInsertQuery = `
INSERT INTO scale_requests (scale_request_id, app_id, release_id, state, old_processes, new_processes, old_tags, new_tags)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING created_at, updated_at`
	scaleRequestCancelQuery = `
WITH updated AS (
	UPDATE scale_requests SET state = 'cancelled', updated_at = now() WHERE app_id = $1 AND release_id = $2 AND state != 'cancelled'
	RETURNING *
)
SELECT scale_request_id, app_id, release_id, state, old_processes, new_processes, old_tags, new_tags, created_at, updated_at
FROM updated
ORDER BY created_at DESC`
	scaleRequestUpdateQuery = `
UPDATE scale_requests SET state = $2, updated_at = now() WHERE scale_request_id = $1
RETURNING updated_at`
	jobListQuery = `
SELECT
  cluster_id, job_id, host_id, app_id, release_id, process_type, state, meta,
  exit_status, host_error, run_at, restarts, created_at, updated_at, args,
  ARRAY(
    SELECT job_volumes.volume_id
    FROM job_volumes
    WHERE job_volumes.job_id = job_cache.job_id
    ORDER BY job_volumes.index
  )
FROM job_cache WHERE app_id = $1 ORDER BY created_at DESC`
	jobListActiveQuery = `
SELECT
  cluster_id, job_id, host_id, app_id, release_id, process_type, state, meta,
  exit_status, host_error, run_at, restarts, created_at, updated_at, args,
  ARRAY(
    SELECT job_volumes.volume_id
    FROM job_volumes
    WHERE job_volumes.job_id = job_cache.job_id
    ORDER BY job_volumes.index
  )
FROM job_cache WHERE state = 'pending' OR state = 'starting' OR state = 'up' OR state = 'stopping' ORDER BY updated_at DESC`
	jobSelectQuery = `
SELECT
  cluster_id, job_id, host_id, app_id, release_id, process_type, state, meta,
  exit_status, host_error, run_at, restarts, created_at, updated_at, args,
  ARRAY(
    SELECT job_volumes.volume_id
    FROM job_volumes
    WHERE job_volumes.job_id = job_cache.job_id
    ORDER BY job_volumes.index
  )
FROM job_cache WHERE job_id = $1`
	jobInsertQuery = `
INSERT INTO job_cache (cluster_id, job_id, host_id, app_id, release_id, process_type, state, meta, exit_status, host_error, run_at, restarts, args)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT (job_id) DO UPDATE
SET cluster_id = $1, host_id = $3, state = $7, exit_status = $9, host_error = $10, run_at = $11, restarts = $12, args = $13, updated_at = now()
RETURNING created_at, updated_at`
	jobVolumeInsertQuery = `
INSERT INTO job_volumes (job_id, volume_id, index) VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT job_volumes_pkey DO UPDATE SET index = $3
	`
	providerListQuery = `
SELECT provider_id, name, url, created_at, updated_at
FROM providers WHERE deleted_at IS NULL ORDER BY created_at DESC`
	providerSelectByNameQuery = `
SELECT provider_id, name, url, created_at, updated_at
FROM providers WHERE deleted_at IS NULL AND name = $1`
	providerSelectByNameOrIDQuery = `
SELECT provider_id, name, url, created_at, updated_at
FROM providers WHERE deleted_at IS NULL AND (provider_id = $1 OR name = $2) LIMIT 1`
	providerInsertQuery = `
INSERT INTO providers (name, url) VALUES ($1, $2)
RETURNING provider_id, created_at, updated_at`
	resourceListQuery = `
SELECT resource_id, provider_id, external_id, env,
  ARRAY(
	SELECT a.app_id
    FROM app_resources a
	WHERE a.resource_id = r.resource_id AND a.deleted_at IS NULL
	ORDER BY a.created_at DESC
  ), created_at
FROM resources r
WHERE deleted_at IS NULL
ORDER BY created_at DESC`
	resourceListByProviderQuery = `
SELECT resource_id, provider_id, external_id, env,
  ARRAY(
	SELECT a.app_id
    FROM app_resources a
	WHERE a.resource_id = r.resource_id AND a.deleted_at IS NULL
	ORDER BY a.created_at DESC
  ), created_at
FROM resources r
WHERE provider_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC`
	resourceListByAppQuery = `
SELECT DISTINCT(r.resource_id), r.provider_id, r.external_id, r.env,
  ARRAY(
    SELECT a.app_id
	FROM app_resources a
	WHERE a.resource_id = r.resource_id AND a.deleted_at IS NULL
	ORDER BY a.created_at DESC
  ), r.created_at
FROM resources r
JOIN app_resources a USING (resource_id)
WHERE a.app_id = $1 AND r.deleted_at IS NULL AND a.deleted_at IS NULL
ORDER BY r.created_at DESC`
	resourceSelectQuery = `
SELECT resource_id, provider_id, external_id, env,
  ARRAY(
    SELECT app_id
	FROM app_resources a
	WHERE a.resource_id = r.resource_id AND a.deleted_at IS NULL
	ORDER BY a.created_at DESC
  ), created_at
FROM resources r
WHERE resource_id = $1 AND deleted_at IS NULL`
	resourceInsertQuery = `
INSERT INTO resources (resource_id, provider_id, external_id, env)
VALUES ($1, $2, $3, $4) RETURNING created_at`
	resourceDeleteQuery = `
UPDATE resources SET deleted_at = now() WHERE resource_id = $1 AND deleted_at IS NULL`
	appResourceInsertAppByNameQuery = `
INSERT INTO app_resources (app_id, resource_id)
VALUES ((SELECT app_id FROM apps WHERE name = $1 AND deleted_at IS NULL), $2)
RETURNING app_id`
	appResourceInsertAppByNameOrIDQuery = `
INSERT INTO app_resources (app_id, resource_id)
VALUES ((SELECT app_id FROM apps WHERE (app_id = $1 OR name = $2) AND deleted_at IS NULL), $3)
RETURNING app_id`
	appResourceDeleteByAppQuery = `
DELETE FROM app_resources WHERE app_id = $1`
	appResourceDeleteByResourceQuery = `
DELETE FROM app_resources WHERE resource_id = $1`
	domainMigrationInsert = `
INSERT INTO domain_migrations (old_domain, domain, old_tls_cert, tls_cert) VALUES ($1, $2, $3, $4) RETURNING migration_id, created_at`
	backupInsert = `
INSERT INTO backups (status, sha512, size, error, completed_at) VALUES ($1, $2, $3, $4, $5) RETURNING backup_id, created_at, updated_at`
	backupUpdate = `
UPDATE backups SET status = $2, sha512 = $3, size = $4, error = $5, completed_at = $6, updated_at = now() WHERE backup_id = $1 RETURNING updated_at`
	backupSelectLatest = `
SELECT backup_id, status, sha512, size, error, created_at, updated_at, completed_at FROM backups WHERE deleted_at IS NULL ORDER BY updated_at DESC LIMIT 1`
	sinkListQuery = `
SELECT sink_id, kind, config, created_at, updated_at FROM sinks WHERE deleted_at IS NULL ORDER BY updated_at DESC`
	sinkListSinceQuery = `
SELECT sink_id, kind, config, created_at, updated_at FROM sinks WHERE updated_at >= $1 AND deleted_at IS NULL ORDER BY updated_at DESC`
	sinkSelectQuery = `
SELECT sink_id, kind, config, created_at, updated_at FROM sinks WHERE sink_id = $1`
	sinkInsertQuery = `
INSERT INTO sinks (sink_id, kind, config) VALUES ($1, $2, $3) RETURNING created_at, updated_at`
	sinkDeleteQuery = `
UPDATE sinks SET deleted_at = now() WHERE sink_id = $1 AND deleted_at IS NULL`
	volumeListQuery = `
SELECT volume_id, host_id, type, state, app_id, release_id, job_id, job_type, path, delete_on_stop, meta, created_at, updated_at, decommissioned_at FROM volumes ORDER BY updated_at DESC`
	volumeAppListQuery = `
SELECT volume_id, host_id, type, state, app_id, release_id, job_id, job_type, path, delete_on_stop, meta, created_at, updated_at, decommissioned_at FROM volumes WHERE app_id = $1 ORDER BY updated_at DESC`
	volumeListSinceQuery = `
SELECT volume_id, host_id, type, state, app_id, release_id, job_id, job_type, path, delete_on_stop, meta, created_at, updated_at, decommissioned_at FROM volumes WHERE updated_at >= $1 ORDER BY updated_at DESC`
	volumeSelectQuery = `
SELECT volume_id, host_id, type, state, app_id, release_id, job_id, job_type, path, delete_on_stop, meta, created_at, updated_at, decommissioned_at FROM volumes WHERE app_id = $1 AND volume_id = $2`
	volumeInsertQuery = `
INSERT INTO volumes (volume_id, host_id, type, state, app_id, release_id, job_id, job_type, path, delete_on_stop, meta) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (volume_id) DO UPDATE SET job_id = $7, updated_at = now()
RETURNING created_at, updated_at`
	volumeDecommissionQuery = `
UPDATE volumes SET updated_at = now(), decommissioned_at = now() WHERE app_id = $1 AND volume_id = $2 RETURNING updated_at, decommissioned_at`
)
