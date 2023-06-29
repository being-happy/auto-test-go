package script

import (
	"auto-test-go/pkg/entities"
	"auto-test-go/pkg/enum"
	"errors"
	"fmt"
	"strings"
)

type LuaScriptRedisHandler struct {
	BaseScripHandler
	function string
}

func NewLuaScriptRedisHandler() *LuaScriptRedisHandler {
	handler := LuaScriptRedisHandler{}
	handler.Name = enum.LuaFuncType_DoRedisExecute
	handler.ScriptType = enum.ScriptType_LuaScript
	handler.FuncType = enum.LuaFuncType_DoRedisExecute
	err := handler.Init()
	if err != nil {
		panic(err)
	}
	return &handler
}

func (l *LuaScriptRedisHandler) Init() error {
	body, err := loadScript(enum.LuaFuncName_DoRedisExecute)
	l.function = body
	return err
}

func (l *LuaScriptRedisHandler) BuildScript(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	if funcCtx.FuncBody == "" || funcCtx.DbName == "" || funcCtx.Password == "" || funcCtx.Host == "" || funcCtx.Port == "" {
		log := fmt.Sprintf("[LuaScriptRedisHandler] One or more of userName,password,host,port is nil!")
		execCtx.AddLogs(log)
		return errors.New(log)
	}

	funcScript := strings.Replace(l.function, "@host", funcCtx.Host, -1)
	funcScript = strings.Replace(funcScript, "@port", funcCtx.Port, -1)
	funcScript = strings.Replace(funcScript, "@dbName", funcCtx.DbName, -1)
	funcScript = strings.Replace(funcScript, "@password", funcCtx.Password, -1)
	return buildScript(execCtx, funcCtx, funcScript, enum.LuaFuncName_DoRedisExecute)
}

func (l *LuaScriptRedisHandler) Execute(execCtx *entities.ExecContext, funcCtx *entities.FuncContext) error {
	return scriptExecute(enum.LuaFuncName_DoRedisExecute, execCtx, funcCtx)
}
