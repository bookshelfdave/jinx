package main

import (
    "os"
    "fmt"
    //"errors"
    //"strconv"
    //"reflect"
)

type ResultGen func(s ...interface{}) interface{}


// type ParsePosition struct {
//     // TODO: change int's to uint64 etc
//     Index int
//     Line  int
//     Col   int
// }

type ParseResult struct {
    Result  interface{}
    Success bool
    Position int   // inclusive
    Length   int   // index not included
}

type ParserState struct {
    R          *JinxReader
    //R        *bufio.Reader
    Position int
    Line     int
    fi       *os.File
}

type ParseFn func(p *Parser, ps *ParserState) *ParseResult

type Parser struct {
    // holds internal parser data
    // ie: in Char('x'), data is 'x'
    data interface{}
    parseFn ParseFn
    Gen ResultGen
}

func (p *Parser) Parse(ps *ParserState) *ParseResult {
    return p.parseFn(p, ps)
}


// TODO: squash ConcatParams + ConcatArray into the same functions
//       just use an array instead of a param list
func ConcatParams(s ...interface{}) interface{} {
    var a string
    for i := range s {
        if v,ok := s[i].(string); ok {
            a += v
        } else {
            fmt.Println("Invalid result")
        }
    }
    return a
}

func IgnoreParams(s ...interface{}) interface{} {
    return ""
}


func ConcatArray(arr ...interface{}) interface{} {
    var a string// probably inefficient
    ss := (arr[0]).([]interface{})
    for i,_ := range ss {
        if v, ok := ss[i].(string); ok {
            a += v
        } else {
            fmt.Println("Invalid type")
        }
    }
    return a
}


func decStringResult(s ...interface{}) interface{} {
    var a string
    for i := range s {
        a += "<<"
        a += s[i].(string)
        a += ">>"
    }
    return a
}

func (ps *ParserState) ParserFromString(s string) {
    ps.R = NewJinxReaderFromString(s)
}

// func (ps *ParserState) FromFile(fn string) {
//     fi, _ := os.Open(fn)
//     // TODO: defer?
//     ps.fi = fi
//     r := NewLAReader(bufio.NewReader(fi))
//     ps.R = r
// }


// func Fail() *Parser {
//     parse := func(p Parser, ps *ParserState) *ParseResult {
//         return &ParseResult{nil, false, ps.Position, 0}
//     }
//     return &Parser{c, parse, g}
// }

func Char(c byte) *Parser {
    return CharWithGen(ConcatParams, c)
}

func CharWithGen(g ResultGen, c byte) *Parser {
    parse := func(p *Parser, ps *ParserState) *ParseResult {
            cdata := p.data.(byte)
            //fmt.Printf("Char(%c) parsing\n", cdata)
            bytes, err := ps.R.Peek(1)
            if len(bytes) != 1 || err != nil {
                // TODO: parse error
                return &ParseResult{nil, false, ps.Position, 0}
            } else if bytes[0] == cdata {
                ps.R.Read(1)
                //fmt.Printf("Char: %c\n", bytes[0])
                pr := &ParseResult{p.Gen(string(bytes)), true, ps.Position, 1}
                ps.Position++
                if bytes[0] == '\n' {
                    ps.Line++
                }
                return pr
            }
            return &ParseResult{nil, false, ps.Position, 0}
    }
    //fmt.Printf("Making a Char parser with %c\n", c)
    return &Parser{c, parse, g}
}

func CharFrom(s string) *Parser {
    return CharFromWithGen(ConcatParams, s)
}

// not rune safe
func CharFromWithGen(g ResultGen, s string) *Parser {
     parse := func(p *Parser, ps *ParserState) *ParseResult {
            sdata := p.data.([]byte)
            bytes, err := ps.R.Peek(1)
            if len(bytes) != 1 || err != nil {
                // TODO: parse error
                return &ParseResult{nil, false, ps.Position, 0}
            }

            for _,c := range sdata {
                if bytes[0] == c {
                    ps.R.Read(1)
                    pr := &ParseResult{p.Gen(string(bytes)), true, ps.Position, 1}
                    ps.Position++
                    if bytes[0] == '\n' {
                        ps.Line++
                    }
                    return pr
                }
            }
            return &ParseResult{nil, false, ps.Position, 0}
    }
    byteArray := []byte(s)
    return &Parser{byteArray, parse, g}
}

func Lower() *Parser {
    return CharFrom("abcdefghijklmnopqrstuvwxyz")
}

func LowerWithGen(g ResultGen) *Parser {
    return CharFromWithGen(g, "abcdefghijklmnopqrstuvwxyz")
}

func Upper() *Parser {
    return CharFrom("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func UpperWithGen(g ResultGen) *Parser {
    return CharFromWithGen(g, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func Letter() *Parser {
    return Alt(Upper(), Lower())
}

func LetterWithGen(g ResultGen) *Parser {
    return Alt(UpperWithGen(g), LowerWithGen(g))
}

func Digit() *Parser {
    return CharFrom("0123456789")
}

func DigitWithGen(g ResultGen) *Parser {
    return CharFromWithGen(g, "0123456789")
}

func Alphanum() *Parser {
    return Alt(Letter(), Digit())
}

func AlphanumWithGen(g ResultGen) *Parser {
    return Alt(LetterWithGen(g), DigitWithGen(g))
}

func Word() *Parser {
    return Many1(Letter())
}

func WordWithGen(g ResultGen) *Parser {
    return Many1WithGen(g, Letter())
}

func Number() *Parser {
    return Many1(Digit())
}

func NumberWithGen(g ResultGen) *Parser {
    return Many1WithGen(g, Digit())
}

func WS() *Parser {
    return CharFrom("\n\t\r")
}

func IgnoreWS() *Parser {
    return Ignore(CharFrom("\n\t\r"))
}

func WSWithGen(g ResultGen) *Parser {
    return CharFromWithGen(g, "\n\t\r")
}

func Many(subparser *Parser) *Parser {
    return ManyWithGen(ConcatArray, subparser)
}

func ManyWithGen(g ResultGen, subparser *Parser) *Parser {
    parse := func(p *Parser, ps *ParserState) *ParseResult {
            //results := make([]*ParseResult,1)
            var totalLen int
            results := make([]interface{},0)
            subparser := p.data.(*Parser)
            for {
                pr := subparser.Parse(ps)
                if pr.Success {
                    totalLen += pr.Length
                    results = append(results, pr.Result)
                } else {
                    break
                }
            }
            return &ParseResult{p.Gen(results), true, ps.Position, totalLen}
    }
    return &Parser{subparser, parse, g}
}

func Many1(subparser *Parser) *Parser {
    return Many1WithGen(ConcatArray, subparser)
}

func Many1WithGen(g ResultGen, subparser *Parser) *Parser {
    parse := func(p *Parser, ps *ParserState) *ParseResult {
            //results := make([]*ParseResult,1)
            var totalLen int
            results := make([]interface{},0)
            subparser := p.data.(*Parser)
            // at least 1
            finalPosition := ps.Position
            pr := subparser.Parse(ps)
            if pr.Success {
                totalLen += pr.Length
                results = append(results, pr.Result)
            } else {
                return &ParseResult{nil, false, ps.Position, 0}
            }
            for {
                pr = subparser.Parse(ps)
                if pr.Success {
                    totalLen += pr.Length
                    results = append(results, pr.Result)
                } else {
                    break
                }
            }
            return &ParseResult{p.Gen(results), true, finalPosition, totalLen}
    }
    return &Parser{subparser, parse, g}
}


func Attempt(subparser *Parser) *Parser {
    return AttemptWithGen(ConcatParams, subparser)
}

func AttemptWithGen(g ResultGen, subparser *Parser) *Parser {
    parse := func(p *Parser, ps *ParserState) *ParseResult {
                    subparser := p.data.(*Parser)
                    preLAPosition := ps.Position
                    rewindTo := ps.R.Offset
                    pr := subparser.Parse(ps)
                    ps.Position = preLAPosition
                    if pr.Success {
                        pr = subparser.Parse(ps)
                        return &ParseResult{p.Gen(pr.Result), true, ps.Position, pr.Length}
                    } else {
                        ps.R.Seek(rewindTo)
                        return &ParseResult{nil, false, ps.Position, 0}
                    }
    }
    return &Parser{subparser, parse, g}
}

func Str(s string) *Parser {
    return StrWithGen(ConcatParams, s)
}

func StrWithGen(g ResultGen, s string) *Parser {
    parse := func(p *Parser, ps *ParserState) *ParseResult {
            sdata := p.data.(string)
            expectedLen := len(s)
            bytes, err := ps.R.Peek(expectedLen)
            if len(bytes) != expectedLen || err != nil {
                // TODO: parse error
                return &ParseResult{nil, false, ps.Position, 0}
            } else if string(bytes) == sdata {
                // Use read for it's side effect on the buffer, ignore the result
                ps.R.Read(expectedLen)
                pr := &ParseResult{p.Gen(string(bytes)), true, ps.Position, expectedLen}
                ps.Position += expectedLen
                // todo: look for newlines
                // if bytes[0] == '\n' {
                //     ps.Line++
                // }
                return pr
            }
            return &ParseResult{nil, false, ps.Position, 0}
    }

    return &Parser{s, parse, g}
}

func seqParser(p *Parser, ps *ParserState) *ParseResult {
    allps := p.data.([]*Parser)
    results := make([]*ParseResult, len(allps))
    raw_results := make([]interface{}, len(allps))
    for i := range allps {
        results[i] = allps[i].Parse(ps)
        //fmt.Printf("Seq: %#v\n", results[i])
        raw_results[i] = results[i].Result
        if !results[i].Success {
            return &ParseResult{nil, false, ps.Position, 0}
        }
    }

    var totalLength int
    for _, i := range results {
        totalLength += i.Length
    }

    return &ParseResult{p.Gen(raw_results...),
                                true,
                                results[0].Position,
                                totalLength }
}

func Seq(parsers ...*Parser) *Parser {
    return &Parser{parsers, seqParser, ConcatParams}
}

func SeqWithGen(g ResultGen, parsers ...*Parser) *Parser {
    return &Parser{parsers, seqParser, g}
}


func altParser(p *Parser, ps *ParserState) *ParseResult {
    allps := p.data.([]*Parser)
    var one_result *ParseResult

    for i := range allps {
        result := (*allps[i]).Parse(ps)
        if result.Success == true {
           one_result = result
           break;
        }
    }

    if one_result == nil {
        return &ParseResult{nil, false, ps.Position, 0}
    }

    return &ParseResult{p.Gen(one_result.Result), true, one_result.Position, one_result.Length}
}

func Alt(parsers ...*Parser) *Parser {
    return &Parser{parsers, altParser, ConcatParams}
}

func AltWithGen(g ResultGen, parsers ...*Parser) *Parser {
    return &Parser{parsers, altParser, g}
}


func proxyParser(p *Parser, ps *ParserState) *ParseResult {
    if p.data == nil {
        return &ParseResult{"Proxy object doesn't have a parser", false, ps.Position, 0}
    }
    subparser := p.data.(*Parser)
    return subparser.Parse(ps)
}

func Proxy() *Parser {
    return &Parser{nil, proxyParser, ConcatParams}
}

func ProxyWithGen(g ResultGen) *Parser {
    return &Parser{nil, proxyParser, g}
}

func ProxySetParser(proxy *Parser, p *Parser) {
    proxy.data = p
}


type betweenData struct {
    first  *Parser
    last   *Parser
    p      *Parser
}

func betweenParser(p *Parser, ps *ParserState) *ParseResult {
    bd := p.data.(*betweenData)
    firstResult := bd.first.Parse(ps) // toss the result if valid
    if firstResult.Success {
        pr := bd.p.Parse(ps)
        lastResult := bd.last.Parse(ps) // toss the result if valid
        if lastResult.Success {
            // note: first length and last length
            return &ParseResult{p.Gen(pr.Result), true, pr.Position, pr.Length}
        } else {
            return &ParseResult{nil, false, ps.Position, 0}
        }
    } else {
        return &ParseResult{nil, false, ps.Position, 0}
    }
}

func Between(first *Parser, p *Parser, last *Parser) *Parser {
    return BetweenWithGen(ConcatParams, first, p, last)
}

func BetweenWithGen(g ResultGen, first *Parser, p *Parser, last *Parser) *Parser {
    d := &betweenData{first, last, p}
    return &Parser{d, betweenParser, g}
}


type sepByData struct {
    p   *Parser
    sep *Parser
}

// The parser sepBy p sep parses zero or more occurrences of p separated
// by sep (in EBNF: (p (sep p)*)?). It returns a list of the results returned by p.

func sepByParser(p *Parser, ps *ParserState) *ParseResult {
    d := p.data.(*sepByData)
    prs := make([]*ParseResult, 0)
    raw_results := make([]interface{}, 0)
    finalPosition := ps.Position
    result := d.p.Parse(ps)
    if result.Success {
        prs = append(prs, result)
        raw_results = append(raw_results, result.Result)
        for {
            if d.sep.Parse(ps).Success {
                result = d.p.Parse(ps)
                if result.Success {
                    prs = append(prs, result)
                    raw_results = append(raw_results, result.Result)
                } else {
                    return &ParseResult{nil, false, finalPosition, 0}
                }
            } else {
                break
            }
        }

        var totalLength int
        for _, i := range prs {
            totalLength += i.Length
        }

        return &ParseResult{p.Gen(raw_results), true, finalPosition, totalLength}
    } else {
        // always succeeds
        return &ParseResult{"", true, ps.Position, 0}
    }
}

func SepBy(p *Parser, sep *Parser) *Parser {
    return SepByWithGen(ConcatArray, p, sep)
}

//sepBy p sep parses a sequence of p separated by sep and returns the results in a list. 
func SepByWithGen(g ResultGen, p *Parser, sep *Parser) *Parser{
    d := &sepByData{p, sep}
    return &Parser{d, sepByParser, g}
}


func ignoreParser(p *Parser, ps *ParserState) *ParseResult {
    if p.data == nil {
        return &ParseResult{"Ignore object doesn't have a parser", false, ps.Position, 0}
    }
    finalPosition := ps.Position
    subparser := p.data.(*Parser)
    result := subparser.Parse(ps)
    if result.Success {
        return &ParseResult{"", true, finalPosition, result.Length}
    } else {
        return &ParseResult{nil, false, ps.Position, 0}
    }
}

func Ignore(subparser *Parser) *Parser {
    return &Parser{subparser, ignoreParser, nil}
}


//sepEndBy parses a sequence of p separated and optionally ended by sep.



func main() {
    //r := strings.NewReader("Hello world")
    // s := "Dave Test"
    // jr := NewJinxReaderFromString(s)
    // fmt.Printf("%c\n", jr.Read(1)[0])
    // fmt.Printf("%c\n", jr.Read(1)[0])
    // fmt.Printf("%c\n", jr.Peek(1)[0])
    // fmt.Printf("%c\n", jr.Read(1)[0])
    // jr.Seek(1)
    // fmt.Printf("%c\n", jr.Read(1)[0])
    // jr.Seek(0)
    // fmt.Printf("%c\n", jr.Read(1)[0])
    // buf := make([]byte, 1)
    // jr.rs.Read(buf)
    // fmt.Printf("%c\n", buf[0])

    // jr.rs.Read(buf)
    // fmt.Printf("%c\n", buf[0])

    // jr.rs.Seek(0,0)
    // jr.rs.Read(buf)
    // fmt.Printf("%c\n", buf[0])

    // rs0 := io.ReadSeeker(r)
    // x := bufio.NewReader(rs0)
    // b, _ := x.ReadByte()
    // fmt.Printf("%c\n", b)

    // rs1 := io.ReadSeeker(r)
    // rs1.Seek(0,0)
    // x1 := bufio.NewReader(rs1)
    // b1, _ := x1.ReadByte()
    // fmt.Printf("%c\n", b1)
    // fmt.Println(s.Parse(ps))
    // d := MSeq(
    //     func(s ...interface{}) interface{} {
    //             r := s[0].(string) + s[1].(string) + s[2].(string) + s[3].(string)
    //             i, _ := strconv.Atoi(r)
    //             return IntNode{i}
    //     }, f, o, o, o)
    // fmt.Printf("%#v\n",d.Parse(ps))

}

