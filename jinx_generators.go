package jinx

import (
	"fmt"
	"reflect"
	"strconv"
)

func GenUnwrap(s interface{}) interface{} {
	if vv, ok := s.([]interface{}); ok {
		result := make([]interface{}, 0)
		for i := range vv {
			v := vv[i]
			fmt.Println(v)
			result = append(result, v)
		}
		return result
	} else if v, ok := s.(interface{}); ok {
		return v
	}
	return ""
}

func GenString(s interface{}) interface{} {
	/// TODO: use a byte buffer etc

	var a string
	if ss, ok := s.([]interface{}); ok {
		for i := range ss {
			if v, ok := ss[i].(string); ok {
				a += v
			} else {
				fmt.Printf("GenString expecting []interface{}, got %s instead\n", reflect.TypeOf(ss[i]))
			}
		}
	} else if ss, ok := s.([]string); ok {
		for i := range ss {
			v := ss[i]
			a += v
		}
	} else {
		// TODO: work on error reporting from Gens
		fmt.Println("Unknown type for GenString")
	}
	return a
}

//func GenDec(s interface{}) interface{} {
//	ss, ok := s.(string)
//	if !ok {
//		return "<<ERROR>>"
//	} else {
//		return "<<" + ss + ">>"
//	}
//}

func GenListOfStrings(s interface{}) interface{} {
	var a []string
	if ss, ok := s.([]interface{}); ok {
		for i := range ss {
			if v, ok := ss[i].(string); ok {
				a = append(a, v)
			} else {
				fmt.Printf("GenListOfStrings expecting []string, got %s instead\n", reflect.TypeOf(ss[i]))
			}
		}
	} else if ss, ok := s.([]string); ok {
		for i := range ss {
			v := ss[i]
			a = append(a, v)
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
	r, _ := strconv.Atoi(str)
	return r
}

func GenDebug(s interface{}) interface{} {
	fmt.Printf("DEBUG: %#v\n", s)
	return s
}

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
