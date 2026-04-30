package handlers

import (
	"fmt"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

const dateLayout = "2006-01-02"

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *recurrenceDTO    `json:"recurrence,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *recurrenceDTO    `json:"recurrence,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Recurrence:  toRecurrenceDTO(task.Recurrence),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

type recurrenceDTO struct {
	Type      string   `json:"type"`
	StartAt   string   `json:"start_at,omitempty"`
	Interval  int      `json:"interval,omitempty"`
	MonthDays []int    `json:"month_days,omitempty"`
	Dates     []string `json:"dates,omitempty"`
	Parity    string   `json:"parity,omitempty"`
}

func toRecurrenceDTO(r *taskdomain.Recurrence) *recurrenceDTO {
	if r == nil {
		return nil
	}

	dto := &recurrenceDTO{
		Type:      string(r.Type),
		Interval:  r.Interval,
		MonthDays: r.MonthDays,
		Parity:    r.Parity,
	}

	if !r.StartAt.IsZero() {
		dto.StartAt = r.StartAt.Format(dateLayout)
	}

	for _, d := range r.Dates {
		dto.Dates = append(dto.Dates, d.Format(dateLayout))
	}

	return dto
}

func toRecurrenceDomain(dto *recurrenceDTO) (*taskdomain.Recurrence, error) {
	if dto == nil {
		return nil, nil
	}

	rec := &taskdomain.Recurrence{
		Type:      taskdomain.RecurrenceType(dto.Type),
		Interval:  dto.Interval,
		MonthDays: dto.MonthDays,
		Parity:    dto.Parity,
	}

	if dto.StartAt != "" {
		t, err := time.ParseInLocation(dateLayout, dto.StartAt, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("Invalid start_at: %w", err)
		}
		rec.StartAt = t
	}

	for _, s := range dto.Dates {
		t, err := time.ParseInLocation(dateLayout, s, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("Invalid date %q: %w", s, err)
		}
		rec.Dates = append(rec.Dates, t)
	}

	return rec, nil
}
