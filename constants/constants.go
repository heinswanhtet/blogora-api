package constants

import (
	"crypto/sha256"
	"reflect"

	"github.com/heinswanhtet/blogora-api/configs"
	"github.com/heinswanhtet/blogora-api/types"
)

const ContextData types.ContextKey = "data"

var AuthorFields []string = getJSONFieldNames(&types.Author{})

var SECRET []byte = func() []byte {
	secret := configs.Envs.SECRET
	h := sha256.New()
	_, err := h.Write([]byte(secret))
	if err != nil {
		panic("SHA256 SECRET failed!")
	}
	sha256SecretByte := h.Sum(nil)
	return sha256SecretByte
}()

func getJSONFieldNames(i any) []string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // to work even if i is passed as ptr
	}
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			fields = append(fields, tag)
		}
	}
	return fields
}
