package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ResourceInfo struct {
	Name           string
	OptionalFields []string
	TestFile       string
	TestedFields   []string
}

func main() {
	resourcesDir := filepath.Join("internal", "resources")
	testsDir := filepath.Join("internal", "resources_acceptance_tests")

	resources, err := analyzeResources(resourcesDir, testsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing resources: %v\n", err)
		os.Exit(1)
	}

	// Report findings
	hasGaps := false
	for _, res := range resources {
		if len(res.OptionalFields) == 0 {
			continue // Skip resources with no optional fields
		}

		missing := findMissingFields(res.OptionalFields, res.TestedFields)
		if len(missing) > 0 {
			hasGaps = true
			fmt.Printf("\n❌ %s:\n", res.Name)
			fmt.Printf("   Optional fields: %s\n", strings.Join(res.OptionalFields, ", "))
			if len(res.TestedFields) > 0 {
				fmt.Printf("   Tested fields:   %s\n", strings.Join(res.TestedFields, ", "))
			} else {
				fmt.Printf("   Tested fields:   (none - no removeOptionalFields test found)\n")
			}
			fmt.Printf("   Missing tests:   %s\n", strings.Join(missing, ", "))
		}
	}

	if !hasGaps {
		fmt.Println("✅ All resources with optional fields have complete test coverage!")
	} else {
		fmt.Println("\nNote: Some fields may be intentionally untested (e.g., Computed+Default fields).")
		os.Exit(1)
	}
}

func analyzeResources(resourcesDir, testsDir string) ([]ResourceInfo, error) {
	files, err := filepath.Glob(filepath.Join(resourcesDir, "*_resource.go"))
	if err != nil {
		return nil, err
	}

	var resources []ResourceInfo
	for _, file := range files {
		// Skip non-resource files
		basename := filepath.Base(file)
		if strings.HasSuffix(basename, "_data_source.go") {
			continue
		}

		resourceName := strings.TrimSuffix(basename, ".go")
		resourceName = strings.TrimSuffix(resourceName, "_resource")

		optionalFields, err := extractOptionalFields(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error parsing %s: %v\n", file, err)
			continue
		}

		testFile := filepath.Join(testsDir, resourceName+"_resource_test.go")
		testedFields := extractTestedFields(testFile)

		resources = append(resources, ResourceInfo{
			Name:           resourceName,
			OptionalFields: optionalFields,
			TestFile:       testFile,
			TestedFields:   testedFields,
		})
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	return resources, nil
}

func extractOptionalFields(filename string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var optionalFields []string

	// Find the Schema method
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "Schema" {
			return true
		}

		// Look for schema.Schema{...} with Attributes map
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			compositeLit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			// Check if this is a map[string]schema.Attribute
			for _, elt := range compositeLit.Elts {
				kvExpr, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				// Get the field name
				key, ok := kvExpr.Key.(*ast.BasicLit)
				if !ok || key.Kind != token.STRING {
					continue
				}
				fieldName := strings.Trim(key.Value, `"`)

				// Skip common computed/internal fields
				if fieldName == "id" || fieldName == "url" || fieldName == "display" {
					continue
				}

				// Check if this attribute has Optional: true
				valueComposite, ok := kvExpr.Value.(*ast.CompositeLit)
				if !ok {
					continue
				}

				hasOptional := false
				hasRequired := false

				for _, field := range valueComposite.Elts {
					fieldKV, ok := field.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					fieldKey, ok := fieldKV.Key.(*ast.Ident)
					if !ok {
						continue
					}

					if fieldKey.Name == "Optional" {
						if ident, ok := fieldKV.Value.(*ast.Ident); ok && ident.Name == "true" {
							hasOptional = true
						}
					} else if fieldKey.Name == "Required" {
						if ident, ok := fieldKV.Value.(*ast.Ident); ok && ident.Name == "true" {
							hasRequired = true
						}
					}
				}

				// Only include Optional-only fields (not Required, not Computed)
				// Actually, include Optional+Computed too since some might need testing
				if hasOptional && !hasRequired {
					optionalFields = append(optionalFields, fieldName)
				}
			}

			return true
		})

		return false // Stop after finding Schema method
	})

	sort.Strings(optionalFields)
	return optionalFields, nil
}

func extractTestedFields(testFile string) []string {
	data, err := os.ReadFile(testFile)
	if err != nil {
		return nil
	}

	content := string(data)

	// Find the removeOptionalFields test
	removeTestPattern := regexp.MustCompile(`func TestAcc\w+Resource_removeOptionalFields\(t \*testing\.T\)`)
	if !removeTestPattern.MatchString(content) {
		return nil // No removeOptionalFields test found
	}

	// Extract TestCheckNoResourceAttr calls within the removeOptionalFields test
	// Find the test function first
	testStartIdx := removeTestPattern.FindStringIndex(content)
	if testStartIdx == nil {
		return nil
	}

	// Find the corresponding closing brace (simplified - just look for next test or EOF)
	testEndIdx := len(content)
	nextTestPattern := regexp.MustCompile(`\nfunc TestAcc`)
	if nextMatch := nextTestPattern.FindStringIndex(content[testStartIdx[1]:]); nextMatch != nil {
		testEndIdx = testStartIdx[1] + nextMatch[0]
	}

	testContent := content[testStartIdx[0]:testEndIdx]

	// Extract field names from TestCheckNoResourceAttr
	// Pattern: resource.TestCheckNoResourceAttr("netbox_xxx.test", "field_name")
	noAttrPattern := regexp.MustCompile(`resource\.TestCheckNoResourceAttr\("netbox_[^"]+",\s*"([^"]+)"\)`)
	matches := noAttrPattern.FindAllStringSubmatch(testContent, -1)

	seen := make(map[string]bool)
	var fields []string
	for _, match := range matches {
		if len(match) > 1 {
			field := match[1]
			if !seen[field] {
				seen[field] = true
				fields = append(fields, field)
			}
		}
	}

	sort.Strings(fields)
	return fields
}

func findMissingFields(optionalFields, testedFields []string) []string {
	tested := make(map[string]bool)
	for _, field := range testedFields {
		tested[field] = true
	}

	var missing []string
	for _, field := range optionalFields {
		if !tested[field] {
			missing = append(missing, field)
		}
	}

	return missing
}
