package task

import "time"

// Occurrences возвращает все даты в диапазоне [from, to] включительно,
// на которые задача должна выполняться согласно настройкам периодичности
func (r *Recurrence) Occurrences(from, to time.Time) []time.Time {
	from = truncateToDay(from)
	to = truncateToDay(to)

	var result []time.Time

	for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
		if r.matches(day) {
			result = append(result, day)
		}
	}

	return result
}

// matches проверяет, попадает ли день day под правило периодичности
// Для daily: вычисляем количество полных дней от StartAt до day
// diff >= 0 исключает даты раньше StartAt
// diff%Interval == 0 означает что день кратен шагу
// Hours()/24 используется вместо встроенного метода, так как time.Duration не имеет Days()
func (r *Recurrence) matches(day time.Time) bool {
	switch r.Type {
	case RecurrenceDaily:
		if r.Interval <= 0 {
			return false
		}
		// сколько прошло дней от начальной даты до текущей
		startAt := truncateToDay(r.StartAt)
		diff := int(day.Sub(startAt).Hours() / 24)
		return diff >= 0 && diff%r.Interval == 0

	case RecurrenceMonthly:
		for _, d := range r.MonthDays {
			if day.Day() == d {
				return true
			}
		}
	case RecurrenceSpecificDates:
		for _, d := range r.Dates {
			if day.Year() == d.Year() && day.Month() == d.Month() && day.Day() == d.Day() {
				return true
			}
		}

	case RecurrenceParity:
		even := day.Day()%2 == 0
		return (r.Parity == "even") == even
	}

	return false

}

// хелпер
func truncateToDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
