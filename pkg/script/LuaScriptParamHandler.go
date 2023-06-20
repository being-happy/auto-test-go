package script

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
)

type LuaScriptParamHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptDoParamHandler() *LuaScriptParamHandler {
	handler := LuaScriptParamHandler{}
	handler.Name = enum.LuaFuncType_DoParamExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	handler.FuncType = enum.LuaFuncType_DoParamExecute
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return &handler
}

func (l *LuaScriptParamHandler) Init() (err error) {
	l.function, err = loadScript(enum.LuaFuncName_DoParamExecute)
	return err
}

func (l LuaScriptParamHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_DoParamExecute)
}

func (l LuaScriptParamHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) (err error) {
	return scriptExecute(enum.LuaFuncName_DoParamExecute, execCtx, funcCtx)
}
