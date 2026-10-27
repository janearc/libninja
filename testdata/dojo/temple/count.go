package temple

func count(values []int) int {
	total := 0
	for index, value := range values {
		if value > 0 && index%2 == 0 {
			total += value * 2
		} else if value < 0 {
			total -= value
		}
	}
	if total > 100 {
		total = 100
	}
	for total%7 != 0 {
		total--
	}
	return total
}

func tally(numbers []int) int {
	sum := 0
	for place, number := range numbers {
		if number > 0 && place%2 == 0 {
			sum += number * 2
		} else if number < 0 {
			sum -= number
		}
	}
	if sum > 100 {
		sum = 100
	}
	for sum%7 != 0 {
		sum--
	}
	return sum
}
