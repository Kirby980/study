package Registry

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type TypeRegistry struct {
	Message map[string]string
	types   map[string]ProtoMessage
	counter int
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		Message: make(map[string]string),
		types:   make(map[string]ProtoMessage),
	}
}

func ProcessValueV1(val interface{}, registry *TypeRegistry, suggestedName string) string {
	switch v := val.(type) {
	case []interface{}:
		if len(v) == 0 {
			return "bytes"
		}
		elemType := ProcessValueV1(v[0], registry, suggestedName)
		return "repeated " + elemType
	case map[string]interface{}:
		return ProcessObjectV1(v, registry, suggestedName)
	case json.Number:
		if _, err := v.Int64(); err == nil {
			return "int32"
		}
		return "double"
	case string:
		return "string"
	case bool:
		return "bool"
	case nil:
		return "google.protobuf.NullValue"
	default:
		return "bytes"
	}
}
func ProcessObjectV1(obj map[string]interface{}, registry *TypeRegistry, suggestedName string) string {
	fields := make([]struct{ Name, Type string }, 0)
	for k, v := range obj {
		if _, err := strconv.Atoi(k); err == nil {
			return "map<string,string>"
		}
		fieldType := ProcessValueV1(v, registry, k)
		fields = append(fields, struct{ Name, Type string }{k, fieldType})
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].Name < fields[j].Name })
	var sig strings.Builder
	for _, f := range fields {
		sig.WriteString(f.Name + ":" + f.Type + ";")
	}
	if name, ok := registry.Message[sig.String()]; ok {
		return name
	}
	messageName := generateMessageName(suggestedName, registry)
	protoField := make([]ProtoField, 0, 20)
	tag := 1
	for _, field := range fields {
		protoField = append(protoField, ProtoField{
			Name: strings.ToLower(field.Name),
			Type: field.Type,
			Tag:  tag,
		})
		tag++
	}
	registry.Message[sig.String()] = messageName
	registry.types[messageName] = ProtoMessage{
		Name:   messageName,
		Fields: protoField,
	}
	return messageName
}

func ProcessValueV2(val *OrderedValue, registry *TypeRegistry, suggestedName string) string {
	switch val.Type {
	case 1: // Object
		return processObjectV2(val.Object, registry, suggestedName)
	case 2: // Array
		if len(val.Array) == 0 {
			return "bytes"
		}
		elemType := ProcessValueV2(val.Array[0], registry, suggestedName)
		return "repeated " + elemType
	default: // Primitive
		return processPrimitive(val.Primitive)
	}
}
func processPrimitive(val interface{}) string {
	switch v := val.(type) {
	case json.Number:
		if _, err := v.Int64(); err == nil {
			return "int32"
		}
		return "double"
	case string:
		return "string"
	case bool:
		return "bool"
	case nil:
		return "google.protobuf.NullValue"
	default:
		return "bytes"
	}
}
func processObjectV2(obj *OrderedObject, registry *TypeRegistry, suggestedName string) string {
	fields := make([]struct{ Name, Type string }, len(obj.Keys))

	// 按原始顺序处理每个字段
	for i, key := range obj.Keys {
		value := obj.Data[key]
		if _, err := strconv.Atoi(key); err == nil {
			return "map<string, string>"
		}
		fields[i] = struct{ Name, Type string }{
			Name: key,
			Type: ProcessValueV2(value, registry, key),
		}
	}

	// 生成排序后的签名
	sortedFields := make([]struct{ Name, Type string }, len(fields))
	copy(sortedFields, fields)
	sort.Slice(sortedFields, func(i, j int) bool {
		return sortedFields[i].Name < sortedFields[j].Name
	})

	var sig strings.Builder
	for _, f := range sortedFields {
		sig.WriteString(f.Name + ":" + f.Type + ";")
	}

	// 类型复用检查
	if name, ok := registry.Message[sig.String()]; ok {
		return name
	}

	// 生成新类型
	messageName := generateMessageName(suggestedName, registry)
	protoFields := make([]ProtoField, len(fields))
	for i, field := range fields {
		protoFields[i] = ProtoField{
			Name: strings.ToLower(field.Name),
			Type: field.Type,
			Tag:  i + 1,
		}
	}

	registry.Message[sig.String()] = messageName
	registry.types[messageName] = ProtoMessage{
		Name:   messageName,
		Fields: protoFields,
	}

	return messageName
}

func generateMessageName(base string, r *TypeRegistry) string {
	name := strings.Title(base)
	if _, exists := r.types[name]; !exists {
		return name
	}
	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s%d", name, i)
		if _, exists := r.types[newName]; !exists {
			return newName
		}
	}
}
