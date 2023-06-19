package script

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
)

type LuaScriptFuncDebuggerHandler struct {
	BaseScripHandler
	function string
}

func NewCommonFuncDebuggerHandler() (handler *LuaScriptFuncDebuggerHandler) {
	handler = &LuaScriptFuncDebuggerHandler{}
	handler.Name = enum.LuaFuncName_DoCommonFunctionExecute
	handler.FuncType = enum.LuaFuncType_DoFuncExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return handler
}

func (l *LuaScriptFuncDebuggerHandler) Init() (err error) {
	l.function, err = loadScript(enum.LuaFuncName_DoCommonFunctionExecute)
	return err
}

func (l LuaScriptFuncDebuggerHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_DoCommonFunctionExecute)
}

func (l LuaScriptFuncDebuggerHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) (err error) {
	return scriptExecute(enum.LuaFuncName_DoCommonFunctionExecute, execCtx, funcCtx)
}
