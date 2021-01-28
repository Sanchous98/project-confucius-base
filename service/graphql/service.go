package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/Sanchous98/project-confucius-base/service/web"
	tools "github.com/bhoriuchi/graphql-go-tools"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"io/ioutil"
	"log"
	"reflect"
)

type (
	// Type aliases for visitor functions
	VisitSchema               = func(*graphql.SchemaConfig, map[string]interface{})
	VisitScalar               = func(*graphql.ScalarConfig, map[string]interface{})
	VisitObject               = func(*graphql.ObjectConfig, map[string]interface{})
	VisitFieldDefinition      = func(*graphql.Field, map[string]interface{})
	VisitArgumentDefinition   = func(*graphql.ArgumentConfig, map[string]interface{})
	VisitInterface            = func(*graphql.InterfaceConfig, map[string]interface{})
	VisitUnion                = func(*graphql.UnionConfig, map[string]interface{})
	VisitEnum                 = func(*graphql.EnumConfig, map[string]interface{})
	VisitEnumValue            = func(*graphql.EnumValueConfig, map[string]interface{})
	VisitInputObject          = func(*graphql.InputObjectConfig, map[string]interface{})
	VisitInputFieldDefinition = func(*graphql.InputObjectFieldConfig, map[string]interface{})

	GraphQL struct {
		config     *config
		directives tools.SchemaDirectiveVisitorMap
		web        *web.Web
	}
)

func (g *GraphQL) Construct(web *web.Web) *GraphQL {
	g.config = new(config)
	err := g.config.Unmarshall()

	if err != nil {
		panic(err)
	}

	g.web = web

	if g.directives == nil {
		g.directives = make(tools.SchemaDirectiveVisitorMap)
	}
	// Add pre-defined directives
	for name, directive := range preDefinedDirectives {
		g.AddDirective(name, directive)
	}

	g.web.Router.GET("/api", g.queryHandler)
	g.web.Router.GET("/graphiql", g.handleGraphiQL(g.config.SchemaPath))

	return g
}

// resolveSchema initializes GraphQL server configuration, based on schema file and predefined resolvers and directives
func (g *GraphQL) resolveSchema(schemaContent []byte) *graphql.Schema {
	schema, err := tools.MakeExecutableSchema(tools.ExecutableSchema{
		TypeDefs:         string(schemaContent),
		SchemaDirectives: g.directives,
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to parse schema, error: %v", err))
	}

	return &schema
}

// handleGraphiQL provides GraphiQL playground
func (g *GraphQL) handleGraphiQL(schemaPath string) fasthttp.RequestHandler {
	b, err := ioutil.ReadFile(schemaPath)

	if err != nil {
		panic(fmt.Errorf("cannot open schema file: %v\n", err))
	}

	return fasthttpadaptor.NewFastHTTPHandler(handler.New(&handler.Config{
		Schema:     g.resolveSchema(b),
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	}))
}

// queryHandler is a function that handles GraphQL queries
func (g *GraphQL) queryHandler(context *fasthttp.RequestCtx) {
	b, err := ioutil.ReadFile(g.config.SchemaPath)
	if err != nil {
		panic("Schema file doesn't exist")
	}

	query := context.Request.Body()

	response := graphql.Do(graphql.Params{
		Schema:        *g.resolveSchema(b),
		RequestString: string(query),
	})

	if len(response.Errors) > 0 {
		log.Printf("Failed to execute graphql operation, errors: %+v", response.Errors)
	}

	//writer.Header().Set("X-CSRF-Token", csrf.Token(context))
	jsonResponse, _ := json.Marshal(response)
	log.Print(context.Write(jsonResponse))
}

// AddDirective adds directive dynamically
func (g *GraphQL) AddDirective(name string, directive interface{}) {
	funcType := reflect.ValueOf(directive)

	if funcType.Kind() != reflect.Func {
		panic("Directive should be type of function")
	}

	newDirective := new(tools.SchemaDirectiveVisitor)

	switch directive.(type) {
	case VisitSchema:
		newDirective.VisitSchema = directive.(VisitSchema)
	case VisitScalar:
		newDirective.VisitScalar = directive.(VisitScalar)
	case VisitObject:
		newDirective.VisitObject = directive.(VisitObject)
	case VisitFieldDefinition:
		newDirective.VisitFieldDefinition = directive.(VisitFieldDefinition)
	case VisitArgumentDefinition:
		newDirective.VisitArgumentDefinition = directive.(VisitArgumentDefinition)
	case VisitInterface:
		newDirective.VisitInterface = directive.(VisitInterface)
	case VisitUnion:
		newDirective.VisitUnion = directive.(VisitUnion)
	case VisitEnum:
		newDirective.VisitEnum = directive.(VisitEnum)
	case VisitEnumValue:
		newDirective.VisitEnumValue = directive.(VisitEnumValue)
	case VisitInputObject:
		newDirective.VisitInputObject = directive.(VisitInputObject)
	case VisitInputFieldDefinition:
		newDirective.VisitInputFieldDefinition = directive.(VisitInputFieldDefinition)
	default:
		panic("Invalid directive definition")
	}

	g.directives[name] = newDirective
}

// DropDirective drops directive dynamically
func (g *GraphQL) DropDirective(name string) {
	if !g.DirectiveExists(name) {
		return
	}

	delete(g.directives, name)
}

// DirectiveExists checks, if directive exists
func (g *GraphQL) DirectiveExists(name string) bool {
	for directiveName := range g.directives {
		if name == directiveName {
			return true
		}
	}

	return false
}
