#include <stdio.h>
#include <stdlib.h>
#include "SQL2JSON.h"
#include "JSON.h"

void mostrar_elemento_json(char* json_str, int indice) {
    JsonResult parseado = ParseJSON(json_str);
    if (!parseado.is_valid) {
        printf("Error al parsear elemento: %s\n", parseado.error);
        FreeJsonResult(&parseado);
        return;
    }
    FreeJsonResult(&parseado);

    JsonArrayResult claves = GetJSONKeys(json_str);
    if (!claves.is_valid) {
        printf("Error al obtener claves: %s\n", claves.error);
        FreeJsonArrayResult(&claves);
        return;
    }

    printf("Elemento %d:\n", indice + 1);
    
    for (int i = 0; i < claves.count; i++) {
        char* clave = claves.items[i];
        JsonResult valor = GetJSONValue(json_str, clave);
        
        printf("  %s: ", clave);
        printf("  %s: ", valor.value);
        printf("\n");
        
        FreeJsonResult(&valor);
    }
    printf("\n");
    
    FreeJsonArrayResult(&claves);
}

void procesar_resultados_multiset(char* json) {
    JsonResult parseado = ParseJSON(json);
    if (!parseado.is_valid) {
        printf("Error al parsear JSON: %s\n", parseado.error);
        FreeJsonResult(&parseado);
        return;
    }
    FreeJsonResult(&parseado);

    // Primero verificamos si es un objeto con múltiples resultsets
    JsonResult results = GetJSONValue(json, "results");
    if (results.is_valid) {
        // Es el formato MultiResultResponse con campo "results"
        printf("Se encontraron múltiples conjuntos de resultados:\n\n");
        
        JsonArrayResult result_sets = GetArrayItems(results.value);
        if (!result_sets.is_valid) {
            printf("Error al obtener conjuntos de resultados: %s\n", result_sets.error);
            FreeJsonResult(&results);
            return;
        }

        for (int set_idx = 0; set_idx < result_sets.count; set_idx++) {
            printf("=== Conjunto de Resultados %d ===\n", set_idx + 1);
            
            JsonArrayResult elementos = GetArrayItems(result_sets.items[set_idx]);
            if (!elementos.is_valid) {
                printf("Error al obtener elementos del conjunto %d: %s\n", set_idx+1, elementos.error);
                continue;
            }

            for (int i = 0; i < elementos.count; i++) {
                mostrar_elemento_json(elementos.items[i], i);
            }
            
            FreeJsonArrayResult(&elementos);
            printf("\n");
        }
        
        FreeJsonArrayResult(&result_sets);
    } else {
        // Formato simple (array u objeto único)
        if (json[0] == '[') {
            JsonArrayResult elementos = GetArrayItems(json);
            if (!elementos.is_valid) {
                printf("Error al obtener elementos: %s\n", elementos.error);
                return;
            }

            printf("Total de elementos: %d\n\n", elementos.count);
            for (int i = 0; i < elementos.count; i++) {
                mostrar_elemento_json(elementos.items[i], i);
            }
            
            FreeJsonArrayResult(&elementos);
        } else {
            mostrar_elemento_json(json, 0);
        }
    }

    if (results.is_valid) {
        FreeJsonResult(&results);
    }
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Uso: %s \"consulta_sql\"\n", argv[0]);
        printf("Ejemplo: %s \"SELECT usuarios, mensajes FROM chat.mensajeria\"\n", argv[0]);
        return 1;
    }

    char* conexion = "root:123456@tcp(192.100.1.210:3306)/mysql";
    char* query = argv[1];
    printf("Ejecutando consulta: %s\n\n", query);
    
    SQLResult resultado = SQLrun(conexion, query, NULL, 0);
    if (resultado.json == NULL) {
        printf("Error al ejecutar la consulta SQL\n");
        return 1;
    }
    printf("Resultado JSON:\n%s\n\n", resultado.json);

    procesar_resultados_multiset(resultado.json);

    FreeSQLResult(&resultado);
    return 0;
}
