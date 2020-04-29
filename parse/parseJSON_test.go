package parse

import (
	"reflect"
	"testing"
)

type testInstance struct {
	path string
	result interface{}
	err error
}

type base struct {
	file string
	result interface{}
	err error
}

type tests struct {
	b base
	instances []testInstance
}

var test = []tests{
	{
		b:base{file:"./test/valid/test.json", err:nil},
		instances:[]testInstance{
			{"1int", 0x123f, nil},
			{"1string", "string", nil},
			{"1string", "string", nil},
			{"1object.2int1", 1234, nil},
			{"1object.2string1\n", "1234", nil},
			{"1object.2object1.3int1", 333, nil},
			{"1object.2object2.string", "23 /.,mb,np	!~ +-_'string\n", nil},
		},
	},
}

func TestLoad(t *testing.T)  {

	for i := 0; i < len(test); i++ {
		json := NewParseJSON(test[i].b.file)
		AssertNotNil(t, json)

		err := json.Load()
		AssertEqual(t, err, test[i].b.err, "")
		//jsonPrint(json.root, t)

		instances := test[i].instances
		for j := 0; j < len(instances); j++ {
			Type := reflect.TypeOf(instances[j].result).Kind()
			if Type == reflect.Int {
				v, err := json.GetIntValue(instances[j].path)
				AssertEqual(t, v, instances[j].result, instances[j].path)
				AssertEqual(t, err, instances[j].err, instances[j].path)
			} else if Type == reflect.String {
				v, err := json.GetStringValue(instances[j].path)
				AssertEqual(t, v, instances[j].result, instances[j].path)
				AssertEqual(t, err, instances[j].err, instances[j].path)
			}
		}
	}
}

func PrintJsonValue(json *jsonValue, t *testing.T, path string)  {
	fields := json.Value.(*jsonValueObject).fields
	for k, v := range fields {
		key := string(k.Value.(jsonValueString))
		if v.Type == Int {
			value := int(*v.Value.(*jsonValueInt))
			t.Errorf("key: %s int: %d", path+key, value)
		} else if v.Type == String {
			value := string(*v.Value.(*jsonValueString))
			t.Errorf("key: %s string: %s", path+key, value)
		} else if v.Type == Object {
			PrintJsonValue(v, t, path + key + ".")
		}
	}
}

func jsonPrint(json *jsonValue, t *testing.T)  {
	t.Helper()

	PrintJsonValue(json, t, "")
}