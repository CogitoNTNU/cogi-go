package projectRepository

const (
	qGetProject     = `SELECT * FROM PROJECTS WHERE project_id = :projectId`
	qGetProjects    = `SELECT * FROM PROJECTS`
	qInsertProject  = `INSERT INTO PROJECTS (project_id, title, github_url, logo, playable, released, active_project, project_url) VALUES (:projectId, :title, :githubUrl, :logo, :playable, :released, :activeProject, :projectUrl)`
	qDeleteProject  = `DELETE FROM PROJECTS WHERE project_id = :projectId`
)