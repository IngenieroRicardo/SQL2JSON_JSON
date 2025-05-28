# SQL2JSON + JSON

Ejemplo de uso de SQL2JSON y JSON para interactuar con bases de datos MySQL.

## Configuración

La cadena de conexión se encuentra en la línea 122 de `main.c`:
```c
char* conexion = "root:123456@tcp(127.0.0.1:3306)/mysql";
```


### 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./JSON.so ./SQL2JSON.so` | `gcc -o main.exe main.c ./JSON.dll ./SQL2JSON.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./JSON.dll ./SQL2JSON.dll` |  |

---

### Ejecución de comando DELETE

```bash
./main.bin "DELETE FROM chat.usuario"
Ejecutando consulta: DELETE FROM chat.usuario

Resultado JSON:
{
  "status": "OK"
}

Elemento 1:
  status: OK
```



### Consulta SELECT simple

```C
./main.bin "SELECT now()"
Ejecutando consulta: SELECT now()

Resultado JSON:
{
  "now()": "2025-05-28 16:22:53"
}

Elemento 1:
  now(): 2025-05-28 16:22:53

```

---

### Manejo de errores

```C
./main.bin "SELECT funcion()"
Ejecutando consulta: SELECT funcion()

Resultado JSON:
{
  "error": "Error en la consulta SQL: Error 1305 (42000): FUNCTION mysql.funcion does not exist"
}

Elemento 1:
  error: Error en la consulta SQL: Error 1305 (42000): FUNCTION mysql.funcion does not exist
```




### 🧪 Ejemplo básico de una unica fila

```C
#include <stdio.h>
#include "SQL2JSON.h"
#include "JSON.h"

// Función para eliminar comillas de una cadena (si las tiene)
char* quitar_comillas(char* str) {
    if (str == NULL){
      return NULL;  
    } 
    size_t len = strlen(str);
    if (len >= 2 && str[0] == '"' && str[len-1] == '"') {
        str[len-1] = '\0';  // Eliminar comilla final
        return str + 1;      // Saltar comilla inicial
    }
    return str;
}

int main() {
    // Configuración de conexión
    char* conexion = "root:123456@tcp(127.0.0.1:3306)/test";
    
    // Consulta SQL con parámetros
    char* query = "select now();";
        
    // Llamar a la función
    char* json = SQLrun(conexion, query, 0, 0);
    
    // Analizar JSON
    JsonResult resultado = ParseJSON(json);
    
    if (!resultado.is_valid) {
        printf("Error: %s\n", resultado.error);
        return 1;
    }
    
    // Mostrar valores sin comillas
    printf("now: %s\n", quitar_comillas(GetJSONValue(json, "now()").value));
        
    // Liberar memoria
    FreeJsonResult(&resultado);
    FreeString(json);
    
    return 0;
}
```

---
