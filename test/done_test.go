package test

import (
	"strings"
	"testing"

	"bitbucket.org/shu_go/gotwant"
	"bitbucket.org/shu_go/scope"
)

func TestDone(t *testing.T) {
	t.Run("NoArgs", func(t *testing.T) {
		result := ""
		s := scope.New(
			func() { result += "begin" },
			func() { result += "end" },
		)
		func() {
			defer s.Done()()
			result += "body"

			gotwant.Test(t, result, "beginbody")
		}()
		gotwant.Test(t, result, "beginbodyend")
	})
	t.Run("One", func(t *testing.T) {
		result := ""
		s := scope.New(
			func(s string) { result += "begin(" + s + ")" },
			func() { result += "end" },
		)
		func() {
			defer s.Done("abc")()
			result += "body"

			gotwant.Test(t, result, "begin(abc)body")
		}()
		gotwant.Test(t, result, "begin(abc)bodyend")
	})
	t.Run("Two", func(t *testing.T) {
		result := ""
		s := scope.New(
			func(s string, n int) { result += "begin(" + strings.Repeat(s, n) + ")" },
			func() { result += "end" },
		)
		func() {
			defer s.Done("abc", 2)()
			result += "body"

			gotwant.Test(t, result, "begin(abcabc)body")
		}()
		gotwant.Test(t, result, "begin(abcabc)bodyend")
	})
	t.Run("Slice", func(t *testing.T) {
		result := ""
		s := scope.New(
			func(s []string) { result += "begin(" + strings.Join(s, ", ") + ")" },
			func() { result += "end" },
		)
		func() {
			defer s.Done([]string{"abc", "def"})()
			result += "body"

			gotwant.Test(t, result, "begin(abc, def)body")
		}()
		gotwant.Test(t, result, "begin(abc, def)bodyend")
	})
}
