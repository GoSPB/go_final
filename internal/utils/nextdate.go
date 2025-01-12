package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"task-traker/internal/models"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Пустое правило повторения")
	}

	startDate, err := parseDate(date)
	if err != nil {
		return "", err
	}

	rule, ruleLiteral, err := parseRepeatRule(repeat)
	if err != nil {
		return "", err
	}

	switch ruleLiteral {
	case "d":
		return calculateNextDateByDays(now, startDate, rule)
	case "y":
		return calculateNextDateByYears(now, startDate)
	default:
		return "", fmt.Errorf("Некорректный литерал правила")
	}
}

func parseDate(date string) (time.Time, error) {
	startDate, err := time.Parse(models.DateFormat, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("Неверный формат даты: %v", err)
	}
	return startDate, nil
}

func parseRepeatRule(repeat string) ([]string, string, error) {
	rule := strings.Split(repeat, " ")
	ruleLiteral := rule[0]

	if ruleLiteral == "d" && len(rule) < 2 {
		return nil, "", fmt.Errorf("Не указано число дней")
	}

	return rule, ruleLiteral, nil
}

func calculateNextDateByDays(now time.Time, startDate time.Time, rule []string) (string, error) {
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
}

func calculateNextDateByYears(now time.Time, startDate time.Time) (string, error) {
	newDate := startDate.AddDate(1, 0, 0)
	for newDate.Before(now) {
		newDate = newDate.AddDate(1, 0, 0)
	}
	return newDate.Format(models.DateFormat), nil
}
