package cron

import (
	"health-checker/env"
	"log"
	"time"
)

type Service struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string    `gorm:"not null"`
	URL           string    `gorm:"not null"`
	LastStatus    string    `gorm:"null"`
	LastCheckedAt time.Time `gorm:"null"`
	UpdatedAt     time.Time `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}

func GetAll(env *env.Env) ([]Service, error) {
	log.Println("Get all services...")

	var records []Service
	res := env.Db.Find(&records)

	if res.Error != nil {
		return nil, res.Error
	}

	return records, nil
}

func Create(env *env.Env, payload Service) error {
	res := env.Db.Create(payload)

	return res.Error
}

func UpdateById(env *env.Env, id string, payload map[string]interface{}) error {
	res := env.Db.Model(&Service{}).Where("id = ?", id).Updates(payload)

	return res.Error
}
