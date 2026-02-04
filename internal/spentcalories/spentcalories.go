package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	str := strings.Split(data, ",")

	if len(str) != 3 {
		return 0, "", 0, errors.New("invalid training data")
	}

	steps, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, "", 0, errors.New("invalid data, the first substring should have int")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("invalid data, the count of steps should be positive")
	}

	activity := str[1]

	walkDuration, err := time.ParseDuration(str[2])
	if err != nil {
		return 0, "", 0, errors.New("invalid data, the second substring should have time.Duration")
	}

	if walkDuration.Seconds() <= 0 {
		return 0, "", 0, errors.New("invalid data, the duration should be positive")
	}

	return steps, activity, walkDuration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * (height * stepLengthCoefficient) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Seconds() <= 0 {
		return 0
	}
	durationInHours := duration.Hours()

	return distance(steps, height) / durationInHours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	avgSpeed := meanSpeed(steps, height, duration)
	durationInHours := duration.Hours()
	var spentCalories float64

	switch activity {
	case "Бег":
		spentCalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	output := fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n", activity, durationInHours, dist, avgSpeed, spentCalories)

	return output, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("invalid data, the count of steps should be positive")
	}

	if weight <= 0 || height <= 0 {
		return 0, errors.New("invalid data, the count of weight and height should be positive")
	}

	if duration.Seconds() <= 0 {
		return 0, errors.New("invalid data, the duration should be positive")
	}

	durationInMinutes := duration.Minutes()
	avgSpeed := meanSpeed(steps, height, duration)

	return (weight * avgSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return runningSpentCalories * walkingCaloriesCoefficient, nil
}
