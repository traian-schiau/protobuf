// Protocol Buffers for Go with Gadgets
//
// Copyright (c) 2020, The GoGo Authors. All rights reserved.
// http://github.com/gogo/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// castwithfixup plugin generates a set of methods, useful in indentifying
// incompatibilities in massages containing at least one "casttypewith" field
//
//		MarshalJSONPB
//		UnmarshalJSONPB
//		MarshalText
//		UnmarshalText

package castwithfixup

import (
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type jsonpb struct {
	*generator.Generator
	generator.PluginImports
}

func NewJsonPB() *jsonpb {
	return &jsonpb{}
}

func (p *jsonpb) Name() string {
	return "castwithfixup"
}

func (p *jsonpb) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *jsonpb) messagNeedsFix(m *generator.Descriptor) bool {
	for _, f := range m.GetField() {
		if gogoproto.IsCastTypeWith(f) {
			return true
		}
	}
	return false
}

func (p *jsonpb) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	jsonpbPkg := p.PluginImports.NewImport("github.com/gogo/protobuf/jsonpb")
	for _, message := range file.Messages() {
		if !gogoproto.HasCastTypeWith(file.FileDescriptorProto, message.DescriptorProto) {
			continue
		}

		ccTypeName := generator.CamelCaseSlice(message.TypeName())

		p.P(`func (m *`, ccTypeName, `) MarshalJSONPB(*`, jsonpbPkg.Use(), `.Marshaler) ([]byte, error) {`)
		p.In()
		p.P(`return nil, fmt.Errorf("MarshalJSONPB is not supported by `, ccTypeName, `(casttypewith)")`)
		p.Out()
		p.P("}")
		p.P()

		p.P(`func (m *`, ccTypeName, `) UnmarshalJSONPB(*`, jsonpbPkg.Use(), `.Unmarshaler, []byte) error {`)
		p.In()
		p.P(`return fmt.Errorf("UnmarshalJSONPB is not supported by `, ccTypeName, `(casttypewith)")`)
		p.Out()
		p.P("}")
		p.P()

		p.P(`func (m *`, ccTypeName, `) MarshalText() ([]byte, error) {`)
		p.In()
		p.P(`return nil, fmt.Errorf("MarshalText is not supported by `, ccTypeName, `(casttypewith)")`)
		p.Out()
		p.P("}")
		p.P()

		p.P(`func (m *`, ccTypeName, `) UnmarshalText([]byte) error {`)
		p.In()
		p.P(`return fmt.Errorf("UnmarshalText is not supported by `, ccTypeName, `(casttypewith)")`)
		p.Out()
		p.P("}")
		p.P()
	}

}

func init() {
	generator.RegisterPlugin(NewJsonPB())
}
