package eventRepository

const (
	qGetEvents = `SELECT event_id, name, start_at, end_at, type, location, description, max_atendees
        FROM EVENT;
    `
	qGetEvent = `SELECT * FROM EVENT WHERE event_id = :eventId`

	qCreateEvent = `INSERT INTO EVENT (event_id, name, start_at, end_at, created_at, updated_at, type, location, description, content, max_atendees) VALUES
        (
            :eventId,
            :name,
            :startAt,
						:endAt,
            :createdAt,
            :updatedAt,
            :type,
            :location,
            :description,
            :content,
            :maxAtendees
        )
    
     `
	qDeleteEvent = `DELETE FROM EVENT WHERE event_id = :eventId`
)
