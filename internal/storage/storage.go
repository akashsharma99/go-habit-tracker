package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"habit-tracker/internal/model"
)

type Storage struct {
	habits map[string]*model.Habit
	file   string
	mutex  sync.RWMutex
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

	file := filepath.Join(dataDir, "habits.json")
	storage := &Storage{
		habits: make(map[string]*model.Habit),
		file:   file,
	}

	if err := storage.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.habits)
}

func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.habits, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, data, 0644)
}

func (s *Storage) AddHabit(habit *model.Habit) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.habits[habit.ID] = habit
	return s.save()
}

func (s *Storage) DeleteHabit(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.habits, id)
	return s.save()
}

func (s *Storage) GetHabits() []*model.Habit {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	habits := make([]*model.Habit, 0, len(s.habits))
	for _, habit := range s.habits {
		if habit.Completions == nil {
			habit.Completions = make(map[string]bool)
		}
		habits = append(habits, habit)
	}

	// Sort habits by creation time to ensure consistent ordering
	sort.Slice(habits, func(i, j int) bool {
		return habits[i].CreatedAt.Before(habits[j].CreatedAt)
	})

	return habits
}

func (s *Storage) UpdateHabit(habit *model.Habit) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.habits[habit.ID] = habit
	return s.save()
}
