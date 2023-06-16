package script

import (
	"auto-test-go/pkg/entities"
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

type CommonFuncWapper struct {
}

func (CommonFuncWapper) wrap(state *lua.LState, functions []string, execCtx *entities.ExecContext) error {
	if functions != nil && len(functions) > 0 {
		for _, f := range functions {
			err := state.DoString(f)
			log := fmt.Sprintf("[CaseScriptHandleRegister] Begin to add common function : %s", f)
			execCtx.AddLogs(log)
			if err != nil {
				log = fmt.Sprintf("[CaseScriptHandleRegister] Add common function to script error, script name: %s, error: %s", f, err.Error())
				execCtx.AddLogs(log)
				return err
			}
		}
	}
	return nil
}
