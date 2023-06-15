package script

import (
	"auto-test-go/pkg/entities"
	"github.com/valyala/fasttemplate"
	"strings"
)

type FastTemplate struct {
}

func (FastTemplate) convertVar(vars map[string]entities.VarValue) map[string]interface{} {
	newVars := map[string]interface{}{}
	for k, v := range vars {
		newVars[k] = v.Value
	}
	return newVars
}

func (FastTemplate) template(source string, vars map[string]interface{}) string {
	t := fasttemplate.New(source, "{@", "}")
	source = t.ExecuteString(vars)
	for keyWord, value := range vars {
		arg1 := "@" + keyWord
		//兼容老的替换方式
		if value == nil {
			value = ""
		}
		if strings.Contains(source, arg1) {
			source = strings.Replace(source, arg1, value.(string), -1)
		}
	}
	return source
}
