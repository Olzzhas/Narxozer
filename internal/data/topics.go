package data

import (
	"database/sql"
	"github.com/olzzhas/narxozer/graph/model"
)

type TopicModel struct {
	DB *sql.DB
}

func (m TopicModel) Insert(topic *model.Topic) (*model.Topic, error) {

	return nil, nil
}

func (m TopicModel) GetAll() ([]*model.Topic, error) {
	return nil, nil
}

func (m TopicModel) GetByID(id int) (*model.Topic, error) {
	return nil, nil
}

func (m TopicModel) Update(topic *model.Topic) (*model.Topic, error) {
	return nil, nil
}

func (m TopicModel) Delete(id int) error {
	return nil
}
