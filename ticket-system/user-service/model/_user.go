package model

import (
	"encoding/json"
	"log"

	"github.com/POABOB/go-microservice/ticket-system/pkg/mysql"
	"github.com/gohouse/gorose/v2"
)

type User struct {
	Id int64 `json:"id"`
	// EDIT IT!
}

type UserModel struct{}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (p *UserModel) getTableName() string {
	return "user"
}

func (p *UserModel) GetUserList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

func (p *UserModel) GetUserByPK(pk int64) (*User, error) {
	conn := mysql.DB()
	if result, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"id": pk}).First(); err == nil {
		// Convert map to JSON
		jsonData, _ := json.Marshal(result)

		// Convert the JSON to a struct
		var entity User
		json.Unmarshal(jsonData, &entity)
		return &entity, err
	} else {
		return nil, err
	}
}

func (p *UserModel) CreateUser(entity *User) error {
	conn := mysql.DB()
	conn.Begin()
	if _, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"id": entity.Id,
		// ...
	}).Insert(); err != nil {
		conn.Rollback()
		log.Printf("Error : %v", err)
		return err
	}

	conn.Commit()
	return nil
}

func (p *UserModel) UpdateUser(entity *User) error {
	conn := mysql.DB()
	conn.Begin()
	if _, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"id": entity.Id,
		// ...
	}).Where("id", entity.Id).Update(); err != nil {
		conn.Rollback()
		log.Printf("Error : %v", err)
		return err
	}

	conn.Commit()
	return nil
}

func (p *UserModel) DeleteUser(entity *User) error {
	conn := mysql.DB()
	conn.Begin()
	if _, err := conn.Table(p.getTableName()).Where("id", entity.Id).Delete(); err != nil {
		conn.Rollback()
		log.Printf("Error : %v", err)
		return err
	}

	conn.Commit()
	return nil
}
