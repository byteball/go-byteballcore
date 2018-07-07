package json

import(
	"encoding/json"
	"log"
)

type(
)

func Stringify(obj interface{}) string {
//	panic("[tbd] json.Stringify")
//	bs, err := json.Marshal(obj)
	bs, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf("json.Stringify: %s", err.Error())
	}
	return string(bs)
}
