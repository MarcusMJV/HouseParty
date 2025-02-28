package services

import (
	"database/sql"
	"sync"
	"time"

	"houseparty.com/models"
	"houseparty.com/storage"
)

func CreateRoom(room *models.Room, userId int64) error {
	room.HostID = userId
	room.CreatedAt = time.Now()
	err := room.Save()
	return err
}

func DeleteRoomByID(roomId string) (*models.RoomResponse, error){
	var room models.Room
	var host models.User

	err := room.GetRoomById(roomId)
	if err != nil {
		return nil, err
	}

	err = host.GetUserById(room.HostID)
	if err != nil {
		return nil, err
	}

	room.Delete()

	return room.ToRoomResponse(host.Username), nil
}

func GetRooms(userId int64) ([]*models.RoomResponse, *models.RoomResponse, error) {
	publicRoomsChan := make(chan []*models.RoomResponse, 1)
	userRoomsChan := make(chan *models.RoomResponse, 1)
	errorChan := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rooms, err := FilterAndGetRooms(userId)
		if err != nil {
			errorChan <- err
			publicRoomsChan <- nil
			return
		}
		publicRoomsChan <- rooms
	}()

	go func() {
		defer wg.Done()
		rooms, err := GetUserRoom(userId)
		if err != nil {
			errorChan <- err
			userRoomsChan <- nil
			return
		}
		userRoomsChan <- rooms
	}()

	wg.Wait()
	close(publicRoomsChan)
	close(userRoomsChan)
	close(errorChan)

	if len(errorChan) > 0 {
		return nil, nil, <-errorChan
	}

	publicRooms := <-publicRoomsChan
	userRooms := <-userRoomsChan

	return publicRooms, userRooms, nil
}

func  GetUserRoom(userId int64) (*models.RoomResponse, error) {
	var room models.RoomResponse

	stmt, err := storage.DB.Prepare(storage.UserRoomsQuery)
	if err != nil {
		return &room, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(userId)

	err = row.Scan(&room.ID, 
		&room.Name, 
		&room.Description, 
		&room.HostID, 
		&room.Public, 
		&room.CreatedAt,
		&room.HostName)
	
	if err == sql.ErrNoRows{
		return &room, nil
	}else if err != nil {
		return &room, err
	}

	return &room, nil
}

func FilterAndGetRooms(userId int64) ([]*models.RoomResponse, error) {
	var rooms []*models.RoomResponse
	rows, err := storage.DB.Query(storage.PublicRoomsQuery, true, userId)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		var username string

		err := rows.Scan(
			&room.ID,
			&room.Name, 
			&room.Description, 
			&room.HostID, 
			&room.Public, 
			&room.CreatedAt, 
			&username,
		)
		
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room.ToRoomResponse(username))
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}
