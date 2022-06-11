package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/ubgo/gofm/publicid"
	"gorm.io/gorm"
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

type ModelNid struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey;"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *ModelNid) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.ID) == 0 {
		u.ID = publicid.Must()
	}
	return
}

type ModelStr struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey;"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *ModelStr) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.ID) == 0 {
		u.ID = publicid.Must()
	}
	return
}

type Plugin interface {
	Version() string
	MigrateDb()
}
