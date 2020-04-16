package gobot

// Command represents a command in a discord bot.
type Command struct {
	Name        string                // The lowercase name of the command.
	Description string                // The description to use in help.
	Category    string                // A virtual group to use when sorting these commands in help.
	Runner      func(context Context) // The command runner.
}
