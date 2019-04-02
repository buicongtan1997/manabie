package mappers

import "github.com/jinzhu/gorm"

type baseMapper struct {
	db *gorm.DB
}

/**
Insert new record to db
 */
func (rcv *baseMapper) Insert(value interface{}) (err error) {
	err = rcv.db.Create(value).Error
	return
}

func (rcv *baseMapper) Save(value interface{}) (err error) {
	err = rcv.db.Save(value).Error
	return
}


/**
Delete record from db
 */
func (rcv *baseMapper) Delete(value interface{}) (err error) {
	err = rcv.db.Unscoped().Delete(value).Error
	return
}

/**
Set column delete_at = now
 */
func (rcv *baseMapper) SoftDelete(value interface{}) (err error) {
	err = rcv.db.Delete(value).Error
	return
}

func (rcv *baseMapper) FindAll(value interface{}) (err error) {
	err = rcv.db.Find(value).Error
	return
}
