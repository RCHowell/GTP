// R. C. Howell

// Adding the `ToCommand` method which serializes the struct to a GTP command
//
// This file is NOT generated and can be edited by hand

// Would be cool to implement these with reflect and Go generics

package gtp

import (
	"fmt"
)

// Command names -- I've tried using enums defined in the protobuf file, but this is cleaner and simpler.
const (
	ProtocolVersion = "protocol_version"
	Name            = "name"
	Version         = "version"
	IsKnownCommand  = "known_command"
	ListCommands    = "list_commands"
	Quit            = "quit"
	SetBoardSize    = "boardsize"
	ClearBoard      = "clear_board"
	SetKomi         = "komi"
	Play            = "play"
	GenMove         = "genmove"
	Undo            = "undo"
)

// This interface abstracts the generated protobuf structs to a single type
type Command interface {
	ToCommand() string
}

func (x *ProtocolVersionRequest) ToCommand() string {
	return noArgCommand(x.Id, ProtocolVersion)
}

func (x *NameRequest) ToCommand() string {
	return noArgCommand(x.Id, Name)
}

func (x *VersionRequest) ToCommand() string {
	return noArgCommand(x.Id, Version)
}

func (x *IsKnownCommandRequest) ToCommand() string {
	if x.Id == 0 {
		return fmt.Sprintf("%s %s%s", IsKnownCommand, x.Command, LF)
	}
	return fmt.Sprintf("%d %s %s%s", x.Id, IsKnownCommand, x.Command, LF)
}

func (x *ListCommandsRequest) ToCommand() string {
	return noArgCommand(x.Id, ListCommands)
}

func (x *QuitRequest) ToCommand() string {
	return noArgCommand(x.Id, Quit)
}

func (x *SetBoardSizeRequest) ToCommand() string {
	if x.Id == 0 {
		return fmt.Sprintf("%s %d%s", SetBoardSize, x.Size, LF)
	}
	return fmt.Sprintf("%d %s %d%s", x.Id, SetBoardSize, x.Size, LF)
}

func (x *ClearBoardRequest) ToCommand() string {
	return noArgCommand(x.Id, ClearBoard)
}

func (x *SetKomiRequest) ToCommand() string {
	if x.Id == 0 {
		return fmt.Sprintf("%s %f%s", SetKomi, x.Komi, LF)
	}
	return fmt.Sprintf("%d %s %f%s", x.Id, SetKomi, x.Komi, LF)
}

func (x *PlayRequest) ToCommand() string {
	if x.Id == 0 {
		return fmt.Sprintf("%s %s%s", Play, x.Move.GTPString(), LF)
	}
	return fmt.Sprintf("%d %s %s%s", x.Id, Play, x.Move.GTPString(), LF)
}

func (x *GenMoveRequest) ToCommand() string {
	if x.Id == 0 {
		return fmt.Sprintf("%s %s%s", GenMove, x.Color, LF)
	}
	return fmt.Sprintf("%d %s %s%s", x.Id, GenMove, x.Color, LF)
}

func (x *UndoRequest) ToCommand() string {
	return noArgCommand(x.Id, ClearBoard)
}

func noArgCommand(id int32, name string) string {
	if id == 0 {
		return fmt.Sprintf("%s%s", name, LF)
	}
	return fmt.Sprintf("%d %s%s", id, name, LF)
}

// TODO separate Place, Pass, Resign
func (x *Move) GTPString() string {
	switch x.Type {
	case Move_PLACE:
		return fmt.Sprintf("%s %s%d", x.Color, columnToLetter(x.Vertex.Column), x.Vertex.Row+1)
	case Move_RESIGN:
		return fmt.Sprintf("%s %s", x.Color, Move_RESIGN)
	default:
		return fmt.Sprintf("%s %s", x.Color, Move_PASS)
	}
}

// [0-18] -> [A-T] \ I
// 0 -> A
// ...
// 7 -> H
// 8 -> J
// ...
// 18 -> T
func columnToLetter(c int32) string {
	if c < 0 || 18 < c {
		return ""
	} else if c < 8 {
		return string('A' + c)
	} else {
		return string('A' + c + 1)
	}
}

// [A-T] \ I -> [0-18]
// A -> 0
// ...
// H -> 7
// J -> 8
// ...
// T -> 18
func letterToColumn(l rune) int32 {
	if l == 'I' || l < 'A' || 'T' < l {
		return -1
	} else if l < 'I' {
		return l - 'A'
	} else {
		return l - 'A' - 1
	}
}
