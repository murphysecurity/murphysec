// Code generated from C:/Users/iseki/working/client/module/luarocks/parser/LuaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // LuaParser

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type LuaParser struct {
	*antlr.BaseParser
}

var LuaParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func luaparserParserInit() {
	staticData := &LuaParserParserStaticData
	staticData.LiteralNames = []string{
		"", "';'", "'='", "'break'", "'goto'", "'do'", "'end'", "'while'", "'repeat'",
		"'until'", "'if'", "'then'", "'elseif'", "'else'", "'for'", "','", "'in'",
		"'function'", "'local'", "'<'", "'>'", "'return'", "'continue'", "'::'",
		"'nil'", "'false'", "'true'", "'.'", "'~'", "'-'", "'#'", "'('", "')'",
		"'not'", "'<<'", "'>>'", "'&'", "'//'", "'%'", "':'", "'<='", "'>='",
		"'and'", "'or'", "'+'", "'*'", "'{'", "'}'", "'['", "']'", "'=='", "'..'",
		"'|'", "'^'", "'/'", "'...'", "'~='",
	}
	staticData.SymbolicNames = []string{
		"", "SEMI", "EQ", "BREAK", "GOTO", "DO", "END", "WHILE", "REPEAT", "UNTIL",
		"IF", "THEN", "ELSEIF", "ELSE", "FOR", "COMMA", "IN", "FUNCTION", "LOCAL",
		"LT", "GT", "RETURN", "CONTINUE", "CC", "NIL", "FALSE", "TRUE", "DOT",
		"SQUIG", "MINUS", "POUND", "OP", "CP", "NOT", "LL", "GG", "AMP", "SS",
		"PER", "COL", "LE", "GE", "AND", "OR", "PLUS", "STAR", "OCU", "CCU",
		"OB", "CB", "EE", "DD", "PIPE", "CARET", "SLASH", "DDD", "SQEQ", "NAME",
		"NORMALSTRING", "CHARSTRING", "LONGSTRING", "INT", "HEX", "FLOAT", "HEX_FLOAT",
		"COMMENT", "WS", "NL", "SHEBANG",
	}
	staticData.RuleNames = []string{
		"start_", "chunk", "block", "stat", "attnamelist", "attrib", "retstat",
		"label", "funcname", "varlist", "namelist", "explist", "exp", "var",
		"prefixexp", "functioncall", "args", "functiondef", "funcbody", "parlist",
		"tableconstructor", "fieldlist", "field", "fieldsep", "number", "string",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 68, 472, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 1, 0,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 5, 2, 59, 8, 2, 10, 2, 12, 2, 62, 9, 2, 1,
		2, 3, 2, 65, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 5, 3, 101, 8, 3, 10, 3, 12, 3, 104, 9, 3, 1, 3, 1, 3, 3, 3, 108,
		8, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3,
		120, 8, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 3, 3, 146, 8, 3, 3, 3, 148, 8, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1,
		4, 5, 4, 155, 8, 4, 10, 4, 12, 4, 158, 9, 4, 1, 5, 1, 5, 1, 5, 3, 5, 163,
		8, 5, 1, 6, 1, 6, 3, 6, 167, 8, 6, 1, 6, 1, 6, 3, 6, 171, 8, 6, 1, 6, 3,
		6, 174, 8, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 5, 8, 183, 8, 8,
		10, 8, 12, 8, 186, 9, 8, 1, 8, 1, 8, 3, 8, 190, 8, 8, 1, 9, 1, 9, 1, 9,
		5, 9, 195, 8, 9, 10, 9, 12, 9, 198, 9, 9, 1, 10, 1, 10, 1, 10, 5, 10, 203,
		8, 10, 10, 10, 12, 10, 206, 9, 10, 1, 11, 1, 11, 1, 11, 5, 11, 211, 8,
		11, 10, 11, 12, 11, 214, 9, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12,
		1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 3, 12, 228, 8, 12, 1, 12, 1,
		12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12,
		1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1,
		12, 1, 12, 5, 12, 254, 8, 12, 10, 12, 12, 12, 257, 9, 12, 1, 13, 1, 13,
		1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 267, 8, 13, 3, 13, 269,
		8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 5, 14, 278, 8,
		14, 10, 14, 12, 14, 281, 9, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14,
		1, 14, 5, 14, 290, 8, 14, 10, 14, 12, 14, 293, 9, 14, 1, 14, 1, 14, 1,
		14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 5, 14, 304, 8, 14, 10, 14,
		12, 14, 307, 9, 14, 3, 14, 309, 8, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15,
		1, 15, 1, 15, 1, 15, 5, 15, 319, 8, 15, 10, 15, 12, 15, 322, 9, 15, 1,
		15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 5, 15,
		334, 8, 15, 10, 15, 12, 15, 337, 9, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1,
		15, 1, 15, 1, 15, 1, 15, 1, 15, 5, 15, 348, 8, 15, 10, 15, 12, 15, 351,
		9, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1,
		15, 1, 15, 1, 15, 5, 15, 365, 8, 15, 10, 15, 12, 15, 368, 9, 15, 1, 15,
		1, 15, 1, 15, 1, 15, 3, 15, 374, 8, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1,
		15, 1, 15, 1, 15, 5, 15, 383, 8, 15, 10, 15, 12, 15, 386, 9, 15, 1, 15,
		1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 5, 15, 396, 8, 15, 10,
		15, 12, 15, 399, 9, 15, 1, 15, 1, 15, 1, 15, 5, 15, 404, 8, 15, 10, 15,
		12, 15, 407, 9, 15, 1, 16, 1, 16, 3, 16, 411, 8, 16, 1, 16, 1, 16, 1, 16,
		3, 16, 416, 8, 16, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1,
		18, 1, 18, 1, 19, 1, 19, 1, 19, 3, 19, 430, 8, 19, 1, 19, 1, 19, 3, 19,
		434, 8, 19, 1, 20, 1, 20, 3, 20, 438, 8, 20, 1, 20, 1, 20, 1, 21, 1, 21,
		1, 21, 1, 21, 5, 21, 446, 8, 21, 10, 21, 12, 21, 449, 9, 21, 1, 21, 3,
		21, 452, 8, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22,
		1, 22, 1, 22, 3, 22, 464, 8, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1,
		25, 1, 25, 0, 2, 24, 30, 26, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22,
		24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 0, 8, 2, 0, 28,
		30, 33, 33, 3, 0, 37, 38, 45, 45, 54, 54, 2, 0, 29, 29, 44, 44, 4, 0, 19,
		20, 40, 41, 50, 50, 56, 56, 3, 0, 28, 28, 34, 36, 52, 52, 2, 0, 1, 1, 15,
		15, 1, 0, 61, 64, 1, 0, 58, 60, 531, 0, 52, 1, 0, 0, 0, 2, 55, 1, 0, 0,
		0, 4, 60, 1, 0, 0, 0, 6, 147, 1, 0, 0, 0, 8, 149, 1, 0, 0, 0, 10, 162,
		1, 0, 0, 0, 12, 170, 1, 0, 0, 0, 14, 175, 1, 0, 0, 0, 16, 179, 1, 0, 0,
		0, 18, 191, 1, 0, 0, 0, 20, 199, 1, 0, 0, 0, 22, 207, 1, 0, 0, 0, 24, 227,
		1, 0, 0, 0, 26, 268, 1, 0, 0, 0, 28, 308, 1, 0, 0, 0, 30, 373, 1, 0, 0,
		0, 32, 415, 1, 0, 0, 0, 34, 417, 1, 0, 0, 0, 36, 420, 1, 0, 0, 0, 38, 433,
		1, 0, 0, 0, 40, 435, 1, 0, 0, 0, 42, 441, 1, 0, 0, 0, 44, 463, 1, 0, 0,
		0, 46, 465, 1, 0, 0, 0, 48, 467, 1, 0, 0, 0, 50, 469, 1, 0, 0, 0, 52, 53,
		3, 2, 1, 0, 53, 54, 5, 0, 0, 1, 54, 1, 1, 0, 0, 0, 55, 56, 3, 4, 2, 0,
		56, 3, 1, 0, 0, 0, 57, 59, 3, 6, 3, 0, 58, 57, 1, 0, 0, 0, 59, 62, 1, 0,
		0, 0, 60, 58, 1, 0, 0, 0, 60, 61, 1, 0, 0, 0, 61, 64, 1, 0, 0, 0, 62, 60,
		1, 0, 0, 0, 63, 65, 3, 12, 6, 0, 64, 63, 1, 0, 0, 0, 64, 65, 1, 0, 0, 0,
		65, 5, 1, 0, 0, 0, 66, 148, 5, 1, 0, 0, 67, 68, 3, 18, 9, 0, 68, 69, 5,
		2, 0, 0, 69, 70, 3, 22, 11, 0, 70, 148, 1, 0, 0, 0, 71, 148, 3, 30, 15,
		0, 72, 148, 3, 14, 7, 0, 73, 148, 5, 3, 0, 0, 74, 75, 5, 4, 0, 0, 75, 148,
		5, 57, 0, 0, 76, 77, 5, 5, 0, 0, 77, 78, 3, 4, 2, 0, 78, 79, 5, 6, 0, 0,
		79, 148, 1, 0, 0, 0, 80, 81, 5, 7, 0, 0, 81, 82, 3, 24, 12, 0, 82, 83,
		5, 5, 0, 0, 83, 84, 3, 4, 2, 0, 84, 85, 5, 6, 0, 0, 85, 148, 1, 0, 0, 0,
		86, 87, 5, 8, 0, 0, 87, 88, 3, 4, 2, 0, 88, 89, 5, 9, 0, 0, 89, 90, 3,
		24, 12, 0, 90, 148, 1, 0, 0, 0, 91, 92, 5, 10, 0, 0, 92, 93, 3, 24, 12,
		0, 93, 94, 5, 11, 0, 0, 94, 102, 3, 4, 2, 0, 95, 96, 5, 12, 0, 0, 96, 97,
		3, 24, 12, 0, 97, 98, 5, 11, 0, 0, 98, 99, 3, 4, 2, 0, 99, 101, 1, 0, 0,
		0, 100, 95, 1, 0, 0, 0, 101, 104, 1, 0, 0, 0, 102, 100, 1, 0, 0, 0, 102,
		103, 1, 0, 0, 0, 103, 107, 1, 0, 0, 0, 104, 102, 1, 0, 0, 0, 105, 106,
		5, 13, 0, 0, 106, 108, 3, 4, 2, 0, 107, 105, 1, 0, 0, 0, 107, 108, 1, 0,
		0, 0, 108, 109, 1, 0, 0, 0, 109, 110, 5, 6, 0, 0, 110, 148, 1, 0, 0, 0,
		111, 112, 5, 14, 0, 0, 112, 113, 5, 57, 0, 0, 113, 114, 5, 2, 0, 0, 114,
		115, 3, 24, 12, 0, 115, 116, 5, 15, 0, 0, 116, 119, 3, 24, 12, 0, 117,
		118, 5, 15, 0, 0, 118, 120, 3, 24, 12, 0, 119, 117, 1, 0, 0, 0, 119, 120,
		1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 122, 5, 5, 0, 0, 122, 123, 3, 4,
		2, 0, 123, 124, 5, 6, 0, 0, 124, 148, 1, 0, 0, 0, 125, 126, 5, 14, 0, 0,
		126, 127, 3, 20, 10, 0, 127, 128, 5, 16, 0, 0, 128, 129, 3, 22, 11, 0,
		129, 130, 5, 5, 0, 0, 130, 131, 3, 4, 2, 0, 131, 132, 5, 6, 0, 0, 132,
		148, 1, 0, 0, 0, 133, 134, 5, 17, 0, 0, 134, 135, 3, 16, 8, 0, 135, 136,
		3, 36, 18, 0, 136, 148, 1, 0, 0, 0, 137, 138, 5, 18, 0, 0, 138, 139, 5,
		17, 0, 0, 139, 140, 5, 57, 0, 0, 140, 148, 3, 36, 18, 0, 141, 142, 5, 18,
		0, 0, 142, 145, 3, 8, 4, 0, 143, 144, 5, 2, 0, 0, 144, 146, 3, 22, 11,
		0, 145, 143, 1, 0, 0, 0, 145, 146, 1, 0, 0, 0, 146, 148, 1, 0, 0, 0, 147,
		66, 1, 0, 0, 0, 147, 67, 1, 0, 0, 0, 147, 71, 1, 0, 0, 0, 147, 72, 1, 0,
		0, 0, 147, 73, 1, 0, 0, 0, 147, 74, 1, 0, 0, 0, 147, 76, 1, 0, 0, 0, 147,
		80, 1, 0, 0, 0, 147, 86, 1, 0, 0, 0, 147, 91, 1, 0, 0, 0, 147, 111, 1,
		0, 0, 0, 147, 125, 1, 0, 0, 0, 147, 133, 1, 0, 0, 0, 147, 137, 1, 0, 0,
		0, 147, 141, 1, 0, 0, 0, 148, 7, 1, 0, 0, 0, 149, 150, 5, 57, 0, 0, 150,
		156, 3, 10, 5, 0, 151, 152, 5, 15, 0, 0, 152, 153, 5, 57, 0, 0, 153, 155,
		3, 10, 5, 0, 154, 151, 1, 0, 0, 0, 155, 158, 1, 0, 0, 0, 156, 154, 1, 0,
		0, 0, 156, 157, 1, 0, 0, 0, 157, 9, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 159,
		160, 5, 19, 0, 0, 160, 161, 5, 57, 0, 0, 161, 163, 5, 20, 0, 0, 162, 159,
		1, 0, 0, 0, 162, 163, 1, 0, 0, 0, 163, 11, 1, 0, 0, 0, 164, 166, 5, 21,
		0, 0, 165, 167, 3, 22, 11, 0, 166, 165, 1, 0, 0, 0, 166, 167, 1, 0, 0,
		0, 167, 171, 1, 0, 0, 0, 168, 171, 5, 3, 0, 0, 169, 171, 5, 22, 0, 0, 170,
		164, 1, 0, 0, 0, 170, 168, 1, 0, 0, 0, 170, 169, 1, 0, 0, 0, 171, 173,
		1, 0, 0, 0, 172, 174, 5, 1, 0, 0, 173, 172, 1, 0, 0, 0, 173, 174, 1, 0,
		0, 0, 174, 13, 1, 0, 0, 0, 175, 176, 5, 23, 0, 0, 176, 177, 5, 57, 0, 0,
		177, 178, 5, 23, 0, 0, 178, 15, 1, 0, 0, 0, 179, 184, 5, 57, 0, 0, 180,
		181, 5, 27, 0, 0, 181, 183, 5, 57, 0, 0, 182, 180, 1, 0, 0, 0, 183, 186,
		1, 0, 0, 0, 184, 182, 1, 0, 0, 0, 184, 185, 1, 0, 0, 0, 185, 189, 1, 0,
		0, 0, 186, 184, 1, 0, 0, 0, 187, 188, 5, 39, 0, 0, 188, 190, 5, 57, 0,
		0, 189, 187, 1, 0, 0, 0, 189, 190, 1, 0, 0, 0, 190, 17, 1, 0, 0, 0, 191,
		196, 3, 26, 13, 0, 192, 193, 5, 15, 0, 0, 193, 195, 3, 26, 13, 0, 194,
		192, 1, 0, 0, 0, 195, 198, 1, 0, 0, 0, 196, 194, 1, 0, 0, 0, 196, 197,
		1, 0, 0, 0, 197, 19, 1, 0, 0, 0, 198, 196, 1, 0, 0, 0, 199, 204, 5, 57,
		0, 0, 200, 201, 5, 15, 0, 0, 201, 203, 5, 57, 0, 0, 202, 200, 1, 0, 0,
		0, 203, 206, 1, 0, 0, 0, 204, 202, 1, 0, 0, 0, 204, 205, 1, 0, 0, 0, 205,
		21, 1, 0, 0, 0, 206, 204, 1, 0, 0, 0, 207, 212, 3, 24, 12, 0, 208, 209,
		5, 15, 0, 0, 209, 211, 3, 24, 12, 0, 210, 208, 1, 0, 0, 0, 211, 214, 1,
		0, 0, 0, 212, 210, 1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 213, 23, 1, 0, 0,
		0, 214, 212, 1, 0, 0, 0, 215, 216, 6, 12, -1, 0, 216, 228, 5, 24, 0, 0,
		217, 228, 5, 25, 0, 0, 218, 228, 5, 26, 0, 0, 219, 228, 3, 48, 24, 0, 220,
		228, 3, 50, 25, 0, 221, 228, 5, 55, 0, 0, 222, 228, 3, 34, 17, 0, 223,
		228, 3, 28, 14, 0, 224, 228, 3, 40, 20, 0, 225, 226, 7, 0, 0, 0, 226, 228,
		3, 24, 12, 8, 227, 215, 1, 0, 0, 0, 227, 217, 1, 0, 0, 0, 227, 218, 1,
		0, 0, 0, 227, 219, 1, 0, 0, 0, 227, 220, 1, 0, 0, 0, 227, 221, 1, 0, 0,
		0, 227, 222, 1, 0, 0, 0, 227, 223, 1, 0, 0, 0, 227, 224, 1, 0, 0, 0, 227,
		225, 1, 0, 0, 0, 228, 255, 1, 0, 0, 0, 229, 230, 10, 9, 0, 0, 230, 231,
		5, 53, 0, 0, 231, 254, 3, 24, 12, 9, 232, 233, 10, 7, 0, 0, 233, 234, 7,
		1, 0, 0, 234, 254, 3, 24, 12, 8, 235, 236, 10, 6, 0, 0, 236, 237, 7, 2,
		0, 0, 237, 254, 3, 24, 12, 7, 238, 239, 10, 5, 0, 0, 239, 240, 5, 51, 0,
		0, 240, 254, 3, 24, 12, 5, 241, 242, 10, 4, 0, 0, 242, 243, 7, 3, 0, 0,
		243, 254, 3, 24, 12, 5, 244, 245, 10, 3, 0, 0, 245, 246, 5, 42, 0, 0, 246,
		254, 3, 24, 12, 4, 247, 248, 10, 2, 0, 0, 248, 249, 5, 43, 0, 0, 249, 254,
		3, 24, 12, 3, 250, 251, 10, 1, 0, 0, 251, 252, 7, 4, 0, 0, 252, 254, 3,
		24, 12, 2, 253, 229, 1, 0, 0, 0, 253, 232, 1, 0, 0, 0, 253, 235, 1, 0,
		0, 0, 253, 238, 1, 0, 0, 0, 253, 241, 1, 0, 0, 0, 253, 244, 1, 0, 0, 0,
		253, 247, 1, 0, 0, 0, 253, 250, 1, 0, 0, 0, 254, 257, 1, 0, 0, 0, 255,
		253, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 25, 1, 0, 0, 0, 257, 255, 1,
		0, 0, 0, 258, 269, 5, 57, 0, 0, 259, 266, 3, 28, 14, 0, 260, 261, 5, 48,
		0, 0, 261, 262, 3, 24, 12, 0, 262, 263, 5, 49, 0, 0, 263, 267, 1, 0, 0,
		0, 264, 265, 5, 27, 0, 0, 265, 267, 5, 57, 0, 0, 266, 260, 1, 0, 0, 0,
		266, 264, 1, 0, 0, 0, 267, 269, 1, 0, 0, 0, 268, 258, 1, 0, 0, 0, 268,
		259, 1, 0, 0, 0, 269, 27, 1, 0, 0, 0, 270, 279, 5, 57, 0, 0, 271, 272,
		5, 48, 0, 0, 272, 273, 3, 24, 12, 0, 273, 274, 5, 49, 0, 0, 274, 278, 1,
		0, 0, 0, 275, 276, 5, 27, 0, 0, 276, 278, 5, 57, 0, 0, 277, 271, 1, 0,
		0, 0, 277, 275, 1, 0, 0, 0, 278, 281, 1, 0, 0, 0, 279, 277, 1, 0, 0, 0,
		279, 280, 1, 0, 0, 0, 280, 309, 1, 0, 0, 0, 281, 279, 1, 0, 0, 0, 282,
		291, 3, 30, 15, 0, 283, 284, 5, 48, 0, 0, 284, 285, 3, 24, 12, 0, 285,
		286, 5, 49, 0, 0, 286, 290, 1, 0, 0, 0, 287, 288, 5, 27, 0, 0, 288, 290,
		5, 57, 0, 0, 289, 283, 1, 0, 0, 0, 289, 287, 1, 0, 0, 0, 290, 293, 1, 0,
		0, 0, 291, 289, 1, 0, 0, 0, 291, 292, 1, 0, 0, 0, 292, 309, 1, 0, 0, 0,
		293, 291, 1, 0, 0, 0, 294, 295, 5, 31, 0, 0, 295, 296, 3, 24, 12, 0, 296,
		305, 5, 32, 0, 0, 297, 298, 5, 48, 0, 0, 298, 299, 3, 24, 12, 0, 299, 300,
		5, 49, 0, 0, 300, 304, 1, 0, 0, 0, 301, 302, 5, 27, 0, 0, 302, 304, 5,
		57, 0, 0, 303, 297, 1, 0, 0, 0, 303, 301, 1, 0, 0, 0, 304, 307, 1, 0, 0,
		0, 305, 303, 1, 0, 0, 0, 305, 306, 1, 0, 0, 0, 306, 309, 1, 0, 0, 0, 307,
		305, 1, 0, 0, 0, 308, 270, 1, 0, 0, 0, 308, 282, 1, 0, 0, 0, 308, 294,
		1, 0, 0, 0, 309, 29, 1, 0, 0, 0, 310, 311, 6, 15, -1, 0, 311, 320, 5, 57,
		0, 0, 312, 313, 5, 48, 0, 0, 313, 314, 3, 24, 12, 0, 314, 315, 5, 49, 0,
		0, 315, 319, 1, 0, 0, 0, 316, 317, 5, 27, 0, 0, 317, 319, 5, 57, 0, 0,
		318, 312, 1, 0, 0, 0, 318, 316, 1, 0, 0, 0, 319, 322, 1, 0, 0, 0, 320,
		318, 1, 0, 0, 0, 320, 321, 1, 0, 0, 0, 321, 323, 1, 0, 0, 0, 322, 320,
		1, 0, 0, 0, 323, 374, 3, 32, 16, 0, 324, 325, 5, 31, 0, 0, 325, 326, 3,
		24, 12, 0, 326, 335, 5, 32, 0, 0, 327, 328, 5, 48, 0, 0, 328, 329, 3, 24,
		12, 0, 329, 330, 5, 49, 0, 0, 330, 334, 1, 0, 0, 0, 331, 332, 5, 27, 0,
		0, 332, 334, 5, 57, 0, 0, 333, 327, 1, 0, 0, 0, 333, 331, 1, 0, 0, 0, 334,
		337, 1, 0, 0, 0, 335, 333, 1, 0, 0, 0, 335, 336, 1, 0, 0, 0, 336, 338,
		1, 0, 0, 0, 337, 335, 1, 0, 0, 0, 338, 339, 3, 32, 16, 0, 339, 374, 1,
		0, 0, 0, 340, 349, 5, 57, 0, 0, 341, 342, 5, 48, 0, 0, 342, 343, 3, 24,
		12, 0, 343, 344, 5, 49, 0, 0, 344, 348, 1, 0, 0, 0, 345, 346, 5, 27, 0,
		0, 346, 348, 5, 57, 0, 0, 347, 341, 1, 0, 0, 0, 347, 345, 1, 0, 0, 0, 348,
		351, 1, 0, 0, 0, 349, 347, 1, 0, 0, 0, 349, 350, 1, 0, 0, 0, 350, 352,
		1, 0, 0, 0, 351, 349, 1, 0, 0, 0, 352, 353, 5, 39, 0, 0, 353, 354, 5, 57,
		0, 0, 354, 374, 3, 32, 16, 0, 355, 356, 5, 31, 0, 0, 356, 357, 3, 24, 12,
		0, 357, 366, 5, 32, 0, 0, 358, 359, 5, 48, 0, 0, 359, 360, 3, 24, 12, 0,
		360, 361, 5, 49, 0, 0, 361, 365, 1, 0, 0, 0, 362, 363, 5, 27, 0, 0, 363,
		365, 5, 57, 0, 0, 364, 358, 1, 0, 0, 0, 364, 362, 1, 0, 0, 0, 365, 368,
		1, 0, 0, 0, 366, 364, 1, 0, 0, 0, 366, 367, 1, 0, 0, 0, 367, 369, 1, 0,
		0, 0, 368, 366, 1, 0, 0, 0, 369, 370, 5, 39, 0, 0, 370, 371, 5, 57, 0,
		0, 371, 372, 3, 32, 16, 0, 372, 374, 1, 0, 0, 0, 373, 310, 1, 0, 0, 0,
		373, 324, 1, 0, 0, 0, 373, 340, 1, 0, 0, 0, 373, 355, 1, 0, 0, 0, 374,
		405, 1, 0, 0, 0, 375, 384, 10, 5, 0, 0, 376, 377, 5, 48, 0, 0, 377, 378,
		3, 24, 12, 0, 378, 379, 5, 49, 0, 0, 379, 383, 1, 0, 0, 0, 380, 381, 5,
		27, 0, 0, 381, 383, 5, 57, 0, 0, 382, 376, 1, 0, 0, 0, 382, 380, 1, 0,
		0, 0, 383, 386, 1, 0, 0, 0, 384, 382, 1, 0, 0, 0, 384, 385, 1, 0, 0, 0,
		385, 387, 1, 0, 0, 0, 386, 384, 1, 0, 0, 0, 387, 404, 3, 32, 16, 0, 388,
		397, 10, 2, 0, 0, 389, 390, 5, 48, 0, 0, 390, 391, 3, 24, 12, 0, 391, 392,
		5, 49, 0, 0, 392, 396, 1, 0, 0, 0, 393, 394, 5, 27, 0, 0, 394, 396, 5,
		57, 0, 0, 395, 389, 1, 0, 0, 0, 395, 393, 1, 0, 0, 0, 396, 399, 1, 0, 0,
		0, 397, 395, 1, 0, 0, 0, 397, 398, 1, 0, 0, 0, 398, 400, 1, 0, 0, 0, 399,
		397, 1, 0, 0, 0, 400, 401, 5, 39, 0, 0, 401, 402, 5, 57, 0, 0, 402, 404,
		3, 32, 16, 0, 403, 375, 1, 0, 0, 0, 403, 388, 1, 0, 0, 0, 404, 407, 1,
		0, 0, 0, 405, 403, 1, 0, 0, 0, 405, 406, 1, 0, 0, 0, 406, 31, 1, 0, 0,
		0, 407, 405, 1, 0, 0, 0, 408, 410, 5, 31, 0, 0, 409, 411, 3, 22, 11, 0,
		410, 409, 1, 0, 0, 0, 410, 411, 1, 0, 0, 0, 411, 412, 1, 0, 0, 0, 412,
		416, 5, 32, 0, 0, 413, 416, 3, 40, 20, 0, 414, 416, 3, 50, 25, 0, 415,
		408, 1, 0, 0, 0, 415, 413, 1, 0, 0, 0, 415, 414, 1, 0, 0, 0, 416, 33, 1,
		0, 0, 0, 417, 418, 5, 17, 0, 0, 418, 419, 3, 36, 18, 0, 419, 35, 1, 0,
		0, 0, 420, 421, 5, 31, 0, 0, 421, 422, 3, 38, 19, 0, 422, 423, 5, 32, 0,
		0, 423, 424, 3, 4, 2, 0, 424, 425, 5, 6, 0, 0, 425, 37, 1, 0, 0, 0, 426,
		429, 3, 20, 10, 0, 427, 428, 5, 15, 0, 0, 428, 430, 5, 55, 0, 0, 429, 427,
		1, 0, 0, 0, 429, 430, 1, 0, 0, 0, 430, 434, 1, 0, 0, 0, 431, 434, 5, 55,
		0, 0, 432, 434, 1, 0, 0, 0, 433, 426, 1, 0, 0, 0, 433, 431, 1, 0, 0, 0,
		433, 432, 1, 0, 0, 0, 434, 39, 1, 0, 0, 0, 435, 437, 5, 46, 0, 0, 436,
		438, 3, 42, 21, 0, 437, 436, 1, 0, 0, 0, 437, 438, 1, 0, 0, 0, 438, 439,
		1, 0, 0, 0, 439, 440, 5, 47, 0, 0, 440, 41, 1, 0, 0, 0, 441, 447, 3, 44,
		22, 0, 442, 443, 3, 46, 23, 0, 443, 444, 3, 44, 22, 0, 444, 446, 1, 0,
		0, 0, 445, 442, 1, 0, 0, 0, 446, 449, 1, 0, 0, 0, 447, 445, 1, 0, 0, 0,
		447, 448, 1, 0, 0, 0, 448, 451, 1, 0, 0, 0, 449, 447, 1, 0, 0, 0, 450,
		452, 3, 46, 23, 0, 451, 450, 1, 0, 0, 0, 451, 452, 1, 0, 0, 0, 452, 43,
		1, 0, 0, 0, 453, 454, 5, 48, 0, 0, 454, 455, 3, 24, 12, 0, 455, 456, 5,
		49, 0, 0, 456, 457, 5, 2, 0, 0, 457, 458, 3, 24, 12, 0, 458, 464, 1, 0,
		0, 0, 459, 460, 5, 57, 0, 0, 460, 461, 5, 2, 0, 0, 461, 464, 3, 24, 12,
		0, 462, 464, 3, 24, 12, 0, 463, 453, 1, 0, 0, 0, 463, 459, 1, 0, 0, 0,
		463, 462, 1, 0, 0, 0, 464, 45, 1, 0, 0, 0, 465, 466, 7, 5, 0, 0, 466, 47,
		1, 0, 0, 0, 467, 468, 7, 6, 0, 0, 468, 49, 1, 0, 0, 0, 469, 470, 7, 7,
		0, 0, 470, 51, 1, 0, 0, 0, 52, 60, 64, 102, 107, 119, 145, 147, 156, 162,
		166, 170, 173, 184, 189, 196, 204, 212, 227, 253, 255, 266, 268, 277, 279,
		289, 291, 303, 305, 308, 318, 320, 333, 335, 347, 349, 364, 366, 373, 382,
		384, 395, 397, 403, 405, 410, 415, 429, 433, 437, 447, 451, 463,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// LuaParserInit initializes any static state used to implement LuaParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewLuaParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func LuaParserInit() {
	staticData := &LuaParserParserStaticData
	staticData.once.Do(luaparserParserInit)
}

// NewLuaParser produces a new parser instance for the optional input antlr.TokenStream.
func NewLuaParser(input antlr.TokenStream) *LuaParser {
	LuaParserInit()
	this := new(LuaParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &LuaParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "LuaParser.g4"

	return this
}

// LuaParser tokens.
const (
	LuaParserEOF          = antlr.TokenEOF
	LuaParserSEMI         = 1
	LuaParserEQ           = 2
	LuaParserBREAK        = 3
	LuaParserGOTO         = 4
	LuaParserDO           = 5
	LuaParserEND          = 6
	LuaParserWHILE        = 7
	LuaParserREPEAT       = 8
	LuaParserUNTIL        = 9
	LuaParserIF           = 10
	LuaParserTHEN         = 11
	LuaParserELSEIF       = 12
	LuaParserELSE         = 13
	LuaParserFOR          = 14
	LuaParserCOMMA        = 15
	LuaParserIN           = 16
	LuaParserFUNCTION     = 17
	LuaParserLOCAL        = 18
	LuaParserLT           = 19
	LuaParserGT           = 20
	LuaParserRETURN       = 21
	LuaParserCONTINUE     = 22
	LuaParserCC           = 23
	LuaParserNIL          = 24
	LuaParserFALSE        = 25
	LuaParserTRUE         = 26
	LuaParserDOT          = 27
	LuaParserSQUIG        = 28
	LuaParserMINUS        = 29
	LuaParserPOUND        = 30
	LuaParserOP           = 31
	LuaParserCP           = 32
	LuaParserNOT          = 33
	LuaParserLL           = 34
	LuaParserGG           = 35
	LuaParserAMP          = 36
	LuaParserSS           = 37
	LuaParserPER          = 38
	LuaParserCOL          = 39
	LuaParserLE           = 40
	LuaParserGE           = 41
	LuaParserAND          = 42
	LuaParserOR           = 43
	LuaParserPLUS         = 44
	LuaParserSTAR         = 45
	LuaParserOCU          = 46
	LuaParserCCU          = 47
	LuaParserOB           = 48
	LuaParserCB           = 49
	LuaParserEE           = 50
	LuaParserDD           = 51
	LuaParserPIPE         = 52
	LuaParserCARET        = 53
	LuaParserSLASH        = 54
	LuaParserDDD          = 55
	LuaParserSQEQ         = 56
	LuaParserNAME         = 57
	LuaParserNORMALSTRING = 58
	LuaParserCHARSTRING   = 59
	LuaParserLONGSTRING   = 60
	LuaParserINT          = 61
	LuaParserHEX          = 62
	LuaParserFLOAT        = 63
	LuaParserHEX_FLOAT    = 64
	LuaParserCOMMENT      = 65
	LuaParserWS           = 66
	LuaParserNL           = 67
	LuaParserSHEBANG      = 68
)

// LuaParser rules.
const (
	LuaParserRULE_start_           = 0
	LuaParserRULE_chunk            = 1
	LuaParserRULE_block            = 2
	LuaParserRULE_stat             = 3
	LuaParserRULE_attnamelist      = 4
	LuaParserRULE_attrib           = 5
	LuaParserRULE_retstat          = 6
	LuaParserRULE_label            = 7
	LuaParserRULE_funcname         = 8
	LuaParserRULE_varlist          = 9
	LuaParserRULE_namelist         = 10
	LuaParserRULE_explist          = 11
	LuaParserRULE_exp              = 12
	LuaParserRULE_var              = 13
	LuaParserRULE_prefixexp        = 14
	LuaParserRULE_functioncall     = 15
	LuaParserRULE_args             = 16
	LuaParserRULE_functiondef      = 17
	LuaParserRULE_funcbody         = 18
	LuaParserRULE_parlist          = 19
	LuaParserRULE_tableconstructor = 20
	LuaParserRULE_fieldlist        = 21
	LuaParserRULE_field            = 22
	LuaParserRULE_fieldsep         = 23
	LuaParserRULE_number           = 24
	LuaParserRULE_string           = 25
)

// IStart_Context is an interface to support dynamic dispatch.
type IStart_Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Chunk() IChunkContext
	EOF() antlr.TerminalNode

	// IsStart_Context differentiates from other interfaces.
	IsStart_Context()
}

type Start_Context struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStart_Context() *Start_Context {
	var p = new(Start_Context)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_start_
	return p
}

func InitEmptyStart_Context(p *Start_Context) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_start_
}

func (*Start_Context) IsStart_Context() {}

func NewStart_Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Start_Context {
	var p = new(Start_Context)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_start_

	return p
}

func (s *Start_Context) GetParser() antlr.Parser { return s.parser }

func (s *Start_Context) Chunk() IChunkContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IChunkContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IChunkContext)
}

func (s *Start_Context) EOF() antlr.TerminalNode {
	return s.GetToken(LuaParserEOF, 0)
}

func (s *Start_Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Start_Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Start_Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterStart_(s)
	}
}

func (s *Start_Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitStart_(s)
	}
}

func (s *Start_Context) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitStart_(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Start_() (localctx IStart_Context) {
	localctx = NewStart_Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, LuaParserRULE_start_)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(52)
		p.Chunk()
	}
	{
		p.SetState(53)
		p.Match(LuaParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IChunkContext is an interface to support dynamic dispatch.
type IChunkContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Block() IBlockContext

	// IsChunkContext differentiates from other interfaces.
	IsChunkContext()
}

type ChunkContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyChunkContext() *ChunkContext {
	var p = new(ChunkContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_chunk
	return p
}

func InitEmptyChunkContext(p *ChunkContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_chunk
}

func (*ChunkContext) IsChunkContext() {}

func NewChunkContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ChunkContext {
	var p = new(ChunkContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_chunk

	return p
}

func (s *ChunkContext) GetParser() antlr.Parser { return s.parser }

func (s *ChunkContext) Block() IBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *ChunkContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChunkContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ChunkContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterChunk(s)
	}
}

func (s *ChunkContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitChunk(s)
	}
}

func (s *ChunkContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitChunk(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Chunk() (localctx IChunkContext) {
	localctx = NewChunkContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, LuaParserRULE_chunk)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(55)
		p.Block()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBlockContext is an interface to support dynamic dispatch.
type IBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStat() []IStatContext
	Stat(i int) IStatContext
	Retstat() IRetstatContext

	// IsBlockContext differentiates from other interfaces.
	IsBlockContext()
}

type BlockContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBlockContext() *BlockContext {
	var p = new(BlockContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_block
	return p
}

func InitEmptyBlockContext(p *BlockContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_block
}

func (*BlockContext) IsBlockContext() {}

func NewBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockContext {
	var p = new(BlockContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_block

	return p
}

func (s *BlockContext) GetParser() antlr.Parser { return s.parser }

func (s *BlockContext) AllStat() []IStatContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatContext); ok {
			len++
		}
	}

	tst := make([]IStatContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatContext); ok {
			tst[i] = t.(IStatContext)
			i++
		}
	}

	return tst
}

func (s *BlockContext) Stat(i int) IStatContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatContext)
}

func (s *BlockContext) Retstat() IRetstatContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRetstatContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRetstatContext)
}

func (s *BlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterBlock(s)
	}
}

func (s *BlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitBlock(s)
	}
}

func (s *BlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitBlock(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Block() (localctx IBlockContext) {
	localctx = NewBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, LuaParserRULE_block)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(60)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(57)
				p.Stat()
			}

		}
		p.SetState(62)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&6291464) != 0 {
		{
			p.SetState(63)
			p.Retstat()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStatContext is an interface to support dynamic dispatch.
type IStatContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SEMI() antlr.TerminalNode
	Varlist() IVarlistContext
	EQ() antlr.TerminalNode
	Explist() IExplistContext
	Functioncall() IFunctioncallContext
	Label() ILabelContext
	BREAK() antlr.TerminalNode
	GOTO() antlr.TerminalNode
	NAME() antlr.TerminalNode
	DO() antlr.TerminalNode
	AllBlock() []IBlockContext
	Block(i int) IBlockContext
	END() antlr.TerminalNode
	WHILE() antlr.TerminalNode
	AllExp() []IExpContext
	Exp(i int) IExpContext
	REPEAT() antlr.TerminalNode
	UNTIL() antlr.TerminalNode
	IF() antlr.TerminalNode
	AllTHEN() []antlr.TerminalNode
	THEN(i int) antlr.TerminalNode
	AllELSEIF() []antlr.TerminalNode
	ELSEIF(i int) antlr.TerminalNode
	ELSE() antlr.TerminalNode
	FOR() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	Namelist() INamelistContext
	IN() antlr.TerminalNode
	FUNCTION() antlr.TerminalNode
	Funcname() IFuncnameContext
	Funcbody() IFuncbodyContext
	LOCAL() antlr.TerminalNode
	Attnamelist() IAttnamelistContext

	// IsStatContext differentiates from other interfaces.
	IsStatContext()
}

type StatContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatContext() *StatContext {
	var p = new(StatContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_stat
	return p
}

func InitEmptyStatContext(p *StatContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_stat
}

func (*StatContext) IsStatContext() {}

func NewStatContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatContext {
	var p = new(StatContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_stat

	return p
}

func (s *StatContext) GetParser() antlr.Parser { return s.parser }

func (s *StatContext) SEMI() antlr.TerminalNode {
	return s.GetToken(LuaParserSEMI, 0)
}

func (s *StatContext) Varlist() IVarlistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarlistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarlistContext)
}

func (s *StatContext) EQ() antlr.TerminalNode {
	return s.GetToken(LuaParserEQ, 0)
}

func (s *StatContext) Explist() IExplistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExplistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExplistContext)
}

func (s *StatContext) Functioncall() IFunctioncallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctioncallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctioncallContext)
}

func (s *StatContext) Label() ILabelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILabelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILabelContext)
}

func (s *StatContext) BREAK() antlr.TerminalNode {
	return s.GetToken(LuaParserBREAK, 0)
}

func (s *StatContext) GOTO() antlr.TerminalNode {
	return s.GetToken(LuaParserGOTO, 0)
}

func (s *StatContext) NAME() antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, 0)
}

func (s *StatContext) DO() antlr.TerminalNode {
	return s.GetToken(LuaParserDO, 0)
}

func (s *StatContext) AllBlock() []IBlockContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBlockContext); ok {
			len++
		}
	}

	tst := make([]IBlockContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBlockContext); ok {
			tst[i] = t.(IBlockContext)
			i++
		}
	}

	return tst
}

func (s *StatContext) Block(i int) IBlockContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *StatContext) END() antlr.TerminalNode {
	return s.GetToken(LuaParserEND, 0)
}

func (s *StatContext) WHILE() antlr.TerminalNode {
	return s.GetToken(LuaParserWHILE, 0)
}

func (s *StatContext) AllExp() []IExpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpContext); ok {
			len++
		}
	}

	tst := make([]IExpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpContext); ok {
			tst[i] = t.(IExpContext)
			i++
		}
	}

	return tst
}

func (s *StatContext) Exp(i int) IExpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *StatContext) REPEAT() antlr.TerminalNode {
	return s.GetToken(LuaParserREPEAT, 0)
}

func (s *StatContext) UNTIL() antlr.TerminalNode {
	return s.GetToken(LuaParserUNTIL, 0)
}

func (s *StatContext) IF() antlr.TerminalNode {
	return s.GetToken(LuaParserIF, 0)
}

func (s *StatContext) AllTHEN() []antlr.TerminalNode {
	return s.GetTokens(LuaParserTHEN)
}

func (s *StatContext) THEN(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserTHEN, i)
}

func (s *StatContext) AllELSEIF() []antlr.TerminalNode {
	return s.GetTokens(LuaParserELSEIF)
}

func (s *StatContext) ELSEIF(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserELSEIF, i)
}

func (s *StatContext) ELSE() antlr.TerminalNode {
	return s.GetToken(LuaParserELSE, 0)
}

func (s *StatContext) FOR() antlr.TerminalNode {
	return s.GetToken(LuaParserFOR, 0)
}

func (s *StatContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCOMMA)
}

func (s *StatContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, i)
}

func (s *StatContext) Namelist() INamelistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamelistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamelistContext)
}

func (s *StatContext) IN() antlr.TerminalNode {
	return s.GetToken(LuaParserIN, 0)
}

func (s *StatContext) FUNCTION() antlr.TerminalNode {
	return s.GetToken(LuaParserFUNCTION, 0)
}

func (s *StatContext) Funcname() IFuncnameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFuncnameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFuncnameContext)
}

func (s *StatContext) Funcbody() IFuncbodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFuncbodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFuncbodyContext)
}

func (s *StatContext) LOCAL() antlr.TerminalNode {
	return s.GetToken(LuaParserLOCAL, 0)
}

func (s *StatContext) Attnamelist() IAttnamelistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAttnamelistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAttnamelistContext)
}

func (s *StatContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterStat(s)
	}
}

func (s *StatContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitStat(s)
	}
}

func (s *StatContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitStat(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Stat() (localctx IStatContext) {
	localctx = NewStatContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, LuaParserRULE_stat)
	var _la int

	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(66)
			p.Match(LuaParserSEMI)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(67)
			p.Varlist()
		}
		{
			p.SetState(68)
			p.Match(LuaParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(69)
			p.Explist()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(71)
			p.functioncall(0)
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(72)
			p.Label()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(73)
			p.Match(LuaParserBREAK)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(74)
			p.Match(LuaParserGOTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(75)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(76)
			p.Match(LuaParserDO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(77)
			p.Block()
		}
		{
			p.SetState(78)
			p.Match(LuaParserEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(80)
			p.Match(LuaParserWHILE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(81)
			p.exp(0)
		}
		{
			p.SetState(82)
			p.Match(LuaParserDO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(83)
			p.Block()
		}
		{
			p.SetState(84)
			p.Match(LuaParserEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(86)
			p.Match(LuaParserREPEAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(87)
			p.Block()
		}
		{
			p.SetState(88)
			p.Match(LuaParserUNTIL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(89)
			p.exp(0)
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(91)
			p.Match(LuaParserIF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(92)
			p.exp(0)
		}
		{
			p.SetState(93)
			p.Match(LuaParserTHEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(94)
			p.Block()
		}
		p.SetState(102)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == LuaParserELSEIF {
			{
				p.SetState(95)
				p.Match(LuaParserELSEIF)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(96)
				p.exp(0)
			}
			{
				p.SetState(97)
				p.Match(LuaParserTHEN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(98)
				p.Block()
			}

			p.SetState(104)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == LuaParserELSE {
			{
				p.SetState(105)
				p.Match(LuaParserELSE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(106)
				p.Block()
			}

		}
		{
			p.SetState(109)
			p.Match(LuaParserEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(111)
			p.Match(LuaParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(112)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(113)
			p.Match(LuaParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(114)
			p.exp(0)
		}
		{
			p.SetState(115)
			p.Match(LuaParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(116)
			p.exp(0)
		}
		p.SetState(119)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == LuaParserCOMMA {
			{
				p.SetState(117)
				p.Match(LuaParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(118)
				p.exp(0)
			}

		}
		{
			p.SetState(121)
			p.Match(LuaParserDO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(122)
			p.Block()
		}
		{
			p.SetState(123)
			p.Match(LuaParserEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(125)
			p.Match(LuaParserFOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(126)
			p.Namelist()
		}
		{
			p.SetState(127)
			p.Match(LuaParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(128)
			p.Explist()
		}
		{
			p.SetState(129)
			p.Match(LuaParserDO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(130)
			p.Block()
		}
		{
			p.SetState(131)
			p.Match(LuaParserEND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(133)
			p.Match(LuaParserFUNCTION)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(134)
			p.Funcname()
		}
		{
			p.SetState(135)
			p.Funcbody()
		}

	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(137)
			p.Match(LuaParserLOCAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(138)
			p.Match(LuaParserFUNCTION)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(139)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(140)
			p.Funcbody()
		}

	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(141)
			p.Match(LuaParserLOCAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(142)
			p.Attnamelist()
		}
		p.SetState(145)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == LuaParserEQ {
			{
				p.SetState(143)
				p.Match(LuaParserEQ)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(144)
				p.Explist()
			}

		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAttnamelistContext is an interface to support dynamic dispatch.
type IAttnamelistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNAME() []antlr.TerminalNode
	NAME(i int) antlr.TerminalNode
	AllAttrib() []IAttribContext
	Attrib(i int) IAttribContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAttnamelistContext differentiates from other interfaces.
	IsAttnamelistContext()
}

type AttnamelistContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttnamelistContext() *AttnamelistContext {
	var p = new(AttnamelistContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_attnamelist
	return p
}

func InitEmptyAttnamelistContext(p *AttnamelistContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_attnamelist
}

func (*AttnamelistContext) IsAttnamelistContext() {}

func NewAttnamelistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttnamelistContext {
	var p = new(AttnamelistContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_attnamelist

	return p
}

func (s *AttnamelistContext) GetParser() antlr.Parser { return s.parser }

func (s *AttnamelistContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(LuaParserNAME)
}

func (s *AttnamelistContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, i)
}

func (s *AttnamelistContext) AllAttrib() []IAttribContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAttribContext); ok {
			len++
		}
	}

	tst := make([]IAttribContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAttribContext); ok {
			tst[i] = t.(IAttribContext)
			i++
		}
	}

	return tst
}

func (s *AttnamelistContext) Attrib(i int) IAttribContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAttribContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAttribContext)
}

func (s *AttnamelistContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCOMMA)
}

func (s *AttnamelistContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, i)
}

func (s *AttnamelistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttnamelistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttnamelistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterAttnamelist(s)
	}
}

func (s *AttnamelistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitAttnamelist(s)
	}
}

func (s *AttnamelistContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitAttnamelist(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Attnamelist() (localctx IAttnamelistContext) {
	localctx = NewAttnamelistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, LuaParserRULE_attnamelist)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(149)
		p.Match(LuaParserNAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(150)
		p.Attrib()
	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LuaParserCOMMA {
		{
			p.SetState(151)
			p.Match(LuaParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(152)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(153)
			p.Attrib()
		}

		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAttribContext is an interface to support dynamic dispatch.
type IAttribContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LT() antlr.TerminalNode
	NAME() antlr.TerminalNode
	GT() antlr.TerminalNode

	// IsAttribContext differentiates from other interfaces.
	IsAttribContext()
}

type AttribContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttribContext() *AttribContext {
	var p = new(AttribContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_attrib
	return p
}

func InitEmptyAttribContext(p *AttribContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_attrib
}

func (*AttribContext) IsAttribContext() {}

func NewAttribContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttribContext {
	var p = new(AttribContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_attrib

	return p
}

func (s *AttribContext) GetParser() antlr.Parser { return s.parser }

func (s *AttribContext) LT() antlr.TerminalNode {
	return s.GetToken(LuaParserLT, 0)
}

func (s *AttribContext) NAME() antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, 0)
}

func (s *AttribContext) GT() antlr.TerminalNode {
	return s.GetToken(LuaParserGT, 0)
}

func (s *AttribContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttribContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttribContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterAttrib(s)
	}
}

func (s *AttribContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitAttrib(s)
	}
}

func (s *AttribContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitAttrib(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Attrib() (localctx IAttribContext) {
	localctx = NewAttribContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, LuaParserRULE_attrib)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(162)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == LuaParserLT {
		{
			p.SetState(159)
			p.Match(LuaParserLT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(160)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(161)
			p.Match(LuaParserGT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRetstatContext is an interface to support dynamic dispatch.
type IRetstatContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RETURN() antlr.TerminalNode
	BREAK() antlr.TerminalNode
	CONTINUE() antlr.TerminalNode
	SEMI() antlr.TerminalNode
	Explist() IExplistContext

	// IsRetstatContext differentiates from other interfaces.
	IsRetstatContext()
}

type RetstatContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRetstatContext() *RetstatContext {
	var p = new(RetstatContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_retstat
	return p
}

func InitEmptyRetstatContext(p *RetstatContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_retstat
}

func (*RetstatContext) IsRetstatContext() {}

func NewRetstatContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RetstatContext {
	var p = new(RetstatContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_retstat

	return p
}

func (s *RetstatContext) GetParser() antlr.Parser { return s.parser }

func (s *RetstatContext) RETURN() antlr.TerminalNode {
	return s.GetToken(LuaParserRETURN, 0)
}

func (s *RetstatContext) BREAK() antlr.TerminalNode {
	return s.GetToken(LuaParserBREAK, 0)
}

func (s *RetstatContext) CONTINUE() antlr.TerminalNode {
	return s.GetToken(LuaParserCONTINUE, 0)
}

func (s *RetstatContext) SEMI() antlr.TerminalNode {
	return s.GetToken(LuaParserSEMI, 0)
}

func (s *RetstatContext) Explist() IExplistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExplistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExplistContext)
}

func (s *RetstatContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RetstatContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RetstatContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterRetstat(s)
	}
}

func (s *RetstatContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitRetstat(s)
	}
}

func (s *RetstatContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitRetstat(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Retstat() (localctx IRetstatContext) {
	localctx = NewRetstatContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, LuaParserRULE_retstat)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case LuaParserRETURN:
		{
			p.SetState(164)
			p.Match(LuaParserRETURN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64((_la-17)) & ^0x3f) == 0 && ((int64(1)<<(_la-17))&280650879957889) != 0 {
			{
				p.SetState(165)
				p.Explist()
			}

		}

	case LuaParserBREAK:
		{
			p.SetState(168)
			p.Match(LuaParserBREAK)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserCONTINUE:
		{
			p.SetState(169)
			p.Match(LuaParserCONTINUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.SetState(173)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == LuaParserSEMI {
		{
			p.SetState(172)
			p.Match(LuaParserSEMI)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILabelContext is an interface to support dynamic dispatch.
type ILabelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCC() []antlr.TerminalNode
	CC(i int) antlr.TerminalNode
	NAME() antlr.TerminalNode

	// IsLabelContext differentiates from other interfaces.
	IsLabelContext()
}

type LabelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelContext() *LabelContext {
	var p = new(LabelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_label
	return p
}

func InitEmptyLabelContext(p *LabelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_label
}

func (*LabelContext) IsLabelContext() {}

func NewLabelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelContext {
	var p = new(LabelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_label

	return p
}

func (s *LabelContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelContext) AllCC() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCC)
}

func (s *LabelContext) CC(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCC, i)
}

func (s *LabelContext) NAME() antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, 0)
}

func (s *LabelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LabelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterLabel(s)
	}
}

func (s *LabelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitLabel(s)
	}
}

func (s *LabelContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitLabel(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Label() (localctx ILabelContext) {
	localctx = NewLabelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, LuaParserRULE_label)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(175)
		p.Match(LuaParserCC)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(176)
		p.Match(LuaParserNAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(177)
		p.Match(LuaParserCC)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFuncnameContext is an interface to support dynamic dispatch.
type IFuncnameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNAME() []antlr.TerminalNode
	NAME(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	COL() antlr.TerminalNode

	// IsFuncnameContext differentiates from other interfaces.
	IsFuncnameContext()
}

type FuncnameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncnameContext() *FuncnameContext {
	var p = new(FuncnameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_funcname
	return p
}

func InitEmptyFuncnameContext(p *FuncnameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_funcname
}

func (*FuncnameContext) IsFuncnameContext() {}

func NewFuncnameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncnameContext {
	var p = new(FuncnameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_funcname

	return p
}

func (s *FuncnameContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncnameContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(LuaParserNAME)
}

func (s *FuncnameContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, i)
}

func (s *FuncnameContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(LuaParserDOT)
}

func (s *FuncnameContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserDOT, i)
}

func (s *FuncnameContext) COL() antlr.TerminalNode {
	return s.GetToken(LuaParserCOL, 0)
}

func (s *FuncnameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncnameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncnameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterFuncname(s)
	}
}

func (s *FuncnameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitFuncname(s)
	}
}

func (s *FuncnameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitFuncname(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Funcname() (localctx IFuncnameContext) {
	localctx = NewFuncnameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, LuaParserRULE_funcname)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(179)
		p.Match(LuaParserNAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LuaParserDOT {
		{
			p.SetState(180)
			p.Match(LuaParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(181)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(186)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == LuaParserCOL {
		{
			p.SetState(187)
			p.Match(LuaParserCOL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(188)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVarlistContext is an interface to support dynamic dispatch.
type IVarlistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllVar_() []IVarContext
	Var_(i int) IVarContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsVarlistContext differentiates from other interfaces.
	IsVarlistContext()
}

type VarlistContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVarlistContext() *VarlistContext {
	var p = new(VarlistContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_varlist
	return p
}

func InitEmptyVarlistContext(p *VarlistContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_varlist
}

func (*VarlistContext) IsVarlistContext() {}

func NewVarlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarlistContext {
	var p = new(VarlistContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_varlist

	return p
}

func (s *VarlistContext) GetParser() antlr.Parser { return s.parser }

func (s *VarlistContext) AllVar_() []IVarContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IVarContext); ok {
			len++
		}
	}

	tst := make([]IVarContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IVarContext); ok {
			tst[i] = t.(IVarContext)
			i++
		}
	}

	return tst
}

func (s *VarlistContext) Var_(i int) IVarContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarContext)
}

func (s *VarlistContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCOMMA)
}

func (s *VarlistContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, i)
}

func (s *VarlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VarlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterVarlist(s)
	}
}

func (s *VarlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitVarlist(s)
	}
}

func (s *VarlistContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitVarlist(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Varlist() (localctx IVarlistContext) {
	localctx = NewVarlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, LuaParserRULE_varlist)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(191)
		p.Var_()
	}
	p.SetState(196)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LuaParserCOMMA {
		{
			p.SetState(192)
			p.Match(LuaParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(193)
			p.Var_()
		}

		p.SetState(198)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamelistContext is an interface to support dynamic dispatch.
type INamelistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNAME() []antlr.TerminalNode
	NAME(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsNamelistContext differentiates from other interfaces.
	IsNamelistContext()
}

type NamelistContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamelistContext() *NamelistContext {
	var p = new(NamelistContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_namelist
	return p
}

func InitEmptyNamelistContext(p *NamelistContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_namelist
}

func (*NamelistContext) IsNamelistContext() {}

func NewNamelistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamelistContext {
	var p = new(NamelistContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_namelist

	return p
}

func (s *NamelistContext) GetParser() antlr.Parser { return s.parser }

func (s *NamelistContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(LuaParserNAME)
}

func (s *NamelistContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, i)
}

func (s *NamelistContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCOMMA)
}

func (s *NamelistContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, i)
}

func (s *NamelistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamelistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamelistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterNamelist(s)
	}
}

func (s *NamelistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitNamelist(s)
	}
}

func (s *NamelistContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitNamelist(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Namelist() (localctx INamelistContext) {
	localctx = NewNamelistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, LuaParserRULE_namelist)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(199)
		p.Match(LuaParserNAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(204)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(200)
				p.Match(LuaParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(201)
				p.Match(LuaParserNAME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(206)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExplistContext is an interface to support dynamic dispatch.
type IExplistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExp() []IExpContext
	Exp(i int) IExpContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExplistContext differentiates from other interfaces.
	IsExplistContext()
}

type ExplistContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExplistContext() *ExplistContext {
	var p = new(ExplistContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_explist
	return p
}

func InitEmptyExplistContext(p *ExplistContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_explist
}

func (*ExplistContext) IsExplistContext() {}

func NewExplistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExplistContext {
	var p = new(ExplistContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_explist

	return p
}

func (s *ExplistContext) GetParser() antlr.Parser { return s.parser }

func (s *ExplistContext) AllExp() []IExpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpContext); ok {
			len++
		}
	}

	tst := make([]IExpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpContext); ok {
			tst[i] = t.(IExpContext)
			i++
		}
	}

	return tst
}

func (s *ExplistContext) Exp(i int) IExpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *ExplistContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCOMMA)
}

func (s *ExplistContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, i)
}

func (s *ExplistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExplistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExplistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterExplist(s)
	}
}

func (s *ExplistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitExplist(s)
	}
}

func (s *ExplistContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitExplist(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Explist() (localctx IExplistContext) {
	localctx = NewExplistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, LuaParserRULE_explist)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(207)
		p.exp(0)
	}
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == LuaParserCOMMA {
		{
			p.SetState(208)
			p.Match(LuaParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(209)
			p.exp(0)
		}

		p.SetState(214)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpContext is an interface to support dynamic dispatch.
type IExpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NIL() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	Number() INumberContext
	String_() IStringContext
	DDD() antlr.TerminalNode
	Functiondef() IFunctiondefContext
	Prefixexp() IPrefixexpContext
	Tableconstructor() ITableconstructorContext
	AllExp() []IExpContext
	Exp(i int) IExpContext
	NOT() antlr.TerminalNode
	POUND() antlr.TerminalNode
	MINUS() antlr.TerminalNode
	SQUIG() antlr.TerminalNode
	CARET() antlr.TerminalNode
	STAR() antlr.TerminalNode
	SLASH() antlr.TerminalNode
	PER() antlr.TerminalNode
	SS() antlr.TerminalNode
	PLUS() antlr.TerminalNode
	DD() antlr.TerminalNode
	LT() antlr.TerminalNode
	GT() antlr.TerminalNode
	LE() antlr.TerminalNode
	GE() antlr.TerminalNode
	SQEQ() antlr.TerminalNode
	EE() antlr.TerminalNode
	AND() antlr.TerminalNode
	OR() antlr.TerminalNode
	AMP() antlr.TerminalNode
	PIPE() antlr.TerminalNode
	LL() antlr.TerminalNode
	GG() antlr.TerminalNode

	// IsExpContext differentiates from other interfaces.
	IsExpContext()
}

type ExpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpContext() *ExpContext {
	var p = new(ExpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_exp
	return p
}

func InitEmptyExpContext(p *ExpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_exp
}

func (*ExpContext) IsExpContext() {}

func NewExpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpContext {
	var p = new(ExpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_exp

	return p
}

func (s *ExpContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpContext) NIL() antlr.TerminalNode {
	return s.GetToken(LuaParserNIL, 0)
}

func (s *ExpContext) FALSE() antlr.TerminalNode {
	return s.GetToken(LuaParserFALSE, 0)
}

func (s *ExpContext) TRUE() antlr.TerminalNode {
	return s.GetToken(LuaParserTRUE, 0)
}

func (s *ExpContext) Number() INumberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumberContext)
}

func (s *ExpContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *ExpContext) DDD() antlr.TerminalNode {
	return s.GetToken(LuaParserDDD, 0)
}

func (s *ExpContext) Functiondef() IFunctiondefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctiondefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctiondefContext)
}

func (s *ExpContext) Prefixexp() IPrefixexpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrefixexpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrefixexpContext)
}

func (s *ExpContext) Tableconstructor() ITableconstructorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableconstructorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableconstructorContext)
}

func (s *ExpContext) AllExp() []IExpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpContext); ok {
			len++
		}
	}

	tst := make([]IExpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpContext); ok {
			tst[i] = t.(IExpContext)
			i++
		}
	}

	return tst
}

func (s *ExpContext) Exp(i int) IExpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *ExpContext) NOT() antlr.TerminalNode {
	return s.GetToken(LuaParserNOT, 0)
}

func (s *ExpContext) POUND() antlr.TerminalNode {
	return s.GetToken(LuaParserPOUND, 0)
}

func (s *ExpContext) MINUS() antlr.TerminalNode {
	return s.GetToken(LuaParserMINUS, 0)
}

func (s *ExpContext) SQUIG() antlr.TerminalNode {
	return s.GetToken(LuaParserSQUIG, 0)
}

func (s *ExpContext) CARET() antlr.TerminalNode {
	return s.GetToken(LuaParserCARET, 0)
}

func (s *ExpContext) STAR() antlr.TerminalNode {
	return s.GetToken(LuaParserSTAR, 0)
}

func (s *ExpContext) SLASH() antlr.TerminalNode {
	return s.GetToken(LuaParserSLASH, 0)
}

func (s *ExpContext) PER() antlr.TerminalNode {
	return s.GetToken(LuaParserPER, 0)
}

func (s *ExpContext) SS() antlr.TerminalNode {
	return s.GetToken(LuaParserSS, 0)
}

func (s *ExpContext) PLUS() antlr.TerminalNode {
	return s.GetToken(LuaParserPLUS, 0)
}

func (s *ExpContext) DD() antlr.TerminalNode {
	return s.GetToken(LuaParserDD, 0)
}

func (s *ExpContext) LT() antlr.TerminalNode {
	return s.GetToken(LuaParserLT, 0)
}

func (s *ExpContext) GT() antlr.TerminalNode {
	return s.GetToken(LuaParserGT, 0)
}

func (s *ExpContext) LE() antlr.TerminalNode {
	return s.GetToken(LuaParserLE, 0)
}

func (s *ExpContext) GE() antlr.TerminalNode {
	return s.GetToken(LuaParserGE, 0)
}

func (s *ExpContext) SQEQ() antlr.TerminalNode {
	return s.GetToken(LuaParserSQEQ, 0)
}

func (s *ExpContext) EE() antlr.TerminalNode {
	return s.GetToken(LuaParserEE, 0)
}

func (s *ExpContext) AND() antlr.TerminalNode {
	return s.GetToken(LuaParserAND, 0)
}

func (s *ExpContext) OR() antlr.TerminalNode {
	return s.GetToken(LuaParserOR, 0)
}

func (s *ExpContext) AMP() antlr.TerminalNode {
	return s.GetToken(LuaParserAMP, 0)
}

func (s *ExpContext) PIPE() antlr.TerminalNode {
	return s.GetToken(LuaParserPIPE, 0)
}

func (s *ExpContext) LL() antlr.TerminalNode {
	return s.GetToken(LuaParserLL, 0)
}

func (s *ExpContext) GG() antlr.TerminalNode {
	return s.GetToken(LuaParserGG, 0)
}

func (s *ExpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterExp(s)
	}
}

func (s *ExpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitExp(s)
	}
}

func (s *ExpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitExp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Exp() (localctx IExpContext) {
	return p.exp(0)
}

func (p *LuaParser) exp(_p int) (localctx IExpContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 24
	p.EnterRecursionRule(localctx, 24, LuaParserRULE_exp, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(227)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case LuaParserNIL:
		{
			p.SetState(216)
			p.Match(LuaParserNIL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserFALSE:
		{
			p.SetState(217)
			p.Match(LuaParserFALSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserTRUE:
		{
			p.SetState(218)
			p.Match(LuaParserTRUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserINT, LuaParserHEX, LuaParserFLOAT, LuaParserHEX_FLOAT:
		{
			p.SetState(219)
			p.Number()
		}

	case LuaParserNORMALSTRING, LuaParserCHARSTRING, LuaParserLONGSTRING:
		{
			p.SetState(220)
			p.String_()
		}

	case LuaParserDDD:
		{
			p.SetState(221)
			p.Match(LuaParserDDD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserFUNCTION:
		{
			p.SetState(222)
			p.Functiondef()
		}

	case LuaParserOP, LuaParserNAME:
		{
			p.SetState(223)
			p.Prefixexp()
		}

	case LuaParserOCU:
		{
			p.SetState(224)
			p.Tableconstructor()
		}

	case LuaParserSQUIG, LuaParserMINUS, LuaParserPOUND, LuaParserNOT:
		{
			p.SetState(225)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&10468982784) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(226)
			p.exp(8)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(255)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(253)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 18, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(229)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
					goto errorExit
				}

				{
					p.SetState(230)
					p.Match(LuaParserCARET)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(231)
					p.exp(9)
				}

			case 2:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(232)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
					goto errorExit
				}
				{
					p.SetState(233)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&18049995198431232) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(234)
					p.exp(8)
				}

			case 3:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(235)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
					goto errorExit
				}
				{
					p.SetState(236)
					_la = p.GetTokenStream().LA(1)

					if !(_la == LuaParserMINUS || _la == LuaParserPLUS) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(237)
					p.exp(7)
				}

			case 4:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(238)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
					goto errorExit
				}

				{
					p.SetState(239)
					p.Match(LuaParserDD)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(240)
					p.exp(5)
				}

			case 5:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(241)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(242)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&73186792481226752) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(243)
					p.exp(5)
				}

			case 6:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(244)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}

				{
					p.SetState(245)
					p.Match(LuaParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(246)
					p.exp(4)
				}

			case 7:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(247)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}

				{
					p.SetState(248)
					p.Match(LuaParserOR)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

				{
					p.SetState(249)
					p.exp(3)
				}

			case 8:
				localctx = NewExpContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_exp)
				p.SetState(250)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(251)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4503720154890240) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(252)
					p.exp(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(257)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IVarContext is an interface to support dynamic dispatch.
type IVarContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NAME() antlr.TerminalNode
	Prefixexp() IPrefixexpContext
	OB() antlr.TerminalNode
	Exp() IExpContext
	CB() antlr.TerminalNode
	DOT() antlr.TerminalNode

	// IsVarContext differentiates from other interfaces.
	IsVarContext()
}

type VarContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVarContext() *VarContext {
	var p = new(VarContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_var
	return p
}

func InitEmptyVarContext(p *VarContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_var
}

func (*VarContext) IsVarContext() {}

func NewVarContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarContext {
	var p = new(VarContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_var

	return p
}

func (s *VarContext) GetParser() antlr.Parser { return s.parser }

func (s *VarContext) NAME() antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, 0)
}

func (s *VarContext) Prefixexp() IPrefixexpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrefixexpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrefixexpContext)
}

func (s *VarContext) OB() antlr.TerminalNode {
	return s.GetToken(LuaParserOB, 0)
}

func (s *VarContext) Exp() IExpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *VarContext) CB() antlr.TerminalNode {
	return s.GetToken(LuaParserCB, 0)
}

func (s *VarContext) DOT() antlr.TerminalNode {
	return s.GetToken(LuaParserDOT, 0)
}

func (s *VarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VarContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterVar(s)
	}
}

func (s *VarContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitVar(s)
	}
}

func (s *VarContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitVar(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Var_() (localctx IVarContext) {
	localctx = NewVarContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, LuaParserRULE_var)
	p.SetState(268)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(258)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(259)
			p.Prefixexp()
		}
		p.SetState(266)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case LuaParserOB:
			{
				p.SetState(260)
				p.Match(LuaParserOB)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(261)
				p.exp(0)
			}
			{
				p.SetState(262)
				p.Match(LuaParserCB)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case LuaParserDOT:
			{
				p.SetState(264)
				p.Match(LuaParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(265)
				p.Match(LuaParserNAME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrefixexpContext is an interface to support dynamic dispatch.
type IPrefixexpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNAME() []antlr.TerminalNode
	NAME(i int) antlr.TerminalNode
	AllOB() []antlr.TerminalNode
	OB(i int) antlr.TerminalNode
	AllExp() []IExpContext
	Exp(i int) IExpContext
	AllCB() []antlr.TerminalNode
	CB(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	Functioncall() IFunctioncallContext
	OP() antlr.TerminalNode
	CP() antlr.TerminalNode

	// IsPrefixexpContext differentiates from other interfaces.
	IsPrefixexpContext()
}

type PrefixexpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrefixexpContext() *PrefixexpContext {
	var p = new(PrefixexpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_prefixexp
	return p
}

func InitEmptyPrefixexpContext(p *PrefixexpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_prefixexp
}

func (*PrefixexpContext) IsPrefixexpContext() {}

func NewPrefixexpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrefixexpContext {
	var p = new(PrefixexpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_prefixexp

	return p
}

func (s *PrefixexpContext) GetParser() antlr.Parser { return s.parser }

func (s *PrefixexpContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(LuaParserNAME)
}

func (s *PrefixexpContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, i)
}

func (s *PrefixexpContext) AllOB() []antlr.TerminalNode {
	return s.GetTokens(LuaParserOB)
}

func (s *PrefixexpContext) OB(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserOB, i)
}

func (s *PrefixexpContext) AllExp() []IExpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpContext); ok {
			len++
		}
	}

	tst := make([]IExpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpContext); ok {
			tst[i] = t.(IExpContext)
			i++
		}
	}

	return tst
}

func (s *PrefixexpContext) Exp(i int) IExpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *PrefixexpContext) AllCB() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCB)
}

func (s *PrefixexpContext) CB(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCB, i)
}

func (s *PrefixexpContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(LuaParserDOT)
}

func (s *PrefixexpContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserDOT, i)
}

func (s *PrefixexpContext) Functioncall() IFunctioncallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctioncallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctioncallContext)
}

func (s *PrefixexpContext) OP() antlr.TerminalNode {
	return s.GetToken(LuaParserOP, 0)
}

func (s *PrefixexpContext) CP() antlr.TerminalNode {
	return s.GetToken(LuaParserCP, 0)
}

func (s *PrefixexpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrefixexpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrefixexpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterPrefixexp(s)
	}
}

func (s *PrefixexpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitPrefixexp(s)
	}
}

func (s *PrefixexpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitPrefixexp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Prefixexp() (localctx IPrefixexpContext) {
	localctx = NewPrefixexpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, LuaParserRULE_prefixexp)
	var _alt int

	p.SetState(308)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(270)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(279)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(277)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}

				switch p.GetTokenStream().LA(1) {
				case LuaParserOB:
					{
						p.SetState(271)
						p.Match(LuaParserOB)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(272)
						p.exp(0)
					}
					{
						p.SetState(273)
						p.Match(LuaParserCB)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				case LuaParserDOT:
					{
						p.SetState(275)
						p.Match(LuaParserDOT)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(276)
						p.Match(LuaParserNAME)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				default:
					p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
					goto errorExit
				}

			}
			p.SetState(281)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(282)
			p.functioncall(0)
		}
		p.SetState(291)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(289)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}

				switch p.GetTokenStream().LA(1) {
				case LuaParserOB:
					{
						p.SetState(283)
						p.Match(LuaParserOB)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(284)
						p.exp(0)
					}
					{
						p.SetState(285)
						p.Match(LuaParserCB)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				case LuaParserDOT:
					{
						p.SetState(287)
						p.Match(LuaParserDOT)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(288)
						p.Match(LuaParserNAME)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				default:
					p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
					goto errorExit
				}

			}
			p.SetState(293)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 25, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(294)
			p.Match(LuaParserOP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(295)
			p.exp(0)
		}
		{
			p.SetState(296)
			p.Match(LuaParserCP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(305)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(303)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}

				switch p.GetTokenStream().LA(1) {
				case LuaParserOB:
					{
						p.SetState(297)
						p.Match(LuaParserOB)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(298)
						p.exp(0)
					}
					{
						p.SetState(299)
						p.Match(LuaParserCB)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				case LuaParserDOT:
					{
						p.SetState(301)
						p.Match(LuaParserDOT)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}
					{
						p.SetState(302)
						p.Match(LuaParserNAME)
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				default:
					p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
					goto errorExit
				}

			}
			p.SetState(307)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctioncallContext is an interface to support dynamic dispatch.
type IFunctioncallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNAME() []antlr.TerminalNode
	NAME(i int) antlr.TerminalNode
	Args() IArgsContext
	AllOB() []antlr.TerminalNode
	OB(i int) antlr.TerminalNode
	AllExp() []IExpContext
	Exp(i int) IExpContext
	AllCB() []antlr.TerminalNode
	CB(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	OP() antlr.TerminalNode
	CP() antlr.TerminalNode
	COL() antlr.TerminalNode
	Functioncall() IFunctioncallContext

	// IsFunctioncallContext differentiates from other interfaces.
	IsFunctioncallContext()
}

type FunctioncallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctioncallContext() *FunctioncallContext {
	var p = new(FunctioncallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_functioncall
	return p
}

func InitEmptyFunctioncallContext(p *FunctioncallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_functioncall
}

func (*FunctioncallContext) IsFunctioncallContext() {}

func NewFunctioncallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctioncallContext {
	var p = new(FunctioncallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_functioncall

	return p
}

func (s *FunctioncallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctioncallContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(LuaParserNAME)
}

func (s *FunctioncallContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, i)
}

func (s *FunctioncallContext) Args() IArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgsContext)
}

func (s *FunctioncallContext) AllOB() []antlr.TerminalNode {
	return s.GetTokens(LuaParserOB)
}

func (s *FunctioncallContext) OB(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserOB, i)
}

func (s *FunctioncallContext) AllExp() []IExpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpContext); ok {
			len++
		}
	}

	tst := make([]IExpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpContext); ok {
			tst[i] = t.(IExpContext)
			i++
		}
	}

	return tst
}

func (s *FunctioncallContext) Exp(i int) IExpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *FunctioncallContext) AllCB() []antlr.TerminalNode {
	return s.GetTokens(LuaParserCB)
}

func (s *FunctioncallContext) CB(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserCB, i)
}

func (s *FunctioncallContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(LuaParserDOT)
}

func (s *FunctioncallContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(LuaParserDOT, i)
}

func (s *FunctioncallContext) OP() antlr.TerminalNode {
	return s.GetToken(LuaParserOP, 0)
}

func (s *FunctioncallContext) CP() antlr.TerminalNode {
	return s.GetToken(LuaParserCP, 0)
}

func (s *FunctioncallContext) COL() antlr.TerminalNode {
	return s.GetToken(LuaParserCOL, 0)
}

func (s *FunctioncallContext) Functioncall() IFunctioncallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctioncallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctioncallContext)
}

func (s *FunctioncallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctioncallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctioncallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterFunctioncall(s)
	}
}

func (s *FunctioncallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitFunctioncall(s)
	}
}

func (s *FunctioncallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitFunctioncall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Functioncall() (localctx IFunctioncallContext) {
	return p.functioncall(0)
}

func (p *LuaParser) functioncall(_p int) (localctx IFunctioncallContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewFunctioncallContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IFunctioncallContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 30
	p.EnterRecursionRule(localctx, 30, LuaParserRULE_functioncall, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(311)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(320)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == LuaParserDOT || _la == LuaParserOB {
			p.SetState(318)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case LuaParserOB:
				{
					p.SetState(312)
					p.Match(LuaParserOB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(313)
					p.exp(0)
				}
				{
					p.SetState(314)
					p.Match(LuaParserCB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case LuaParserDOT:
				{
					p.SetState(316)
					p.Match(LuaParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(317)
					p.Match(LuaParserNAME)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(322)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(323)
			p.Args()
		}

	case 2:
		{
			p.SetState(324)
			p.Match(LuaParserOP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(325)
			p.exp(0)
		}
		{
			p.SetState(326)
			p.Match(LuaParserCP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(335)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == LuaParserDOT || _la == LuaParserOB {
			p.SetState(333)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case LuaParserOB:
				{
					p.SetState(327)
					p.Match(LuaParserOB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(328)
					p.exp(0)
				}
				{
					p.SetState(329)
					p.Match(LuaParserCB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case LuaParserDOT:
				{
					p.SetState(331)
					p.Match(LuaParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(332)
					p.Match(LuaParserNAME)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(337)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(338)
			p.Args()
		}

	case 3:
		{
			p.SetState(340)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(349)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == LuaParserDOT || _la == LuaParserOB {
			p.SetState(347)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case LuaParserOB:
				{
					p.SetState(341)
					p.Match(LuaParserOB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(342)
					p.exp(0)
				}
				{
					p.SetState(343)
					p.Match(LuaParserCB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case LuaParserDOT:
				{
					p.SetState(345)
					p.Match(LuaParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(346)
					p.Match(LuaParserNAME)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(351)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(352)
			p.Match(LuaParserCOL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(353)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(354)
			p.Args()
		}

	case 4:
		{
			p.SetState(355)
			p.Match(LuaParserOP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(356)
			p.exp(0)
		}
		{
			p.SetState(357)
			p.Match(LuaParserCP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(366)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == LuaParserDOT || _la == LuaParserOB {
			p.SetState(364)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case LuaParserOB:
				{
					p.SetState(358)
					p.Match(LuaParserOB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(359)
					p.exp(0)
				}
				{
					p.SetState(360)
					p.Match(LuaParserCB)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case LuaParserDOT:
				{
					p.SetState(362)
					p.Match(LuaParserDOT)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(363)
					p.Match(LuaParserNAME)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(368)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(369)
			p.Match(LuaParserCOL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(370)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(371)
			p.Args()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(405)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(403)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFunctioncallContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_functioncall)
				p.SetState(375)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
					goto errorExit
				}
				p.SetState(384)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				for _la == LuaParserDOT || _la == LuaParserOB {
					p.SetState(382)
					p.GetErrorHandler().Sync(p)
					if p.HasError() {
						goto errorExit
					}

					switch p.GetTokenStream().LA(1) {
					case LuaParserOB:
						{
							p.SetState(376)
							p.Match(LuaParserOB)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}
						{
							p.SetState(377)
							p.exp(0)
						}
						{
							p.SetState(378)
							p.Match(LuaParserCB)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}

					case LuaParserDOT:
						{
							p.SetState(380)
							p.Match(LuaParserDOT)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}
						{
							p.SetState(381)
							p.Match(LuaParserNAME)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}

					default:
						p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
						goto errorExit
					}

					p.SetState(386)
					p.GetErrorHandler().Sync(p)
					if p.HasError() {
						goto errorExit
					}
					_la = p.GetTokenStream().LA(1)
				}
				{
					p.SetState(387)
					p.Args()
				}

			case 2:
				localctx = NewFunctioncallContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, LuaParserRULE_functioncall)
				p.SetState(388)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				p.SetState(397)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				for _la == LuaParserDOT || _la == LuaParserOB {
					p.SetState(395)
					p.GetErrorHandler().Sync(p)
					if p.HasError() {
						goto errorExit
					}

					switch p.GetTokenStream().LA(1) {
					case LuaParserOB:
						{
							p.SetState(389)
							p.Match(LuaParserOB)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}
						{
							p.SetState(390)
							p.exp(0)
						}
						{
							p.SetState(391)
							p.Match(LuaParserCB)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}

					case LuaParserDOT:
						{
							p.SetState(393)
							p.Match(LuaParserDOT)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}
						{
							p.SetState(394)
							p.Match(LuaParserNAME)
							if p.HasError() {
								// Recognition error - abort rule
								goto errorExit
							}
						}

					default:
						p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
						goto errorExit
					}

					p.SetState(399)
					p.GetErrorHandler().Sync(p)
					if p.HasError() {
						goto errorExit
					}
					_la = p.GetTokenStream().LA(1)
				}
				{
					p.SetState(400)
					p.Match(LuaParserCOL)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(401)
					p.Match(LuaParserNAME)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(402)
					p.Args()
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(407)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArgsContext is an interface to support dynamic dispatch.
type IArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OP() antlr.TerminalNode
	CP() antlr.TerminalNode
	Explist() IExplistContext
	Tableconstructor() ITableconstructorContext
	String_() IStringContext

	// IsArgsContext differentiates from other interfaces.
	IsArgsContext()
}

type ArgsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgsContext() *ArgsContext {
	var p = new(ArgsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_args
	return p
}

func InitEmptyArgsContext(p *ArgsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_args
}

func (*ArgsContext) IsArgsContext() {}

func NewArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgsContext {
	var p = new(ArgsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_args

	return p
}

func (s *ArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgsContext) OP() antlr.TerminalNode {
	return s.GetToken(LuaParserOP, 0)
}

func (s *ArgsContext) CP() antlr.TerminalNode {
	return s.GetToken(LuaParserCP, 0)
}

func (s *ArgsContext) Explist() IExplistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExplistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExplistContext)
}

func (s *ArgsContext) Tableconstructor() ITableconstructorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableconstructorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableconstructorContext)
}

func (s *ArgsContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *ArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterArgs(s)
	}
}

func (s *ArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitArgs(s)
	}
}

func (s *ArgsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitArgs(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Args() (localctx IArgsContext) {
	localctx = NewArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, LuaParserRULE_args)
	var _la int

	p.SetState(415)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case LuaParserOP:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(408)
			p.Match(LuaParserOP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(410)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64((_la-17)) & ^0x3f) == 0 && ((int64(1)<<(_la-17))&280650879957889) != 0 {
			{
				p.SetState(409)
				p.Explist()
			}

		}
		{
			p.SetState(412)
			p.Match(LuaParserCP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserOCU:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(413)
			p.Tableconstructor()
		}

	case LuaParserNORMALSTRING, LuaParserCHARSTRING, LuaParserLONGSTRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(414)
			p.String_()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctiondefContext is an interface to support dynamic dispatch.
type IFunctiondefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNCTION() antlr.TerminalNode
	Funcbody() IFuncbodyContext

	// IsFunctiondefContext differentiates from other interfaces.
	IsFunctiondefContext()
}

type FunctiondefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctiondefContext() *FunctiondefContext {
	var p = new(FunctiondefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_functiondef
	return p
}

func InitEmptyFunctiondefContext(p *FunctiondefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_functiondef
}

func (*FunctiondefContext) IsFunctiondefContext() {}

func NewFunctiondefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctiondefContext {
	var p = new(FunctiondefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_functiondef

	return p
}

func (s *FunctiondefContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctiondefContext) FUNCTION() antlr.TerminalNode {
	return s.GetToken(LuaParserFUNCTION, 0)
}

func (s *FunctiondefContext) Funcbody() IFuncbodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFuncbodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFuncbodyContext)
}

func (s *FunctiondefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctiondefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctiondefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterFunctiondef(s)
	}
}

func (s *FunctiondefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitFunctiondef(s)
	}
}

func (s *FunctiondefContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitFunctiondef(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Functiondef() (localctx IFunctiondefContext) {
	localctx = NewFunctiondefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, LuaParserRULE_functiondef)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(417)
		p.Match(LuaParserFUNCTION)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(418)
		p.Funcbody()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFuncbodyContext is an interface to support dynamic dispatch.
type IFuncbodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OP() antlr.TerminalNode
	Parlist() IParlistContext
	CP() antlr.TerminalNode
	Block() IBlockContext
	END() antlr.TerminalNode

	// IsFuncbodyContext differentiates from other interfaces.
	IsFuncbodyContext()
}

type FuncbodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncbodyContext() *FuncbodyContext {
	var p = new(FuncbodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_funcbody
	return p
}

func InitEmptyFuncbodyContext(p *FuncbodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_funcbody
}

func (*FuncbodyContext) IsFuncbodyContext() {}

func NewFuncbodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncbodyContext {
	var p = new(FuncbodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_funcbody

	return p
}

func (s *FuncbodyContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncbodyContext) OP() antlr.TerminalNode {
	return s.GetToken(LuaParserOP, 0)
}

func (s *FuncbodyContext) Parlist() IParlistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParlistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParlistContext)
}

func (s *FuncbodyContext) CP() antlr.TerminalNode {
	return s.GetToken(LuaParserCP, 0)
}

func (s *FuncbodyContext) Block() IBlockContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBlockContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBlockContext)
}

func (s *FuncbodyContext) END() antlr.TerminalNode {
	return s.GetToken(LuaParserEND, 0)
}

func (s *FuncbodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncbodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FuncbodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterFuncbody(s)
	}
}

func (s *FuncbodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitFuncbody(s)
	}
}

func (s *FuncbodyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitFuncbody(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Funcbody() (localctx IFuncbodyContext) {
	localctx = NewFuncbodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, LuaParserRULE_funcbody)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(420)
		p.Match(LuaParserOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(421)
		p.Parlist()
	}
	{
		p.SetState(422)
		p.Match(LuaParserCP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(423)
		p.Block()
	}
	{
		p.SetState(424)
		p.Match(LuaParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParlistContext is an interface to support dynamic dispatch.
type IParlistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Namelist() INamelistContext
	COMMA() antlr.TerminalNode
	DDD() antlr.TerminalNode

	// IsParlistContext differentiates from other interfaces.
	IsParlistContext()
}

type ParlistContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParlistContext() *ParlistContext {
	var p = new(ParlistContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_parlist
	return p
}

func InitEmptyParlistContext(p *ParlistContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_parlist
}

func (*ParlistContext) IsParlistContext() {}

func NewParlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParlistContext {
	var p = new(ParlistContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_parlist

	return p
}

func (s *ParlistContext) GetParser() antlr.Parser { return s.parser }

func (s *ParlistContext) Namelist() INamelistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamelistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamelistContext)
}

func (s *ParlistContext) COMMA() antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, 0)
}

func (s *ParlistContext) DDD() antlr.TerminalNode {
	return s.GetToken(LuaParserDDD, 0)
}

func (s *ParlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterParlist(s)
	}
}

func (s *ParlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitParlist(s)
	}
}

func (s *ParlistContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitParlist(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Parlist() (localctx IParlistContext) {
	localctx = NewParlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, LuaParserRULE_parlist)
	var _la int

	p.SetState(433)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case LuaParserNAME:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(426)
			p.Namelist()
		}
		p.SetState(429)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == LuaParserCOMMA {
			{
				p.SetState(427)
				p.Match(LuaParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(428)
				p.Match(LuaParserDDD)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case LuaParserDDD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(431)
			p.Match(LuaParserDDD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case LuaParserCP:
		p.EnterOuterAlt(localctx, 3)

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITableconstructorContext is an interface to support dynamic dispatch.
type ITableconstructorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OCU() antlr.TerminalNode
	CCU() antlr.TerminalNode
	Fieldlist() IFieldlistContext

	// IsTableconstructorContext differentiates from other interfaces.
	IsTableconstructorContext()
}

type TableconstructorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTableconstructorContext() *TableconstructorContext {
	var p = new(TableconstructorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_tableconstructor
	return p
}

func InitEmptyTableconstructorContext(p *TableconstructorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_tableconstructor
}

func (*TableconstructorContext) IsTableconstructorContext() {}

func NewTableconstructorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TableconstructorContext {
	var p = new(TableconstructorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_tableconstructor

	return p
}

func (s *TableconstructorContext) GetParser() antlr.Parser { return s.parser }

func (s *TableconstructorContext) OCU() antlr.TerminalNode {
	return s.GetToken(LuaParserOCU, 0)
}

func (s *TableconstructorContext) CCU() antlr.TerminalNode {
	return s.GetToken(LuaParserCCU, 0)
}

func (s *TableconstructorContext) Fieldlist() IFieldlistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldlistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldlistContext)
}

func (s *TableconstructorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableconstructorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TableconstructorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterTableconstructor(s)
	}
}

func (s *TableconstructorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitTableconstructor(s)
	}
}

func (s *TableconstructorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitTableconstructor(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Tableconstructor() (localctx ITableconstructorContext) {
	localctx = NewTableconstructorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, LuaParserRULE_tableconstructor)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(435)
		p.Match(LuaParserOCU)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(437)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64((_la-17)) & ^0x3f) == 0 && ((int64(1)<<(_la-17))&280653027441537) != 0 {
		{
			p.SetState(436)
			p.Fieldlist()
		}

	}
	{
		p.SetState(439)
		p.Match(LuaParserCCU)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldlistContext is an interface to support dynamic dispatch.
type IFieldlistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllField() []IFieldContext
	Field(i int) IFieldContext
	AllFieldsep() []IFieldsepContext
	Fieldsep(i int) IFieldsepContext

	// IsFieldlistContext differentiates from other interfaces.
	IsFieldlistContext()
}

type FieldlistContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldlistContext() *FieldlistContext {
	var p = new(FieldlistContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_fieldlist
	return p
}

func InitEmptyFieldlistContext(p *FieldlistContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_fieldlist
}

func (*FieldlistContext) IsFieldlistContext() {}

func NewFieldlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldlistContext {
	var p = new(FieldlistContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_fieldlist

	return p
}

func (s *FieldlistContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldlistContext) AllField() []IFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldContext); ok {
			len++
		}
	}

	tst := make([]IFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldContext); ok {
			tst[i] = t.(IFieldContext)
			i++
		}
	}

	return tst
}

func (s *FieldlistContext) Field(i int) IFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *FieldlistContext) AllFieldsep() []IFieldsepContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldsepContext); ok {
			len++
		}
	}

	tst := make([]IFieldsepContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldsepContext); ok {
			tst[i] = t.(IFieldsepContext)
			i++
		}
	}

	return tst
}

func (s *FieldlistContext) Fieldsep(i int) IFieldsepContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldsepContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldsepContext)
}

func (s *FieldlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterFieldlist(s)
	}
}

func (s *FieldlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitFieldlist(s)
	}
}

func (s *FieldlistContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitFieldlist(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Fieldlist() (localctx IFieldlistContext) {
	localctx = NewFieldlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, LuaParserRULE_fieldlist)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(441)
		p.Field()
	}
	p.SetState(447)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 49, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(442)
				p.Fieldsep()
			}
			{
				p.SetState(443)
				p.Field()
			}

		}
		p.SetState(449)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 49, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(451)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == LuaParserSEMI || _la == LuaParserCOMMA {
		{
			p.SetState(450)
			p.Fieldsep()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OB() antlr.TerminalNode
	AllExp() []IExpContext
	Exp(i int) IExpContext
	CB() antlr.TerminalNode
	EQ() antlr.TerminalNode
	NAME() antlr.TerminalNode

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_field
	return p
}

func InitEmptyFieldContext(p *FieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_field
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) OB() antlr.TerminalNode {
	return s.GetToken(LuaParserOB, 0)
}

func (s *FieldContext) AllExp() []IExpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpContext); ok {
			len++
		}
	}

	tst := make([]IExpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpContext); ok {
			tst[i] = t.(IExpContext)
			i++
		}
	}

	return tst
}

func (s *FieldContext) Exp(i int) IExpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpContext)
}

func (s *FieldContext) CB() antlr.TerminalNode {
	return s.GetToken(LuaParserCB, 0)
}

func (s *FieldContext) EQ() antlr.TerminalNode {
	return s.GetToken(LuaParserEQ, 0)
}

func (s *FieldContext) NAME() antlr.TerminalNode {
	return s.GetToken(LuaParserNAME, 0)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterField(s)
	}
}

func (s *FieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitField(s)
	}
}

func (s *FieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Field() (localctx IFieldContext) {
	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, LuaParserRULE_field)
	p.SetState(463)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(453)
			p.Match(LuaParserOB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(454)
			p.exp(0)
		}
		{
			p.SetState(455)
			p.Match(LuaParserCB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(456)
			p.Match(LuaParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(457)
			p.exp(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(459)
			p.Match(LuaParserNAME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(460)
			p.Match(LuaParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(461)
			p.exp(0)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(462)
			p.exp(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldsepContext is an interface to support dynamic dispatch.
type IFieldsepContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMA() antlr.TerminalNode
	SEMI() antlr.TerminalNode

	// IsFieldsepContext differentiates from other interfaces.
	IsFieldsepContext()
}

type FieldsepContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldsepContext() *FieldsepContext {
	var p = new(FieldsepContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_fieldsep
	return p
}

func InitEmptyFieldsepContext(p *FieldsepContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_fieldsep
}

func (*FieldsepContext) IsFieldsepContext() {}

func NewFieldsepContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldsepContext {
	var p = new(FieldsepContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_fieldsep

	return p
}

func (s *FieldsepContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldsepContext) COMMA() antlr.TerminalNode {
	return s.GetToken(LuaParserCOMMA, 0)
}

func (s *FieldsepContext) SEMI() antlr.TerminalNode {
	return s.GetToken(LuaParserSEMI, 0)
}

func (s *FieldsepContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldsepContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldsepContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterFieldsep(s)
	}
}

func (s *FieldsepContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitFieldsep(s)
	}
}

func (s *FieldsepContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitFieldsep(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Fieldsep() (localctx IFieldsepContext) {
	localctx = NewFieldsepContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, LuaParserRULE_fieldsep)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(465)
		_la = p.GetTokenStream().LA(1)

		if !(_la == LuaParserSEMI || _la == LuaParserCOMMA) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INumberContext is an interface to support dynamic dispatch.
type INumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT() antlr.TerminalNode
	HEX() antlr.TerminalNode
	FLOAT() antlr.TerminalNode
	HEX_FLOAT() antlr.TerminalNode

	// IsNumberContext differentiates from other interfaces.
	IsNumberContext()
}

type NumberContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumberContext() *NumberContext {
	var p = new(NumberContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_number
	return p
}

func InitEmptyNumberContext(p *NumberContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_number
}

func (*NumberContext) IsNumberContext() {}

func NewNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumberContext {
	var p = new(NumberContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_number

	return p
}

func (s *NumberContext) GetParser() antlr.Parser { return s.parser }

func (s *NumberContext) INT() antlr.TerminalNode {
	return s.GetToken(LuaParserINT, 0)
}

func (s *NumberContext) HEX() antlr.TerminalNode {
	return s.GetToken(LuaParserHEX, 0)
}

func (s *NumberContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(LuaParserFLOAT, 0)
}

func (s *NumberContext) HEX_FLOAT() antlr.TerminalNode {
	return s.GetToken(LuaParserHEX_FLOAT, 0)
}

func (s *NumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterNumber(s)
	}
}

func (s *NumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitNumber(s)
	}
}

func (s *NumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitNumber(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) Number() (localctx INumberContext) {
	localctx = NewNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, LuaParserRULE_number)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(467)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-61)) & ^0x3f) == 0 && ((int64(1)<<(_la-61))&15) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStringContext is an interface to support dynamic dispatch.
type IStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NORMALSTRING() antlr.TerminalNode
	CHARSTRING() antlr.TerminalNode
	LONGSTRING() antlr.TerminalNode

	// IsStringContext differentiates from other interfaces.
	IsStringContext()
}

type StringContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringContext() *StringContext {
	var p = new(StringContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_string
	return p
}

func InitEmptyStringContext(p *StringContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = LuaParserRULE_string
}

func (*StringContext) IsStringContext() {}

func NewStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringContext {
	var p = new(StringContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = LuaParserRULE_string

	return p
}

func (s *StringContext) GetParser() antlr.Parser { return s.parser }

func (s *StringContext) NORMALSTRING() antlr.TerminalNode {
	return s.GetToken(LuaParserNORMALSTRING, 0)
}

func (s *StringContext) CHARSTRING() antlr.TerminalNode {
	return s.GetToken(LuaParserCHARSTRING, 0)
}

func (s *StringContext) LONGSTRING() antlr.TerminalNode {
	return s.GetToken(LuaParserLONGSTRING, 0)
}

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(LuaParserListener); ok {
		listenerT.ExitString(s)
	}
}

func (s *StringContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case LuaParserVisitor:
		return t.VisitString(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *LuaParser) String_() (localctx IStringContext) {
	localctx = NewStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, LuaParserRULE_string)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(469)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2017612633061982208) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *LuaParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 12:
		var t *ExpContext = nil
		if localctx != nil {
			t = localctx.(*ExpContext)
		}
		return p.Exp_Sempred(t, predIndex)

	case 15:
		var t *FunctioncallContext = nil
		if localctx != nil {
			t = localctx.(*FunctioncallContext)
		}
		return p.Functioncall_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *LuaParser) Exp_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *LuaParser) Functioncall_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 8:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
