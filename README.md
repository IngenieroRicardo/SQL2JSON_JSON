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
[
  {
    "now()": "2025-05-28 13:50:42"
  }
]

Total de elementos: 1

Elemento 1:
  now(): 2025-05-28 13:50:42

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




---
