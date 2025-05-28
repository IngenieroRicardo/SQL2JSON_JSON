# SQL2JSON + JSON
Ejemplo usando SQL2JSON y JSON.
String de conexion en Linea 84 main.c: `char* conexion = "root:123456@tcp(192.100.1.210:3306)/mysql";`


### 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./JSON.so ./SQL2JSON.so` | `gcc -o main.exe main.c ./JSON.dll ./SQL2JSON.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./JSON.dll ./SQL2JSON.dll` |  |

---

### Respuesta de ejecucion

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



### Respuesta de consulta

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

### Respuesta erronea

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
