package stdlib

import (
	"github.com/graphql-go/graphql"
)

var (
	fakeWeb              = new(GraphQL)
	testInvalidDirective string
	testDirectives       = map[string]interface{}{
		"testSchemaDirective":               testSchemaDirective,
		"testScalarDirective":               testScalarDirective,
		"testObjectDirective":               testObjectDirective,
		"testFieldDefinitionDirective":      testFieldDefinitionDirective,
		"testArgumentDefinitionDirective":   testArgumentDefinitionDirective,
		"testInterfaceDirective":            testInterfaceDirective,
		"testUnionDirective":                testUnionDirective,
		"testEnumDirective":                 testEnumDirective,
		"testEnumValueDirective":            testEnumValueDirective,
		"testInputObjectDirective":          testInputObjectDirective,
		"testInputFieldDefinitionDirective": testInputFieldDefinitionDirective,
	}
)

func testSchemaDirective(*graphql.SchemaConfig, map[string]interface{})                         {}
func testScalarDirective(*graphql.ScalarConfig, map[string]interface{})                         {}
func testObjectDirective(*graphql.ObjectConfig, map[string]interface{})                         {}
func testFieldDefinitionDirective(*graphql.Field, map[string]interface{})                       {}
func testArgumentDefinitionDirective(*graphql.ArgumentConfig, map[string]interface{})           {}
func testInterfaceDirective(*graphql.InterfaceConfig, map[string]interface{})                   {}
func testUnionDirective(*graphql.UnionConfig, map[string]interface{})                           {}
func testEnumDirective(*graphql.EnumConfig, map[string]interface{})                             {}
func testEnumValueDirective(*graphql.EnumValueConfig, map[string]interface{})                   {}
func testInputObjectDirective(*graphql.InputObjectConfig, map[string]interface{})               {}
func testInputFieldDefinitionDirective(*graphql.InputObjectFieldConfig, map[string]interface{}) {}
func testInvalidFuncDirective()                                                                 {}

//func TestAddResolvers(t *testing.T) {
//	fakeWeb.directives = make(tools.SchemaDirectiveVisitorMap)
//
//	for name, directive := range testDirectives {
//		fakeWeb.AddDirective(name, directive)
//	}
//
//	assert.True(t, fakeWeb.DirectiveExists("testSchemaDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testScalarDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testObjectDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testFieldDefinitionDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testArgumentDefinitionDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testInterfaceDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testUnionDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testEnumDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testEnumValueDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testInputObjectDirective"))
//	assert.True(t, fakeWeb.DirectiveExists("testInputFieldDefinitionDirective"))
//	assert.Panics(t, func() {
//		fakeWeb.AddDirective("testInvalidDirective", testInvalidDirective)
//	})
//	assert.Panics(t, func() {
//		fakeWeb.AddDirective("testInvalidFuncDirective", testInvalidFuncDirective)
//	})
//}
//
//func TestDropDirective(t *testing.T) {
//	TestAddResolvers(t)
//	fakeWeb.DropDirective("testSchemaDirective")
//	assert.False(t, fakeWeb.DirectiveExists("testSchemaDirective"))
//}
