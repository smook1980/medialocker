package genswagger

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/utils"
)

//NewGenerator returns an initialized instance of a JavaScript Client Generator
func NewGenerator(options ...Option) *Generator {
	g := &Generator{}

	for _, option := range options {
		option(g)
	}

	return g
}

// Generator is the swagger code generator.
type Generator struct {
	API      *design.APIDefinition // The API definition
	OutDir   string                // Path to output directory
	genfiles []string              // Generated files
}

// Generate is the generator entry point called by the meta generator.
func Generate() (files []string, err error) {
	var outDir, ver string
	set := flag.NewFlagSet("swagger", flag.PanicOnError)
	set.StringVar(&outDir, "out", "", "")
	set.StringVar(&ver, "version", "", "")
	set.String("design", "", "")
	set.Bool("force", false, "")
	set.Bool("notest", false, "")
	set.Parse(os.Args[1:])

	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	g := &Generator{OutDir: outDir, API: design.Design}

	return g.Generate()
}

// Generate produces the skeleton main.
func (g *Generator) Generate() (_ []string, err error) {
	go utils.Catch(nil, func() { g.Cleanup() })

	defer func() {
		if err != nil {
			g.Cleanup()
		}
	}()

	s, err := New(g.API)
	if err != nil {
		return nil, err
	}

	swaggerDir := filepath.Join(g.OutDir, "swagger")
	os.RemoveAll(swaggerDir)
	if err = os.MkdirAll(swaggerDir, 0755); err != nil {
		return nil, err
	}
	g.genfiles = append(g.genfiles, swaggerDir)

	// JSON
	rawJSON, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	swaggerFile := filepath.Join(swaggerDir, "swagger.json")
	if err := ioutil.WriteFile(swaggerFile, rawJSON, 0644); err != nil {
		return nil, err
	}
	g.genfiles = append(g.genfiles, swaggerFile)

	// YAML
	var yamlSource interface{}
	if err = json.Unmarshal(rawJSON, &yamlSource); err != nil {
		return nil, err
	}

	rawYAML, err := yaml.Marshal(yamlSource)
	if err != nil {
		return nil, err
	}
	swaggerFile = filepath.Join(swaggerDir, "swagger.yaml")
	if err := ioutil.WriteFile(swaggerFile, rawYAML, 0644); err != nil {
		return nil, err
	}
	g.genfiles = append(g.genfiles, swaggerFile)

	return g.genfiles, nil
}

// Cleanup removes all the files generated by this generator during the last invokation of Generate.
func (g *Generator) Cleanup() {
	for _, f := range g.genfiles {
		os.Remove(f)
	}
	g.genfiles = nil
}