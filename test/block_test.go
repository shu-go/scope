package test

import (
	"strings"
	"testing"

	"bitbucket.org/shu_go/gotwant"
	"bitbucket.org/shu_go/scope"
)

func TestBlockBegin(t *testing.T) {
	t.Run("NoArgs", func(t *testing.T) {
		result := ""
		s := scope.New(
			func() { result += "begin" },
			func() { result += "end" },
		)
		s.Block(func() {
			result += "body"
		})
		gotwant.Test(t, result, "beginbodyend")
	})
	t.Run("One", func(t *testing.T) {
		result := ""
		s := scope.New(
			func(s string) { result += "begin(" + s + ")" },
			func() { result += "end" },
		)
		s.Block(func() {
			result += "body"
		}, "aaa")
		gotwant.Test(t, result, "begin(aaa)bodyend")
	})
	t.Run("Two", func(t *testing.T) {
		result := ""
		s := scope.New(
			func(s string, n int) { result += "begin(" + strings.Repeat(s, n) + ")" },
			func() { result += "end" },
		)
		s.Block(func() {
			result += "body"
		}, "abc", 2)
		gotwant.Test(t, result, "begin(abcabc)bodyend")
	})
	t.Run("Mismatch", func(t *testing.T) {
		defer func() {
			recover()
		}()

		result := ""
		s := scope.New(
			func(s string) { result += "begin(" + s + ")" },
			func() { result += "end" },
		)
		s.Block(func() {
			result += "body"
		}, "abc", 2)

		t.Fatal("not panicked")
	})
}

func TestBlockBeginToBody(t *testing.T) {
	t.Run("NoArgs", func(t *testing.T) {
		result := ""
		s := scope.New(
			func() { result += "begin" },
			func() { result += "end" },
		)
		s.Block(func() {
			result += "body"
		})
		gotwant.Test(t, result, "beginbodyend")
	})
	t.Run("OneString", func(t *testing.T) {
		result := ""
		s := scope.New(
			func() string { return "begin" },
			func() { result += "end" },
		)
		s.Block(func(beginstring string) {
			result += beginstring + "body"
		})
		gotwant.Test(t, result, "beginbodyend")
	})
	t.Run("Two", func(t *testing.T) {
		result := ""
		s := scope.New(
			func() (int, string) { return 3, "begin" },
			func() { result += "end" },
		)
		s.Block(func(beginint int, beginstring string) {
			result += strings.Repeat(beginstring, beginint) + "body"
		})
		gotwant.Test(t, result, "beginbeginbeginbodyend")
	})
	t.Run("OneWithoutErr", func(t *testing.T) {
		result := ""
		s := scope.New(
			func() (string, error) { return "begin", nil },
			func() { result += "end" },
		)
		s.Block(func(beginstring string) {
			result += beginstring + "body"
		})
		gotwant.Test(t, result, "beginbodyend")
	})
}
