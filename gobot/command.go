package gobot

// Command represents a command in a discord bot.
type Command struct {
	Name        string // The lowercase name of the command.
	Description string // The description to use in help.
	Category    string // A virtual group to use when sorting these commands in a helper.

	// The checks to run. If you wish to send an error message, do so in the
	// check itself before returning false.
	Checks []func(context *Context) bool
	Runner func(context *Context) // The command runner.
}
