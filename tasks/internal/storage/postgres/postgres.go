package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tasks/internal/config"
)

type Storage struct {
	Conn *pgxpool.Pool
}

type Task struct {
	Description string
	Title       string
	IsDone      bool
	ID          int
}

func NewConnection(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		// log.Fatalf завершает программу, лучше возвращать ошибку
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return &Storage{Conn: dbpool}, nil
}

func (s *Storage) GetTasks() ([]Task, error) {
	rows, err := s.Conn.Query(context.Background(), "SELECT description, title, is_done, id FROM task")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Description, &task.Title, &task.IsDone, &task.ID); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) CreateTask(title, description string) error {
	query := "INSERT INTO task (description, title, is_done) VALUES ($1, $2, $3)"
	_, err := s.Conn.Exec(context.Background(), query, description, title, false)
	return err
}

func (s *Storage) DeleteTask(taskID string) error {
	query := "DELETE FROM task WHERE id = $1"
	_, err := s.Conn.Exec(context.Background(), query, taskID)
	return err
}

func (s *Storage) GetTaskByID(taskID string) (Task, error) {
	var task Task
	query := "SELECT description, title, is_done, id FROM task WHERE id = $1"
	err := s.Conn.QueryRow(context.Background(), query, taskID).Scan(
		&task.Description, &task.Title, &task.IsDone, &task.ID,
	)
	return task, err
}

func (s *Storage) UpdateTask(taskID, title, description string, isDone bool) error {
	query := "UPDATE task SET title = $1, description = $2, is_done = $3 WHERE id = $4"
	_, err := s.Conn.Exec(context.Background(), query, title, description, isDone, taskID)
	return err
}
