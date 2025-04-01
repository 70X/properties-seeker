package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type DB struct {
	*gorm.DB
	config *Config
}

func (*DB) getEnv(key, defaultValue string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	if defaultValue == "" {
		return "", fmt.Errorf("env %s is empty and no default value is set", key)
	}
	return defaultValue, nil
}
func (s *DB) setConfigFromEnv() error {
	port, _ := s.getEnv("DB_PORT", "5432")
	host, _ := s.getEnv("DB_HOST", "localhost")
	user, err := s.getEnv("DB_USER", "")
	if err != nil {
		return fmt.Errorf("env %s is not set: %v", "DB_USER", err)
	}
	pwd, err := s.getEnv("DB_PASSWORD", "")
	if err != nil {
		return fmt.Errorf("env %s is not set: %v", "DB_PASSWORD", err)
	}
	dbname, err := s.getEnv("DB_NAME", "")
	if err != nil {
		return fmt.Errorf("env %s is not set: %v", "DB_NAME", err)
	}
	config := &Config{
		Host: host, Port: port, User: user, Password: pwd, DBName: dbname,
	}
	log.Printf("config: %v", config)
	s.config = config
	return nil
}

func (s *DB) setInstance() error {
	err := s.initialize(s.config)
	if err != nil {
		return fmt.Errorf("failed to take database conn: %v", err)
	}
	return nil
}

func (s *DB) initialize(config *Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	s.DB = db
	return nil
}

func NewDb(config *Config) *DB {
	s := &DB{config: config}
	if config == nil {
		err := s.setConfigFromEnv()
		if err != nil {
			panic(fmt.Sprintf("DB config is missing %v", err))
		}
	}
	s.setInstance()
	return s
}
