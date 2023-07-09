package apollo_sdk

import (
	"fmt"
	"reflect"
	"strconv"

	"arclinks-go/library/log"

	"github.com/mitchellh/mapstructure"
	"github.com/zouyx/agollo"
)

var (
	// Fields to hot update listen.
	ApolloListenFields = make(map[string]*Field)
)

type Field struct {
	Name          string
	ApolloKeyName string
	Value         reflect.Value
	Type          reflect.Type
}

func MapApolloConfig(section string, v interface{}) error {
	fields, err := getReflectFields(section, v)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return nil
	}

	err = save(v, fields)
	if err != nil {
		return fmt.Errorf("save err: %v", err)
	}

	for apolloKeyName, field := range fields {
		ApolloListenFields[apolloKeyName] = field
	}

	return nil
}

func getReflectFields(section string, v interface{}) (map[string]*Field, error) {
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Type must be a struct")
	}

	result := make(map[string]*Field)
	fieldCnt := typeOf.NumField()
	for i := 0; i < fieldCnt; i++ {
		name := typeOf.Field(i).Name
		apolloKeyName := fmt.Sprintf("%s.%s", section, name)
		result[apolloKeyName] = &Field{
			Name:          name,
			ApolloKeyName: apolloKeyName,
			Type:          typeOf.Field(i).Type,
			Value:         valueOf.Field(i),
		}
	}

	return result, nil
}

func save(v interface{}, fields map[string]*Field) error {
	configValues := make(map[string]interface{})
	for apolloKeyName, field := range fields {
		var value interface{}
		switch field.Type.Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			value = ApolloServer.GetIntValue(apolloKeyName, 0)
		case reflect.Float32,
			reflect.Float64:
			value = ApolloServer.GetFloatValue(apolloKeyName, 0)
		case reflect.String:
			value = ApolloServer.GetStringValue(apolloKeyName, "")
		case reflect.Bool:
			value = ApolloServer.GetBoolValue(apolloKeyName, false)
		default:
			return fmt.Errorf("Current field type is not be supported")
		}
		configValues[field.Name] = value
	}
	if len(configValues) == 0 {
		return nil
	}

	if err := mapstructure.Decode(configValues, v); err != nil {
		return err
	}

	return nil
}

// TriggerApolloHotUpdateListen triggers hot update listen of apollo.
func TriggerApolloHotUpdateListen(apolloServerUrl, appId, namespaceName string) {
	if len(ApolloListenFields) == 0 {
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				ApolloServer.logger.Errorf("TriggerApolloHotUpdateListen Recover err: %v", err)
			}
		}()

		NewApolloServer(apolloServerUrl)
		apolloLog := &ApolloLogger{Logger: log.BusinessLog}
		err := ApolloServer.SetLog(apolloLog).Start(appId, namespaceName)
		if err != nil {
			ApolloServer.logger.Errorf("InitApolloServer Recover err: %v", err)
			return
		}

		for event := range ApolloServer.ListenChangeEvent() {
			for apolloKeyName, value := range event.Changes {
				if value.ChangeType != agollo.MODIFIED {
					continue
				}
				if field, ok := ApolloListenFields[apolloKeyName]; ok {
					if !field.Value.CanSet() {
						ApolloServer.logger.Errorf("Apollo field can't set, field: %v", field)
						continue
					}
					err := field.SetNewValue(value.NewValue)
					if err != nil {
						ApolloServer.logger.Errorf("Apollo field.SetNewValue err: %v field: %v value: %v", err, field, value)
					}
				}
			}
		}
	}()
}

func (f *Field) SetNewValue(newValue string) error {
	switch f.Type.Kind() {
	case reflect.Int:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(againValue))
	case reflect.Int8:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(int8(againValue)))
	case reflect.Int16:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(int16(againValue)))
	case reflect.Int32:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(int32(againValue)))
	case reflect.Int64:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(int64(againValue)))
	case reflect.Uint:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(uint(againValue)))
	case reflect.Uint8:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(uint8(againValue)))
	case reflect.Uint16:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(uint16(againValue)))
	case reflect.Uint32:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(uint32(againValue)))
	case reflect.Uint64:
		againValue, err := strconv.Atoi(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(uint64(againValue)))
	case reflect.Float32:
		againValue, err := strconv.ParseFloat(newValue, 32)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(againValue))
	case reflect.Float64:
		againValue, err := strconv.ParseFloat(newValue, 64)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(againValue))
	case reflect.String:
		f.Value.Set(reflect.ValueOf(newValue))
	case reflect.Bool:
		againValue, err := strconv.ParseBool(newValue)
		if err != nil {
			return err
		}
		f.Value.Set(reflect.ValueOf(againValue))
	default:
		return fmt.Errorf("Unkonwn field type")
	}

	return nil
}
