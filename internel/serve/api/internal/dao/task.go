package dao

import (
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/types"
	"fmt"
	"gorm.io/gorm"
)

var TaskDao TaskModel

type TaskModel struct {
}

func (t *TaskModel) Create(d model.Task) error {
	return db.Model(&model.Task{}).Create(&d).Error
}

func (t *TaskModel) Delete(id int) error {
	return db.Model(&model.Task{}).Delete("id", id).Error
}

func (t *TaskModel) Update(d model.Task) error {
	return db.Model(&model.Task{}).Where("id", d.ID).Updates(d).Error
}

func (t *TaskModel) UpdateStatus(id int, status int) error {
	return db.Model(&model.Task{}).Where("id", id).Update("status", status).Error
}

func (t *TaskModel) Find(ids []int) ([]model.Task, error) {
	data := make([]model.Task, 0)
	_db := db.Model(&model.Task{})
	if len(ids) != 0 {
		_db = _db.Where("id IN ?", ids)
	}
	err := _db.Where("status != ?", model.StatusSuccess).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TaskModel) Exist(data string) (*model.Task, error) {
	findTask := &model.Task{}
	if err := db.Model(&model.Task{}).Where("data = ?", data).First(findTask).Error; err != nil {
		return nil, err
	}

	return findTask, nil
}

func (t *TaskModel) parseWhere(where map[string]any) *gorm.DB {
	_db := db.Model(&model.Task{})
	for key, value := range where {
		if value == nil {
			continue
		}
		switch key {
		case "type":
			if value == "all" {
				continue
			}
		case "video_type":

		}
		sql := fmt.Sprintf("%s = ?", key)
		_db = _db.Where(sql, value)
	}

	return _db
}

func (t *TaskModel) Count(turner types.PageTurner) (int64, error) {
	where, _, _ := turner.ParseMysql()

	_db := t.parseWhere(where)
	var count int64
	if err := _db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (t *TaskModel) List(turner types.PageTurner) ([]model.Task, error) {
	where, offset, limit := turner.ParseMysql()

	_db := t.parseWhere(where)
	var data []model.Task
	err := _db.Offset(offset).Limit(limit).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
