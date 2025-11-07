package sponsorRepository

const (
	qGetSponsor     = `SELECT * FROM SPONSORS WHERE sponsor_id = :sponsorId`
	qGetSponsors    = `SELECT * FROM SPONSORS`
	qInsertSponsor  = `INSERT INTO SPONSORS (sponsor_id, name, logo, website, description, level) VALUES (:sponsorId, :name, :logo, :website, :description, :level)`
	qDeleteSponsor  = `DELETE FROM SPONSORS WHERE sponsor_id = :sponsorId`
)