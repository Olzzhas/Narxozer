package data

import (
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/olzzhas/narxozer/graph/model"
)

type TopicModel struct {
	DB    *sql.DB
	Redis *redis.Client
}

// Insert добавляет новый топик в базу данных и возвращает его
func (m TopicModel) Insert(topic *model.Topic) (*model.Topic, error) {
	query := `
		INSERT INTO topics (title, content, image_url, author_id, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at`

	err := m.DB.QueryRow(query, topic.Title, topic.Content, topic.ImageURL, topic.Author.ID).Scan(&topic.ID, &topic.CreatedAt)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// GetAll возвращает все топики из базы данных
func (m TopicModel) GetAll() ([]*model.Topic, error) {
	query := `SELECT id, title, content, image_url, author_id, created_at, updated_at FROM topics`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*model.Topic
	for rows.Next() {
		var topic model.Topic
		topic.Author = &model.User{}
		err := rows.Scan(&topic.ID, &topic.Title, &topic.Content, &topic.ImageURL, &topic.Author.ID, &topic.CreatedAt, &topic.UpdatedAt)
		if err != nil {
			return nil, err
		}
		topics = append(topics, &topic)
	}

	return topics, nil
}

// GetByID возвращает топик по его ID
func (m TopicModel) GetByID(id int) (*model.Topic, error) {
	query := `
		SELECT id, title, content, image_url, author_id, created_at, updated_at, likes
		FROM topics 
		WHERE id = $1`
	topic := &model.Topic{}
	topic.Author = &model.User{}
	err := m.DB.QueryRow(query, id).Scan(
		&topic.ID,
		&topic.Title,
		&topic.Content,
		&topic.ImageURL,
		&topic.Author.ID,
		&topic.CreatedAt,
		&topic.UpdatedAt,
		&topic.Likes, // Добавлено поле likes
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// Update обновляет данные топика и возвращает обновленный топик
func (m TopicModel) Update(topic *model.Topic) (*model.Topic, error) {
	query := `
		UPDATE topics
		SET title = $1, content = $2, image_url = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, content, image_url, author_id, created_at, updated_at`

	err := m.DB.QueryRow(query, topic.Title, topic.Content, topic.ImageURL, topic.ID).Scan(
		&topic.ID, &topic.Title, &topic.Content, &topic.ImageURL, &topic.Author.ID, &topic.CreatedAt, &topic.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// Delete удаляет топик по его ID
func (m TopicModel) Delete(id int) error {
	query := `DELETE FROM topics WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
