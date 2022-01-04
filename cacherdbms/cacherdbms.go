package cacherdbms

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ubgo/gofm/gormdb"
	"github.com/ubgo/goutil"
	"gorm.io/gorm"
)

// Cache ...
type Cache struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	Key       string `gorm:"type:varchar(100);unique_index"`
	Value     string
	Expires   int
}

func (Cache) TableName() string {
	return getTableName()
}

func getTableName() string {
	return goutil.Env("PKG_CACHERDBMS_TABLENAME", "caches")
}

type Rdbms struct {
	Config
}

func (pkg Rdbms) Version() string {
	return "0.01"
}

type Config struct {
	GormDB gormdb.GormDB
}

func (pkg Rdbms) MigrateDb() {
	pkg.GormDB.DB.AutoMigrate(&Cache{})
}

// New initialize
func New(config Config) *Rdbms {
	rdbms := &Rdbms{Config: config}
	return rdbms
}

func (a *Rdbms) Put(key string, val interface{}, ttl int) (bool, error) {
	p, err := json.Marshal(val)
	if err != nil {
		return false, err
	}

	expire := int(time.Now().UnixNano()/int64(time.Second) + int64(ttl))

	var entity Cache
	res := a.Config.GormDB.DB.First(&entity, &Cache{Key: key})
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		entity := &Cache{
			Key:     key,
			Value:   string(p),
			Expires: expire,
		}

		err1 := a.Config.GormDB.DB.Create(entity).Error
		if err1 != nil {
			return false, err
		}
	} else {
		entity.Value = string(p)
		entity.Expires = expire
		a.Config.GormDB.DB.Save(&entity)
	}

	return true, nil
}

func (a *Rdbms) Get(key string) interface{} {
	var entity Cache
	res := a.Config.GormDB.DB.First(&entity, &Cache{Key: key})

	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	now := int(time.Now().UnixNano() / int64(time.Second))
	if now > entity.Expires {
		a.Del(key)
		return nil
	}

	var v string
	err := json.Unmarshal([]byte(entity.Value), &v)
	if err != nil {
		return entity.Value
	}
	// val, _ := strconv.Unquote(entity.Value)

	return v
}

func (a *Rdbms) Del(key string) {
	a.Config.GormDB.DB.Where("key = ?", key).Delete(&Cache{})
}

func (a *Rdbms) Flush() {
	a.Config.GormDB.DB.Exec("Truncate TABLE " + getTableName() + " RESTART IDENTITY;")
}
