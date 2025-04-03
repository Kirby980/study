package Registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
)

type OrderedValue struct {
	Type      int             // 0: primitive, 1: object, 2: array
	Primitive interface{}     // 基础类型值
	Object    *OrderedObject  // 对象类型值
	Array     []*OrderedValue // 数组类型值
}
type OrderedObject struct {
	Keys []string // 保留键顺序
	Data map[string]*OrderedValue
}

// RemoveJSONComment 去除注释格式化json
func RemoveJSONComment(jsonStr string) string {
	reMultiline := regexp.MustCompile(`(?s)/\*.*?\*/`)
	cleaned := reMultiline.ReplaceAllString(jsonStr, "")

	// 单行注释处理（排除合法URL）
	reSingleLine := regexp.MustCompile(`(?m)(^|[^:])//.*$`)
	cleaned = reSingleLine.ReplaceAllString(cleaned, "$1")

	cleaned = regexp.MustCompile(`,(\s*[}\]])`).ReplaceAllString(cleaned, "$1")
	return cleaned
}

// ParseJson 解析json字段，同时按照json字段顺序存储
func ParseJson(jsonStr string) (*OrderedValue, error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
	decoder.UseNumber()
	return parseValue(decoder)
}
func parseValue(decoder *json.Decoder) (*OrderedValue, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			return parseObject(decoder)
		case '[':
			return parseArray(decoder)
		default:
			return nil, errors.New("unexpected delimiter")
		}
	case json.Number, string, bool, nil:
		return &OrderedValue{Type: 0, Primitive: t}, nil
	default:
		return nil, errors.New("unexpected delimiter")
	}
}
func parseObject(decoder *json.Decoder) (*OrderedValue, error) {
	obj := &OrderedObject{
		Keys: make([]string, 0),
		Data: make(map[string]*OrderedValue),
	}

	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		keyStr := key.(string)
		value, err := parseValue(decoder)
		if err != nil {
			return nil, err
		}

		obj.Keys = append(obj.Keys, keyStr)
		obj.Data[keyStr] = value
	}

	// 消耗结束的'}'
	decoder.Token()
	return &OrderedValue{Type: 1, Object: obj}, nil
}

func parseArray(decoder *json.Decoder) (*OrderedValue, error) {
	var arr []*OrderedValue
	for decoder.More() {
		value, err := parseValue(decoder)
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}
	// 消耗结束的']'
	decoder.Token()
	return &OrderedValue{Type: 2, Array: arr}, nil
}
