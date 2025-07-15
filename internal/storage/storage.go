package storage

import (
	"database/sql"
	"log"
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
			completed_date TEXT,
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

	_, err = tx.Exec(`
		INSERT INTO completions (habit_id, completed_date, is_completed)
		VALUES (?, ?, ?)
	`, habit.ID, time.Now().Format("2006-01-02"), false)
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
	habitMap := make(map[string]*model.Habit)
	habits := make([]*model.Habit, 0)

	rows, err := s.db.Query(`
		SELECT h.id, h.name, h.created_at, c.completed_date, c.is_completed
		FROM habits h
		LEFT JOIN completions c ON h.id = c.habit_id
		ORDER BY h.created_at DESC
	`)
	if err != nil {
		return habits
	}
	defer rows.Close()

	for rows.Next() {
		var habitID, name, completedDate string
		var createdAt time.Time
		var isCompleted bool

		if err := rows.Scan(&habitID, &name, &createdAt, &completedDate, &isCompleted); err != nil {
			continue
		}

		if _, ok := habitMap[habitID]; !ok {
			habitMap[habitID] = &model.Habit{
				ID:          habitID,
				Name:        name,
				CreatedAt:   createdAt,
				Completions: []model.Completion{},
			}
			habits = append(habits, habitMap[habitID])
		}

		habitMap[habitID].Completions = append(habitMap[habitID].Completions, model.Completion{
			Date:      completedDate,
			Completed: isCompleted,
		})
	}
	// log the loaded habits for debugging
	for _, habit := range habitMap {
		log.Printf("Loaded habit: %s with completions: %v", habit.Name, habit.Completions)
	}

	return habits
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
