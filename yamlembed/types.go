package yamlembed

import (
	"strings"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"-"`
}

type Bar struct {
	I      int64 `yaml:"i,omitempty"`
	B      string
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty"`
	F      []any    `yaml:"f"`
}

func (b *Bar) MarshalYAML() (interface{}, error) {
	return struct {
		I  int64 `yaml:"i,omitempty"`
		B  string
		OI []string `yaml:"oi,omitempty"`
		F  []any    `yaml:"f,flow"`
	}{
		I:  b.I,
		B:  b.B,
		OI: b.OI,
		F:  b.F,
	}, nil
}

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var bar struct {
		I  int64 `yaml:"i,omitempty"`
		B  string
		OI []string `yaml:"oi,omitempty"`
		F  []any    `yaml:"f"`
	}
	err := unmarshal(&bar)
	if err != nil {
		return err
	}
	b.I = bar.I
	b.B = bar.B
	b.UpperB = strings.ToUpper(bar.B)
	b.OI = bar.OI
	b.F = bar.F
	return nil
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func (b *Baz) MarshalYAML() (interface{}, error) {
	return struct {
		Foo Foo `yaml:",inline"`
		Bar struct {
			I  int64 `yaml:"i,omitempty"`
			B  string
			OI []string `yaml:"oi,omitempty"`
			F  []any    `yaml:"f,flow"`
		} `yaml:",inline"`
	}{
		Foo: b.Foo,
		Bar: struct {
			I  int64 `yaml:"i,omitempty"`
			B  string
			OI []string `yaml:"oi,omitempty"`
			F  []any    `yaml:"f,flow"`
		}{
			I:  b.Bar.I,
			B:  b.Bar.B,
			OI: b.Bar.OI,
			F:  b.Bar.F,
		},
	}, nil
}

func (b *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var baz struct {
		Foo Foo `yaml:",inline"`
		Bar struct {
			I  int64 `yaml:"i,omitempty"`
			B  string
			OI []string `yaml:"oi,omitempty"`
			F  []any    `yaml:"f"`
		} `yaml:",inline"`
	}
	err := unmarshal(&baz)
	if err != nil {
		return err
	}
	b.Foo = baz.Foo
	b.Bar.I = baz.Bar.I
	b.Bar.B = baz.Bar.B
	b.Bar.UpperB = strings.ToUpper(baz.Bar.B)
	b.Bar.OI = baz.Bar.OI
	b.Bar.F = baz.Bar.F
	return nil
}
