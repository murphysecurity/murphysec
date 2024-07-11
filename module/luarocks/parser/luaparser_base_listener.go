// Code generated from C:/Users/iseki/working/client/module/luarocks/parser/LuaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LuaParser

import "github.com/antlr4-go/antlr/v4"

// BaseLuaParserListener is a complete listener for a parse tree produced by LuaParser.
type BaseLuaParserListener struct{}

var _ LuaParserListener = &BaseLuaParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseLuaParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseLuaParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseLuaParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseLuaParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart_ is called when production start_ is entered.
func (s *BaseLuaParserListener) EnterStart_(ctx *Start_Context) {}

// ExitStart_ is called when production start_ is exited.
func (s *BaseLuaParserListener) ExitStart_(ctx *Start_Context) {}

// EnterChunk is called when production chunk is entered.
func (s *BaseLuaParserListener) EnterChunk(ctx *ChunkContext) {}

// ExitChunk is called when production chunk is exited.
func (s *BaseLuaParserListener) ExitChunk(ctx *ChunkContext) {}

// EnterBlock is called when production block is entered.
func (s *BaseLuaParserListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BaseLuaParserListener) ExitBlock(ctx *BlockContext) {}

// EnterStat is called when production stat is entered.
func (s *BaseLuaParserListener) EnterStat(ctx *StatContext) {}

// ExitStat is called when production stat is exited.
func (s *BaseLuaParserListener) ExitStat(ctx *StatContext) {}

// EnterAttnamelist is called when production attnamelist is entered.
func (s *BaseLuaParserListener) EnterAttnamelist(ctx *AttnamelistContext) {}

// ExitAttnamelist is called when production attnamelist is exited.
func (s *BaseLuaParserListener) ExitAttnamelist(ctx *AttnamelistContext) {}

// EnterAttrib is called when production attrib is entered.
func (s *BaseLuaParserListener) EnterAttrib(ctx *AttribContext) {}

// ExitAttrib is called when production attrib is exited.
func (s *BaseLuaParserListener) ExitAttrib(ctx *AttribContext) {}

// EnterRetstat is called when production retstat is entered.
func (s *BaseLuaParserListener) EnterRetstat(ctx *RetstatContext) {}

// ExitRetstat is called when production retstat is exited.
func (s *BaseLuaParserListener) ExitRetstat(ctx *RetstatContext) {}

// EnterLabel is called when production label is entered.
func (s *BaseLuaParserListener) EnterLabel(ctx *LabelContext) {}

// ExitLabel is called when production label is exited.
func (s *BaseLuaParserListener) ExitLabel(ctx *LabelContext) {}

// EnterFuncname is called when production funcname is entered.
func (s *BaseLuaParserListener) EnterFuncname(ctx *FuncnameContext) {}

// ExitFuncname is called when production funcname is exited.
func (s *BaseLuaParserListener) ExitFuncname(ctx *FuncnameContext) {}

// EnterVarlist is called when production varlist is entered.
func (s *BaseLuaParserListener) EnterVarlist(ctx *VarlistContext) {}

// ExitVarlist is called when production varlist is exited.
func (s *BaseLuaParserListener) ExitVarlist(ctx *VarlistContext) {}

// EnterNamelist is called when production namelist is entered.
func (s *BaseLuaParserListener) EnterNamelist(ctx *NamelistContext) {}

// ExitNamelist is called when production namelist is exited.
func (s *BaseLuaParserListener) ExitNamelist(ctx *NamelistContext) {}

// EnterExplist is called when production explist is entered.
func (s *BaseLuaParserListener) EnterExplist(ctx *ExplistContext) {}

// ExitExplist is called when production explist is exited.
func (s *BaseLuaParserListener) ExitExplist(ctx *ExplistContext) {}

// EnterExp is called when production exp is entered.
func (s *BaseLuaParserListener) EnterExp(ctx *ExpContext) {}

// ExitExp is called when production exp is exited.
func (s *BaseLuaParserListener) ExitExp(ctx *ExpContext) {}

// EnterVar is called when production var is entered.
func (s *BaseLuaParserListener) EnterVar(ctx *VarContext) {}

// ExitVar is called when production var is exited.
func (s *BaseLuaParserListener) ExitVar(ctx *VarContext) {}

// EnterPrefixexp is called when production prefixexp is entered.
func (s *BaseLuaParserListener) EnterPrefixexp(ctx *PrefixexpContext) {}

// ExitPrefixexp is called when production prefixexp is exited.
func (s *BaseLuaParserListener) ExitPrefixexp(ctx *PrefixexpContext) {}

// EnterFunctioncall is called when production functioncall is entered.
func (s *BaseLuaParserListener) EnterFunctioncall(ctx *FunctioncallContext) {}

// ExitFunctioncall is called when production functioncall is exited.
func (s *BaseLuaParserListener) ExitFunctioncall(ctx *FunctioncallContext) {}

// EnterArgs is called when production args is entered.
func (s *BaseLuaParserListener) EnterArgs(ctx *ArgsContext) {}

// ExitArgs is called when production args is exited.
func (s *BaseLuaParserListener) ExitArgs(ctx *ArgsContext) {}

// EnterFunctiondef is called when production functiondef is entered.
func (s *BaseLuaParserListener) EnterFunctiondef(ctx *FunctiondefContext) {}

// ExitFunctiondef is called when production functiondef is exited.
func (s *BaseLuaParserListener) ExitFunctiondef(ctx *FunctiondefContext) {}

// EnterFuncbody is called when production funcbody is entered.
func (s *BaseLuaParserListener) EnterFuncbody(ctx *FuncbodyContext) {}

// ExitFuncbody is called when production funcbody is exited.
func (s *BaseLuaParserListener) ExitFuncbody(ctx *FuncbodyContext) {}

// EnterParlist is called when production parlist is entered.
func (s *BaseLuaParserListener) EnterParlist(ctx *ParlistContext) {}

// ExitParlist is called when production parlist is exited.
func (s *BaseLuaParserListener) ExitParlist(ctx *ParlistContext) {}

// EnterTableconstructor is called when production tableconstructor is entered.
func (s *BaseLuaParserListener) EnterTableconstructor(ctx *TableconstructorContext) {}

// ExitTableconstructor is called when production tableconstructor is exited.
func (s *BaseLuaParserListener) ExitTableconstructor(ctx *TableconstructorContext) {}

// EnterFieldlist is called when production fieldlist is entered.
func (s *BaseLuaParserListener) EnterFieldlist(ctx *FieldlistContext) {}

// ExitFieldlist is called when production fieldlist is exited.
func (s *BaseLuaParserListener) ExitFieldlist(ctx *FieldlistContext) {}

// EnterField is called when production field is entered.
func (s *BaseLuaParserListener) EnterField(ctx *FieldContext) {}

// ExitField is called when production field is exited.
func (s *BaseLuaParserListener) ExitField(ctx *FieldContext) {}

// EnterFieldsep is called when production fieldsep is entered.
func (s *BaseLuaParserListener) EnterFieldsep(ctx *FieldsepContext) {}

// ExitFieldsep is called when production fieldsep is exited.
func (s *BaseLuaParserListener) ExitFieldsep(ctx *FieldsepContext) {}

// EnterNumber is called when production number is entered.
func (s *BaseLuaParserListener) EnterNumber(ctx *NumberContext) {}

// ExitNumber is called when production number is exited.
func (s *BaseLuaParserListener) ExitNumber(ctx *NumberContext) {}

// EnterString is called when production string is entered.
func (s *BaseLuaParserListener) EnterString(ctx *StringContext) {}

// ExitString is called when production string is exited.
func (s *BaseLuaParserListener) ExitString(ctx *StringContext) {}
