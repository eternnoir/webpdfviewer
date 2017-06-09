package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

type DbMgr struct {
	dbpath string
}

func NewDbMgr(path string) (*DbMgr, error) {
	db := &DbMgr{dbpath: path}
	log.Warnf("Init database %s", path)
	if err := db.CheckDbExist(); err != nil {
		return nil, err
	}
	return db, nil
}

func (dm *DbMgr) CheckDbExist() error {
	db, err := gorm.Open("sqlite3", dm.dbpath)
	if err != nil {
		return err
	}
	defer db.Close()
	db.AutoMigrate(&ViewRecord{})
	return nil
}

func (dm *DbMgr) InsertRecord(filename string) {
	log.Infof("Start Insert record. %s", filename)
	db, err := gorm.Open("sqlite3", dm.dbpath)
	if err != nil {
		log.Errorf("Cannot open database. %s", dm.dbpath)
		return
	}
	defer db.Close()
	vr := &ViewRecord{FileName: filename, RecordTime: time.Now()}
	db.NewRecord(vr) // => returns `true` as primary key is blank

	err = db.Create(vr).Error
	if err != nil {
		log.Errorf("Insert %s Record Error.%s", filename, err.Error())
		return
	}
	log.Infof("Insert Record %s success.", filename)
}
