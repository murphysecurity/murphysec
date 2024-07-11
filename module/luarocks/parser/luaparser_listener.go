// Code generated from C:/Users/iseki/working/client/module/luarocks/parser/LuaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LuaParser

import "github.com/antlr4-go/antlr/v4"

// LuaParserListener is a complete listener for a parse tree produced by LuaParser.
type LuaParserListener interface {
	antlr.ParseTreeListener

	// EnterStart_ is called when entering the start_ production.
	EnterStart_(c *Start_Context)

	// EnterChunk is called when entering the chunk production.
	EnterChunk(c *ChunkContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterStat is called when entering the stat production.
	EnterStat(c *StatContext)

	// EnterAttnamelist is called when entering the attnamelist production.
	EnterAttnamelist(c *AttnamelistContext)

	// EnterAttrib is called when entering the attrib production.
	EnterAttrib(c *AttribContext)

	// EnterRetstat is called when entering the retstat production.
	EnterRetstat(c *RetstatContext)

	// EnterLabel is called when entering the label production.
	EnterLabel(c *LabelContext)

	// EnterFuncname is called when entering the funcname production.
	EnterFuncname(c *FuncnameContext)

	// EnterVarlist is called when entering the varlist production.
	EnterVarlist(c *VarlistContext)

	// EnterNamelist is called when entering the namelist production.
	EnterNamelist(c *NamelistContext)

	// EnterExplist is called when entering the explist production.
	EnterExplist(c *ExplistContext)

	// EnterExp is called when entering the exp production.
	EnterExp(c *ExpContext)

	// EnterVar is called when entering the var production.
	EnterVar(c *VarContext)

	// EnterPrefixexp is called when entering the prefixexp production.
	EnterPrefixexp(c *PrefixexpContext)

	// EnterFunctioncall is called when entering the functioncall production.
	EnterFunctioncall(c *FunctioncallContext)

	// EnterArgs is called when entering the args production.
	EnterArgs(c *ArgsContext)

	// EnterFunctiondef is called when entering the functiondef production.
	EnterFunctiondef(c *FunctiondefContext)

	// EnterFuncbody is called when entering the funcbody production.
	EnterFuncbody(c *FuncbodyContext)

	// EnterParlist is called when entering the parlist production.
	EnterParlist(c *ParlistContext)

	// EnterTableconstructor is called when entering the tableconstructor production.
	EnterTableconstructor(c *TableconstructorContext)

	// EnterFieldlist is called when entering the fieldlist production.
	EnterFieldlist(c *FieldlistContext)

	// EnterField is called when entering the field production.
	EnterField(c *FieldContext)

	// EnterFieldsep is called when entering the fieldsep production.
	EnterFieldsep(c *FieldsepContext)

	// EnterNumber is called when entering the number production.
	EnterNumber(c *NumberContext)

	// EnterString is called when entering the string production.
	EnterString(c *StringContext)

	// ExitStart_ is called when exiting the start_ production.
	ExitStart_(c *Start_Context)

	// ExitChunk is called when exiting the chunk production.
	ExitChunk(c *ChunkContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitStat is called when exiting the stat production.
	ExitStat(c *StatContext)

	// ExitAttnamelist is called when exiting the attnamelist production.
	ExitAttnamelist(c *AttnamelistContext)

	// ExitAttrib is called when exiting the attrib production.
	ExitAttrib(c *AttribContext)

	// ExitRetstat is called when exiting the retstat production.
	ExitRetstat(c *RetstatContext)

	// ExitLabel is called when exiting the label production.
	ExitLabel(c *LabelContext)

	// ExitFuncname is called when exiting the funcname production.
	ExitFuncname(c *FuncnameContext)

	// ExitVarlist is called when exiting the varlist production.
	ExitVarlist(c *VarlistContext)

	// ExitNamelist is called when exiting the namelist production.
	ExitNamelist(c *NamelistContext)

	// ExitExplist is called when exiting the explist production.
	ExitExplist(c *ExplistContext)

	// ExitExp is called when exiting the exp production.
	ExitExp(c *ExpContext)

	// ExitVar is called when exiting the var production.
	ExitVar(c *VarContext)

	// ExitPrefixexp is called when exiting the prefixexp production.
	ExitPrefixexp(c *PrefixexpContext)

	// ExitFunctioncall is called when exiting the functioncall production.
	ExitFunctioncall(c *FunctioncallContext)

	// ExitArgs is called when exiting the args production.
	ExitArgs(c *ArgsContext)

	// ExitFunctiondef is called when exiting the functiondef production.
	ExitFunctiondef(c *FunctiondefContext)

	// ExitFuncbody is called when exiting the funcbody production.
	ExitFuncbody(c *FuncbodyContext)

	// ExitParlist is called when exiting the parlist production.
	ExitParlist(c *ParlistContext)

	// ExitTableconstructor is called when exiting the tableconstructor production.
	ExitTableconstructor(c *TableconstructorContext)

	// ExitFieldlist is called when exiting the fieldlist production.
	ExitFieldlist(c *FieldlistContext)

	// ExitField is called when exiting the field production.
	ExitField(c *FieldContext)

	// ExitFieldsep is called when exiting the fieldsep production.
	ExitFieldsep(c *FieldsepContext)

	// ExitNumber is called when exiting the number production.
	ExitNumber(c *NumberContext)

	// ExitString is called when exiting the string production.
	ExitString(c *StringContext)
}
