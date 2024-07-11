// Code generated from C:/Users/iseki/working/client/module/luarocks/parser/LuaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LuaParser

import "github.com/antlr4-go/antlr/v4"

type BaseLuaParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseLuaParserVisitor) VisitStart_(ctx *Start_Context) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitChunk(ctx *ChunkContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitStat(ctx *StatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitAttnamelist(ctx *AttnamelistContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitAttrib(ctx *AttribContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitRetstat(ctx *RetstatContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitLabel(ctx *LabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitFuncname(ctx *FuncnameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitVarlist(ctx *VarlistContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitNamelist(ctx *NamelistContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitExplist(ctx *ExplistContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitExp(ctx *ExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitVar(ctx *VarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitPrefixexp(ctx *PrefixexpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitFunctioncall(ctx *FunctioncallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitArgs(ctx *ArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitFunctiondef(ctx *FunctiondefContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitFuncbody(ctx *FuncbodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitParlist(ctx *ParlistContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitTableconstructor(ctx *TableconstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitFieldlist(ctx *FieldlistContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitField(ctx *FieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitFieldsep(ctx *FieldsepContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitNumber(ctx *NumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseLuaParserVisitor) VisitString(ctx *StringContext) interface{} {
	return v.VisitChildren(ctx)
}
