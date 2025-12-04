package repository

const (
	createEventQuery = `
		INSERT INTO events (id,
		                    name,
		                    date,
		                    total_seats,
		                    booking_lifetime,
		                    requires_payment_confirmation,
		                    created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	getEventByIDQuery = `
	SELECT id,
	       name,
	       date,
	       total_seats,
	       reserved_seats,
	       booked_seats,
	       booking_lifetime,
	       requires_payment_confirmation,
	       created_at
	FROM events
	WHERE id = $1
`
	createUserQuery = `
	INSERT INTO users (id,
	                   name,
	                   email,
	                   telegram_id,
	                   created_at)
	VALUES ($1, $2, $3, $4, $5)
`
	getUserByIDQuery = `
	SELECT id, 
	       name,
	       email,
	       telegram_id,
	       created_at
	FROM users
	WHERE id = $1
`
	getUserByEmailQuery = `
	SELECT id,
	       name,
	       email,
	       telegram_id,
	       created_at
	FROM users
	WHERE email = $1
`
	selectEventForUpdateQuery = `
	SELECT id,
	       name,
	       date,
	       total_seats,
	       reserved_seats,
	       booked_seats,
	       booking_lifetime,
	       requires_payment_confirmation,
	       created_at
	FROM events
	WHERE id = $1
	FOR UPDATE
`

	insertBookingQuery = `
	INSERT INTO bookings (id,
	                      event_id,
	                      user_id,
	                      status,
	                      deadline,
	                      created_at,
	                      updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	updateReservedSeatsQuery = `
	UPDATE events 
	SET reserved_seats = reserved_seats + 1 
	WHERE id = $1
`
	updateBookedSeatsQuery = `
	UPDATE events 
	SET booked_seats = booked_seats + 1 
	WHERE id = $1
`
	selectBookingForUpdateQuery = `
	SELECT id,
	       event_id,
	       user_id,
	       status,
	       deadline,
	       created_at,
	       updated_at
	FROM bookings
	WHERE id = $1
	FOR UPDATE
`
	updateBookingStatusQuery = `
	UPDATE bookings 
	SET status = $1, updated_at = NOW() 
	WHERE id = $2
`
	updateEventSeatsQuery = `
	UPDATE events
	SET reserved_seats = reserved_seats - 1,
	    booked_seats = booked_seats + 1
	WHERE id = $1
`
	countUserBookingsQuery = `
	SELECT 1
	FROM bookings
	WHERE event_id = $1 
	  AND user_id = $2 
	  AND status IN ('reserved', 'confirmed')
	LIMIT 1
`

	decreaseBookingSeatsQuery = `
	UPDATE events
	SET reserved_seats = reserved_seats - 1
	WHERE id = $1
`
	listEventsQuery = `
	SELECT id,
		   name,
		   date,
		   total_seats,
		   reserved_seats,
		   booked_seats,
		   booking_lifetime,
		   requires_payment_confirmation,
		   created_at
	FROM events
	ORDER BY date
	`

	getBookingsByEventQuery = `
	SELECT id,
	       event_id,
	       user_id,
	       status,
	       deadline,
	       created_at,
	       updated_at
	FROM bookings
	WHERE event_id = $1
	ORDER BY created_at DESC
`

	getBookingByIDQuery = `
	SELECT id,
	       event_id,
	       user_id,
	       status,
	       deadline,
	       created_at,
	       updated_at
	FROM bookings
	WHERE id = $1
`

	getExpiredReservedBookingsQuery = `
	SELECT id,
	       event_id,
	       user_id,
	       status,
	       deadline,
	       created_at,
	       updated_at
	FROM bookings
	WHERE status = 'reserved' AND deadline <= NOW()
	ORDER BY deadline 
`
)
