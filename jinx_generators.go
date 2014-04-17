package jinx

import (
    "fmt"
    "reflect"
)

// TODO: squash ConcatParams + ConcatArray into the same functions
//       just use an array instead of a param list
func ConcatParams(s interface{}) interface{} {
    var a string
    ss,ok := s.([]interface{})
    if !ok {
        return ""
    }
    for i := range ss {
        if v,ok := ss[i].(string); ok {
            a += v
        } else {
            fmt.Println("Invalid result")
        }
    }
    return a
}

// func decStringResult(s ...interface{}) interface{} {
//     var a string
//     for i := range s {
//         a += "<<"
//         a += s[i].(string)
//         a += ">>"
//     }
//     return a
// }


func IgnoreParams(s interface{}) interface{} {
    return ""
}

func Identity(s interface{}) interface{} {
    return s
}

func ConcatArray(s interface{}) interface{} {
    fmt.Println(reflect.TypeOf(s))
    ss,ok := s.([]interface{})
    if !ok {
        return ""
    }
    var a string// TODO: inefficient
    s0 := (ss[0]).([]interface{})
    for i,_ := range s0 {
        if v, ok := s0[i].(string); ok {
            a += v
        } else {
            fmt.Println("Invalid type")
        }
    }
    return a
}