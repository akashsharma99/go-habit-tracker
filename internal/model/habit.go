package model

import "time"

type Habit struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Completions map[string]bool `json:"completions"` // date -> completed
	CreatedAt   time.Time       `json:"created_at"`
}

func NewHabit(name string) *Habit {
	return &Habit{
		ID:          time.Now().Format("20060102150405"),
		Name:        name,
		Completions: make(map[string]bool),
		CreatedAt:   time.Now(),
	}
}

func (h *Habit) ToggleToday() {
	today := time.Now().Format("2006-01-02")
	if h.Completions == nil {
		h.Completions = make(map[string]bool)
	}
	h.Completions[today] = !h.Completions[today]
}

func (h *Habit) IsCompletedToday() bool {
	today := time.Now().Format("2006-01-02")
	if h.Completions == nil {
		h.Completions = make(map[string]bool)
		return false
	}
	return h.Completions[today]
}

func (h *Habit) GetCompletionRate() float64 {
	if len(h.Completions) == 0 {
		return 0.0
	}
	completed := 0
	for _, isCompleted := range h.Completions {
		if isCompleted {
			completed++
		}
	}
	return float64(completed) / float64(len(h.Completions)) * 100
}
