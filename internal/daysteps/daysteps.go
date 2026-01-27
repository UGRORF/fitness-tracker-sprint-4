package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	str := strings.Split(data, ",")
	if len(str) != 2 {
		return 0, 0, errors.New("invalid data")
	}

	steps, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, 0, errors.New("invalid data, the first substring should have int")
	}

	if steps <= 0 {
		return 0, 0, errors.New("invalid data, the count of steps should be positive")
	}

	walkDuration, err := time.ParseDuration(str[1])
	if err != nil {
		return 0, 0, errors.New("invalid data, the second substring should have time.Duration")
	}

	if walkDuration.Seconds() <= 0 {
		return 0, 0, errors.New("invalid data, the duration should be positive")
	}

	return steps, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		log.Println(err)
		return ""
	}

	distance := float64(steps) * stepLength / float64(mInKm)

	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\n"+
		"Дистанция составила %.2f км.\n"+
		"Вы сожгли %.2f ккал.\n", steps, distance, spentCalories)
}
