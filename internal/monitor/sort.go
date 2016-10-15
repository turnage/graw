package monitor

type byCreationTime struct {
	things []redditThing
}

func (b byCreationTime) Len() int { return len(b.things) }

// Returns true if things[i] should precede things[j], which is true if it is a
// younger redditThing (more recent creation time).
func (b byCreationTime) Less(i, j int) bool {
	return b.things[i].GetCreatedUtc() > b.things[j].GetCreatedUtc()
}

func (b byCreationTime) Swap(i, j int) {
	thing := b.things[i]
	b.things[i] = b.things[j]
	b.things[j] = thing
}
