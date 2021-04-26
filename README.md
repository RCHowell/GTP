## About

Somewhat of a follow-up to [GTP Protobuf Definition - 11](https://github.com/RCHowell/rchowell.github.io/issues/11) 
because I posted the full 
protobuf definition, but didn't want to publish the serializing and deserializing extensions. The extensions are 
written in the same Go package as the generated code so that I can "extend" the generated structs with receiver 
functions. I've included a small Makefile.

```
make gen
make clean
```

I had been working on a project/exercise, that I'll open source some of, which involved communicating with Go (Baduk)
engines. Go (Baduk) engines communicate with [Go Text Protocol](http://www.lysator.liu.se/~gunnar/gtp/) which is 
sent as plaintext to the engine's stdin or some network interface. I did not want to use
plaintext in my 
program, so I chose to use protobuf for GTP and Go/Baduk structures and wrap Go engines with a gRPC service.

My project's Go (Baduk) central server operates with the generated Go (lang) code and communicates with the 
AIs/Engines e.g.
GnuGo and KataGo with the generated gRPC stub because the engines are wrapped by the generated server. This allows 
for communicating between the central server and various Go (Baduk) AIs/Engines using structured input and output as
well as gives me the benefits of gRPC and the generated stub/client.

## GTP SerDe

### Parsing / Deserialization
The gRPC wrapper server receives output from the engine as plaintext and parses these a response objects to return 
to the caller.
```golang
// Parsing a generic response
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

// Parsing a particular and more complicated response
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
```

### Encoding / Serialization
The gRPC wrapper server receives the callers request objects and calls "ToCommand" to get the plaintext GTP command to 
send to the Engine's input.
```golang
// Example for `Play`
func (x *PlayRequest) ToCommand() string {
    if x.Id == 0 {
        return fmt.Sprintf("%s %s%s", Play, x.Move.GTPString(), LF)
    }
    return fmt.Sprintf("%d %s %s%s", x.Id, Play, x.Move.GTPString(), LF)
}

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
```

## GTP gRPC Server
```protobuf
// Example Request Response Objects
message SetKomiRequest {
  int32 id = 1;
  float komi = 2;
}

message SetKomiResponse {
  int32 id = 1;
  Error error = 2;
}

service GTP {

  // Version of the GTP Protocol
  rpc ProtocolVersion (ProtocolVersionRequest) returns (ProtocolVersionResponse) {}

  // E.g. “GNU Go”, “GoLois”, “Many Faces of Go”. The name does not include any version information. Use `version`.
  rpc Name (NameRequest) returns (NameResponse) {}

  // E.g. “3.1.33”, “10.5”. Engines without a sense of version number return the empty string.
  rpc Version (VersionRequest) returns (VersionResponse) {}

  // Returns “true” if the command is known by the engine, “false” otherwise
  rpc KnownCommand (KnownCommandRequest) returns (KnownCommandResponse) {}

  // Lists all known commands, including required ones and private extensions.
  rpc ListCommands (ListCommandsRequest) returns (ListCommandsResponse) {}

  // The session is terminated and the connection is closed.
  rpc Quit (QuitRequest) returns (QuitResponse) {}

  // Changes the board size. If the engine cannot handle the new size, fails with the error message ”unacceptable size”.
  rpc BoardSize (BoardSizeRequest) returns (BoardSizeResponse) {}

  // Clears the board, captured stones are reset, and move history is reset
  rpc ClearBoard (ClearBoardRequest) returns (ClearBoardResponse) {}

  // Changes the Komi
  rpc Komi (KomiRequest) returns (KomiResponse) {}

  // Plays the given move
  rpc Play (PlayRequest) returns (PlayResponse) {}

  // Asks the engine to generate a move, it will play it, and will return what was played
  rpc GenMove (GenMoveRequest) returns (GenMoveResponse) {}

  // The board and captured stones are reset to the previous move
  rpc Undo(UndoRequest) returns(UndoResponse) {}

}
```

