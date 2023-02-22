package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

func ConvertBodyToMap(reqBody io.ReadCloser) (map[string]interface{}, error) {
	var checking map[string]interface{}
	body, errBody := ioutil.ReadAll(reqBody)
	if errBody != nil {
		fmt.Println("err errBody != nil")
		return nil, errBody
	}
	if err := json.Unmarshal([]byte(body), &checking); err != nil {
		fmt.Println("err ConvertBodyToMap ", err.Error())
		return nil, err
	}
	return checking, nil
}

func CheckingSettingRequest(listString []string, listMap map[string]interface{}) (bool, interface{}) {
	fmt.Println(listMap)
	for _, element := range listString {
		// element is the element from someSlice for where we are
		if _, ok := listMap[element]; !ok {
			//do something here
			return false, element
		}
	}
	return true, "Success"
}

func FormDataToInteger(data interface{}) (int, bool) {
	dataString, ok := data.([]string)
	if !ok {
		PrintCommand("FormDataToInteger err ok => ", ok)
		return -1, false
	}
	resultData, err := strconv.Atoi(dataString[0])
	if err != nil {
		PrintCommand("FormDataToInteger err resultData => ", resultData, "- Data => ", data)
		return -1, false
	}
	return resultData, true
}

func FormDataToFloat64(data interface{}) (float64, bool) {
	dataString, ok := data.([]string)
	if !ok {
		return -1, false
	}
	resultData, err := strconv.ParseFloat(dataString[0], 64)
	if err != nil {
		return -1, false
	}
	return resultData, true
}

func FormDataToBool(data interface{}) bool {
	dataString, ok := data.([]string)
	if !ok {
		fmt.Println("is_ecc ", ok)
		return false
	}
	resultData, err := strconv.ParseBool(dataString[0])
	if err != nil {
		fmt.Println("is_ecc pars ", dataString)
		return false
	}
	return resultData
}
