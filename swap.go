package common

import "encoding/json"

// 将 from类型转换进成to类型 前提是 from/to 拥有相同的字段
func SwapTo(from, to interface{}) error {
	dataByte, err := json.Marshal(from)
	if err != nil {
		return err
	}

	return json.Unmarshal(dataByte, to)
}
