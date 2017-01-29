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
		//fmt.Println("CELL: ", redisCells[i]);
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

	if len(Rows) == 0 || len(Columns) == 0{
		fmt.Println("Redis columns or rows had zero value - making default table.")
		defaultTableCalc()
		InitRedisColsAndRows()
	}

}

func defaultTableCalc(){
	keysToDelete := []string{"columns", "rows"};
	RedisClient.Del(keysToDelete...);

	defaultRowsAndColumns := []string{"0","1","2"}
	redisDefaultRnC := make([]interface{}, len(defaultRowsAndColumns));
	for i, v := range defaultRowsAndColumns{
		redisDefaultRnC[i] = v;
	}
	RedisClient.RPush("columns", redisDefaultRnC...);
	RedisClient.RPush("rows", redisDefaultRnC...);

	valueToSet := []string{
		"CELL_0_0", "", "CELL_0_1", "", "CELL_0_2", "",
		"CELL_1_0", "", "CELL_1_1", "", "CELL_1_2", "",
		"CELL_2_0", "", "CELL_2_1", "", "CELL_2_2", "",
	};
	redisCellValues := make([]interface{}, len(valueToSet));
	for i, v := range valueToSet{
		redisCellValues[i] = v;
	}

	RedisClient.MSet(redisCellValues...);
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

func RemoveFromList(w http.ResponseWriter, req *http.Request){
	statusCode := http.StatusOK;
	responseMessage := "Ok";

	type ListName struct {
		Name	string	`json:"name,omitempty"`
	};

	var redisList ListName;
	result := json.NewDecoder(req.Body).Decode(&redisList);

	if result != nil {
		responseMessage = "Error decoding data.";
		statusCode = http.StatusBadRequest;
	}

	if redisList.Name != "columns" && redisList.Name != "rows"{
		responseMessage = "Unknown list name";
		statusCode = http.StatusBadRequest;
	} else{
		RedisClient.RPop(redisList.Name)
		if redisList.Name == "columns"{
			Columns = Columns[:len(Columns)-1];
		} else{
			Rows = Rows[:len(Rows)-1];
		}
	}

	fmt.Println("Columns: ", Columns);
	fmt.Println("Rows: ", Rows);

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(responseMessage);
}

func AddToList(w http.ResponseWriter, req *http.Request){
	statusCode := http.StatusOK;
	responseMessage := "Ok";

	type ListName struct {
		Name	string	`json:"name,omitempty"`
		Value	string	`json:"value,omitempty"`
	};

	var redisList ListName;
	result := json.NewDecoder(req.Body).Decode(&redisList);

	if result != nil {
		responseMessage = "Error decoding data.";
		statusCode = http.StatusBadRequest;
	}

	if redisList.Name != "columns" && redisList.Name != "rows"{
		responseMessage = "Unknown list name";
		statusCode = http.StatusBadRequest;
	} else{
		redisValue := make([]interface{}, 1);
		redisValue[0] = redisList.Value;
		RedisClient.RPush(redisList.Name, redisValue...);
		if redisList.Name == "columns"{
			Columns = append(Columns, redisList.Value)
		} else{
			Rows = append(Rows, redisList.Value)
		}
	}

	fmt.Println("Columns: ", Columns);
	fmt.Println("Rows: ", Rows);

	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	w.WriteHeader(statusCode);
	json.NewEncoder(w).Encode(responseMessage);
}


