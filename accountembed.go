package graw

import (
	"github.com/turnage/graw/internal/api/account"
)

// Account implements account.Account. Embed it in your bot, and graw will
// initialize it prior to the run; then your bot can call the methods of Account
// to make logged-in actions.
type Account struct {
	account.Account
}

// Casting to this interface is how graw identifies when a bot has embedded the
// Account struct, and how it initializes the implementation.
type accountEmbed interface {
	// grawSetImpl sets the implementation of account.Account for the embed.
	grawSetImpl(account.Account)
}

func (a *Account) grawSetImpl(acc account.Account) {
	a.Account = acc
}
