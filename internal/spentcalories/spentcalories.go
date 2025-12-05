package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	activity3Parts := strings.Split(data, ",")
	if len(activity3Parts) != 3 {
		return 0, "", 0, fmt.Errorf("ошибка разделения данных на 3 части! Количество частей на данный момент: %d", len(activity3Parts))
	}
	stepsPart, err := strconv.Atoi(activity3Parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepsPart <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	if activity3Parts[1] == "" {
		return 0, "", 0, fmt.Errorf("пустая строка")
	}
	timePart, err := time.ParseDuration(activity3Parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка разбора времени: %v", err)
	}
	if timePart <= 0 {
		return 0, "", 0, fmt.Errorf("количество времени должно быть больше нуля")
	}
	return stepsPart, activity3Parts[1], timePart, nil
}

func distance(steps int, height float64) float64 {
	calcLenStep := height * stepLengthCoefficient
	distStep := float64(steps) * calcLenStep
	distStepInKm := distStep / mInKm
	return distStepInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Seconds() <= 0 {
		return 0
	}
	calcDistance := distance(steps, height)
	if calcDistance <= 0 {
		return 0
	}
	return calcDistance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsAct, typeOfAct, timeAct, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("ошибка расчета параметров тренировки")
	}
	calcCalories := 0.0
	switch typeOfAct {
	case "Ходьба":
		calcCalories, err = WalkingSpentCalories(stepsAct, weight, height, timeAct)
		if err != nil {
			return "", err
		}
	case "Бег":
		calcCalories, err = RunningSpentCalories(stepsAct, weight, height, timeAct)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	distanceTotal := distance(stepsAct, height)
	meanSpeedTotal := meanSpeed(stepsAct, height, timeAct)
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		typeOfAct, timeAct.Hours(), distanceTotal, meanSpeedTotal, calcCalories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration.Seconds() <= 0 {
		return 0, fmt.Errorf("ошибка в расчётах калорий во время бега. Значения шагов, продолжительности, веса или роста должны быть больше нуля")
	}
	calcMeanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	runCalories := (weight * calcMeanSpeed * durationInMinutes) / minInH
	if runCalories <= 0 {
		return 0, fmt.Errorf("ошибка расчета калорий при беге")
	}
	return runCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration.Seconds() < 0 {
		return 0, fmt.Errorf("ошибка в расчётах калорий во время ходьбы. Значения шагов, продолжительности, веса или роста должны быть больше нуля")
	}
	calcMeanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calcRunSpCalories := (weight * calcMeanSpeed * durationInMinutes) / minInH
	if calcRunSpCalories <= 0 {
		return 0, fmt.Errorf("ошибка расчета калорий при хотьбе")
	}
	return calcRunSpCalories * walkingCaloriesCoefficient, nil
}
