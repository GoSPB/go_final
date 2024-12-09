package repeat

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GoSPB/go_final/internal/models"
)

// NextDate вычисляет следующую дату выполнения задачи на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Пустое правило повторения")
	}

	// Парсим исходную дату
	startDate, err := time.Parse(models.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("Неверный формат даты: %v", err)
	}

	// Разделяем правило на компоненты
	rule := strings.Split(repeat, " ")
	ruleLiteral := rule[0]

	switch ruleLiteral {
	case "d":
		// Правило: перенос на указанное количество дней
		if len(rule) < 2 {
			return "", fmt.Errorf("Не указано число дней")
		}
		daysN, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("Неверное число дней: %v", err)
		}
		if daysN > 400 {
			return "", fmt.Errorf("Число дней не может превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysN)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysN)
		}
		return newDate.Format(models.DateFormat), nil

	case "y":
		// Правило: ежегодное повторение
		newDate := startDate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(models.DateFormat), nil

	default:
		// Неподдерживаемое правило
		return "", fmt.Errorf("Неподдерживаемое правило повторения: %s", ruleLiteral)
	}
}
