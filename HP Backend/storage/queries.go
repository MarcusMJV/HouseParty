package storage

const RetrieveHashPasswordQuery = "SELECT id, username, email, password FROM users WHERE email = ? OR username = ? LIMIT 1"

const SaveUserQuery = `INSERT INTO users(email, password, username) VALUES(?, ?, ?)`

const SaveRoomQuery = `INSERT INTO rooms(id, name, description, host_id, public, created_at) VALUES(?, ?, ?, ?, ?, ?)`

const GetRoomByIdQuery = `SELECT id, name, description, host_id, public, created_at FROM rooms WHERE id = ?`

const GetUserByIdQuery = `SELECT id, username, email FROM users WHERE id = ?`

const DeleteRoomQuery = `DELETE FROM rooms WHERE id = ?`

const PublicRoomsQuery = `
SELECT 
    rooms.id, 
    rooms.name, 
    rooms.description, 
    rooms.host_id, 
    rooms.public, 
    rooms.created_at, 
    users.username
FROM 
    rooms 
JOIN 
    users 
ON 
    rooms.host_id = users.id 
WHERE 
    rooms.public = ? 
AND 
    rooms.host_id != ?`

const UserRoomsQuery = `
SELECT 
    rooms.id, 
    rooms.name, 
    rooms.description, 
    rooms.host_id, 
    rooms.public, 
    rooms.created_at, 
    users.username
FROM 
    rooms 
JOIN 
    users 
ON 
    rooms.host_id = users.id 
WHERE 
    rooms.host_id = ?`

const CheckInUseQuery = `
	SELECT 
		CASE 
			WHEN username = ? THEN 'username'
			WHEN email = ? THEN 'email'
		END AS conflict_field
	FROM users
	WHERE username = ? OR email = ?
	LIMIT 1;
`