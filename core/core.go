package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Core struct {
}

func (pkg Core) Version() string {
	return "0.0.1"
}

func New() Core {
	getAndSetConfig()
	return Core{}
}

func getAndSetConfig() {
	viper.SetConfigName("default") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

type Model struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Plugin interface {
	Version() string
	MigrateDb()
}
