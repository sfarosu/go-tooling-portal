package helper

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func DisableDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func KeysCleanup() {
	// cleanup the tmp folder; key age deletion is defined in the time.Afterfunc function
	cleanupCmd := exec.Command("find", "web/tmp/", "-type", "f", "-name", "id_rsa*", "-mmin", "+0", "-exec", "rm", "{}", ";")
	outputCleanupCmd, err := cleanupCmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + string(outputCleanupCmd))
	} else {
		log.Println("SSH keys stored locally have been purged !")
	}
}

func RandomString(size int, Uppercase bool, Lowercase bool, Numbers bool, Specials bool) (string, error) {
	// error handling for the 2 inputs, string size can't be lower than 4 or higher than 64 AND at least one option should be selected
	if size < 4 {
		return "", errors.New("string length must be at least 4 chars long")
	}
	if size > 64 {
		return "", errors.New("string length most not exceed 64 chars")
	}
	if Uppercase == false && Lowercase == false && Numbers == false && Specials == false {
		return "", errors.New("at least one of the categories must be chosen")
	}

	rand.Seed(time.Now().UnixNano()) // used so that the shuffled result is NOT always the same

	optionsActive := 0 // we use this to count how many choices the user has selected

	// define the categories/slices and their content
	var uppercase = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var lowercase = []rune("abcdefghijklmnopqrstuvwxyz")
	var numbers = []rune("0123456789")
	var specials = []rune(";#$%&'()*+,-.:;<=>?@[]^_`{|}~")

	// selectedList contains all the categories/chars types the user selects
	var selectedList []rune
	if Uppercase == true {
		selectedList = append(selectedList, uppercase...)
		optionsActive++
	}
	if Lowercase == true {
		selectedList = append(selectedList, lowercase...)
		optionsActive++
	}
	if Numbers == true {
		selectedList = append(selectedList, numbers...)
		optionsActive++
	}
	if Specials == true {
		selectedList = append(selectedList, specials...)
		optionsActive++
	}

	// partialResult makes sure that at least ONE element is added from each category the user selects
	var partialResult []rune
	if Uppercase == true {
		partialResult = append(partialResult, uppercase[rand.Intn(len(uppercase))])
	}
	if Lowercase == true {
		partialResult = append(partialResult, lowercase[rand.Intn(len(lowercase))])
	}
	if Numbers == true {
		partialResult = append(partialResult, numbers[rand.Intn(len(numbers))])
	}
	if Specials == true {
		partialResult = append(partialResult, specials[rand.Intn(len(specials))])
	}

	// finalResult is composed of 2 slices and because append always adds the second slice to the end of the first one we use the last FOR to randomize everything
	finalResult := make([]rune, size-optionsActive)
	for i := range finalResult {
		finalResult[i] = selectedList[rand.Intn(len(selectedList))]
	}
	finalResult = append(finalResult, partialResult...)

	// now that the finalResult is complete we shuffle everything
	for x := len(finalResult) - 1; x > 0; x-- {
		y := rand.Intn(x + 1)
		finalResult[x], finalResult[y] = finalResult[y], finalResult[x]
	}
	return string(finalResult), nil
}

func ReadFile(filePath string) []byte {
	byteData, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		log.Println("error reading file: ", errReadFile)
	}
	return byteData
}
