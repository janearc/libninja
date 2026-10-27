package temple

// edgeOne and edgeTwo are one body of exactly smallestBlade nodes, forty,
// forged twice: a blade.
func edgeOne(times int) int {
	total := 0
	for step := 0; step < times; step++ {
		total += step * 2
	}
	if total > 10 {
		total -= 3
	}
	total *= times
	total = total - times
	return (total + times)
}

func edgeTwo(times int) int {
	total := 0
	for step := 0; step < times; step++ {
		total += step * 2
	}
	if total > 10 {
		total -= 3
	}
	total *= times
	total = total - times
	return (total + times)
}

// nearOne and nearTwo are one body of thirty-nine nodes, forged twice:
// too short to be a blade.
func nearOne(times int) int {
	total := 0
	for step := 0; step < times; step++ {
		total += step * 2
	}
	if total > 10 {
		total -= 3
	}
	total *= times
	total = total - times
	return total + times
}

func nearTwo(times int) int {
	total := 0
	for step := 0; step < times; step++ {
		total += step * 2
	}
	if total > 10 {
		total -= 3
	}
	total *= times
	total = total - times
	return total + times
}

// orderOne and orderTwo differ only in which local each step moves, so
// they are two bodies, not one forged twice.
func orderOne(times int) int {
	first, second := times*2, times+7
	for step := 0; step < times; step++ {
		first += step
		second -= step
	}
	first *= times
	if first > second {
		return first - second
	}
	return second
}

func orderTwo(times int) int {
	first, second := times*2, times+7
	for step := 0; step < times; step++ {
		second += step
		first -= step
	}
	first *= times
	if first > second {
		return first - second
	}
	return second
}
