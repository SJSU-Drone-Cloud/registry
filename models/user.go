package models

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	InvalidLogin  = errors.New("Invalid Login")
	UserNameTaken = errors.New("Username taken.")
)

type User struct {
	id int64
}

func (user *User) GetId() (int64, error) {
	return user.id, nil
}

func (user *User) GetUserName() (string, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return client.HGet(key, "username").Result()
}

func (user *User) GetHash() ([]byte, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return client.HGet(key, "hash").Bytes()
}

func (user *User) Authenticate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return InvalidLogin
	} else if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(username string) (*User, error) {
	id, err := client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, InvalidLogin
	} else if err != nil {
		return nil, err
	}
	return GetUserById(id)
}

func NewUser(user string, hash []byte) (*User, error) {
	exists, err := client.HExists("user:by-username", user).Result()
	if exists {
		return nil, UserNameTaken
	} else if err != nil {
		return nil, err
	}

	id, err := client.Incr("user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline() //optimization, no need to wait for a return, sendc commands in batch basically
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", user)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", user, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	u := User{id}
	return &u, err
}

func RegisterUser(user string, pass string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	if err != nil {
		return err
	}
	_, err = NewUser(user, hash)

	return err
}

func AuthenticateUser(username string, pass string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, user.Authenticate(pass)
}

func GetUserById(id int64) (*User, error) {
	return &User{id}, nil
}
