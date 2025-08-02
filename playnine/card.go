package playnine

// Card is just the value of the card. Possible values are 0-12 and -5.
type Card int8

const (
	// CardHoleInOne is a special card with a value of -5, and doesn't
	CardHoleInOne = Card(-5)
)
