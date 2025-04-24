package coding

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"st/norm"
	"time"

	anotherRand "golang.org/x/exp/rand"
)

// Исправление одного бита ошибки на основе синдрома ошибки
func fixMistake(mistake byte, data byte) byte {
	switch mistake {
	case 1, 2, 4, 8, 16, 32, 64:
		return data ^ mistake
	default:
		return data
	}
}

// Кодирование одного 4-битного блока (вход — 1 байт, выход — 7-битный код с проверкой на ошибки)
func rcrEncode(data byte) byte {
	num := data << 3 // смещаем влево на 3 бита, оставляя место для контрольных
	buf := num
	pol := byte(11 << 3) // полином для деления

	for i := 0; i < 4; i++ {
		if bits.Len8(buf) < bits.Len8(pol) {
			pol >>= 1
			continue
		}
		buf ^= pol
		pol >>= 1
	}
	return num | buf // объединяем исходные данные с остатком
}

// Декодирование одного блока данных, возвращает байт и флаг валидности
func decode(data []byte) (byte, bool) {
	valid := true

	for i := range data {
		pol := byte(11) << 3
		buf := data[i]
		var res byte

		for j := 0; j < 4; j++ {
			if bits.Len8(buf) < 4 {
				res = buf
				break
			}
			if bits.Len8(buf) < bits.Len8(pol) {
				pol >>= 1
				continue
			}
			res = buf ^ pol
			buf = res
			pol >>= 1
		}

		// Ошибка не найдена
		if res == 0 {
			data[i] = data[i] >> 3
		} else {
			// Попытка исправления — если более 1 ошибки, помечаем как невалидный
			if bits.OnesCount8(res) > 1 {
				valid = false
			}
			data[i] = fixMistake(res, data[i]) >> 3
		}
	}

	// Склеиваем два полубайта в один байт
	return data[0]<<4 | data[1], valid
}

// Разделение байта на два полубайта и кодирование каждого
func encode(data byte) []byte {
	b1 := data & 240 >> 4 // первые 4 бита
	b2 := data & 15       // вторые 4 бита
	return []byte{rcrEncode(b1), rcrEncode(b2)}
}

// Внесение 1–3 случайных ошибок в 1 байт с вероятностной моделью
func makeAdvancedMistakes(data []byte, byteIndex int, randSrc anotherRand.Source) {
	r := rand.New(randSrc)

	chance := r.Intn(100)
	numErrors := 1
	if chance < 10 {
		numErrors = 3 // 10% шанс на 3 ошибки
	} else if chance < 30 {
		numErrors = 2 // 20% шанс на 2 ошибки
	}

	for i := 0; i < numErrors; i++ {
		bitToFlip := r.Intn(7)
		byteInPair := r.Intn(2)
		data[byteInPair] ^= (1 << bitToFlip) // инвертируем случайный бит
	}
}

// Обработка полного сообщения: кодирование, внесение ошибок, декодирование
func ProcessMessage(msg string, randSrc anotherRand.Source) (res string, err error) {
	// Случайная позиция для внесения ошибок
	errorPos := norm.GenerateNormalInt(0, min(len(msg), 100), 130/8, 13, randSrc)

	var processedMsg []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 5% вероятность потери всего сообщения
	if r.Intn(100)%20 == 0 {
		fmt.Println("lost message ")
		return "", errors.New("lost message")
	}

	encodedMsg := []byte{}
	messageWithErrors := []byte{}
	decodedMessage := []byte{}

	for i := 0; i < len(msg); i++ {
		data := encode(msg[i]) // кодируем байт
		encodedMsg = append(encodedMsg, data...)

		if i == errorPos {
			makeAdvancedMistakes(data, i, randSrc) // вносим ошибки в выбранный байт
		}

		messageWithErrors = append(messageWithErrors, data...)

		decodedByte, valid := decode(data) // декодируем

		if !valid {
			fmt.Printf("Не удалось восстановить байт %d. Байт утерян.\n", i)
			return "", errors.New("uncorrectable block error")
		}

		decodedMessage = append(decodedMessage, decodedByte)
		processedMsg = append(processedMsg, decodedByte)
	}

	// Вывод всех этапов в консоль
	fmt.Println("Исходное сообщение: ", msg)
	fmt.Println("Закодированное сообщение: ", string(encodedMsg))
	fmt.Println("Сообщение с ошибкой: ", string(messageWithErrors))
	fmt.Println("Декодированное сообщение: ", string(decodedMessage))

	return string(processedMsg), nil
}

// Возвращает минимальное значение двух чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
