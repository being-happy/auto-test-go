package script

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
)

type CommonFuncDebuggerHandler struct {
	BaseScripHandler
	function string
}

func NewCommonFuncDebuggerHandler() (handler *LuaScriptBaseHandler, err error) {
	handler = &LuaScriptBaseHandler{}
	handler.Name = enum.LuaFuncName_DoCommonFunctionExecute
	handler.FuncType = enum.LuaFuncType_DoFuncExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	err = handler.Init()
	return handler, err
}

func (l *CommonFuncDebuggerHandler) Init() (err error) {
	l.function, err = loadScript(enum.LuaFuncName_DoCommonFunctionExecute)
	return err
}

func (l CommonFuncDebuggerHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_DoCommonFunctionExecute)
}

func (l CommonFuncDebuggerHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) (err error) {
	return scriptExecute(enum.LuaFuncName_DoCommonFunctionExecute, execCtx, funcCtx)
}
