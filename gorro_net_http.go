package gorro

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

type queryRegex struct {
	regexIndex    int
	subregIndexes []int
}

type NetHttpRouter struct {
	regex           string
	compiledRegex   *regexp.Regexp
	regexMap        map[int]queryRegex
	handlers        []map[string]Handler
	notFoundHandler func(w http.ResponseWriter, r *http.Request)
	errorHandler    func(w http.ResponseWriter, r *http.Request, err error)
	ran             bool
}

func NewRouter() *NetHttpRouter {
	handlers := make([]map[string]Handler, 0, 0)
	return &NetHttpRouter{"", nil, nil, handlers, nil, nil, false}
}

func (dr *NetHttpRouter) Register(regex string, handlers HandlersMap) error {
	var err error
	log.Printf("Registering route: %s", regex)
	for m := range handlers {
		err = ensureValidMethod(m)
		if err != nil {
			return err
		}
	}
	dr.regex = dr.regex + separator(dr) + "(" + regex + ")"
	dr.handlers = append(dr.handlers, handlers)
	dr.regexMap = getPatterns([]rune(dr.regex))
	if dr.ran {
		dr.compiledRegex = nil
	}
	return nil
}

func (dr *NetHttpRouter) Route(w http.ResponseWriter, r *http.Request) error {
	dr.ran = true
	log.Printf("Received: %s", r.URL.Path)

	start := time.Now()
	results := regEx(dr).FindStringSubmatch(r.URL.Path)

	if len(results) != 0 {
		regexIndex, subRegexes := getRegexIndexAndSubPatterns(results, dr.regexMap)

		if regexIndex != -1 {
			log.Printf("Finding the route took %s", time.Since(start))
			start = time.Now()

			if handler, ok := dr.handlers[subRegexes.regexIndex][r.Method]; ok {
				err := handler(w, toRequest(r, dr.compiledRegex, results, subRegexes.subregIndexes, dr.regex))

				if err != nil {
					handleError(err, dr, w, r)
				}
			}

			log.Printf("Executing the handler took %s", time.Since(start))
		} else {
			log.Printf("No handlers found for %s", r.RequestURI)
		}
	} else {
		handleNotFound(dr, w, r)
	}
	return nil
}

func (dr *NetHttpRouter) OnNotFound(handler func(w http.ResponseWriter, r *http.Request)) {
	dr.notFoundHandler = handler
}

func (dr *NetHttpRouter) OnError(handler func(w http.ResponseWriter, r *http.Request, err error)) {
	dr.errorHandler = handler
}

func handleNotFound(dr *NetHttpRouter, w http.ResponseWriter, r *http.Request) {
	if dr.notFoundHandler != nil {
		dr.notFoundHandler(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func handleError(err error, dr *NetHttpRouter, w http.ResponseWriter, r *http.Request) {
	if dr.errorHandler != nil {
		dr.errorHandler(w, r, err)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRegexIndexAndSubPatterns(results []string, patterns map[int]queryRegex) (int, queryRegex) {
	for i := 1; i < len(results); i++ {
		if results[i] != "" {
			return i, patterns[i]
		}
	}
	return -1, queryRegex{}
}

func getPatterns(regex []rune) map[int]queryRegex {
	var ret map[int]queryRegex = make(map[int]queryRegex)
	var currentIndex = 0
	var subregexIndex = 0
	for i := 0; i < len(regex); i++ {
		if regex[i] == '(' {
			subregexIndex++
			subregexes, pos, s := getFlattened(regex, i, subregexIndex)
			i = pos
			ret[subregexIndex] = queryRegex{currentIndex, subregexes}
			currentIndex++
			subregexIndex = s
		}
		i++
	}
	return ret
}

func getFlattened(regex []rune, from int, regexIndex int) ([]int, int, int) {
	var subRegexesIndex []int = make([]int, 0, 5)
	var length int = len(regex)
	var openCount = 1
	var i int
	for i = from + 1; i < length; {
		if regex[i] == '(' {
			openCount++
			isCapturing := i+2 < length && (regex[i+1] != '?' || regex[i+2] != ':')
			if isCapturing {
				regexIndex++
				subRegexesIndex = append(subRegexesIndex, regexIndex)
			}
		} else if regex[i] == ')' {
			openCount--
			if openCount == 0 {
				return subRegexesIndex, i, regexIndex
			}
		}
		i++
	}
	return subRegexesIndex, i, regexIndex
}

func toRequest(hr *http.Request, regex *regexp.Regexp, results []string, paramsIndexes []int, rexStr string) *Request {
	namedParams := make(map[string]string)
	params := make([]string, 0)

	subexpNames := regex.SubexpNames()

	for i := 0; i < len(paramsIndexes); i++ {
		if subexpNames[paramsIndexes[i]] != "" {
			namedParams[subexpNames[paramsIndexes[i]]] = results[paramsIndexes[i]]
		}
		params = append(params, results[paramsIndexes[i]])
	}

	r := Request{*hr, namedParams, params, results, rexStr}
	return &r
}

func regEx(dr *NetHttpRouter) *regexp.Regexp {
	if dr.compiledRegex == nil {
		dr.compiledRegex = regexp.MustCompile(dr.regex)
	}
	return dr.compiledRegex
}

func separator(dr *NetHttpRouter) string {
	if dr.regex == "" {
		return ""
	}
	return "|"
}

func ensureValidMethod(m string) error {
	if m != http.MethodConnect && m != http.MethodDelete &&
		m != http.MethodGet && m != http.MethodHead &&
		m != http.MethodOptions && m != http.MethodPatch &&
		m != http.MethodPost && m != http.MethodPut &&
		m != http.MethodTrace {
		return errors.New(fmt.Sprintf("Invalid http method received: %s", m))
	}
	return nil
}
