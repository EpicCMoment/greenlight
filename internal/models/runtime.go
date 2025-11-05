package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil

}

func (r *Runtime) UnmarshalJSON(data []byte) error {

	strData := string(data)

	strData, err := strconv.Unquote(strData)

	if err != nil {
		return err
	}

	runtimeString := strings.Split(strData, " ")[0]

	val, err := strconv.ParseInt(runtimeString, 10, 32)

	if err != nil {
		return err
	}

	*r = Runtime(val)

	return nil

}
