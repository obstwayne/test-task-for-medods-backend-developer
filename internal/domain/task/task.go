package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      Status      `json:"status"`
	Recurrence  *Recurrence `json:"recurrence,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type RecurrenceType string

const (
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceParity        RecurrenceType = "parity"
)

type Recurrence struct {
	Type      RecurrenceType `json:"type"`
	StartAt   time.Time      `json:"start_at"`             // точка отсчета
	Interval  int            `json:"interval,omitempty"`   // For daily (e.g., every 2 days)
	MonthDays []int          `json:"month_days,omitempty"` // For monthly
	Dates     []time.Time    `json:"dates,omitempty"`      // For specific dates
	Parity    string         `json:"parity,omitempty"`     // For parity (e.g., "even", "odd")
}

func (r *Recurrence) Valid() bool {
	switch r.Type {
	case RecurrenceDaily:
		return r.Interval > 0 && !r.StartAt.IsZero()
	case RecurrenceMonthly:
		if len(r.MonthDays) == 0 {
			return false
		}
		for _, d := range r.MonthDays {
			if d < 1 || d > 30 {
				return false
			}
		}
		return true
	case RecurrenceSpecificDates:
		return len(r.Dates) > 0
	case RecurrenceParity:
		return r.Parity == "even" || r.Parity == "odd"
	default:
		return false
	}
}
