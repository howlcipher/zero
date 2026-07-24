import sys

with open("zero.go", "r") as f:
    content = f.read()

target = """			typeHints := make(map[string]string)
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
			returnTypeStr := " " + returnType
			if returnType == "void" {
				returnTypeStr = ""
			}

			bodyNode := handlerNode.Children[len(handlerNode.Children)-1]
			bodyCode := generateStatement(bodyNode, "", 0)
			funcsCode += fmt.Sprintf("//line %s:%d\\nfunc %s(%s)%s {\\n%s\\n}\\n\\n", handlerNode.Filename, handlerNode.Line, name, argsStr, returnTypeStr, bodyCode)
			continue"""

replacement = """			typeHints := make(map[string]string)
			var typeParams []string
			for j := 3; j < len(handlerNode.Children)-1; j++ {
				cfgNode := handlerNode.Children[j]
				if cfgNode.Type == "List" && len(cfgNode.Children) >= 3 && cfgNode.Children[0].Value == "type_hint" {
					varName := cfgNode.Children[1].Value
					varType := cfgNode.Children[2].Value
					typeHints[varName] = varType
				} else if cfgNode.Type == "List" && len(cfgNode.Children) >= 2 && cfgNode.Children[0].Value == "type_param" {
					typeParams = append(typeParams, cfgNode.Children[1].Value)
				}
			}

			typeParamsStr := ""
			if len(typeParams) > 0 {
				var typed []string
				for _, tp := range typeParams {
					typed = append(typed, tp+" any")
				}
				typeParamsStr = "[" + strings.Join(typed, ", ") + "]"
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
			returnTypeStr := " " + returnType
			if returnType == "void" {
				returnTypeStr = ""
			}

			bodyNode := handlerNode.Children[len(handlerNode.Children)-1]
			bodyCode := generateStatement(bodyNode, "", 0)
			funcsCode += fmt.Sprintf("//line %s:%d\\nfunc %s%s(%s)%s {\\n%s\\n}\\n\\n", handlerNode.Filename, handlerNode.Line, name, typeParamsStr, argsStr, returnTypeStr, bodyCode)
			continue"""

if target in content:
    content = content.replace(target, replacement)
    with open("zero.go", "w") as f:
        f.write(content)
    print("Replaced successfully")
else:
    print("Target not found")
