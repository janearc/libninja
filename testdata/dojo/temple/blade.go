package temple

// Cutter cuts and has a name. Only Blade is one: steel cuts but has no
// name, so steel.Cut keeps Cutter's promise through Blade, which embeds
// it.
type Cutter interface {
	Cut() int
	Named() string
}

// steel cuts, and the blade carries it.
type steel struct{}

// Cut is promoted into Blade, and answers Cutter there.
func (steel) Cut() int { return 4 }

// Blade embeds steel and gives it a name.
type Blade struct{ steel }

// Named is Blade's own half of Cutter.
func (Blade) Named() string { return "blade" }

// Spare is a well no one outside draws from: the well's alone, not waste
// as well.
var Spare = 2
