package norm

import (
	"math"

	anotherRand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// GenerateNormalInt генерирует целое число в диапазоне [min, max]
// по нормальному распределению с заданным средним значением (mean) и стандартным отклонением (stddev).
// Используется кастомный источник случайных чисел (rand.Source) для воспроизводимости или повышения случайности.
func GenerateNormalInt(min, max int, mean, stddev float64, src anotherRand.Source) int {
	// Создаем распределение N(mean, stddev^2)
	normal := distuv.Normal{
		Mu:    mean,   // среднее значение
		Sigma: stddev, // стандартное отклонение
		Src:   src,    // источник случайных чисел
	}

	for {
		value := normal.Rand()         // генерируем значение с плавающей точкой
		intValue := int(math.Round(value)) // округляем до ближайшего целого

		// Проверяем попадание в диапазон
		if intValue >= min && intValue <= max {
			return intValue
		}
		// Если сгенерированное число не входит в диапазон — пробуем снова
	}
}
