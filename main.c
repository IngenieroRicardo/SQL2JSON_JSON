#include <stdio.h>
#include <stdlib.h>
#include "SQL2JSON.h"  // Biblioteca para convertir resultados SQL a JSON
#include "JSON.h"      // Biblioteca para manipulación de JSON

/*
 * Función para imprimir una cadena sin las comillas dobles que la rodean (si existen)
 * @param str Cadena a imprimir (puede estar entre comillas)
 */
void imprimir_sin_comillas(char* str) {
    if (str == NULL) return;
    
    size_t len = strlen(str);
    // Verificar si la cadena comienza y termina con comillas dobles
    if (len >= 2 && str[0] == '"' && str[len-1] == '"') {
        // Imprimir la cadena sin las comillas iniciales y finales
        printf("%.*s", (int)(len-2), str+1);
    } else {
        printf("%s", str);
    }
}

/*
 * Función para mostrar un elemento JSON con formato legible
 * @param json_str Cadena JSON a mostrar
 * @param indice Número de índice del elemento (para mostrar en la salida)
 */
void mostrar_elemento_json(char* json_str, int indice) {
    // Primero validamos que el JSON sea correcto
    JsonResult parseado = ParseJSON(json_str);
    if (!parseado.is_valid) {
        printf("Error al parsear elemento: %s\n", parseado.error);
        FreeJsonResult(&parseado);
        return;
    }
    FreeJsonResult(&parseado);

    // Obtenemos todas las claves del objeto JSON
    JsonArrayResult claves = GetJSONKeys(json_str);
    if (!claves.is_valid) {
        printf("Error al obtener claves: %s\n", claves.error);
        FreeJsonArrayResult(&claves);
        return;
    }

    printf("Elemento %d:\n", indice + 1);
    
    // Iteramos sobre cada clave para mostrar su valor
    for (int i = 0; i < claves.count; i++) {
        char* clave = claves.items[i];
        JsonResult valor = GetJSONValue(json_str, clave);
        
        printf("  %s: ", clave);
        imprimir_sin_comillas(valor.value);
        printf("\n");
        
        FreeJsonResult(&valor);
    }
    printf("\n");
    
    FreeJsonArrayResult(&claves);
}

/*
 * Función principal que demuestra el flujo completo:
 * 1. Recibe la consulta SQL como argumento
 * 2. Conecta a una base de datos SQL
 * 3. Ejecuta la consulta
 * 4. Convierte el resultado a JSON
 * 5. Parsea y muestra el JSON de forma legible
 * 
 * Uso: ./programa "SELECT * FROM tabla"
 */
int main(int argc, char *argv[]) {
    // Validamos que se haya proporcionado la consulta SQL como argumento
    if (argc < 2) {
        printf("Uso: %s \"consulta_sql\"\n", argv[0]);
        printf("Ejemplo: %s \"SELECT usuarios, mensajes FROM chat.mensajeria\"\n", argv[0]);
        return 1;
    }

    // Configuración de conexión a la base de datos
    // Formato: usuario:contraseña@protocolo(ip:puerto)/basedatos
    char* conexion = "root:123456@tcp(127.0.0.1:3306)/mysql";
    
    // La consulta SQL se toma del primer argumento
    char* query = argv[1];
    printf("Ejecutando consulta: %s\n\n", query);
    
    // Ejecutamos la consulta y obtenemos los resultados en formato JSON
    char* json = SQLrun(conexion, query, 0, 0);
    if (json == NULL) {
        printf("Error al ejecutar la consulta SQL\n");
        return 1;
    }
    printf("Resultado JSON:\n%s\n\n", json);

    // Validamos que el JSON obtenido sea correcto
    JsonResult parseado = ParseJSON(json);
    if (!parseado.is_valid) {
        printf("Error al parsear JSON: %s\n", parseado.error);
        FreeJsonResult(&parseado);
        FreeString(json);
        return 1;
    }
    FreeJsonResult(&parseado);
    
    // Determinamos si el JSON es un array (comienza con '[') o un objeto simple
    int es_array = (json[0] == '[');
    
    if (es_array) {
        // Si es un array, obtenemos su longitud
        JsonResult longitud = GetArrayLength(json);
        if (!longitud.is_valid) {
            printf("Error al obtener longitud: %s\n", longitud.error);
            FreeJsonResult(&longitud);
            FreeString(json);
            return 1;
        }
        
        int num_elementos = atoi(longitud.value);
        FreeJsonResult(&longitud);
        
        printf("Total de elementos: %d\n\n", num_elementos);
        
        // Iteramos sobre cada elemento del array
        for (int i = 0; i < num_elementos; i++) {
            JsonResult elemento = GetArrayItem(json, i);
            if (!elemento.is_valid) {
                printf("Error al obtener elemento %d: %s\n", i, elemento.error);
                continue;
            }
            
            // Mostramos cada elemento con formato legible
            mostrar_elemento_json(elemento.value, i);
            
            FreeJsonResult(&elemento);
        }
    } else {
        // Si no es array, mostramos el objeto JSON directamente
        mostrar_elemento_json(json, 0);
    }

    // Liberamos la memoria utilizada por el JSON
    FreeString(json);
    
    return 0;
}
