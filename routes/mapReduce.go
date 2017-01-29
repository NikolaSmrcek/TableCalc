package routes

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"unicode"
	"io/ioutil"
	"os"
)

var WorkingDirectory string;

func LoaderFunction(w http.ResponseWriter, req *http.Request){
	statusCode := http.StatusOK;
	responseMessage := "Ok";

	fmt.Println("Starting to split string.");


	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// TODO see what's wrong with the file path
	err = ioutil.WriteFile(WorkingDirectory+handler.Filename, data, 0777)
	if err != nil {
		fmt.Println(err)
	}

	stringData := string(data)
	fmt.Println("Text that we will decode.", stringData);

	replacer := strings.NewReplacer(",", " , ", ".", " . ", ";", " ; ", "!", " ! ", "?", " ? ", ":", " : ")
	stringData = replacer.Replace(stringData)
	//Saving words and special characters to the redis List
	words := _fields(stringData)
	//punctuation
	// TODO make a seprate function call
	redisValues := make([]interface{}, len(words));
	for i, v := range words{
		redisValues[i] = v;
	}
	RedisClient.RPush("loadedElements", redisValues...);
	var totalElements interface{} = len(words);
	RedisClient.HSet("stats","totalElements",totalElements);

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(responseMessage);
}

func _fields(s string) []string {
	return _fieldsFunc(s, unicode.IsSpace)
}

func _fieldsFunc(s string, f func(rune) bool) []string {
	// First count the fields.
	n := 0
	inField := false
	for _, r := range s {
		wasInField := inField
		inField = !f(r)
		if inField && !wasInField {
			n++
		}
	}

	// Now create them.
	a := make([]string, n)
	na := 0
	fieldStart := -1 // Set to -1 when looking for start of field.
	for i, r := range s {
		if f(r) {
			if fieldStart >= 0 {
				a[na] = strings.ToLower(s[fieldStart:i])
				na++
				fieldStart = -1
			}
		} else if fieldStart == -1 {
			fieldStart = i
		}
	}
	if fieldStart >= 0 { // Last field might end at EOF.
		a[na] = strings.ToLower(s[fieldStart:])
	}
	return a
}

func InitReduceMapper(){
	WorkingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	WorkingDirectory += "/files/"
	fmt.Println(WorkingDirectory)
}