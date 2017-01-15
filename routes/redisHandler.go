package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
)

type FrontEndValue struct {
	Value	string	`json:"value,omitempty"`
}

func GetRedisValueByKey(w http.ResponseWriter, req *http.Request){
	statusCode := http.StatusOK;

	params := mux.Vars(req);
	key := params["key"];

	val, err := RedisClient.Get(key).Result()

	if err != nil {
		//panic(err);
		statusCode = http.StatusBadRequest;
		val = "Value not found";
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(val);

}

func SetRedisValueByKey(w http.ResponseWriter, req *http.Request){
	statusCode := http.StatusOK;
	responseMessage := "Ok";

	params := mux.Vars(req);
	key := params["key"];

	var frontend FrontEndValue;

	result := json.NewDecoder(req.Body).Decode(&frontend);
	fmt.Println("result: ", result);
	fmt.Println("req.Body: ", req.Body);
	if result != nil {
		responseMessage = "Error decoding data.";
		statusCode = http.StatusBadRequest;
	}

	err := RedisClient.Set(key, frontend.Value, 0).Err()
	if err != nil {
		panic(err);
		responseMessage = "Error writing to redis storage.";
		statusCode = http.StatusBadRequest;
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(responseMessage);
}