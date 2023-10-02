package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type wordComposition struct {
	letters [33]int
}

var (
	ErrNonRussianLetter = errors.New("there is non-russian letter")
)

func findAnagrams(dictionary []string) (map[string][]string, error) {
	mainMap := make(map[string][]string)
	compToFirstWord := make(map[wordComposition]string)

	for i, str := range dictionary {
		dictionary[i] = strings.ToLower(str)
	}

	dictionary = filterUnique(dictionary)

	for i, word := range dictionary {
		c := wordComposition{}

		for _, r := range word {
			if err := insertRune(&c, r); err != nil {
				return nil, fmt.Errorf("error on word #%d-\"%s\": %w", i, word, err)
			}
		}

		var mainMapKey string
		var ok bool
		if mainMapKey, ok = compToFirstWord[c]; !ok {
			compToFirstWord[c] = word
			mainMapKey = word
		}

		mainMap[mainMapKey] = append(mainMap[mainMapKey], word)
	}

	for key, slice := range mainMap {
		if len(slice) <= 1 {
			delete(mainMap, key)
			continue
		}

		slices.Sort(slice)
	}

	return mainMap, nil
}

func filterUnique(dictionary []string) []string {
	m := make(map[string]struct{})

	newSlice := make([]string, 0)
	for _, elem := range dictionary {
		if _, ok := m[elem]; !ok {
			m[elem] = struct{}{}
			newSlice = append(newSlice, elem)
		}
	}

	return newSlice
}

func insertRune(wordComposition *wordComposition, r rune) error {
	if r-'а' < 0 || r-'а' > 32 {
		return ErrNonRussianLetter
	}

	wordComposition.letters[r-'а']++
	return nil
}
