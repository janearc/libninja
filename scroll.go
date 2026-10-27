package libninja

import "slices"

// scroll is one rule. It holds the name a vow uses for it and its title.
// It holds its law in the master's words, and the way out. It holds its
// line from the film The Ninja Showdown, with the line's timecode. The
// film is from 1988 and in the public domain.
type scroll struct {
	name, title, law, way, filmed string
}

// taught is where a scroll comes in the teaching: its place in scrolls.
func (each scroll) taught() int {
	return slices.Index(scrolls(), each)
}

var gate = scroll{"gate", "the scroll of the open gate",
	"A ninja does not leave his sword at the temple door.",
	"Close the gate. Make the field private, let the type be made by " +
		"a constructor that may refuse, and let those who need the " +
		"value ask a method for it.",
	`00:14:22  "If I fail, never reveal your Silver Ninja identity."`}

var stranger = scroll{"stranger", "the scroll of the stranger's house",
	"A ninja does not build in another's house.",
	"Ask the package that owns the type for one, through its " +
		"constructor, and let it refuse what it cannot be.",
	`00:14:29  "This is my fate, and mine alone."`}

var well = scroll{"well", "the scroll of the open well",
	"A ninja does not leave his well uncovered, for anyone may " +
		"poison it.",
	"Make it a function that returns the value. Then no one can " +
		"change it behind your back.",
	`01:06:22  "He must always keep a clear mind."`}

var waste = scroll{"waste", "the scroll of wasted power",
	"A ninja shows only the weapons he must carry.",
	"No one outside the package calls it. Make it private.",
	`00:39:25  "He must never waste his powers."`}

var practice = scroll{"practice", "the scroll of the practice sword",
	"A ninja trains with the weapons he will carry, and no others.",
	"Only tests call it from outside. Make it private, and test " +
		"through what the package offers. If that is hard, the " +
		"difficulty is the lesson: it is telling you the interface. " +
		"A " +
		"package named for its tests, like httptest, is no practice " +
		"sword: only tests may import it, and some test must.",
	`00:12:37  "Master, it's practice time."`}

var dojo = scroll{"dojo", "the scroll of the dojo wall",
	"A ninja is tested in the outside world, not behind the dojo wall.",
	"Move the test to the package's _test package, and test through " +
		"the real interface.",
	`00:21:49  "Tomorrow, you all will start your training in the ` +
		`outside world."`}

var decoy = scroll{"decoy", "the scroll of the decoy",
	"A ninja trains against the decoy he was given, not one he carves.",
	"You have mocks. If you can't get it done with the mocks, maybe you " +
		"have a design problem. Hand the code an interface, and the " +
		"test the mock mockgen made for it.",
	""}

var blade = scroll{"blade", "the scroll of the borrowed blade",
	"A ninja forges a blade once.",
	"Find what each place needs, build it once where it belongs, and " +
		"have each place ask for it. Then put the copies down.",
	`00:39:08  "Violence is not the only way to solve a problem for a ` +
		`ninja. Patience is sometimes even stronger."`}

var stance = scroll{"stance", "the scroll of the stance",
	"Now, you must do it smoothly and with grace.",
	"Break it into steps that each fit on one screen, about fifty " +
		"lines.",
	`00:12:55  "Correct your stance."`}

var truename = scroll{"name", "the scroll of the true name",
	"A ninja knows each thing by its name.",
	"Give it a word. Only i, k, x, y and z may stand alone, and a test's " +
		"t, typed *testing.T or rapid's *T in a test file. A " +
		"colour's channels are libtheme's to name; ask it for the " +
		"colour.",
	`00:22:01  "Tony, stay here."`}

var trueAge = scroll{"age", "the scroll of the true age",
	"A ninja knows when he was forged, and says so when asked.",
	"Declare build and built, the stamp game passes, and make the " +
		"main a pzm tool, which answers --age with them, so anyone " +
		"can ask a binary which commit it is.",
	""}

var buried = scroll{"buried", "the scroll of the buried path",
	"A ninja does not bury his things where only he can dig.",
	"A path written into the code is true on one machine. Take it " +
		"from the config, or from the directory the tool is given.",
	""}

var editorless = scroll{"editor", "the scroll of the unmarked house",
	"A ninja's house says how it is kept, so every guest keeps it alike.",
	"Tell the editors how the project is kept: an .editorconfig " +
		"that covers its go files, or the project settings of the " +
		"editors it is written in, kept in the repository.",
	""}

var nested = scroll{"nest", "the scroll of the nested blades",
	"A ninja draws one blade at a time.",
	"Give the inner call's result a name, on a line of its own, and " +
		"hand the name on. Two calls deep is the most: f(g(x)).",
	""}

var nameless = scroll{"literal", "the scroll of the nameless blade",
	"A ninja names every blade he carries.",
	"Give it a name: a function, or a method on a type that holds " +
		"what it needs. A goroutine runs a named method; an api that " +
		"takes a function is handed a named one.",
	""}

var place = scroll{"index", "the scroll of the numbered place",
	"A ninja knows what is kept in each place, not where it is.",
	"Use the stack's own ways for a stack, a struct with named fields " +
		"for a fixed set of parts, and where a place must be named, " +
		"declare it once, with why.",
	""}

var bare = scroll{"number", "the scroll of the bare number",
	"A ninja says what each number is and where it came from.",
	"Name it in a declaration that says why, or say what it does with " +
		"a small named function written once: rest(share), not " +
		"1 - share; positive(x), not x > 0. 0 and 1 are numbers too.",
	""}

var tableInCode = scroll{"table", "the scroll of the table in code",
	"A ninja keeps his records in the archive, not on his sword.",
	"Put the table in a data file, pinned by its source if it was " +
		"measured, and read it; a name does not make a table into " +
		"anything but data in code.",
	""}

var kernelHeld = scroll{"kernel", "the scroll of the kernel's own numbers",
	"A ninja's blade takes its measure from the forge, not from itself.",
	"Declare the number once in go, named and saying why, and hand it " +
		"to the kernel through the prelude, so go and the kernel can " +
		"never hold two of it.",
	""}

var sourceless = scroll{"source", "the scroll of the sourceless number",
	"A ninja knows where each of his numbers came from.",
	"A measured number goes in the research data, pulled from its " +
		"source; a picked number goes in the picked numbers, with " +
		"who chose it. Code asks for both by name. A number true by " +
		"definition stays, and its comment says it is defined.",
	""}

var crowdedVows = scroll{"crowd", "the scroll of the crowded vows",
	"A ninja who vows at every step has stopped training.",
	"A ninja has pride and does not run from a fight. Where a scroll " +
		"itself is wrong, bring it to the master, so it can be " +
		"rewritten. A vow does not pay for a vow.",
	""}

var emptyVow = scroll{"vow", "the scroll of the empty vow",
	"A vow without words is no vow.",
	"Name the scroll you break and say why, on the vow's own line.",
	`00:14:33  "Promise me you will never forget."`}

// scrolls is every scroll, in the order they are taught.
func scrolls() []scroll {
	return []scroll{gate, stranger, well, waste, practice, dojo, decoy,
		blade, stance, truename, trueAge, buried, editorless, nested,
		nameless, place, bare, tableInCode, kernelHeld, sourceless,
		emptyVow, crowdedVows}
}

// named is the scroll a vow names, or false.
func named(name string) (scroll, bool) {
	for _, each := range scrolls() {
		if each.name == name {
			return each, true
		}
	}
	return scroll{}, false
}
