package datatype

// 项目配置

// type NullTime string

// 查询的时候解析成字符串
// func (n *NullTime) Scan(value interface{}) error {
// 	bytes, ok := value.([]byte)
// 	if !ok {
// 		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
// 	}
// 	if len(bytes) != 0 {
// 		n = &string(bytes)
// 	}
// 	return nil
// }

// // 保存的时候转换成datatime格式
// func (n NullTime) Value() (driver.Value, error) {
// 	str, err := json.Marshal(j)
// 	if err != nil {
// 		return string(str), err
// 	}
// 	return string(str), nil
// }
