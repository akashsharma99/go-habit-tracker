package model

import (
	"time"

	"github.com/google/uuid"
)

type Completion struct {
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

type Habit struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Completions []Completion `json:"completions"`
	CreatedAt   time.Time    `json:"created_at"`
}

func NewHabit(name string) *Habit {
	return &Habit{
		ID:          uuid.New().String(),
		Name:        name,
		Completions: []Completion{},
		CreatedAt:   time.Now(),
	}
}

func (h *Habit) ToggleToday() {
	today := time.Now().Format("2006-01-02")
	for i, c := range h.Completions {
		if c.Date == today {
			h.Completions[i].Completed = !c.Completed
			return
		}
	}
	// If not found, add new completion for today
	h.Completions = append(h.Completions, Completion{Date: today, Completed: true})
}

func (h *Habit) IsCompletedToday() bool {
	today := time.Now().Format("2006-01-02")
	for _, c := range h.Completions {
		if c.Date == today {
			return c.Completed
		}
	}
	return false
}

func (h *Habit) GetCompletionRate() float64 {
	if len(h.Completions) == 0 {
		return 0.0
	}
	completed := 0
	for _, c := range h.Completions {
		if c.Completed {
			completed++
		}
	}
	return float64(completed) / float64(len(h.Completions)) * 100
}
