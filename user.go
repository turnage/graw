package graw

// User defines behaviors of a Reddit user account, which bots will have access
// to. These behaviors are not expected to follow through immediately; they will
// be put in queue and executed as soon as possible.
type User interface {}
