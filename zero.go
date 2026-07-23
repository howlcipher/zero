//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenLParen TokenType = "LPAREN"
	TokenRParen TokenType = "RPAREN"
	TokenSymbol TokenType = "SYMBOL"
	TokenInt    TokenType = "INT"
	TokenString TokenType = "STRING"
	TokenEOF    TokenType = "EOF"
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type ErrorOutput struct {
	Reason string `json:"reason"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

func reportError(reason string, line, column int) {
	errOut := ErrorOutput{Reason: reason, Line: line, Column: column}
	b, _ := json.Marshal(errOut)
	fmt.Println(string(b))
	os.Exit(1)
}

// Lexer
type Lexer struct {
	input  string
	pos    int
	line   int
	column int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, line: 1, column: 1}
}

func (l *Lexer) nextChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := rune(l.input[l.pos])
	l.pos++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) peekChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) NextToken() Token {
	for {
		ch := l.peekChar()
		if ch == 0 {
			return Token{Type: TokenEOF, Line: l.line, Column: l.column}
		}
		if unicode.IsSpace(ch) {
			l.nextChar()
			continue
		}
		break
	}

	startLine := l.line
	startCol := l.column
	ch := l.nextChar()

	if ch == '(' {
		return Token{Type: TokenLParen, Value: "(", Line: startLine, Column: startCol}
	}
	if ch == ')' {
		return Token{Type: TokenRParen, Value: ")", Line: startLine, Column: startCol}
	}
	if ch == '"' {
		raw := "\""
		for {
			nextCh := l.nextChar()
			if nextCh == 0 {
				reportError("Unterminated string", startLine, startCol)
			}
			if nextCh == '\\' {
				escapedCh := l.nextChar()
				if escapedCh == 0 {
					reportError("Unterminated string escape", startLine, startCol)
				}
				if escapedCh == '\'' {
					raw += "'"
				} else {
					raw += "\\" + string(escapedCh)
				}
				continue
			}
			raw += string(nextCh)
			if nextCh == '"' {
				break
			}
		}
		val, err := strconv.Unquote(raw)
		if err != nil {
			reportError("Invalid string literal: "+err.Error(), startLine, startCol)
		}
		return Token{Type: TokenString, Value: val, Line: startLine, Column: startCol}
	}
	if unicode.IsDigit(ch) {
		val := string(ch)
		for unicode.IsDigit(l.peekChar()) {
			val += string(l.nextChar())
		}
		return Token{Type: TokenInt, Value: val, Line: startLine, Column: startCol}
	}
	if unicode.IsLetter(ch) || ch == '_' || ch == '/' || ch == '-' || ch == '=' || ch == '.' || ch == '+' || ch == '*' || ch == '<' || ch == '>' || ch == '!' {
		val := string(ch)
		for unicode.IsLetter(l.peekChar()) || unicode.IsDigit(l.peekChar()) || l.peekChar() == '_' || l.peekChar() == '/' || l.peekChar() == '-' || l.peekChar() == '=' || l.peekChar() == '.' || l.peekChar() == '+' || l.peekChar() == '*' || l.peekChar() == '<' || l.peekChar() == '>' || l.peekChar() == '!' {
			val += string(l.nextChar())
		}
		return Token{Type: TokenSymbol, Value: val, Line: startLine, Column: startCol}
	}

	reportError(fmt.Sprintf("Unexpected character: %c", ch), startLine, startCol)
	return Token{}
}

// AST
type Node struct {
	Type     string
	Value    string
	Children []*Node
	Line     int
	Column   int
	Filename string
}

// Parser
type Parser struct {
	lexer    *Lexer
	cur      Token
	filename string
}

func NewParser(lexer *Lexer, filename string) *Parser {
	p := &Parser{lexer: lexer, filename: filename}
	p.cur = p.lexer.NextToken()
	return p
}

func (p *Parser) parseExpression() *Node {
	if p.cur.Type == TokenLParen {
		node := &Node{Type: "List", Line: p.cur.Line, Column: p.cur.Column, Filename: p.filename}
		p.cur = p.lexer.NextToken() // consume '('
		for p.cur.Type != TokenRParen && p.cur.Type != TokenEOF {
			node.Children = append(node.Children, p.parseExpression())
		}
		if p.cur.Type != TokenRParen {
			reportError("Expected ')'", p.cur.Line, p.cur.Column)
		}
		p.cur = p.lexer.NextToken() // consume ')'
		return node
	}
	if p.cur.Type == TokenSymbol || p.cur.Type == TokenInt || p.cur.Type == TokenString {
		node := &Node{Type: string(p.cur.Type), Value: p.cur.Value, Line: p.cur.Line, Column: p.cur.Column, Filename: p.filename}
		p.cur = p.lexer.NextToken()
		return node
	}
	reportError(fmt.Sprintf("Unexpected token: %s", p.cur.Value), p.cur.Line, p.cur.Column)
	return nil
}

// AST Utilities
func copyNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	clone := &Node{Type: n.Type, Value: n.Value, Line: n.Line, Column: n.Column, Filename: n.Filename}
	for _, child := range n.Children {
		clone.Children = append(clone.Children, copyNode(child))
	}
	return clone
}

func replaceNext(node *Node, replacement *Node) {
	if node == nil {
		return
	}
	for i, child := range node.Children {
		if child.Type == "List" && len(child.Children) == 1 && child.Children[0].Value == "next" {
			node.Children[i] = replacement
		} else {
			replaceNext(child, replacement)
		}
	}
}

func renameVar(node *Node, oldName, newName string) {
	if node == nil {
		return
	}
	if node.Type == "SYMBOL" {
		if node.Value == oldName {
			node.Value = newName
		} else if strings.HasPrefix(node.Value, oldName+".") {
			node.Value = newName + node.Value[len(oldName):]
		}
	}
	for _, child := range node.Children {
		renameVar(child, oldName, newName)
	}
}

// Code Generator
func generateCode(node *Node) (string, string) {
	if node.Type != "List" || len(node.Children) == 0 {
		reportError("Expected list at root", node.Line, node.Column)
	}
	head := node.Children[0]
	if head.Type != "SYMBOL" || (head.Value != "http_server" && head.Value != "cli_app") {
		reportError("Expected http_server or cli_app as root symbol", head.Line, head.Column)
	}

	isCliApp := head.Value == "cli_app"

	var portNode *Node
	var startIndex int
	if isCliApp {
		startIndex = 1
	} else {
		if len(node.Children) < 3 {
			reportError("http_server expects at least a port and 1 route", head.Line, head.Column)
		}
		portNode = node.Children[1]
		if portNode.Type != "INT" {
			reportError("Expected integer for port", portNode.Line, portNode.Column)
		}
		startIndex = 2
	}

	var funcsCode string
	var routesCode string
	var cliCode string
	var testCode string
	var extraImports []string
	defaultImports := map[string]bool{
		"bytes":         true,
		"database/sql":  true,
		"encoding/json": true,
		"fmt":           true,
		"io":            true,
		"net/http":      true,
		"os":            true,
		"os/exec":       true,
		"regexp":        true,
		"strconv":       true,
		"strings":       true,
		"time":          true,
	}
	seenImports := make(map[string]bool)

	for i := startIndex; i < len(node.Children); i++ {
		handlerNode := node.Children[i]
		if handlerNode.Type != "List" || len(handlerNode.Children) == 0 {
			reportError("Expected route, defun, struct, import, test, or middleware definition", handlerNode.Line, handlerNode.Column)
		}

		head := handlerNode.Children[0].Value

		if head == "intent" {
			continue
		}

		if head == "test" {
			if len(handlerNode.Children) < 3 {
				reportError("test expects (test \"description\" body...)", handlerNode.Line, handlerNode.Column)
			}
			descNode := handlerNode.Children[1]
			if descNode.Type != "STRING" {
				reportError("test description must be a string", descNode.Line, descNode.Column)
			}
			desc := descNode.Value
			safeDesc := ""
			lastWasUnderscore := false
			for _, r := range desc {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					safeDesc += string(r)
					lastWasUnderscore = false
				} else {
					if !lastWasUnderscore {
						safeDesc += "_"
						lastWasUnderscore = true
					}
				}
			}
			safeDesc = strings.Trim(safeDesc, "_")
			testFuncName := "Test"
			if len(safeDesc) > 0 {
				testFuncName += "_" + safeDesc
			}

			var testBodyCode string
			for j := 2; j < len(handlerNode.Children); j++ {
				testBodyCode += generateStatement(handlerNode.Children[j], "", 0) + "\n"
			}
			testCode += fmt.Sprintf("func %s(t *testing.T) {\n%s\n}\n\n", testFuncName, testBodyCode)
			continue
		}

		if head == "import" {
			if len(handlerNode.Children) != 2 {
				reportError("import expects (import \"pkg\")", handlerNode.Line, handlerNode.Column)
			}
			pkgNode := handlerNode.Children[1]
			if pkgNode.Type != "STRING" {
				reportError("import package must be a string", pkgNode.Line, pkgNode.Column)
			}
			pkg := pkgNode.Value
			if !defaultImports[pkg] && !seenImports[pkg] {
				seenImports[pkg] = true
				extraImports = append(extraImports, pkg)
			}
			continue
		}

		if head == "struct" {
			if len(handlerNode.Children) < 2 {
				reportError("struct expects (struct Name (field type)...)", handlerNode.Line, handlerNode.Column)
			}
			name := handlerNode.Children[1].Value
			funcsCode += fmt.Sprintf("type %s struct {\n", name)
			for j := 2; j < len(handlerNode.Children); j++ {
				fieldNode := handlerNode.Children[j]
				if fieldNode.Type != "List" || len(fieldNode.Children) != 2 {
					reportError("struct field expects (name type)", fieldNode.Line, fieldNode.Column)
				}
				fieldName := fieldNode.Children[0].Value
				fieldType := fieldNode.Children[1].Value
				if len(fieldName) > 0 {
					fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]
				}
				funcsCode += fmt.Sprintf("\t%s %s\n", fieldName, fieldType)
			}
			funcsCode += "}\n\n"
			continue
		}

		if head == "defun" {
			if len(handlerNode.Children) < 4 {
				reportError("defun expects (defun name (args) body)", handlerNode.Line, handlerNode.Column)
			}
			name := handlerNode.Children[1].Value
			argsNode := handlerNode.Children[2]

			typeHints := make(map[string]string)
			for j := 3; j < len(handlerNode.Children)-1; j++ {
				cfgNode := handlerNode.Children[j]
				if cfgNode.Type == "List" && len(cfgNode.Children) >= 3 && cfgNode.Children[0].Value == "type_hint" {
					varName := cfgNode.Children[1].Value
					varType := cfgNode.Children[2].Value
					typeHints[varName] = varType
				}
			}

			var argsList []string
			for _, arg := range argsNode.Children {
				argType := "string"
				if t, ok := typeHints[arg.Value]; ok {
					argType = t
				}
				argsList = append(argsList, arg.Value+" "+argType)
			}
			argsStr := strings.Join(argsList, ", ")

			returnType := "string"
			if t, ok := typeHints["return"]; ok {
				returnType = t
			}

			bodyNode := handlerNode.Children[len(handlerNode.Children)-1]
			bodyCode := generateStatement(bodyNode, "", 0)
			funcsCode += fmt.Sprintf("//line %s:%d\nfunc %s(%s) %s {\n%s\n}\n\n", handlerNode.Filename, handlerNode.Line, name, argsStr, returnType, bodyCode)
			continue
		}

		if head == "route" {
			if len(handlerNode.Children) != 3 {
				reportError("route expects (route path handler)", handlerNode.Line, handlerNode.Column)
			}
			pathNode := handlerNode.Children[1]
			if pathNode.Type != "STRING" {
				reportError("route path must be a string", pathNode.Line, pathNode.Column)
			}
			reqNodeList := handlerNode.Children[2].Children[1]
			if reqNodeList.Type != "List" || len(reqNodeList.Children) != 1 {
				reportError("Expected exactly 1 argument in lambda (req)", reqNodeList.Line, reqNodeList.Column)
			}
			reqVar := reqNodeList.Children[0].Value
			bodyNode := handlerNode.Children[2].Children[2]
			bodyCode := generateStatement(bodyNode, reqVar, 0)
			routesCode += fmt.Sprintf(`	http.HandleFunc(%q, func(w http.ResponseWriter, %s *http.Request) {
%s
	})
`, pathNode.Value, reqVar, bodyCode)
			continue
		}

		if head == "middleware" {
			if len(handlerNode.Children) < 3 {
				reportError("middleware expects (middleware (lambda (req) body) routes...)", handlerNode.Line, handlerNode.Column)
			}
			lambdaNode := handlerNode.Children[1]
			if lambdaNode.Type != "List" || len(lambdaNode.Children) != 3 || lambdaNode.Children[0].Value != "lambda" {
				reportError("middleware expects a lambda", lambdaNode.Line, lambdaNode.Column)
			}
			reqNodeList := lambdaNode.Children[1]
			if reqNodeList.Type != "List" || len(reqNodeList.Children) != 1 {
				reportError("middleware lambda expects exactly 1 argument", reqNodeList.Line, reqNodeList.Column)
			}
			mwReqVar := reqNodeList.Children[0].Value
			mwBodyNode := lambdaNode.Children[2]

			for j := 2; j < len(handlerNode.Children); j++ {
				routeNode := handlerNode.Children[j]
				if routeNode.Type != "List" || len(routeNode.Children) == 0 || routeNode.Children[0].Value != "route" {
					reportError("middleware block can only contain routes", routeNode.Line, routeNode.Column)
				}
				if len(routeNode.Children) != 3 {
					reportError("route expects (route path handler)", routeNode.Line, routeNode.Column)
				}
				pathNode := routeNode.Children[1]
				if pathNode.Type != "STRING" {
					reportError("route path must be a string", pathNode.Line, pathNode.Column)
				}

				routeLambdaNode := routeNode.Children[2]
				routeReqList := routeLambdaNode.Children[1]
				routeReqVar := routeReqList.Children[0].Value
				routeBodyNode := routeLambdaNode.Children[2]

				clonedMwBody := copyNode(mwBodyNode)
				clonedRouteBody := copyNode(routeBodyNode)
				if routeReqVar != mwReqVar {
					renameVar(clonedRouteBody, routeReqVar, mwReqVar)
				}

				replaceNext(clonedMwBody, clonedRouteBody)
				combinedCode := generateStatement(clonedMwBody, mwReqVar, 0)

				routesCode += fmt.Sprintf(`	http.HandleFunc(%q, func(w http.ResponseWriter, %s *http.Request) {
%s
	})
`, pathNode.Value, mwReqVar, combinedCode)
			}
			continue
		}

		if isCliApp {
			// For cli_app, unhandled blocks are treated as statements executed in main
			cliCode += generateStatement(handlerNode, "", 0) + "\n"
			continue
		}

		reportError("Expected route, defun, struct, import, test, or middleware block", handlerNode.Line, handlerNode.Column)
	}

	code := `package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
`
	for _, imp := range extraImports {
		code += fmt.Sprintf("\t%q\n", imp)
	}
	code += `)
`
	code += funcsCode
	code += `func main() {
	var _ = sql.Open
	var _ = os.Getenv
	var _ = json.Marshal
	var _ = io.ReadAll
	var _ = bytes.NewBuffer
	var _ = http.DefaultClient
	var _ = exec.Command
	var _ = regexp.MatchString
	var _ = strings.Split
	var _ = time.Sleep
	var _ = strconv.Atoi
	var _ = fmt.Println
`
	if isCliApp {
		code += cliCode
		code += "}\n"
	} else {
		code += routesCode
		code += fmt.Sprintf(`	
	fmt.Println("Starting server on port %s...")
	if err := http.ListenAndServe(":%s", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
`, portNode.Value, portNode.Value)
	}

	if testCode != "" {
		fullTestCode := `package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
`
		for _, imp := range extraImports {
			parts := strings.Split(imp, "/")
			pkgName := parts[len(parts)-1]
			if strings.Contains(testCode, pkgName+".") {
				fullTestCode += fmt.Sprintf("\t%q\n", imp)
			}
		}
		fullTestCode += `)

var _ = sql.Open
var _ = os.Getenv
var _ = json.Marshal
var _ = io.ReadAll
var _ = bytes.NewBuffer
var _ = http.DefaultClient
var _ = exec.Command
var _ = regexp.MatchString
var _ = strings.Split
var _ = time.Sleep
var _ = strconv.Atoi
var _ = fmt.Println

` + testCode
		testCode = fullTestCode
	}

	return code, testCode
}

func generateStatement(node *Node, reqVar string, depth int) string {
	code := generateStatementRaw(node, reqVar, depth)
	if node.Type != "List" || len(node.Children) == 0 {
		return code
	}
	head := node.Children[0].Value
	switch head {
	case "return", "res_json", "res", "let", "do", "try_let", "spawn", "if", "print", "db_connect", "sql_query", "append", "map_set", "map_delete", "for", "sleep", "write_file", "mkdir", "exec", "while", "match", "set", "call", "cli_args":
		if node.Filename != "" {
			return fmt.Sprintf("//line %s:%d\n%s", node.Filename, node.Line, code)
		}
	}
	return code
}

func generateStatementRaw(node *Node, reqVar string, depth int) string {
	if depth > 1000 {
		reportError("AST too deep: exceeded maximum nesting limit of 1000", node.Line, node.Column)
	}
	if node.Type == "STRING" {
		return fmt.Sprintf("%q", node.Value)
	}
	if node.Type == "SYMBOL" || node.Type == "INT" {
		return node.Value
	}
	if node.Type != "List" || len(node.Children) == 0 {
		reportError("Expected list for statement", node.Line, node.Column)
	}
	head := node.Children[0].Value
	if head == "intent" {
		return ""
	}
	if head == "return" {
		if len(node.Children) != 2 {
			reportError("return expects (return val)", node.Line, node.Column)
		}
		valNode := node.Children[1]
		if valNode.Type == "STRING" {
			return fmt.Sprintf("		return %q", valNode.Value)
		} else {
			return fmt.Sprintf("		return %s", generateStatementRaw(valNode, reqVar, depth+1))
		}
	} else if head == "res_json" {
		if len(node.Children) != 3 {
			reportError("res_json expects (res_json status data)", node.Line, node.Column)
		}
		status := node.Children[1].Value
		dataNode := node.Children[2]
		dataVar := dataNode.Value
		if dataNode.Type == "STRING" {
			dataVar = fmt.Sprintf("%q", dataVar)
		}
		return fmt.Sprintf(`		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(%s)
		_ = json.NewEncoder(w).Encode(%s)`, status, dataVar)
	} else if head == "res" {
		if len(node.Children) != 4 {
			reportError("res expects status, contentType, and body", node.Line, node.Column)
		}
		status := node.Children[1].Value
		contentType := node.Children[2].Value
		resBody := node.Children[3].Value
		if node.Children[3].Type == "SYMBOL" || (node.Children[3].Type == "List" && node.Children[3].Children[0].Value == "call") {
			var bodyStr string
			if node.Children[3].Type == "List" {
				funcName := node.Children[3].Children[1].Value
				var args []string
				for j := 2; j < len(node.Children[3].Children); j++ {
					if node.Children[3].Children[j].Type == "STRING" {
						args = append(args, fmt.Sprintf("%q", node.Children[3].Children[j].Value))
					} else {
						args = append(args, node.Children[3].Children[j].Value)
					}
				}
				bodyStr = fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", "))
			} else {
				bodyStr = resBody
			}
			return fmt.Sprintf(`		w.Header().Set("Content-Type", %q)
		w.WriteHeader(%s)
		fmt.Fprint(w, %s)`, contentType, status, bodyStr)
		} else {
			return fmt.Sprintf(`		w.Header().Set("Content-Type", %q)
		w.WriteHeader(%s)
		fmt.Fprint(w, %q)`, contentType, status, resBody)
		}
	} else if head == "let" {
		if len(node.Children) < 3 {
			reportError("let expects (let (var val) body)", node.Line, node.Column)
		}
		var letPrefix strings.Builder
		letPrefix.WriteString("		{\n")
		declaredVars := make(map[string]bool)

		curr := node
		for curr.Type == "List" && len(curr.Children) == 3 && curr.Children[0].Value == "let" {
			binds := curr.Children[1]
			if binds.Type != "List" || len(binds.Children) != 2 {
				reportError("let binding expects (var val)", binds.Line, binds.Column)
			}
			varName := binds.Children[0].Value
			valNode := binds.Children[1]
			var valStr string
			if valNode.Type == "STRING" {
				valStr = fmt.Sprintf("%q", valNode.Value)
			} else if valNode.Type == "List" && len(valNode.Children) > 0 {
				funcName := valNode.Children[0].Value
				if funcName == "call" {
					var args []string
					for j := 2; j < len(valNode.Children); j++ {
						if valNode.Children[j].Type == "STRING" {
							args = append(args, fmt.Sprintf("%q", valNode.Children[j].Value))
						} else {
							args = append(args, valNode.Children[j].Value)
						}
					}
					valStr = fmt.Sprintf("%s(%s)", valNode.Children[1].Value, strings.Join(args, ", "))
				} else if funcName == "list" {
					var items []string
					for j := 1; j < len(valNode.Children); j++ {
						if valNode.Children[j].Type == "STRING" {
							items = append(items, fmt.Sprintf("%q", valNode.Children[j].Value))
						} else {
							items = append(items, valNode.Children[j].Value)
						}
					}
					valStr = fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
				} else if funcName == "dict" {
					var pairs []string
					for j := 1; j < len(valNode.Children); j++ {
						pair := valNode.Children[j]
						if pair.Type == "List" && len(pair.Children) == 2 {
							k := pair.Children[0].Value
							if pair.Children[0].Type == "STRING" {
								k = fmt.Sprintf("%q", k)
							}
							v := pair.Children[1].Value
							if pair.Children[1].Type == "STRING" {
								v = fmt.Sprintf("%q", v)
							}
							pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
						}
					}
					valStr = fmt.Sprintf("map[string]string{%s}", strings.Join(pairs, ", "))
				} else if funcName == "env" {
					if len(valNode.Children) != 2 {
						reportError("env expects (env \"KEY\")", valNode.Line, valNode.Column)
					}
					keyNode := valNode.Children[1]
					if keyNode.Type == "STRING" {
						valStr = fmt.Sprintf("os.Getenv(%q)", keyNode.Value)
					}
				} else if funcName == "parse_json" {
					if len(valNode.Children) != 3 {
						reportError("parse_json expects (parse_json Type body)", valNode.Line, valNode.Column)
					}
					// Handled downstream
				} else {
					valStr = generateStatement(valNode, reqVar, depth+1)
				}
			} else {
				valStr = generateStatement(valNode, reqVar, depth+1)
			}

			if valNode.Type == "List" && len(valNode.Children) > 0 && valNode.Children[0].Value == "parse_json" {
				targetType := valNode.Children[1].Value
				bodyVar := valNode.Children[2].Value
				if bodyVar == "req.body" {
					bodyVar = reqVar + ".Body"
				}
				letPrefix.WriteString(fmt.Sprintf("			var %s %s\n			_ = json.NewDecoder(%s).Decode(&%s)\n			_ = %s\n", varName, targetType, bodyVar, varName, varName))
				declaredVars[varName] = true
			} else {
				if declaredVars[varName] {
					letPrefix.WriteString(fmt.Sprintf("			%s = %s\n			_ = %s\n", varName, valStr, varName))
				} else {
					letPrefix.WriteString(fmt.Sprintf("			%s := %s\n			_ = %s\n", varName, valStr, varName))
					declaredVars[varName] = true
				}
			}

			curr = curr.Children[2]
		}

		bodyCode := generateStatement(curr, reqVar, depth+1)
		return fmt.Sprintf("%s%s\n		}", letPrefix.String(), bodyCode)
	} else if head == "do" {
		var stmts string
		for j := 1; j < len(node.Children); j++ {
			stmts += generateStatement(node.Children[j], reqVar, depth+1) + "\n"
		}
		return fmt.Sprintf("		{\n%s\n		}", stmts)
	} else if head == "try_let" {
		if len(node.Children) != 4 {
			reportError("try_let expects (try_let (var val) (catch err catchBody) successBody)", node.Line, node.Column)
		}
		binds := node.Children[1]
		if binds.Type != "List" || len(binds.Children) != 2 {
			reportError("try_let binding expects (var val)", binds.Line, binds.Column)
		}
		varName := binds.Children[0].Value
		valNode := binds.Children[1]

		catchNode := node.Children[2]
		if catchNode.Type != "List" || len(catchNode.Children) != 3 || catchNode.Children[0].Value != "catch" {
			reportError("try_let catch expects (catch errVar catchBody)", catchNode.Line, catchNode.Column)
		}
		errVar := catchNode.Children[1].Value
		catchBodyCode := generateStatement(catchNode.Children[2], reqVar, depth+1)
		successBodyCode := generateStatement(node.Children[3], reqVar, depth+1)

		if valNode.Type == "List" && len(valNode.Children) > 0 && valNode.Children[0].Value == "parse_json" {
			targetType := valNode.Children[1].Value
			bodyVar := valNode.Children[2].Value
			if bodyVar == "req.body" {
				bodyVar = reqVar + ".Body"
			}
			return fmt.Sprintf(`		{
			var %s %s
			if %s := json.NewDecoder(%s).Decode(&%s); %s != nil {
%s
			} else {
				_ = %s
%s
			}
		}`, varName, targetType, errVar, bodyVar, varName, errVar, catchBodyCode, varName, successBodyCode)
		}

		valStr := generateStatement(valNode, reqVar, depth+1)
		return fmt.Sprintf(`		{
			%s, %s := %s
			if %s != nil {
%s
			} else {
				_ = %s
%s
			}
		}`, varName, errVar, valStr, errVar, catchBodyCode, varName, successBodyCode)
	} else if head == "spawn" {
		if len(node.Children) != 2 {
			reportError("spawn expects (spawn (lambda () body))", node.Line, node.Column)
		}
		lambdaNode := node.Children[1]
		if lambdaNode.Type != "List" || len(lambdaNode.Children) != 3 || lambdaNode.Children[0].Value != "lambda" {
			reportError("spawn expects a lambda", lambdaNode.Line, lambdaNode.Column)
		}
		argsNode := lambdaNode.Children[1]
		if argsNode.Type != "List" || len(argsNode.Children) != 0 {
			reportError("spawn lambda expects no arguments ()", argsNode.Line, argsNode.Column)
		}
		bodyCode := generateStatement(lambdaNode.Children[2], reqVar, depth+1)
		return fmt.Sprintf("		go func() {\n%s\n		}()", bodyCode)
	} else if head == "if" {
		if len(node.Children) != 3 && len(node.Children) != 4 {
			reportError("if expects (if cond then) or (if cond then else)", node.Line, node.Column)
		}
		condNode := node.Children[1]
		if condNode.Type != "List" || len(condNode.Children) != 3 {
			reportError("cond expects (op a b)", condNode.Line, condNode.Column)
		}
		op := condNode.Children[0].Value
		if op == "=" {
			op = "=="
		}
		if op != "==" && op != "!=" && op != "<" && op != ">" && op != "<=" && op != ">=" {
			reportError("unsupported operator in if cond: "+op, condNode.Line, condNode.Column)
		}
		left := condNode.Children[1].Value
		if left == "req.method" {
			left = reqVar + ".Method"
		}
		right := condNode.Children[2]
		rightStr := right.Value
		if right.Type == "STRING" {
			rightStr = fmt.Sprintf("%q", rightStr)
		}

		thenCode := generateStatement(node.Children[2], reqVar, depth+1)

		if len(node.Children) == 3 {
			return fmt.Sprintf(`		if %s %s %s {
%s
		}`, left, op, rightStr, thenCode)
		}

		elseCode := generateStatement(node.Children[3], reqVar, depth+1)

		return fmt.Sprintf(`		if %s %s %s {
%s
		} else {
%s
		}`, left, op, rightStr, thenCode, elseCode)
	} else if head == "print" {
		var args []string
		for j := 1; j < len(node.Children); j++ {
			args = append(args, generateStatement(node.Children[j], reqVar, depth+1))
		}
		return fmt.Sprintf("		fmt.Println(%s)", strings.Join(args, ", "))
	} else if head == "db_connect" {
		if len(node.Children) != 4 {
			reportError("db_connect expects (db_connect var driver dsn)", node.Line, node.Column)
		}
		varName := node.Children[1].Value
		driverNode := node.Children[2]
		dsnNode := node.Children[3]
		return fmt.Sprintf("		%s, _ := sql.Open(%q, %q)\n		_ = %s", varName, driverNode.Value, dsnNode.Value, varName)
	} else if head == "sql_query" {
		if len(node.Children) != 3 {
			reportError("sql_query expects (sql_query db query)", node.Line, node.Column)
		}
		dbVar := node.Children[1].Value
		queryNode := node.Children[2]
		queryStr := queryNode.Value
		if queryNode.Type == "STRING" {
			queryStr = fmt.Sprintf("%q", queryStr)
		}
		return fmt.Sprintf("		%s.Query(%s)", dbVar, queryStr)
	} else if head == "append" {
		if len(node.Children) != 3 {
			reportError("append expects (append list item)", node.Line, node.Column)
		}
		listNode := node.Children[1]
		if listNode.Type != "SYMBOL" {
			reportError("append requires a symbol for list", node.Line, node.Column)
		}
		itemStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("		%s = append(%s, %s)", listNode.Value, listNode.Value, itemStr)
	} else if head == "map_set" {
		if len(node.Children) != 4 {
			reportError("map_set expects (map_set dict key val)", node.Line, node.Column)
		}
		dictNode := node.Children[1]
		if dictNode.Type != "SYMBOL" {
			reportError("map_set requires a symbol for dict", node.Line, node.Column)
		}
		keyStr := generateStatement(node.Children[2], reqVar, depth+1)
		valStr := generateStatement(node.Children[3], reqVar, depth+1)
		return fmt.Sprintf("		%s[%s] = %s", dictNode.Value, keyStr, valStr)
	} else if head == "map_delete" {
		if len(node.Children) != 3 {
			reportError("map_delete expects (map_delete dict key)", node.Line, node.Column)
		}
		dictNode := node.Children[1]
		if dictNode.Type != "SYMBOL" {
			reportError("map_delete requires a symbol for dict", node.Line, node.Column)
		}
		keyStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("		delete(%s, %s)", dictNode.Value, keyStr)
	} else if head == "for" {
		if len(node.Children) != 4 {
			reportError("for expects (for item list body)", node.Line, node.Column)
		}
		itemNode := node.Children[1].Value
		listNode := node.Children[2].Value
		bodyCode := generateStatement(node.Children[3], reqVar, depth+1)
		return fmt.Sprintf(`		for _, %s := range %s {
			_ = %s
%s
		}`, itemNode, listNode, itemNode, bodyCode)
	} else if head == "while" {
		if len(node.Children) != 3 {
			reportError("while expects (while cond body)", node.Line, node.Column)
		}
		condNode := node.Children[1]
		if condNode.Type != "List" || len(condNode.Children) != 3 {
			reportError("while cond expects (op a b)", condNode.Line, condNode.Column)
		}
		op := condNode.Children[0].Value
		if op == "=" {
			op = "=="
		}
		if op != "==" && op != "!=" && op != "<" && op != ">" && op != "<=" && op != ">=" {
			reportError("unsupported operator in while cond: "+op, condNode.Line, condNode.Column)
		}
		left := condNode.Children[1].Value
		if left == "req.method" {
			left = reqVar + ".Method"
		}
		right := condNode.Children[2].Value
		if condNode.Children[2].Type == "STRING" {
			right = fmt.Sprintf("%q", right)
		}
		bodyCode := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf(`		for %s %s %s {
%s
		}`, left, op, right, bodyCode)
	} else if head == "set" {
		if len(node.Children) != 3 {
			reportError("set expects (set var val)", node.Line, node.Column)
		}
		varStr := generateStatement(node.Children[1], reqVar, depth+1)
		valStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("		%s = %s", varStr, valStr)
	} else if head == "match" {
		if len(node.Children) < 3 {
			reportError("match expects (match var (val body)...)", node.Line, node.Column)
		}
		varStr := generateStatement(node.Children[1], reqVar, depth+1)
		
		var casesStr string
		for j := 2; j < len(node.Children); j++ {
			caseNode := node.Children[j]
			if caseNode.Type != "List" || len(caseNode.Children) != 2 {
				reportError("match case expects (val body)", caseNode.Line, caseNode.Column)
			}
			caseValNode := caseNode.Children[0]
			caseValStr := caseValNode.Value
			
			if caseValNode.Type == "SYMBOL" && caseValStr == "default" {
				caseValStr = "default"
			} else if caseValNode.Type == "STRING" {
				caseValStr = fmt.Sprintf("%q", caseValStr)
			}
			
			caseBodyCode := generateStatement(caseNode.Children[1], reqVar, depth+1)
			if caseValStr == "default" {
				casesStr += fmt.Sprintf("		default:\n%s\n", caseBodyCode)
			} else {
				casesStr += fmt.Sprintf("		case %s:\n%s\n", caseValStr, caseBodyCode)
			}
		}
		return fmt.Sprintf("		switch %s {\n%s		}", varStr, casesStr)
	} else if head == "sleep" {
		if len(node.Children) != 2 {
			reportError("sleep expects (sleep ms)", node.Line, node.Column)
		}
		msStr := generateStatement(node.Children[1], reqVar, depth+1)
		return fmt.Sprintf("		time.Sleep(time.Duration(%s) * time.Millisecond)", msStr)
	} else if head == "read_file" {
		if len(node.Children) != 2 {
			reportError("read_file expects (read_file path)", node.Line, node.Column)
		}
		pathStr := generateStatement(node.Children[1], reqVar, depth+1)
		return fmt.Sprintf("os.ReadFile(%s)", pathStr)
	} else if head == "write_file" {
		if len(node.Children) != 3 {
			reportError("write_file expects (write_file path data)", node.Line, node.Column)
		}
		pathStr := generateStatement(node.Children[1], reqVar, depth+1)
		dataStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("		os.WriteFile(%s, []byte(%s), 0644)", pathStr, dataStr)
	} else if head == "mkdir" {
		if len(node.Children) != 2 {
			reportError("mkdir expects (mkdir path)", node.Line, node.Column)
		}
		pathStr := generateStatement(node.Children[1], reqVar, depth+1)
		return fmt.Sprintf("		os.MkdirAll(%s, 0755)", pathStr)
	} else if head == "exec" {
		if len(node.Children) < 2 {
			reportError("exec expects (exec cmd args...)", node.Line, node.Column)
		}
		cmdStr := generateStatement(node.Children[1], reqVar, depth+1)
		var args []string
		for j := 2; j < len(node.Children); j++ {
			args = append(args, generateStatement(node.Children[j], reqVar, depth+1))
		}
		return fmt.Sprintf("func() ([]byte, error) { return exec.Command(%s, %s).CombinedOutput() }()", cmdStr, strings.Join(args, ", "))
	} else if head == "str_split" {
		if len(node.Children) != 3 {
			reportError("str_split expects (str_split s sep)", node.Line, node.Column)
		}
		sStr := generateStatement(node.Children[1], reqVar, depth+1)
		sepStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("strings.Split(%s, %s)", sStr, sepStr)
	} else if head == "str_join" {
		if len(node.Children) != 3 {
			reportError("str_join expects (str_join list sep)", node.Line, node.Column)
		}
		listStr := generateStatement(node.Children[1], reqVar, depth+1)
		sepStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("strings.Join(%s, %s)", listStr, sepStr)
	} else if head == "regex_match" {
		if len(node.Children) != 3 {
			reportError("regex_match expects (regex_match pattern s)", node.Line, node.Column)
		}
		patStr := generateStatement(node.Children[1], reqVar, depth+1)
		sStr := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("regexp.MatchString(%s, %s)", patStr, sStr)
	} else if head == "rate_limit" {
		if len(node.Children) != 3 {
			reportError("rate_limit expects (rate_limit \"10/s\" body)", node.Line, node.Column)
		}
		rateStr := node.Children[1].Value
		bodyCode := generateStatement(node.Children[2], reqVar, depth+1)
		// simple implementation: "10/s" -> sleep 100ms
		ms := 1000
		if strings.HasSuffix(rateStr, "/s") {
			n, _ := strconv.Atoi(strings.TrimSuffix(rateStr, "/s"))
			if n > 0 { ms = 1000 / n }
		}
		return fmt.Sprintf(`		{
			time.Sleep(%d * time.Millisecond)
			%s
		}`, ms, bodyCode)
	} else if head == "retry" {
		if len(node.Children) != 3 {
			reportError("retry expects (retry times body)", node.Line, node.Column)
		}
		timesStr := generateStatement(node.Children[1], reqVar, depth+1)
		bodyCode := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf(`		for i := 0; i < %s; i++ {
			%s
		}`, timesStr, bodyCode)
	} else if head == "fetch" {
		if len(node.Children) != 3 {
			reportError("fetch expects (fetch url method)", node.Line, node.Column)
		}
		urlStr := generateStatement(node.Children[1], reqVar, depth+1)
		methodStr := generateStatement(node.Children[2], reqVar, depth+1)

		return fmt.Sprintf(`func() ([]byte, error) {
			req, err := http.NewRequest(%s, %s, nil)
			if err != nil { return nil, err }
			resp, err := http.DefaultClient.Do(req)
			if err != nil { return nil, err }
			defer resp.Body.Close()
			return io.ReadAll(resp.Body)
		}()`, methodStr, urlStr)
	} else if head == "llm_generate" {
		if len(node.Children) < 2 {
			reportError("llm_generate expects (llm_generate prompt [model])", node.Line, node.Column)
		}
		promptStr := generateStatement(node.Children[1], reqVar, depth+1)
		modelStr := `"llama3"`
		if len(node.Children) >= 3 {
			modelStr = generateStatement(node.Children[2], reqVar, depth+1)
		}

		return fmt.Sprintf(`func() (string, error) {
			reqBody, _ := json.Marshal(map[string]any{
				"model":  %s,
				"prompt": %s,
				"stream": false,
			})
			resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(reqBody))
			if err != nil { return "", err }
			defer resp.Body.Close()
			var res struct {
				Response string `+"`json:\"response\"`"+`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil { return "", err }
			return res.Response, nil
		}()`, modelStr, promptStr)
	} else if head == "fuzzy_cast" {
		if len(node.Children) < 3 {
			reportError("fuzzy_cast expects (fuzzy_cast Type var [model])", node.Line, node.Column)
		}
		typeStr := node.Children[1].Value
		varStr := generateStatement(node.Children[2], reqVar, depth+1)
		modelStr := `"llama3"`
		if len(node.Children) >= 4 {
			modelStr = generateStatement(node.Children[3], reqVar, depth+1)
		}

		return fmt.Sprintf(`func() (%s, error) {
			var out %s
			reqBody, _ := json.Marshal(map[string]any{
				"model":  %s,
				"prompt": fmt.Sprintf("Coerce this input into a valid JSON object matching the requested schema. Reply strictly with the JSON object and nothing else.\nInput: %%s", %s),
				"stream": false,
				"format": "json",
			})
			resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(reqBody))
			if err != nil { return out, err }
			defer resp.Body.Close()
			var res struct {
				Response string `+"`json:\"response\"`"+`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil { return out, err }
			err = json.Unmarshal([]byte(res.Response), &out)
			return out, err
		}()`, typeStr, typeStr, modelStr, varStr)
	} else if head == "assert_semantic" {
		if len(node.Children) != 3 {
			reportError("assert_semantic expects (assert_semantic var \"condition\")", node.Line, node.Column)
		}
		varStr := generateStatement(node.Children[1], reqVar, depth+1)
		condStr := generateStatement(node.Children[2], reqVar, depth+1)
		
		return fmt.Sprintf(`func() bool {
			reqBody, _ := json.Marshal(map[string]any{
				"model":  "llama3",
				"prompt": fmt.Sprintf("Does this input satisfy the condition: '%%s'? Reply strictly with 'true' or 'false' and nothing else.\nInput: %%s", %s, %s),
				"stream": false,
			})
			resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(reqBody))
			if err != nil { return false }
			defer resp.Body.Close()
			var res struct {
				Response string `+"`json:\"response\"`"+`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil { return false }
			return strings.TrimSpace(strings.ToLower(res.Response)) == "true"
		}()`, condStr, varStr)
	} else if head == "+" || head == "-" || head == "*" || head == "/" || head == "<" || head == ">" || head == "and" || head == "or" || head == "==" {
		if len(node.Children) != 3 {
			reportError(fmt.Sprintf("%s expects 2 arguments", head), node.Line, node.Column)
		}
		op := head
		if op == "and" {
			op = "&&"
		}
		if op == "or" {
			op = "||"
		}
		arg1 := generateStatement(node.Children[1], reqVar, depth+1)
		arg2 := generateStatement(node.Children[2], reqVar, depth+1)
		return fmt.Sprintf("(%s %s %s)", arg1, op, arg2)
	} else if head == "call" {
		if len(node.Children) < 2 {
			reportError("call expects (call func args...)", node.Line, node.Column)
		}
		funcName := node.Children[1].Value
		var args []string
		for j := 2; j < len(node.Children); j++ {
			argNode := node.Children[j]
			if argNode.Type == "STRING" {
				args = append(args, fmt.Sprintf("%q", argNode.Value))
			} else {
				args = append(args, argNode.Value)
			}
		}
		return fmt.Sprintf("		%s(%s)", funcName, strings.Join(args, ", "))
	} else if head == "cli_args" {
		if len(node.Children) == 1 {
			return "os.Args[1:]"
		} else if len(node.Children) == 2 {
			idxStr := generateStatement(node.Children[1], reqVar, depth+1)
			return fmt.Sprintf("func() string { _idx, _ := strconv.Atoi(fmt.Sprint(%s)); if len(os.Args) > _idx+1 { return os.Args[_idx+1] }; return \"\" }()", idxStr)
		} else {
			reportError("cli_args expects (cli_args) or (cli_args index)", node.Line, node.Column)
		}
	}
	reportError(fmt.Sprintf("Unknown statement: %s", head), node.Line, node.Column)
	return ""
}

func expandIncludes(node *Node, baseDir string, depth int) {
	if depth > 100 {
		reportError("Include depth exceeded (circular include?)", node.Line, node.Column)
	}
	if node.Type != "List" {
		return
	}
	var newChildren []*Node
	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]
		if child.Type == "List" && len(child.Children) == 2 && child.Children[0].Type == "SYMBOL" && child.Children[0].Value == "include" {
			filenameNode := child.Children[1]
			if filenameNode.Type != "STRING" {
				reportError("include expects a string filename", child.Line, child.Column)
			}
			filename := filenameNode.Value
			fullPath := filepath.Join(baseDir, filename)

			content, err := os.ReadFile(fullPath)
			if err != nil {
				reportError(fmt.Sprintf("Failed to read included file %q: %v", filename, err), child.Line, child.Column)
			}

			lexer := NewLexer(string(content))
			parser := NewParser(lexer, filepath.Base(fullPath))
			includedAst := parser.parseExpression()

			if parser.cur.Type != TokenEOF {
				reportError(fmt.Sprintf("Unexpected tokens after EOF in included file %q", filename), parser.cur.Line, parser.cur.Column)
			}

			expandIncludes(includedAst, filepath.Dir(fullPath), depth+1)

			if includedAst.Type == "List" && len(includedAst.Children) > 0 && includedAst.Children[0].Value == "module" {
				newChildren = append(newChildren, includedAst.Children[1:]...)
			} else {
				newChildren = append(newChildren, includedAst)
			}
		} else {
			expandIncludes(child, baseDir, depth)
			newChildren = append(newChildren, child)
		}
	}
	node.Children = newChildren
}

func applyPatches(node *Node) {
	if node == nil || node.Type != "List" {
		return
	}
	var newChildren []*Node
	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]
		if child.Type == "List" && len(child.Children) == 3 && child.Children[0].Type == "SYMBOL" && child.Children[0].Value == "patch" {
			funcNameNode := child.Children[1]
			if funcNameNode.Type != "SYMBOL" {
				reportError("patch expects a symbol for function name", funcNameNode.Line, funcNameNode.Column)
			}
			funcName := funcNameNode.Value
			newBody := child.Children[2]

			found := false
			for j := 0; j < len(newChildren); j++ {
				target := newChildren[j]
				if target.Type == "List" && len(target.Children) >= 4 && target.Children[0].Type == "SYMBOL" && target.Children[0].Value == "defun" && target.Children[1].Value == funcName {
					target.Children[len(target.Children)-1] = copyNode(newBody)
					found = true
					break
				}
			}
			if !found {
				for j := i + 1; j < len(node.Children); j++ {
					target := node.Children[j]
					if target.Type == "List" && len(target.Children) >= 4 && target.Children[0].Type == "SYMBOL" && target.Children[0].Value == "defun" && target.Children[1].Value == funcName {
						target.Children[len(target.Children)-1] = copyNode(newBody)
						found = true
						break
					}
				}
			}
			if !found {
				reportError(fmt.Sprintf("patch target function %q not found", funcName), child.Line, child.Column)
			}
		} else {
			applyPatches(child)
			newChildren = append(newChildren, child)
		}
	}
	node.Children = newChildren
}

func applyWithContext(node *Node, ctxVars []*Node) *Node {
	if node == nil {
		return nil
	}
	if node.Type != "List" || len(node.Children) == 0 {
		return node
	}

	head := node.Children[0].Value
	if head == "with_context" {
		if len(node.Children) < 3 {
			reportError("with_context expects (with_context vars body...)", node.Line, node.Column)
		}
		varsNode := node.Children[1]
		var newCtxVars []*Node
		newCtxVars = append(newCtxVars, ctxVars...)
		if varsNode.Type == "SYMBOL" {
			newCtxVars = append(newCtxVars, varsNode)
		} else if varsNode.Type == "List" {
			for _, v := range varsNode.Children {
				if v.Type != "SYMBOL" {
					reportError("with_context vars must be symbols", v.Line, v.Column)
				}
				newCtxVars = append(newCtxVars, v)
			}
		} else {
			reportError("with_context expects symbol or list of symbols", varsNode.Line, varsNode.Column)
		}

		doNode := &Node{Type: "List", Line: node.Line, Column: node.Column, Filename: node.Filename}
		doNode.Children = append(doNode.Children, &Node{Type: "SYMBOL", Value: "do", Line: node.Line, Column: node.Column})
		for i := 2; i < len(node.Children); i++ {
			doNode.Children = append(doNode.Children, applyWithContext(node.Children[i], newCtxVars))
		}
		return doNode
	}

	if head == "call" && len(ctxVars) > 0 {
		newNode := &Node{Type: "List", Line: node.Line, Column: node.Column, Filename: node.Filename}
		newNode.Children = append(newNode.Children, node.Children[0])
		if len(node.Children) > 1 {
			newNode.Children = append(newNode.Children, node.Children[1])
		}
		for _, cv := range ctxVars {
			newNode.Children = append(newNode.Children, copyNode(cv))
		}
		for i := 2; i < len(node.Children); i++ {
			newNode.Children = append(newNode.Children, applyWithContext(node.Children[i], ctxVars))
		}
		return newNode
	}

	newNode := &Node{Type: "List", Line: node.Line, Column: node.Column, Filename: node.Filename}
	newNode.Children = append(newNode.Children, node.Children[0])
	for i := 1; i < len(node.Children); i++ {
		newNode.Children = append(newNode.Children, applyWithContext(node.Children[i], ctxVars))
	}
	return newNode
}

func main() {
	if len(os.Args) < 2 {
		reportError("Missing file argument", 0, 0)
	}
	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		reportError(fmt.Sprintf("Cannot read file: %v", err), 0, 0)
	}

	lexer := NewLexer(string(content))
	parser := NewParser(lexer, filepath.Base(os.Args[1]))
	ast := parser.parseExpression()

	if parser.cur.Type != TokenEOF {
		reportError("Unexpected tokens after EOF", parser.cur.Line, parser.cur.Column)
	}

	expandIncludes(ast, filepath.Dir(os.Args[1]), 0)

	applyPatches(ast)
	ast = applyWithContext(ast, nil)

	goCode, testCode := generateCode(ast)

	err = os.WriteFile("server.go", []byte(goCode), 0644)
	if err != nil {
		reportError(fmt.Sprintf("Failed to write server.go: %v", err), 0, 0)
	}

	if testCode != "" {
		err = os.WriteFile("server_test.go", []byte(testCode), 0644)
		if err != nil {
			reportError(fmt.Sprintf("Failed to write server_test.go: %v", err), 0, 0)
		}
	} else {
		os.Remove("server_test.go")
	}
}
