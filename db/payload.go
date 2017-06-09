package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ViewRecord struct {
	gorm.Model
	FileName   string
	RecordTime time.Time
}
