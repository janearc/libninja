package libninja_test

import (
	"fmt"
	"strings"
	"testing"
)

// thing is the struct the feint tests hand back, in a house's file.
const thing = "// Thing is a thing.\ntype Thing struct{ ID string }\n\n"

// emptied is how the feint says an empty struct came back with an
// error, from the type's name on.
const emptied = " is an empty value handed back with an error. " +
	"A caller may use it unchecked. On a refusal return nothing usable: " +
	"a pointer, nil; a slice, nil; or change the result."

// settled is how the feint says a struct with fields set came back
// with an error, from the type's name on.
const settled = " has fields set and is handed back with an error. " +
	"It looks like real data, so a caller who skips the check uses it. " +
	"Check the error first and build only on success: " +
	"if err != nil { return nil, err }. " +
	"On a refusal return nothing usable: " +
	"a pointer, nil; a slice, nil; or change the result."

// once is how many hits one bad return makes, defined so.
const once = 1

// twice is how many hits two bad results in one return make, defined
// so.
const twice = 2

// feintsIn is what the feint scroll saw in a house with one more
// file, and the reasons of the vows it carries for the scroll.
func feintsIn(test *testing.T, body string) ([]string, []string) {
	test.Helper()
	dir := house(test, ".", map[string]string{
		".editorconfig": "root = true\n\n[*.go]\nindent_style = tab\n",
		"more.go":       "package house\n\n" + body,
	})
	_, facts := thrown(test, dir)
	var hits, vowed []string
	for _, each := range facts["hits"] {
		if each["scroll"] == "feint" && each["file"] == "more.go" {
			hits = append(hits, each["seen"].(string))
		}
	}
	for _, each := range facts["vows"] {
		if each["scroll"] == "feint" && each["file"] == "more.go" {
			vowed = append(vowed, each["seen"].(string))
		}
	}
	return hits, vowed
}

// feintHits is only the hits of feintsIn.
func feintHits(test *testing.T, body string) []string {
	test.Helper()
	hits, _ := feintsIn(test, body)
	return hits
}

// wantOne checks that the one hit said exactly what is wanted.
func wantOne(test *testing.T, hits []string, want string) {
	test.Helper()
	if strings.Join(hits, "|") != want {
		test.Errorf("hits %q, want %q", hits, want)
	}
}

func TestAnEmptyStructWithAnErrorIsAHit(test *testing.T) {
	body := thing +
		"func refuse(err error) (Thing, error) { return Thing{}, err }\n"
	wantOne(test, feintHits(test, body), "Thing{}"+emptied)
}

func TestAnEmptyStructOfAnotherPackageIsAHit(test *testing.T) {
	body := "import \"time\"\n\n" +
		"func refuse(err error) (time.Time, error) {\n" +
		"\treturn time.Time{}, err\n}\n"
	wantOne(test, feintHits(test, body), "time.Time{}"+emptied)
}

func TestAnEmptyGenericStructWithAnErrorIsAHit(test *testing.T) {
	body := "// Schema is a schema.\n" +
		"type Schema[T any] struct{ Name string }\n\n" +
		"func refuse[T any](err error) (Schema[T], error) {\n" +
		"\treturn Schema[T]{}, err\n}\n"
	wantOne(test, feintHits(test, body), "Schema[T]{}"+emptied)
}

func TestAStructNamedOnlyByItsFieldsIsFoundByItsType(test *testing.T) {
	body := "func refuse(err error) (struct{ ID string }, error) {\n" +
		"\treturn struct{ ID string }{}, err\n}\n"
	wantOne(test, feintHits(test, body),
		"struct{ID string}{}"+emptied)
}

func TestAStructBehindAnAliasIsAStruct(test *testing.T) {
	body := thing + "type Alias = Thing\n\n" +
		"func refuse(err error) (Alias, error) { return Alias{}, err }\n"
	wantOne(test, feintHits(test, body), "Thing{}"+emptied)
}

func TestAStructBehindANamedTypeIsAStruct(test *testing.T) {
	body := thing + "type Made Thing\n\n" +
		"func refuse(err error) (Made, error) { return Made{}, err }\n"
	wantOne(test, feintHits(test, body), "Made{}"+emptied)
}

func TestAStructOnlyDeclaredIsAHit(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\tvar made Thing\n\treturn made, err\n}\n"
	wantOne(test, feintHits(test, body),
		"made is a Thing declared and never set, handed back with an "+
			"error. A caller may use it unchecked. On a refusal return "+
			"nothing usable: a pointer, nil; a slice, nil; or change "+
			"the result.")
}

func TestAStructDeclaredInAGroupAndNeverSetIsAHit(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\tvar (\n\t\tmade Thing\n\t\tother Thing\n\t)\n" +
		"\tother.ID = \"set\"\n\t_ = other\n\treturn made, err\n}\n"
	if hits := feintHits(test, body); len(hits) != once {
		test.Errorf("one of two declared, one set: %q", hits)
	}
}

func TestAStructWithFieldsSetIsAHit(test *testing.T) {
	body := thing + "func refuse(id string, err error) (Thing, error) {\n" +
		"\treturn Thing{ID: id}, err\n}\n"
	wantOne(test, feintHits(test, body), "Thing{...}"+settled)
}

func TestAStructWithPositionalFieldsIsAHit(test *testing.T) {
	body := thing + "func refuse(id string, err error) (Thing, error) {\n" +
		"\treturn Thing{id}, err\n}\n"
	wantOne(test, feintHits(test, body), "Thing{...}"+settled)
}

func TestTheStructMayStandAmongOtherResults(test *testing.T) {
	body := thing + "func refuse(err error) (string, Thing, bool, error) {\n" +
		"\treturn \"\", Thing{}, false, err\n}\n"
	wantOne(test, feintHits(test, body), "Thing{}"+emptied)
}

func TestTwoStructsAreTwoHits(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, Thing, error) {\n" +
		"\treturn Thing{}, Thing{}, err\n}\n"
	if hits := feintHits(test, body); len(hits) != twice {
		test.Errorf("two empty structs: %q", hits)
	}
}

func TestAFunctionLiteralIsJudgedByItsOwnResults(test *testing.T) {
	body := thing + "func outer(err error) (*Thing, error) {\n" +
		"\tinner := func() (Thing, error) { return Thing{}, err }\n" +
		"\t_, _ = inner()\n\treturn nil, err\n}\n"
	wantOne(test, feintHits(test, body), "Thing{}"+emptied)
}

func TestAFunctionLiteralWithNoErrorDoesNotLendItsResults(test *testing.T) {
	body := thing + "func outer(err error) (Thing, error) {\n" +
		"\tinner := func() Thing { return Thing{} }\n" +
		"\treturn inner(), err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a literal that returns no error: %q", hits)
	}
}

func TestAMethodIsJudgedLikeAFunction(test *testing.T) {
	body := thing + "func (each Thing) Refuse(err error) (Thing, error) {\n" +
		"\treturn Thing{}, err\n}\n"
	wantOne(test, feintHits(test, body), "Thing{}"+emptied)
}

func TestBasicTypesWithAnErrorAreFine(test *testing.T) {
	body := "func number(err error) (int, error) { return 0, err }\n\n" +
		"func text(err error) (string, error) { return \"\", err }\n\n" +
		"func truth(err error) (bool, error) { return false, err }\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("basic types: %q", hits)
	}
}

func TestNilsWithAnErrorAreTheFix(test *testing.T) {
	body := thing + "type Shape interface{ Area() float64 }\n\n" +
		"func pointer(err error) (*Thing, error) { return nil, err }\n\n" +
		"func slice(err error) ([]Thing, error) { return nil, err }\n\n" +
		"func table(err error) (map[string]Thing, error) {\n" +
		"\treturn nil, err\n}\n\n" +
		"func shape(err error) (Shape, error) { return nil, err }\n\n" +
		"func verb(err error) (func(), error) { return nil, err }\n\n" +
		"func pipe(err error) (chan Thing, error) { return nil, err }\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("nil results: %q", hits)
	}
}

func TestAPointerToAnEmptyStructIsNotAStruct(test *testing.T) {
	body := thing + "func refuse(err error) (*Thing, error) {\n" +
		"\treturn &Thing{}, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a pointer: %q", hits)
	}
}

func TestANilErrorIsNotARefusal(test *testing.T) {
	body := thing + "func fine() (Thing, error) { return Thing{}, nil }\n\n" +
		"func filled(id string) (Thing, error) {\n" +
		"\treturn Thing{ID: id}, nil\n}\n\n" +
		"func kept() (Thing, error) {\n" +
		"\tvar made Thing\n\treturn made, nil\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("nil error: %q", hits)
	}
}

func TestAStructWithoutAnErrorIsFine(test *testing.T) {
	body := thing + "func fine() Thing { return Thing{} }\n\n" +
		"func pair() (Thing, bool) { return Thing{}, false }\n\n" +
		"func backwards(err error) (error, Thing) {\n" +
		"\treturn err, Thing{}\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("no error last: %q", hits)
	}
}

func TestAnErrorTypeThatIsNotErrorIsLeftAlone(test *testing.T) {
	body := thing + "// Failure is a failure.\ntype Failure struct{}\n\n" +
		"// Error says it.\nfunc (Failure) Error() string { return \"\" }\n\n" +
		"func refuse(failure Failure) (Thing, Failure) {\n" +
		"\treturn Thing{}, failure\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a result that is not the error type: %q", hits)
	}
}

func TestAValueFromACallIsLeftAlone(test *testing.T) {
	body := thing + "func build() (Thing, error) { return Thing{}, nil }\n\n" +
		"func again() (Thing, error) { return build() }\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a call handed on: %q", hits)
	}
}

func TestANakedReturnIsLeftAlone(test *testing.T) {
	body := thing + "func named(err error) (made Thing, failed error) {\n" +
		"\tfailed = err\n\treturn\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a naked return: %q", hits)
	}
}

func TestAZeroOfATypeParameterIsLeftAlone(test *testing.T) {
	body := "func refuse[T any](err error) (T, error) {\n" +
		"\tvar zero T\n\treturn zero, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a type parameter: %q", hits)
	}
}

func TestADeclaredStructThenSetIsLeftAlone(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\tvar made Thing\n\tmade = Thing{ID: \"a\"}\n" +
		"\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("assigned whole: %q", hits)
	}
}

func TestADeclaredStructWithAFieldSetIsLeftAlone(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\tvar made Thing\n\tmade.ID = \"a\"\n\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a field assigned: %q", hits)
	}
}

func TestADeclaredStructHandedByAddressIsLeftAlone(test *testing.T) {
	body := thing + "func fill(into *Thing) { into.ID = \"a\" }\n\n" +
		"func refuse(err error) (Thing, error) {\n" +
		"\tvar made Thing\n\tfill(&made)\n\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("filled through its address: %q", hits)
	}
}

func TestADeclaredStructWithAMethodCalledIsLeftAlone(test *testing.T) {
	body := thing + "func (into *Thing) Fill() { into.ID = \"a\" }\n\n" +
		"func refuse(err error) (Thing, error) {\n" +
		"\tvar made Thing\n\tmade.Fill()\n\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("filled by a method: %q", hits)
	}
}

func TestADeclaredStructSetInAnEarlierBranchIsLeftAlone(test *testing.T) {
	body := thing + "func refuse(ok bool, err error) (Thing, error) {\n" +
		"\tvar made Thing\n\tif ok {\n\t\tmade = Thing{ID: \"a\"}\n\t}\n" +
		"\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("set in a branch: %q", hits)
	}
}

func TestADeclaredStructSetOnlyAfterTheReturnStillHits(test *testing.T) {
	body := thing + "func refuse(ok bool, err error) (Thing, error) {\n" +
		"\tvar made Thing\n\tif ok {\n\t\treturn made, err\n\t}\n" +
		"\tmade = Thing{ID: \"a\"}\n\treturn made, nil\n}\n"
	if hits := feintHits(test, body); len(hits) != once {
		test.Errorf("set after the refusal: %q", hits)
	}
}

func TestAVariableGivenAValueAtItsDeclarationIsLeftAlone(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\tvar made = Thing{ID: \"a\"}\n\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("declared with a value: %q", hits)
	}
}

func TestAParameterIsLeftAlone(test *testing.T) {
	body := thing + "func refuse(made Thing, err error) (Thing, error) {\n" +
		"\treturn made, err\n}\n"
	if hits := feintHits(test, body); hits != nil {
		test.Errorf("a parameter handed back: %q", hits)
	}
}

func TestAVowAboveTheReturnCoversTheHit(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\t//ninja:vow feint the wire wants the struct whole\n" +
		"\treturn Thing{}, err\n}\n"
	hits, vowed := feintsIn(test, body)
	if hits != nil {
		test.Errorf("a vowed return was still a hit: %q", hits)
	}
	wantOne(test, vowed, "Thing{}"+emptied)
}

func TestAVowAtTheEndOfTheLineCoversTheHit(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\treturn Thing{}, err //ninja:vow feint the wire wants it\n}\n"
	hits, vowed := feintsIn(test, body)
	if hits != nil {
		test.Errorf("a vowed return was still a hit: %q", hits)
	}
	wantOne(test, vowed, "Thing{}"+emptied)
}

func TestAVowOverAReturnThatIsFineCoversNothing(test *testing.T) {
	body := thing + "func refuse(err error) (*Thing, error) {\n" +
		"\t//ninja:vow feint the wire wants the struct whole\n" +
		"\treturn nil, err\n}\n"
	seen := hitsIn(test, "vow", body)
	if !strings.Contains(strings.Join(seen, "|"), "covers nothing") {
		test.Errorf("a vow over a return that breaks nothing: %q", seen)
	}
}

func TestTheDojoVowsAndLeavesAloneWhatItShould(test *testing.T) {
	_, facts := thrown(test, "testdata/dojo")
	file := "temple/feint.go"
	if !holds(facts["vows"], "feint", file, "feintThing{} is an") {
		test.Error("the vowed empty return is not among the vows")
	}
	var lines []string
	for _, each := range facts["hits"] {
		if each["scroll"] == "feint" && each["file"] == file {
			lines = append(lines, fmt.Sprint(each["line"]))
		}
	}
	if got := strings.Join(lines, " "); got != dojoFeints {
		test.Errorf("the dojo's feints are on lines %q, want %q",
			got, dojoFeints)
	}
}

// dojoFeints are the lines of testdata/dojo/temple/feint.go that
// break the scroll: the empty literal, the declared and never set
// variable, and the literal with fields set. The vowed return, the
// basic type, the nils, the nil error and the variable set are not
// among them.
const dojoFeints = "11 16 20"

func TestAVowForAnotherScrollDoesNotCoverIt(test *testing.T) {
	body := thing + "func refuse(err error) (Thing, error) {\n" +
		"\t//ninja:vow gate this is another scroll altogether\n" +
		"\treturn Thing{}, err\n}\n"
	if hits := feintHits(test, body); len(hits) != once {
		test.Errorf("a vow for the wrong scroll: %q", hits)
	}
}
