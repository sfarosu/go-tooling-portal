package main

import (
    "net/http"
    "errors"
    "math/rand"
    "time"
    "strconv"
    "log"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    passGen = promauto.NewCounter(prometheus.CounterOpts{
        Name: "passwords_generated_total",
        Help: "The total number of generated passwords",
    })
)

func passgen(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Redirect(w, r, "/passgen", http.StatusSeeOther)
    }
    log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
    tpl.ExecuteTemplate(w, "passgen.html", nil)
}

func passgenProcess(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/passgen", http.StatusSeeOther)
    }

    // Get the checkbox status from the user and convert it to bool before sending it to the randomString() function
    upper := r.FormValue("uppercase")
    upperBool := false; _ = upperBool
    if len(upper) != 0 {
        upperBool = true
    }

    lower := r.FormValue("lowercase")
    lowerBool := false; _ = lowerBool
    if len(lower) != 0 {
        lowerBool = true
    }

    numbers := r.FormValue("numbers")
    numbersBool := false; _ = numbersBool
    if len(numbers) != 0 {
        numbersBool = true
    }

    symbols := r.FormValue("symbols")
    symbolsBool := false; _ = symbolsBool
    if len(symbols) != 0 {
        symbolsBool = true
    }

    length := r.FormValue("number")

    lengthInt, err := strconv.Atoi(length) // convert length string to int
    generatedRamdomPassord, err := randomString(lengthInt, upperBool, lowerBool, numbersBool, symbolsBool)
    if err != nil {
        log.Println(err)
    }

    data := struct {
        Uppercase string
        Lowercase string
        Numbers string
        Symbols string
        Length string
        Result string
    }{
        Uppercase: upper,
        Lowercase: lower,
        Numbers: numbers,
        Symbols: symbols,
        Length: length,
        Result: generatedRamdomPassord,
    }
    log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
    tpl.ExecuteTemplate(w, "passgen-process.html", data)

    passGen.Inc()
}

func randomString(size int, Uppercase bool, Lowercase bool, Numbers bool, Specials bool) (string, error) {
    // error handling for the 2 inputs, string size can't be lower than 4 or higher than 64 AND at least one option should be selected
    if size < 4 {
        return "", errors.New("String length must be at least 4 chars long")
    }
    if size > 64 {
        return "", errors.New("String length most not exceed 64 chars")
    }
    if (Uppercase == false && Lowercase == false && Numbers == false && Specials == false) {
        return "", errors.New("At least one of the categories must be chosen")
    }

    rand.Seed(time.Now().UnixNano()) // used so that the shuffled result is NOT always the same

    optionsActive := 0 // we use this to count how many choices the user has selected

    // define the categories/slices and their content
    var uppercase = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    var lowercase = []rune("abcdefghijklmnopqrstuvwxyz")
    var numbers = []rune("0123456789")
    var specials = []rune(";#$%&'()*+,-.:;<=>?@[]^_`{|}~")

    // selectedlist contains all the categories/chars types the user selects
    var selectedList = []rune{}
    if Uppercase == true {
        selectedList = append(selectedList, uppercase...)
        optionsActive ++
    }
    if Lowercase == true {
        selectedList = append(selectedList, lowercase...)
        optionsActive ++
    }
    if Numbers == true {
        selectedList = append(selectedList, numbers...)
        optionsActive ++
    }
    if Specials == true {
        selectedList = append(selectedList, specials...)
        optionsActive ++
    }

    // partialResult makes sure that at least ONE element is added from each category the user selects
    var partialResult = []rune{}
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

    // finalresult is composed of 2 slices and because append always adds the second slice to the end of the first one we use the last FOR to randomize everything
    finalResult := make([]rune, size - optionsActive)
    for i := range finalResult {
        finalResult[i] = selectedList[rand.Intn(len(selectedList))]
    }
    finalResult = append(finalResult, partialResult ...)

    // now that the finalresult is complete we shuffle everything
    for x := len(finalResult) - 1; x > 0; x-- {
        y := rand.Intn(x + 1)
        finalResult[x], finalResult[y] = finalResult[y], finalResult[x]
    }
    return string(finalResult), nil
}