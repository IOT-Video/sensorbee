package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint16

const (
	ruleUnknown pegRule = iota
	ruleSingleStatement
	ruleStatementWithRest
	ruleStatementWithoutRest
	ruleStatement
	ruleSourceStmt
	ruleSinkStmt
	ruleStateStmt
	ruleStreamStmt
	ruleSelectStmt
	ruleSelectUnionStmt
	ruleCreateStreamAsSelectStmt
	ruleCreateStreamAsSelectUnionStmt
	ruleCreateSourceStmt
	ruleCreateSinkStmt
	ruleCreateStateStmt
	ruleUpdateStateStmt
	ruleUpdateSourceStmt
	ruleUpdateSinkStmt
	ruleInsertIntoFromStmt
	rulePauseSourceStmt
	ruleResumeSourceStmt
	ruleRewindSourceStmt
	ruleDropSourceStmt
	ruleDropStreamStmt
	ruleDropSinkStmt
	ruleDropStateStmt
	ruleLoadStateStmt
	ruleLoadStateOrCreateStmt
	ruleSaveStateStmt
	ruleEvalStmt
	ruleEmitter
	ruleEmitterOptions
	ruleEmitterOptionCombinations
	ruleEmitterLimit
	ruleEmitterSample
	ruleCountBasedSampling
	ruleRandomizedSampling
	ruleTimeBasedSampling
	ruleTimeBasedSamplingSeconds
	ruleTimeBasedSamplingMilliseconds
	ruleProjections
	ruleProjection
	ruleAliasExpression
	ruleWindowedFrom
	ruleInterval
	ruleTimeInterval
	ruleTuplesInterval
	ruleRelations
	ruleFilter
	ruleGrouping
	ruleGroupList
	ruleHaving
	ruleRelationLike
	ruleAliasedStreamWindow
	ruleStreamWindow
	ruleStreamLike
	ruleUDSFFuncApp
	ruleCapacitySpecOpt
	ruleSheddingSpecOpt
	ruleSheddingOption
	ruleSourceSinkSpecs
	ruleUpdateSourceSinkSpecs
	ruleSetOptSpecs
	ruleStateTagOpt
	ruleSourceSinkParam
	ruleSourceSinkParamVal
	ruleParamLiteral
	ruleParamArrayExpr
	ruleParamMapExpr
	ruleParamKeyValuePair
	rulePausedOpt
	ruleExpressionOrWildcard
	ruleExpression
	ruleorExpr
	ruleandExpr
	rulenotExpr
	rulecomparisonExpr
	ruleotherOpExpr
	ruleisExpr
	ruletermExpr
	ruleproductExpr
	ruleminusExpr
	rulecastExpr
	rulebaseExpr
	ruleFuncTypeCast
	ruleFuncApp
	ruleFuncAppWithOrderBy
	ruleFuncAppWithoutOrderBy
	ruleFuncParams
	ruleParamsOrder
	ruleSortedExpression
	ruleOrderDirectionOpt
	ruleArrayExpr
	ruleMapExpr
	ruleKeyValuePair
	ruleCase
	ruleConditionCase
	ruleExpressionCase
	ruleWhenThenPair
	ruleLiteral
	ruleComparisonOp
	ruleOtherOp
	ruleIsOp
	rulePlusMinusOp
	ruleMultDivOp
	ruleStream
	ruleRowMeta
	ruleRowTimestamp
	ruleRowValue
	ruleNumericLiteral
	ruleNonNegativeNumericLiteral
	ruleFloatLiteral
	ruleFunction
	ruleNullLiteral
	ruleBooleanLiteral
	ruleTRUE
	ruleFALSE
	ruleWildcard
	ruleStringLiteral
	ruleISTREAM
	ruleDSTREAM
	ruleRSTREAM
	ruleTUPLES
	ruleSECONDS
	ruleMILLISECONDS
	ruleWait
	ruleDropOldest
	ruleDropNewest
	ruleStreamIdentifier
	ruleSourceSinkType
	ruleSourceSinkParamKey
	rulePaused
	ruleUnpaused
	ruleAscending
	ruleDescending
	ruleType
	ruleBool
	ruleInt
	ruleFloat
	ruleString
	ruleBlob
	ruleTimestamp
	ruleArray
	ruleMap
	ruleOr
	ruleAnd
	ruleNot
	ruleEqual
	ruleLess
	ruleLessOrEqual
	ruleGreater
	ruleGreaterOrEqual
	ruleNotEqual
	ruleConcat
	ruleIs
	ruleIsNot
	rulePlus
	ruleMinus
	ruleMultiply
	ruleDivide
	ruleModulo
	ruleUnaryMinus
	ruleIdentifier
	ruleTargetIdentifier
	ruleident
	rulejsonGetPath
	rulejsonSetPath
	rulejsonPathHead
	rulejsonGetPathNonHead
	rulejsonSetPathNonHead
	rulejsonMapSingleLevel
	rulejsonMapMultipleLevel
	rulejsonMapAccessString
	rulejsonMapAccessBracket
	ruledoubleQuotedString
	rulejsonArrayAccess
	rulejsonNonNegativeArrayAccess
	rulejsonArraySlice
	rulejsonArrayPartialSlice
	rulejsonArrayFullSlice
	rulespElem
	rulesp
	rulespOpt
	rulecomment
	rulefinalComment
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
	ruleAction39
	ruleAction40
	ruleAction41
	ruleAction42
	ruleAction43
	ruleAction44
	ruleAction45
	ruleAction46
	ruleAction47
	ruleAction48
	ruleAction49
	ruleAction50
	ruleAction51
	ruleAction52
	ruleAction53
	ruleAction54
	ruleAction55
	ruleAction56
	ruleAction57
	ruleAction58
	ruleAction59
	ruleAction60
	ruleAction61
	ruleAction62
	ruleAction63
	ruleAction64
	ruleAction65
	ruleAction66
	ruleAction67
	ruleAction68
	ruleAction69
	ruleAction70
	ruleAction71
	ruleAction72
	ruleAction73
	ruleAction74
	ruleAction75
	ruleAction76
	ruleAction77
	ruleAction78
	ruleAction79
	ruleAction80
	ruleAction81
	ruleAction82
	ruleAction83
	ruleAction84
	ruleAction85
	ruleAction86
	ruleAction87
	ruleAction88
	ruleAction89
	ruleAction90
	ruleAction91
	ruleAction92
	ruleAction93
	ruleAction94
	ruleAction95
	ruleAction96
	ruleAction97
	ruleAction98
	ruleAction99
	ruleAction100
	ruleAction101
	ruleAction102
	ruleAction103
	ruleAction104
	ruleAction105
	ruleAction106
	ruleAction107
	ruleAction108
	ruleAction109
	ruleAction110
	ruleAction111
	ruleAction112
	ruleAction113
	ruleAction114
	ruleAction115
	ruleAction116
	ruleAction117
	ruleAction118
	ruleAction119
	ruleAction120
	ruleAction121
	ruleAction122
	ruleAction123
	ruleAction124
	ruleAction125
	ruleAction126
	ruleAction127
	ruleAction128
	ruleAction129
	ruleAction130
	ruleAction131
	ruleAction132

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"SingleStatement",
	"StatementWithRest",
	"StatementWithoutRest",
	"Statement",
	"SourceStmt",
	"SinkStmt",
	"StateStmt",
	"StreamStmt",
	"SelectStmt",
	"SelectUnionStmt",
	"CreateStreamAsSelectStmt",
	"CreateStreamAsSelectUnionStmt",
	"CreateSourceStmt",
	"CreateSinkStmt",
	"CreateStateStmt",
	"UpdateStateStmt",
	"UpdateSourceStmt",
	"UpdateSinkStmt",
	"InsertIntoFromStmt",
	"PauseSourceStmt",
	"ResumeSourceStmt",
	"RewindSourceStmt",
	"DropSourceStmt",
	"DropStreamStmt",
	"DropSinkStmt",
	"DropStateStmt",
	"LoadStateStmt",
	"LoadStateOrCreateStmt",
	"SaveStateStmt",
	"EvalStmt",
	"Emitter",
	"EmitterOptions",
	"EmitterOptionCombinations",
	"EmitterLimit",
	"EmitterSample",
	"CountBasedSampling",
	"RandomizedSampling",
	"TimeBasedSampling",
	"TimeBasedSamplingSeconds",
	"TimeBasedSamplingMilliseconds",
	"Projections",
	"Projection",
	"AliasExpression",
	"WindowedFrom",
	"Interval",
	"TimeInterval",
	"TuplesInterval",
	"Relations",
	"Filter",
	"Grouping",
	"GroupList",
	"Having",
	"RelationLike",
	"AliasedStreamWindow",
	"StreamWindow",
	"StreamLike",
	"UDSFFuncApp",
	"CapacitySpecOpt",
	"SheddingSpecOpt",
	"SheddingOption",
	"SourceSinkSpecs",
	"UpdateSourceSinkSpecs",
	"SetOptSpecs",
	"StateTagOpt",
	"SourceSinkParam",
	"SourceSinkParamVal",
	"ParamLiteral",
	"ParamArrayExpr",
	"ParamMapExpr",
	"ParamKeyValuePair",
	"PausedOpt",
	"ExpressionOrWildcard",
	"Expression",
	"orExpr",
	"andExpr",
	"notExpr",
	"comparisonExpr",
	"otherOpExpr",
	"isExpr",
	"termExpr",
	"productExpr",
	"minusExpr",
	"castExpr",
	"baseExpr",
	"FuncTypeCast",
	"FuncApp",
	"FuncAppWithOrderBy",
	"FuncAppWithoutOrderBy",
	"FuncParams",
	"ParamsOrder",
	"SortedExpression",
	"OrderDirectionOpt",
	"ArrayExpr",
	"MapExpr",
	"KeyValuePair",
	"Case",
	"ConditionCase",
	"ExpressionCase",
	"WhenThenPair",
	"Literal",
	"ComparisonOp",
	"OtherOp",
	"IsOp",
	"PlusMinusOp",
	"MultDivOp",
	"Stream",
	"RowMeta",
	"RowTimestamp",
	"RowValue",
	"NumericLiteral",
	"NonNegativeNumericLiteral",
	"FloatLiteral",
	"Function",
	"NullLiteral",
	"BooleanLiteral",
	"TRUE",
	"FALSE",
	"Wildcard",
	"StringLiteral",
	"ISTREAM",
	"DSTREAM",
	"RSTREAM",
	"TUPLES",
	"SECONDS",
	"MILLISECONDS",
	"Wait",
	"DropOldest",
	"DropNewest",
	"StreamIdentifier",
	"SourceSinkType",
	"SourceSinkParamKey",
	"Paused",
	"Unpaused",
	"Ascending",
	"Descending",
	"Type",
	"Bool",
	"Int",
	"Float",
	"String",
	"Blob",
	"Timestamp",
	"Array",
	"Map",
	"Or",
	"And",
	"Not",
	"Equal",
	"Less",
	"LessOrEqual",
	"Greater",
	"GreaterOrEqual",
	"NotEqual",
	"Concat",
	"Is",
	"IsNot",
	"Plus",
	"Minus",
	"Multiply",
	"Divide",
	"Modulo",
	"UnaryMinus",
	"Identifier",
	"TargetIdentifier",
	"ident",
	"jsonGetPath",
	"jsonSetPath",
	"jsonPathHead",
	"jsonGetPathNonHead",
	"jsonSetPathNonHead",
	"jsonMapSingleLevel",
	"jsonMapMultipleLevel",
	"jsonMapAccessString",
	"jsonMapAccessBracket",
	"doubleQuotedString",
	"jsonArrayAccess",
	"jsonNonNegativeArrayAccess",
	"jsonArraySlice",
	"jsonArrayPartialSlice",
	"jsonArrayFullSlice",
	"spElem",
	"sp",
	"spOpt",
	"comment",
	"finalComment",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",
	"Action40",
	"Action41",
	"Action42",
	"Action43",
	"Action44",
	"Action45",
	"Action46",
	"Action47",
	"Action48",
	"Action49",
	"Action50",
	"Action51",
	"Action52",
	"Action53",
	"Action54",
	"Action55",
	"Action56",
	"Action57",
	"Action58",
	"Action59",
	"Action60",
	"Action61",
	"Action62",
	"Action63",
	"Action64",
	"Action65",
	"Action66",
	"Action67",
	"Action68",
	"Action69",
	"Action70",
	"Action71",
	"Action72",
	"Action73",
	"Action74",
	"Action75",
	"Action76",
	"Action77",
	"Action78",
	"Action79",
	"Action80",
	"Action81",
	"Action82",
	"Action83",
	"Action84",
	"Action85",
	"Action86",
	"Action87",
	"Action88",
	"Action89",
	"Action90",
	"Action91",
	"Action92",
	"Action93",
	"Action94",
	"Action95",
	"Action96",
	"Action97",
	"Action98",
	"Action99",
	"Action100",
	"Action101",
	"Action102",
	"Action103",
	"Action104",
	"Action105",
	"Action106",
	"Action107",
	"Action108",
	"Action109",
	"Action110",
	"Action111",
	"Action112",
	"Action113",
	"Action114",
	"Action115",
	"Action116",
	"Action117",
	"Action118",
	"Action119",
	"Action120",
	"Action121",
	"Action122",
	"Action123",
	"Action124",
	"Action125",
	"Action126",
	"Action127",
	"Action128",
	"Action129",
	"Action130",
	"Action131",
	"Action132",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next uint32, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/*func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2 * len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}*/

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type bqlPegBackend struct {
	parseStack

	Buffer string
	buffer []rune
	rules  [320]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range []rune(buffer) {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *bqlPegBackend
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *bqlPegBackend) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *bqlPegBackend) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *bqlPegBackend) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.IncludeTrailingWhitespace(begin, end)

		case ruleAction1:

			p.IncludeTrailingWhitespace(begin, end)

		case ruleAction2:

			p.AssembleSelect()

		case ruleAction3:

			p.AssembleSelectUnion(begin, end)

		case ruleAction4:

			p.AssembleCreateStreamAsSelect()

		case ruleAction5:

			p.AssembleCreateStreamAsSelectUnion()

		case ruleAction6:

			p.AssembleCreateSource()

		case ruleAction7:

			p.AssembleCreateSink()

		case ruleAction8:

			p.AssembleCreateState()

		case ruleAction9:

			p.AssembleUpdateState()

		case ruleAction10:

			p.AssembleUpdateSource()

		case ruleAction11:

			p.AssembleUpdateSink()

		case ruleAction12:

			p.AssembleInsertIntoFrom()

		case ruleAction13:

			p.AssemblePauseSource()

		case ruleAction14:

			p.AssembleResumeSource()

		case ruleAction15:

			p.AssembleRewindSource()

		case ruleAction16:

			p.AssembleDropSource()

		case ruleAction17:

			p.AssembleDropStream()

		case ruleAction18:

			p.AssembleDropSink()

		case ruleAction19:

			p.AssembleDropState()

		case ruleAction20:

			p.AssembleLoadState()

		case ruleAction21:

			p.AssembleLoadStateOrCreate()

		case ruleAction22:

			p.AssembleSaveState()

		case ruleAction23:

			p.AssembleEval(begin, end)

		case ruleAction24:

			p.AssembleEmitter()

		case ruleAction25:

			p.AssembleEmitterOptions(begin, end)

		case ruleAction26:

			p.AssembleEmitterLimit()

		case ruleAction27:

			p.AssembleEmitterSampling(CountBasedSampling, 1)

		case ruleAction28:

			p.AssembleEmitterSampling(RandomizedSampling, 1)

		case ruleAction29:

			p.AssembleEmitterSampling(TimeBasedSampling, 1)

		case ruleAction30:

			p.AssembleEmitterSampling(TimeBasedSampling, 0.001)

		case ruleAction31:

			p.AssembleProjections(begin, end)

		case ruleAction32:

			p.AssembleAlias()

		case ruleAction33:

			// This is *always* executed, even if there is no
			// FROM clause present in the statement.
			p.AssembleWindowedFrom(begin, end)

		case ruleAction34:

			p.AssembleInterval()

		case ruleAction35:

			p.AssembleInterval()

		case ruleAction36:

			// This is *always* executed, even if there is no
			// WHERE clause present in the statement.
			p.AssembleFilter(begin, end)

		case ruleAction37:

			// This is *always* executed, even if there is no
			// GROUP BY clause present in the statement.
			p.AssembleGrouping(begin, end)

		case ruleAction38:

			// This is *always* executed, even if there is no
			// HAVING clause present in the statement.
			p.AssembleHaving(begin, end)

		case ruleAction39:

			p.EnsureAliasedStreamWindow()

		case ruleAction40:

			p.AssembleAliasedStreamWindow()

		case ruleAction41:

			p.AssembleStreamWindow()

		case ruleAction42:

			p.AssembleUDSFFuncApp()

		case ruleAction43:

			p.EnsureCapacitySpec(begin, end)

		case ruleAction44:

			p.EnsureSheddingSpec(begin, end)

		case ruleAction45:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction46:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction47:

			p.AssembleSourceSinkSpecs(begin, end)

		case ruleAction48:

			p.EnsureIdentifier(begin, end)

		case ruleAction49:

			p.AssembleSourceSinkParam()

		case ruleAction50:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction51:

			p.AssembleMap(begin, end)

		case ruleAction52:

			p.AssembleKeyValuePair()

		case ruleAction53:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction54:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction55:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction56:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction57:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction58:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction59:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction60:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction61:

			p.AssembleBinaryOperation(begin, end)

		case ruleAction62:

			p.AssembleUnaryPrefixOperation(begin, end)

		case ruleAction63:

			p.AssembleTypeCast(begin, end)

		case ruleAction64:

			p.AssembleTypeCast(begin, end)

		case ruleAction65:

			p.AssembleFuncApp()

		case ruleAction66:

			p.AssembleExpressions(begin, end)
			p.AssembleFuncApp()

		case ruleAction67:

			p.AssembleExpressions(begin, end)

		case ruleAction68:

			p.AssembleExpressions(begin, end)

		case ruleAction69:

			p.AssembleSortedExpression()

		case ruleAction70:

			p.EnsureKeywordPresent(begin, end)

		case ruleAction71:

			p.AssembleExpressions(begin, end)
			p.AssembleArray()

		case ruleAction72:

			p.AssembleMap(begin, end)

		case ruleAction73:

			p.AssembleKeyValuePair()

		case ruleAction74:

			p.AssembleConditionCase(begin, end)

		case ruleAction75:

			p.AssembleExpressionCase(begin, end)

		case ruleAction76:

			p.AssembleWhenThenPair()

		case ruleAction77:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStream(substr))

		case ruleAction78:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))

		case ruleAction79:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewRowValue(substr))

		case ruleAction80:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction81:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewNumericLiteral(substr))

		case ruleAction82:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewFloatLiteral(substr))

		case ruleAction83:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, FuncName(substr))

		case ruleAction84:

			p.PushComponent(begin, end, NewNullLiteral())

		case ruleAction85:

			p.PushComponent(begin, end, NewBoolLiteral(true))

		case ruleAction86:

			p.PushComponent(begin, end, NewBoolLiteral(false))

		case ruleAction87:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewWildcard(substr))

		case ruleAction88:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, NewStringLiteral(substr))

		case ruleAction89:

			p.PushComponent(begin, end, Istream)

		case ruleAction90:

			p.PushComponent(begin, end, Dstream)

		case ruleAction91:

			p.PushComponent(begin, end, Rstream)

		case ruleAction92:

			p.PushComponent(begin, end, Tuples)

		case ruleAction93:

			p.PushComponent(begin, end, Seconds)

		case ruleAction94:

			p.PushComponent(begin, end, Milliseconds)

		case ruleAction95:

			p.PushComponent(begin, end, Wait)

		case ruleAction96:

			p.PushComponent(begin, end, DropOldest)

		case ruleAction97:

			p.PushComponent(begin, end, DropNewest)

		case ruleAction98:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, StreamIdentifier(substr))

		case ruleAction99:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkType(substr))

		case ruleAction100:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, SourceSinkParamKey(substr))

		case ruleAction101:

			p.PushComponent(begin, end, Yes)

		case ruleAction102:

			p.PushComponent(begin, end, No)

		case ruleAction103:

			p.PushComponent(begin, end, Yes)

		case ruleAction104:

			p.PushComponent(begin, end, No)

		case ruleAction105:

			p.PushComponent(begin, end, Bool)

		case ruleAction106:

			p.PushComponent(begin, end, Int)

		case ruleAction107:

			p.PushComponent(begin, end, Float)

		case ruleAction108:

			p.PushComponent(begin, end, String)

		case ruleAction109:

			p.PushComponent(begin, end, Blob)

		case ruleAction110:

			p.PushComponent(begin, end, Timestamp)

		case ruleAction111:

			p.PushComponent(begin, end, Array)

		case ruleAction112:

			p.PushComponent(begin, end, Map)

		case ruleAction113:

			p.PushComponent(begin, end, Or)

		case ruleAction114:

			p.PushComponent(begin, end, And)

		case ruleAction115:

			p.PushComponent(begin, end, Not)

		case ruleAction116:

			p.PushComponent(begin, end, Equal)

		case ruleAction117:

			p.PushComponent(begin, end, Less)

		case ruleAction118:

			p.PushComponent(begin, end, LessOrEqual)

		case ruleAction119:

			p.PushComponent(begin, end, Greater)

		case ruleAction120:

			p.PushComponent(begin, end, GreaterOrEqual)

		case ruleAction121:

			p.PushComponent(begin, end, NotEqual)

		case ruleAction122:

			p.PushComponent(begin, end, Concat)

		case ruleAction123:

			p.PushComponent(begin, end, Is)

		case ruleAction124:

			p.PushComponent(begin, end, IsNot)

		case ruleAction125:

			p.PushComponent(begin, end, Plus)

		case ruleAction126:

			p.PushComponent(begin, end, Minus)

		case ruleAction127:

			p.PushComponent(begin, end, Multiply)

		case ruleAction128:

			p.PushComponent(begin, end, Divide)

		case ruleAction129:

			p.PushComponent(begin, end, Modulo)

		case ruleAction130:

			p.PushComponent(begin, end, UnaryMinus)

		case ruleAction131:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		case ruleAction132:

			substr := string([]rune(buffer)[begin:end])
			p.PushComponent(begin, end, Identifier(substr))

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *bqlPegBackend) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens32{tree: make([]token32, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 SingleStatement <- <(spOpt (StatementWithRest / StatementWithoutRest) !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rulespOpt]() {
					goto l0
				}
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[ruleStatementWithRest]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleStatementWithoutRest]() {
						goto l0
					}
				}
			l2:
				{
					position4, tokenIndex4, depth4 := position, tokenIndex, depth
					if !matchDot() {
						goto l4
					}
					goto l0
				l4:
					position, tokenIndex, depth = position4, tokenIndex4, depth4
				}
				depth--
				add(ruleSingleStatement, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 StatementWithRest <- <(<(Statement spOpt ';' spOpt)> .* Action0)> */
		func() bool {
			position5, tokenIndex5, depth5 := position, tokenIndex, depth
			{
				position6 := position
				depth++
				{
					position7 := position
					depth++
					if !_rules[ruleStatement]() {
						goto l5
					}
					if !_rules[rulespOpt]() {
						goto l5
					}
					if buffer[position] != rune(';') {
						goto l5
					}
					position++
					if !_rules[rulespOpt]() {
						goto l5
					}
					depth--
					add(rulePegText, position7)
				}
			l8:
				{
					position9, tokenIndex9, depth9 := position, tokenIndex, depth
					if !matchDot() {
						goto l9
					}
					goto l8
				l9:
					position, tokenIndex, depth = position9, tokenIndex9, depth9
				}
				if !_rules[ruleAction0]() {
					goto l5
				}
				depth--
				add(ruleStatementWithRest, position6)
			}
			return true
		l5:
			position, tokenIndex, depth = position5, tokenIndex5, depth5
			return false
		},
		/* 2 StatementWithoutRest <- <(<(Statement spOpt)> Action1)> */
		func() bool {
			position10, tokenIndex10, depth10 := position, tokenIndex, depth
			{
				position11 := position
				depth++
				{
					position12 := position
					depth++
					if !_rules[ruleStatement]() {
						goto l10
					}
					if !_rules[rulespOpt]() {
						goto l10
					}
					depth--
					add(rulePegText, position12)
				}
				if !_rules[ruleAction1]() {
					goto l10
				}
				depth--
				add(ruleStatementWithoutRest, position11)
			}
			return true
		l10:
			position, tokenIndex, depth = position10, tokenIndex10, depth10
			return false
		},
		/* 3 Statement <- <(SelectUnionStmt / SelectStmt / SourceStmt / SinkStmt / StateStmt / StreamStmt / EvalStmt)> */
		func() bool {
			position13, tokenIndex13, depth13 := position, tokenIndex, depth
			{
				position14 := position
				depth++
				{
					position15, tokenIndex15, depth15 := position, tokenIndex, depth
					if !_rules[ruleSelectUnionStmt]() {
						goto l16
					}
					goto l15
				l16:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleSelectStmt]() {
						goto l17
					}
					goto l15
				l17:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleSourceStmt]() {
						goto l18
					}
					goto l15
				l18:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleSinkStmt]() {
						goto l19
					}
					goto l15
				l19:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleStateStmt]() {
						goto l20
					}
					goto l15
				l20:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleStreamStmt]() {
						goto l21
					}
					goto l15
				l21:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
					if !_rules[ruleEvalStmt]() {
						goto l13
					}
				}
			l15:
				depth--
				add(ruleStatement, position14)
			}
			return true
		l13:
			position, tokenIndex, depth = position13, tokenIndex13, depth13
			return false
		},
		/* 4 SourceStmt <- <(CreateSourceStmt / UpdateSourceStmt / DropSourceStmt / PauseSourceStmt / ResumeSourceStmt / RewindSourceStmt)> */
		func() bool {
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if !_rules[ruleCreateSourceStmt]() {
						goto l25
					}
					goto l24
				l25:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleUpdateSourceStmt]() {
						goto l26
					}
					goto l24
				l26:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleDropSourceStmt]() {
						goto l27
					}
					goto l24
				l27:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[rulePauseSourceStmt]() {
						goto l28
					}
					goto l24
				l28:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleResumeSourceStmt]() {
						goto l29
					}
					goto l24
				l29:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
					if !_rules[ruleRewindSourceStmt]() {
						goto l22
					}
				}
			l24:
				depth--
				add(ruleSourceStmt, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 5 SinkStmt <- <(CreateSinkStmt / UpdateSinkStmt / DropSinkStmt)> */
		func() bool {
			position30, tokenIndex30, depth30 := position, tokenIndex, depth
			{
				position31 := position
				depth++
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					if !_rules[ruleCreateSinkStmt]() {
						goto l33
					}
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if !_rules[ruleUpdateSinkStmt]() {
						goto l34
					}
					goto l32
				l34:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if !_rules[ruleDropSinkStmt]() {
						goto l30
					}
				}
			l32:
				depth--
				add(ruleSinkStmt, position31)
			}
			return true
		l30:
			position, tokenIndex, depth = position30, tokenIndex30, depth30
			return false
		},
		/* 6 StateStmt <- <(CreateStateStmt / UpdateStateStmt / DropStateStmt / LoadStateOrCreateStmt / LoadStateStmt / SaveStateStmt)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				{
					position37, tokenIndex37, depth37 := position, tokenIndex, depth
					if !_rules[ruleCreateStateStmt]() {
						goto l38
					}
					goto l37
				l38:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleUpdateStateStmt]() {
						goto l39
					}
					goto l37
				l39:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleDropStateStmt]() {
						goto l40
					}
					goto l37
				l40:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleLoadStateOrCreateStmt]() {
						goto l41
					}
					goto l37
				l41:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleLoadStateStmt]() {
						goto l42
					}
					goto l37
				l42:
					position, tokenIndex, depth = position37, tokenIndex37, depth37
					if !_rules[ruleSaveStateStmt]() {
						goto l35
					}
				}
			l37:
				depth--
				add(ruleStateStmt, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 7 StreamStmt <- <(CreateStreamAsSelectUnionStmt / CreateStreamAsSelectStmt / DropStreamStmt / InsertIntoFromStmt)> */
		func() bool {
			position43, tokenIndex43, depth43 := position, tokenIndex, depth
			{
				position44 := position
				depth++
				{
					position45, tokenIndex45, depth45 := position, tokenIndex, depth
					if !_rules[ruleCreateStreamAsSelectUnionStmt]() {
						goto l46
					}
					goto l45
				l46:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleCreateStreamAsSelectStmt]() {
						goto l47
					}
					goto l45
				l47:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleDropStreamStmt]() {
						goto l48
					}
					goto l45
				l48:
					position, tokenIndex, depth = position45, tokenIndex45, depth45
					if !_rules[ruleInsertIntoFromStmt]() {
						goto l43
					}
				}
			l45:
				depth--
				add(ruleStreamStmt, position44)
			}
			return true
		l43:
			position, tokenIndex, depth = position43, tokenIndex43, depth43
			return false
		},
		/* 8 SelectStmt <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') Emitter Projections WindowedFrom Filter Grouping Having Action2)> */
		func() bool {
			position49, tokenIndex49, depth49 := position, tokenIndex, depth
			{
				position50 := position
				depth++
				{
					position51, tokenIndex51, depth51 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l52
					}
					position++
					goto l51
				l52:
					position, tokenIndex, depth = position51, tokenIndex51, depth51
					if buffer[position] != rune('S') {
						goto l49
					}
					position++
				}
			l51:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l54
					}
					position++
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					if buffer[position] != rune('E') {
						goto l49
					}
					position++
				}
			l53:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l56
					}
					position++
					goto l55
				l56:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
					if buffer[position] != rune('L') {
						goto l49
					}
					position++
				}
			l55:
				{
					position57, tokenIndex57, depth57 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex, depth = position57, tokenIndex57, depth57
					if buffer[position] != rune('E') {
						goto l49
					}
					position++
				}
			l57:
				{
					position59, tokenIndex59, depth59 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l60
					}
					position++
					goto l59
				l60:
					position, tokenIndex, depth = position59, tokenIndex59, depth59
					if buffer[position] != rune('C') {
						goto l49
					}
					position++
				}
			l59:
				{
					position61, tokenIndex61, depth61 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l62
					}
					position++
					goto l61
				l62:
					position, tokenIndex, depth = position61, tokenIndex61, depth61
					if buffer[position] != rune('T') {
						goto l49
					}
					position++
				}
			l61:
				if !_rules[ruleEmitter]() {
					goto l49
				}
				if !_rules[ruleProjections]() {
					goto l49
				}
				if !_rules[ruleWindowedFrom]() {
					goto l49
				}
				if !_rules[ruleFilter]() {
					goto l49
				}
				if !_rules[ruleGrouping]() {
					goto l49
				}
				if !_rules[ruleHaving]() {
					goto l49
				}
				if !_rules[ruleAction2]() {
					goto l49
				}
				depth--
				add(ruleSelectStmt, position50)
			}
			return true
		l49:
			position, tokenIndex, depth = position49, tokenIndex49, depth49
			return false
		},
		/* 9 SelectUnionStmt <- <(<(SelectStmt (sp (('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N')) sp (('a' / 'A') ('l' / 'L') ('l' / 'L')) sp SelectStmt)+)> Action3)> */
		func() bool {
			position63, tokenIndex63, depth63 := position, tokenIndex, depth
			{
				position64 := position
				depth++
				{
					position65 := position
					depth++
					if !_rules[ruleSelectStmt]() {
						goto l63
					}
					if !_rules[rulesp]() {
						goto l63
					}
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l69
						}
						position++
						goto l68
					l69:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
						if buffer[position] != rune('U') {
							goto l63
						}
						position++
					}
				l68:
					{
						position70, tokenIndex70, depth70 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l71
						}
						position++
						goto l70
					l71:
						position, tokenIndex, depth = position70, tokenIndex70, depth70
						if buffer[position] != rune('N') {
							goto l63
						}
						position++
					}
				l70:
					{
						position72, tokenIndex72, depth72 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l73
						}
						position++
						goto l72
					l73:
						position, tokenIndex, depth = position72, tokenIndex72, depth72
						if buffer[position] != rune('I') {
							goto l63
						}
						position++
					}
				l72:
					{
						position74, tokenIndex74, depth74 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l75
						}
						position++
						goto l74
					l75:
						position, tokenIndex, depth = position74, tokenIndex74, depth74
						if buffer[position] != rune('O') {
							goto l63
						}
						position++
					}
				l74:
					{
						position76, tokenIndex76, depth76 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l77
						}
						position++
						goto l76
					l77:
						position, tokenIndex, depth = position76, tokenIndex76, depth76
						if buffer[position] != rune('N') {
							goto l63
						}
						position++
					}
				l76:
					if !_rules[rulesp]() {
						goto l63
					}
					{
						position78, tokenIndex78, depth78 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l79
						}
						position++
						goto l78
					l79:
						position, tokenIndex, depth = position78, tokenIndex78, depth78
						if buffer[position] != rune('A') {
							goto l63
						}
						position++
					}
				l78:
					{
						position80, tokenIndex80, depth80 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l81
						}
						position++
						goto l80
					l81:
						position, tokenIndex, depth = position80, tokenIndex80, depth80
						if buffer[position] != rune('L') {
							goto l63
						}
						position++
					}
				l80:
					{
						position82, tokenIndex82, depth82 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l83
						}
						position++
						goto l82
					l83:
						position, tokenIndex, depth = position82, tokenIndex82, depth82
						if buffer[position] != rune('L') {
							goto l63
						}
						position++
					}
				l82:
					if !_rules[rulesp]() {
						goto l63
					}
					if !_rules[ruleSelectStmt]() {
						goto l63
					}
				l66:
					{
						position67, tokenIndex67, depth67 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l67
						}
						{
							position84, tokenIndex84, depth84 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l85
							}
							position++
							goto l84
						l85:
							position, tokenIndex, depth = position84, tokenIndex84, depth84
							if buffer[position] != rune('U') {
								goto l67
							}
							position++
						}
					l84:
						{
							position86, tokenIndex86, depth86 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l87
							}
							position++
							goto l86
						l87:
							position, tokenIndex, depth = position86, tokenIndex86, depth86
							if buffer[position] != rune('N') {
								goto l67
							}
							position++
						}
					l86:
						{
							position88, tokenIndex88, depth88 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l89
							}
							position++
							goto l88
						l89:
							position, tokenIndex, depth = position88, tokenIndex88, depth88
							if buffer[position] != rune('I') {
								goto l67
							}
							position++
						}
					l88:
						{
							position90, tokenIndex90, depth90 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l91
							}
							position++
							goto l90
						l91:
							position, tokenIndex, depth = position90, tokenIndex90, depth90
							if buffer[position] != rune('O') {
								goto l67
							}
							position++
						}
					l90:
						{
							position92, tokenIndex92, depth92 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l93
							}
							position++
							goto l92
						l93:
							position, tokenIndex, depth = position92, tokenIndex92, depth92
							if buffer[position] != rune('N') {
								goto l67
							}
							position++
						}
					l92:
						if !_rules[rulesp]() {
							goto l67
						}
						{
							position94, tokenIndex94, depth94 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l95
							}
							position++
							goto l94
						l95:
							position, tokenIndex, depth = position94, tokenIndex94, depth94
							if buffer[position] != rune('A') {
								goto l67
							}
							position++
						}
					l94:
						{
							position96, tokenIndex96, depth96 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l97
							}
							position++
							goto l96
						l97:
							position, tokenIndex, depth = position96, tokenIndex96, depth96
							if buffer[position] != rune('L') {
								goto l67
							}
							position++
						}
					l96:
						{
							position98, tokenIndex98, depth98 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l99
							}
							position++
							goto l98
						l99:
							position, tokenIndex, depth = position98, tokenIndex98, depth98
							if buffer[position] != rune('L') {
								goto l67
							}
							position++
						}
					l98:
						if !_rules[rulesp]() {
							goto l67
						}
						if !_rules[ruleSelectStmt]() {
							goto l67
						}
						goto l66
					l67:
						position, tokenIndex, depth = position67, tokenIndex67, depth67
					}
					depth--
					add(rulePegText, position65)
				}
				if !_rules[ruleAction3]() {
					goto l63
				}
				depth--
				add(ruleSelectUnionStmt, position64)
			}
			return true
		l63:
			position, tokenIndex, depth = position63, tokenIndex63, depth63
			return false
		},
		/* 10 CreateStreamAsSelectStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectStmt Action4)> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
				{
					position102, tokenIndex102, depth102 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l103
					}
					position++
					goto l102
				l103:
					position, tokenIndex, depth = position102, tokenIndex102, depth102
					if buffer[position] != rune('C') {
						goto l100
					}
					position++
				}
			l102:
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('R') {
						goto l100
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('E') {
						goto l100
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('T') {
						goto l100
					}
					position++
				}
			l110:
				{
					position112, tokenIndex112, depth112 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l113
					}
					position++
					goto l112
				l113:
					position, tokenIndex, depth = position112, tokenIndex112, depth112
					if buffer[position] != rune('E') {
						goto l100
					}
					position++
				}
			l112:
				if !_rules[rulesp]() {
					goto l100
				}
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('S') {
						goto l100
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('T') {
						goto l100
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('R') {
						goto l100
					}
					position++
				}
			l118:
				{
					position120, tokenIndex120, depth120 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l121
					}
					position++
					goto l120
				l121:
					position, tokenIndex, depth = position120, tokenIndex120, depth120
					if buffer[position] != rune('E') {
						goto l100
					}
					position++
				}
			l120:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('M') {
						goto l100
					}
					position++
				}
			l124:
				if !_rules[rulesp]() {
					goto l100
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l100
				}
				if !_rules[rulesp]() {
					goto l100
				}
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('A') {
						goto l100
					}
					position++
				}
			l126:
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l129
					}
					position++
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if buffer[position] != rune('S') {
						goto l100
					}
					position++
				}
			l128:
				if !_rules[rulesp]() {
					goto l100
				}
				if !_rules[ruleSelectStmt]() {
					goto l100
				}
				if !_rules[ruleAction4]() {
					goto l100
				}
				depth--
				add(ruleCreateStreamAsSelectStmt, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 11 CreateStreamAsSelectUnionStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier sp (('a' / 'A') ('s' / 'S')) sp SelectUnionStmt Action5)> */
		func() bool {
			position130, tokenIndex130, depth130 := position, tokenIndex, depth
			{
				position131 := position
				depth++
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('C') {
						goto l130
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('R') {
						goto l130
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l136:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l139
					}
					position++
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l138:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('T') {
						goto l130
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l142:
				if !_rules[rulesp]() {
					goto l130
				}
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('S') {
						goto l130
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('T') {
						goto l130
					}
					position++
				}
			l146:
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l149
					}
					position++
					goto l148
				l149:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
					if buffer[position] != rune('R') {
						goto l130
					}
					position++
				}
			l148:
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('E') {
						goto l130
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('M') {
						goto l130
					}
					position++
				}
			l154:
				if !_rules[rulesp]() {
					goto l130
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l130
				}
				if !_rules[rulesp]() {
					goto l130
				}
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l157
					}
					position++
					goto l156
				l157:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
					if buffer[position] != rune('A') {
						goto l130
					}
					position++
				}
			l156:
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('S') {
						goto l130
					}
					position++
				}
			l158:
				if !_rules[rulesp]() {
					goto l130
				}
				if !_rules[ruleSelectUnionStmt]() {
					goto l130
				}
				if !_rules[ruleAction5]() {
					goto l130
				}
				depth--
				add(ruleCreateStreamAsSelectUnionStmt, position131)
			}
			return true
		l130:
			position, tokenIndex, depth = position130, tokenIndex130, depth130
			return false
		},
		/* 12 CreateSourceStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') PausedOpt sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action6)> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('C') {
						goto l160
					}
					position++
				}
			l162:
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l165
					}
					position++
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if buffer[position] != rune('R') {
						goto l160
					}
					position++
				}
			l164:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('A') {
						goto l160
					}
					position++
				}
			l168:
				{
					position170, tokenIndex170, depth170 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l171
					}
					position++
					goto l170
				l171:
					position, tokenIndex, depth = position170, tokenIndex170, depth170
					if buffer[position] != rune('T') {
						goto l160
					}
					position++
				}
			l170:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l172:
				if !_rules[rulePausedOpt]() {
					goto l160
				}
				if !_rules[rulesp]() {
					goto l160
				}
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('S') {
						goto l160
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('O') {
						goto l160
					}
					position++
				}
			l176:
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l179
					}
					position++
					goto l178
				l179:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
					if buffer[position] != rune('U') {
						goto l160
					}
					position++
				}
			l178:
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('R') {
						goto l160
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('C') {
						goto l160
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l184:
				if !_rules[rulesp]() {
					goto l160
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l160
				}
				if !_rules[rulesp]() {
					goto l160
				}
				{
					position186, tokenIndex186, depth186 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position186, tokenIndex186, depth186
					if buffer[position] != rune('T') {
						goto l160
					}
					position++
				}
			l186:
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('Y') {
						goto l160
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('P') {
						goto l160
					}
					position++
				}
			l190:
				{
					position192, tokenIndex192, depth192 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l193
					}
					position++
					goto l192
				l193:
					position, tokenIndex, depth = position192, tokenIndex192, depth192
					if buffer[position] != rune('E') {
						goto l160
					}
					position++
				}
			l192:
				if !_rules[rulesp]() {
					goto l160
				}
				if !_rules[ruleSourceSinkType]() {
					goto l160
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l160
				}
				if !_rules[ruleAction6]() {
					goto l160
				}
				depth--
				add(ruleCreateSourceStmt, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 13 CreateSinkStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action7)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('C') {
						goto l194
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('R') {
						goto l194
					}
					position++
				}
			l198:
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l201
					}
					position++
					goto l200
				l201:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
					if buffer[position] != rune('E') {
						goto l194
					}
					position++
				}
			l200:
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('A') {
						goto l194
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('E') {
						goto l194
					}
					position++
				}
			l206:
				if !_rules[rulesp]() {
					goto l194
				}
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l209
					}
					position++
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if buffer[position] != rune('S') {
						goto l194
					}
					position++
				}
			l208:
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('I') {
						goto l194
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('N') {
						goto l194
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('K') {
						goto l194
					}
					position++
				}
			l214:
				if !_rules[rulesp]() {
					goto l194
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l194
				}
				if !_rules[rulesp]() {
					goto l194
				}
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('T') {
						goto l194
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('Y') {
						goto l194
					}
					position++
				}
			l218:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l221
					}
					position++
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if buffer[position] != rune('P') {
						goto l194
					}
					position++
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('E') {
						goto l194
					}
					position++
				}
			l222:
				if !_rules[rulesp]() {
					goto l194
				}
				if !_rules[ruleSourceSinkType]() {
					goto l194
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l194
				}
				if !_rules[ruleAction7]() {
					goto l194
				}
				depth--
				add(ruleCreateSinkStmt, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 14 CreateStateStmt <- <(('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType SourceSinkSpecs Action8)> */
		func() bool {
			position224, tokenIndex224, depth224 := position, tokenIndex, depth
			{
				position225 := position
				depth++
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('C') {
						goto l224
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('R') {
						goto l224
					}
					position++
				}
			l228:
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l231
					}
					position++
					goto l230
				l231:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l230:
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l233
					}
					position++
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					if buffer[position] != rune('A') {
						goto l224
					}
					position++
				}
			l232:
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l235
					}
					position++
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l237
					}
					position++
					goto l236
				l237:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l236:
				if !_rules[rulesp]() {
					goto l224
				}
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l239
					}
					position++
					goto l238
				l239:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
					if buffer[position] != rune('S') {
						goto l224
					}
					position++
				}
			l238:
				{
					position240, tokenIndex240, depth240 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l241
					}
					position++
					goto l240
				l241:
					position, tokenIndex, depth = position240, tokenIndex240, depth240
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l240:
				{
					position242, tokenIndex242, depth242 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l243
					}
					position++
					goto l242
				l243:
					position, tokenIndex, depth = position242, tokenIndex242, depth242
					if buffer[position] != rune('A') {
						goto l224
					}
					position++
				}
			l242:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l245
					}
					position++
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l244:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l247
					}
					position++
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l246:
				if !_rules[rulesp]() {
					goto l224
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l224
				}
				if !_rules[rulesp]() {
					goto l224
				}
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l249
					}
					position++
					goto l248
				l249:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
					if buffer[position] != rune('T') {
						goto l224
					}
					position++
				}
			l248:
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l251
					}
					position++
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if buffer[position] != rune('Y') {
						goto l224
					}
					position++
				}
			l250:
				{
					position252, tokenIndex252, depth252 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l253
					}
					position++
					goto l252
				l253:
					position, tokenIndex, depth = position252, tokenIndex252, depth252
					if buffer[position] != rune('P') {
						goto l224
					}
					position++
				}
			l252:
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l255
					}
					position++
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					if buffer[position] != rune('E') {
						goto l224
					}
					position++
				}
			l254:
				if !_rules[rulesp]() {
					goto l224
				}
				if !_rules[ruleSourceSinkType]() {
					goto l224
				}
				if !_rules[ruleSourceSinkSpecs]() {
					goto l224
				}
				if !_rules[ruleAction8]() {
					goto l224
				}
				depth--
				add(ruleCreateStateStmt, position225)
			}
			return true
		l224:
			position, tokenIndex, depth = position224, tokenIndex224, depth224
			return false
		},
		/* 15 UpdateStateStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action9)> */
		func() bool {
			position256, tokenIndex256, depth256 := position, tokenIndex, depth
			{
				position257 := position
				depth++
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l259
					}
					position++
					goto l258
				l259:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
					if buffer[position] != rune('U') {
						goto l256
					}
					position++
				}
			l258:
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l261
					}
					position++
					goto l260
				l261:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
					if buffer[position] != rune('P') {
						goto l256
					}
					position++
				}
			l260:
				{
					position262, tokenIndex262, depth262 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l263
					}
					position++
					goto l262
				l263:
					position, tokenIndex, depth = position262, tokenIndex262, depth262
					if buffer[position] != rune('D') {
						goto l256
					}
					position++
				}
			l262:
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l265
					}
					position++
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l264:
				{
					position266, tokenIndex266, depth266 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l267
					}
					position++
					goto l266
				l267:
					position, tokenIndex, depth = position266, tokenIndex266, depth266
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l266:
				{
					position268, tokenIndex268, depth268 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l269
					}
					position++
					goto l268
				l269:
					position, tokenIndex, depth = position268, tokenIndex268, depth268
					if buffer[position] != rune('E') {
						goto l256
					}
					position++
				}
			l268:
				if !_rules[rulesp]() {
					goto l256
				}
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l271
					}
					position++
					goto l270
				l271:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
					if buffer[position] != rune('S') {
						goto l256
					}
					position++
				}
			l270:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l273
					}
					position++
					goto l272
				l273:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l272:
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l275
					}
					position++
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('A') {
						goto l256
					}
					position++
				}
			l274:
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l277
					}
					position++
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					if buffer[position] != rune('T') {
						goto l256
					}
					position++
				}
			l276:
				{
					position278, tokenIndex278, depth278 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l279
					}
					position++
					goto l278
				l279:
					position, tokenIndex, depth = position278, tokenIndex278, depth278
					if buffer[position] != rune('E') {
						goto l256
					}
					position++
				}
			l278:
				if !_rules[rulesp]() {
					goto l256
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l256
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l256
				}
				if !_rules[ruleAction9]() {
					goto l256
				}
				depth--
				add(ruleUpdateStateStmt, position257)
			}
			return true
		l256:
			position, tokenIndex, depth = position256, tokenIndex256, depth256
			return false
		},
		/* 16 UpdateSourceStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier UpdateSourceSinkSpecs Action10)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
				}
			l282:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l285
					}
					position++
					goto l284
				l285:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
					if buffer[position] != rune('P') {
						goto l280
					}
					position++
				}
			l284:
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l287
					}
					position++
					goto l286
				l287:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
					if buffer[position] != rune('D') {
						goto l280
					}
					position++
				}
			l286:
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l289
					}
					position++
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					if buffer[position] != rune('A') {
						goto l280
					}
					position++
				}
			l288:
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l291
					}
					position++
					goto l290
				l291:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
					if buffer[position] != rune('T') {
						goto l280
					}
					position++
				}
			l290:
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l293
					}
					position++
					goto l292
				l293:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
					if buffer[position] != rune('E') {
						goto l280
					}
					position++
				}
			l292:
				if !_rules[rulesp]() {
					goto l280
				}
				{
					position294, tokenIndex294, depth294 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l295
					}
					position++
					goto l294
				l295:
					position, tokenIndex, depth = position294, tokenIndex294, depth294
					if buffer[position] != rune('S') {
						goto l280
					}
					position++
				}
			l294:
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l297
					}
					position++
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					if buffer[position] != rune('O') {
						goto l280
					}
					position++
				}
			l296:
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l299
					}
					position++
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
				}
			l298:
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l301
					}
					position++
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if buffer[position] != rune('R') {
						goto l280
					}
					position++
				}
			l300:
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l303
					}
					position++
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if buffer[position] != rune('C') {
						goto l280
					}
					position++
				}
			l302:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l305
					}
					position++
					goto l304
				l305:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
					if buffer[position] != rune('E') {
						goto l280
					}
					position++
				}
			l304:
				if !_rules[rulesp]() {
					goto l280
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l280
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l280
				}
				if !_rules[ruleAction10]() {
					goto l280
				}
				depth--
				add(ruleUpdateSourceStmt, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 17 UpdateSinkStmt <- <(('u' / 'U') ('p' / 'P') ('d' / 'D') ('a' / 'A') ('t' / 'T') ('e' / 'E') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier UpdateSourceSinkSpecs Action11)> */
		func() bool {
			position306, tokenIndex306, depth306 := position, tokenIndex, depth
			{
				position307 := position
				depth++
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l309
					}
					position++
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					if buffer[position] != rune('U') {
						goto l306
					}
					position++
				}
			l308:
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l311
					}
					position++
					goto l310
				l311:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
					if buffer[position] != rune('P') {
						goto l306
					}
					position++
				}
			l310:
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l313
					}
					position++
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					if buffer[position] != rune('D') {
						goto l306
					}
					position++
				}
			l312:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l315
					}
					position++
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					if buffer[position] != rune('A') {
						goto l306
					}
					position++
				}
			l314:
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l317
					}
					position++
					goto l316
				l317:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
					if buffer[position] != rune('T') {
						goto l306
					}
					position++
				}
			l316:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l319
					}
					position++
					goto l318
				l319:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
					if buffer[position] != rune('E') {
						goto l306
					}
					position++
				}
			l318:
				if !_rules[rulesp]() {
					goto l306
				}
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l321
					}
					position++
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					if buffer[position] != rune('S') {
						goto l306
					}
					position++
				}
			l320:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l323
					}
					position++
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if buffer[position] != rune('I') {
						goto l306
					}
					position++
				}
			l322:
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l325
					}
					position++
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if buffer[position] != rune('N') {
						goto l306
					}
					position++
				}
			l324:
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l327
					}
					position++
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					if buffer[position] != rune('K') {
						goto l306
					}
					position++
				}
			l326:
				if !_rules[rulesp]() {
					goto l306
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l306
				}
				if !_rules[ruleUpdateSourceSinkSpecs]() {
					goto l306
				}
				if !_rules[ruleAction11]() {
					goto l306
				}
				depth--
				add(ruleUpdateSinkStmt, position307)
			}
			return true
		l306:
			position, tokenIndex, depth = position306, tokenIndex306, depth306
			return false
		},
		/* 18 InsertIntoFromStmt <- <(('i' / 'I') ('n' / 'N') ('s' / 'S') ('e' / 'E') ('r' / 'R') ('t' / 'T') sp (('i' / 'I') ('n' / 'N') ('t' / 'T') ('o' / 'O')) sp StreamIdentifier sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp StreamIdentifier Action12)> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l331
					}
					position++
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					if buffer[position] != rune('I') {
						goto l328
					}
					position++
				}
			l330:
				{
					position332, tokenIndex332, depth332 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l333
					}
					position++
					goto l332
				l333:
					position, tokenIndex, depth = position332, tokenIndex332, depth332
					if buffer[position] != rune('N') {
						goto l328
					}
					position++
				}
			l332:
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l335
					}
					position++
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if buffer[position] != rune('S') {
						goto l328
					}
					position++
				}
			l334:
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l337
					}
					position++
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if buffer[position] != rune('E') {
						goto l328
					}
					position++
				}
			l336:
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l339
					}
					position++
					goto l338
				l339:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
					if buffer[position] != rune('R') {
						goto l328
					}
					position++
				}
			l338:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l341
					}
					position++
					goto l340
				l341:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
					if buffer[position] != rune('T') {
						goto l328
					}
					position++
				}
			l340:
				if !_rules[rulesp]() {
					goto l328
				}
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l343
					}
					position++
					goto l342
				l343:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
					if buffer[position] != rune('I') {
						goto l328
					}
					position++
				}
			l342:
				{
					position344, tokenIndex344, depth344 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l345
					}
					position++
					goto l344
				l345:
					position, tokenIndex, depth = position344, tokenIndex344, depth344
					if buffer[position] != rune('N') {
						goto l328
					}
					position++
				}
			l344:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l347
					}
					position++
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					if buffer[position] != rune('T') {
						goto l328
					}
					position++
				}
			l346:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l349
					}
					position++
					goto l348
				l349:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
					if buffer[position] != rune('O') {
						goto l328
					}
					position++
				}
			l348:
				if !_rules[rulesp]() {
					goto l328
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l328
				}
				if !_rules[rulesp]() {
					goto l328
				}
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l351
					}
					position++
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if buffer[position] != rune('F') {
						goto l328
					}
					position++
				}
			l350:
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l353
					}
					position++
					goto l352
				l353:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
					if buffer[position] != rune('R') {
						goto l328
					}
					position++
				}
			l352:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l355
					}
					position++
					goto l354
				l355:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
					if buffer[position] != rune('O') {
						goto l328
					}
					position++
				}
			l354:
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l357
					}
					position++
					goto l356
				l357:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
					if buffer[position] != rune('M') {
						goto l328
					}
					position++
				}
			l356:
				if !_rules[rulesp]() {
					goto l328
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l328
				}
				if !_rules[ruleAction12]() {
					goto l328
				}
				depth--
				add(ruleInsertIntoFromStmt, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 19 PauseSourceStmt <- <(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action13)> */
		func() bool {
			position358, tokenIndex358, depth358 := position, tokenIndex, depth
			{
				position359 := position
				depth++
				{
					position360, tokenIndex360, depth360 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l361
					}
					position++
					goto l360
				l361:
					position, tokenIndex, depth = position360, tokenIndex360, depth360
					if buffer[position] != rune('P') {
						goto l358
					}
					position++
				}
			l360:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l363
					}
					position++
					goto l362
				l363:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if buffer[position] != rune('A') {
						goto l358
					}
					position++
				}
			l362:
				{
					position364, tokenIndex364, depth364 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l365
					}
					position++
					goto l364
				l365:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if buffer[position] != rune('U') {
						goto l358
					}
					position++
				}
			l364:
				{
					position366, tokenIndex366, depth366 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l367
					}
					position++
					goto l366
				l367:
					position, tokenIndex, depth = position366, tokenIndex366, depth366
					if buffer[position] != rune('S') {
						goto l358
					}
					position++
				}
			l366:
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l369
					}
					position++
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					if buffer[position] != rune('E') {
						goto l358
					}
					position++
				}
			l368:
				if !_rules[rulesp]() {
					goto l358
				}
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l371
					}
					position++
					goto l370
				l371:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
					if buffer[position] != rune('S') {
						goto l358
					}
					position++
				}
			l370:
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l373
					}
					position++
					goto l372
				l373:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
					if buffer[position] != rune('O') {
						goto l358
					}
					position++
				}
			l372:
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l375
					}
					position++
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					if buffer[position] != rune('U') {
						goto l358
					}
					position++
				}
			l374:
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l377
					}
					position++
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					if buffer[position] != rune('R') {
						goto l358
					}
					position++
				}
			l376:
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l379
					}
					position++
					goto l378
				l379:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if buffer[position] != rune('C') {
						goto l358
					}
					position++
				}
			l378:
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l381
					}
					position++
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if buffer[position] != rune('E') {
						goto l358
					}
					position++
				}
			l380:
				if !_rules[rulesp]() {
					goto l358
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l358
				}
				if !_rules[ruleAction13]() {
					goto l358
				}
				depth--
				add(rulePauseSourceStmt, position359)
			}
			return true
		l358:
			position, tokenIndex, depth = position358, tokenIndex358, depth358
			return false
		},
		/* 20 ResumeSourceStmt <- <(('r' / 'R') ('e' / 'E') ('s' / 'S') ('u' / 'U') ('m' / 'M') ('e' / 'E') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action14)> */
		func() bool {
			position382, tokenIndex382, depth382 := position, tokenIndex, depth
			{
				position383 := position
				depth++
				{
					position384, tokenIndex384, depth384 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l385
					}
					position++
					goto l384
				l385:
					position, tokenIndex, depth = position384, tokenIndex384, depth384
					if buffer[position] != rune('R') {
						goto l382
					}
					position++
				}
			l384:
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l387
					}
					position++
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					if buffer[position] != rune('E') {
						goto l382
					}
					position++
				}
			l386:
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l389
					}
					position++
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					if buffer[position] != rune('S') {
						goto l382
					}
					position++
				}
			l388:
				{
					position390, tokenIndex390, depth390 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l391
					}
					position++
					goto l390
				l391:
					position, tokenIndex, depth = position390, tokenIndex390, depth390
					if buffer[position] != rune('U') {
						goto l382
					}
					position++
				}
			l390:
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l393
					}
					position++
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if buffer[position] != rune('M') {
						goto l382
					}
					position++
				}
			l392:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l395
					}
					position++
					goto l394
				l395:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
					if buffer[position] != rune('E') {
						goto l382
					}
					position++
				}
			l394:
				if !_rules[rulesp]() {
					goto l382
				}
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('S') {
						goto l382
					}
					position++
				}
			l396:
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l399
					}
					position++
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if buffer[position] != rune('O') {
						goto l382
					}
					position++
				}
			l398:
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					if buffer[position] != rune('U') {
						goto l382
					}
					position++
				}
			l400:
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l403
					}
					position++
					goto l402
				l403:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
					if buffer[position] != rune('R') {
						goto l382
					}
					position++
				}
			l402:
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l405
					}
					position++
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					if buffer[position] != rune('C') {
						goto l382
					}
					position++
				}
			l404:
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l407
					}
					position++
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if buffer[position] != rune('E') {
						goto l382
					}
					position++
				}
			l406:
				if !_rules[rulesp]() {
					goto l382
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l382
				}
				if !_rules[ruleAction14]() {
					goto l382
				}
				depth--
				add(ruleResumeSourceStmt, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 21 RewindSourceStmt <- <(('r' / 'R') ('e' / 'E') ('w' / 'W') ('i' / 'I') ('n' / 'N') ('d' / 'D') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action15)> */
		func() bool {
			position408, tokenIndex408, depth408 := position, tokenIndex, depth
			{
				position409 := position
				depth++
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l411
					}
					position++
					goto l410
				l411:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
					if buffer[position] != rune('R') {
						goto l408
					}
					position++
				}
			l410:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l413
					}
					position++
					goto l412
				l413:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
					if buffer[position] != rune('E') {
						goto l408
					}
					position++
				}
			l412:
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('W') {
						goto l408
					}
					position++
				}
			l414:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l417
					}
					position++
					goto l416
				l417:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
					if buffer[position] != rune('I') {
						goto l408
					}
					position++
				}
			l416:
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l419
					}
					position++
					goto l418
				l419:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
					if buffer[position] != rune('N') {
						goto l408
					}
					position++
				}
			l418:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l421
					}
					position++
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('D') {
						goto l408
					}
					position++
				}
			l420:
				if !_rules[rulesp]() {
					goto l408
				}
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l423
					}
					position++
					goto l422
				l423:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
					if buffer[position] != rune('S') {
						goto l408
					}
					position++
				}
			l422:
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l425
					}
					position++
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if buffer[position] != rune('O') {
						goto l408
					}
					position++
				}
			l424:
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l427
					}
					position++
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if buffer[position] != rune('U') {
						goto l408
					}
					position++
				}
			l426:
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l429
					}
					position++
					goto l428
				l429:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
					if buffer[position] != rune('R') {
						goto l408
					}
					position++
				}
			l428:
				{
					position430, tokenIndex430, depth430 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l431
					}
					position++
					goto l430
				l431:
					position, tokenIndex, depth = position430, tokenIndex430, depth430
					if buffer[position] != rune('C') {
						goto l408
					}
					position++
				}
			l430:
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('E') {
						goto l408
					}
					position++
				}
			l432:
				if !_rules[rulesp]() {
					goto l408
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l408
				}
				if !_rules[ruleAction15]() {
					goto l408
				}
				depth--
				add(ruleRewindSourceStmt, position409)
			}
			return true
		l408:
			position, tokenIndex, depth = position408, tokenIndex408, depth408
			return false
		},
		/* 22 DropSourceStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('c' / 'C') ('e' / 'E')) sp StreamIdentifier Action16)> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l437
					}
					position++
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('D') {
						goto l434
					}
					position++
				}
			l436:
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('R') {
						goto l434
					}
					position++
				}
			l438:
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if buffer[position] != rune('O') {
						goto l434
					}
					position++
				}
			l440:
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
					if buffer[position] != rune('P') {
						goto l434
					}
					position++
				}
			l442:
				if !_rules[rulesp]() {
					goto l434
				}
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l445
					}
					position++
					goto l444
				l445:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
					if buffer[position] != rune('S') {
						goto l434
					}
					position++
				}
			l444:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('O') {
						goto l434
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('U') {
						goto l434
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('R') {
						goto l434
					}
					position++
				}
			l450:
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l453
					}
					position++
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if buffer[position] != rune('C') {
						goto l434
					}
					position++
				}
			l452:
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l455
					}
					position++
					goto l454
				l455:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
					if buffer[position] != rune('E') {
						goto l434
					}
					position++
				}
			l454:
				if !_rules[rulesp]() {
					goto l434
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l434
				}
				if !_rules[ruleAction16]() {
					goto l434
				}
				depth--
				add(ruleDropSourceStmt, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		/* 23 DropStreamStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M')) sp StreamIdentifier Action17)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l459
					}
					position++
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if buffer[position] != rune('D') {
						goto l456
					}
					position++
				}
			l458:
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l461
					}
					position++
					goto l460
				l461:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
					if buffer[position] != rune('R') {
						goto l456
					}
					position++
				}
			l460:
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if buffer[position] != rune('O') {
						goto l456
					}
					position++
				}
			l462:
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l465
					}
					position++
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if buffer[position] != rune('P') {
						goto l456
					}
					position++
				}
			l464:
				if !_rules[rulesp]() {
					goto l456
				}
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l467
					}
					position++
					goto l466
				l467:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
					if buffer[position] != rune('S') {
						goto l456
					}
					position++
				}
			l466:
				{
					position468, tokenIndex468, depth468 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l469
					}
					position++
					goto l468
				l469:
					position, tokenIndex, depth = position468, tokenIndex468, depth468
					if buffer[position] != rune('T') {
						goto l456
					}
					position++
				}
			l468:
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l471
					}
					position++
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if buffer[position] != rune('R') {
						goto l456
					}
					position++
				}
			l470:
				{
					position472, tokenIndex472, depth472 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l473
					}
					position++
					goto l472
				l473:
					position, tokenIndex, depth = position472, tokenIndex472, depth472
					if buffer[position] != rune('E') {
						goto l456
					}
					position++
				}
			l472:
				{
					position474, tokenIndex474, depth474 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l475
					}
					position++
					goto l474
				l475:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if buffer[position] != rune('A') {
						goto l456
					}
					position++
				}
			l474:
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l477
					}
					position++
					goto l476
				l477:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if buffer[position] != rune('M') {
						goto l456
					}
					position++
				}
			l476:
				if !_rules[rulesp]() {
					goto l456
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l456
				}
				if !_rules[ruleAction17]() {
					goto l456
				}
				depth--
				add(ruleDropStreamStmt, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 24 DropSinkStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('i' / 'I') ('n' / 'N') ('k' / 'K')) sp StreamIdentifier Action18)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				{
					position480, tokenIndex480, depth480 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l481
					}
					position++
					goto l480
				l481:
					position, tokenIndex, depth = position480, tokenIndex480, depth480
					if buffer[position] != rune('D') {
						goto l478
					}
					position++
				}
			l480:
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l483
					}
					position++
					goto l482
				l483:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
					if buffer[position] != rune('R') {
						goto l478
					}
					position++
				}
			l482:
				{
					position484, tokenIndex484, depth484 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l485
					}
					position++
					goto l484
				l485:
					position, tokenIndex, depth = position484, tokenIndex484, depth484
					if buffer[position] != rune('O') {
						goto l478
					}
					position++
				}
			l484:
				{
					position486, tokenIndex486, depth486 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l487
					}
					position++
					goto l486
				l487:
					position, tokenIndex, depth = position486, tokenIndex486, depth486
					if buffer[position] != rune('P') {
						goto l478
					}
					position++
				}
			l486:
				if !_rules[rulesp]() {
					goto l478
				}
				{
					position488, tokenIndex488, depth488 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l489
					}
					position++
					goto l488
				l489:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
					if buffer[position] != rune('S') {
						goto l478
					}
					position++
				}
			l488:
				{
					position490, tokenIndex490, depth490 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l491
					}
					position++
					goto l490
				l491:
					position, tokenIndex, depth = position490, tokenIndex490, depth490
					if buffer[position] != rune('I') {
						goto l478
					}
					position++
				}
			l490:
				{
					position492, tokenIndex492, depth492 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l493
					}
					position++
					goto l492
				l493:
					position, tokenIndex, depth = position492, tokenIndex492, depth492
					if buffer[position] != rune('N') {
						goto l478
					}
					position++
				}
			l492:
				{
					position494, tokenIndex494, depth494 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l495
					}
					position++
					goto l494
				l495:
					position, tokenIndex, depth = position494, tokenIndex494, depth494
					if buffer[position] != rune('K') {
						goto l478
					}
					position++
				}
			l494:
				if !_rules[rulesp]() {
					goto l478
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l478
				}
				if !_rules[ruleAction18]() {
					goto l478
				}
				depth--
				add(ruleDropSinkStmt, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 25 DropStateStmt <- <(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier Action19)> */
		func() bool {
			position496, tokenIndex496, depth496 := position, tokenIndex, depth
			{
				position497 := position
				depth++
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l499
					}
					position++
					goto l498
				l499:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
					if buffer[position] != rune('D') {
						goto l496
					}
					position++
				}
			l498:
				{
					position500, tokenIndex500, depth500 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l501
					}
					position++
					goto l500
				l501:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if buffer[position] != rune('R') {
						goto l496
					}
					position++
				}
			l500:
				{
					position502, tokenIndex502, depth502 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l503
					}
					position++
					goto l502
				l503:
					position, tokenIndex, depth = position502, tokenIndex502, depth502
					if buffer[position] != rune('O') {
						goto l496
					}
					position++
				}
			l502:
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l505
					}
					position++
					goto l504
				l505:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
					if buffer[position] != rune('P') {
						goto l496
					}
					position++
				}
			l504:
				if !_rules[rulesp]() {
					goto l496
				}
				{
					position506, tokenIndex506, depth506 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l507
					}
					position++
					goto l506
				l507:
					position, tokenIndex, depth = position506, tokenIndex506, depth506
					if buffer[position] != rune('S') {
						goto l496
					}
					position++
				}
			l506:
				{
					position508, tokenIndex508, depth508 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l509
					}
					position++
					goto l508
				l509:
					position, tokenIndex, depth = position508, tokenIndex508, depth508
					if buffer[position] != rune('T') {
						goto l496
					}
					position++
				}
			l508:
				{
					position510, tokenIndex510, depth510 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l511
					}
					position++
					goto l510
				l511:
					position, tokenIndex, depth = position510, tokenIndex510, depth510
					if buffer[position] != rune('A') {
						goto l496
					}
					position++
				}
			l510:
				{
					position512, tokenIndex512, depth512 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l513
					}
					position++
					goto l512
				l513:
					position, tokenIndex, depth = position512, tokenIndex512, depth512
					if buffer[position] != rune('T') {
						goto l496
					}
					position++
				}
			l512:
				{
					position514, tokenIndex514, depth514 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l515
					}
					position++
					goto l514
				l515:
					position, tokenIndex, depth = position514, tokenIndex514, depth514
					if buffer[position] != rune('E') {
						goto l496
					}
					position++
				}
			l514:
				if !_rules[rulesp]() {
					goto l496
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l496
				}
				if !_rules[ruleAction19]() {
					goto l496
				}
				depth--
				add(ruleDropStateStmt, position497)
			}
			return true
		l496:
			position, tokenIndex, depth = position496, tokenIndex496, depth496
			return false
		},
		/* 26 LoadStateStmt <- <(('l' / 'L') ('o' / 'O') ('a' / 'A') ('d' / 'D') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier sp (('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E')) sp SourceSinkType StateTagOpt SetOptSpecs Action20)> */
		func() bool {
			position516, tokenIndex516, depth516 := position, tokenIndex, depth
			{
				position517 := position
				depth++
				{
					position518, tokenIndex518, depth518 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l519
					}
					position++
					goto l518
				l519:
					position, tokenIndex, depth = position518, tokenIndex518, depth518
					if buffer[position] != rune('L') {
						goto l516
					}
					position++
				}
			l518:
				{
					position520, tokenIndex520, depth520 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l521
					}
					position++
					goto l520
				l521:
					position, tokenIndex, depth = position520, tokenIndex520, depth520
					if buffer[position] != rune('O') {
						goto l516
					}
					position++
				}
			l520:
				{
					position522, tokenIndex522, depth522 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l523
					}
					position++
					goto l522
				l523:
					position, tokenIndex, depth = position522, tokenIndex522, depth522
					if buffer[position] != rune('A') {
						goto l516
					}
					position++
				}
			l522:
				{
					position524, tokenIndex524, depth524 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l525
					}
					position++
					goto l524
				l525:
					position, tokenIndex, depth = position524, tokenIndex524, depth524
					if buffer[position] != rune('D') {
						goto l516
					}
					position++
				}
			l524:
				if !_rules[rulesp]() {
					goto l516
				}
				{
					position526, tokenIndex526, depth526 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l527
					}
					position++
					goto l526
				l527:
					position, tokenIndex, depth = position526, tokenIndex526, depth526
					if buffer[position] != rune('S') {
						goto l516
					}
					position++
				}
			l526:
				{
					position528, tokenIndex528, depth528 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l529
					}
					position++
					goto l528
				l529:
					position, tokenIndex, depth = position528, tokenIndex528, depth528
					if buffer[position] != rune('T') {
						goto l516
					}
					position++
				}
			l528:
				{
					position530, tokenIndex530, depth530 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l531
					}
					position++
					goto l530
				l531:
					position, tokenIndex, depth = position530, tokenIndex530, depth530
					if buffer[position] != rune('A') {
						goto l516
					}
					position++
				}
			l530:
				{
					position532, tokenIndex532, depth532 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l533
					}
					position++
					goto l532
				l533:
					position, tokenIndex, depth = position532, tokenIndex532, depth532
					if buffer[position] != rune('T') {
						goto l516
					}
					position++
				}
			l532:
				{
					position534, tokenIndex534, depth534 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l535
					}
					position++
					goto l534
				l535:
					position, tokenIndex, depth = position534, tokenIndex534, depth534
					if buffer[position] != rune('E') {
						goto l516
					}
					position++
				}
			l534:
				if !_rules[rulesp]() {
					goto l516
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l516
				}
				if !_rules[rulesp]() {
					goto l516
				}
				{
					position536, tokenIndex536, depth536 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l537
					}
					position++
					goto l536
				l537:
					position, tokenIndex, depth = position536, tokenIndex536, depth536
					if buffer[position] != rune('T') {
						goto l516
					}
					position++
				}
			l536:
				{
					position538, tokenIndex538, depth538 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l539
					}
					position++
					goto l538
				l539:
					position, tokenIndex, depth = position538, tokenIndex538, depth538
					if buffer[position] != rune('Y') {
						goto l516
					}
					position++
				}
			l538:
				{
					position540, tokenIndex540, depth540 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l541
					}
					position++
					goto l540
				l541:
					position, tokenIndex, depth = position540, tokenIndex540, depth540
					if buffer[position] != rune('P') {
						goto l516
					}
					position++
				}
			l540:
				{
					position542, tokenIndex542, depth542 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l543
					}
					position++
					goto l542
				l543:
					position, tokenIndex, depth = position542, tokenIndex542, depth542
					if buffer[position] != rune('E') {
						goto l516
					}
					position++
				}
			l542:
				if !_rules[rulesp]() {
					goto l516
				}
				if !_rules[ruleSourceSinkType]() {
					goto l516
				}
				if !_rules[ruleStateTagOpt]() {
					goto l516
				}
				if !_rules[ruleSetOptSpecs]() {
					goto l516
				}
				if !_rules[ruleAction20]() {
					goto l516
				}
				depth--
				add(ruleLoadStateStmt, position517)
			}
			return true
		l516:
			position, tokenIndex, depth = position516, tokenIndex516, depth516
			return false
		},
		/* 27 LoadStateOrCreateStmt <- <(LoadStateStmt sp (('o' / 'O') ('r' / 'R')) sp (('c' / 'C') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp (('i' / 'I') ('f' / 'F')) sp (('n' / 'N') ('o' / 'O') ('t' / 'T')) sp ((('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') ('d' / 'D')) / (('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S'))) SourceSinkSpecs Action21)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				if !_rules[ruleLoadStateStmt]() {
					goto l544
				}
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position546, tokenIndex546, depth546 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l547
					}
					position++
					goto l546
				l547:
					position, tokenIndex, depth = position546, tokenIndex546, depth546
					if buffer[position] != rune('O') {
						goto l544
					}
					position++
				}
			l546:
				{
					position548, tokenIndex548, depth548 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l549
					}
					position++
					goto l548
				l549:
					position, tokenIndex, depth = position548, tokenIndex548, depth548
					if buffer[position] != rune('R') {
						goto l544
					}
					position++
				}
			l548:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position550, tokenIndex550, depth550 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l551
					}
					position++
					goto l550
				l551:
					position, tokenIndex, depth = position550, tokenIndex550, depth550
					if buffer[position] != rune('C') {
						goto l544
					}
					position++
				}
			l550:
				{
					position552, tokenIndex552, depth552 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l553
					}
					position++
					goto l552
				l553:
					position, tokenIndex, depth = position552, tokenIndex552, depth552
					if buffer[position] != rune('R') {
						goto l544
					}
					position++
				}
			l552:
				{
					position554, tokenIndex554, depth554 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l555
					}
					position++
					goto l554
				l555:
					position, tokenIndex, depth = position554, tokenIndex554, depth554
					if buffer[position] != rune('E') {
						goto l544
					}
					position++
				}
			l554:
				{
					position556, tokenIndex556, depth556 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l557
					}
					position++
					goto l556
				l557:
					position, tokenIndex, depth = position556, tokenIndex556, depth556
					if buffer[position] != rune('A') {
						goto l544
					}
					position++
				}
			l556:
				{
					position558, tokenIndex558, depth558 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l559
					}
					position++
					goto l558
				l559:
					position, tokenIndex, depth = position558, tokenIndex558, depth558
					if buffer[position] != rune('T') {
						goto l544
					}
					position++
				}
			l558:
				{
					position560, tokenIndex560, depth560 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l561
					}
					position++
					goto l560
				l561:
					position, tokenIndex, depth = position560, tokenIndex560, depth560
					if buffer[position] != rune('E') {
						goto l544
					}
					position++
				}
			l560:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position562, tokenIndex562, depth562 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l563
					}
					position++
					goto l562
				l563:
					position, tokenIndex, depth = position562, tokenIndex562, depth562
					if buffer[position] != rune('I') {
						goto l544
					}
					position++
				}
			l562:
				{
					position564, tokenIndex564, depth564 := position, tokenIndex, depth
					if buffer[position] != rune('f') {
						goto l565
					}
					position++
					goto l564
				l565:
					position, tokenIndex, depth = position564, tokenIndex564, depth564
					if buffer[position] != rune('F') {
						goto l544
					}
					position++
				}
			l564:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position566, tokenIndex566, depth566 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l567
					}
					position++
					goto l566
				l567:
					position, tokenIndex, depth = position566, tokenIndex566, depth566
					if buffer[position] != rune('N') {
						goto l544
					}
					position++
				}
			l566:
				{
					position568, tokenIndex568, depth568 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l569
					}
					position++
					goto l568
				l569:
					position, tokenIndex, depth = position568, tokenIndex568, depth568
					if buffer[position] != rune('O') {
						goto l544
					}
					position++
				}
			l568:
				{
					position570, tokenIndex570, depth570 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l571
					}
					position++
					goto l570
				l571:
					position, tokenIndex, depth = position570, tokenIndex570, depth570
					if buffer[position] != rune('T') {
						goto l544
					}
					position++
				}
			l570:
				if !_rules[rulesp]() {
					goto l544
				}
				{
					position572, tokenIndex572, depth572 := position, tokenIndex, depth
					{
						position574, tokenIndex574, depth574 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l575
						}
						position++
						goto l574
					l575:
						position, tokenIndex, depth = position574, tokenIndex574, depth574
						if buffer[position] != rune('S') {
							goto l573
						}
						position++
					}
				l574:
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l577
						}
						position++
						goto l576
					l577:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						if buffer[position] != rune('A') {
							goto l573
						}
						position++
					}
				l576:
					{
						position578, tokenIndex578, depth578 := position, tokenIndex, depth
						if buffer[position] != rune('v') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
						if buffer[position] != rune('V') {
							goto l573
						}
						position++
					}
				l578:
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
						if buffer[position] != rune('E') {
							goto l573
						}
						position++
					}
				l580:
					{
						position582, tokenIndex582, depth582 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l583
						}
						position++
						goto l582
					l583:
						position, tokenIndex, depth = position582, tokenIndex582, depth582
						if buffer[position] != rune('D') {
							goto l573
						}
						position++
					}
				l582:
					goto l572
				l573:
					position, tokenIndex, depth = position572, tokenIndex572, depth572
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
						if buffer[position] != rune('E') {
							goto l544
						}
						position++
					}
				l584:
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if buffer[position] != rune('x') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
						if buffer[position] != rune('X') {
							goto l544
						}
						position++
					}
				l586:
					{
						position588, tokenIndex588, depth588 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
						if buffer[position] != rune('I') {
							goto l544
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
						if buffer[position] != rune('S') {
							goto l544
						}
						position++
					}
				l590:
					{
						position592, tokenIndex592, depth592 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position592, tokenIndex592, depth592
						if buffer[position] != rune('T') {
							goto l544
						}
						position++
					}
				l592:
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('S') {
							goto l544
						}
						position++
					}
				l594:
				}
			l572:
				if !_rules[ruleSourceSinkSpecs]() {
					goto l544
				}
				if !_rules[ruleAction21]() {
					goto l544
				}
				depth--
				add(ruleLoadStateOrCreateStmt, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 28 SaveStateStmt <- <(('s' / 'S') ('a' / 'A') ('v' / 'V') ('e' / 'E') sp (('s' / 'S') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('e' / 'E')) sp StreamIdentifier StateTagOpt Action22)> */
		func() bool {
			position596, tokenIndex596, depth596 := position, tokenIndex, depth
			{
				position597 := position
				depth++
				{
					position598, tokenIndex598, depth598 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l599
					}
					position++
					goto l598
				l599:
					position, tokenIndex, depth = position598, tokenIndex598, depth598
					if buffer[position] != rune('S') {
						goto l596
					}
					position++
				}
			l598:
				{
					position600, tokenIndex600, depth600 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l601
					}
					position++
					goto l600
				l601:
					position, tokenIndex, depth = position600, tokenIndex600, depth600
					if buffer[position] != rune('A') {
						goto l596
					}
					position++
				}
			l600:
				{
					position602, tokenIndex602, depth602 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l603
					}
					position++
					goto l602
				l603:
					position, tokenIndex, depth = position602, tokenIndex602, depth602
					if buffer[position] != rune('V') {
						goto l596
					}
					position++
				}
			l602:
				{
					position604, tokenIndex604, depth604 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l605
					}
					position++
					goto l604
				l605:
					position, tokenIndex, depth = position604, tokenIndex604, depth604
					if buffer[position] != rune('E') {
						goto l596
					}
					position++
				}
			l604:
				if !_rules[rulesp]() {
					goto l596
				}
				{
					position606, tokenIndex606, depth606 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l607
					}
					position++
					goto l606
				l607:
					position, tokenIndex, depth = position606, tokenIndex606, depth606
					if buffer[position] != rune('S') {
						goto l596
					}
					position++
				}
			l606:
				{
					position608, tokenIndex608, depth608 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l609
					}
					position++
					goto l608
				l609:
					position, tokenIndex, depth = position608, tokenIndex608, depth608
					if buffer[position] != rune('T') {
						goto l596
					}
					position++
				}
			l608:
				{
					position610, tokenIndex610, depth610 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l611
					}
					position++
					goto l610
				l611:
					position, tokenIndex, depth = position610, tokenIndex610, depth610
					if buffer[position] != rune('A') {
						goto l596
					}
					position++
				}
			l610:
				{
					position612, tokenIndex612, depth612 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l613
					}
					position++
					goto l612
				l613:
					position, tokenIndex, depth = position612, tokenIndex612, depth612
					if buffer[position] != rune('T') {
						goto l596
					}
					position++
				}
			l612:
				{
					position614, tokenIndex614, depth614 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l615
					}
					position++
					goto l614
				l615:
					position, tokenIndex, depth = position614, tokenIndex614, depth614
					if buffer[position] != rune('E') {
						goto l596
					}
					position++
				}
			l614:
				if !_rules[rulesp]() {
					goto l596
				}
				if !_rules[ruleStreamIdentifier]() {
					goto l596
				}
				if !_rules[ruleStateTagOpt]() {
					goto l596
				}
				if !_rules[ruleAction22]() {
					goto l596
				}
				depth--
				add(ruleSaveStateStmt, position597)
			}
			return true
		l596:
			position, tokenIndex, depth = position596, tokenIndex596, depth596
			return false
		},
		/* 29 EvalStmt <- <(('e' / 'E') ('v' / 'V') ('a' / 'A') ('l' / 'L') sp Expression <(sp (('o' / 'O') ('n' / 'N')) sp MapExpr)?> Action23)> */
		func() bool {
			position616, tokenIndex616, depth616 := position, tokenIndex, depth
			{
				position617 := position
				depth++
				{
					position618, tokenIndex618, depth618 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l619
					}
					position++
					goto l618
				l619:
					position, tokenIndex, depth = position618, tokenIndex618, depth618
					if buffer[position] != rune('E') {
						goto l616
					}
					position++
				}
			l618:
				{
					position620, tokenIndex620, depth620 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l621
					}
					position++
					goto l620
				l621:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
					if buffer[position] != rune('V') {
						goto l616
					}
					position++
				}
			l620:
				{
					position622, tokenIndex622, depth622 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l623
					}
					position++
					goto l622
				l623:
					position, tokenIndex, depth = position622, tokenIndex622, depth622
					if buffer[position] != rune('A') {
						goto l616
					}
					position++
				}
			l622:
				{
					position624, tokenIndex624, depth624 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l625
					}
					position++
					goto l624
				l625:
					position, tokenIndex, depth = position624, tokenIndex624, depth624
					if buffer[position] != rune('L') {
						goto l616
					}
					position++
				}
			l624:
				if !_rules[rulesp]() {
					goto l616
				}
				if !_rules[ruleExpression]() {
					goto l616
				}
				{
					position626 := position
					depth++
					{
						position627, tokenIndex627, depth627 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l627
						}
						{
							position629, tokenIndex629, depth629 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l630
							}
							position++
							goto l629
						l630:
							position, tokenIndex, depth = position629, tokenIndex629, depth629
							if buffer[position] != rune('O') {
								goto l627
							}
							position++
						}
					l629:
						{
							position631, tokenIndex631, depth631 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l632
							}
							position++
							goto l631
						l632:
							position, tokenIndex, depth = position631, tokenIndex631, depth631
							if buffer[position] != rune('N') {
								goto l627
							}
							position++
						}
					l631:
						if !_rules[rulesp]() {
							goto l627
						}
						if !_rules[ruleMapExpr]() {
							goto l627
						}
						goto l628
					l627:
						position, tokenIndex, depth = position627, tokenIndex627, depth627
					}
				l628:
					depth--
					add(rulePegText, position626)
				}
				if !_rules[ruleAction23]() {
					goto l616
				}
				depth--
				add(ruleEvalStmt, position617)
			}
			return true
		l616:
			position, tokenIndex, depth = position616, tokenIndex616, depth616
			return false
		},
		/* 30 Emitter <- <(sp (ISTREAM / DSTREAM / RSTREAM) EmitterOptions Action24)> */
		func() bool {
			position633, tokenIndex633, depth633 := position, tokenIndex, depth
			{
				position634 := position
				depth++
				if !_rules[rulesp]() {
					goto l633
				}
				{
					position635, tokenIndex635, depth635 := position, tokenIndex, depth
					if !_rules[ruleISTREAM]() {
						goto l636
					}
					goto l635
				l636:
					position, tokenIndex, depth = position635, tokenIndex635, depth635
					if !_rules[ruleDSTREAM]() {
						goto l637
					}
					goto l635
				l637:
					position, tokenIndex, depth = position635, tokenIndex635, depth635
					if !_rules[ruleRSTREAM]() {
						goto l633
					}
				}
			l635:
				if !_rules[ruleEmitterOptions]() {
					goto l633
				}
				if !_rules[ruleAction24]() {
					goto l633
				}
				depth--
				add(ruleEmitter, position634)
			}
			return true
		l633:
			position, tokenIndex, depth = position633, tokenIndex633, depth633
			return false
		},
		/* 31 EmitterOptions <- <(<(spOpt '[' spOpt EmitterOptionCombinations spOpt ']')?> Action25)> */
		func() bool {
			position638, tokenIndex638, depth638 := position, tokenIndex, depth
			{
				position639 := position
				depth++
				{
					position640 := position
					depth++
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l641
						}
						if buffer[position] != rune('[') {
							goto l641
						}
						position++
						if !_rules[rulespOpt]() {
							goto l641
						}
						if !_rules[ruleEmitterOptionCombinations]() {
							goto l641
						}
						if !_rules[rulespOpt]() {
							goto l641
						}
						if buffer[position] != rune(']') {
							goto l641
						}
						position++
						goto l642
					l641:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
					}
				l642:
					depth--
					add(rulePegText, position640)
				}
				if !_rules[ruleAction25]() {
					goto l638
				}
				depth--
				add(ruleEmitterOptions, position639)
			}
			return true
		l638:
			position, tokenIndex, depth = position638, tokenIndex638, depth638
			return false
		},
		/* 32 EmitterOptionCombinations <- <(EmitterLimit / (EmitterSample sp EmitterLimit) / EmitterSample)> */
		func() bool {
			position643, tokenIndex643, depth643 := position, tokenIndex, depth
			{
				position644 := position
				depth++
				{
					position645, tokenIndex645, depth645 := position, tokenIndex, depth
					if !_rules[ruleEmitterLimit]() {
						goto l646
					}
					goto l645
				l646:
					position, tokenIndex, depth = position645, tokenIndex645, depth645
					if !_rules[ruleEmitterSample]() {
						goto l647
					}
					if !_rules[rulesp]() {
						goto l647
					}
					if !_rules[ruleEmitterLimit]() {
						goto l647
					}
					goto l645
				l647:
					position, tokenIndex, depth = position645, tokenIndex645, depth645
					if !_rules[ruleEmitterSample]() {
						goto l643
					}
				}
			l645:
				depth--
				add(ruleEmitterOptionCombinations, position644)
			}
			return true
		l643:
			position, tokenIndex, depth = position643, tokenIndex643, depth643
			return false
		},
		/* 33 EmitterLimit <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') sp NumericLiteral Action26)> */
		func() bool {
			position648, tokenIndex648, depth648 := position, tokenIndex, depth
			{
				position649 := position
				depth++
				{
					position650, tokenIndex650, depth650 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l651
					}
					position++
					goto l650
				l651:
					position, tokenIndex, depth = position650, tokenIndex650, depth650
					if buffer[position] != rune('L') {
						goto l648
					}
					position++
				}
			l650:
				{
					position652, tokenIndex652, depth652 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l653
					}
					position++
					goto l652
				l653:
					position, tokenIndex, depth = position652, tokenIndex652, depth652
					if buffer[position] != rune('I') {
						goto l648
					}
					position++
				}
			l652:
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l655
					}
					position++
					goto l654
				l655:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
					if buffer[position] != rune('M') {
						goto l648
					}
					position++
				}
			l654:
				{
					position656, tokenIndex656, depth656 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l657
					}
					position++
					goto l656
				l657:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
					if buffer[position] != rune('I') {
						goto l648
					}
					position++
				}
			l656:
				{
					position658, tokenIndex658, depth658 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l659
					}
					position++
					goto l658
				l659:
					position, tokenIndex, depth = position658, tokenIndex658, depth658
					if buffer[position] != rune('T') {
						goto l648
					}
					position++
				}
			l658:
				if !_rules[rulesp]() {
					goto l648
				}
				if !_rules[ruleNumericLiteral]() {
					goto l648
				}
				if !_rules[ruleAction26]() {
					goto l648
				}
				depth--
				add(ruleEmitterLimit, position649)
			}
			return true
		l648:
			position, tokenIndex, depth = position648, tokenIndex648, depth648
			return false
		},
		/* 34 EmitterSample <- <(CountBasedSampling / RandomizedSampling / TimeBasedSampling)> */
		func() bool {
			position660, tokenIndex660, depth660 := position, tokenIndex, depth
			{
				position661 := position
				depth++
				{
					position662, tokenIndex662, depth662 := position, tokenIndex, depth
					if !_rules[ruleCountBasedSampling]() {
						goto l663
					}
					goto l662
				l663:
					position, tokenIndex, depth = position662, tokenIndex662, depth662
					if !_rules[ruleRandomizedSampling]() {
						goto l664
					}
					goto l662
				l664:
					position, tokenIndex, depth = position662, tokenIndex662, depth662
					if !_rules[ruleTimeBasedSampling]() {
						goto l660
					}
				}
			l662:
				depth--
				add(ruleEmitterSample, position661)
			}
			return true
		l660:
			position, tokenIndex, depth = position660, tokenIndex660, depth660
			return false
		},
		/* 35 CountBasedSampling <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp NumericLiteral spOpt '-'? spOpt ((('s' / 'S') ('t' / 'T')) / (('n' / 'N') ('d' / 'D')) / (('r' / 'R') ('d' / 'D')) / (('t' / 'T') ('h' / 'H'))) sp (('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E')) Action27)> */
		func() bool {
			position665, tokenIndex665, depth665 := position, tokenIndex, depth
			{
				position666 := position
				depth++
				{
					position667, tokenIndex667, depth667 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l668
					}
					position++
					goto l667
				l668:
					position, tokenIndex, depth = position667, tokenIndex667, depth667
					if buffer[position] != rune('E') {
						goto l665
					}
					position++
				}
			l667:
				{
					position669, tokenIndex669, depth669 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l670
					}
					position++
					goto l669
				l670:
					position, tokenIndex, depth = position669, tokenIndex669, depth669
					if buffer[position] != rune('V') {
						goto l665
					}
					position++
				}
			l669:
				{
					position671, tokenIndex671, depth671 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l672
					}
					position++
					goto l671
				l672:
					position, tokenIndex, depth = position671, tokenIndex671, depth671
					if buffer[position] != rune('E') {
						goto l665
					}
					position++
				}
			l671:
				{
					position673, tokenIndex673, depth673 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l674
					}
					position++
					goto l673
				l674:
					position, tokenIndex, depth = position673, tokenIndex673, depth673
					if buffer[position] != rune('R') {
						goto l665
					}
					position++
				}
			l673:
				{
					position675, tokenIndex675, depth675 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l676
					}
					position++
					goto l675
				l676:
					position, tokenIndex, depth = position675, tokenIndex675, depth675
					if buffer[position] != rune('Y') {
						goto l665
					}
					position++
				}
			l675:
				if !_rules[rulesp]() {
					goto l665
				}
				if !_rules[ruleNumericLiteral]() {
					goto l665
				}
				if !_rules[rulespOpt]() {
					goto l665
				}
				{
					position677, tokenIndex677, depth677 := position, tokenIndex, depth
					if buffer[position] != rune('-') {
						goto l677
					}
					position++
					goto l678
				l677:
					position, tokenIndex, depth = position677, tokenIndex677, depth677
				}
			l678:
				if !_rules[rulespOpt]() {
					goto l665
				}
				{
					position679, tokenIndex679, depth679 := position, tokenIndex, depth
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l682
						}
						position++
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						if buffer[position] != rune('S') {
							goto l680
						}
						position++
					}
				l681:
					{
						position683, tokenIndex683, depth683 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l684
						}
						position++
						goto l683
					l684:
						position, tokenIndex, depth = position683, tokenIndex683, depth683
						if buffer[position] != rune('T') {
							goto l680
						}
						position++
					}
				l683:
					goto l679
				l680:
					position, tokenIndex, depth = position679, tokenIndex679, depth679
					{
						position686, tokenIndex686, depth686 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l687
						}
						position++
						goto l686
					l687:
						position, tokenIndex, depth = position686, tokenIndex686, depth686
						if buffer[position] != rune('N') {
							goto l685
						}
						position++
					}
				l686:
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l689
						}
						position++
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if buffer[position] != rune('D') {
							goto l685
						}
						position++
					}
				l688:
					goto l679
				l685:
					position, tokenIndex, depth = position679, tokenIndex679, depth679
					{
						position691, tokenIndex691, depth691 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l692
						}
						position++
						goto l691
					l692:
						position, tokenIndex, depth = position691, tokenIndex691, depth691
						if buffer[position] != rune('R') {
							goto l690
						}
						position++
					}
				l691:
					{
						position693, tokenIndex693, depth693 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l694
						}
						position++
						goto l693
					l694:
						position, tokenIndex, depth = position693, tokenIndex693, depth693
						if buffer[position] != rune('D') {
							goto l690
						}
						position++
					}
				l693:
					goto l679
				l690:
					position, tokenIndex, depth = position679, tokenIndex679, depth679
					{
						position695, tokenIndex695, depth695 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l696
						}
						position++
						goto l695
					l696:
						position, tokenIndex, depth = position695, tokenIndex695, depth695
						if buffer[position] != rune('T') {
							goto l665
						}
						position++
					}
				l695:
					{
						position697, tokenIndex697, depth697 := position, tokenIndex, depth
						if buffer[position] != rune('h') {
							goto l698
						}
						position++
						goto l697
					l698:
						position, tokenIndex, depth = position697, tokenIndex697, depth697
						if buffer[position] != rune('H') {
							goto l665
						}
						position++
					}
				l697:
				}
			l679:
				if !_rules[rulesp]() {
					goto l665
				}
				{
					position699, tokenIndex699, depth699 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l700
					}
					position++
					goto l699
				l700:
					position, tokenIndex, depth = position699, tokenIndex699, depth699
					if buffer[position] != rune('T') {
						goto l665
					}
					position++
				}
			l699:
				{
					position701, tokenIndex701, depth701 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l702
					}
					position++
					goto l701
				l702:
					position, tokenIndex, depth = position701, tokenIndex701, depth701
					if buffer[position] != rune('U') {
						goto l665
					}
					position++
				}
			l701:
				{
					position703, tokenIndex703, depth703 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l704
					}
					position++
					goto l703
				l704:
					position, tokenIndex, depth = position703, tokenIndex703, depth703
					if buffer[position] != rune('P') {
						goto l665
					}
					position++
				}
			l703:
				{
					position705, tokenIndex705, depth705 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l706
					}
					position++
					goto l705
				l706:
					position, tokenIndex, depth = position705, tokenIndex705, depth705
					if buffer[position] != rune('L') {
						goto l665
					}
					position++
				}
			l705:
				{
					position707, tokenIndex707, depth707 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l708
					}
					position++
					goto l707
				l708:
					position, tokenIndex, depth = position707, tokenIndex707, depth707
					if buffer[position] != rune('E') {
						goto l665
					}
					position++
				}
			l707:
				if !_rules[ruleAction27]() {
					goto l665
				}
				depth--
				add(ruleCountBasedSampling, position666)
			}
			return true
		l665:
			position, tokenIndex, depth = position665, tokenIndex665, depth665
			return false
		},
		/* 36 RandomizedSampling <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') sp (FloatLiteral / NumericLiteral) spOpt '%' Action28)> */
		func() bool {
			position709, tokenIndex709, depth709 := position, tokenIndex, depth
			{
				position710 := position
				depth++
				{
					position711, tokenIndex711, depth711 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l712
					}
					position++
					goto l711
				l712:
					position, tokenIndex, depth = position711, tokenIndex711, depth711
					if buffer[position] != rune('S') {
						goto l709
					}
					position++
				}
			l711:
				{
					position713, tokenIndex713, depth713 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l714
					}
					position++
					goto l713
				l714:
					position, tokenIndex, depth = position713, tokenIndex713, depth713
					if buffer[position] != rune('A') {
						goto l709
					}
					position++
				}
			l713:
				{
					position715, tokenIndex715, depth715 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l716
					}
					position++
					goto l715
				l716:
					position, tokenIndex, depth = position715, tokenIndex715, depth715
					if buffer[position] != rune('M') {
						goto l709
					}
					position++
				}
			l715:
				{
					position717, tokenIndex717, depth717 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l718
					}
					position++
					goto l717
				l718:
					position, tokenIndex, depth = position717, tokenIndex717, depth717
					if buffer[position] != rune('P') {
						goto l709
					}
					position++
				}
			l717:
				{
					position719, tokenIndex719, depth719 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l720
					}
					position++
					goto l719
				l720:
					position, tokenIndex, depth = position719, tokenIndex719, depth719
					if buffer[position] != rune('L') {
						goto l709
					}
					position++
				}
			l719:
				{
					position721, tokenIndex721, depth721 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l722
					}
					position++
					goto l721
				l722:
					position, tokenIndex, depth = position721, tokenIndex721, depth721
					if buffer[position] != rune('E') {
						goto l709
					}
					position++
				}
			l721:
				if !_rules[rulesp]() {
					goto l709
				}
				{
					position723, tokenIndex723, depth723 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l724
					}
					goto l723
				l724:
					position, tokenIndex, depth = position723, tokenIndex723, depth723
					if !_rules[ruleNumericLiteral]() {
						goto l709
					}
				}
			l723:
				if !_rules[rulespOpt]() {
					goto l709
				}
				if buffer[position] != rune('%') {
					goto l709
				}
				position++
				if !_rules[ruleAction28]() {
					goto l709
				}
				depth--
				add(ruleRandomizedSampling, position710)
			}
			return true
		l709:
			position, tokenIndex, depth = position709, tokenIndex709, depth709
			return false
		},
		/* 37 TimeBasedSampling <- <(TimeBasedSamplingSeconds / TimeBasedSamplingMilliseconds)> */
		func() bool {
			position725, tokenIndex725, depth725 := position, tokenIndex, depth
			{
				position726 := position
				depth++
				{
					position727, tokenIndex727, depth727 := position, tokenIndex, depth
					if !_rules[ruleTimeBasedSamplingSeconds]() {
						goto l728
					}
					goto l727
				l728:
					position, tokenIndex, depth = position727, tokenIndex727, depth727
					if !_rules[ruleTimeBasedSamplingMilliseconds]() {
						goto l725
					}
				}
			l727:
				depth--
				add(ruleTimeBasedSampling, position726)
			}
			return true
		l725:
			position, tokenIndex, depth = position725, tokenIndex725, depth725
			return false
		},
		/* 38 TimeBasedSamplingSeconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action29)> */
		func() bool {
			position729, tokenIndex729, depth729 := position, tokenIndex, depth
			{
				position730 := position
				depth++
				{
					position731, tokenIndex731, depth731 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l732
					}
					position++
					goto l731
				l732:
					position, tokenIndex, depth = position731, tokenIndex731, depth731
					if buffer[position] != rune('E') {
						goto l729
					}
					position++
				}
			l731:
				{
					position733, tokenIndex733, depth733 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l734
					}
					position++
					goto l733
				l734:
					position, tokenIndex, depth = position733, tokenIndex733, depth733
					if buffer[position] != rune('V') {
						goto l729
					}
					position++
				}
			l733:
				{
					position735, tokenIndex735, depth735 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l736
					}
					position++
					goto l735
				l736:
					position, tokenIndex, depth = position735, tokenIndex735, depth735
					if buffer[position] != rune('E') {
						goto l729
					}
					position++
				}
			l735:
				{
					position737, tokenIndex737, depth737 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l738
					}
					position++
					goto l737
				l738:
					position, tokenIndex, depth = position737, tokenIndex737, depth737
					if buffer[position] != rune('R') {
						goto l729
					}
					position++
				}
			l737:
				{
					position739, tokenIndex739, depth739 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l740
					}
					position++
					goto l739
				l740:
					position, tokenIndex, depth = position739, tokenIndex739, depth739
					if buffer[position] != rune('Y') {
						goto l729
					}
					position++
				}
			l739:
				if !_rules[rulesp]() {
					goto l729
				}
				{
					position741, tokenIndex741, depth741 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l742
					}
					goto l741
				l742:
					position, tokenIndex, depth = position741, tokenIndex741, depth741
					if !_rules[ruleNumericLiteral]() {
						goto l729
					}
				}
			l741:
				if !_rules[rulesp]() {
					goto l729
				}
				{
					position743, tokenIndex743, depth743 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l744
					}
					position++
					goto l743
				l744:
					position, tokenIndex, depth = position743, tokenIndex743, depth743
					if buffer[position] != rune('S') {
						goto l729
					}
					position++
				}
			l743:
				{
					position745, tokenIndex745, depth745 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l746
					}
					position++
					goto l745
				l746:
					position, tokenIndex, depth = position745, tokenIndex745, depth745
					if buffer[position] != rune('E') {
						goto l729
					}
					position++
				}
			l745:
				{
					position747, tokenIndex747, depth747 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l748
					}
					position++
					goto l747
				l748:
					position, tokenIndex, depth = position747, tokenIndex747, depth747
					if buffer[position] != rune('C') {
						goto l729
					}
					position++
				}
			l747:
				{
					position749, tokenIndex749, depth749 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l750
					}
					position++
					goto l749
				l750:
					position, tokenIndex, depth = position749, tokenIndex749, depth749
					if buffer[position] != rune('O') {
						goto l729
					}
					position++
				}
			l749:
				{
					position751, tokenIndex751, depth751 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l752
					}
					position++
					goto l751
				l752:
					position, tokenIndex, depth = position751, tokenIndex751, depth751
					if buffer[position] != rune('N') {
						goto l729
					}
					position++
				}
			l751:
				{
					position753, tokenIndex753, depth753 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l754
					}
					position++
					goto l753
				l754:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					if buffer[position] != rune('D') {
						goto l729
					}
					position++
				}
			l753:
				{
					position755, tokenIndex755, depth755 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l756
					}
					position++
					goto l755
				l756:
					position, tokenIndex, depth = position755, tokenIndex755, depth755
					if buffer[position] != rune('S') {
						goto l729
					}
					position++
				}
			l755:
				if !_rules[ruleAction29]() {
					goto l729
				}
				depth--
				add(ruleTimeBasedSamplingSeconds, position730)
			}
			return true
		l729:
			position, tokenIndex, depth = position729, tokenIndex729, depth729
			return false
		},
		/* 39 TimeBasedSamplingMilliseconds <- <(('e' / 'E') ('v' / 'V') ('e' / 'E') ('r' / 'R') ('y' / 'Y') sp (FloatLiteral / NumericLiteral) sp (('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S')) Action30)> */
		func() bool {
			position757, tokenIndex757, depth757 := position, tokenIndex, depth
			{
				position758 := position
				depth++
				{
					position759, tokenIndex759, depth759 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l760
					}
					position++
					goto l759
				l760:
					position, tokenIndex, depth = position759, tokenIndex759, depth759
					if buffer[position] != rune('E') {
						goto l757
					}
					position++
				}
			l759:
				{
					position761, tokenIndex761, depth761 := position, tokenIndex, depth
					if buffer[position] != rune('v') {
						goto l762
					}
					position++
					goto l761
				l762:
					position, tokenIndex, depth = position761, tokenIndex761, depth761
					if buffer[position] != rune('V') {
						goto l757
					}
					position++
				}
			l761:
				{
					position763, tokenIndex763, depth763 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l764
					}
					position++
					goto l763
				l764:
					position, tokenIndex, depth = position763, tokenIndex763, depth763
					if buffer[position] != rune('E') {
						goto l757
					}
					position++
				}
			l763:
				{
					position765, tokenIndex765, depth765 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l766
					}
					position++
					goto l765
				l766:
					position, tokenIndex, depth = position765, tokenIndex765, depth765
					if buffer[position] != rune('R') {
						goto l757
					}
					position++
				}
			l765:
				{
					position767, tokenIndex767, depth767 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l768
					}
					position++
					goto l767
				l768:
					position, tokenIndex, depth = position767, tokenIndex767, depth767
					if buffer[position] != rune('Y') {
						goto l757
					}
					position++
				}
			l767:
				if !_rules[rulesp]() {
					goto l757
				}
				{
					position769, tokenIndex769, depth769 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l770
					}
					goto l769
				l770:
					position, tokenIndex, depth = position769, tokenIndex769, depth769
					if !_rules[ruleNumericLiteral]() {
						goto l757
					}
				}
			l769:
				if !_rules[rulesp]() {
					goto l757
				}
				{
					position771, tokenIndex771, depth771 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l772
					}
					position++
					goto l771
				l772:
					position, tokenIndex, depth = position771, tokenIndex771, depth771
					if buffer[position] != rune('M') {
						goto l757
					}
					position++
				}
			l771:
				{
					position773, tokenIndex773, depth773 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l774
					}
					position++
					goto l773
				l774:
					position, tokenIndex, depth = position773, tokenIndex773, depth773
					if buffer[position] != rune('I') {
						goto l757
					}
					position++
				}
			l773:
				{
					position775, tokenIndex775, depth775 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l776
					}
					position++
					goto l775
				l776:
					position, tokenIndex, depth = position775, tokenIndex775, depth775
					if buffer[position] != rune('L') {
						goto l757
					}
					position++
				}
			l775:
				{
					position777, tokenIndex777, depth777 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l778
					}
					position++
					goto l777
				l778:
					position, tokenIndex, depth = position777, tokenIndex777, depth777
					if buffer[position] != rune('L') {
						goto l757
					}
					position++
				}
			l777:
				{
					position779, tokenIndex779, depth779 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l780
					}
					position++
					goto l779
				l780:
					position, tokenIndex, depth = position779, tokenIndex779, depth779
					if buffer[position] != rune('I') {
						goto l757
					}
					position++
				}
			l779:
				{
					position781, tokenIndex781, depth781 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l782
					}
					position++
					goto l781
				l782:
					position, tokenIndex, depth = position781, tokenIndex781, depth781
					if buffer[position] != rune('S') {
						goto l757
					}
					position++
				}
			l781:
				{
					position783, tokenIndex783, depth783 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l784
					}
					position++
					goto l783
				l784:
					position, tokenIndex, depth = position783, tokenIndex783, depth783
					if buffer[position] != rune('E') {
						goto l757
					}
					position++
				}
			l783:
				{
					position785, tokenIndex785, depth785 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l786
					}
					position++
					goto l785
				l786:
					position, tokenIndex, depth = position785, tokenIndex785, depth785
					if buffer[position] != rune('C') {
						goto l757
					}
					position++
				}
			l785:
				{
					position787, tokenIndex787, depth787 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l788
					}
					position++
					goto l787
				l788:
					position, tokenIndex, depth = position787, tokenIndex787, depth787
					if buffer[position] != rune('O') {
						goto l757
					}
					position++
				}
			l787:
				{
					position789, tokenIndex789, depth789 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l790
					}
					position++
					goto l789
				l790:
					position, tokenIndex, depth = position789, tokenIndex789, depth789
					if buffer[position] != rune('N') {
						goto l757
					}
					position++
				}
			l789:
				{
					position791, tokenIndex791, depth791 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l792
					}
					position++
					goto l791
				l792:
					position, tokenIndex, depth = position791, tokenIndex791, depth791
					if buffer[position] != rune('D') {
						goto l757
					}
					position++
				}
			l791:
				{
					position793, tokenIndex793, depth793 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l794
					}
					position++
					goto l793
				l794:
					position, tokenIndex, depth = position793, tokenIndex793, depth793
					if buffer[position] != rune('S') {
						goto l757
					}
					position++
				}
			l793:
				if !_rules[ruleAction30]() {
					goto l757
				}
				depth--
				add(ruleTimeBasedSamplingMilliseconds, position758)
			}
			return true
		l757:
			position, tokenIndex, depth = position757, tokenIndex757, depth757
			return false
		},
		/* 40 Projections <- <(<(sp Projection (spOpt ',' spOpt Projection)*)> Action31)> */
		func() bool {
			position795, tokenIndex795, depth795 := position, tokenIndex, depth
			{
				position796 := position
				depth++
				{
					position797 := position
					depth++
					if !_rules[rulesp]() {
						goto l795
					}
					if !_rules[ruleProjection]() {
						goto l795
					}
				l798:
					{
						position799, tokenIndex799, depth799 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l799
						}
						if buffer[position] != rune(',') {
							goto l799
						}
						position++
						if !_rules[rulespOpt]() {
							goto l799
						}
						if !_rules[ruleProjection]() {
							goto l799
						}
						goto l798
					l799:
						position, tokenIndex, depth = position799, tokenIndex799, depth799
					}
					depth--
					add(rulePegText, position797)
				}
				if !_rules[ruleAction31]() {
					goto l795
				}
				depth--
				add(ruleProjections, position796)
			}
			return true
		l795:
			position, tokenIndex, depth = position795, tokenIndex795, depth795
			return false
		},
		/* 41 Projection <- <(AliasExpression / ExpressionOrWildcard)> */
		func() bool {
			position800, tokenIndex800, depth800 := position, tokenIndex, depth
			{
				position801 := position
				depth++
				{
					position802, tokenIndex802, depth802 := position, tokenIndex, depth
					if !_rules[ruleAliasExpression]() {
						goto l803
					}
					goto l802
				l803:
					position, tokenIndex, depth = position802, tokenIndex802, depth802
					if !_rules[ruleExpressionOrWildcard]() {
						goto l800
					}
				}
			l802:
				depth--
				add(ruleProjection, position801)
			}
			return true
		l800:
			position, tokenIndex, depth = position800, tokenIndex800, depth800
			return false
		},
		/* 42 AliasExpression <- <(ExpressionOrWildcard sp (('a' / 'A') ('s' / 'S')) sp TargetIdentifier Action32)> */
		func() bool {
			position804, tokenIndex804, depth804 := position, tokenIndex, depth
			{
				position805 := position
				depth++
				if !_rules[ruleExpressionOrWildcard]() {
					goto l804
				}
				if !_rules[rulesp]() {
					goto l804
				}
				{
					position806, tokenIndex806, depth806 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l807
					}
					position++
					goto l806
				l807:
					position, tokenIndex, depth = position806, tokenIndex806, depth806
					if buffer[position] != rune('A') {
						goto l804
					}
					position++
				}
			l806:
				{
					position808, tokenIndex808, depth808 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l809
					}
					position++
					goto l808
				l809:
					position, tokenIndex, depth = position808, tokenIndex808, depth808
					if buffer[position] != rune('S') {
						goto l804
					}
					position++
				}
			l808:
				if !_rules[rulesp]() {
					goto l804
				}
				if !_rules[ruleTargetIdentifier]() {
					goto l804
				}
				if !_rules[ruleAction32]() {
					goto l804
				}
				depth--
				add(ruleAliasExpression, position805)
			}
			return true
		l804:
			position, tokenIndex, depth = position804, tokenIndex804, depth804
			return false
		},
		/* 43 WindowedFrom <- <(<(sp (('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M')) sp Relations)?> Action33)> */
		func() bool {
			position810, tokenIndex810, depth810 := position, tokenIndex, depth
			{
				position811 := position
				depth++
				{
					position812 := position
					depth++
					{
						position813, tokenIndex813, depth813 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l813
						}
						{
							position815, tokenIndex815, depth815 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l816
							}
							position++
							goto l815
						l816:
							position, tokenIndex, depth = position815, tokenIndex815, depth815
							if buffer[position] != rune('F') {
								goto l813
							}
							position++
						}
					l815:
						{
							position817, tokenIndex817, depth817 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l818
							}
							position++
							goto l817
						l818:
							position, tokenIndex, depth = position817, tokenIndex817, depth817
							if buffer[position] != rune('R') {
								goto l813
							}
							position++
						}
					l817:
						{
							position819, tokenIndex819, depth819 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l820
							}
							position++
							goto l819
						l820:
							position, tokenIndex, depth = position819, tokenIndex819, depth819
							if buffer[position] != rune('O') {
								goto l813
							}
							position++
						}
					l819:
						{
							position821, tokenIndex821, depth821 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l822
							}
							position++
							goto l821
						l822:
							position, tokenIndex, depth = position821, tokenIndex821, depth821
							if buffer[position] != rune('M') {
								goto l813
							}
							position++
						}
					l821:
						if !_rules[rulesp]() {
							goto l813
						}
						if !_rules[ruleRelations]() {
							goto l813
						}
						goto l814
					l813:
						position, tokenIndex, depth = position813, tokenIndex813, depth813
					}
				l814:
					depth--
					add(rulePegText, position812)
				}
				if !_rules[ruleAction33]() {
					goto l810
				}
				depth--
				add(ruleWindowedFrom, position811)
			}
			return true
		l810:
			position, tokenIndex, depth = position810, tokenIndex810, depth810
			return false
		},
		/* 44 Interval <- <(TimeInterval / TuplesInterval)> */
		func() bool {
			position823, tokenIndex823, depth823 := position, tokenIndex, depth
			{
				position824 := position
				depth++
				{
					position825, tokenIndex825, depth825 := position, tokenIndex, depth
					if !_rules[ruleTimeInterval]() {
						goto l826
					}
					goto l825
				l826:
					position, tokenIndex, depth = position825, tokenIndex825, depth825
					if !_rules[ruleTuplesInterval]() {
						goto l823
					}
				}
			l825:
				depth--
				add(ruleInterval, position824)
			}
			return true
		l823:
			position, tokenIndex, depth = position823, tokenIndex823, depth823
			return false
		},
		/* 45 TimeInterval <- <((FloatLiteral / NumericLiteral) sp (SECONDS / MILLISECONDS) Action34)> */
		func() bool {
			position827, tokenIndex827, depth827 := position, tokenIndex, depth
			{
				position828 := position
				depth++
				{
					position829, tokenIndex829, depth829 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l830
					}
					goto l829
				l830:
					position, tokenIndex, depth = position829, tokenIndex829, depth829
					if !_rules[ruleNumericLiteral]() {
						goto l827
					}
				}
			l829:
				if !_rules[rulesp]() {
					goto l827
				}
				{
					position831, tokenIndex831, depth831 := position, tokenIndex, depth
					if !_rules[ruleSECONDS]() {
						goto l832
					}
					goto l831
				l832:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					if !_rules[ruleMILLISECONDS]() {
						goto l827
					}
				}
			l831:
				if !_rules[ruleAction34]() {
					goto l827
				}
				depth--
				add(ruleTimeInterval, position828)
			}
			return true
		l827:
			position, tokenIndex, depth = position827, tokenIndex827, depth827
			return false
		},
		/* 46 TuplesInterval <- <(NumericLiteral sp TUPLES Action35)> */
		func() bool {
			position833, tokenIndex833, depth833 := position, tokenIndex, depth
			{
				position834 := position
				depth++
				if !_rules[ruleNumericLiteral]() {
					goto l833
				}
				if !_rules[rulesp]() {
					goto l833
				}
				if !_rules[ruleTUPLES]() {
					goto l833
				}
				if !_rules[ruleAction35]() {
					goto l833
				}
				depth--
				add(ruleTuplesInterval, position834)
			}
			return true
		l833:
			position, tokenIndex, depth = position833, tokenIndex833, depth833
			return false
		},
		/* 47 Relations <- <(RelationLike (spOpt ',' spOpt RelationLike)*)> */
		func() bool {
			position835, tokenIndex835, depth835 := position, tokenIndex, depth
			{
				position836 := position
				depth++
				if !_rules[ruleRelationLike]() {
					goto l835
				}
			l837:
				{
					position838, tokenIndex838, depth838 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l838
					}
					if buffer[position] != rune(',') {
						goto l838
					}
					position++
					if !_rules[rulespOpt]() {
						goto l838
					}
					if !_rules[ruleRelationLike]() {
						goto l838
					}
					goto l837
				l838:
					position, tokenIndex, depth = position838, tokenIndex838, depth838
				}
				depth--
				add(ruleRelations, position836)
			}
			return true
		l835:
			position, tokenIndex, depth = position835, tokenIndex835, depth835
			return false
		},
		/* 48 Filter <- <(<(sp (('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E')) sp Expression)?> Action36)> */
		func() bool {
			position839, tokenIndex839, depth839 := position, tokenIndex, depth
			{
				position840 := position
				depth++
				{
					position841 := position
					depth++
					{
						position842, tokenIndex842, depth842 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l842
						}
						{
							position844, tokenIndex844, depth844 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l845
							}
							position++
							goto l844
						l845:
							position, tokenIndex, depth = position844, tokenIndex844, depth844
							if buffer[position] != rune('W') {
								goto l842
							}
							position++
						}
					l844:
						{
							position846, tokenIndex846, depth846 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l847
							}
							position++
							goto l846
						l847:
							position, tokenIndex, depth = position846, tokenIndex846, depth846
							if buffer[position] != rune('H') {
								goto l842
							}
							position++
						}
					l846:
						{
							position848, tokenIndex848, depth848 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l849
							}
							position++
							goto l848
						l849:
							position, tokenIndex, depth = position848, tokenIndex848, depth848
							if buffer[position] != rune('E') {
								goto l842
							}
							position++
						}
					l848:
						{
							position850, tokenIndex850, depth850 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l851
							}
							position++
							goto l850
						l851:
							position, tokenIndex, depth = position850, tokenIndex850, depth850
							if buffer[position] != rune('R') {
								goto l842
							}
							position++
						}
					l850:
						{
							position852, tokenIndex852, depth852 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l853
							}
							position++
							goto l852
						l853:
							position, tokenIndex, depth = position852, tokenIndex852, depth852
							if buffer[position] != rune('E') {
								goto l842
							}
							position++
						}
					l852:
						if !_rules[rulesp]() {
							goto l842
						}
						if !_rules[ruleExpression]() {
							goto l842
						}
						goto l843
					l842:
						position, tokenIndex, depth = position842, tokenIndex842, depth842
					}
				l843:
					depth--
					add(rulePegText, position841)
				}
				if !_rules[ruleAction36]() {
					goto l839
				}
				depth--
				add(ruleFilter, position840)
			}
			return true
		l839:
			position, tokenIndex, depth = position839, tokenIndex839, depth839
			return false
		},
		/* 49 Grouping <- <(<(sp (('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P')) sp (('b' / 'B') ('y' / 'Y')) sp GroupList)?> Action37)> */
		func() bool {
			position854, tokenIndex854, depth854 := position, tokenIndex, depth
			{
				position855 := position
				depth++
				{
					position856 := position
					depth++
					{
						position857, tokenIndex857, depth857 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l857
						}
						{
							position859, tokenIndex859, depth859 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l860
							}
							position++
							goto l859
						l860:
							position, tokenIndex, depth = position859, tokenIndex859, depth859
							if buffer[position] != rune('G') {
								goto l857
							}
							position++
						}
					l859:
						{
							position861, tokenIndex861, depth861 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l862
							}
							position++
							goto l861
						l862:
							position, tokenIndex, depth = position861, tokenIndex861, depth861
							if buffer[position] != rune('R') {
								goto l857
							}
							position++
						}
					l861:
						{
							position863, tokenIndex863, depth863 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l864
							}
							position++
							goto l863
						l864:
							position, tokenIndex, depth = position863, tokenIndex863, depth863
							if buffer[position] != rune('O') {
								goto l857
							}
							position++
						}
					l863:
						{
							position865, tokenIndex865, depth865 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l866
							}
							position++
							goto l865
						l866:
							position, tokenIndex, depth = position865, tokenIndex865, depth865
							if buffer[position] != rune('U') {
								goto l857
							}
							position++
						}
					l865:
						{
							position867, tokenIndex867, depth867 := position, tokenIndex, depth
							if buffer[position] != rune('p') {
								goto l868
							}
							position++
							goto l867
						l868:
							position, tokenIndex, depth = position867, tokenIndex867, depth867
							if buffer[position] != rune('P') {
								goto l857
							}
							position++
						}
					l867:
						if !_rules[rulesp]() {
							goto l857
						}
						{
							position869, tokenIndex869, depth869 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l870
							}
							position++
							goto l869
						l870:
							position, tokenIndex, depth = position869, tokenIndex869, depth869
							if buffer[position] != rune('B') {
								goto l857
							}
							position++
						}
					l869:
						{
							position871, tokenIndex871, depth871 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l872
							}
							position++
							goto l871
						l872:
							position, tokenIndex, depth = position871, tokenIndex871, depth871
							if buffer[position] != rune('Y') {
								goto l857
							}
							position++
						}
					l871:
						if !_rules[rulesp]() {
							goto l857
						}
						if !_rules[ruleGroupList]() {
							goto l857
						}
						goto l858
					l857:
						position, tokenIndex, depth = position857, tokenIndex857, depth857
					}
				l858:
					depth--
					add(rulePegText, position856)
				}
				if !_rules[ruleAction37]() {
					goto l854
				}
				depth--
				add(ruleGrouping, position855)
			}
			return true
		l854:
			position, tokenIndex, depth = position854, tokenIndex854, depth854
			return false
		},
		/* 50 GroupList <- <(Expression (spOpt ',' spOpt Expression)*)> */
		func() bool {
			position873, tokenIndex873, depth873 := position, tokenIndex, depth
			{
				position874 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l873
				}
			l875:
				{
					position876, tokenIndex876, depth876 := position, tokenIndex, depth
					if !_rules[rulespOpt]() {
						goto l876
					}
					if buffer[position] != rune(',') {
						goto l876
					}
					position++
					if !_rules[rulespOpt]() {
						goto l876
					}
					if !_rules[ruleExpression]() {
						goto l876
					}
					goto l875
				l876:
					position, tokenIndex, depth = position876, tokenIndex876, depth876
				}
				depth--
				add(ruleGroupList, position874)
			}
			return true
		l873:
			position, tokenIndex, depth = position873, tokenIndex873, depth873
			return false
		},
		/* 51 Having <- <(<(sp (('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G')) sp Expression)?> Action38)> */
		func() bool {
			position877, tokenIndex877, depth877 := position, tokenIndex, depth
			{
				position878 := position
				depth++
				{
					position879 := position
					depth++
					{
						position880, tokenIndex880, depth880 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l880
						}
						{
							position882, tokenIndex882, depth882 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l883
							}
							position++
							goto l882
						l883:
							position, tokenIndex, depth = position882, tokenIndex882, depth882
							if buffer[position] != rune('H') {
								goto l880
							}
							position++
						}
					l882:
						{
							position884, tokenIndex884, depth884 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l885
							}
							position++
							goto l884
						l885:
							position, tokenIndex, depth = position884, tokenIndex884, depth884
							if buffer[position] != rune('A') {
								goto l880
							}
							position++
						}
					l884:
						{
							position886, tokenIndex886, depth886 := position, tokenIndex, depth
							if buffer[position] != rune('v') {
								goto l887
							}
							position++
							goto l886
						l887:
							position, tokenIndex, depth = position886, tokenIndex886, depth886
							if buffer[position] != rune('V') {
								goto l880
							}
							position++
						}
					l886:
						{
							position888, tokenIndex888, depth888 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l889
							}
							position++
							goto l888
						l889:
							position, tokenIndex, depth = position888, tokenIndex888, depth888
							if buffer[position] != rune('I') {
								goto l880
							}
							position++
						}
					l888:
						{
							position890, tokenIndex890, depth890 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l891
							}
							position++
							goto l890
						l891:
							position, tokenIndex, depth = position890, tokenIndex890, depth890
							if buffer[position] != rune('N') {
								goto l880
							}
							position++
						}
					l890:
						{
							position892, tokenIndex892, depth892 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l893
							}
							position++
							goto l892
						l893:
							position, tokenIndex, depth = position892, tokenIndex892, depth892
							if buffer[position] != rune('G') {
								goto l880
							}
							position++
						}
					l892:
						if !_rules[rulesp]() {
							goto l880
						}
						if !_rules[ruleExpression]() {
							goto l880
						}
						goto l881
					l880:
						position, tokenIndex, depth = position880, tokenIndex880, depth880
					}
				l881:
					depth--
					add(rulePegText, position879)
				}
				if !_rules[ruleAction38]() {
					goto l877
				}
				depth--
				add(ruleHaving, position878)
			}
			return true
		l877:
			position, tokenIndex, depth = position877, tokenIndex877, depth877
			return false
		},
		/* 52 RelationLike <- <(AliasedStreamWindow / (StreamWindow Action39))> */
		func() bool {
			position894, tokenIndex894, depth894 := position, tokenIndex, depth
			{
				position895 := position
				depth++
				{
					position896, tokenIndex896, depth896 := position, tokenIndex, depth
					if !_rules[ruleAliasedStreamWindow]() {
						goto l897
					}
					goto l896
				l897:
					position, tokenIndex, depth = position896, tokenIndex896, depth896
					if !_rules[ruleStreamWindow]() {
						goto l894
					}
					if !_rules[ruleAction39]() {
						goto l894
					}
				}
			l896:
				depth--
				add(ruleRelationLike, position895)
			}
			return true
		l894:
			position, tokenIndex, depth = position894, tokenIndex894, depth894
			return false
		},
		/* 53 AliasedStreamWindow <- <(StreamWindow sp (('a' / 'A') ('s' / 'S')) sp Identifier Action40)> */
		func() bool {
			position898, tokenIndex898, depth898 := position, tokenIndex, depth
			{
				position899 := position
				depth++
				if !_rules[ruleStreamWindow]() {
					goto l898
				}
				if !_rules[rulesp]() {
					goto l898
				}
				{
					position900, tokenIndex900, depth900 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l901
					}
					position++
					goto l900
				l901:
					position, tokenIndex, depth = position900, tokenIndex900, depth900
					if buffer[position] != rune('A') {
						goto l898
					}
					position++
				}
			l900:
				{
					position902, tokenIndex902, depth902 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l903
					}
					position++
					goto l902
				l903:
					position, tokenIndex, depth = position902, tokenIndex902, depth902
					if buffer[position] != rune('S') {
						goto l898
					}
					position++
				}
			l902:
				if !_rules[rulesp]() {
					goto l898
				}
				if !_rules[ruleIdentifier]() {
					goto l898
				}
				if !_rules[ruleAction40]() {
					goto l898
				}
				depth--
				add(ruleAliasedStreamWindow, position899)
			}
			return true
		l898:
			position, tokenIndex, depth = position898, tokenIndex898, depth898
			return false
		},
		/* 54 StreamWindow <- <(StreamLike spOpt '[' spOpt (('r' / 'R') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('e' / 'E')) sp Interval CapacitySpecOpt SheddingSpecOpt spOpt ']' Action41)> */
		func() bool {
			position904, tokenIndex904, depth904 := position, tokenIndex, depth
			{
				position905 := position
				depth++
				if !_rules[ruleStreamLike]() {
					goto l904
				}
				if !_rules[rulespOpt]() {
					goto l904
				}
				if buffer[position] != rune('[') {
					goto l904
				}
				position++
				if !_rules[rulespOpt]() {
					goto l904
				}
				{
					position906, tokenIndex906, depth906 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l907
					}
					position++
					goto l906
				l907:
					position, tokenIndex, depth = position906, tokenIndex906, depth906
					if buffer[position] != rune('R') {
						goto l904
					}
					position++
				}
			l906:
				{
					position908, tokenIndex908, depth908 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l909
					}
					position++
					goto l908
				l909:
					position, tokenIndex, depth = position908, tokenIndex908, depth908
					if buffer[position] != rune('A') {
						goto l904
					}
					position++
				}
			l908:
				{
					position910, tokenIndex910, depth910 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l911
					}
					position++
					goto l910
				l911:
					position, tokenIndex, depth = position910, tokenIndex910, depth910
					if buffer[position] != rune('N') {
						goto l904
					}
					position++
				}
			l910:
				{
					position912, tokenIndex912, depth912 := position, tokenIndex, depth
					if buffer[position] != rune('g') {
						goto l913
					}
					position++
					goto l912
				l913:
					position, tokenIndex, depth = position912, tokenIndex912, depth912
					if buffer[position] != rune('G') {
						goto l904
					}
					position++
				}
			l912:
				{
					position914, tokenIndex914, depth914 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l915
					}
					position++
					goto l914
				l915:
					position, tokenIndex, depth = position914, tokenIndex914, depth914
					if buffer[position] != rune('E') {
						goto l904
					}
					position++
				}
			l914:
				if !_rules[rulesp]() {
					goto l904
				}
				if !_rules[ruleInterval]() {
					goto l904
				}
				if !_rules[ruleCapacitySpecOpt]() {
					goto l904
				}
				if !_rules[ruleSheddingSpecOpt]() {
					goto l904
				}
				if !_rules[rulespOpt]() {
					goto l904
				}
				if buffer[position] != rune(']') {
					goto l904
				}
				position++
				if !_rules[ruleAction41]() {
					goto l904
				}
				depth--
				add(ruleStreamWindow, position905)
			}
			return true
		l904:
			position, tokenIndex, depth = position904, tokenIndex904, depth904
			return false
		},
		/* 55 StreamLike <- <(UDSFFuncApp / Stream)> */
		func() bool {
			position916, tokenIndex916, depth916 := position, tokenIndex, depth
			{
				position917 := position
				depth++
				{
					position918, tokenIndex918, depth918 := position, tokenIndex, depth
					if !_rules[ruleUDSFFuncApp]() {
						goto l919
					}
					goto l918
				l919:
					position, tokenIndex, depth = position918, tokenIndex918, depth918
					if !_rules[ruleStream]() {
						goto l916
					}
				}
			l918:
				depth--
				add(ruleStreamLike, position917)
			}
			return true
		l916:
			position, tokenIndex, depth = position916, tokenIndex916, depth916
			return false
		},
		/* 56 UDSFFuncApp <- <(FuncAppWithoutOrderBy Action42)> */
		func() bool {
			position920, tokenIndex920, depth920 := position, tokenIndex, depth
			{
				position921 := position
				depth++
				if !_rules[ruleFuncAppWithoutOrderBy]() {
					goto l920
				}
				if !_rules[ruleAction42]() {
					goto l920
				}
				depth--
				add(ruleUDSFFuncApp, position921)
			}
			return true
		l920:
			position, tokenIndex, depth = position920, tokenIndex920, depth920
			return false
		},
		/* 57 CapacitySpecOpt <- <(<(spOpt ',' spOpt (('b' / 'B') ('u' / 'U') ('f' / 'F') ('f' / 'F') ('e' / 'E') ('r' / 'R')) sp (('s' / 'S') ('i' / 'I') ('z' / 'Z') ('e' / 'E')) sp NonNegativeNumericLiteral)?> Action43)> */
		func() bool {
			position922, tokenIndex922, depth922 := position, tokenIndex, depth
			{
				position923 := position
				depth++
				{
					position924 := position
					depth++
					{
						position925, tokenIndex925, depth925 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l925
						}
						if buffer[position] != rune(',') {
							goto l925
						}
						position++
						if !_rules[rulespOpt]() {
							goto l925
						}
						{
							position927, tokenIndex927, depth927 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l928
							}
							position++
							goto l927
						l928:
							position, tokenIndex, depth = position927, tokenIndex927, depth927
							if buffer[position] != rune('B') {
								goto l925
							}
							position++
						}
					l927:
						{
							position929, tokenIndex929, depth929 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l930
							}
							position++
							goto l929
						l930:
							position, tokenIndex, depth = position929, tokenIndex929, depth929
							if buffer[position] != rune('U') {
								goto l925
							}
							position++
						}
					l929:
						{
							position931, tokenIndex931, depth931 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l932
							}
							position++
							goto l931
						l932:
							position, tokenIndex, depth = position931, tokenIndex931, depth931
							if buffer[position] != rune('F') {
								goto l925
							}
							position++
						}
					l931:
						{
							position933, tokenIndex933, depth933 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l934
							}
							position++
							goto l933
						l934:
							position, tokenIndex, depth = position933, tokenIndex933, depth933
							if buffer[position] != rune('F') {
								goto l925
							}
							position++
						}
					l933:
						{
							position935, tokenIndex935, depth935 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l936
							}
							position++
							goto l935
						l936:
							position, tokenIndex, depth = position935, tokenIndex935, depth935
							if buffer[position] != rune('E') {
								goto l925
							}
							position++
						}
					l935:
						{
							position937, tokenIndex937, depth937 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l938
							}
							position++
							goto l937
						l938:
							position, tokenIndex, depth = position937, tokenIndex937, depth937
							if buffer[position] != rune('R') {
								goto l925
							}
							position++
						}
					l937:
						if !_rules[rulesp]() {
							goto l925
						}
						{
							position939, tokenIndex939, depth939 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l940
							}
							position++
							goto l939
						l940:
							position, tokenIndex, depth = position939, tokenIndex939, depth939
							if buffer[position] != rune('S') {
								goto l925
							}
							position++
						}
					l939:
						{
							position941, tokenIndex941, depth941 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l942
							}
							position++
							goto l941
						l942:
							position, tokenIndex, depth = position941, tokenIndex941, depth941
							if buffer[position] != rune('I') {
								goto l925
							}
							position++
						}
					l941:
						{
							position943, tokenIndex943, depth943 := position, tokenIndex, depth
							if buffer[position] != rune('z') {
								goto l944
							}
							position++
							goto l943
						l944:
							position, tokenIndex, depth = position943, tokenIndex943, depth943
							if buffer[position] != rune('Z') {
								goto l925
							}
							position++
						}
					l943:
						{
							position945, tokenIndex945, depth945 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l946
							}
							position++
							goto l945
						l946:
							position, tokenIndex, depth = position945, tokenIndex945, depth945
							if buffer[position] != rune('E') {
								goto l925
							}
							position++
						}
					l945:
						if !_rules[rulesp]() {
							goto l925
						}
						if !_rules[ruleNonNegativeNumericLiteral]() {
							goto l925
						}
						goto l926
					l925:
						position, tokenIndex, depth = position925, tokenIndex925, depth925
					}
				l926:
					depth--
					add(rulePegText, position924)
				}
				if !_rules[ruleAction43]() {
					goto l922
				}
				depth--
				add(ruleCapacitySpecOpt, position923)
			}
			return true
		l922:
			position, tokenIndex, depth = position922, tokenIndex922, depth922
			return false
		},
		/* 58 SheddingSpecOpt <- <(<(spOpt ',' spOpt SheddingOption sp (('i' / 'I') ('f' / 'F')) sp (('f' / 'F') ('u' / 'U') ('l' / 'L') ('l' / 'L')))?> Action44)> */
		func() bool {
			position947, tokenIndex947, depth947 := position, tokenIndex, depth
			{
				position948 := position
				depth++
				{
					position949 := position
					depth++
					{
						position950, tokenIndex950, depth950 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l950
						}
						if buffer[position] != rune(',') {
							goto l950
						}
						position++
						if !_rules[rulespOpt]() {
							goto l950
						}
						if !_rules[ruleSheddingOption]() {
							goto l950
						}
						if !_rules[rulesp]() {
							goto l950
						}
						{
							position952, tokenIndex952, depth952 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l953
							}
							position++
							goto l952
						l953:
							position, tokenIndex, depth = position952, tokenIndex952, depth952
							if buffer[position] != rune('I') {
								goto l950
							}
							position++
						}
					l952:
						{
							position954, tokenIndex954, depth954 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l955
							}
							position++
							goto l954
						l955:
							position, tokenIndex, depth = position954, tokenIndex954, depth954
							if buffer[position] != rune('F') {
								goto l950
							}
							position++
						}
					l954:
						if !_rules[rulesp]() {
							goto l950
						}
						{
							position956, tokenIndex956, depth956 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l957
							}
							position++
							goto l956
						l957:
							position, tokenIndex, depth = position956, tokenIndex956, depth956
							if buffer[position] != rune('F') {
								goto l950
							}
							position++
						}
					l956:
						{
							position958, tokenIndex958, depth958 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l959
							}
							position++
							goto l958
						l959:
							position, tokenIndex, depth = position958, tokenIndex958, depth958
							if buffer[position] != rune('U') {
								goto l950
							}
							position++
						}
					l958:
						{
							position960, tokenIndex960, depth960 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l961
							}
							position++
							goto l960
						l961:
							position, tokenIndex, depth = position960, tokenIndex960, depth960
							if buffer[position] != rune('L') {
								goto l950
							}
							position++
						}
					l960:
						{
							position962, tokenIndex962, depth962 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l963
							}
							position++
							goto l962
						l963:
							position, tokenIndex, depth = position962, tokenIndex962, depth962
							if buffer[position] != rune('L') {
								goto l950
							}
							position++
						}
					l962:
						goto l951
					l950:
						position, tokenIndex, depth = position950, tokenIndex950, depth950
					}
				l951:
					depth--
					add(rulePegText, position949)
				}
				if !_rules[ruleAction44]() {
					goto l947
				}
				depth--
				add(ruleSheddingSpecOpt, position948)
			}
			return true
		l947:
			position, tokenIndex, depth = position947, tokenIndex947, depth947
			return false
		},
		/* 59 SheddingOption <- <(Wait / DropOldest / DropNewest)> */
		func() bool {
			position964, tokenIndex964, depth964 := position, tokenIndex, depth
			{
				position965 := position
				depth++
				{
					position966, tokenIndex966, depth966 := position, tokenIndex, depth
					if !_rules[ruleWait]() {
						goto l967
					}
					goto l966
				l967:
					position, tokenIndex, depth = position966, tokenIndex966, depth966
					if !_rules[ruleDropOldest]() {
						goto l968
					}
					goto l966
				l968:
					position, tokenIndex, depth = position966, tokenIndex966, depth966
					if !_rules[ruleDropNewest]() {
						goto l964
					}
				}
			l966:
				depth--
				add(ruleSheddingOption, position965)
			}
			return true
		l964:
			position, tokenIndex, depth = position964, tokenIndex964, depth964
			return false
		},
		/* 60 SourceSinkSpecs <- <(<(sp (('w' / 'W') ('i' / 'I') ('t' / 'T') ('h' / 'H')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action45)> */
		func() bool {
			position969, tokenIndex969, depth969 := position, tokenIndex, depth
			{
				position970 := position
				depth++
				{
					position971 := position
					depth++
					{
						position972, tokenIndex972, depth972 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l972
						}
						{
							position974, tokenIndex974, depth974 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l975
							}
							position++
							goto l974
						l975:
							position, tokenIndex, depth = position974, tokenIndex974, depth974
							if buffer[position] != rune('W') {
								goto l972
							}
							position++
						}
					l974:
						{
							position976, tokenIndex976, depth976 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l977
							}
							position++
							goto l976
						l977:
							position, tokenIndex, depth = position976, tokenIndex976, depth976
							if buffer[position] != rune('I') {
								goto l972
							}
							position++
						}
					l976:
						{
							position978, tokenIndex978, depth978 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l979
							}
							position++
							goto l978
						l979:
							position, tokenIndex, depth = position978, tokenIndex978, depth978
							if buffer[position] != rune('T') {
								goto l972
							}
							position++
						}
					l978:
						{
							position980, tokenIndex980, depth980 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l981
							}
							position++
							goto l980
						l981:
							position, tokenIndex, depth = position980, tokenIndex980, depth980
							if buffer[position] != rune('H') {
								goto l972
							}
							position++
						}
					l980:
						if !_rules[rulesp]() {
							goto l972
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l972
						}
					l982:
						{
							position983, tokenIndex983, depth983 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l983
							}
							if buffer[position] != rune(',') {
								goto l983
							}
							position++
							if !_rules[rulespOpt]() {
								goto l983
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l983
							}
							goto l982
						l983:
							position, tokenIndex, depth = position983, tokenIndex983, depth983
						}
						goto l973
					l972:
						position, tokenIndex, depth = position972, tokenIndex972, depth972
					}
				l973:
					depth--
					add(rulePegText, position971)
				}
				if !_rules[ruleAction45]() {
					goto l969
				}
				depth--
				add(ruleSourceSinkSpecs, position970)
			}
			return true
		l969:
			position, tokenIndex, depth = position969, tokenIndex969, depth969
			return false
		},
		/* 61 UpdateSourceSinkSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)> Action46)> */
		func() bool {
			position984, tokenIndex984, depth984 := position, tokenIndex, depth
			{
				position985 := position
				depth++
				{
					position986 := position
					depth++
					if !_rules[rulesp]() {
						goto l984
					}
					{
						position987, tokenIndex987, depth987 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l988
						}
						position++
						goto l987
					l988:
						position, tokenIndex, depth = position987, tokenIndex987, depth987
						if buffer[position] != rune('S') {
							goto l984
						}
						position++
					}
				l987:
					{
						position989, tokenIndex989, depth989 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l990
						}
						position++
						goto l989
					l990:
						position, tokenIndex, depth = position989, tokenIndex989, depth989
						if buffer[position] != rune('E') {
							goto l984
						}
						position++
					}
				l989:
					{
						position991, tokenIndex991, depth991 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l992
						}
						position++
						goto l991
					l992:
						position, tokenIndex, depth = position991, tokenIndex991, depth991
						if buffer[position] != rune('T') {
							goto l984
						}
						position++
					}
				l991:
					if !_rules[rulesp]() {
						goto l984
					}
					if !_rules[ruleSourceSinkParam]() {
						goto l984
					}
				l993:
					{
						position994, tokenIndex994, depth994 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l994
						}
						if buffer[position] != rune(',') {
							goto l994
						}
						position++
						if !_rules[rulespOpt]() {
							goto l994
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l994
						}
						goto l993
					l994:
						position, tokenIndex, depth = position994, tokenIndex994, depth994
					}
					depth--
					add(rulePegText, position986)
				}
				if !_rules[ruleAction46]() {
					goto l984
				}
				depth--
				add(ruleUpdateSourceSinkSpecs, position985)
			}
			return true
		l984:
			position, tokenIndex, depth = position984, tokenIndex984, depth984
			return false
		},
		/* 62 SetOptSpecs <- <(<(sp (('s' / 'S') ('e' / 'E') ('t' / 'T')) sp SourceSinkParam (spOpt ',' spOpt SourceSinkParam)*)?> Action47)> */
		func() bool {
			position995, tokenIndex995, depth995 := position, tokenIndex, depth
			{
				position996 := position
				depth++
				{
					position997 := position
					depth++
					{
						position998, tokenIndex998, depth998 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l998
						}
						{
							position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1001
							}
							position++
							goto l1000
						l1001:
							position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
							if buffer[position] != rune('S') {
								goto l998
							}
							position++
						}
					l1000:
						{
							position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1003
							}
							position++
							goto l1002
						l1003:
							position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
							if buffer[position] != rune('E') {
								goto l998
							}
							position++
						}
					l1002:
						{
							position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1005
							}
							position++
							goto l1004
						l1005:
							position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
							if buffer[position] != rune('T') {
								goto l998
							}
							position++
						}
					l1004:
						if !_rules[rulesp]() {
							goto l998
						}
						if !_rules[ruleSourceSinkParam]() {
							goto l998
						}
					l1006:
						{
							position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1007
							}
							if buffer[position] != rune(',') {
								goto l1007
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1007
							}
							if !_rules[ruleSourceSinkParam]() {
								goto l1007
							}
							goto l1006
						l1007:
							position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
						}
						goto l999
					l998:
						position, tokenIndex, depth = position998, tokenIndex998, depth998
					}
				l999:
					depth--
					add(rulePegText, position997)
				}
				if !_rules[ruleAction47]() {
					goto l995
				}
				depth--
				add(ruleSetOptSpecs, position996)
			}
			return true
		l995:
			position, tokenIndex, depth = position995, tokenIndex995, depth995
			return false
		},
		/* 63 StateTagOpt <- <(<(sp (('t' / 'T') ('a' / 'A') ('g' / 'G')) sp Identifier)?> Action48)> */
		func() bool {
			position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
			{
				position1009 := position
				depth++
				{
					position1010 := position
					depth++
					{
						position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1011
						}
						{
							position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1014
							}
							position++
							goto l1013
						l1014:
							position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
							if buffer[position] != rune('T') {
								goto l1011
							}
							position++
						}
					l1013:
						{
							position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1016
							}
							position++
							goto l1015
						l1016:
							position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
							if buffer[position] != rune('A') {
								goto l1011
							}
							position++
						}
					l1015:
						{
							position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
							if buffer[position] != rune('g') {
								goto l1018
							}
							position++
							goto l1017
						l1018:
							position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
							if buffer[position] != rune('G') {
								goto l1011
							}
							position++
						}
					l1017:
						if !_rules[rulesp]() {
							goto l1011
						}
						if !_rules[ruleIdentifier]() {
							goto l1011
						}
						goto l1012
					l1011:
						position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
					}
				l1012:
					depth--
					add(rulePegText, position1010)
				}
				if !_rules[ruleAction48]() {
					goto l1008
				}
				depth--
				add(ruleStateTagOpt, position1009)
			}
			return true
		l1008:
			position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
			return false
		},
		/* 64 SourceSinkParam <- <(SourceSinkParamKey spOpt '=' spOpt SourceSinkParamVal Action49)> */
		func() bool {
			position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
			{
				position1020 := position
				depth++
				if !_rules[ruleSourceSinkParamKey]() {
					goto l1019
				}
				if !_rules[rulespOpt]() {
					goto l1019
				}
				if buffer[position] != rune('=') {
					goto l1019
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1019
				}
				if !_rules[ruleSourceSinkParamVal]() {
					goto l1019
				}
				if !_rules[ruleAction49]() {
					goto l1019
				}
				depth--
				add(ruleSourceSinkParam, position1020)
			}
			return true
		l1019:
			position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
			return false
		},
		/* 65 SourceSinkParamVal <- <(ParamLiteral / ParamArrayExpr / ParamMapExpr)> */
		func() bool {
			position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
			{
				position1022 := position
				depth++
				{
					position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
					if !_rules[ruleParamLiteral]() {
						goto l1024
					}
					goto l1023
				l1024:
					position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
					if !_rules[ruleParamArrayExpr]() {
						goto l1025
					}
					goto l1023
				l1025:
					position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
					if !_rules[ruleParamMapExpr]() {
						goto l1021
					}
				}
			l1023:
				depth--
				add(ruleSourceSinkParamVal, position1022)
			}
			return true
		l1021:
			position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
			return false
		},
		/* 66 ParamLiteral <- <(BooleanLiteral / Literal)> */
		func() bool {
			position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
			{
				position1027 := position
				depth++
				{
					position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
					if !_rules[ruleBooleanLiteral]() {
						goto l1029
					}
					goto l1028
				l1029:
					position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
					if !_rules[ruleLiteral]() {
						goto l1026
					}
				}
			l1028:
				depth--
				add(ruleParamLiteral, position1027)
			}
			return true
		l1026:
			position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
			return false
		},
		/* 67 ParamArrayExpr <- <(<('[' spOpt (ParamLiteral (',' spOpt ParamLiteral)*)? spOpt ','? spOpt ']')> Action50)> */
		func() bool {
			position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
			{
				position1031 := position
				depth++
				{
					position1032 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1030
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1030
					}
					{
						position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
						if !_rules[ruleParamLiteral]() {
							goto l1033
						}
					l1035:
						{
							position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
							if buffer[position] != rune(',') {
								goto l1036
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1036
							}
							if !_rules[ruleParamLiteral]() {
								goto l1036
							}
							goto l1035
						l1036:
							position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
						}
						goto l1034
					l1033:
						position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
					}
				l1034:
					if !_rules[rulespOpt]() {
						goto l1030
					}
					{
						position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1037
						}
						position++
						goto l1038
					l1037:
						position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
					}
				l1038:
					if !_rules[rulespOpt]() {
						goto l1030
					}
					if buffer[position] != rune(']') {
						goto l1030
					}
					position++
					depth--
					add(rulePegText, position1032)
				}
				if !_rules[ruleAction50]() {
					goto l1030
				}
				depth--
				add(ruleParamArrayExpr, position1031)
			}
			return true
		l1030:
			position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
			return false
		},
		/* 68 ParamMapExpr <- <(<('{' spOpt (ParamKeyValuePair (spOpt ',' spOpt ParamKeyValuePair)*)? spOpt '}')> Action51)> */
		func() bool {
			position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
			{
				position1040 := position
				depth++
				{
					position1041 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1039
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1039
					}
					{
						position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
						if !_rules[ruleParamKeyValuePair]() {
							goto l1042
						}
					l1044:
						{
							position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1045
							}
							if buffer[position] != rune(',') {
								goto l1045
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1045
							}
							if !_rules[ruleParamKeyValuePair]() {
								goto l1045
							}
							goto l1044
						l1045:
							position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
						}
						goto l1043
					l1042:
						position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
					}
				l1043:
					if !_rules[rulespOpt]() {
						goto l1039
					}
					if buffer[position] != rune('}') {
						goto l1039
					}
					position++
					depth--
					add(rulePegText, position1041)
				}
				if !_rules[ruleAction51]() {
					goto l1039
				}
				depth--
				add(ruleParamMapExpr, position1040)
			}
			return true
		l1039:
			position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
			return false
		},
		/* 69 ParamKeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ParamLiteral)> Action52)> */
		func() bool {
			position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
			{
				position1047 := position
				depth++
				{
					position1048 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1046
					}
					if !_rules[rulespOpt]() {
						goto l1046
					}
					if buffer[position] != rune(':') {
						goto l1046
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1046
					}
					if !_rules[ruleParamLiteral]() {
						goto l1046
					}
					depth--
					add(rulePegText, position1048)
				}
				if !_rules[ruleAction52]() {
					goto l1046
				}
				depth--
				add(ruleParamKeyValuePair, position1047)
			}
			return true
		l1046:
			position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
			return false
		},
		/* 70 PausedOpt <- <(<(sp (Paused / Unpaused))?> Action53)> */
		func() bool {
			position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
			{
				position1050 := position
				depth++
				{
					position1051 := position
					depth++
					{
						position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1052
						}
						{
							position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
							if !_rules[rulePaused]() {
								goto l1055
							}
							goto l1054
						l1055:
							position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
							if !_rules[ruleUnpaused]() {
								goto l1052
							}
						}
					l1054:
						goto l1053
					l1052:
						position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
					}
				l1053:
					depth--
					add(rulePegText, position1051)
				}
				if !_rules[ruleAction53]() {
					goto l1049
				}
				depth--
				add(rulePausedOpt, position1050)
			}
			return true
		l1049:
			position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
			return false
		},
		/* 71 ExpressionOrWildcard <- <(Wildcard / Expression)> */
		func() bool {
			position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
			{
				position1057 := position
				depth++
				{
					position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
					if !_rules[ruleWildcard]() {
						goto l1059
					}
					goto l1058
				l1059:
					position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
					if !_rules[ruleExpression]() {
						goto l1056
					}
				}
			l1058:
				depth--
				add(ruleExpressionOrWildcard, position1057)
			}
			return true
		l1056:
			position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
			return false
		},
		/* 72 Expression <- <orExpr> */
		func() bool {
			position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
			{
				position1061 := position
				depth++
				if !_rules[ruleorExpr]() {
					goto l1060
				}
				depth--
				add(ruleExpression, position1061)
			}
			return true
		l1060:
			position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
			return false
		},
		/* 73 orExpr <- <(<(andExpr (sp Or sp andExpr)*)> Action54)> */
		func() bool {
			position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
			{
				position1063 := position
				depth++
				{
					position1064 := position
					depth++
					if !_rules[ruleandExpr]() {
						goto l1062
					}
				l1065:
					{
						position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1066
						}
						if !_rules[ruleOr]() {
							goto l1066
						}
						if !_rules[rulesp]() {
							goto l1066
						}
						if !_rules[ruleandExpr]() {
							goto l1066
						}
						goto l1065
					l1066:
						position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
					}
					depth--
					add(rulePegText, position1064)
				}
				if !_rules[ruleAction54]() {
					goto l1062
				}
				depth--
				add(ruleorExpr, position1063)
			}
			return true
		l1062:
			position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
			return false
		},
		/* 74 andExpr <- <(<(notExpr (sp And sp notExpr)*)> Action55)> */
		func() bool {
			position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
			{
				position1068 := position
				depth++
				{
					position1069 := position
					depth++
					if !_rules[rulenotExpr]() {
						goto l1067
					}
				l1070:
					{
						position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1071
						}
						if !_rules[ruleAnd]() {
							goto l1071
						}
						if !_rules[rulesp]() {
							goto l1071
						}
						if !_rules[rulenotExpr]() {
							goto l1071
						}
						goto l1070
					l1071:
						position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
					}
					depth--
					add(rulePegText, position1069)
				}
				if !_rules[ruleAction55]() {
					goto l1067
				}
				depth--
				add(ruleandExpr, position1068)
			}
			return true
		l1067:
			position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
			return false
		},
		/* 75 notExpr <- <(<((Not sp)? comparisonExpr)> Action56)> */
		func() bool {
			position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
			{
				position1073 := position
				depth++
				{
					position1074 := position
					depth++
					{
						position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
						if !_rules[ruleNot]() {
							goto l1075
						}
						if !_rules[rulesp]() {
							goto l1075
						}
						goto l1076
					l1075:
						position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
					}
				l1076:
					if !_rules[rulecomparisonExpr]() {
						goto l1072
					}
					depth--
					add(rulePegText, position1074)
				}
				if !_rules[ruleAction56]() {
					goto l1072
				}
				depth--
				add(rulenotExpr, position1073)
			}
			return true
		l1072:
			position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
			return false
		},
		/* 76 comparisonExpr <- <(<(otherOpExpr (spOpt ComparisonOp spOpt otherOpExpr)?)> Action57)> */
		func() bool {
			position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
			{
				position1078 := position
				depth++
				{
					position1079 := position
					depth++
					if !_rules[ruleotherOpExpr]() {
						goto l1077
					}
					{
						position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1080
						}
						if !_rules[ruleComparisonOp]() {
							goto l1080
						}
						if !_rules[rulespOpt]() {
							goto l1080
						}
						if !_rules[ruleotherOpExpr]() {
							goto l1080
						}
						goto l1081
					l1080:
						position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
					}
				l1081:
					depth--
					add(rulePegText, position1079)
				}
				if !_rules[ruleAction57]() {
					goto l1077
				}
				depth--
				add(rulecomparisonExpr, position1078)
			}
			return true
		l1077:
			position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
			return false
		},
		/* 77 otherOpExpr <- <(<(isExpr (spOpt OtherOp spOpt isExpr)*)> Action58)> */
		func() bool {
			position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
			{
				position1083 := position
				depth++
				{
					position1084 := position
					depth++
					if !_rules[ruleisExpr]() {
						goto l1082
					}
				l1085:
					{
						position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1086
						}
						if !_rules[ruleOtherOp]() {
							goto l1086
						}
						if !_rules[rulespOpt]() {
							goto l1086
						}
						if !_rules[ruleisExpr]() {
							goto l1086
						}
						goto l1085
					l1086:
						position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
					}
					depth--
					add(rulePegText, position1084)
				}
				if !_rules[ruleAction58]() {
					goto l1082
				}
				depth--
				add(ruleotherOpExpr, position1083)
			}
			return true
		l1082:
			position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
			return false
		},
		/* 78 isExpr <- <(<(termExpr (sp IsOp sp NullLiteral)?)> Action59)> */
		func() bool {
			position1087, tokenIndex1087, depth1087 := position, tokenIndex, depth
			{
				position1088 := position
				depth++
				{
					position1089 := position
					depth++
					if !_rules[ruletermExpr]() {
						goto l1087
					}
					{
						position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1090
						}
						if !_rules[ruleIsOp]() {
							goto l1090
						}
						if !_rules[rulesp]() {
							goto l1090
						}
						if !_rules[ruleNullLiteral]() {
							goto l1090
						}
						goto l1091
					l1090:
						position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
					}
				l1091:
					depth--
					add(rulePegText, position1089)
				}
				if !_rules[ruleAction59]() {
					goto l1087
				}
				depth--
				add(ruleisExpr, position1088)
			}
			return true
		l1087:
			position, tokenIndex, depth = position1087, tokenIndex1087, depth1087
			return false
		},
		/* 79 termExpr <- <(<(productExpr (spOpt PlusMinusOp spOpt productExpr)*)> Action60)> */
		func() bool {
			position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
			{
				position1093 := position
				depth++
				{
					position1094 := position
					depth++
					if !_rules[ruleproductExpr]() {
						goto l1092
					}
				l1095:
					{
						position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1096
						}
						if !_rules[rulePlusMinusOp]() {
							goto l1096
						}
						if !_rules[rulespOpt]() {
							goto l1096
						}
						if !_rules[ruleproductExpr]() {
							goto l1096
						}
						goto l1095
					l1096:
						position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
					}
					depth--
					add(rulePegText, position1094)
				}
				if !_rules[ruleAction60]() {
					goto l1092
				}
				depth--
				add(ruletermExpr, position1093)
			}
			return true
		l1092:
			position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
			return false
		},
		/* 80 productExpr <- <(<(minusExpr (spOpt MultDivOp spOpt minusExpr)*)> Action61)> */
		func() bool {
			position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
			{
				position1098 := position
				depth++
				{
					position1099 := position
					depth++
					if !_rules[ruleminusExpr]() {
						goto l1097
					}
				l1100:
					{
						position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1101
						}
						if !_rules[ruleMultDivOp]() {
							goto l1101
						}
						if !_rules[rulespOpt]() {
							goto l1101
						}
						if !_rules[ruleminusExpr]() {
							goto l1101
						}
						goto l1100
					l1101:
						position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
					}
					depth--
					add(rulePegText, position1099)
				}
				if !_rules[ruleAction61]() {
					goto l1097
				}
				depth--
				add(ruleproductExpr, position1098)
			}
			return true
		l1097:
			position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
			return false
		},
		/* 81 minusExpr <- <(<((UnaryMinus spOpt)? castExpr)> Action62)> */
		func() bool {
			position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
			{
				position1103 := position
				depth++
				{
					position1104 := position
					depth++
					{
						position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
						if !_rules[ruleUnaryMinus]() {
							goto l1105
						}
						if !_rules[rulespOpt]() {
							goto l1105
						}
						goto l1106
					l1105:
						position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
					}
				l1106:
					if !_rules[rulecastExpr]() {
						goto l1102
					}
					depth--
					add(rulePegText, position1104)
				}
				if !_rules[ruleAction62]() {
					goto l1102
				}
				depth--
				add(ruleminusExpr, position1103)
			}
			return true
		l1102:
			position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
			return false
		},
		/* 82 castExpr <- <(<(baseExpr (spOpt (':' ':') spOpt Type)?)> Action63)> */
		func() bool {
			position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
			{
				position1108 := position
				depth++
				{
					position1109 := position
					depth++
					if !_rules[rulebaseExpr]() {
						goto l1107
					}
					{
						position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1110
						}
						if buffer[position] != rune(':') {
							goto l1110
						}
						position++
						if buffer[position] != rune(':') {
							goto l1110
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1110
						}
						if !_rules[ruleType]() {
							goto l1110
						}
						goto l1111
					l1110:
						position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
					}
				l1111:
					depth--
					add(rulePegText, position1109)
				}
				if !_rules[ruleAction63]() {
					goto l1107
				}
				depth--
				add(rulecastExpr, position1108)
			}
			return true
		l1107:
			position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
			return false
		},
		/* 83 baseExpr <- <(('(' spOpt Expression spOpt ')') / MapExpr / BooleanLiteral / NullLiteral / Case / RowMeta / FuncTypeCast / FuncApp / RowValue / ArrayExpr / Literal)> */
		func() bool {
			position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
			{
				position1113 := position
				depth++
				{
					position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l1115
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1115
					}
					if !_rules[ruleExpression]() {
						goto l1115
					}
					if !_rules[rulespOpt]() {
						goto l1115
					}
					if buffer[position] != rune(')') {
						goto l1115
					}
					position++
					goto l1114
				l1115:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleMapExpr]() {
						goto l1116
					}
					goto l1114
				l1116:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleBooleanLiteral]() {
						goto l1117
					}
					goto l1114
				l1117:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleNullLiteral]() {
						goto l1118
					}
					goto l1114
				l1118:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleCase]() {
						goto l1119
					}
					goto l1114
				l1119:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleRowMeta]() {
						goto l1120
					}
					goto l1114
				l1120:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleFuncTypeCast]() {
						goto l1121
					}
					goto l1114
				l1121:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleFuncApp]() {
						goto l1122
					}
					goto l1114
				l1122:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleRowValue]() {
						goto l1123
					}
					goto l1114
				l1123:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleArrayExpr]() {
						goto l1124
					}
					goto l1114
				l1124:
					position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
					if !_rules[ruleLiteral]() {
						goto l1112
					}
				}
			l1114:
				depth--
				add(rulebaseExpr, position1113)
			}
			return true
		l1112:
			position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
			return false
		},
		/* 84 FuncTypeCast <- <(<(('c' / 'C') ('a' / 'A') ('s' / 'S') ('t' / 'T') spOpt '(' spOpt Expression sp (('a' / 'A') ('s' / 'S')) sp Type spOpt ')')> Action64)> */
		func() bool {
			position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
			{
				position1126 := position
				depth++
				{
					position1127 := position
					depth++
					{
						position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1129
						}
						position++
						goto l1128
					l1129:
						position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
						if buffer[position] != rune('C') {
							goto l1125
						}
						position++
					}
				l1128:
					{
						position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1131
						}
						position++
						goto l1130
					l1131:
						position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
						if buffer[position] != rune('A') {
							goto l1125
						}
						position++
					}
				l1130:
					{
						position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1133
						}
						position++
						goto l1132
					l1133:
						position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
						if buffer[position] != rune('S') {
							goto l1125
						}
						position++
					}
				l1132:
					{
						position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1135
						}
						position++
						goto l1134
					l1135:
						position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
						if buffer[position] != rune('T') {
							goto l1125
						}
						position++
					}
				l1134:
					if !_rules[rulespOpt]() {
						goto l1125
					}
					if buffer[position] != rune('(') {
						goto l1125
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1125
					}
					if !_rules[ruleExpression]() {
						goto l1125
					}
					if !_rules[rulesp]() {
						goto l1125
					}
					{
						position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1137
						}
						position++
						goto l1136
					l1137:
						position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
						if buffer[position] != rune('A') {
							goto l1125
						}
						position++
					}
				l1136:
					{
						position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1139
						}
						position++
						goto l1138
					l1139:
						position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
						if buffer[position] != rune('S') {
							goto l1125
						}
						position++
					}
				l1138:
					if !_rules[rulesp]() {
						goto l1125
					}
					if !_rules[ruleType]() {
						goto l1125
					}
					if !_rules[rulespOpt]() {
						goto l1125
					}
					if buffer[position] != rune(')') {
						goto l1125
					}
					position++
					depth--
					add(rulePegText, position1127)
				}
				if !_rules[ruleAction64]() {
					goto l1125
				}
				depth--
				add(ruleFuncTypeCast, position1126)
			}
			return true
		l1125:
			position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
			return false
		},
		/* 85 FuncApp <- <(FuncAppWithOrderBy / FuncAppWithoutOrderBy)> */
		func() bool {
			position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
			{
				position1141 := position
				depth++
				{
					position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
					if !_rules[ruleFuncAppWithOrderBy]() {
						goto l1143
					}
					goto l1142
				l1143:
					position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
					if !_rules[ruleFuncAppWithoutOrderBy]() {
						goto l1140
					}
				}
			l1142:
				depth--
				add(ruleFuncApp, position1141)
			}
			return true
		l1140:
			position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
			return false
		},
		/* 86 FuncAppWithOrderBy <- <(Function spOpt '(' spOpt FuncParams sp ParamsOrder spOpt ')' Action65)> */
		func() bool {
			position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
			{
				position1145 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1144
				}
				if !_rules[rulespOpt]() {
					goto l1144
				}
				if buffer[position] != rune('(') {
					goto l1144
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1144
				}
				if !_rules[ruleFuncParams]() {
					goto l1144
				}
				if !_rules[rulesp]() {
					goto l1144
				}
				if !_rules[ruleParamsOrder]() {
					goto l1144
				}
				if !_rules[rulespOpt]() {
					goto l1144
				}
				if buffer[position] != rune(')') {
					goto l1144
				}
				position++
				if !_rules[ruleAction65]() {
					goto l1144
				}
				depth--
				add(ruleFuncAppWithOrderBy, position1145)
			}
			return true
		l1144:
			position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
			return false
		},
		/* 87 FuncAppWithoutOrderBy <- <(Function spOpt '(' spOpt FuncParams <spOpt> ')' Action66)> */
		func() bool {
			position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
			{
				position1147 := position
				depth++
				if !_rules[ruleFunction]() {
					goto l1146
				}
				if !_rules[rulespOpt]() {
					goto l1146
				}
				if buffer[position] != rune('(') {
					goto l1146
				}
				position++
				if !_rules[rulespOpt]() {
					goto l1146
				}
				if !_rules[ruleFuncParams]() {
					goto l1146
				}
				{
					position1148 := position
					depth++
					if !_rules[rulespOpt]() {
						goto l1146
					}
					depth--
					add(rulePegText, position1148)
				}
				if buffer[position] != rune(')') {
					goto l1146
				}
				position++
				if !_rules[ruleAction66]() {
					goto l1146
				}
				depth--
				add(ruleFuncAppWithoutOrderBy, position1147)
			}
			return true
		l1146:
			position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
			return false
		},
		/* 88 FuncParams <- <(<(ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)?> Action67)> */
		func() bool {
			position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
			{
				position1150 := position
				depth++
				{
					position1151 := position
					depth++
					{
						position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1152
						}
					l1154:
						{
							position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1155
							}
							if buffer[position] != rune(',') {
								goto l1155
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1155
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1155
							}
							goto l1154
						l1155:
							position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
						}
						goto l1153
					l1152:
						position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
					}
				l1153:
					depth--
					add(rulePegText, position1151)
				}
				if !_rules[ruleAction67]() {
					goto l1149
				}
				depth--
				add(ruleFuncParams, position1150)
			}
			return true
		l1149:
			position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
			return false
		},
		/* 89 ParamsOrder <- <(<(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') sp (('b' / 'B') ('y' / 'Y')) sp SortedExpression (spOpt ',' spOpt SortedExpression)*)> Action68)> */
		func() bool {
			position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
			{
				position1157 := position
				depth++
				{
					position1158 := position
					depth++
					{
						position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1160
						}
						position++
						goto l1159
					l1160:
						position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
						if buffer[position] != rune('O') {
							goto l1156
						}
						position++
					}
				l1159:
					{
						position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1162
						}
						position++
						goto l1161
					l1162:
						position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
						if buffer[position] != rune('R') {
							goto l1156
						}
						position++
					}
				l1161:
					{
						position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1164
						}
						position++
						goto l1163
					l1164:
						position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
						if buffer[position] != rune('D') {
							goto l1156
						}
						position++
					}
				l1163:
					{
						position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1166
						}
						position++
						goto l1165
					l1166:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						if buffer[position] != rune('E') {
							goto l1156
						}
						position++
					}
				l1165:
					{
						position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1168
						}
						position++
						goto l1167
					l1168:
						position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
						if buffer[position] != rune('R') {
							goto l1156
						}
						position++
					}
				l1167:
					if !_rules[rulesp]() {
						goto l1156
					}
					{
						position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1170
						}
						position++
						goto l1169
					l1170:
						position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
						if buffer[position] != rune('B') {
							goto l1156
						}
						position++
					}
				l1169:
					{
						position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1172
						}
						position++
						goto l1171
					l1172:
						position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
						if buffer[position] != rune('Y') {
							goto l1156
						}
						position++
					}
				l1171:
					if !_rules[rulesp]() {
						goto l1156
					}
					if !_rules[ruleSortedExpression]() {
						goto l1156
					}
				l1173:
					{
						position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
						if !_rules[rulespOpt]() {
							goto l1174
						}
						if buffer[position] != rune(',') {
							goto l1174
						}
						position++
						if !_rules[rulespOpt]() {
							goto l1174
						}
						if !_rules[ruleSortedExpression]() {
							goto l1174
						}
						goto l1173
					l1174:
						position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
					}
					depth--
					add(rulePegText, position1158)
				}
				if !_rules[ruleAction68]() {
					goto l1156
				}
				depth--
				add(ruleParamsOrder, position1157)
			}
			return true
		l1156:
			position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
			return false
		},
		/* 90 SortedExpression <- <(Expression OrderDirectionOpt Action69)> */
		func() bool {
			position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
			{
				position1176 := position
				depth++
				if !_rules[ruleExpression]() {
					goto l1175
				}
				if !_rules[ruleOrderDirectionOpt]() {
					goto l1175
				}
				if !_rules[ruleAction69]() {
					goto l1175
				}
				depth--
				add(ruleSortedExpression, position1176)
			}
			return true
		l1175:
			position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
			return false
		},
		/* 91 OrderDirectionOpt <- <(<(sp (Ascending / Descending))?> Action70)> */
		func() bool {
			position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
			{
				position1178 := position
				depth++
				{
					position1179 := position
					depth++
					{
						position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1180
						}
						{
							position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
							if !_rules[ruleAscending]() {
								goto l1183
							}
							goto l1182
						l1183:
							position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
							if !_rules[ruleDescending]() {
								goto l1180
							}
						}
					l1182:
						goto l1181
					l1180:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
					}
				l1181:
					depth--
					add(rulePegText, position1179)
				}
				if !_rules[ruleAction70]() {
					goto l1177
				}
				depth--
				add(ruleOrderDirectionOpt, position1178)
			}
			return true
		l1177:
			position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
			return false
		},
		/* 92 ArrayExpr <- <(<('[' spOpt (ExpressionOrWildcard (spOpt ',' spOpt ExpressionOrWildcard)*)? spOpt ','? spOpt ']')> Action71)> */
		func() bool {
			position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
			{
				position1185 := position
				depth++
				{
					position1186 := position
					depth++
					if buffer[position] != rune('[') {
						goto l1184
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1184
					}
					{
						position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
						if !_rules[ruleExpressionOrWildcard]() {
							goto l1187
						}
					l1189:
						{
							position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1190
							}
							if buffer[position] != rune(',') {
								goto l1190
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1190
							}
							if !_rules[ruleExpressionOrWildcard]() {
								goto l1190
							}
							goto l1189
						l1190:
							position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
						}
						goto l1188
					l1187:
						position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
					}
				l1188:
					if !_rules[rulespOpt]() {
						goto l1184
					}
					{
						position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
						if buffer[position] != rune(',') {
							goto l1191
						}
						position++
						goto l1192
					l1191:
						position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
					}
				l1192:
					if !_rules[rulespOpt]() {
						goto l1184
					}
					if buffer[position] != rune(']') {
						goto l1184
					}
					position++
					depth--
					add(rulePegText, position1186)
				}
				if !_rules[ruleAction71]() {
					goto l1184
				}
				depth--
				add(ruleArrayExpr, position1185)
			}
			return true
		l1184:
			position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
			return false
		},
		/* 93 MapExpr <- <(<('{' spOpt (KeyValuePair (spOpt ',' spOpt KeyValuePair)*)? spOpt '}')> Action72)> */
		func() bool {
			position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
			{
				position1194 := position
				depth++
				{
					position1195 := position
					depth++
					if buffer[position] != rune('{') {
						goto l1193
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1193
					}
					{
						position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
						if !_rules[ruleKeyValuePair]() {
							goto l1196
						}
					l1198:
						{
							position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
							if !_rules[rulespOpt]() {
								goto l1199
							}
							if buffer[position] != rune(',') {
								goto l1199
							}
							position++
							if !_rules[rulespOpt]() {
								goto l1199
							}
							if !_rules[ruleKeyValuePair]() {
								goto l1199
							}
							goto l1198
						l1199:
							position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
						}
						goto l1197
					l1196:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
					}
				l1197:
					if !_rules[rulespOpt]() {
						goto l1193
					}
					if buffer[position] != rune('}') {
						goto l1193
					}
					position++
					depth--
					add(rulePegText, position1195)
				}
				if !_rules[ruleAction72]() {
					goto l1193
				}
				depth--
				add(ruleMapExpr, position1194)
			}
			return true
		l1193:
			position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
			return false
		},
		/* 94 KeyValuePair <- <(<(StringLiteral spOpt ':' spOpt ExpressionOrWildcard)> Action73)> */
		func() bool {
			position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
			{
				position1201 := position
				depth++
				{
					position1202 := position
					depth++
					if !_rules[ruleStringLiteral]() {
						goto l1200
					}
					if !_rules[rulespOpt]() {
						goto l1200
					}
					if buffer[position] != rune(':') {
						goto l1200
					}
					position++
					if !_rules[rulespOpt]() {
						goto l1200
					}
					if !_rules[ruleExpressionOrWildcard]() {
						goto l1200
					}
					depth--
					add(rulePegText, position1202)
				}
				if !_rules[ruleAction73]() {
					goto l1200
				}
				depth--
				add(ruleKeyValuePair, position1201)
			}
			return true
		l1200:
			position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
			return false
		},
		/* 95 Case <- <(ConditionCase / ExpressionCase)> */
		func() bool {
			position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
			{
				position1204 := position
				depth++
				{
					position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
					if !_rules[ruleConditionCase]() {
						goto l1206
					}
					goto l1205
				l1206:
					position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
					if !_rules[ruleExpressionCase]() {
						goto l1203
					}
				}
			l1205:
				depth--
				add(ruleCase, position1204)
			}
			return true
		l1203:
			position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
			return false
		},
		/* 96 ConditionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action74)> */
		func() bool {
			position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
			{
				position1208 := position
				depth++
				{
					position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1210
					}
					position++
					goto l1209
				l1210:
					position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
					if buffer[position] != rune('C') {
						goto l1207
					}
					position++
				}
			l1209:
				{
					position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1212
					}
					position++
					goto l1211
				l1212:
					position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
					if buffer[position] != rune('A') {
						goto l1207
					}
					position++
				}
			l1211:
				{
					position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1214
					}
					position++
					goto l1213
				l1214:
					position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
					if buffer[position] != rune('S') {
						goto l1207
					}
					position++
				}
			l1213:
				{
					position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1216
					}
					position++
					goto l1215
				l1216:
					position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
					if buffer[position] != rune('E') {
						goto l1207
					}
					position++
				}
			l1215:
				{
					position1217 := position
					depth++
					if !_rules[rulesp]() {
						goto l1207
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1207
					}
				l1218:
					{
						position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1219
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1219
						}
						goto l1218
					l1219:
						position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
					}
					{
						position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1220
						}
						{
							position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1223
							}
							position++
							goto l1222
						l1223:
							position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
							if buffer[position] != rune('E') {
								goto l1220
							}
							position++
						}
					l1222:
						{
							position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1225
							}
							position++
							goto l1224
						l1225:
							position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
							if buffer[position] != rune('L') {
								goto l1220
							}
							position++
						}
					l1224:
						{
							position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1227
							}
							position++
							goto l1226
						l1227:
							position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
							if buffer[position] != rune('S') {
								goto l1220
							}
							position++
						}
					l1226:
						{
							position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1229
							}
							position++
							goto l1228
						l1229:
							position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
							if buffer[position] != rune('E') {
								goto l1220
							}
							position++
						}
					l1228:
						if !_rules[rulesp]() {
							goto l1220
						}
						if !_rules[ruleExpression]() {
							goto l1220
						}
						goto l1221
					l1220:
						position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
					}
				l1221:
					if !_rules[rulesp]() {
						goto l1207
					}
					{
						position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1231
						}
						position++
						goto l1230
					l1231:
						position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
						if buffer[position] != rune('E') {
							goto l1207
						}
						position++
					}
				l1230:
					{
						position1232, tokenIndex1232, depth1232 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1233
						}
						position++
						goto l1232
					l1233:
						position, tokenIndex, depth = position1232, tokenIndex1232, depth1232
						if buffer[position] != rune('N') {
							goto l1207
						}
						position++
					}
				l1232:
					{
						position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1235
						}
						position++
						goto l1234
					l1235:
						position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
						if buffer[position] != rune('D') {
							goto l1207
						}
						position++
					}
				l1234:
					depth--
					add(rulePegText, position1217)
				}
				if !_rules[ruleAction74]() {
					goto l1207
				}
				depth--
				add(ruleConditionCase, position1208)
			}
			return true
		l1207:
			position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
			return false
		},
		/* 97 ExpressionCase <- <(('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') sp Expression <((sp WhenThenPair)+ (sp (('e' / 'E') ('l' / 'L') ('s' / 'S') ('e' / 'E')) sp Expression)? sp (('e' / 'E') ('n' / 'N') ('d' / 'D')))> Action75)> */
		func() bool {
			position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
			{
				position1237 := position
				depth++
				{
					position1238, tokenIndex1238, depth1238 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1239
					}
					position++
					goto l1238
				l1239:
					position, tokenIndex, depth = position1238, tokenIndex1238, depth1238
					if buffer[position] != rune('C') {
						goto l1236
					}
					position++
				}
			l1238:
				{
					position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1241
					}
					position++
					goto l1240
				l1241:
					position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
					if buffer[position] != rune('A') {
						goto l1236
					}
					position++
				}
			l1240:
				{
					position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1243
					}
					position++
					goto l1242
				l1243:
					position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
					if buffer[position] != rune('S') {
						goto l1236
					}
					position++
				}
			l1242:
				{
					position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1245
					}
					position++
					goto l1244
				l1245:
					position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
					if buffer[position] != rune('E') {
						goto l1236
					}
					position++
				}
			l1244:
				if !_rules[rulesp]() {
					goto l1236
				}
				if !_rules[ruleExpression]() {
					goto l1236
				}
				{
					position1246 := position
					depth++
					if !_rules[rulesp]() {
						goto l1236
					}
					if !_rules[ruleWhenThenPair]() {
						goto l1236
					}
				l1247:
					{
						position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1248
						}
						if !_rules[ruleWhenThenPair]() {
							goto l1248
						}
						goto l1247
					l1248:
						position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
					}
					{
						position1249, tokenIndex1249, depth1249 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l1249
						}
						{
							position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1252
							}
							position++
							goto l1251
						l1252:
							position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
							if buffer[position] != rune('E') {
								goto l1249
							}
							position++
						}
					l1251:
						{
							position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1254
							}
							position++
							goto l1253
						l1254:
							position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
							if buffer[position] != rune('L') {
								goto l1249
							}
							position++
						}
					l1253:
						{
							position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1256
							}
							position++
							goto l1255
						l1256:
							position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
							if buffer[position] != rune('S') {
								goto l1249
							}
							position++
						}
					l1255:
						{
							position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1258
							}
							position++
							goto l1257
						l1258:
							position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
							if buffer[position] != rune('E') {
								goto l1249
							}
							position++
						}
					l1257:
						if !_rules[rulesp]() {
							goto l1249
						}
						if !_rules[ruleExpression]() {
							goto l1249
						}
						goto l1250
					l1249:
						position, tokenIndex, depth = position1249, tokenIndex1249, depth1249
					}
				l1250:
					if !_rules[rulesp]() {
						goto l1236
					}
					{
						position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1260
						}
						position++
						goto l1259
					l1260:
						position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
						if buffer[position] != rune('E') {
							goto l1236
						}
						position++
					}
				l1259:
					{
						position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1262
						}
						position++
						goto l1261
					l1262:
						position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
						if buffer[position] != rune('N') {
							goto l1236
						}
						position++
					}
				l1261:
					{
						position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1264
						}
						position++
						goto l1263
					l1264:
						position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
						if buffer[position] != rune('D') {
							goto l1236
						}
						position++
					}
				l1263:
					depth--
					add(rulePegText, position1246)
				}
				if !_rules[ruleAction75]() {
					goto l1236
				}
				depth--
				add(ruleExpressionCase, position1237)
			}
			return true
		l1236:
			position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
			return false
		},
		/* 98 WhenThenPair <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('n' / 'N') sp Expression sp (('t' / 'T') ('h' / 'H') ('e' / 'E') ('n' / 'N')) sp ExpressionOrWildcard Action76)> */
		func() bool {
			position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
			{
				position1266 := position
				depth++
				{
					position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
					if buffer[position] != rune('w') {
						goto l1268
					}
					position++
					goto l1267
				l1268:
					position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
					if buffer[position] != rune('W') {
						goto l1265
					}
					position++
				}
			l1267:
				{
					position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
					if buffer[position] != rune('h') {
						goto l1270
					}
					position++
					goto l1269
				l1270:
					position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
					if buffer[position] != rune('H') {
						goto l1265
					}
					position++
				}
			l1269:
				{
					position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1272
					}
					position++
					goto l1271
				l1272:
					position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
					if buffer[position] != rune('E') {
						goto l1265
					}
					position++
				}
			l1271:
				{
					position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1274
					}
					position++
					goto l1273
				l1274:
					position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
					if buffer[position] != rune('N') {
						goto l1265
					}
					position++
				}
			l1273:
				if !_rules[rulesp]() {
					goto l1265
				}
				if !_rules[ruleExpression]() {
					goto l1265
				}
				if !_rules[rulesp]() {
					goto l1265
				}
				{
					position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1276
					}
					position++
					goto l1275
				l1276:
					position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
					if buffer[position] != rune('T') {
						goto l1265
					}
					position++
				}
			l1275:
				{
					position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
					if buffer[position] != rune('h') {
						goto l1278
					}
					position++
					goto l1277
				l1278:
					position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
					if buffer[position] != rune('H') {
						goto l1265
					}
					position++
				}
			l1277:
				{
					position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l1280
					}
					position++
					goto l1279
				l1280:
					position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
					if buffer[position] != rune('E') {
						goto l1265
					}
					position++
				}
			l1279:
				{
					position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1282
					}
					position++
					goto l1281
				l1282:
					position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
					if buffer[position] != rune('N') {
						goto l1265
					}
					position++
				}
			l1281:
				if !_rules[rulesp]() {
					goto l1265
				}
				if !_rules[ruleExpressionOrWildcard]() {
					goto l1265
				}
				if !_rules[ruleAction76]() {
					goto l1265
				}
				depth--
				add(ruleWhenThenPair, position1266)
			}
			return true
		l1265:
			position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
			return false
		},
		/* 99 Literal <- <(FloatLiteral / NumericLiteral / StringLiteral)> */
		func() bool {
			position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
			{
				position1284 := position
				depth++
				{
					position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l1286
					}
					goto l1285
				l1286:
					position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
					if !_rules[ruleNumericLiteral]() {
						goto l1287
					}
					goto l1285
				l1287:
					position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
					if !_rules[ruleStringLiteral]() {
						goto l1283
					}
				}
			l1285:
				depth--
				add(ruleLiteral, position1284)
			}
			return true
		l1283:
			position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
			return false
		},
		/* 100 ComparisonOp <- <(Equal / NotEqual / LessOrEqual / Less / GreaterOrEqual / Greater / NotEqual)> */
		func() bool {
			position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
			{
				position1289 := position
				depth++
				{
					position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
					if !_rules[ruleEqual]() {
						goto l1291
					}
					goto l1290
				l1291:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if !_rules[ruleNotEqual]() {
						goto l1292
					}
					goto l1290
				l1292:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if !_rules[ruleLessOrEqual]() {
						goto l1293
					}
					goto l1290
				l1293:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if !_rules[ruleLess]() {
						goto l1294
					}
					goto l1290
				l1294:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if !_rules[ruleGreaterOrEqual]() {
						goto l1295
					}
					goto l1290
				l1295:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if !_rules[ruleGreater]() {
						goto l1296
					}
					goto l1290
				l1296:
					position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
					if !_rules[ruleNotEqual]() {
						goto l1288
					}
				}
			l1290:
				depth--
				add(ruleComparisonOp, position1289)
			}
			return true
		l1288:
			position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
			return false
		},
		/* 101 OtherOp <- <Concat> */
		func() bool {
			position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
			{
				position1298 := position
				depth++
				if !_rules[ruleConcat]() {
					goto l1297
				}
				depth--
				add(ruleOtherOp, position1298)
			}
			return true
		l1297:
			position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
			return false
		},
		/* 102 IsOp <- <(IsNot / Is)> */
		func() bool {
			position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
			{
				position1300 := position
				depth++
				{
					position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
					if !_rules[ruleIsNot]() {
						goto l1302
					}
					goto l1301
				l1302:
					position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
					if !_rules[ruleIs]() {
						goto l1299
					}
				}
			l1301:
				depth--
				add(ruleIsOp, position1300)
			}
			return true
		l1299:
			position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
			return false
		},
		/* 103 PlusMinusOp <- <(Plus / Minus)> */
		func() bool {
			position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
			{
				position1304 := position
				depth++
				{
					position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
					if !_rules[rulePlus]() {
						goto l1306
					}
					goto l1305
				l1306:
					position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
					if !_rules[ruleMinus]() {
						goto l1303
					}
				}
			l1305:
				depth--
				add(rulePlusMinusOp, position1304)
			}
			return true
		l1303:
			position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
			return false
		},
		/* 104 MultDivOp <- <(Multiply / Divide / Modulo)> */
		func() bool {
			position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
			{
				position1308 := position
				depth++
				{
					position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
					if !_rules[ruleMultiply]() {
						goto l1310
					}
					goto l1309
				l1310:
					position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
					if !_rules[ruleDivide]() {
						goto l1311
					}
					goto l1309
				l1311:
					position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
					if !_rules[ruleModulo]() {
						goto l1307
					}
				}
			l1309:
				depth--
				add(ruleMultDivOp, position1308)
			}
			return true
		l1307:
			position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
			return false
		},
		/* 105 Stream <- <(<ident> Action77)> */
		func() bool {
			position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
			{
				position1313 := position
				depth++
				{
					position1314 := position
					depth++
					if !_rules[ruleident]() {
						goto l1312
					}
					depth--
					add(rulePegText, position1314)
				}
				if !_rules[ruleAction77]() {
					goto l1312
				}
				depth--
				add(ruleStream, position1313)
			}
			return true
		l1312:
			position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
			return false
		},
		/* 106 RowMeta <- <RowTimestamp> */
		func() bool {
			position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
			{
				position1316 := position
				depth++
				if !_rules[ruleRowTimestamp]() {
					goto l1315
				}
				depth--
				add(ruleRowMeta, position1316)
			}
			return true
		l1315:
			position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
			return false
		},
		/* 107 RowTimestamp <- <(<((ident ':')? ('t' 's' '(' ')'))> Action78)> */
		func() bool {
			position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
			{
				position1318 := position
				depth++
				{
					position1319 := position
					depth++
					{
						position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1320
						}
						if buffer[position] != rune(':') {
							goto l1320
						}
						position++
						goto l1321
					l1320:
						position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
					}
				l1321:
					if buffer[position] != rune('t') {
						goto l1317
					}
					position++
					if buffer[position] != rune('s') {
						goto l1317
					}
					position++
					if buffer[position] != rune('(') {
						goto l1317
					}
					position++
					if buffer[position] != rune(')') {
						goto l1317
					}
					position++
					depth--
					add(rulePegText, position1319)
				}
				if !_rules[ruleAction78]() {
					goto l1317
				}
				depth--
				add(ruleRowTimestamp, position1318)
			}
			return true
		l1317:
			position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
			return false
		},
		/* 108 RowValue <- <(<((ident ':' !':')? jsonGetPath)> Action79)> */
		func() bool {
			position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
			{
				position1323 := position
				depth++
				{
					position1324 := position
					depth++
					{
						position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1325
						}
						if buffer[position] != rune(':') {
							goto l1325
						}
						position++
						{
							position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1327
							}
							position++
							goto l1325
						l1327:
							position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
						}
						goto l1326
					l1325:
						position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
					}
				l1326:
					if !_rules[rulejsonGetPath]() {
						goto l1322
					}
					depth--
					add(rulePegText, position1324)
				}
				if !_rules[ruleAction79]() {
					goto l1322
				}
				depth--
				add(ruleRowValue, position1323)
			}
			return true
		l1322:
			position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
			return false
		},
		/* 109 NumericLiteral <- <(<('-'? [0-9]+)> Action80)> */
		func() bool {
			position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
			{
				position1329 := position
				depth++
				{
					position1330 := position
					depth++
					{
						position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1331
						}
						position++
						goto l1332
					l1331:
						position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
					}
				l1332:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1328
					}
					position++
				l1333:
					{
						position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1334
						}
						position++
						goto l1333
					l1334:
						position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
					}
					depth--
					add(rulePegText, position1330)
				}
				if !_rules[ruleAction80]() {
					goto l1328
				}
				depth--
				add(ruleNumericLiteral, position1329)
			}
			return true
		l1328:
			position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
			return false
		},
		/* 110 NonNegativeNumericLiteral <- <(<[0-9]+> Action81)> */
		func() bool {
			position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
			{
				position1336 := position
				depth++
				{
					position1337 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1335
					}
					position++
				l1338:
					{
						position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1339
						}
						position++
						goto l1338
					l1339:
						position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
					}
					depth--
					add(rulePegText, position1337)
				}
				if !_rules[ruleAction81]() {
					goto l1335
				}
				depth--
				add(ruleNonNegativeNumericLiteral, position1336)
			}
			return true
		l1335:
			position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
			return false
		},
		/* 111 FloatLiteral <- <(<('-'? [0-9]+ '.' [0-9]+)> Action82)> */
		func() bool {
			position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
			{
				position1341 := position
				depth++
				{
					position1342 := position
					depth++
					{
						position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1343
						}
						position++
						goto l1344
					l1343:
						position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
					}
				l1344:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1340
					}
					position++
				l1345:
					{
						position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1346
						}
						position++
						goto l1345
					l1346:
						position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
					}
					if buffer[position] != rune('.') {
						goto l1340
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1340
					}
					position++
				l1347:
					{
						position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1348
						}
						position++
						goto l1347
					l1348:
						position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
					}
					depth--
					add(rulePegText, position1342)
				}
				if !_rules[ruleAction82]() {
					goto l1340
				}
				depth--
				add(ruleFloatLiteral, position1341)
			}
			return true
		l1340:
			position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
			return false
		},
		/* 112 Function <- <(<ident> Action83)> */
		func() bool {
			position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
			{
				position1350 := position
				depth++
				{
					position1351 := position
					depth++
					if !_rules[ruleident]() {
						goto l1349
					}
					depth--
					add(rulePegText, position1351)
				}
				if !_rules[ruleAction83]() {
					goto l1349
				}
				depth--
				add(ruleFunction, position1350)
			}
			return true
		l1349:
			position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
			return false
		},
		/* 113 NullLiteral <- <(<(('n' / 'N') ('u' / 'U') ('l' / 'L') ('l' / 'L'))> Action84)> */
		func() bool {
			position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
			{
				position1353 := position
				depth++
				{
					position1354 := position
					depth++
					{
						position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1356
						}
						position++
						goto l1355
					l1356:
						position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
						if buffer[position] != rune('N') {
							goto l1352
						}
						position++
					}
				l1355:
					{
						position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1358
						}
						position++
						goto l1357
					l1358:
						position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
						if buffer[position] != rune('U') {
							goto l1352
						}
						position++
					}
				l1357:
					{
						position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1360
						}
						position++
						goto l1359
					l1360:
						position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
						if buffer[position] != rune('L') {
							goto l1352
						}
						position++
					}
				l1359:
					{
						position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1362
						}
						position++
						goto l1361
					l1362:
						position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
						if buffer[position] != rune('L') {
							goto l1352
						}
						position++
					}
				l1361:
					depth--
					add(rulePegText, position1354)
				}
				if !_rules[ruleAction84]() {
					goto l1352
				}
				depth--
				add(ruleNullLiteral, position1353)
			}
			return true
		l1352:
			position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
			return false
		},
		/* 114 BooleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
			{
				position1364 := position
				depth++
				{
					position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
					if !_rules[ruleTRUE]() {
						goto l1366
					}
					goto l1365
				l1366:
					position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
					if !_rules[ruleFALSE]() {
						goto l1363
					}
				}
			l1365:
				depth--
				add(ruleBooleanLiteral, position1364)
			}
			return true
		l1363:
			position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
			return false
		},
		/* 115 TRUE <- <(<(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E'))> Action85)> */
		func() bool {
			position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
			{
				position1368 := position
				depth++
				{
					position1369 := position
					depth++
					{
						position1370, tokenIndex1370, depth1370 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1371
						}
						position++
						goto l1370
					l1371:
						position, tokenIndex, depth = position1370, tokenIndex1370, depth1370
						if buffer[position] != rune('T') {
							goto l1367
						}
						position++
					}
				l1370:
					{
						position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1373
						}
						position++
						goto l1372
					l1373:
						position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
						if buffer[position] != rune('R') {
							goto l1367
						}
						position++
					}
				l1372:
					{
						position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1375
						}
						position++
						goto l1374
					l1375:
						position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
						if buffer[position] != rune('U') {
							goto l1367
						}
						position++
					}
				l1374:
					{
						position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1377
						}
						position++
						goto l1376
					l1377:
						position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
						if buffer[position] != rune('E') {
							goto l1367
						}
						position++
					}
				l1376:
					depth--
					add(rulePegText, position1369)
				}
				if !_rules[ruleAction85]() {
					goto l1367
				}
				depth--
				add(ruleTRUE, position1368)
			}
			return true
		l1367:
			position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
			return false
		},
		/* 116 FALSE <- <(<(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E'))> Action86)> */
		func() bool {
			position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
			{
				position1379 := position
				depth++
				{
					position1380 := position
					depth++
					{
						position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1382
						}
						position++
						goto l1381
					l1382:
						position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
						if buffer[position] != rune('F') {
							goto l1378
						}
						position++
					}
				l1381:
					{
						position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1384
						}
						position++
						goto l1383
					l1384:
						position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
						if buffer[position] != rune('A') {
							goto l1378
						}
						position++
					}
				l1383:
					{
						position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1386
						}
						position++
						goto l1385
					l1386:
						position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
						if buffer[position] != rune('L') {
							goto l1378
						}
						position++
					}
				l1385:
					{
						position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1388
						}
						position++
						goto l1387
					l1388:
						position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
						if buffer[position] != rune('S') {
							goto l1378
						}
						position++
					}
				l1387:
					{
						position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1390
						}
						position++
						goto l1389
					l1390:
						position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
						if buffer[position] != rune('E') {
							goto l1378
						}
						position++
					}
				l1389:
					depth--
					add(rulePegText, position1380)
				}
				if !_rules[ruleAction86]() {
					goto l1378
				}
				depth--
				add(ruleFALSE, position1379)
			}
			return true
		l1378:
			position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
			return false
		},
		/* 117 Wildcard <- <(<((ident ':' !':')? '*')> Action87)> */
		func() bool {
			position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
			{
				position1392 := position
				depth++
				{
					position1393 := position
					depth++
					{
						position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
						if !_rules[ruleident]() {
							goto l1394
						}
						if buffer[position] != rune(':') {
							goto l1394
						}
						position++
						{
							position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l1396
							}
							position++
							goto l1394
						l1396:
							position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
						}
						goto l1395
					l1394:
						position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
					}
				l1395:
					if buffer[position] != rune('*') {
						goto l1391
					}
					position++
					depth--
					add(rulePegText, position1393)
				}
				if !_rules[ruleAction87]() {
					goto l1391
				}
				depth--
				add(ruleWildcard, position1392)
			}
			return true
		l1391:
			position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
			return false
		},
		/* 118 StringLiteral <- <(<('"' (('"' '"') / (!'"' .))* '"')> Action88)> */
		func() bool {
			position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
			{
				position1398 := position
				depth++
				{
					position1399 := position
					depth++
					if buffer[position] != rune('"') {
						goto l1397
					}
					position++
				l1400:
					{
						position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
						{
							position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l1403
							}
							position++
							if buffer[position] != rune('"') {
								goto l1403
							}
							position++
							goto l1402
						l1403:
							position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
							{
								position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
								if buffer[position] != rune('"') {
									goto l1404
								}
								position++
								goto l1401
							l1404:
								position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
							}
							if !matchDot() {
								goto l1401
							}
						}
					l1402:
						goto l1400
					l1401:
						position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
					}
					if buffer[position] != rune('"') {
						goto l1397
					}
					position++
					depth--
					add(rulePegText, position1399)
				}
				if !_rules[ruleAction88]() {
					goto l1397
				}
				depth--
				add(ruleStringLiteral, position1398)
			}
			return true
		l1397:
			position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
			return false
		},
		/* 119 ISTREAM <- <(<(('i' / 'I') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action89)> */
		func() bool {
			position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
			{
				position1406 := position
				depth++
				{
					position1407 := position
					depth++
					{
						position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1409
						}
						position++
						goto l1408
					l1409:
						position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
						if buffer[position] != rune('I') {
							goto l1405
						}
						position++
					}
				l1408:
					{
						position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1411
						}
						position++
						goto l1410
					l1411:
						position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
						if buffer[position] != rune('S') {
							goto l1405
						}
						position++
					}
				l1410:
					{
						position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1413
						}
						position++
						goto l1412
					l1413:
						position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
						if buffer[position] != rune('T') {
							goto l1405
						}
						position++
					}
				l1412:
					{
						position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1415
						}
						position++
						goto l1414
					l1415:
						position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
						if buffer[position] != rune('R') {
							goto l1405
						}
						position++
					}
				l1414:
					{
						position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1417
						}
						position++
						goto l1416
					l1417:
						position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
						if buffer[position] != rune('E') {
							goto l1405
						}
						position++
					}
				l1416:
					{
						position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1419
						}
						position++
						goto l1418
					l1419:
						position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
						if buffer[position] != rune('A') {
							goto l1405
						}
						position++
					}
				l1418:
					{
						position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1421
						}
						position++
						goto l1420
					l1421:
						position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
						if buffer[position] != rune('M') {
							goto l1405
						}
						position++
					}
				l1420:
					depth--
					add(rulePegText, position1407)
				}
				if !_rules[ruleAction89]() {
					goto l1405
				}
				depth--
				add(ruleISTREAM, position1406)
			}
			return true
		l1405:
			position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
			return false
		},
		/* 120 DSTREAM <- <(<(('d' / 'D') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action90)> */
		func() bool {
			position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
			{
				position1423 := position
				depth++
				{
					position1424 := position
					depth++
					{
						position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1426
						}
						position++
						goto l1425
					l1426:
						position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
						if buffer[position] != rune('D') {
							goto l1422
						}
						position++
					}
				l1425:
					{
						position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1428
						}
						position++
						goto l1427
					l1428:
						position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
						if buffer[position] != rune('S') {
							goto l1422
						}
						position++
					}
				l1427:
					{
						position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1430
						}
						position++
						goto l1429
					l1430:
						position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
						if buffer[position] != rune('T') {
							goto l1422
						}
						position++
					}
				l1429:
					{
						position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1432
						}
						position++
						goto l1431
					l1432:
						position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
						if buffer[position] != rune('R') {
							goto l1422
						}
						position++
					}
				l1431:
					{
						position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1434
						}
						position++
						goto l1433
					l1434:
						position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
						if buffer[position] != rune('E') {
							goto l1422
						}
						position++
					}
				l1433:
					{
						position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1436
						}
						position++
						goto l1435
					l1436:
						position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
						if buffer[position] != rune('A') {
							goto l1422
						}
						position++
					}
				l1435:
					{
						position1437, tokenIndex1437, depth1437 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1438
						}
						position++
						goto l1437
					l1438:
						position, tokenIndex, depth = position1437, tokenIndex1437, depth1437
						if buffer[position] != rune('M') {
							goto l1422
						}
						position++
					}
				l1437:
					depth--
					add(rulePegText, position1424)
				}
				if !_rules[ruleAction90]() {
					goto l1422
				}
				depth--
				add(ruleDSTREAM, position1423)
			}
			return true
		l1422:
			position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
			return false
		},
		/* 121 RSTREAM <- <(<(('r' / 'R') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('a' / 'A') ('m' / 'M'))> Action91)> */
		func() bool {
			position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
			{
				position1440 := position
				depth++
				{
					position1441 := position
					depth++
					{
						position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1443
						}
						position++
						goto l1442
					l1443:
						position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
						if buffer[position] != rune('R') {
							goto l1439
						}
						position++
					}
				l1442:
					{
						position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1445
						}
						position++
						goto l1444
					l1445:
						position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
						if buffer[position] != rune('S') {
							goto l1439
						}
						position++
					}
				l1444:
					{
						position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1447
						}
						position++
						goto l1446
					l1447:
						position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
						if buffer[position] != rune('T') {
							goto l1439
						}
						position++
					}
				l1446:
					{
						position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1449
						}
						position++
						goto l1448
					l1449:
						position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
						if buffer[position] != rune('R') {
							goto l1439
						}
						position++
					}
				l1448:
					{
						position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1451
						}
						position++
						goto l1450
					l1451:
						position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
						if buffer[position] != rune('E') {
							goto l1439
						}
						position++
					}
				l1450:
					{
						position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1453
						}
						position++
						goto l1452
					l1453:
						position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
						if buffer[position] != rune('A') {
							goto l1439
						}
						position++
					}
				l1452:
					{
						position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1455
						}
						position++
						goto l1454
					l1455:
						position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
						if buffer[position] != rune('M') {
							goto l1439
						}
						position++
					}
				l1454:
					depth--
					add(rulePegText, position1441)
				}
				if !_rules[ruleAction91]() {
					goto l1439
				}
				depth--
				add(ruleRSTREAM, position1440)
			}
			return true
		l1439:
			position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
			return false
		},
		/* 122 TUPLES <- <(<(('t' / 'T') ('u' / 'U') ('p' / 'P') ('l' / 'L') ('e' / 'E') ('s' / 'S'))> Action92)> */
		func() bool {
			position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
			{
				position1457 := position
				depth++
				{
					position1458 := position
					depth++
					{
						position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1460
						}
						position++
						goto l1459
					l1460:
						position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
						if buffer[position] != rune('T') {
							goto l1456
						}
						position++
					}
				l1459:
					{
						position1461, tokenIndex1461, depth1461 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1462
						}
						position++
						goto l1461
					l1462:
						position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
						if buffer[position] != rune('U') {
							goto l1456
						}
						position++
					}
				l1461:
					{
						position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1464
						}
						position++
						goto l1463
					l1464:
						position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
						if buffer[position] != rune('P') {
							goto l1456
						}
						position++
					}
				l1463:
					{
						position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1466
						}
						position++
						goto l1465
					l1466:
						position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
						if buffer[position] != rune('L') {
							goto l1456
						}
						position++
					}
				l1465:
					{
						position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1468
						}
						position++
						goto l1467
					l1468:
						position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
						if buffer[position] != rune('E') {
							goto l1456
						}
						position++
					}
				l1467:
					{
						position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1470
						}
						position++
						goto l1469
					l1470:
						position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
						if buffer[position] != rune('S') {
							goto l1456
						}
						position++
					}
				l1469:
					depth--
					add(rulePegText, position1458)
				}
				if !_rules[ruleAction92]() {
					goto l1456
				}
				depth--
				add(ruleTUPLES, position1457)
			}
			return true
		l1456:
			position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
			return false
		},
		/* 123 SECONDS <- <(<(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action93)> */
		func() bool {
			position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
			{
				position1472 := position
				depth++
				{
					position1473 := position
					depth++
					{
						position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1475
						}
						position++
						goto l1474
					l1475:
						position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
						if buffer[position] != rune('S') {
							goto l1471
						}
						position++
					}
				l1474:
					{
						position1476, tokenIndex1476, depth1476 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1477
						}
						position++
						goto l1476
					l1477:
						position, tokenIndex, depth = position1476, tokenIndex1476, depth1476
						if buffer[position] != rune('E') {
							goto l1471
						}
						position++
					}
				l1476:
					{
						position1478, tokenIndex1478, depth1478 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1479
						}
						position++
						goto l1478
					l1479:
						position, tokenIndex, depth = position1478, tokenIndex1478, depth1478
						if buffer[position] != rune('C') {
							goto l1471
						}
						position++
					}
				l1478:
					{
						position1480, tokenIndex1480, depth1480 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1481
						}
						position++
						goto l1480
					l1481:
						position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
						if buffer[position] != rune('O') {
							goto l1471
						}
						position++
					}
				l1480:
					{
						position1482, tokenIndex1482, depth1482 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1483
						}
						position++
						goto l1482
					l1483:
						position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
						if buffer[position] != rune('N') {
							goto l1471
						}
						position++
					}
				l1482:
					{
						position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1485
						}
						position++
						goto l1484
					l1485:
						position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
						if buffer[position] != rune('D') {
							goto l1471
						}
						position++
					}
				l1484:
					{
						position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1487
						}
						position++
						goto l1486
					l1487:
						position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
						if buffer[position] != rune('S') {
							goto l1471
						}
						position++
					}
				l1486:
					depth--
					add(rulePegText, position1473)
				}
				if !_rules[ruleAction93]() {
					goto l1471
				}
				depth--
				add(ruleSECONDS, position1472)
			}
			return true
		l1471:
			position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
			return false
		},
		/* 124 MILLISECONDS <- <(<(('m' / 'M') ('i' / 'I') ('l' / 'L') ('l' / 'L') ('i' / 'I') ('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S'))> Action94)> */
		func() bool {
			position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
			{
				position1489 := position
				depth++
				{
					position1490 := position
					depth++
					{
						position1491, tokenIndex1491, depth1491 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1492
						}
						position++
						goto l1491
					l1492:
						position, tokenIndex, depth = position1491, tokenIndex1491, depth1491
						if buffer[position] != rune('M') {
							goto l1488
						}
						position++
					}
				l1491:
					{
						position1493, tokenIndex1493, depth1493 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1494
						}
						position++
						goto l1493
					l1494:
						position, tokenIndex, depth = position1493, tokenIndex1493, depth1493
						if buffer[position] != rune('I') {
							goto l1488
						}
						position++
					}
				l1493:
					{
						position1495, tokenIndex1495, depth1495 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1496
						}
						position++
						goto l1495
					l1496:
						position, tokenIndex, depth = position1495, tokenIndex1495, depth1495
						if buffer[position] != rune('L') {
							goto l1488
						}
						position++
					}
				l1495:
					{
						position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1498
						}
						position++
						goto l1497
					l1498:
						position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
						if buffer[position] != rune('L') {
							goto l1488
						}
						position++
					}
				l1497:
					{
						position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1500
						}
						position++
						goto l1499
					l1500:
						position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
						if buffer[position] != rune('I') {
							goto l1488
						}
						position++
					}
				l1499:
					{
						position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1502
						}
						position++
						goto l1501
					l1502:
						position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
						if buffer[position] != rune('S') {
							goto l1488
						}
						position++
					}
				l1501:
					{
						position1503, tokenIndex1503, depth1503 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1504
						}
						position++
						goto l1503
					l1504:
						position, tokenIndex, depth = position1503, tokenIndex1503, depth1503
						if buffer[position] != rune('E') {
							goto l1488
						}
						position++
					}
				l1503:
					{
						position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1506
						}
						position++
						goto l1505
					l1506:
						position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
						if buffer[position] != rune('C') {
							goto l1488
						}
						position++
					}
				l1505:
					{
						position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1508
						}
						position++
						goto l1507
					l1508:
						position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
						if buffer[position] != rune('O') {
							goto l1488
						}
						position++
					}
				l1507:
					{
						position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1510
						}
						position++
						goto l1509
					l1510:
						position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
						if buffer[position] != rune('N') {
							goto l1488
						}
						position++
					}
				l1509:
					{
						position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1512
						}
						position++
						goto l1511
					l1512:
						position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
						if buffer[position] != rune('D') {
							goto l1488
						}
						position++
					}
				l1511:
					{
						position1513, tokenIndex1513, depth1513 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1514
						}
						position++
						goto l1513
					l1514:
						position, tokenIndex, depth = position1513, tokenIndex1513, depth1513
						if buffer[position] != rune('S') {
							goto l1488
						}
						position++
					}
				l1513:
					depth--
					add(rulePegText, position1490)
				}
				if !_rules[ruleAction94]() {
					goto l1488
				}
				depth--
				add(ruleMILLISECONDS, position1489)
			}
			return true
		l1488:
			position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
			return false
		},
		/* 125 Wait <- <(<(('w' / 'W') ('a' / 'A') ('i' / 'I') ('t' / 'T'))> Action95)> */
		func() bool {
			position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
			{
				position1516 := position
				depth++
				{
					position1517 := position
					depth++
					{
						position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
						if buffer[position] != rune('w') {
							goto l1519
						}
						position++
						goto l1518
					l1519:
						position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
						if buffer[position] != rune('W') {
							goto l1515
						}
						position++
					}
				l1518:
					{
						position1520, tokenIndex1520, depth1520 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1521
						}
						position++
						goto l1520
					l1521:
						position, tokenIndex, depth = position1520, tokenIndex1520, depth1520
						if buffer[position] != rune('A') {
							goto l1515
						}
						position++
					}
				l1520:
					{
						position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1523
						}
						position++
						goto l1522
					l1523:
						position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
						if buffer[position] != rune('I') {
							goto l1515
						}
						position++
					}
				l1522:
					{
						position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1525
						}
						position++
						goto l1524
					l1525:
						position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
						if buffer[position] != rune('T') {
							goto l1515
						}
						position++
					}
				l1524:
					depth--
					add(rulePegText, position1517)
				}
				if !_rules[ruleAction95]() {
					goto l1515
				}
				depth--
				add(ruleWait, position1516)
			}
			return true
		l1515:
			position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
			return false
		},
		/* 126 DropOldest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('o' / 'O') ('l' / 'L') ('d' / 'D') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action96)> */
		func() bool {
			position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
			{
				position1527 := position
				depth++
				{
					position1528 := position
					depth++
					{
						position1529, tokenIndex1529, depth1529 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1530
						}
						position++
						goto l1529
					l1530:
						position, tokenIndex, depth = position1529, tokenIndex1529, depth1529
						if buffer[position] != rune('D') {
							goto l1526
						}
						position++
					}
				l1529:
					{
						position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1532
						}
						position++
						goto l1531
					l1532:
						position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
						if buffer[position] != rune('R') {
							goto l1526
						}
						position++
					}
				l1531:
					{
						position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1534
						}
						position++
						goto l1533
					l1534:
						position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
						if buffer[position] != rune('O') {
							goto l1526
						}
						position++
					}
				l1533:
					{
						position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1536
						}
						position++
						goto l1535
					l1536:
						position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
						if buffer[position] != rune('P') {
							goto l1526
						}
						position++
					}
				l1535:
					if !_rules[rulesp]() {
						goto l1526
					}
					{
						position1537, tokenIndex1537, depth1537 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1538
						}
						position++
						goto l1537
					l1538:
						position, tokenIndex, depth = position1537, tokenIndex1537, depth1537
						if buffer[position] != rune('O') {
							goto l1526
						}
						position++
					}
				l1537:
					{
						position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1540
						}
						position++
						goto l1539
					l1540:
						position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
						if buffer[position] != rune('L') {
							goto l1526
						}
						position++
					}
				l1539:
					{
						position1541, tokenIndex1541, depth1541 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1542
						}
						position++
						goto l1541
					l1542:
						position, tokenIndex, depth = position1541, tokenIndex1541, depth1541
						if buffer[position] != rune('D') {
							goto l1526
						}
						position++
					}
				l1541:
					{
						position1543, tokenIndex1543, depth1543 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1544
						}
						position++
						goto l1543
					l1544:
						position, tokenIndex, depth = position1543, tokenIndex1543, depth1543
						if buffer[position] != rune('E') {
							goto l1526
						}
						position++
					}
				l1543:
					{
						position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1546
						}
						position++
						goto l1545
					l1546:
						position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
						if buffer[position] != rune('S') {
							goto l1526
						}
						position++
					}
				l1545:
					{
						position1547, tokenIndex1547, depth1547 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1548
						}
						position++
						goto l1547
					l1548:
						position, tokenIndex, depth = position1547, tokenIndex1547, depth1547
						if buffer[position] != rune('T') {
							goto l1526
						}
						position++
					}
				l1547:
					depth--
					add(rulePegText, position1528)
				}
				if !_rules[ruleAction96]() {
					goto l1526
				}
				depth--
				add(ruleDropOldest, position1527)
			}
			return true
		l1526:
			position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
			return false
		},
		/* 127 DropNewest <- <(<(('d' / 'D') ('r' / 'R') ('o' / 'O') ('p' / 'P') sp (('n' / 'N') ('e' / 'E') ('w' / 'W') ('e' / 'E') ('s' / 'S') ('t' / 'T')))> Action97)> */
		func() bool {
			position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
			{
				position1550 := position
				depth++
				{
					position1551 := position
					depth++
					{
						position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1553
						}
						position++
						goto l1552
					l1553:
						position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
						if buffer[position] != rune('D') {
							goto l1549
						}
						position++
					}
				l1552:
					{
						position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1555
						}
						position++
						goto l1554
					l1555:
						position, tokenIndex, depth = position1554, tokenIndex1554, depth1554
						if buffer[position] != rune('R') {
							goto l1549
						}
						position++
					}
				l1554:
					{
						position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1557
						}
						position++
						goto l1556
					l1557:
						position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
						if buffer[position] != rune('O') {
							goto l1549
						}
						position++
					}
				l1556:
					{
						position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1559
						}
						position++
						goto l1558
					l1559:
						position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
						if buffer[position] != rune('P') {
							goto l1549
						}
						position++
					}
				l1558:
					if !_rules[rulesp]() {
						goto l1549
					}
					{
						position1560, tokenIndex1560, depth1560 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1561
						}
						position++
						goto l1560
					l1561:
						position, tokenIndex, depth = position1560, tokenIndex1560, depth1560
						if buffer[position] != rune('N') {
							goto l1549
						}
						position++
					}
				l1560:
					{
						position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1563
						}
						position++
						goto l1562
					l1563:
						position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
						if buffer[position] != rune('E') {
							goto l1549
						}
						position++
					}
				l1562:
					{
						position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
						if buffer[position] != rune('w') {
							goto l1565
						}
						position++
						goto l1564
					l1565:
						position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
						if buffer[position] != rune('W') {
							goto l1549
						}
						position++
					}
				l1564:
					{
						position1566, tokenIndex1566, depth1566 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1567
						}
						position++
						goto l1566
					l1567:
						position, tokenIndex, depth = position1566, tokenIndex1566, depth1566
						if buffer[position] != rune('E') {
							goto l1549
						}
						position++
					}
				l1566:
					{
						position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1569
						}
						position++
						goto l1568
					l1569:
						position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
						if buffer[position] != rune('S') {
							goto l1549
						}
						position++
					}
				l1568:
					{
						position1570, tokenIndex1570, depth1570 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1571
						}
						position++
						goto l1570
					l1571:
						position, tokenIndex, depth = position1570, tokenIndex1570, depth1570
						if buffer[position] != rune('T') {
							goto l1549
						}
						position++
					}
				l1570:
					depth--
					add(rulePegText, position1551)
				}
				if !_rules[ruleAction97]() {
					goto l1549
				}
				depth--
				add(ruleDropNewest, position1550)
			}
			return true
		l1549:
			position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
			return false
		},
		/* 128 StreamIdentifier <- <(<ident> Action98)> */
		func() bool {
			position1572, tokenIndex1572, depth1572 := position, tokenIndex, depth
			{
				position1573 := position
				depth++
				{
					position1574 := position
					depth++
					if !_rules[ruleident]() {
						goto l1572
					}
					depth--
					add(rulePegText, position1574)
				}
				if !_rules[ruleAction98]() {
					goto l1572
				}
				depth--
				add(ruleStreamIdentifier, position1573)
			}
			return true
		l1572:
			position, tokenIndex, depth = position1572, tokenIndex1572, depth1572
			return false
		},
		/* 129 SourceSinkType <- <(<ident> Action99)> */
		func() bool {
			position1575, tokenIndex1575, depth1575 := position, tokenIndex, depth
			{
				position1576 := position
				depth++
				{
					position1577 := position
					depth++
					if !_rules[ruleident]() {
						goto l1575
					}
					depth--
					add(rulePegText, position1577)
				}
				if !_rules[ruleAction99]() {
					goto l1575
				}
				depth--
				add(ruleSourceSinkType, position1576)
			}
			return true
		l1575:
			position, tokenIndex, depth = position1575, tokenIndex1575, depth1575
			return false
		},
		/* 130 SourceSinkParamKey <- <(<ident> Action100)> */
		func() bool {
			position1578, tokenIndex1578, depth1578 := position, tokenIndex, depth
			{
				position1579 := position
				depth++
				{
					position1580 := position
					depth++
					if !_rules[ruleident]() {
						goto l1578
					}
					depth--
					add(rulePegText, position1580)
				}
				if !_rules[ruleAction100]() {
					goto l1578
				}
				depth--
				add(ruleSourceSinkParamKey, position1579)
			}
			return true
		l1578:
			position, tokenIndex, depth = position1578, tokenIndex1578, depth1578
			return false
		},
		/* 131 Paused <- <(<(('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action101)> */
		func() bool {
			position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
			{
				position1582 := position
				depth++
				{
					position1583 := position
					depth++
					{
						position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1585
						}
						position++
						goto l1584
					l1585:
						position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
						if buffer[position] != rune('P') {
							goto l1581
						}
						position++
					}
				l1584:
					{
						position1586, tokenIndex1586, depth1586 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1587
						}
						position++
						goto l1586
					l1587:
						position, tokenIndex, depth = position1586, tokenIndex1586, depth1586
						if buffer[position] != rune('A') {
							goto l1581
						}
						position++
					}
				l1586:
					{
						position1588, tokenIndex1588, depth1588 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1589
						}
						position++
						goto l1588
					l1589:
						position, tokenIndex, depth = position1588, tokenIndex1588, depth1588
						if buffer[position] != rune('U') {
							goto l1581
						}
						position++
					}
				l1588:
					{
						position1590, tokenIndex1590, depth1590 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1591
						}
						position++
						goto l1590
					l1591:
						position, tokenIndex, depth = position1590, tokenIndex1590, depth1590
						if buffer[position] != rune('S') {
							goto l1581
						}
						position++
					}
				l1590:
					{
						position1592, tokenIndex1592, depth1592 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1593
						}
						position++
						goto l1592
					l1593:
						position, tokenIndex, depth = position1592, tokenIndex1592, depth1592
						if buffer[position] != rune('E') {
							goto l1581
						}
						position++
					}
				l1592:
					{
						position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1595
						}
						position++
						goto l1594
					l1595:
						position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
						if buffer[position] != rune('D') {
							goto l1581
						}
						position++
					}
				l1594:
					depth--
					add(rulePegText, position1583)
				}
				if !_rules[ruleAction101]() {
					goto l1581
				}
				depth--
				add(rulePaused, position1582)
			}
			return true
		l1581:
			position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
			return false
		},
		/* 132 Unpaused <- <(<(('u' / 'U') ('n' / 'N') ('p' / 'P') ('a' / 'A') ('u' / 'U') ('s' / 'S') ('e' / 'E') ('d' / 'D'))> Action102)> */
		func() bool {
			position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
			{
				position1597 := position
				depth++
				{
					position1598 := position
					depth++
					{
						position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1600
						}
						position++
						goto l1599
					l1600:
						position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
						if buffer[position] != rune('U') {
							goto l1596
						}
						position++
					}
				l1599:
					{
						position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1602
						}
						position++
						goto l1601
					l1602:
						position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
						if buffer[position] != rune('N') {
							goto l1596
						}
						position++
					}
				l1601:
					{
						position1603, tokenIndex1603, depth1603 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1604
						}
						position++
						goto l1603
					l1604:
						position, tokenIndex, depth = position1603, tokenIndex1603, depth1603
						if buffer[position] != rune('P') {
							goto l1596
						}
						position++
					}
				l1603:
					{
						position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1606
						}
						position++
						goto l1605
					l1606:
						position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
						if buffer[position] != rune('A') {
							goto l1596
						}
						position++
					}
				l1605:
					{
						position1607, tokenIndex1607, depth1607 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l1608
						}
						position++
						goto l1607
					l1608:
						position, tokenIndex, depth = position1607, tokenIndex1607, depth1607
						if buffer[position] != rune('U') {
							goto l1596
						}
						position++
					}
				l1607:
					{
						position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1610
						}
						position++
						goto l1609
					l1610:
						position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
						if buffer[position] != rune('S') {
							goto l1596
						}
						position++
					}
				l1609:
					{
						position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1612
						}
						position++
						goto l1611
					l1612:
						position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
						if buffer[position] != rune('E') {
							goto l1596
						}
						position++
					}
				l1611:
					{
						position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1614
						}
						position++
						goto l1613
					l1614:
						position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
						if buffer[position] != rune('D') {
							goto l1596
						}
						position++
					}
				l1613:
					depth--
					add(rulePegText, position1598)
				}
				if !_rules[ruleAction102]() {
					goto l1596
				}
				depth--
				add(ruleUnpaused, position1597)
			}
			return true
		l1596:
			position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
			return false
		},
		/* 133 Ascending <- <(<(('a' / 'A') ('s' / 'S') ('c' / 'C'))> Action103)> */
		func() bool {
			position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
			{
				position1616 := position
				depth++
				{
					position1617 := position
					depth++
					{
						position1618, tokenIndex1618, depth1618 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1619
						}
						position++
						goto l1618
					l1619:
						position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
						if buffer[position] != rune('A') {
							goto l1615
						}
						position++
					}
				l1618:
					{
						position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1621
						}
						position++
						goto l1620
					l1621:
						position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
						if buffer[position] != rune('S') {
							goto l1615
						}
						position++
					}
				l1620:
					{
						position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1623
						}
						position++
						goto l1622
					l1623:
						position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
						if buffer[position] != rune('C') {
							goto l1615
						}
						position++
					}
				l1622:
					depth--
					add(rulePegText, position1617)
				}
				if !_rules[ruleAction103]() {
					goto l1615
				}
				depth--
				add(ruleAscending, position1616)
			}
			return true
		l1615:
			position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
			return false
		},
		/* 134 Descending <- <(<(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C'))> Action104)> */
		func() bool {
			position1624, tokenIndex1624, depth1624 := position, tokenIndex, depth
			{
				position1625 := position
				depth++
				{
					position1626 := position
					depth++
					{
						position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1628
						}
						position++
						goto l1627
					l1628:
						position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
						if buffer[position] != rune('D') {
							goto l1624
						}
						position++
					}
				l1627:
					{
						position1629, tokenIndex1629, depth1629 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1630
						}
						position++
						goto l1629
					l1630:
						position, tokenIndex, depth = position1629, tokenIndex1629, depth1629
						if buffer[position] != rune('E') {
							goto l1624
						}
						position++
					}
				l1629:
					{
						position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1632
						}
						position++
						goto l1631
					l1632:
						position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
						if buffer[position] != rune('S') {
							goto l1624
						}
						position++
					}
				l1631:
					{
						position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l1634
						}
						position++
						goto l1633
					l1634:
						position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
						if buffer[position] != rune('C') {
							goto l1624
						}
						position++
					}
				l1633:
					depth--
					add(rulePegText, position1626)
				}
				if !_rules[ruleAction104]() {
					goto l1624
				}
				depth--
				add(ruleDescending, position1625)
			}
			return true
		l1624:
			position, tokenIndex, depth = position1624, tokenIndex1624, depth1624
			return false
		},
		/* 135 Type <- <(Bool / Int / Float / String / Blob / Timestamp / Array / Map)> */
		func() bool {
			position1635, tokenIndex1635, depth1635 := position, tokenIndex, depth
			{
				position1636 := position
				depth++
				{
					position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
					if !_rules[ruleBool]() {
						goto l1638
					}
					goto l1637
				l1638:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleInt]() {
						goto l1639
					}
					goto l1637
				l1639:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleFloat]() {
						goto l1640
					}
					goto l1637
				l1640:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleString]() {
						goto l1641
					}
					goto l1637
				l1641:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleBlob]() {
						goto l1642
					}
					goto l1637
				l1642:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleTimestamp]() {
						goto l1643
					}
					goto l1637
				l1643:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleArray]() {
						goto l1644
					}
					goto l1637
				l1644:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if !_rules[ruleMap]() {
						goto l1635
					}
				}
			l1637:
				depth--
				add(ruleType, position1636)
			}
			return true
		l1635:
			position, tokenIndex, depth = position1635, tokenIndex1635, depth1635
			return false
		},
		/* 136 Bool <- <(<(('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L'))> Action105)> */
		func() bool {
			position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
			{
				position1646 := position
				depth++
				{
					position1647 := position
					depth++
					{
						position1648, tokenIndex1648, depth1648 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1649
						}
						position++
						goto l1648
					l1649:
						position, tokenIndex, depth = position1648, tokenIndex1648, depth1648
						if buffer[position] != rune('B') {
							goto l1645
						}
						position++
					}
				l1648:
					{
						position1650, tokenIndex1650, depth1650 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1651
						}
						position++
						goto l1650
					l1651:
						position, tokenIndex, depth = position1650, tokenIndex1650, depth1650
						if buffer[position] != rune('O') {
							goto l1645
						}
						position++
					}
				l1650:
					{
						position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1653
						}
						position++
						goto l1652
					l1653:
						position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
						if buffer[position] != rune('O') {
							goto l1645
						}
						position++
					}
				l1652:
					{
						position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1655
						}
						position++
						goto l1654
					l1655:
						position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
						if buffer[position] != rune('L') {
							goto l1645
						}
						position++
					}
				l1654:
					depth--
					add(rulePegText, position1647)
				}
				if !_rules[ruleAction105]() {
					goto l1645
				}
				depth--
				add(ruleBool, position1646)
			}
			return true
		l1645:
			position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
			return false
		},
		/* 137 Int <- <(<(('i' / 'I') ('n' / 'N') ('t' / 'T'))> Action106)> */
		func() bool {
			position1656, tokenIndex1656, depth1656 := position, tokenIndex, depth
			{
				position1657 := position
				depth++
				{
					position1658 := position
					depth++
					{
						position1659, tokenIndex1659, depth1659 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1660
						}
						position++
						goto l1659
					l1660:
						position, tokenIndex, depth = position1659, tokenIndex1659, depth1659
						if buffer[position] != rune('I') {
							goto l1656
						}
						position++
					}
				l1659:
					{
						position1661, tokenIndex1661, depth1661 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1662
						}
						position++
						goto l1661
					l1662:
						position, tokenIndex, depth = position1661, tokenIndex1661, depth1661
						if buffer[position] != rune('N') {
							goto l1656
						}
						position++
					}
				l1661:
					{
						position1663, tokenIndex1663, depth1663 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1664
						}
						position++
						goto l1663
					l1664:
						position, tokenIndex, depth = position1663, tokenIndex1663, depth1663
						if buffer[position] != rune('T') {
							goto l1656
						}
						position++
					}
				l1663:
					depth--
					add(rulePegText, position1658)
				}
				if !_rules[ruleAction106]() {
					goto l1656
				}
				depth--
				add(ruleInt, position1657)
			}
			return true
		l1656:
			position, tokenIndex, depth = position1656, tokenIndex1656, depth1656
			return false
		},
		/* 138 Float <- <(<(('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))> Action107)> */
		func() bool {
			position1665, tokenIndex1665, depth1665 := position, tokenIndex, depth
			{
				position1666 := position
				depth++
				{
					position1667 := position
					depth++
					{
						position1668, tokenIndex1668, depth1668 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l1669
						}
						position++
						goto l1668
					l1669:
						position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
						if buffer[position] != rune('F') {
							goto l1665
						}
						position++
					}
				l1668:
					{
						position1670, tokenIndex1670, depth1670 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1671
						}
						position++
						goto l1670
					l1671:
						position, tokenIndex, depth = position1670, tokenIndex1670, depth1670
						if buffer[position] != rune('L') {
							goto l1665
						}
						position++
					}
				l1670:
					{
						position1672, tokenIndex1672, depth1672 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1673
						}
						position++
						goto l1672
					l1673:
						position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
						if buffer[position] != rune('O') {
							goto l1665
						}
						position++
					}
				l1672:
					{
						position1674, tokenIndex1674, depth1674 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1675
						}
						position++
						goto l1674
					l1675:
						position, tokenIndex, depth = position1674, tokenIndex1674, depth1674
						if buffer[position] != rune('A') {
							goto l1665
						}
						position++
					}
				l1674:
					{
						position1676, tokenIndex1676, depth1676 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1677
						}
						position++
						goto l1676
					l1677:
						position, tokenIndex, depth = position1676, tokenIndex1676, depth1676
						if buffer[position] != rune('T') {
							goto l1665
						}
						position++
					}
				l1676:
					depth--
					add(rulePegText, position1667)
				}
				if !_rules[ruleAction107]() {
					goto l1665
				}
				depth--
				add(ruleFloat, position1666)
			}
			return true
		l1665:
			position, tokenIndex, depth = position1665, tokenIndex1665, depth1665
			return false
		},
		/* 139 String <- <(<(('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G'))> Action108)> */
		func() bool {
			position1678, tokenIndex1678, depth1678 := position, tokenIndex, depth
			{
				position1679 := position
				depth++
				{
					position1680 := position
					depth++
					{
						position1681, tokenIndex1681, depth1681 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1682
						}
						position++
						goto l1681
					l1682:
						position, tokenIndex, depth = position1681, tokenIndex1681, depth1681
						if buffer[position] != rune('S') {
							goto l1678
						}
						position++
					}
				l1681:
					{
						position1683, tokenIndex1683, depth1683 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1684
						}
						position++
						goto l1683
					l1684:
						position, tokenIndex, depth = position1683, tokenIndex1683, depth1683
						if buffer[position] != rune('T') {
							goto l1678
						}
						position++
					}
				l1683:
					{
						position1685, tokenIndex1685, depth1685 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1686
						}
						position++
						goto l1685
					l1686:
						position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
						if buffer[position] != rune('R') {
							goto l1678
						}
						position++
					}
				l1685:
					{
						position1687, tokenIndex1687, depth1687 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1688
						}
						position++
						goto l1687
					l1688:
						position, tokenIndex, depth = position1687, tokenIndex1687, depth1687
						if buffer[position] != rune('I') {
							goto l1678
						}
						position++
					}
				l1687:
					{
						position1689, tokenIndex1689, depth1689 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1690
						}
						position++
						goto l1689
					l1690:
						position, tokenIndex, depth = position1689, tokenIndex1689, depth1689
						if buffer[position] != rune('N') {
							goto l1678
						}
						position++
					}
				l1689:
					{
						position1691, tokenIndex1691, depth1691 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l1692
						}
						position++
						goto l1691
					l1692:
						position, tokenIndex, depth = position1691, tokenIndex1691, depth1691
						if buffer[position] != rune('G') {
							goto l1678
						}
						position++
					}
				l1691:
					depth--
					add(rulePegText, position1680)
				}
				if !_rules[ruleAction108]() {
					goto l1678
				}
				depth--
				add(ruleString, position1679)
			}
			return true
		l1678:
			position, tokenIndex, depth = position1678, tokenIndex1678, depth1678
			return false
		},
		/* 140 Blob <- <(<(('b' / 'B') ('l' / 'L') ('o' / 'O') ('b' / 'B'))> Action109)> */
		func() bool {
			position1693, tokenIndex1693, depth1693 := position, tokenIndex, depth
			{
				position1694 := position
				depth++
				{
					position1695 := position
					depth++
					{
						position1696, tokenIndex1696, depth1696 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1697
						}
						position++
						goto l1696
					l1697:
						position, tokenIndex, depth = position1696, tokenIndex1696, depth1696
						if buffer[position] != rune('B') {
							goto l1693
						}
						position++
					}
				l1696:
					{
						position1698, tokenIndex1698, depth1698 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l1699
						}
						position++
						goto l1698
					l1699:
						position, tokenIndex, depth = position1698, tokenIndex1698, depth1698
						if buffer[position] != rune('L') {
							goto l1693
						}
						position++
					}
				l1698:
					{
						position1700, tokenIndex1700, depth1700 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1701
						}
						position++
						goto l1700
					l1701:
						position, tokenIndex, depth = position1700, tokenIndex1700, depth1700
						if buffer[position] != rune('O') {
							goto l1693
						}
						position++
					}
				l1700:
					{
						position1702, tokenIndex1702, depth1702 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l1703
						}
						position++
						goto l1702
					l1703:
						position, tokenIndex, depth = position1702, tokenIndex1702, depth1702
						if buffer[position] != rune('B') {
							goto l1693
						}
						position++
					}
				l1702:
					depth--
					add(rulePegText, position1695)
				}
				if !_rules[ruleAction109]() {
					goto l1693
				}
				depth--
				add(ruleBlob, position1694)
			}
			return true
		l1693:
			position, tokenIndex, depth = position1693, tokenIndex1693, depth1693
			return false
		},
		/* 141 Timestamp <- <(<(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('m' / 'M') ('p' / 'P'))> Action110)> */
		func() bool {
			position1704, tokenIndex1704, depth1704 := position, tokenIndex, depth
			{
				position1705 := position
				depth++
				{
					position1706 := position
					depth++
					{
						position1707, tokenIndex1707, depth1707 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1708
						}
						position++
						goto l1707
					l1708:
						position, tokenIndex, depth = position1707, tokenIndex1707, depth1707
						if buffer[position] != rune('T') {
							goto l1704
						}
						position++
					}
				l1707:
					{
						position1709, tokenIndex1709, depth1709 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1710
						}
						position++
						goto l1709
					l1710:
						position, tokenIndex, depth = position1709, tokenIndex1709, depth1709
						if buffer[position] != rune('I') {
							goto l1704
						}
						position++
					}
				l1709:
					{
						position1711, tokenIndex1711, depth1711 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1712
						}
						position++
						goto l1711
					l1712:
						position, tokenIndex, depth = position1711, tokenIndex1711, depth1711
						if buffer[position] != rune('M') {
							goto l1704
						}
						position++
					}
				l1711:
					{
						position1713, tokenIndex1713, depth1713 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l1714
						}
						position++
						goto l1713
					l1714:
						position, tokenIndex, depth = position1713, tokenIndex1713, depth1713
						if buffer[position] != rune('E') {
							goto l1704
						}
						position++
					}
				l1713:
					{
						position1715, tokenIndex1715, depth1715 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1716
						}
						position++
						goto l1715
					l1716:
						position, tokenIndex, depth = position1715, tokenIndex1715, depth1715
						if buffer[position] != rune('S') {
							goto l1704
						}
						position++
					}
				l1715:
					{
						position1717, tokenIndex1717, depth1717 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1718
						}
						position++
						goto l1717
					l1718:
						position, tokenIndex, depth = position1717, tokenIndex1717, depth1717
						if buffer[position] != rune('T') {
							goto l1704
						}
						position++
					}
				l1717:
					{
						position1719, tokenIndex1719, depth1719 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1720
						}
						position++
						goto l1719
					l1720:
						position, tokenIndex, depth = position1719, tokenIndex1719, depth1719
						if buffer[position] != rune('A') {
							goto l1704
						}
						position++
					}
				l1719:
					{
						position1721, tokenIndex1721, depth1721 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1722
						}
						position++
						goto l1721
					l1722:
						position, tokenIndex, depth = position1721, tokenIndex1721, depth1721
						if buffer[position] != rune('M') {
							goto l1704
						}
						position++
					}
				l1721:
					{
						position1723, tokenIndex1723, depth1723 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1724
						}
						position++
						goto l1723
					l1724:
						position, tokenIndex, depth = position1723, tokenIndex1723, depth1723
						if buffer[position] != rune('P') {
							goto l1704
						}
						position++
					}
				l1723:
					depth--
					add(rulePegText, position1706)
				}
				if !_rules[ruleAction110]() {
					goto l1704
				}
				depth--
				add(ruleTimestamp, position1705)
			}
			return true
		l1704:
			position, tokenIndex, depth = position1704, tokenIndex1704, depth1704
			return false
		},
		/* 142 Array <- <(<(('a' / 'A') ('r' / 'R') ('r' / 'R') ('a' / 'A') ('y' / 'Y'))> Action111)> */
		func() bool {
			position1725, tokenIndex1725, depth1725 := position, tokenIndex, depth
			{
				position1726 := position
				depth++
				{
					position1727 := position
					depth++
					{
						position1728, tokenIndex1728, depth1728 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1729
						}
						position++
						goto l1728
					l1729:
						position, tokenIndex, depth = position1728, tokenIndex1728, depth1728
						if buffer[position] != rune('A') {
							goto l1725
						}
						position++
					}
				l1728:
					{
						position1730, tokenIndex1730, depth1730 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1731
						}
						position++
						goto l1730
					l1731:
						position, tokenIndex, depth = position1730, tokenIndex1730, depth1730
						if buffer[position] != rune('R') {
							goto l1725
						}
						position++
					}
				l1730:
					{
						position1732, tokenIndex1732, depth1732 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1733
						}
						position++
						goto l1732
					l1733:
						position, tokenIndex, depth = position1732, tokenIndex1732, depth1732
						if buffer[position] != rune('R') {
							goto l1725
						}
						position++
					}
				l1732:
					{
						position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1735
						}
						position++
						goto l1734
					l1735:
						position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
						if buffer[position] != rune('A') {
							goto l1725
						}
						position++
					}
				l1734:
					{
						position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l1737
						}
						position++
						goto l1736
					l1737:
						position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
						if buffer[position] != rune('Y') {
							goto l1725
						}
						position++
					}
				l1736:
					depth--
					add(rulePegText, position1727)
				}
				if !_rules[ruleAction111]() {
					goto l1725
				}
				depth--
				add(ruleArray, position1726)
			}
			return true
		l1725:
			position, tokenIndex, depth = position1725, tokenIndex1725, depth1725
			return false
		},
		/* 143 Map <- <(<(('m' / 'M') ('a' / 'A') ('p' / 'P'))> Action112)> */
		func() bool {
			position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
			{
				position1739 := position
				depth++
				{
					position1740 := position
					depth++
					{
						position1741, tokenIndex1741, depth1741 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l1742
						}
						position++
						goto l1741
					l1742:
						position, tokenIndex, depth = position1741, tokenIndex1741, depth1741
						if buffer[position] != rune('M') {
							goto l1738
						}
						position++
					}
				l1741:
					{
						position1743, tokenIndex1743, depth1743 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1744
						}
						position++
						goto l1743
					l1744:
						position, tokenIndex, depth = position1743, tokenIndex1743, depth1743
						if buffer[position] != rune('A') {
							goto l1738
						}
						position++
					}
				l1743:
					{
						position1745, tokenIndex1745, depth1745 := position, tokenIndex, depth
						if buffer[position] != rune('p') {
							goto l1746
						}
						position++
						goto l1745
					l1746:
						position, tokenIndex, depth = position1745, tokenIndex1745, depth1745
						if buffer[position] != rune('P') {
							goto l1738
						}
						position++
					}
				l1745:
					depth--
					add(rulePegText, position1740)
				}
				if !_rules[ruleAction112]() {
					goto l1738
				}
				depth--
				add(ruleMap, position1739)
			}
			return true
		l1738:
			position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
			return false
		},
		/* 144 Or <- <(<(('o' / 'O') ('r' / 'R'))> Action113)> */
		func() bool {
			position1747, tokenIndex1747, depth1747 := position, tokenIndex, depth
			{
				position1748 := position
				depth++
				{
					position1749 := position
					depth++
					{
						position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1751
						}
						position++
						goto l1750
					l1751:
						position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
						if buffer[position] != rune('O') {
							goto l1747
						}
						position++
					}
				l1750:
					{
						position1752, tokenIndex1752, depth1752 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l1753
						}
						position++
						goto l1752
					l1753:
						position, tokenIndex, depth = position1752, tokenIndex1752, depth1752
						if buffer[position] != rune('R') {
							goto l1747
						}
						position++
					}
				l1752:
					depth--
					add(rulePegText, position1749)
				}
				if !_rules[ruleAction113]() {
					goto l1747
				}
				depth--
				add(ruleOr, position1748)
			}
			return true
		l1747:
			position, tokenIndex, depth = position1747, tokenIndex1747, depth1747
			return false
		},
		/* 145 And <- <(<(('a' / 'A') ('n' / 'N') ('d' / 'D'))> Action114)> */
		func() bool {
			position1754, tokenIndex1754, depth1754 := position, tokenIndex, depth
			{
				position1755 := position
				depth++
				{
					position1756 := position
					depth++
					{
						position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l1758
						}
						position++
						goto l1757
					l1758:
						position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
						if buffer[position] != rune('A') {
							goto l1754
						}
						position++
					}
				l1757:
					{
						position1759, tokenIndex1759, depth1759 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1760
						}
						position++
						goto l1759
					l1760:
						position, tokenIndex, depth = position1759, tokenIndex1759, depth1759
						if buffer[position] != rune('N') {
							goto l1754
						}
						position++
					}
				l1759:
					{
						position1761, tokenIndex1761, depth1761 := position, tokenIndex, depth
						if buffer[position] != rune('d') {
							goto l1762
						}
						position++
						goto l1761
					l1762:
						position, tokenIndex, depth = position1761, tokenIndex1761, depth1761
						if buffer[position] != rune('D') {
							goto l1754
						}
						position++
					}
				l1761:
					depth--
					add(rulePegText, position1756)
				}
				if !_rules[ruleAction114]() {
					goto l1754
				}
				depth--
				add(ruleAnd, position1755)
			}
			return true
		l1754:
			position, tokenIndex, depth = position1754, tokenIndex1754, depth1754
			return false
		},
		/* 146 Not <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T'))> Action115)> */
		func() bool {
			position1763, tokenIndex1763, depth1763 := position, tokenIndex, depth
			{
				position1764 := position
				depth++
				{
					position1765 := position
					depth++
					{
						position1766, tokenIndex1766, depth1766 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1767
						}
						position++
						goto l1766
					l1767:
						position, tokenIndex, depth = position1766, tokenIndex1766, depth1766
						if buffer[position] != rune('N') {
							goto l1763
						}
						position++
					}
				l1766:
					{
						position1768, tokenIndex1768, depth1768 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1769
						}
						position++
						goto l1768
					l1769:
						position, tokenIndex, depth = position1768, tokenIndex1768, depth1768
						if buffer[position] != rune('O') {
							goto l1763
						}
						position++
					}
				l1768:
					{
						position1770, tokenIndex1770, depth1770 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1771
						}
						position++
						goto l1770
					l1771:
						position, tokenIndex, depth = position1770, tokenIndex1770, depth1770
						if buffer[position] != rune('T') {
							goto l1763
						}
						position++
					}
				l1770:
					depth--
					add(rulePegText, position1765)
				}
				if !_rules[ruleAction115]() {
					goto l1763
				}
				depth--
				add(ruleNot, position1764)
			}
			return true
		l1763:
			position, tokenIndex, depth = position1763, tokenIndex1763, depth1763
			return false
		},
		/* 147 Equal <- <(<'='> Action116)> */
		func() bool {
			position1772, tokenIndex1772, depth1772 := position, tokenIndex, depth
			{
				position1773 := position
				depth++
				{
					position1774 := position
					depth++
					if buffer[position] != rune('=') {
						goto l1772
					}
					position++
					depth--
					add(rulePegText, position1774)
				}
				if !_rules[ruleAction116]() {
					goto l1772
				}
				depth--
				add(ruleEqual, position1773)
			}
			return true
		l1772:
			position, tokenIndex, depth = position1772, tokenIndex1772, depth1772
			return false
		},
		/* 148 Less <- <(<'<'> Action117)> */
		func() bool {
			position1775, tokenIndex1775, depth1775 := position, tokenIndex, depth
			{
				position1776 := position
				depth++
				{
					position1777 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1775
					}
					position++
					depth--
					add(rulePegText, position1777)
				}
				if !_rules[ruleAction117]() {
					goto l1775
				}
				depth--
				add(ruleLess, position1776)
			}
			return true
		l1775:
			position, tokenIndex, depth = position1775, tokenIndex1775, depth1775
			return false
		},
		/* 149 LessOrEqual <- <(<('<' '=')> Action118)> */
		func() bool {
			position1778, tokenIndex1778, depth1778 := position, tokenIndex, depth
			{
				position1779 := position
				depth++
				{
					position1780 := position
					depth++
					if buffer[position] != rune('<') {
						goto l1778
					}
					position++
					if buffer[position] != rune('=') {
						goto l1778
					}
					position++
					depth--
					add(rulePegText, position1780)
				}
				if !_rules[ruleAction118]() {
					goto l1778
				}
				depth--
				add(ruleLessOrEqual, position1779)
			}
			return true
		l1778:
			position, tokenIndex, depth = position1778, tokenIndex1778, depth1778
			return false
		},
		/* 150 Greater <- <(<'>'> Action119)> */
		func() bool {
			position1781, tokenIndex1781, depth1781 := position, tokenIndex, depth
			{
				position1782 := position
				depth++
				{
					position1783 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1781
					}
					position++
					depth--
					add(rulePegText, position1783)
				}
				if !_rules[ruleAction119]() {
					goto l1781
				}
				depth--
				add(ruleGreater, position1782)
			}
			return true
		l1781:
			position, tokenIndex, depth = position1781, tokenIndex1781, depth1781
			return false
		},
		/* 151 GreaterOrEqual <- <(<('>' '=')> Action120)> */
		func() bool {
			position1784, tokenIndex1784, depth1784 := position, tokenIndex, depth
			{
				position1785 := position
				depth++
				{
					position1786 := position
					depth++
					if buffer[position] != rune('>') {
						goto l1784
					}
					position++
					if buffer[position] != rune('=') {
						goto l1784
					}
					position++
					depth--
					add(rulePegText, position1786)
				}
				if !_rules[ruleAction120]() {
					goto l1784
				}
				depth--
				add(ruleGreaterOrEqual, position1785)
			}
			return true
		l1784:
			position, tokenIndex, depth = position1784, tokenIndex1784, depth1784
			return false
		},
		/* 152 NotEqual <- <(<(('!' '=') / ('<' '>'))> Action121)> */
		func() bool {
			position1787, tokenIndex1787, depth1787 := position, tokenIndex, depth
			{
				position1788 := position
				depth++
				{
					position1789 := position
					depth++
					{
						position1790, tokenIndex1790, depth1790 := position, tokenIndex, depth
						if buffer[position] != rune('!') {
							goto l1791
						}
						position++
						if buffer[position] != rune('=') {
							goto l1791
						}
						position++
						goto l1790
					l1791:
						position, tokenIndex, depth = position1790, tokenIndex1790, depth1790
						if buffer[position] != rune('<') {
							goto l1787
						}
						position++
						if buffer[position] != rune('>') {
							goto l1787
						}
						position++
					}
				l1790:
					depth--
					add(rulePegText, position1789)
				}
				if !_rules[ruleAction121]() {
					goto l1787
				}
				depth--
				add(ruleNotEqual, position1788)
			}
			return true
		l1787:
			position, tokenIndex, depth = position1787, tokenIndex1787, depth1787
			return false
		},
		/* 153 Concat <- <(<('|' '|')> Action122)> */
		func() bool {
			position1792, tokenIndex1792, depth1792 := position, tokenIndex, depth
			{
				position1793 := position
				depth++
				{
					position1794 := position
					depth++
					if buffer[position] != rune('|') {
						goto l1792
					}
					position++
					if buffer[position] != rune('|') {
						goto l1792
					}
					position++
					depth--
					add(rulePegText, position1794)
				}
				if !_rules[ruleAction122]() {
					goto l1792
				}
				depth--
				add(ruleConcat, position1793)
			}
			return true
		l1792:
			position, tokenIndex, depth = position1792, tokenIndex1792, depth1792
			return false
		},
		/* 154 Is <- <(<(('i' / 'I') ('s' / 'S'))> Action123)> */
		func() bool {
			position1795, tokenIndex1795, depth1795 := position, tokenIndex, depth
			{
				position1796 := position
				depth++
				{
					position1797 := position
					depth++
					{
						position1798, tokenIndex1798, depth1798 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1799
						}
						position++
						goto l1798
					l1799:
						position, tokenIndex, depth = position1798, tokenIndex1798, depth1798
						if buffer[position] != rune('I') {
							goto l1795
						}
						position++
					}
				l1798:
					{
						position1800, tokenIndex1800, depth1800 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1801
						}
						position++
						goto l1800
					l1801:
						position, tokenIndex, depth = position1800, tokenIndex1800, depth1800
						if buffer[position] != rune('S') {
							goto l1795
						}
						position++
					}
				l1800:
					depth--
					add(rulePegText, position1797)
				}
				if !_rules[ruleAction123]() {
					goto l1795
				}
				depth--
				add(ruleIs, position1796)
			}
			return true
		l1795:
			position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
			return false
		},
		/* 155 IsNot <- <(<(('i' / 'I') ('s' / 'S') sp (('n' / 'N') ('o' / 'O') ('t' / 'T')))> Action124)> */
		func() bool {
			position1802, tokenIndex1802, depth1802 := position, tokenIndex, depth
			{
				position1803 := position
				depth++
				{
					position1804 := position
					depth++
					{
						position1805, tokenIndex1805, depth1805 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l1806
						}
						position++
						goto l1805
					l1806:
						position, tokenIndex, depth = position1805, tokenIndex1805, depth1805
						if buffer[position] != rune('I') {
							goto l1802
						}
						position++
					}
				l1805:
					{
						position1807, tokenIndex1807, depth1807 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l1808
						}
						position++
						goto l1807
					l1808:
						position, tokenIndex, depth = position1807, tokenIndex1807, depth1807
						if buffer[position] != rune('S') {
							goto l1802
						}
						position++
					}
				l1807:
					if !_rules[rulesp]() {
						goto l1802
					}
					{
						position1809, tokenIndex1809, depth1809 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l1810
						}
						position++
						goto l1809
					l1810:
						position, tokenIndex, depth = position1809, tokenIndex1809, depth1809
						if buffer[position] != rune('N') {
							goto l1802
						}
						position++
					}
				l1809:
					{
						position1811, tokenIndex1811, depth1811 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l1812
						}
						position++
						goto l1811
					l1812:
						position, tokenIndex, depth = position1811, tokenIndex1811, depth1811
						if buffer[position] != rune('O') {
							goto l1802
						}
						position++
					}
				l1811:
					{
						position1813, tokenIndex1813, depth1813 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l1814
						}
						position++
						goto l1813
					l1814:
						position, tokenIndex, depth = position1813, tokenIndex1813, depth1813
						if buffer[position] != rune('T') {
							goto l1802
						}
						position++
					}
				l1813:
					depth--
					add(rulePegText, position1804)
				}
				if !_rules[ruleAction124]() {
					goto l1802
				}
				depth--
				add(ruleIsNot, position1803)
			}
			return true
		l1802:
			position, tokenIndex, depth = position1802, tokenIndex1802, depth1802
			return false
		},
		/* 156 Plus <- <(<'+'> Action125)> */
		func() bool {
			position1815, tokenIndex1815, depth1815 := position, tokenIndex, depth
			{
				position1816 := position
				depth++
				{
					position1817 := position
					depth++
					if buffer[position] != rune('+') {
						goto l1815
					}
					position++
					depth--
					add(rulePegText, position1817)
				}
				if !_rules[ruleAction125]() {
					goto l1815
				}
				depth--
				add(rulePlus, position1816)
			}
			return true
		l1815:
			position, tokenIndex, depth = position1815, tokenIndex1815, depth1815
			return false
		},
		/* 157 Minus <- <(<'-'> Action126)> */
		func() bool {
			position1818, tokenIndex1818, depth1818 := position, tokenIndex, depth
			{
				position1819 := position
				depth++
				{
					position1820 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1818
					}
					position++
					depth--
					add(rulePegText, position1820)
				}
				if !_rules[ruleAction126]() {
					goto l1818
				}
				depth--
				add(ruleMinus, position1819)
			}
			return true
		l1818:
			position, tokenIndex, depth = position1818, tokenIndex1818, depth1818
			return false
		},
		/* 158 Multiply <- <(<'*'> Action127)> */
		func() bool {
			position1821, tokenIndex1821, depth1821 := position, tokenIndex, depth
			{
				position1822 := position
				depth++
				{
					position1823 := position
					depth++
					if buffer[position] != rune('*') {
						goto l1821
					}
					position++
					depth--
					add(rulePegText, position1823)
				}
				if !_rules[ruleAction127]() {
					goto l1821
				}
				depth--
				add(ruleMultiply, position1822)
			}
			return true
		l1821:
			position, tokenIndex, depth = position1821, tokenIndex1821, depth1821
			return false
		},
		/* 159 Divide <- <(<'/'> Action128)> */
		func() bool {
			position1824, tokenIndex1824, depth1824 := position, tokenIndex, depth
			{
				position1825 := position
				depth++
				{
					position1826 := position
					depth++
					if buffer[position] != rune('/') {
						goto l1824
					}
					position++
					depth--
					add(rulePegText, position1826)
				}
				if !_rules[ruleAction128]() {
					goto l1824
				}
				depth--
				add(ruleDivide, position1825)
			}
			return true
		l1824:
			position, tokenIndex, depth = position1824, tokenIndex1824, depth1824
			return false
		},
		/* 160 Modulo <- <(<'%'> Action129)> */
		func() bool {
			position1827, tokenIndex1827, depth1827 := position, tokenIndex, depth
			{
				position1828 := position
				depth++
				{
					position1829 := position
					depth++
					if buffer[position] != rune('%') {
						goto l1827
					}
					position++
					depth--
					add(rulePegText, position1829)
				}
				if !_rules[ruleAction129]() {
					goto l1827
				}
				depth--
				add(ruleModulo, position1828)
			}
			return true
		l1827:
			position, tokenIndex, depth = position1827, tokenIndex1827, depth1827
			return false
		},
		/* 161 UnaryMinus <- <(<'-'> Action130)> */
		func() bool {
			position1830, tokenIndex1830, depth1830 := position, tokenIndex, depth
			{
				position1831 := position
				depth++
				{
					position1832 := position
					depth++
					if buffer[position] != rune('-') {
						goto l1830
					}
					position++
					depth--
					add(rulePegText, position1832)
				}
				if !_rules[ruleAction130]() {
					goto l1830
				}
				depth--
				add(ruleUnaryMinus, position1831)
			}
			return true
		l1830:
			position, tokenIndex, depth = position1830, tokenIndex1830, depth1830
			return false
		},
		/* 162 Identifier <- <(<ident> Action131)> */
		func() bool {
			position1833, tokenIndex1833, depth1833 := position, tokenIndex, depth
			{
				position1834 := position
				depth++
				{
					position1835 := position
					depth++
					if !_rules[ruleident]() {
						goto l1833
					}
					depth--
					add(rulePegText, position1835)
				}
				if !_rules[ruleAction131]() {
					goto l1833
				}
				depth--
				add(ruleIdentifier, position1834)
			}
			return true
		l1833:
			position, tokenIndex, depth = position1833, tokenIndex1833, depth1833
			return false
		},
		/* 163 TargetIdentifier <- <(<('*' / jsonSetPath)> Action132)> */
		func() bool {
			position1836, tokenIndex1836, depth1836 := position, tokenIndex, depth
			{
				position1837 := position
				depth++
				{
					position1838 := position
					depth++
					{
						position1839, tokenIndex1839, depth1839 := position, tokenIndex, depth
						if buffer[position] != rune('*') {
							goto l1840
						}
						position++
						goto l1839
					l1840:
						position, tokenIndex, depth = position1839, tokenIndex1839, depth1839
						if !_rules[rulejsonSetPath]() {
							goto l1836
						}
					}
				l1839:
					depth--
					add(rulePegText, position1838)
				}
				if !_rules[ruleAction132]() {
					goto l1836
				}
				depth--
				add(ruleTargetIdentifier, position1837)
			}
			return true
		l1836:
			position, tokenIndex, depth = position1836, tokenIndex1836, depth1836
			return false
		},
		/* 164 ident <- <(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)> */
		func() bool {
			position1841, tokenIndex1841, depth1841 := position, tokenIndex, depth
			{
				position1842 := position
				depth++
				{
					position1843, tokenIndex1843, depth1843 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l1844
					}
					position++
					goto l1843
				l1844:
					position, tokenIndex, depth = position1843, tokenIndex1843, depth1843
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l1841
					}
					position++
				}
			l1843:
			l1845:
				{
					position1846, tokenIndex1846, depth1846 := position, tokenIndex, depth
					{
						position1847, tokenIndex1847, depth1847 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1848
						}
						position++
						goto l1847
					l1848:
						position, tokenIndex, depth = position1847, tokenIndex1847, depth1847
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1849
						}
						position++
						goto l1847
					l1849:
						position, tokenIndex, depth = position1847, tokenIndex1847, depth1847
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1850
						}
						position++
						goto l1847
					l1850:
						position, tokenIndex, depth = position1847, tokenIndex1847, depth1847
						if buffer[position] != rune('_') {
							goto l1846
						}
						position++
					}
				l1847:
					goto l1845
				l1846:
					position, tokenIndex, depth = position1846, tokenIndex1846, depth1846
				}
				depth--
				add(ruleident, position1842)
			}
			return true
		l1841:
			position, tokenIndex, depth = position1841, tokenIndex1841, depth1841
			return false
		},
		/* 165 jsonGetPath <- <(jsonPathHead jsonGetPathNonHead*)> */
		func() bool {
			position1851, tokenIndex1851, depth1851 := position, tokenIndex, depth
			{
				position1852 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1851
				}
			l1853:
				{
					position1854, tokenIndex1854, depth1854 := position, tokenIndex, depth
					if !_rules[rulejsonGetPathNonHead]() {
						goto l1854
					}
					goto l1853
				l1854:
					position, tokenIndex, depth = position1854, tokenIndex1854, depth1854
				}
				depth--
				add(rulejsonGetPath, position1852)
			}
			return true
		l1851:
			position, tokenIndex, depth = position1851, tokenIndex1851, depth1851
			return false
		},
		/* 166 jsonSetPath <- <(jsonPathHead jsonSetPathNonHead*)> */
		func() bool {
			position1855, tokenIndex1855, depth1855 := position, tokenIndex, depth
			{
				position1856 := position
				depth++
				if !_rules[rulejsonPathHead]() {
					goto l1855
				}
			l1857:
				{
					position1858, tokenIndex1858, depth1858 := position, tokenIndex, depth
					if !_rules[rulejsonSetPathNonHead]() {
						goto l1858
					}
					goto l1857
				l1858:
					position, tokenIndex, depth = position1858, tokenIndex1858, depth1858
				}
				depth--
				add(rulejsonSetPath, position1856)
			}
			return true
		l1855:
			position, tokenIndex, depth = position1855, tokenIndex1855, depth1855
			return false
		},
		/* 167 jsonPathHead <- <(jsonMapAccessString / jsonMapAccessBracket)> */
		func() bool {
			position1859, tokenIndex1859, depth1859 := position, tokenIndex, depth
			{
				position1860 := position
				depth++
				{
					position1861, tokenIndex1861, depth1861 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1862
					}
					goto l1861
				l1862:
					position, tokenIndex, depth = position1861, tokenIndex1861, depth1861
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1859
					}
				}
			l1861:
				depth--
				add(rulejsonPathHead, position1860)
			}
			return true
		l1859:
			position, tokenIndex, depth = position1859, tokenIndex1859, depth1859
			return false
		},
		/* 168 jsonGetPathNonHead <- <(jsonMapMultipleLevel / jsonMapSingleLevel / jsonArrayFullSlice / jsonArrayPartialSlice / jsonArraySlice / jsonArrayAccess)> */
		func() bool {
			position1863, tokenIndex1863, depth1863 := position, tokenIndex, depth
			{
				position1864 := position
				depth++
				{
					position1865, tokenIndex1865, depth1865 := position, tokenIndex, depth
					if !_rules[rulejsonMapMultipleLevel]() {
						goto l1866
					}
					goto l1865
				l1866:
					position, tokenIndex, depth = position1865, tokenIndex1865, depth1865
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1867
					}
					goto l1865
				l1867:
					position, tokenIndex, depth = position1865, tokenIndex1865, depth1865
					if !_rules[rulejsonArrayFullSlice]() {
						goto l1868
					}
					goto l1865
				l1868:
					position, tokenIndex, depth = position1865, tokenIndex1865, depth1865
					if !_rules[rulejsonArrayPartialSlice]() {
						goto l1869
					}
					goto l1865
				l1869:
					position, tokenIndex, depth = position1865, tokenIndex1865, depth1865
					if !_rules[rulejsonArraySlice]() {
						goto l1870
					}
					goto l1865
				l1870:
					position, tokenIndex, depth = position1865, tokenIndex1865, depth1865
					if !_rules[rulejsonArrayAccess]() {
						goto l1863
					}
				}
			l1865:
				depth--
				add(rulejsonGetPathNonHead, position1864)
			}
			return true
		l1863:
			position, tokenIndex, depth = position1863, tokenIndex1863, depth1863
			return false
		},
		/* 169 jsonSetPathNonHead <- <(jsonMapSingleLevel / jsonNonNegativeArrayAccess)> */
		func() bool {
			position1871, tokenIndex1871, depth1871 := position, tokenIndex, depth
			{
				position1872 := position
				depth++
				{
					position1873, tokenIndex1873, depth1873 := position, tokenIndex, depth
					if !_rules[rulejsonMapSingleLevel]() {
						goto l1874
					}
					goto l1873
				l1874:
					position, tokenIndex, depth = position1873, tokenIndex1873, depth1873
					if !_rules[rulejsonNonNegativeArrayAccess]() {
						goto l1871
					}
				}
			l1873:
				depth--
				add(rulejsonSetPathNonHead, position1872)
			}
			return true
		l1871:
			position, tokenIndex, depth = position1871, tokenIndex1871, depth1871
			return false
		},
		/* 170 jsonMapSingleLevel <- <(('.' jsonMapAccessString) / jsonMapAccessBracket)> */
		func() bool {
			position1875, tokenIndex1875, depth1875 := position, tokenIndex, depth
			{
				position1876 := position
				depth++
				{
					position1877, tokenIndex1877, depth1877 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1878
					}
					position++
					if !_rules[rulejsonMapAccessString]() {
						goto l1878
					}
					goto l1877
				l1878:
					position, tokenIndex, depth = position1877, tokenIndex1877, depth1877
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1875
					}
				}
			l1877:
				depth--
				add(rulejsonMapSingleLevel, position1876)
			}
			return true
		l1875:
			position, tokenIndex, depth = position1875, tokenIndex1875, depth1875
			return false
		},
		/* 171 jsonMapMultipleLevel <- <('.' '.' (jsonMapAccessString / jsonMapAccessBracket))> */
		func() bool {
			position1879, tokenIndex1879, depth1879 := position, tokenIndex, depth
			{
				position1880 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1879
				}
				position++
				if buffer[position] != rune('.') {
					goto l1879
				}
				position++
				{
					position1881, tokenIndex1881, depth1881 := position, tokenIndex, depth
					if !_rules[rulejsonMapAccessString]() {
						goto l1882
					}
					goto l1881
				l1882:
					position, tokenIndex, depth = position1881, tokenIndex1881, depth1881
					if !_rules[rulejsonMapAccessBracket]() {
						goto l1879
					}
				}
			l1881:
				depth--
				add(rulejsonMapMultipleLevel, position1880)
			}
			return true
		l1879:
			position, tokenIndex, depth = position1879, tokenIndex1879, depth1879
			return false
		},
		/* 172 jsonMapAccessString <- <<(([a-z] / [A-Z]) ([a-z] / [A-Z] / [0-9] / '_')*)>> */
		func() bool {
			position1883, tokenIndex1883, depth1883 := position, tokenIndex, depth
			{
				position1884 := position
				depth++
				{
					position1885 := position
					depth++
					{
						position1886, tokenIndex1886, depth1886 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l1887
						}
						position++
						goto l1886
					l1887:
						position, tokenIndex, depth = position1886, tokenIndex1886, depth1886
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1883
						}
						position++
					}
				l1886:
				l1888:
					{
						position1889, tokenIndex1889, depth1889 := position, tokenIndex, depth
						{
							position1890, tokenIndex1890, depth1890 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1891
							}
							position++
							goto l1890
						l1891:
							position, tokenIndex, depth = position1890, tokenIndex1890, depth1890
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1892
							}
							position++
							goto l1890
						l1892:
							position, tokenIndex, depth = position1890, tokenIndex1890, depth1890
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1893
							}
							position++
							goto l1890
						l1893:
							position, tokenIndex, depth = position1890, tokenIndex1890, depth1890
							if buffer[position] != rune('_') {
								goto l1889
							}
							position++
						}
					l1890:
						goto l1888
					l1889:
						position, tokenIndex, depth = position1889, tokenIndex1889, depth1889
					}
					depth--
					add(rulePegText, position1885)
				}
				depth--
				add(rulejsonMapAccessString, position1884)
			}
			return true
		l1883:
			position, tokenIndex, depth = position1883, tokenIndex1883, depth1883
			return false
		},
		/* 173 jsonMapAccessBracket <- <('[' doubleQuotedString ']')> */
		func() bool {
			position1894, tokenIndex1894, depth1894 := position, tokenIndex, depth
			{
				position1895 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1894
				}
				position++
				if !_rules[ruledoubleQuotedString]() {
					goto l1894
				}
				if buffer[position] != rune(']') {
					goto l1894
				}
				position++
				depth--
				add(rulejsonMapAccessBracket, position1895)
			}
			return true
		l1894:
			position, tokenIndex, depth = position1894, tokenIndex1894, depth1894
			return false
		},
		/* 174 doubleQuotedString <- <('"' <(('"' '"') / (!'"' .))*> '"')> */
		func() bool {
			position1896, tokenIndex1896, depth1896 := position, tokenIndex, depth
			{
				position1897 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1896
				}
				position++
				{
					position1898 := position
					depth++
				l1899:
					{
						position1900, tokenIndex1900, depth1900 := position, tokenIndex, depth
						{
							position1901, tokenIndex1901, depth1901 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l1902
							}
							position++
							if buffer[position] != rune('"') {
								goto l1902
							}
							position++
							goto l1901
						l1902:
							position, tokenIndex, depth = position1901, tokenIndex1901, depth1901
							{
								position1903, tokenIndex1903, depth1903 := position, tokenIndex, depth
								if buffer[position] != rune('"') {
									goto l1903
								}
								position++
								goto l1900
							l1903:
								position, tokenIndex, depth = position1903, tokenIndex1903, depth1903
							}
							if !matchDot() {
								goto l1900
							}
						}
					l1901:
						goto l1899
					l1900:
						position, tokenIndex, depth = position1900, tokenIndex1900, depth1900
					}
					depth--
					add(rulePegText, position1898)
				}
				if buffer[position] != rune('"') {
					goto l1896
				}
				position++
				depth--
				add(ruledoubleQuotedString, position1897)
			}
			return true
		l1896:
			position, tokenIndex, depth = position1896, tokenIndex1896, depth1896
			return false
		},
		/* 175 jsonArrayAccess <- <('[' <('-'? [0-9]+)> ']')> */
		func() bool {
			position1904, tokenIndex1904, depth1904 := position, tokenIndex, depth
			{
				position1905 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1904
				}
				position++
				{
					position1906 := position
					depth++
					{
						position1907, tokenIndex1907, depth1907 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1907
						}
						position++
						goto l1908
					l1907:
						position, tokenIndex, depth = position1907, tokenIndex1907, depth1907
					}
				l1908:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1904
					}
					position++
				l1909:
					{
						position1910, tokenIndex1910, depth1910 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1910
						}
						position++
						goto l1909
					l1910:
						position, tokenIndex, depth = position1910, tokenIndex1910, depth1910
					}
					depth--
					add(rulePegText, position1906)
				}
				if buffer[position] != rune(']') {
					goto l1904
				}
				position++
				depth--
				add(rulejsonArrayAccess, position1905)
			}
			return true
		l1904:
			position, tokenIndex, depth = position1904, tokenIndex1904, depth1904
			return false
		},
		/* 176 jsonNonNegativeArrayAccess <- <('[' <[0-9]+> ']')> */
		func() bool {
			position1911, tokenIndex1911, depth1911 := position, tokenIndex, depth
			{
				position1912 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1911
				}
				position++
				{
					position1913 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1911
					}
					position++
				l1914:
					{
						position1915, tokenIndex1915, depth1915 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1915
						}
						position++
						goto l1914
					l1915:
						position, tokenIndex, depth = position1915, tokenIndex1915, depth1915
					}
					depth--
					add(rulePegText, position1913)
				}
				if buffer[position] != rune(']') {
					goto l1911
				}
				position++
				depth--
				add(rulejsonNonNegativeArrayAccess, position1912)
			}
			return true
		l1911:
			position, tokenIndex, depth = position1911, tokenIndex1911, depth1911
			return false
		},
		/* 177 jsonArraySlice <- <('[' <('-'? [0-9]+ ':' '-'? [0-9]+ (':' '-'? [0-9]+)?)> ']')> */
		func() bool {
			position1916, tokenIndex1916, depth1916 := position, tokenIndex, depth
			{
				position1917 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1916
				}
				position++
				{
					position1918 := position
					depth++
					{
						position1919, tokenIndex1919, depth1919 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1919
						}
						position++
						goto l1920
					l1919:
						position, tokenIndex, depth = position1919, tokenIndex1919, depth1919
					}
				l1920:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1916
					}
					position++
				l1921:
					{
						position1922, tokenIndex1922, depth1922 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1922
						}
						position++
						goto l1921
					l1922:
						position, tokenIndex, depth = position1922, tokenIndex1922, depth1922
					}
					if buffer[position] != rune(':') {
						goto l1916
					}
					position++
					{
						position1923, tokenIndex1923, depth1923 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l1923
						}
						position++
						goto l1924
					l1923:
						position, tokenIndex, depth = position1923, tokenIndex1923, depth1923
					}
				l1924:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1916
					}
					position++
				l1925:
					{
						position1926, tokenIndex1926, depth1926 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1926
						}
						position++
						goto l1925
					l1926:
						position, tokenIndex, depth = position1926, tokenIndex1926, depth1926
					}
					{
						position1927, tokenIndex1927, depth1927 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1927
						}
						position++
						{
							position1929, tokenIndex1929, depth1929 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1929
							}
							position++
							goto l1930
						l1929:
							position, tokenIndex, depth = position1929, tokenIndex1929, depth1929
						}
					l1930:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1927
						}
						position++
					l1931:
						{
							position1932, tokenIndex1932, depth1932 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1932
							}
							position++
							goto l1931
						l1932:
							position, tokenIndex, depth = position1932, tokenIndex1932, depth1932
						}
						goto l1928
					l1927:
						position, tokenIndex, depth = position1927, tokenIndex1927, depth1927
					}
				l1928:
					depth--
					add(rulePegText, position1918)
				}
				if buffer[position] != rune(']') {
					goto l1916
				}
				position++
				depth--
				add(rulejsonArraySlice, position1917)
			}
			return true
		l1916:
			position, tokenIndex, depth = position1916, tokenIndex1916, depth1916
			return false
		},
		/* 178 jsonArrayPartialSlice <- <('[' <((':' '-'? [0-9]+) / ('-'? [0-9]+ ':'))> ']')> */
		func() bool {
			position1933, tokenIndex1933, depth1933 := position, tokenIndex, depth
			{
				position1934 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1933
				}
				position++
				{
					position1935 := position
					depth++
					{
						position1936, tokenIndex1936, depth1936 := position, tokenIndex, depth
						if buffer[position] != rune(':') {
							goto l1937
						}
						position++
						{
							position1938, tokenIndex1938, depth1938 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1938
							}
							position++
							goto l1939
						l1938:
							position, tokenIndex, depth = position1938, tokenIndex1938, depth1938
						}
					l1939:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1937
						}
						position++
					l1940:
						{
							position1941, tokenIndex1941, depth1941 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1941
							}
							position++
							goto l1940
						l1941:
							position, tokenIndex, depth = position1941, tokenIndex1941, depth1941
						}
						goto l1936
					l1937:
						position, tokenIndex, depth = position1936, tokenIndex1936, depth1936
						{
							position1942, tokenIndex1942, depth1942 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1942
							}
							position++
							goto l1943
						l1942:
							position, tokenIndex, depth = position1942, tokenIndex1942, depth1942
						}
					l1943:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1933
						}
						position++
					l1944:
						{
							position1945, tokenIndex1945, depth1945 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1945
							}
							position++
							goto l1944
						l1945:
							position, tokenIndex, depth = position1945, tokenIndex1945, depth1945
						}
						if buffer[position] != rune(':') {
							goto l1933
						}
						position++
					}
				l1936:
					depth--
					add(rulePegText, position1935)
				}
				if buffer[position] != rune(']') {
					goto l1933
				}
				position++
				depth--
				add(rulejsonArrayPartialSlice, position1934)
			}
			return true
		l1933:
			position, tokenIndex, depth = position1933, tokenIndex1933, depth1933
			return false
		},
		/* 179 jsonArrayFullSlice <- <('[' ':' ']')> */
		func() bool {
			position1946, tokenIndex1946, depth1946 := position, tokenIndex, depth
			{
				position1947 := position
				depth++
				if buffer[position] != rune('[') {
					goto l1946
				}
				position++
				if buffer[position] != rune(':') {
					goto l1946
				}
				position++
				if buffer[position] != rune(']') {
					goto l1946
				}
				position++
				depth--
				add(rulejsonArrayFullSlice, position1947)
			}
			return true
		l1946:
			position, tokenIndex, depth = position1946, tokenIndex1946, depth1946
			return false
		},
		/* 180 spElem <- <(' ' / '\t' / '\n' / '\r' / comment / finalComment)> */
		func() bool {
			position1948, tokenIndex1948, depth1948 := position, tokenIndex, depth
			{
				position1949 := position
				depth++
				{
					position1950, tokenIndex1950, depth1950 := position, tokenIndex, depth
					if buffer[position] != rune(' ') {
						goto l1951
					}
					position++
					goto l1950
				l1951:
					position, tokenIndex, depth = position1950, tokenIndex1950, depth1950
					if buffer[position] != rune('\t') {
						goto l1952
					}
					position++
					goto l1950
				l1952:
					position, tokenIndex, depth = position1950, tokenIndex1950, depth1950
					if buffer[position] != rune('\n') {
						goto l1953
					}
					position++
					goto l1950
				l1953:
					position, tokenIndex, depth = position1950, tokenIndex1950, depth1950
					if buffer[position] != rune('\r') {
						goto l1954
					}
					position++
					goto l1950
				l1954:
					position, tokenIndex, depth = position1950, tokenIndex1950, depth1950
					if !_rules[rulecomment]() {
						goto l1955
					}
					goto l1950
				l1955:
					position, tokenIndex, depth = position1950, tokenIndex1950, depth1950
					if !_rules[rulefinalComment]() {
						goto l1948
					}
				}
			l1950:
				depth--
				add(rulespElem, position1949)
			}
			return true
		l1948:
			position, tokenIndex, depth = position1948, tokenIndex1948, depth1948
			return false
		},
		/* 181 sp <- <spElem+> */
		func() bool {
			position1956, tokenIndex1956, depth1956 := position, tokenIndex, depth
			{
				position1957 := position
				depth++
				if !_rules[rulespElem]() {
					goto l1956
				}
			l1958:
				{
					position1959, tokenIndex1959, depth1959 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1959
					}
					goto l1958
				l1959:
					position, tokenIndex, depth = position1959, tokenIndex1959, depth1959
				}
				depth--
				add(rulesp, position1957)
			}
			return true
		l1956:
			position, tokenIndex, depth = position1956, tokenIndex1956, depth1956
			return false
		},
		/* 182 spOpt <- <spElem*> */
		func() bool {
			{
				position1961 := position
				depth++
			l1962:
				{
					position1963, tokenIndex1963, depth1963 := position, tokenIndex, depth
					if !_rules[rulespElem]() {
						goto l1963
					}
					goto l1962
				l1963:
					position, tokenIndex, depth = position1963, tokenIndex1963, depth1963
				}
				depth--
				add(rulespOpt, position1961)
			}
			return true
		},
		/* 183 comment <- <('-' '-' (!('\r' / '\n') .)* ('\r' / '\n'))> */
		func() bool {
			position1964, tokenIndex1964, depth1964 := position, tokenIndex, depth
			{
				position1965 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1964
				}
				position++
				if buffer[position] != rune('-') {
					goto l1964
				}
				position++
			l1966:
				{
					position1967, tokenIndex1967, depth1967 := position, tokenIndex, depth
					{
						position1968, tokenIndex1968, depth1968 := position, tokenIndex, depth
						{
							position1969, tokenIndex1969, depth1969 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1970
							}
							position++
							goto l1969
						l1970:
							position, tokenIndex, depth = position1969, tokenIndex1969, depth1969
							if buffer[position] != rune('\n') {
								goto l1968
							}
							position++
						}
					l1969:
						goto l1967
					l1968:
						position, tokenIndex, depth = position1968, tokenIndex1968, depth1968
					}
					if !matchDot() {
						goto l1967
					}
					goto l1966
				l1967:
					position, tokenIndex, depth = position1967, tokenIndex1967, depth1967
				}
				{
					position1971, tokenIndex1971, depth1971 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1972
					}
					position++
					goto l1971
				l1972:
					position, tokenIndex, depth = position1971, tokenIndex1971, depth1971
					if buffer[position] != rune('\n') {
						goto l1964
					}
					position++
				}
			l1971:
				depth--
				add(rulecomment, position1965)
			}
			return true
		l1964:
			position, tokenIndex, depth = position1964, tokenIndex1964, depth1964
			return false
		},
		/* 184 finalComment <- <('-' '-' (!('\r' / '\n') .)* !.)> */
		func() bool {
			position1973, tokenIndex1973, depth1973 := position, tokenIndex, depth
			{
				position1974 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1973
				}
				position++
				if buffer[position] != rune('-') {
					goto l1973
				}
				position++
			l1975:
				{
					position1976, tokenIndex1976, depth1976 := position, tokenIndex, depth
					{
						position1977, tokenIndex1977, depth1977 := position, tokenIndex, depth
						{
							position1978, tokenIndex1978, depth1978 := position, tokenIndex, depth
							if buffer[position] != rune('\r') {
								goto l1979
							}
							position++
							goto l1978
						l1979:
							position, tokenIndex, depth = position1978, tokenIndex1978, depth1978
							if buffer[position] != rune('\n') {
								goto l1977
							}
							position++
						}
					l1978:
						goto l1976
					l1977:
						position, tokenIndex, depth = position1977, tokenIndex1977, depth1977
					}
					if !matchDot() {
						goto l1976
					}
					goto l1975
				l1976:
					position, tokenIndex, depth = position1976, tokenIndex1976, depth1976
				}
				{
					position1980, tokenIndex1980, depth1980 := position, tokenIndex, depth
					if !matchDot() {
						goto l1980
					}
					goto l1973
				l1980:
					position, tokenIndex, depth = position1980, tokenIndex1980, depth1980
				}
				depth--
				add(rulefinalComment, position1974)
			}
			return true
		l1973:
			position, tokenIndex, depth = position1973, tokenIndex1973, depth1973
			return false
		},
		nil,
		/* 187 Action0 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 188 Action1 <- <{
		    p.IncludeTrailingWhitespace(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 189 Action2 <- <{
		    p.AssembleSelect()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 190 Action3 <- <{
		    p.AssembleSelectUnion(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 191 Action4 <- <{
		    p.AssembleCreateStreamAsSelect()
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 192 Action5 <- <{
		    p.AssembleCreateStreamAsSelectUnion()
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 193 Action6 <- <{
		    p.AssembleCreateSource()
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 194 Action7 <- <{
		    p.AssembleCreateSink()
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 195 Action8 <- <{
		    p.AssembleCreateState()
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 196 Action9 <- <{
		    p.AssembleUpdateState()
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 197 Action10 <- <{
		    p.AssembleUpdateSource()
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 198 Action11 <- <{
		    p.AssembleUpdateSink()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 199 Action12 <- <{
		    p.AssembleInsertIntoFrom()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 200 Action13 <- <{
		    p.AssemblePauseSource()
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 201 Action14 <- <{
		    p.AssembleResumeSource()
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 202 Action15 <- <{
		    p.AssembleRewindSource()
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 203 Action16 <- <{
		    p.AssembleDropSource()
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 204 Action17 <- <{
		    p.AssembleDropStream()
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 205 Action18 <- <{
		    p.AssembleDropSink()
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 206 Action19 <- <{
		    p.AssembleDropState()
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 207 Action20 <- <{
		    p.AssembleLoadState()
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 208 Action21 <- <{
		    p.AssembleLoadStateOrCreate()
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 209 Action22 <- <{
		    p.AssembleSaveState()
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 210 Action23 <- <{
		    p.AssembleEval(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 211 Action24 <- <{
		    p.AssembleEmitter()
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 212 Action25 <- <{
		    p.AssembleEmitterOptions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 213 Action26 <- <{
		    p.AssembleEmitterLimit()
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 214 Action27 <- <{
		    p.AssembleEmitterSampling(CountBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 215 Action28 <- <{
		    p.AssembleEmitterSampling(RandomizedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 216 Action29 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 1)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 217 Action30 <- <{
		    p.AssembleEmitterSampling(TimeBasedSampling, 0.001)
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 218 Action31 <- <{
		    p.AssembleProjections(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 219 Action32 <- <{
		    p.AssembleAlias()
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 220 Action33 <- <{
		    // This is *always* executed, even if there is no
		    // FROM clause present in the statement.
		    p.AssembleWindowedFrom(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 221 Action34 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 222 Action35 <- <{
		    p.AssembleInterval()
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 223 Action36 <- <{
		    // This is *always* executed, even if there is no
		    // WHERE clause present in the statement.
		    p.AssembleFilter(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 224 Action37 <- <{
		    // This is *always* executed, even if there is no
		    // GROUP BY clause present in the statement.
		    p.AssembleGrouping(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 225 Action38 <- <{
		    // This is *always* executed, even if there is no
		    // HAVING clause present in the statement.
		    p.AssembleHaving(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 226 Action39 <- <{
		    p.EnsureAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 227 Action40 <- <{
		    p.AssembleAliasedStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 228 Action41 <- <{
		    p.AssembleStreamWindow()
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 229 Action42 <- <{
		    p.AssembleUDSFFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 230 Action43 <- <{
		    p.EnsureCapacitySpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 231 Action44 <- <{
		    p.EnsureSheddingSpec(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 232 Action45 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 233 Action46 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 234 Action47 <- <{
		    p.AssembleSourceSinkSpecs(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 235 Action48 <- <{
		    p.EnsureIdentifier(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 236 Action49 <- <{
		    p.AssembleSourceSinkParam()
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
		/* 237 Action50 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction50, position)
			}
			return true
		},
		/* 238 Action51 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction51, position)
			}
			return true
		},
		/* 239 Action52 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction52, position)
			}
			return true
		},
		/* 240 Action53 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction53, position)
			}
			return true
		},
		/* 241 Action54 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction54, position)
			}
			return true
		},
		/* 242 Action55 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction55, position)
			}
			return true
		},
		/* 243 Action56 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction56, position)
			}
			return true
		},
		/* 244 Action57 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction57, position)
			}
			return true
		},
		/* 245 Action58 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction58, position)
			}
			return true
		},
		/* 246 Action59 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction59, position)
			}
			return true
		},
		/* 247 Action60 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction60, position)
			}
			return true
		},
		/* 248 Action61 <- <{
		    p.AssembleBinaryOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction61, position)
			}
			return true
		},
		/* 249 Action62 <- <{
		    p.AssembleUnaryPrefixOperation(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction62, position)
			}
			return true
		},
		/* 250 Action63 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction63, position)
			}
			return true
		},
		/* 251 Action64 <- <{
		    p.AssembleTypeCast(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction64, position)
			}
			return true
		},
		/* 252 Action65 <- <{
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction65, position)
			}
			return true
		},
		/* 253 Action66 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleFuncApp()
		}> */
		func() bool {
			{
				add(ruleAction66, position)
			}
			return true
		},
		/* 254 Action67 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction67, position)
			}
			return true
		},
		/* 255 Action68 <- <{
		    p.AssembleExpressions(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction68, position)
			}
			return true
		},
		/* 256 Action69 <- <{
		    p.AssembleSortedExpression()
		}> */
		func() bool {
			{
				add(ruleAction69, position)
			}
			return true
		},
		/* 257 Action70 <- <{
		    p.EnsureKeywordPresent(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction70, position)
			}
			return true
		},
		/* 258 Action71 <- <{
		    p.AssembleExpressions(begin, end)
		    p.AssembleArray()
		}> */
		func() bool {
			{
				add(ruleAction71, position)
			}
			return true
		},
		/* 259 Action72 <- <{
		    p.AssembleMap(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction72, position)
			}
			return true
		},
		/* 260 Action73 <- <{
		    p.AssembleKeyValuePair()
		}> */
		func() bool {
			{
				add(ruleAction73, position)
			}
			return true
		},
		/* 261 Action74 <- <{
		    p.AssembleConditionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction74, position)
			}
			return true
		},
		/* 262 Action75 <- <{
		    p.AssembleExpressionCase(begin, end)
		}> */
		func() bool {
			{
				add(ruleAction75, position)
			}
			return true
		},
		/* 263 Action76 <- <{
		    p.AssembleWhenThenPair()
		}> */
		func() bool {
			{
				add(ruleAction76, position)
			}
			return true
		},
		/* 264 Action77 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStream(substr))
		}> */
		func() bool {
			{
				add(ruleAction77, position)
			}
			return true
		},
		/* 265 Action78 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowMeta(substr, TimestampMeta))
		}> */
		func() bool {
			{
				add(ruleAction78, position)
			}
			return true
		},
		/* 266 Action79 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewRowValue(substr))
		}> */
		func() bool {
			{
				add(ruleAction79, position)
			}
			return true
		},
		/* 267 Action80 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction80, position)
			}
			return true
		},
		/* 268 Action81 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewNumericLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction81, position)
			}
			return true
		},
		/* 269 Action82 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewFloatLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction82, position)
			}
			return true
		},
		/* 270 Action83 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, FuncName(substr))
		}> */
		func() bool {
			{
				add(ruleAction83, position)
			}
			return true
		},
		/* 271 Action84 <- <{
		    p.PushComponent(begin, end, NewNullLiteral())
		}> */
		func() bool {
			{
				add(ruleAction84, position)
			}
			return true
		},
		/* 272 Action85 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(true))
		}> */
		func() bool {
			{
				add(ruleAction85, position)
			}
			return true
		},
		/* 273 Action86 <- <{
		    p.PushComponent(begin, end, NewBoolLiteral(false))
		}> */
		func() bool {
			{
				add(ruleAction86, position)
			}
			return true
		},
		/* 274 Action87 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewWildcard(substr))
		}> */
		func() bool {
			{
				add(ruleAction87, position)
			}
			return true
		},
		/* 275 Action88 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, NewStringLiteral(substr))
		}> */
		func() bool {
			{
				add(ruleAction88, position)
			}
			return true
		},
		/* 276 Action89 <- <{
		    p.PushComponent(begin, end, Istream)
		}> */
		func() bool {
			{
				add(ruleAction89, position)
			}
			return true
		},
		/* 277 Action90 <- <{
		    p.PushComponent(begin, end, Dstream)
		}> */
		func() bool {
			{
				add(ruleAction90, position)
			}
			return true
		},
		/* 278 Action91 <- <{
		    p.PushComponent(begin, end, Rstream)
		}> */
		func() bool {
			{
				add(ruleAction91, position)
			}
			return true
		},
		/* 279 Action92 <- <{
		    p.PushComponent(begin, end, Tuples)
		}> */
		func() bool {
			{
				add(ruleAction92, position)
			}
			return true
		},
		/* 280 Action93 <- <{
		    p.PushComponent(begin, end, Seconds)
		}> */
		func() bool {
			{
				add(ruleAction93, position)
			}
			return true
		},
		/* 281 Action94 <- <{
		    p.PushComponent(begin, end, Milliseconds)
		}> */
		func() bool {
			{
				add(ruleAction94, position)
			}
			return true
		},
		/* 282 Action95 <- <{
		    p.PushComponent(begin, end, Wait)
		}> */
		func() bool {
			{
				add(ruleAction95, position)
			}
			return true
		},
		/* 283 Action96 <- <{
		    p.PushComponent(begin, end, DropOldest)
		}> */
		func() bool {
			{
				add(ruleAction96, position)
			}
			return true
		},
		/* 284 Action97 <- <{
		    p.PushComponent(begin, end, DropNewest)
		}> */
		func() bool {
			{
				add(ruleAction97, position)
			}
			return true
		},
		/* 285 Action98 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, StreamIdentifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction98, position)
			}
			return true
		},
		/* 286 Action99 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkType(substr))
		}> */
		func() bool {
			{
				add(ruleAction99, position)
			}
			return true
		},
		/* 287 Action100 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, SourceSinkParamKey(substr))
		}> */
		func() bool {
			{
				add(ruleAction100, position)
			}
			return true
		},
		/* 288 Action101 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction101, position)
			}
			return true
		},
		/* 289 Action102 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction102, position)
			}
			return true
		},
		/* 290 Action103 <- <{
		    p.PushComponent(begin, end, Yes)
		}> */
		func() bool {
			{
				add(ruleAction103, position)
			}
			return true
		},
		/* 291 Action104 <- <{
		    p.PushComponent(begin, end, No)
		}> */
		func() bool {
			{
				add(ruleAction104, position)
			}
			return true
		},
		/* 292 Action105 <- <{
		    p.PushComponent(begin, end, Bool)
		}> */
		func() bool {
			{
				add(ruleAction105, position)
			}
			return true
		},
		/* 293 Action106 <- <{
		    p.PushComponent(begin, end, Int)
		}> */
		func() bool {
			{
				add(ruleAction106, position)
			}
			return true
		},
		/* 294 Action107 <- <{
		    p.PushComponent(begin, end, Float)
		}> */
		func() bool {
			{
				add(ruleAction107, position)
			}
			return true
		},
		/* 295 Action108 <- <{
		    p.PushComponent(begin, end, String)
		}> */
		func() bool {
			{
				add(ruleAction108, position)
			}
			return true
		},
		/* 296 Action109 <- <{
		    p.PushComponent(begin, end, Blob)
		}> */
		func() bool {
			{
				add(ruleAction109, position)
			}
			return true
		},
		/* 297 Action110 <- <{
		    p.PushComponent(begin, end, Timestamp)
		}> */
		func() bool {
			{
				add(ruleAction110, position)
			}
			return true
		},
		/* 298 Action111 <- <{
		    p.PushComponent(begin, end, Array)
		}> */
		func() bool {
			{
				add(ruleAction111, position)
			}
			return true
		},
		/* 299 Action112 <- <{
		    p.PushComponent(begin, end, Map)
		}> */
		func() bool {
			{
				add(ruleAction112, position)
			}
			return true
		},
		/* 300 Action113 <- <{
		    p.PushComponent(begin, end, Or)
		}> */
		func() bool {
			{
				add(ruleAction113, position)
			}
			return true
		},
		/* 301 Action114 <- <{
		    p.PushComponent(begin, end, And)
		}> */
		func() bool {
			{
				add(ruleAction114, position)
			}
			return true
		},
		/* 302 Action115 <- <{
		    p.PushComponent(begin, end, Not)
		}> */
		func() bool {
			{
				add(ruleAction115, position)
			}
			return true
		},
		/* 303 Action116 <- <{
		    p.PushComponent(begin, end, Equal)
		}> */
		func() bool {
			{
				add(ruleAction116, position)
			}
			return true
		},
		/* 304 Action117 <- <{
		    p.PushComponent(begin, end, Less)
		}> */
		func() bool {
			{
				add(ruleAction117, position)
			}
			return true
		},
		/* 305 Action118 <- <{
		    p.PushComponent(begin, end, LessOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction118, position)
			}
			return true
		},
		/* 306 Action119 <- <{
		    p.PushComponent(begin, end, Greater)
		}> */
		func() bool {
			{
				add(ruleAction119, position)
			}
			return true
		},
		/* 307 Action120 <- <{
		    p.PushComponent(begin, end, GreaterOrEqual)
		}> */
		func() bool {
			{
				add(ruleAction120, position)
			}
			return true
		},
		/* 308 Action121 <- <{
		    p.PushComponent(begin, end, NotEqual)
		}> */
		func() bool {
			{
				add(ruleAction121, position)
			}
			return true
		},
		/* 309 Action122 <- <{
		    p.PushComponent(begin, end, Concat)
		}> */
		func() bool {
			{
				add(ruleAction122, position)
			}
			return true
		},
		/* 310 Action123 <- <{
		    p.PushComponent(begin, end, Is)
		}> */
		func() bool {
			{
				add(ruleAction123, position)
			}
			return true
		},
		/* 311 Action124 <- <{
		    p.PushComponent(begin, end, IsNot)
		}> */
		func() bool {
			{
				add(ruleAction124, position)
			}
			return true
		},
		/* 312 Action125 <- <{
		    p.PushComponent(begin, end, Plus)
		}> */
		func() bool {
			{
				add(ruleAction125, position)
			}
			return true
		},
		/* 313 Action126 <- <{
		    p.PushComponent(begin, end, Minus)
		}> */
		func() bool {
			{
				add(ruleAction126, position)
			}
			return true
		},
		/* 314 Action127 <- <{
		    p.PushComponent(begin, end, Multiply)
		}> */
		func() bool {
			{
				add(ruleAction127, position)
			}
			return true
		},
		/* 315 Action128 <- <{
		    p.PushComponent(begin, end, Divide)
		}> */
		func() bool {
			{
				add(ruleAction128, position)
			}
			return true
		},
		/* 316 Action129 <- <{
		    p.PushComponent(begin, end, Modulo)
		}> */
		func() bool {
			{
				add(ruleAction129, position)
			}
			return true
		},
		/* 317 Action130 <- <{
		    p.PushComponent(begin, end, UnaryMinus)
		}> */
		func() bool {
			{
				add(ruleAction130, position)
			}
			return true
		},
		/* 318 Action131 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction131, position)
			}
			return true
		},
		/* 319 Action132 <- <{
		    substr := string([]rune(buffer)[begin:end])
		    p.PushComponent(begin, end, Identifier(substr))
		}> */
		func() bool {
			{
				add(ruleAction132, position)
			}
			return true
		},
	}
	p.rules = _rules
}
