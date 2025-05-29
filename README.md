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


### Consulta SELECT a tabla

```C
/main.bin "SELECT usuarios, mensajes FROM chat.mensajeria"
Ejecutando consulta: SELECT usuarios, mensajes FROM chat.mensajeria

Resultado JSON:
[
  {
    "mensajes": "Hola",
    "usuarios": "Ricardo"
  },
  {
    "mensajes": "hola",
    "usuarios": "yo"
  },
  {
    "mensajes": "hola",
    "usuarios": "mike"
  }
]

Total de elementos: 3

Elemento 1:
  mensajes: Hola
  usuarios: Ricardo

Elemento 2:
  mensajes: hola
  usuarios: yo

Elemento 3:
  mensajes: hola
  usuarios: mike

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




### 🧪 Ejemplo básico (consulta de una única fila)

```C
#include <stdio.h>
#include "SQL2JSON.h"
#include "JSON.h"

// Función para eliminar comillas de una cadena (si las tiene)
char* quitar_comillas(char* str) {
    if (str == NULL){
      return "";  
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
    char* conexion = "root:123456@tcp(192.100.1.210:3306)/test";
    
    // Consulta SQL con parámetros
    char* query = "select now();";
        
    // Llamar a la función
    char* json = SQLrun(conexion, query, 0, 0);
    
    // Mostrar valores sin comillas
    printf("now: %s\n", quitar_comillas(GetJSONValue(json, "now()").value));
        
    // Liberar memoria
    FreeString(json);
    
    return 0;
}
```

---
