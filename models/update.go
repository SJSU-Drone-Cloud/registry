package models

import (
	"fmt"
	"strconv"
)

type Update struct {
	id int64
}

func NewUpdate(userID int64, body string) (*Update, error) {
	id, err := client.Incr("update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline() //optimization, no need to wait for a return, sendc commands in batch basically
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userID)
	pipe.HSet(key, "body", body)
	userUpdatesKey := fmt.Sprintf("user:%d:updates", userID)
	pipe.LPush(userUpdatesKey, id) //pushes it to a specific user key of user:Name:updates
	pipe.LPush("updates", id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	u := Update{id}
	return &u, err
}

func (update *Update) GetBody() (string, error) {
	key := fmt.Sprintf("update:%d", update.id)
	return client.HGet(key, "body").Result()
}

func (update *Update) GetUser() (*User, error) {
	key := fmt.Sprintf("update:%d", update.id)
	userId, err := client.HGet(key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserById(userId)

}

func GetAllUpdates() ([]*Update, error) {
	return queryUpdates("updates")
}

func GetUpdates(userID int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userID)
	fmt.Println("getting key", key)
	return queryUpdates(key)
}

func queryUpdates(key string) ([]*Update, error) {
	updateIds, err := client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIds))
	for i, strID := range updateIds {
		id, err := strconv.Atoi((strID))
		if err != nil {
			return nil, err
		}
		updates[i] = &Update{int64(id)}
	}
	return updates, nil
}

func PostUpdates(userID int64, body string) error {
	_, err := NewUpdate(userID, body)
	return err
}
