// R. C. Howell

// Parser to convert GTP responses to the respective Protobuf Objects
//
// This file is NOT generated and CAN be edited by hand.
// Parse commands return an error if parsing failed, not if the command failed

// Replace this with a generated parser from a GTP grammar -- which I don't think exists

package gtp

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Response struct {
	Id       int32
	Error    error
	Elements []string
}

// Some tokens which could be eventually used in a grammar
const (
	LF    = string(10) // \n Line Feed
	CR    = string(13) // Carriage Return
	HT    = string(9)  // Horizontal Tab
	SPACE = " "
	EQ    = '='
	ERR   = '?'
)

// matches all ASCII control characters 0-32 other than HT and LF
var controlChars = regexp.MustCompile("[\\x00-\\x08\\x0b-\\x1f\\x7f]")

func parseResponse(i, command string, minElements int) (*Response, error) {
	clean := strings.ReplaceAll(controlChars.ReplaceAllString(i, ""), HT, SPACE)
	// tokens ~technically~ aren't split on a SPACE
	tokens := strings.Split(clean, SPACE)
	l := len(tokens)
	id, err := parseCommandId(tokens[0])
	if err != nil {
		return nil, err
	}
	// handle empty responses
	if l < 2 {
		return &Response{Id: int32(id)}, nil
	}
	if l-1 < minElements {
		return nil, fmt.Errorf("%s response requires %d response elements", command, minElements)
	}
	succeeded, err := isSuccessful(clean)
	if err != nil {
		return nil, err
	}
	if !succeeded {
		return &Response{
			Id:    int32(id),
			Error: errors.New(strings.Join(tokens[1:l], SPACE)),
		}, nil
	}
	return &Response{
		Id:       int32(id),
		Elements: tokens[1:l],
	}, nil
}

//	=id response\n\n
//	=id\n\n
//	= response\n\n
//	=\n\n
//	?id error_message\n\n
//	? error_message\n\n
func isSuccessful(i string) (bool, error) {
	if i == "" {
		return false, errors.New("cannot parse command success status from empty string")
	}
	switch i[0] {
	case EQ:
		return true, nil
	case ERR:
		return false, nil
	}
	return false, fmt.Errorf("unknown command success status %s", string(i[0]))
}

//   =
//   =N
//   ?
//   ?N
func parseCommandId(i string) (int, error) {
	n := len(i)
	switch n {
	case 0:
		return 0, errors.New("cannot parse command id from empty string")
	case 1:
		return 0, nil
	default:
		return strconv.Atoi(i[1:n])
	}
}

// TODO good candidates for reflect generation

//	= N
//	=ID N
func ParseProtocolVersionResponse(i string) (*ProtocolVersionResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 1)
	if err != nil {
		return nil, err
	}
	v, err := strconv.Atoi(res.Elements[1])
	if err != nil {
		return nil, err
	}
	return &ProtocolVersionResponse{
		Id:      res.Id,
		Version: int32(v),
	}, nil
}

//	= string*
//	=ID string*
func ParseNameResponse(i string) (*NameResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 1)
	if err != nil {
		return nil, err
	}
	return &NameResponse{
		Id:   res.Id,
		Name: strings.Join(res.Elements, SPACE),
	}, nil
}

//	= string*
//	=ID string*
func ParseVersionResponse(i string) (*VersionResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 1)
	if err != nil {
		return nil, err
	}
	return &VersionResponse{
		Id:      res.Id,
		Version: strings.Join(res.Elements, SPACE),
	}, nil
}

func ParseIsKnownCommandResponse(i string) (*IsKnownCommandResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 1)
	if err != nil {
		return nil, err
	}
	k, err := strconv.ParseBool(res.Elements[1])
	if err != nil {
		return nil, err
	}
	return &IsKnownCommandResponse{
		Id:    res.Id,
		Known: k,
	}, nil
}

func ParseListCommandsResponse(i string) (*ListCommandsResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 0)
	if err != nil {
		return nil, err
	}
	return &ListCommandsResponse{
		Id:       res.Id,
		Commands: res.Elements,
	}, nil
}

func ParseQuitResponse(i string) (*QuitResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 0)
	if err != nil {
		return nil, err
	}
	return &QuitResponse{Id: res.Id}, nil
}

func ParseSetBoardSizeResponse(i string) (*SetBoardSizeResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 0)
	if err != nil {
		return nil, err
	}
	o := &SetBoardSizeResponse{Id: res.Id}
	if res.Error != nil {
		o.Error = &Error{Message: res.Error.Error()}
	}
	return o, nil
}

func ParseClearBoardResponse(i string) (*ClearBoardResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 0)
	if err != nil {
		return nil, err
	}
	return &ClearBoardResponse{Id: res.Id}, nil
}

func ParseSetKomiResponse(i string) (*SetKomiResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 0)
	if err != nil {
		return nil, err
	}
	o := &SetKomiResponse{Id: res.Id}
	if res.Error != nil {
		o.Error = &Error{Message: res.Error.Error()}
	}
	return o, nil
}

func ParsePlayResponse(i string) (*PlayResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 0)
	if err != nil {
		return nil, err
	}
	o := &PlayResponse{Id: res.Id}
	if res.Error != nil {
		o.Error = &Error{Message: res.Error.Error()}
	}
	return o, nil
}

func ParseGenMoveResponse(i string) (*GenMoveResponse, error) {
	res, err := parseResponse(i, ProtocolVersion, 1)
	if err != nil {
		return nil, err
	}
	t := strings.ToUpper(res.Elements[0])
	switch t {
	case Move_PASS.String():
		return &GenMoveResponse{Id: res.Id, Move: &Move{Type: Move_PASS}}, nil
	case Move_RESIGN.String():
		return &GenMoveResponse{Id: res.Id, Move: &Move{Type: Move_RESIGN}}, nil
	}
	v, err := ParseVertex(t)
	if err != nil {
		return nil, err
	}
	return &GenMoveResponse{
		Id: res.Id,
		Move: &Move{
			Type:   Move_PLACE,
			Vertex: v,
		},
	}, nil
}

func ParseUndoResponse(i string) (*UndoResponse, error) {
	res, err := parseResponse(i, Undo, 0)
	if err != nil {
		return nil, err
	}
	o := &UndoResponse{Id: res.Id}
	if res.Error != nil {
		o.Error = &Error{Message: res.Error.Error()}
	}
	return &UndoResponse{Id: res.Id}, nil
}

func ParseColor(i string) Color {
	switch strings.TrimSpace(strings.ToUpper(i)) {
	case "W", Color_WHITE.String():
		return Color_WHITE
	case "B", Color_BLACK.String():
		return Color_BLACK
	}
	return Color_EMPTY
}

// TODO include range checks
//	StringInt
// 	GTP uses letters for columns and 1-index integers for rows. This method returns 0-index row, columns
//	Ex: A12 -> 11,0
//		a2 -> 1,0
//		j7 -> 6,9
func ParseVertex(i string) (*Vertex, error) {
	l := len(i)
	if l < 2 {
		return nil, errors.New("vertex input has fewer than 2 characters")
	}
	r, err := strconv.Atoi(i[1:l])
	if err != nil {
		return nil, err
	}
	return &Vertex{
		Row:    int32(r - 1),
		Column: letterToColumn(rune(strings.ToUpper(i)[0])),
	}, nil
}

// Split function to split on two line feeds
func SplitOnDoubleLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := 1; i < len(data); i++ {
		if data[i-1] == '\n' && data[i] == '\n' {
			return i + 1, data[0 : i-1], nil
		}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
