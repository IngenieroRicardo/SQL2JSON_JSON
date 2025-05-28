package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
	_ "github.com/go-sql-driver/mysql"
)

// ErrorResponse representa una respuesta de error estandarizada
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse representa una respuesta exitosa para operaciones sin resultado
type SuccessResponse struct {
	Status string `json:"status"`
}

//export SQLrun
func SQLrun(conexion *C.char, query *C.char, args **C.char, argCount C.int) *C.char {
	goConexion := C.GoString(conexion)
	goQuery := C.GoString(query)

	var goArgs []interface{}
	if argCount > 0 {
		argSlice := (*[1 << 30]*C.char)(unsafe.Pointer(args))[:argCount:argCount]
		for _, arg := range argSlice {
			argStr := C.GoString(arg)

			switch {
			case strings.HasPrefix(argStr, "int::"):
				intVal, err := strconv.ParseInt(argStr[5:], 10, 64)
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Error parseando entero: %s", argStr[5:]))
				}
				goArgs = append(goArgs, intVal)

			case strings.HasPrefix(argStr, "float::"), strings.HasPrefix(argStr, "double::"):
				prefixLen := 7
				if strings.HasPrefix(argStr, "double::") {
					prefixLen = 8
				}
				floatVal, err := strconv.ParseFloat(argStr[prefixLen:], 64)
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Error parseando float: %s", argStr[prefixLen:]))
				}
				goArgs = append(goArgs, floatVal)

			case strings.HasPrefix(argStr, "bool::"):
				boolVal, err := strconv.ParseBool(argStr[6:])
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Error parseando booleano: %s", argStr[6:]))
				}
				goArgs = append(goArgs, boolVal)

			case strings.HasPrefix(argStr, "null::"):
				goArgs = append(goArgs, nil)

			case strings.HasPrefix(argStr, "blob::"):
				data, err := base64.StdEncoding.DecodeString(argStr[6:])
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Error decodificando blob: %v", err))
				}
				goArgs = append(goArgs, data)

			default:
				goArgs = append(goArgs, argStr)
			}
		}
	}

	_, result := sqlRunInternal(goConexion, goQuery, goArgs...)
	return C.CString(result)
}

func sqlRunInternal(conexion, query string, args ...any) (int, string) {
	db, err := sql.Open("mysql", conexion)
	if err != nil {
		return 1, createErrorJSON(fmt.Sprintf("Error al abrir conexión: %v", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return 1, createErrorJSON(fmt.Sprintf("Error al conectar a la base de datos: %v", err))
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return 1, createErrorJSON(fmt.Sprintf("Error en la consulta SQL: %v", err))
	}
	defer rows.Close()

	var flatResults []map[string]interface{}

	for {
		columns, err := rows.Columns()
		if err != nil {
			return 1, createErrorJSON(fmt.Sprintf("Error al obtener columnas: %v", err))
		}

		colTypes, err := rows.ColumnTypes()
		if err != nil {
			return 1, createErrorJSON(fmt.Sprintf("Error al obtener tipos de columna: %v", err))
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				return 1, createErrorJSON(fmt.Sprintf("Error al escanear fila: %v", err))
			}

			rowData := make(map[string]interface{})
			for i, col := range values {
				colName := columns[i]
				if col == nil {
					rowData[colName] = nil
				} else {
					if strings.Contains(colTypes[i].DatabaseTypeName(), "BLOB") {
						rowData[colName] = base64.StdEncoding.EncodeToString(col)
					} else {
						strValue := string(col)
						rowData[colName] = strings.ReplaceAll(strValue, "\"", "'")
					}
				}
			}
			flatResults = append(flatResults, rowData)
		}

		if err = rows.Err(); err != nil {
			return 1, createErrorJSON(fmt.Sprintf("Error después de iterar filas: %v", err))
		}

		if !rows.NextResultSet() {
			break
		}
	}

	if len(flatResults) == 0 && isNonReturningQuery(query) {
		return 0, createSuccessJSON()
	}

	jsonData, err := json.MarshalIndent(flatResults, "", "  ")
	if err != nil {
		return 1, createErrorJSON(fmt.Sprintf("Error al convertir resultados a JSON: %v", err))
	}

	return 0, string(jsonData)
}

func isNonReturningQuery(query string) bool {
	queryUpper := strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(queryUpper, "INSERT ") ||
		strings.HasPrefix(queryUpper, "UPDATE ") ||
		strings.HasPrefix(queryUpper, "DELETE ") ||
		strings.HasPrefix(queryUpper, "DROP ") ||
		strings.HasPrefix(queryUpper, "CREATE ") ||
		strings.HasPrefix(queryUpper, "ALTER ") ||
		strings.HasPrefix(queryUpper, "TRUNCATE ") ||
		strings.HasPrefix(queryUpper, "CALL ")
}

func createErrorResponse(message string) *C.char {
	return C.CString(createErrorJSON(message))
}

func createErrorJSON(message string) string {
	errResp := ErrorResponse{Error: message}
	jsonData, _ := json.MarshalIndent(errResp, "", "  ")
	return string(jsonData)
}

func createSuccessJSON() string {
	successResp := SuccessResponse{Status: "OK"}
	jsonData, _ := json.MarshalIndent(successResp, "", "  ")
	return string(jsonData)
}

//export FreeString
func FreeString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func main() {}