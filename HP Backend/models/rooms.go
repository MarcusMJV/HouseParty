package models

import (
	"time"

	"github.com/google/uuid"
	"houseparty.com/storage"
)

type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	HostID      int64     `json:"host_id"`
	Public      bool      `json:"public"`
	CreatedAt   time.Time `json:"created_at"`
}


type RoomResponse struct {
    Room
	HostName string `json:"host_name"`
}

func (r *Room) ToRoomResponse(name string) *RoomResponse {
    return &RoomResponse{
        Room: *r,
		HostName: name,
    }
}

func (r *Room) Delete() error {
	stmt, err := storage.DB.Prepare(storage.DeleteRoomQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Room) Save() error {
	unique_id := uuid.New().String()

	stmt, err := storage.DB.Prepare(storage.SaveRoomQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(unique_id, r.Name, r.Description, r.HostID, r.Public, r.CreatedAt)
	if err != nil {
		return err
	}
	r.ID = unique_id

	return nil
}

func (r *Room) GetRoomById(id string) error {
	stmt, err := storage.DB.Prepare(storage.GetRoomByIdQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&r.ID, &r.Name, &r.Description, &r.HostID, &r.Public, &r.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
