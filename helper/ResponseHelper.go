package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
