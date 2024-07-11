// Code generated from C:/Users/iseki/working/client/module/luarocks/parser/LuaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LuaParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by LuaParser.
type LuaParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by LuaParser#start_.
	VisitStart_(ctx *Start_Context) interface{}

	// Visit a parse tree produced by LuaParser#chunk.
	VisitChunk(ctx *ChunkContext) interface{}

	// Visit a parse tree produced by LuaParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by LuaParser#stat.
	VisitStat(ctx *StatContext) interface{}

	// Visit a parse tree produced by LuaParser#attnamelist.
	VisitAttnamelist(ctx *AttnamelistContext) interface{}

	// Visit a parse tree produced by LuaParser#attrib.
	VisitAttrib(ctx *AttribContext) interface{}

	// Visit a parse tree produced by LuaParser#retstat.
	VisitRetstat(ctx *RetstatContext) interface{}

	// Visit a parse tree produced by LuaParser#label.
	VisitLabel(ctx *LabelContext) interface{}

	// Visit a parse tree produced by LuaParser#funcname.
	VisitFuncname(ctx *FuncnameContext) interface{}

	// Visit a parse tree produced by LuaParser#varlist.
	VisitVarlist(ctx *VarlistContext) interface{}

	// Visit a parse tree produced by LuaParser#namelist.
	VisitNamelist(ctx *NamelistContext) interface{}

	// Visit a parse tree produced by LuaParser#explist.
	VisitExplist(ctx *ExplistContext) interface{}

	// Visit a parse tree produced by LuaParser#exp.
	VisitExp(ctx *ExpContext) interface{}

	// Visit a parse tree produced by LuaParser#var.
	VisitVar(ctx *VarContext) interface{}

	// Visit a parse tree produced by LuaParser#prefixexp.
	VisitPrefixexp(ctx *PrefixexpContext) interface{}

	// Visit a parse tree produced by LuaParser#functioncall.
	VisitFunctioncall(ctx *FunctioncallContext) interface{}

	// Visit a parse tree produced by LuaParser#args.
	VisitArgs(ctx *ArgsContext) interface{}

	// Visit a parse tree produced by LuaParser#functiondef.
	VisitFunctiondef(ctx *FunctiondefContext) interface{}

	// Visit a parse tree produced by LuaParser#funcbody.
	VisitFuncbody(ctx *FuncbodyContext) interface{}

	// Visit a parse tree produced by LuaParser#parlist.
	VisitParlist(ctx *ParlistContext) interface{}

	// Visit a parse tree produced by LuaParser#tableconstructor.
	VisitTableconstructor(ctx *TableconstructorContext) interface{}

	// Visit a parse tree produced by LuaParser#fieldlist.
	VisitFieldlist(ctx *FieldlistContext) interface{}

	// Visit a parse tree produced by LuaParser#field.
	VisitField(ctx *FieldContext) interface{}

	// Visit a parse tree produced by LuaParser#fieldsep.
	VisitFieldsep(ctx *FieldsepContext) interface{}

	// Visit a parse tree produced by LuaParser#number.
	VisitNumber(ctx *NumberContext) interface{}

	// Visit a parse tree produced by LuaParser#string.
	VisitString(ctx *StringContext) interface{}
}
