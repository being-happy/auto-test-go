package script

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
)

type LuaScriptJudgmentHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptJudgmentHandler() *LuaScriptJudgmentHandler {
	handler := LuaScriptJudgmentHandler{}
	handler.Name = enum.LuaFuncName_DoJudgmentExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	handler.FuncType = enum.LuaFuncType_DoJudgmentExecute
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return &handler
}

func (l *LuaScriptJudgmentHandler) Init() error {
	body, err := loadScript(enum.LuaFuncName_DoJudgmentExecute)
	l.function = body
	return err
}

func (l LuaScriptJudgmentHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return buildScript(execCtx, funcCtx, l.function, enum.LuaFuncName_DoJudgmentExecute)
}

func (l LuaScriptJudgmentHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) (err error) {
	return scriptExecute(enum.LuaFuncName_DoJudgmentExecute, execCtx, funcCtx)
}
