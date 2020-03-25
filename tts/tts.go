package tts

import (
	"bytes"
	"fmt"
	"sort"

	pgs "github.com/lyft/protoc-gen-star"
)

// New tts plugin
func New() pgs.Module {
	return &tts{
		ModuleBase: &pgs.ModuleBase{},
		pkgs:       map[pgs.Package]*protoPackage{},
	}
}

type tts struct {
	*pgs.ModuleBase
	pkgs map[pgs.Package]*protoPackage
}

// Name is the identifier used to identify the module. This value is
// automatically attached to the BuildContext associated with the ModuleBase.
func (t *tts) Name() string { return "tts" }

type protoPackage struct {
	Name     string
	Imports  []*importData
	Services []*serviceData
	Messages []*messageData
	Enums    []*enumData
}

func (p *protoPackage) sort() {
	p.sortImports()
	p.sortMessages()
	p.sortEnums()
}

func (p *protoPackage) sortImports() {
	sort.Slice(p.Imports, func(i, j int) bool {
		return p.Imports[i].Name < p.Imports[j].Name
	})
	for _, imp := range p.Imports {
		sort.Slice(imp.Declarations, func(i, j int) bool {
			return imp.Declarations[i] < imp.Declarations[j]
		})
	}
}

func (p *protoPackage) sortMessages() {
	sort.Slice(p.Messages, func(i, j int) bool {
		return p.Messages[i].Name < p.Messages[j].Name
	})
}

func (p *protoPackage) sortEnums() {
	sort.Slice(p.Enums, func(i, j int) bool {
		return p.Enums[i].Name < p.Enums[j].Name
	})
}

func (t *tts) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		pkg := t.addPackage(f.Package())
		for _, s := range f.Services() {
			svc := t.createService(s)
			pkg.Services = append(pkg.Services, svc)
		}
	}
	t.AddGeneratorFile("twirp.ts", twirpFile)
	for _, pkg := range t.pkgs {
		pkg.sort()
		var buf bytes.Buffer
		err := tmpl.Execute(&buf, pkg)
		if err != nil {
			panic(err)
		}
		t.AddGeneratorFile(pkg.Name+".ts", buf.String())
	}
	return t.Artifacts()
}

func (t *tts) addPackage(p pgs.Package) *protoPackage {
	pkgName := p.ProtoName().String()
	if pkg, ok := t.pkgs[p]; ok {
		return pkg
	}
	t.pkgs[p] = &protoPackage{Name: pkgName}
	return t.pkgs[p]
}

func (t *tts) createService(s pgs.Service) *serviceData {
	var sd serviceData
	sd.Name = s.Name().String()
	sd.Doc = getDoc(s.SourceCodeInfo().LeadingComments(), 0)
	sd.PathPrefix = fmt.Sprintf(
		"/twirp/%s.%s/",
		s.Package().ProtoName().String(),
		sd.Name,
	)
	for _, m := range s.Methods() {
		t.visitMessage(m.Package(), m.Input(), true)
		t.visitMessage(m.Package(), m.Output(), false)

		md := createMethod(m)
		sd.Methods = append(sd.Methods, md)
	}
	return &sd
}

func (t *tts) visitMessage(from pgs.Package, m pgs.Message, optional bool) {
	if from != m.Package() {
		imp := t.createImportForPackage(from, m.Package())
		t.addDeclarationForImport(imp, m.Name().String())
		t.addPackage(m.Package())
		t.visitMessage(m.Package(), m, optional)
	}
	if t.messageVisited(m) {
		return
	}
	msg := &messageData{
		Name:     m.Name().String(),
		Doc:      getDoc(m.SourceCodeInfo().LeadingComments(), 0),
		Optional: optional,
	}
	t.pkgs[m.Package()].Messages = append(t.pkgs[m.Package()].Messages, msg)
	for _, f := range m.Fields() {
		mf := createField(f)
		switch {
		case mf.IsEnum:
			t.visitEnum(m.Package(), pgsEnumFromField(f))
		case mf.IsMessage:
			t.visitMessage(m.Package(), pgsMsgFromField(f), optional)
		}
		msg.Fields = append(msg.Fields, mf)
	}
}

func (t *tts) visitEnum(from pgs.Package, e pgs.Enum) {
	if from != e.Package() {
		imp := t.createImportForPackage(from, e.Package())
		t.addDeclarationForImport(imp, e.Name().String())
		t.addPackage(e.Package())
		t.visitEnum(e.Package(), e)
	}
	if t.enumVisited(e) {
		return
	}
	var ed enumData
	ed.Name = e.Name().String()
	for _, v := range e.Values() {
		ed.Values = append(ed.Values, v.Name().String())
	}
	t.pkgs[e.Package()].Enums = append(t.pkgs[e.Package()].Enums, &ed)
}

func pgsMsgFromField(f pgs.Field) pgs.Message {
	if f.Type().IsRepeated() {
		return f.Type().Element().Embed()
	}
	return f.Type().Embed()
}

func pgsEnumFromField(f pgs.Field) pgs.Enum {
	if f.Type().IsRepeated() {
		return f.Type().Element().Enum()
	}
	return f.Type().Enum()
}

func (t *tts) messageVisited(m pgs.Message) bool {
	for _, visited := range t.pkgs[m.Package()].Messages {
		if visited.Name == m.Name().String() {
			return true
		}
	}
	return false
}

func (t *tts) enumVisited(e pgs.Enum) bool {
	for _, visited := range t.pkgs[e.Package()].Enums {
		if visited.Name == e.Name().String() {
			return true
		}
	}
	return false
}

func (t *tts) createImportForPackage(importer, imported pgs.Package) *importData {
	if imp, ok := t.hasImport(importer, imported); ok {
		return imp
	}
	imp := &importData{Name: imported.ProtoName().String()}
	t.pkgs[importer].Imports = append(t.pkgs[importer].Imports, imp)
	return imp
}

func (t *tts) hasImport(pkg pgs.Package, imported pgs.Package) (*importData, bool) {
	importerPackage := t.pkgs[pkg]
	for _, imp := range importerPackage.Imports {
		if imp.Name == imported.ProtoName().String() {
			return imp, true
		}
	}
	return nil, false
}

func (t *tts) addDeclarationForImport(imp *importData, decl string) {
	for _, declared := range imp.Declarations {
		if declared == decl {
			return
		}
	}
	imp.Declarations = append(imp.Declarations, decl)
}
