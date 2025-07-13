package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"habit-tracker/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".habit-tracker")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "habits.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS habits (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS completions (
			habit_id TEXT,
			completed_date DATE,
			is_completed BOOLEAN,
			PRIMARY KEY (habit_id, completed_date),
			FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddHabit(habit *model.Habit) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO habits (id, name, created_at)
		VALUES (?, ?, ?)
	`, habit.ID, habit.Name, habit.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) DeleteHabit(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete completions first due to foreign key constraint
	_, err = tx.Exec("DELETE FROM completions WHERE habit_id = ?", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM habits WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetHabits() []*model.Habit {
	habits := make([]*model.Habit, 0)
	rows, err := s.db.Query(`
		SELECT h.id, h.name, h.created_at,
			c.completed_date, CASE WHEN c.is_completed = 1 THEN 1 ELSE 0 END as is_completed
		FROM habits h
		LEFT JOIN completions c ON h.id = c.habit_id
		ORDER BY h.created_at, c.completed_date
	`)
	if err != nil {
		return habits
	}
	defer rows.Close()

	habitMap := make(map[string]*model.Habit)
	for rows.Next() {
		var id, name string
		var createdAt time.Time
		var completedDate sql.NullString
		var isCompleted sql.NullBool

		err := rows.Scan(&id, &name, &createdAt, &completedDate, &isCompleted)
		if err != nil {
			continue
		}

		habit, exists := habitMap[id]
		if !exists {
			habit = &model.Habit{
				ID:          id,
				Name:        name,
				CreatedAt:   createdAt,
				Completions: []model.Completion{},
			}
			habitMap[id] = habit
			habits = append(habits, habit)
		}

		if completedDate.Valid && isCompleted.Valid {
			if date, err := time.Parse("2006-01-02", completedDate.String); err == nil {
				habit.Completions = append(habit.Completions, model.Completion{
					Date:      date.Format("2006-01-02"),
					Completed: isCompleted.Bool,
				})
			}
		}
	}

	return habits
}

func (s *Storage) UpdateHabit(habit *model.Habit) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update habit name
	_, err = tx.Exec(`
		UPDATE habits 
		SET name = ?
		WHERE id = ?
	`, habit.Name, habit.ID)
	if err != nil {
		return err
	}

	// Delete old completions
	_, err = tx.Exec("DELETE FROM completions WHERE habit_id = ?", habit.ID)
	if err != nil {
		return err
	}

	// Insert new completions
	for _, c := range habit.Completions {
		if date, err := time.Parse("2006-01-02", c.Date); err == nil {
			_, err = tx.Exec(`
				INSERT INTO completions (habit_id, completed_date, is_completed)
				VALUES (?, ?, ?)
			`, habit.ID, date.Format("2006-01-02"), c.Completed)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) UpdateCompletion(habitID string, date string, completed bool) error {
	completedInt := 0
	if completed {
		completedInt = 1
	}
	_, err := s.db.Exec(`
		INSERT INTO completions (habit_id, completed_date, is_completed)
		VALUES (?, ?, ?)
		ON CONFLICT(habit_id, completed_date) DO UPDATE SET is_completed=excluded.is_completed
	`, habitID, date, completedInt)
	return err
}
