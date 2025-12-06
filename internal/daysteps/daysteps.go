package daysteps

import (
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
	activity2Parts := strings.Split(data, ",")
	if len(activity2Parts) != 2 {
		return 0, 0, fmt.Errorf("ошибка разделения данных на 2 части! Количество частей на данный момент: %d", len(activity2Parts))
	}
	stepsPart, err := strconv.Atoi(activity2Parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка в данных о количестве шагов")
	}
	if stepsPart <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	timePart, err := time.ParseDuration(activity2Parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка в данных о количестве времени")
	}
	if timePart.Seconds() <= 0 {
		return 0, 0, fmt.Errorf("количество времени должно быть больше нуля")
	}
	return stepsPart, timePart, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка парсинга строки"), err)
		log.Println(fmt.Errorf("ошибка парсинга строки"), err)
		return ""
	}
	if duration <= 0 {
		fmt.Printf("продолжительность должна быть больше нуля")
		log.Printf("продолжительность должна быть больше нуля\n")
		return ""
	}
	if steps <= 0 {
		fmt.Printf("количество шагов должно быть больше нуля.")
		log.Printf("количество шагов должно быть больше нуля.\n")
		return ""
	}
	distInMeters := float64(steps) * stepLength
	distInKm := distInMeters / mInKm
	totalCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distInKm, totalCalories)
}
