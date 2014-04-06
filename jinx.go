package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "errors"
    //"strconv"
    //"reflect"
)

type LAReader struct {
    R *bufio.Reader
    la_mode bool
    la_offset int
}

func (b *LAReader) Read(l int) ([]byte, error) {
    if b.la_mode {
        // lookahead mode doesn't move the fd position
        return b.FakeRead(l)
    } else {
         buf := make([]byte, l)
         n, err := b.R.Read(buf)
         if err != nil {
            return nil, err
         } else if n != l {
            return nil, errors.New("Not enough data available to Read")
         } else {
            // fmt.Printf(">>>>> READ %c\n", buf[0])
            // fmt.Println(buf)
            return buf, nil
         }
    }
}

func (b *LAReader) FakeRead(l int) ([]byte, error) {
    var bs []byte
    var err error

    // lookahead mode, add the number of bytes we've
    // already "read" (LA'd) to the size of the peek
    // This is to handle multiple calls to Peek() during
    // la_mode == true
    bs, err = b.R.Peek(l+b.la_offset)
    bs = bs[b.la_offset:]
    // the line below is the only diff
    b.la_offset += l

    if err != nil {
        return nil, err
    } else if len(bs) != l {
        return nil, errors.New("Not enough data available to Peek")
    } else {
        //fmt.Printf(">>>>> FAKE %c\n", bs[0])
        return bs, nil
    }
}

func (b *LAReader) Peek(l int) ([]byte, error) {
    var bs []byte
    var err error
    if b.la_mode {
        // lookahead mode, add the number of bytes we've
        // already "read" (LA'd) to the size of the peek
        // This is to handle multiple calls to Peek() during
        // la_mode == true
        bs, err = b.R.Peek(l+b.la_offset)
        bs = bs[b.la_offset:]
        //b.la_offset += l
    } else {
        bs, err = b.R.Peek(l)
    }
    if err != nil {
        return nil, err
    } else if len(bs) != l {
        return nil, errors.New("Not enough data available to Peek")
    } else {
        //fmt.Printf(">>>>> PEEK %c\n", bs[0])
        return bs, nil
    }
}

func (b *LAReader) StartLA() {
    b.la_mode = true
    b.la_offset = 0
}

func (b *LAReader) StopLA() {
    b.la_mode = false
    b.la_offset = 0
}

func NewLAReader(br *bufio.Reader) *LAReader {
    return &LAReader{br,false,0}
}


type ResultGen func(s ...interface{}) interface{}


type ParseResult struct {
    Result  interface{}
    Success bool
    Position int   // inclusive
    Length   int   // index not included
}

type ParserState struct {
    LAR        *LAReader
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
    r := strings.NewReader(s)
    ps.LAR = NewLAReader(bufio.NewReader(r))
}

func (ps *ParserState) FromFile(fn string) {
    fi, _ := os.Open(fn)
    // TODO: defer?
    ps.fi = fi
    r := NewLAReader(bufio.NewReader(fi))
    ps.LAR = r
}


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
            bytes, err := ps.LAR.Peek(1)
            if len(bytes) != 1 || err != nil {
                // TODO: parse error
                return &ParseResult{nil, false, ps.Position, 0}
            } else if bytes[0] == cdata {
                ps.LAR.Read(1)
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
            bytes, err := ps.LAR.Peek(1)
            if len(bytes) != 1 || err != nil {
                // TODO: parse error
                return &ParseResult{nil, false, ps.Position, 0}
            }

            for _,c := range sdata {
                if bytes[0] == c {
                    ps.LAR.Read(1)
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
                    ps.LAR.StartLA()
                    pr := subparser.Parse(ps)
                    //fmt.Printf("%#v\n", pr)
                    ps.LAR.StopLA()
                    ps.Position = preLAPosition
                    if pr.Success {
                        pr = subparser.Parse(ps)
                        return &ParseResult{p.Gen(pr.Result), true, ps.Position, pr.Length}
                    } else {
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
            bytes, err := ps.LAR.Peek(expectedLen)
            if len(bytes) != expectedLen || err != nil {
                // TODO: parse error
                return &ParseResult{nil, false, ps.Position, 0}
            } else if string(bytes) == sdata {
                // Use read for it's side effect on the buffer, ignore the result
                ps.LAR.Read(expectedLen)
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

//sepEndBy parses a sequence of p separated and optionally ended by sep.

func main() {
    fmt.Println("Hello world")
    // fmt.Println(s.Parse(ps))
    // d := MSeq(
    //     func(s ...interface{}) interface{} {
    //             r := s[0].(string) + s[1].(string) + s[2].(string) + s[3].(string)
    //             i, _ := strconv.Atoi(r)
    //             return IntNode{i}
    //     }, f, o, o, o)
    // fmt.Printf("%#v\n",d.Parse(ps))

}

