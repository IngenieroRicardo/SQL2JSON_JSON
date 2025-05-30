#include <stdio.h>
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
        
        printf("  %s: %s\n", clave, valor.value);
        
        FreeJsonResult(&valor);
    }
    printf("\n");
    
    FreeJsonArrayResult(&claves);
}

void procesar_resultados_multiset(char* json) {
    // Primero verificamos si es un array de arrays (multi-resultset)
    if (json[0] == '[') {
        JsonArrayResult outer_array = GetArrayItems(json);
        if (!outer_array.is_valid) {
            printf("Error al parsear JSON principal: %s\n", outer_array.error);
            return;
        }

        printf("Se encontraron %d conjuntos de resultados:\n\n", outer_array.count);

        for (int set_idx = 0; set_idx < outer_array.count; set_idx++) {
            char* current_set = outer_array.items[set_idx];
            
            // Verificamos si este item es un array (resultset) u objeto individual
            if (current_set[0] == '[') {
                // Es un resultset (array de objetos)
                printf("=== Conjunto de Resultados %d ===\n", set_idx + 1);
                
                JsonArrayResult elementos = GetArrayItems(current_set);
                if (!elementos.is_valid) {
                    printf("Error al obtener elementos del conjunto %d: %s\n", set_idx+1, elementos.error);
                    continue;
                }

                for (int i = 0; i < elementos.count; i++) {
                    mostrar_elemento_json(elementos.items[i], i);
                }
                
                FreeJsonArrayResult(&elementos);
            } else {
                // Es un objeto individual (un solo registro)
                printf("=== Resultado Individual %d ===\n", set_idx + 1);
                mostrar_elemento_json(current_set, 0);
            }
            printf("\n");
        }
        
        FreeJsonArrayResult(&outer_array);
    } else {
        // Formato simple (objeto único)
        mostrar_elemento_json(json, 0);
    }
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Uso: %s \"consulta_sql\"\n", argv[0]);
        printf("Ejemplo: %s \"SELECT usuarios, mensajes FROM chat.mensajeria\"\n", argv[0]);
        return 1;
    }

    char* conexion = "root:123456@tcp(127.0.0.1:3306)/mysql";
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
