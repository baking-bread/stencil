package funcs

func Length(input []interface{}) int {
	return len(input)
}

func Pick(input []interface{}, index int) string {

	// casting items from an interface to a string because yaml only provides interface types
	var array = make([]string, len(input))

	for index, value := range input {
		array[index] = value.(string)
	}

	return array[index]
}
