package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/oziev02/checklist-microservices/internal/db/domain/entities"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(host, port, user, password, dbname string) (*PostgresRepository, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping PostgerSQL: %v", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("Failed to create tables: %v", err)
	}

	return &PostgresRepository{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Таблица пользователей
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            avatar TEXT,
            description TEXT,
            twofa_enabled BOOLEAN DEFAULT FALSE,
            twofa_secret TEXT
        )
    `)
	if err != nil {
		return err
	}

	// Таблица соцсетей (для хранения socials пользователя)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_socials (
		    user_id UUID PREFERENCES users(id),
		    social_key TEXT,
		    social_value TEXT,
		    PRIMARY KEY (user_id, social_key)
		)
	`)
	if err != nil {
		return nil
	}

	// Таблица задач
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id UUID PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT,
            done BOOLEAN DEFAULT FALSE,
            user_id UUID REFERENCES users(id)
        )
    `)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email, password string) (*entities.User, error) {
	user := entities.NewUser(email, password)
	query := `
        INSERT INTO users (id, email, password, avatar, description, twofa_enabled, twofa_secret)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Avatar, user.Description, user.TwoFAEnabled, user.TwoFASecret)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %v", err)
	}
	return user, nil
}

func (r *PostgresRepository) UpdateProfile(ctx context.Context, userID, avatar, description string, socials map[string]string) (*entities.User, error) {
	// Обновляем основные поля профиля
	query := `
		UPDATE users
		SET avatar = $2, description = $3
		WHERE id = $1
		RETURING id, email, password, avatar, twofa_enabled, twofa_secret
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, userID, avatar, description).Scan(
		&userID, &user.Email, &user.Password, &user.Avatar, &user.Description, &user.TwoFAEnabled, &user.TwoFASecret
		)
	if err != nil {
		return nil, fmt.Errorf("Failed to update profile: %v", err)
	}

	// Обновляем соцсети
	// Сначала удаляем старые записи
	_, err = r.db.ExecContext(ctx, "DELETE FROM user_socials WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete old socials: %v", err)
	}

	// Добавляем новые записи
	for key, value := range socials {
		_, err = r.db.ExecContext(ctx, "INSERT INTO user_socials (user_id, social_key, social_value) VALUES ($1, $2, $3)", userID, key, value)
		if err != nil {
			return nil, fmt.Errorf("Failed to insert social: %v", err)
		}
	}

	user.Socials = socials
	return user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	query := `
        SELECT id, email, password, avatar, description, twofa_enabled, twofa_secret
        FROM users
        WHERE email = $1
    `
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Avatar, &user.Description, &user.TwoFAEnabled, &user.TwoFASecret,
		)
	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get user by email: %v", err)
	}

	// Получаем соцсети
	user.Socials = make(map[string]string)
	rows, err := r.db.QueryContext(ctx, "SELECT social_key, social_value FROM user_socials WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get socials: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("Failed to scan socials: %v", err)
		}
		user.Socials[key] = value
	}

	return user, nil
}

func (r *PostgresRepository) CreateTask(ctx context.Context, title, content, userID string) (*entities.Task, error) {
	task := entities.NewTask(title, content, uuid.MustParse(userID))
	query := `
        INSERT INTO tasks (id, title, content, done, user_id)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecContext(ctx, query, task.ID, task.Title, task.Content, task.Done, task.UserID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create task: %v", err)
	}
	return task, nil
}

func (r *PostgresRepository) ListTasks(ctx context.Context, userID string) ([]*entities.Task, error) {
	query := `
        SELECT id, title, content, done, user_id
        FROM tasks
        WHERE user_id = $1
    `
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to list tasks: %v", err)
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		task := &entities.Task{}
		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.Done, &task.UserID); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *PostgresRepository) DeleteTask(ctx context.Context, taskID string) error {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("Failed to delete task: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresRepository) MarkTaskDone(ctx context.Context, taskID string) (*entities.Task, error) {
	query := `
        UPDATE tasks
        SET done = TRUE
        WHERE id = $1
        RETURNING id, title, content, done, user_id
    `
	task := &entities.Task{}
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(&task.ID, &task.Title, &task.Content, &task.Done, &task.UserID)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to mark task as done: %v", err)
	}
	return task, nil
}