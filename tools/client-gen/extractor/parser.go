package extractor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

// InterfaceInfo holds information about an interface
type InterfaceInfo struct {
	Name    string       `json:"name"`
	Methods []MethodInfo `json:"methods"`
}

// MethodInfo holds information about a method in an interface
type MethodInfo struct {
	Name      string      `json:"name"`
	Arguments []ParamInfo `json:"arguments"`
	Results   []TypeInfo  `json:"results"`
}

// ParamInfo holds information about a parameter
type ParamInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// TypeInfo holds information about a type
type TypeInfo struct {
	Type string `json:"type"`
}

// ExtractInterfaces extracts interface information from the given file
func ExtractInterfaces(fset *token.FileSet, file *ast.File) []InterfaceInfo {
	var interfaces []InterfaceInfo

	// Inspect the AST and extract interfaces
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for type declarations within GenDecl nodes
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Check if it's an interface
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			// Create a new InterfaceInfo
			iface := InterfaceInfo{
				Name: typeSpec.Name.Name,
			}

			// Extract methods from the interface
			if interfaceType.Methods != nil {
				for _, method := range interfaceType.Methods.List {
					// Skip embedded interfaces
					if len(method.Names) == 0 {
						continue
					}

					methodInfo := MethodInfo{
						Name: method.Names[0].Name,
					}

					// Get the function type
					funcType, ok := method.Type.(*ast.FuncType)
					if !ok {
						continue
					}

					// Extract arguments
					if funcType.Params != nil {
						for _, param := range funcType.Params.List {
							typeStr := typeToString(fset, param.Type)

							// A single field can have multiple names (e.g., a, b int)
							if len(param.Names) > 0 {
								for _, name := range param.Names {
									paramInfo := ParamInfo{
										Name: name.Name,
										Type: typeStr,
									}
									methodInfo.Arguments = append(methodInfo.Arguments, paramInfo)
								}
							} else {
								// Unnamed parameter
								paramInfo := ParamInfo{
									Name: "",
									Type: typeStr,
								}
								methodInfo.Arguments = append(methodInfo.Arguments, paramInfo)
							}
						}
					}

					// Extract results
					if funcType.Results != nil {
						for _, result := range funcType.Results.List {
							typeStr := typeToString(fset, result.Type)

							// A single result field can represent multiple results of the same type
							if len(result.Names) > 0 {
								for _, name := range result.Names {
									resultInfo := TypeInfo{
										Type: fmt.Sprintf("%s %s", name.Name, typeStr),
									}
									methodInfo.Results = append(methodInfo.Results, resultInfo)
								}
							} else {
								// Unnamed result
								resultInfo := TypeInfo{
									Type: typeStr,
								}
								methodInfo.Results = append(methodInfo.Results, resultInfo)
							}
						}
					}

					iface.Methods = append(iface.Methods, methodInfo)
				}
			}

			interfaces = append(interfaces, iface)
		}
		return true
	})

	return interfaces
}

// typeToString converts an AST type to a string representation
func typeToString(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	err := format.Node(&buf, fset, expr)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return strings.TrimSpace(buf.String())
}
