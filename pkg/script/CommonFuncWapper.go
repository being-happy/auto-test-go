package script

import (
	"auto-test-go/pkg/entities"
	"fmt"
)

type CommonFuncWapper struct {
}

func (CommonFuncWapper) wrap(functions []string, execCtx *entities.ExecContext) string {
	funcScripts := ""
	if functions != nil && len(functions) > 0 {
		for _, f := range functions {
			log := fmt.Sprintf("[CaseScriptHandleRegister] Begin to add common function : %s", f)
			if funcScripts != "" {
				funcScripts = funcScripts + "\n" + f
			} else {
				funcScripts = f
			}
			execCtx.AddLogs(log)
		}
	}
	return funcScripts
}
