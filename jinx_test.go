package jinx
import (
    "testing"
    "strconv"
    "fmt"
   // "reflect"
    )

func assertRead(jr *JinxReader, s string, t *testing.T) bool {
    n := len(s)
    bs, _ := jr.Read(n)
    if string(bs) == s {
        return true
    } else {
        t.Error("Fail!")
        return false
    }
}

func assertPeek(jr *JinxReader, s string, t *testing.T) bool {
    n := len(s)
    bs, _ := jr.Peek(n)
    if string(bs) == s {
        return true
    } else {
        t.Error("Fail!")
        return false
    }
}

func TestJinxReader(t *testing.T) {
    s := "123456789"
    jr := NewJinxReaderFromString(s)
    assertRead(jr, "1", t)
    assertRead(jr, "2", t)
    assertPeek(jr, "345", t)
    assertRead(jr, "345", t)
    jr.Seek(0)
    assertRead(jr, "123", t)
    jr.Seek(0)
    assertRead(jr, "123", t)
}

// Most of these tests use the ConcatParams result generator, which
// returns a concatenated string of parse results

func TestSimpleChar(t *testing.T) {
    fmt.Println("Testing Jinx...")
    {
        ps := new(ParserState)
        ps.ParserFromString("123")

        p0 := Char('1')
        result := p0.Parse(ps)
        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "1") {
            t.Error("Expected 1")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 1) {
         t.Error("Expected Length == 1")
        }
    }

    {
        ps := new(ParserState)
        ps.ParserFromString("123")
        // instead of returning the string "1", an int is returned
        p0 := Char('1').WithGen(func (s interface{}) interface{} {
                str := s.(string)
                r,_ := strconv.Atoi(str)
                return r
            })
        result := p0.Parse(ps)
        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != 1) {
            t.Error("Expected 1")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 1) {
         t.Error("Expected Length == 1")
        }
    }

}

func TestMultipleChar(t *testing.T) {
    // Parser state should advance when parsing multiple characters
    ps := new(ParserState)
    ps.ParserFromString("123")

    p0 := Char('1')
    p1 := Char('2')
    result0 := p0.Parse(ps)
    result1 := p1.Parse(ps)

    if(result0.Success == false) {
        t.Error("Expected parse success")
    }
    if(result0.Result != "1") {
        t.Error("Expected 1")
    }
    if(result0.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result0.Length != 1) {
     t.Error("Expected Length == 1")
    }


    if(result1.Success == false) {
        t.Error("Expected parse success")
    }
    if(result1.Result != "2") {
        t.Error("Expected 2")
    }
    if(result1.Position != 1) {
        t.Error("Expected Position == 1")
    }
    if(result1.Length != 1) {
     t.Error("Expected Length == 1")
    }

}

func TestSimpleStr(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("foobarbaz")
    foo := Str("foo")
    bar := Str("bar")

    result := foo.Parse(ps)
    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "foo") {
        t.Error("Expected foo")
    }
    if(result.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result.Length != 3) {
        t.Error("Expected Length == 3")
    }

    // test that a string advances the parser state
    result2 := bar.Parse(ps)
    if(result2.Success == false) {
        t.Error("Expected parse success")
    }
    if(result2.Result != "bar") {
        t.Error("Expected bar")
    }
    if(result2.Position != 3) {
        t.Error("Expected Position == 0")
    }
    if(result2.Length != 3) {
        t.Error("Expected Length == 3")
    }
}

func TestSimpleSeq(t *testing.T) {
    {
        ps := new(ParserState)
        ps.ParserFromString("123")

        p0 := Char('1')
        p1 := Char('2')
        p2 := Char('3')
        s  := Seq(p0, p1, p2).WithGen(GenString)
        result := s.Parse(ps)

        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "123") {
            t.Error("Expected 123")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 3) {
            t.Error("Expected Length == 3")
        }
    }

    {
        ps := new(ParserState)
        ps.ParserFromString("123")

        p0 := Char('1')
        p1 := Char('2')
        p2 := Char('3')
        // this generator concats all the chars into a string and then converts
        // into an int
        s  := Seq(p0, p1, p2).PipeGen(GenString,GenStringToInt)
        result := s.Parse(ps)

        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != 123) {
            t.Error("Expected 123")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 3) {
            t.Error("Expected Length == 3")
        }
    }
}


func TestSimpleAlt(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("123")
    a := Char('x')
    b := Char('1')
    alt := Alt(a, b)
    result := alt.Parse(ps)

    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "1") {
        t.Error("Expected 1")
    }
    if(result.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result.Length != 1) {
     t.Error("Expected Length == 1")
    }
}

func TestNestedAlt(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("123")

    alt := Alt(Char('x'), Alt(Char('y'), Char('1')))
    result := alt.Parse(ps)
    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "1") {
        t.Error("Expected 1")
    }
    if(result.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result.Length != 1) {
     t.Error("Expected Length == 1")
    }
}


func TestNestedStrAlt(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("foo123")

    alt := Alt(Str("foo"), Alt(Char('y'), Char('1')))
    result := alt.Parse(ps)

    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "foo") {
        t.Error("Expected foo")
    }
    if(result.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result.Length != 3) {
     t.Error("Expected Length == 3")
    }
}

func TestStrSeq(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("foobarbaz")

    alt := Seq(Str("foo"), Str("bar"), Alt(Str("BAZ"), Str("baz"))).WithGen(GenString)
    result := alt.Parse(ps)

    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "foobarbaz") {
        t.Error("Expected foobarbaz")
    }
    if(result.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result.Length != 9) {
     t.Error("Expected Length == 9")
    }
}


func TestCharFrom(t *testing.T) {
    {
        ps := new(ParserState)
        ps.ParserFromString("zya123")
        oneOf := CharFrom("abcdefghijklmnopqrstuvwxyz")
        result0 := oneOf.Parse(ps)
        result1 := oneOf.Parse(ps)
        if(result0.Success == false) {
            t.Error("Expected parse success")
        }
        if(result0.Result != "z") {
            t.Error("Expected z")
        }
        if(result0.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result0.Length != 1) {
         t.Error("Expected Length == 1")
        }

        if(result1.Success == false) {
            t.Error("Expected parse success")
        }
        if(result1.Result != "y") {
            t.Error("Expected z")
        }
        if(result1.Position != 1) {
            t.Error("Expected Position == 1")
        }
        if(result1.Length != 1) {
         t.Error("Expected Length == 1")
        }
    }
}

// func TestCharFrom(t *testing.T) {
//     ps := new(ParserState)
//     ps.ParserFromString("abcdefghijklmnopqrstuvwxyz")
//     result := oneOf.Parse(ps)

// }

func TestMany(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("12345abcdef")
    m := Many(Digit())
    result := m.Parse(ps)

    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "12345") {
        t.Error("Expected 12345")
    }
    if(result.Position != 5) {
        t.Error("Expected Position == 5")
    }
    if(result.Length != 5) {
     t.Error("Expected Length == 5")
    }
}

func TestMany1(t *testing.T) {
    {
        ps := new(ParserState)
        ps.ParserFromString("12345abcdef")
        m := Many1(Digit())
        result := m.Parse(ps)

        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "12345") {
            t.Error("Expected 12345")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 5) {
         t.Error("Expected Length == 5")
        }
    }

    {
        ps := new(ParserState)
        ps.ParserFromString("abcd")
        m := Many1(Digit())
        result := m.Parse(ps)

        if(result.Success == true) {
            t.Error("Expected parse failure")
        }
        if(result.Result != nil) {
            t.Error("Expected nil")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 0) {
         t.Error("Expected Length == 0")
        }
    }

    {
        ps := new(ParserState)
        ps.ParserFromString("1a")
        m := Many1(Digit())
        result := m.Parse(ps)

        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "1") {
            t.Error("Expected 1")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 1) {
         t.Error("Expected Length == 1")
        }
    }
}


func TestManyAlt(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("12345abcde")
    m := Many(Alt(Digit(), Letter()))
    result := m.Parse(ps)

    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "12345abcde") {
        t.Error("Expected 12345abcde")
    }
    if(result.Position != 10) {
        t.Error("Expected Position == 5")
    }
    if(result.Length != 10) {
     t.Error("Expected Length == 5")
    }
}

// TODO: ATTEMPT IS BROKEN!
func TestAttempt(t *testing.T) {
    {
        ps := new(ParserState)
        ps.ParserFromString("abc")
        fc := Seq(Char('f'), Char('c'))
        ab := Seq(Char('a'), Char('b'))
        m := Alt(Attempt(fc), ab).WithGen(GenString)
        result := m.Parse(ps)

        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "ab") {
            t.Error("Expected ab")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 2) {
         t.Error("Expected Length == 2")
        }
    }

    {
        ps := new(ParserState)
        ps.ParserFromString("a")
        p := Attempt(Char('a'))
        result := p.Parse(ps)

        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "a") {
            t.Error("Expected a")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 1) {
         t.Error("Expected Length == 1")
        }
    }
}

func TestProxy(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("foobar")
    p := Proxy()
    s := Str("foo")
    ProxySetParser(p, s)
    result := p.Parse(ps)

    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "foo") {
        t.Error("Expected foo")
    }
    if(result.Position != 0) {
        t.Error("Expected Position == 0")
    }
    if(result.Length != 3) {
     t.Error("Expected Length == 3")
    }
}

func TestBetween(t *testing.T) {
    ps := new(ParserState)
    ps.ParserFromString("[foo]")
    s := Str("foo")
    b := Between(Char('['), s, Char(']'))
    result := b.Parse(ps)
    if(result.Success == false) {
        t.Error("Expected parse success")
    }
    if(result.Result != "foo") {
        t.Error("Expected foo")
    }
    if(result.Position != 1) {
        t.Error("Expected Position == 1")
    }
    if(result.Length != 3) {
     t.Error("Expected Length == 3")
    }
}

func TestSepBy(t *testing.T) {
    // {
    //     ps := new(ParserState)
    //     ps.ParserFromString(",2,3,4")
    //     sep := Char(',')
    //     p   := Digit()
    //     next := Attempt(Seq(sep, p))
    //     fmt.Println(next.Parse(ps))
    // }

    {
        ps := new(ParserState)
        ps.ParserFromString("1,2,3,4")
        digits := SepBy(Digit(), Char(',')).WithGen(GenString)
        result := digits.Parse(ps)
        fmt.Println(result)
        if(result.Success == false) {
            t.Error("Expected parse success")
        }
        if(result.Result != "1234") {
            t.Error("Expected 1234")
        }
        if(result.Position != 0) {
            t.Error("Expected Position == 0")
        }
        if(result.Length != 4) {
            t.Error("Expected Length == 4")
        }
    }

    // {
    //     // TODO: should this be a failure?
    //     ps := new(ParserState)
    //     ps.ParserFromString("1,2,")
    //     digits := SepBy(Digit(), Char(','))
    //     result := digits.Parse(ps)
    //     if(result.Success == true) {
    //         t.Error("Expected parse failure")
    //     }
    //     if(result.Result != nil) {
    //         t.Error("Expected nil")
    //     }
    //     if(result.Position != 0) {
    //         t.Error("Expected Position == 0")
    //     }
    //     if(result.Length != 0) {
    //      t.Error("Expected Length == 0")
    //     }
    // }

    // {
    //     ps := new(ParserState)
    //     // sep by 0 or more, passes if if it doesn't find any matches
    //     ps.ParserFromString("xyz")
    //     digits := SepBy(Digit(), Char(','))
    //     result := digits.Parse(ps)
    //     if(result.Success == false) {
    //         t.Error("Expected parse success")
    //     }
    //     if(result.Result != "") {
    //         t.Error("Expected empty result")
    //     }
    //     if(result.Position != 0) {
    //         t.Error("Expected Position == 0")
    //     }
    //     if(result.Length != 0) {
    //      t.Error("Expected Length == 4")
    //     }
    // }

    // {
    //     ps := new(ParserState)
    //     ps.ParserFromString("[1,2,3,4]")
    //     digits := Between(Char('['),
    //                 SepBy(Digit(), Char(',')),
    //             Char(']'))
    //     result := digits.Parse(ps)
    //     if(result.Success == false) {
    //         t.Error("Expected parse success")
    //     }
    //     if(result.Result != "1234") {
    //         t.Error("Expected 1234")
    //     }
    //     if(result.Position != 1) {
    //         t.Error("Expected Position == 1")
    //     }
    //     if(result.Length != 4) {
    //      t.Error("Expected Length == 4")
    //     }
    // }

    // {
    //     ps := new(ParserState)
    //     ps.ParserFromString("[1,2,3,4]")
    //     g := func (arr ...interface{}) interface{} {
    //         ss := (arr[0]).([]interface{})
    //         intList := make([]int, len(ss))
    //         for i,_ := range ss {
    //             if v, ok := ss[i].(string); ok {
    //                 intval, _ := strconv.Atoi(v)
    //                 intList[i] = intval
    //             } else {
    //                 fmt.Println("Invalid type")
    //             }
    //         }
    //         return intList
    //     }
    //     digitList := SepByWithGen(g,Digit(), Char(','))
    //     digits := Between(Char('['), digitList, Char(']'))
    //     result := digits.Parse(ps)
    //     if(result.Success == false) {
    //         t.Error("Expected parse success")
    //     }
    //     if(result.Result != "1234") {
    //         t.Error("Expected 1234")
    //     }
    //     if(result.Position != 1) {
    //         t.Error("Expected Position == 1")
    //     }
    //     if(result.Length != 4) {
    //      t.Error("Expected Length == 4")
    //     }
    // }
}


// func TestParseCSV(t *testing.T) {
//     ps := new(ParserState)
//     ps.ParserFromString("foo,bar,baz\n1,2,3\n4,5,6\n7,8,9")
//     // TODO: not checking for WS at the end of the header row
//     //       ... but the parse succeeds???
//     header := Seq(SepBy(Word(), Char(',')))
//     line   := Seq(SepBy(Number(), Char(',')), IgnoreWS())
//     lines  := Many(line)

//     csv    := Seq(header, lines)
//     result := csv.Parse(ps)
//     fmt.Println(result)

//     if(result.Success == false) {
//         t.Error("Expected parse success")
//     }
//     if(result.Result != "foobarbaz123456") {
//         t.Error("Expected foobarbaz123456")
//     }
//     if(result.Position != 0) {
//         t.Error("Expected Position == 0")
//     }
//     if(result.Length != 18) {
//      t.Error("Expected Length == 18")
//     }
// }

func TestSelect(t *testing.T) {
    {
        gs := GenSelect(1,3)
        data := []interface{}{0,1,2,3,4,5}
        result := gs(data).([]interface{})
        if len(result) != 2 {
            t.Error("result slice invalid size")
        }
        if result[0] != 1 {
            t.Error("Expected 1")
        }

        if result[1] != 3 {
            t.Error("Expected 3")
        }
    }

    {
        gs := GenSelect(3)
        data := []interface{}{0,1,2,3,4,5}
        result := gs(data).([]interface{})
        if len(result) != 1 {
            t.Error("result slice invalid size")
        }
        if result[0] != 3 {
            t.Error("Expected 3")
        }
    }
}


func TestSelect1(t *testing.T) {
    {
        gs := GenSelect1(2)
        data := []interface{}{0,1,2,3,4,5}
        result := gs(data)
        if result != 2 {
            t.Error("Expected 2")
        }
    }

    {
        gs := GenSelect1(2)
        data := []interface{}{"foo","bar","baz","foobar"}
        result := gs(data)
        if result != "baz" {
            t.Error("Expected baz")
        }
    }
}
