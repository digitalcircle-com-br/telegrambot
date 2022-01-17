package lib

import (
	"github.com/digitalcircle-com-br/envfile"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GDB struct {
	*gorm.DB
}

var db *GDB

func (d *GDB) FindOrCreate(c interface{}) error {
	err := d.Where(c).First(c).Error
	if err == gorm.ErrRecordNotFound {
		return d.Save(c).Error
	}

	return err

}

func (d *GDB) Init() error {
	var err error
	dsn := envfile.Must("DSN")
	d.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = d.Exec("select 1+1").Error
	if err != nil {
		return err
	}

	return nil
}

func (d *GDB) AutoMigrates(i ...interface{}) error {
	for _, v := range i {
		err := db.AutoMigrate(v)
		if err != nil {
			return err
		}
	}
	return nil
}
