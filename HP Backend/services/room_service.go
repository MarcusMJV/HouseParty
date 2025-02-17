package services

import (
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

func GetRooms(userId int64) ([]*models.RoomResponse, []*models.RoomResponse, error) {
	publicRoomsChan := make(chan []*models.RoomResponse, 1)
	userRoomsChan := make(chan []*models.RoomResponse, 1)
	errorChan := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rooms, err := FilterAndGetRooms(storage.PublicRoomsQuery, true, userId)
		if err != nil {
			errorChan <- err
			publicRoomsChan <- nil
			return
		}
		publicRoomsChan <- rooms
	}()

	go func() {
		defer wg.Done()
		rooms, err := FilterAndGetRooms(storage.UserRoomsQuery, userId)
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

func FilterAndGetRooms(query string, values ...interface{}) ([]*models.RoomResponse, error) {
	var rooms []*models.RoomResponse
	rows, err := storage.DB.Query(query, values...)
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
