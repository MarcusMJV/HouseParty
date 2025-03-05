package storage

const RetrieveHashPasswordQuery = "SELECT id, username, email, password FROM users WHERE email = ? OR username = ? LIMIT 1"

const SaveUserQuery = `INSERT INTO users(email, password, username) VALUES(?, ?, ?)`

const SaveRoomQuery = `INSERT INTO rooms(id, name, description, host_id, public, created_at) VALUES(?, ?, ?, ?, ?, ?)`

const GetRoomByIdQuery = `SELECT id, name, description, host_id, public, created_at FROM rooms WHERE id = ?`

const GetUserByIdQuery = `SELECT id, username, email, spotify_connected FROM users WHERE id = ?`

const DeleteRoomQuery = `DELETE FROM rooms WHERE id = ?`

const GetSpotifyToken = `SELECT access_token, token_type, scope, expires_in, refresh_token, time_issued, user_id FROM token WHERE user_id = ?`

const DeleteTokenQuery = "DELETE FROM token WHERE user_id = ?"

const ActivateSpotifyQuery = "UPDATE users SET spotify_connected = ? WHERE id = ?"

const SaveTokenQuery = `
INSERT INTO token(access_token, token_type, scope, expires_in, refresh_token, time_issued, user_id) 
VALUES(?, ?, ?, ?, ?, ?, ?)`

const UpdateTokenQuery = `
UPDATE token 
SET access_token = ?, 
    token_type = ?, 
    scope = ?, 
    expires_in = ?, 
    time_issued = ? 
WHERE user_id = ?`

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