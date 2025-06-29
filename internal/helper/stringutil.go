package helper

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomString generates a random string of a specified size using the selected character types
func RandomString(size int, Uppercase bool, Lowercase bool, Numbers bool, Specials bool) string {
	//TODO analize if we can tranform this function into a method of a struct that contains the options
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

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

	// partialResult makes sure that at least ONE element is added from each category the user selects
	var partialResult []rune
	if Uppercase {
		partialResult = append(partialResult, uppercase[rand.Intn(len(uppercase))])
	}
	if Lowercase {
		partialResult = append(partialResult, lowercase[rand.Intn(len(lowercase))])
	}
	if Numbers {
		partialResult = append(partialResult, numbers[rand.Intn(len(numbers))])
	}
	if Specials {
		partialResult = append(partialResult, specials[rand.Intn(len(specials))])
	}

	// finalResult is composed of 2 slices and because append always adds the second slice to the end of the first one, we use the last FOR to randomize everything
	finalResult := make([]rune, size-optionsActive)
	for i := range finalResult {
		finalResult[i] = selectedList[rand.Intn(len(selectedList))]
	}
	finalResult = append(finalResult, partialResult...)

	// now that the finalResult is complete we shuffle everything
	for x := len(finalResult) - 1; x > 0; x-- {
		y := rng.Intn(x + 1)
		finalResult[x], finalResult[y] = finalResult[y], finalResult[x]
	}
	return string(finalResult)
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
