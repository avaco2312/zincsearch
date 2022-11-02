## Prueba técnica de interés: Indexar BD de emails en ZincSearch.

**Definición del problema:**

Crear una interfaz para buscar información en una base de datos de correos electrónicos.

La primera parte es indexar la base de datos de prueba y la segunda construir una interfaz para consultarla.

Parte 1: Indexar base de datos de correos electrónicos

Primero descargar la base de datos de correos de Enron Corp:
http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz (423MB)

Después escribir un programa que indexe sus contenidos en la herramienta ZincSearch: https://zincsearch.com/

Parte 2: Visualizador
Crea una interfaz simple para buscar los contenidos.

**Solución, parte 1 y algo de la parte 2**

Para la parte 1 ver main.go en /index. La solución es concurrente y su eficiencia está limitada por el nivel de I/O alcanzable al leer los archivos (517423 archivos, 2.23 Gb). Solución "quick and dirty" :=)

Para "algo" de la parte 2: main.go en /query. Hace la búsqueda de un término (embebido). Debe ampliarse a un servidor de API que permita consultas generales y crear una interfase de usuario.