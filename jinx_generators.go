package jinx

import (
    "fmt"
    "strconv"
    "reflect"
)

func GenString(s interface{}) interface{} {
    var a string
    ss,ok := s.([]interface{})
    if !ok {
        return ""
    }
    for i := range ss {
        if v,ok := ss[i].(string); ok {
            a += v
        } else {
            fmt.Printf("GenString expecting []string, got %s instead\n", reflect.TypeOf(ss[i]))
        }
    }
    return a
}

func GenStringToInt(s interface{}) interface{} {
    str, ok := s.(string)
    if !ok {
        // TODO: error handling in Gens
        fmt.Println("Error in GenStringToInt")
        return nil
    }
    r,_ := strconv.Atoi(str)
    return r
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


func GenIgnoreParams(s interface{}) interface{} {
    return ""
}

func GenIdentity(s interface{}) interface{} {
    return s
}

func GenSelect(idxs ...int) func(interface{}) interface{} {
    return func(s interface{}) interface{} {
        num_results := len(idxs)
        results := make([]interface{}, num_results)
        ss, ok := s.([]interface{})
        if !ok {
            // TODO
            fmt.Println("Type conversion failure")
            return "FAIL"
        }
        result_position := 0
        for i := range idxs {
            results[result_position] = ss[idxs[i]]
            result_position++
        }
        return results
    }
}

func GenSelect1(idx int) func(interface{}) interface{} {
    return func(s interface{}) interface{} {
        ss, ok := s.([]interface{})
        if !ok {
            return "FAIL"
        }
        return ss[idx]
    }
}

