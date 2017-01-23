package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"fmt"
)

type TableCalcCell struct {
	Value	string	`json:"value,omitempty"`
	Column	int	`json:"column,omitempty"`
	Row	int	`json:"row,omitempty"`
	RedisKey string	`json:"redisKey,omitempty"`
}

var Columns []string;
var Rows []string;

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

	var redisCells []TableCalcCell;
	result := json.NewDecoder(req.Body).Decode(&redisCells);
	if result != nil {
		responseMessage = "Error decoding data.";
		statusCode = http.StatusBadRequest;
	}

	for i:=0;i < len(redisCells); i++{
		fmt.Println("CELL: ", redisCells[i]);
		err := RedisClient.Set(redisCells[i].RedisKey, redisCells[i].Value, 0).Err()
		if err != nil {
			panic(err);
			responseMessage = "Error writing to redis storage.";
			statusCode = http.StatusBadRequest;
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(responseMessage);
}

//Used for setting initial rows and columns as global variables.
func InitRedisColsAndRows(){

	columns, err := RedisClient.LRange("columns",0,-1).Result();
	if err != nil {
		panic(err);
	}
	Columns = columns;

	rows, err := RedisClient.LRange("rows",0,-1).Result();
	if err != nil {
		panic(err);
	}
	Rows = rows;
}

//We are writing values in redis as the CELL_<ROW>_<COLUMN>
func GetRedisData(w http.ResponseWriter, req *http.Request){
	statusCode := http.StatusOK;
	//responseMessage := "Ok";

	numRows := len(Rows);
	numColumns := len(Columns);

	responseValues := make([][]string, numRows);

	for i:=0; i < numRows; i++{
		redisQueryKeys := []string{};
		stringRow := strconv.Itoa(i);

		for j:=0; j < numColumns; j++{
			redisQueryKeys = append(redisQueryKeys, "CELL_" + stringRow + "_" + strconv.Itoa(j));
		}

		result := RedisClient.MGet(redisQueryKeys...);
		if result.Err() != nil {
			panic(result.Err());
		}
		values := result.Val();
		responseValues[i] = make([]string, numColumns);

		for j:=0; j < len(values); j++{
			if values[j] == nil{
				responseValues[i][j] = "";
			} else{
				responseValues[i][j] = values[j].(string);
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(responseValues);
}