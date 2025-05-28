package main

/*
#include <stdlib.h>
#include <string.h>

typedef struct {
    char* value;
    int is_valid;
    char* error;
} JsonResult;

typedef struct {
    char** items;
    int count;
    int is_valid;
    char* error;
} JsonArrayResult;
*/
import "C"
import (
	"encoding/json"
	"unsafe"
	"strconv"
	"fmt"
	"strings"
)

//export ParseJSON
func ParseJSON(jsonStr *C.char) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	var data interface{}
	err := json.Unmarshal([]byte(goStr), &data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(jsonBytes))
	result.error = nil
	return result
}

//export GetJSONValue
func GetJSONValue(jsonStr *C.char, key *C.char) C.JsonResult {
	goJsonStr := C.GoString(jsonStr)
	goKey := C.GoString(key)
	var result C.JsonResult

	var data map[string]interface{}
	err := json.Unmarshal([]byte(goJsonStr), &data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	value, exists := data[goKey]
	if !exists {
		result.is_valid = 0
		result.error = C.CString("key not found")
		return result
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(jsonBytes))
	result.error = nil
	return result
}

//export GetArrayLength
func GetArrayLength(jsonStr *C.char) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	if len(goStr) == 0 || goStr[0] != '[' {
		result.is_valid = 0
		result.error = C.CString("not a JSON array")
		return result
	}

	var arr []interface{}
	err := json.Unmarshal([]byte(goStr), &arr)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	lengthStr := strconv.Itoa(len(arr))
	
	result.is_valid = 1
	result.value = C.CString(lengthStr)
	result.error = nil
	return result
}

//export GetArrayItem
func GetArrayItem(jsonStr *C.char, index int) C.JsonResult {
	goStr := C.GoString(jsonStr)
	var result C.JsonResult

	var arr []interface{}
	err := json.Unmarshal([]byte(goStr), &arr)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	if index < 0 || index >= len(arr) {
		result.is_valid = 0
		result.error = C.CString("index out of bounds")
		return result
	}

	itemBytes, err := json.Marshal(arr[index])
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	result.is_valid = 1
	result.value = C.CString(string(itemBytes))
	result.error = nil
	return result
}

//export FreeJsonResult
func FreeJsonResult(result *C.JsonResult) {
	if result.value != nil {
		C.free(unsafe.Pointer(result.value))
	}
	if result.error != nil {
		C.free(unsafe.Pointer(result.error))
	}
}

//export GetJSONKeys
func GetJSONKeys(jsonStr *C.char) C.JsonArrayResult {
	goJsonStr := C.GoString(jsonStr)
	var result C.JsonArrayResult

	var data map[string]interface{}
	err := json.Unmarshal([]byte(goJsonStr), &data)
	if err != nil {
		result.is_valid = 0
		result.error = C.CString(err.Error())
		return result
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	// Allocate C memory for the array
	cArray := C.malloc(C.size_t(len(keys)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	cKeys := (*[1<<30 - 1]*C.char)(cArray)

	for i, key := range keys {
		cKeys[i] = C.CString(key)
	}

	result.is_valid = 1
	result.items = (**C.char)(cArray)
	result.count = C.int(len(keys))
	result.error = nil
	return result
}

//export FreeJsonArrayResult
func FreeJsonArrayResult(result *C.JsonArrayResult) {
	if result.items != nil {
		// Convert to Go slice to free each string
		cKeys := (*[1<<30]*C.char)(unsafe.Pointer(result.items))[:result.count:result.count]
		for i := 0; i < int(result.count); i++ {
			C.free(unsafe.Pointer(cKeys[i]))
		}
		C.free(unsafe.Pointer(result.items))
	}
	if result.error != nil {
		C.free(unsafe.Pointer(result.error))
	}
}

//export GetJSONValueByPath
func GetJSONValueByPath(jsonStr *C.char, path *C.char) C.JsonResult {
    goJsonStr := C.GoString(jsonStr)
    goPath := C.GoString(path)
    var result C.JsonResult

    var data interface{}
    err := json.Unmarshal([]byte(goJsonStr), &data)
    if err != nil {
        result.is_valid = 0
        result.error = C.CString(err.Error())
        return result
    }

    current := data
    pathParts := strings.Split(goPath, ".")
    for _, part := range pathParts {
        if part == "" {
            continue
        }

        switch v := current.(type) {
        case map[string]interface{}:
            val, exists := v[part]
            if !exists {
                result.is_valid = 0
                result.error = C.CString(fmt.Sprintf("path '%s' not found", part))
                return result
            }
            current = val
        case []interface{}:
            index, err := strconv.Atoi(part)
            if err != nil || index < 0 || index >= len(v) {
                result.is_valid = 0
                result.error = C.CString(fmt.Sprintf("invalid array index '%s'", part))
                return result
            }
            current = v[index]
        default:
            result.is_valid = 0
            result.error = C.CString(fmt.Sprintf("cannot traverse path '%s'", part))
            return result
        }
    }

    jsonBytes, err := json.Marshal(current)
    if err != nil {
        result.is_valid = 0
        result.error = C.CString(err.Error())
        return result
    }

    result.is_valid = 1
    result.value = C.CString(string(jsonBytes))
    result.error = nil
    return result
}

//export GetArrayItems
func GetArrayItems(jsonStr *C.char) C.JsonArrayResult {
    goStr := C.GoString(jsonStr)
    var result C.JsonArrayResult

    // Verificar si es un array JSON válido
    if len(goStr) == 0 || goStr[0] != '[' {
        result.is_valid = 0
        result.error = C.CString("not a JSON array")
        return result
    }

    var arr []interface{}
    err := json.Unmarshal([]byte(goStr), &arr)
    if err != nil {
        result.is_valid = 0
        result.error = C.CString(err.Error())
        return result
    }

    // Allocate memory for the C array of strings
    cArray := C.malloc(C.size_t(len(arr)) * C.size_t(unsafe.Sizeof(uintptr(0))))
    cItems := (*[1<<30 - 1]*C.char)(unsafe.Pointer(cArray))

    for i, item := range arr {
        itemBytes, err := json.Marshal(item)
        if err != nil {
            // Liberar memoria ya asignada si hay error
            for j := 0; j < i; j++ {
                C.free(unsafe.Pointer(cItems[j]))
            }
            C.free(cArray)
            
            result.is_valid = 0
            result.error = C.CString(err.Error())
            return result
        }
        cItems[i] = C.CString(string(itemBytes))
    }

    result.is_valid = 1
    result.items = (**C.char)(cArray)
    result.count = C.int(len(arr))
    result.error = nil
    return result
}


func main() {}