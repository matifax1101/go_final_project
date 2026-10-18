package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateOnly.After(nowOnly)
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("пустое правило повторения")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", fmt.Errorf("интервал должен быть от 1 до 400")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}

	case "w", "m":
		return "", fmt.Errorf("правило %q не поддерживается", parts[0])

	default:
		return "", fmt.Errorf("неизвестное правило повторения: %s", parts[0])
	}
}
