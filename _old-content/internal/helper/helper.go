package helper

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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

func RandomString(size int, Uppercase bool, Lowercase bool, Numbers bool, Specials bool) (string, error) {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// we use this to count how many choices the user has selected
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

	// finalResult is composed of 2 slices and because append always adds the second slice to the end of the first one we use the last FOR to randomize everything
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
	return string(finalResult), nil
}

func ReadFile(filePath string) []byte {
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("error reading file: ", err)
	}
	return byteData
}

func PrettyJSON(insertedText string) bytes.Buffer {
	var pretty bytes.Buffer

	err := json.Indent(&pretty, []byte(insertedText), "", "    ")
	if err != nil {
		log.Println("error indenting JSON: ", err)
	}

	return pretty
}

func MarshalJSON(data interface{}) []byte {
	byteData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("failed to marshal JSON: ", err)
	}
	return byteData
}

func UnmarshalJSON(byteData []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(byteData), &jsonData)
	if err != nil {
		log.Println("error unmarshaling JSON: ", err)
	}
	return jsonData, err
}

func MarshalYAML(data interface{}) []byte {
	var byteData bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&byteData)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&data)
	if err != nil {
		log.Println("error marshaling YAML: ", err)
	}
	return byteData.Bytes()
}

func UnmarshalYAML(byteData []byte) (map[string]interface{}, error) {
	var yamlData map[string]interface{}
	err := yaml.Unmarshal([]byte(byteData), &yamlData)
	if err != nil {
		log.Println("error unmarshaling YAML: ", err)
	}
	return yamlData, err
}

func AddSecondDigit(number int) string {
	if number < 10 {
		return "0" + strconv.Itoa(number)
	} else {
		return strconv.Itoa(number)
	}
}

func GetNumberDigitsAmmount(number int64) int {
	count := 0
	for number > 0 {
		number = number / 10
		count = count + 1
	}

	return count
}

func CurrentWorkingDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error establishing current working directory: %v", err)
	}
	return cwd
}
