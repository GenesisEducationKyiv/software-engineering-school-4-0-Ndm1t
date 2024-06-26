package container

import (
	"gorm.io/gorm"
)

type IContainer interface {
	GetDatabase() *gorm.DB
}

type Container struct {
	db *gorm.DB
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		db: db,
	}
}

func (c *Container) GetDatabase() *gorm.DB {
	return c.db
}
