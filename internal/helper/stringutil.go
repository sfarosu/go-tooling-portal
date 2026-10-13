package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// RandomString generates a random string of a specified size using the selected character types
func RandomString(size int, Uppercase bool, Lowercase bool, Numbers bool, Specials bool) (string, error) {
	// optionsActive counts how many choices the user has selected
	optionsActive := 0

	// define the categories/slices and their content
	var uppercase = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var lowercase = []rune("abcdefghijklmnopqrstuvwxyz")
	var numbers = []rune("0123456789")
	var specials = []rune(";#$%&'()*+,-.:;<=>?@[]^_`{|}~")

	// selectedList contains all the categories/chars types the user selects
	var selectedList []rune
	if Uppercase {
		selectedList = append(selectedList, uppercase...)
		optionsActive++
	}
	if Lowercase {
		selectedList = append(selectedList, lowercase...)
		optionsActive++
	}
	if Numbers {
		selectedList = append(selectedList, numbers...)
		optionsActive++
	}
	if Specials {
		selectedList = append(selectedList, specials...)
		optionsActive++
	}

	// validate that the requested size is at least equal to the number of selected categories
	if size < optionsActive {
		return "", fmt.Errorf("password length [%v] must be at least equal to number of selected character classes [%v]", size, optionsActive)
	}

	// partialResult makes sure that at least ONE element is added from each category the user selects
	var partialResult []rune
	if Uppercase {
		char, err := randomRune(uppercase)
		if err != nil {
			return "", err
		}
		partialResult = append(partialResult, char)
	}
	if Lowercase {
		char, err := randomRune(lowercase)
		if err != nil {
			return "", err
		}
		partialResult = append(partialResult, char)
	}
	if Numbers {
		char, err := randomRune(numbers)
		if err != nil {
			return "", err
		}
		partialResult = append(partialResult, char)
	}
	if Specials {
		char, err := randomRune(specials)
		if err != nil {
			return "", err
		}
		partialResult = append(partialResult, char)
	}

	// finalResult is composed of 2 slices and because append always adds the second slice to the end of the first one, we use the last FOR to randomize everything
	finalResult := make([]rune, size-optionsActive)
	for i := range finalResult {
		char, err := randomRune(selectedList)
		if err != nil {
			return "", err
		}
		finalResult[i] = char
	}
	finalResult = append(finalResult, partialResult...)

	// now that the finalResult is complete we shuffle everything
	for x := len(finalResult) - 1; x > 0; x-- {
		index, err := randomIndex(x + 1)
		if err != nil {
			return "", err
		}
		finalResult[x], finalResult[index] = finalResult[index], finalResult[x]
	}
	return string(finalResult), nil
}

func randomRune(characters []rune) (rune, error) {
	index, err := randomIndex(len(characters))
	if err != nil {
		return 0, err
	}
	return characters[index], nil
}

func randomIndex(limit int) (int, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(limit)))
	if err != nil {
		return 0, fmt.Errorf("failed to generate secure random value: %w", err)
	}
	return int(index.Int64()), nil
}

// AddSecondDigit adds a leading zero to single-digit numbers
func AddSecondDigit(number int) string {
	return fmt.Sprintf("%02d", number)
}

// GetNumberDigitsAmmount returns the number of digits in a given integer
func GetNumberDigitsAmmount(number int64) int {
	count := 0
	for number > 0 {
		number = number / 10
		count = count + 1
	}

	return count
}
