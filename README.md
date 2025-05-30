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


### Consulta SELECT a tabla

```C
./main.bin "SELECT usuarios, mensajes FROM chat.mensajeria"
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

int main() {
    // Ejemplo de conexión y consulta
    char* conexion = "root:123456@tcp(127.0.0.1:3306)/test";
    char* query = "SELECT now();";
    
    SQLResult resultado = SQLrun(conexion, query, NULL, 0);
    
    if (resultado.is_error) {
        printf("Error: %s\n", resultado.json);
        return 1;
    } else if (resultado.is_empty) {
        printf("Consulta ejecutada pero no retornó datos\n");
        printf("JSON: %s\n", resultado.json); // Mostrará {"status":"OK"} o []
        return 1;
    }

    JsonResult now = GetJSONValueByPath(resultado.json, "0.now()");
    
    // Mostrar valores sin comillas
    printf("Now: %s\n", now.value);
    
    // Liberar memoria
    FreeJsonResult(&now);
    FreeSQLResult(&resultado);
    
    return 0;
}
```

---
