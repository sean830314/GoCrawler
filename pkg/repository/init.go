package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	model "github.com/sean830314/GoCrawler/pkg/model/admin"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Role{}, &model.User{})

	return db.Debug(), nil
}

type repository struct {
	Role model.RoleRepository
	User model.UserRepository
}

func New(db *gorm.DB) repository {
	return repository{
		Role: NewRoleRepository(db),
		User: NewUserRepository(db),
	}
}
