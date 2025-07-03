package dedent_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	. "github.com/onsi/gomega"
)

func TestPipe(t *testing.T) {
	g := NewWithT(t)

	a := dedent.Pipe(`
		|package main
	`)
	g.Expect(a).To(Equal("package main\n"))

	b := dedent.Pipe(`
	`)
	g.Expect(b).To(Equal(""))

	c := dedent.Pipe(` `)
	g.Expect(c).To(Equal(""))

	d := dedent.Pipe(`
		|package main`)
	g.Expect(d).To(Equal("package main"))

	e := dedent.Pipe(` |package main
	`)
	g.Expect(e).To(Equal("package main\n"))

	f := dedent.Pipe(` |package main`)
	g.Expect(f).To(Equal("package main"))
}
