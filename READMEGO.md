READMEGO.md

# 🧅 buscadoronion.go

Versión en Go del buscador de sitios `.onion` (inspirado en buscaonion/buscadonion.py).

## 🚀 Requisitos

- Go 1.18 o superior
- Acceso a internet (para buscar en motores públicos)
- (Opcional) Tor si quieres navegar los .onion directamente desde Go

## 📦 Instalación de dependencias

```bash
go get github.com/PuerkitoBio/goquery
```

## 🏃‍♂️ Ejecución

```bash
go run main.go
```

## 🛠️ Funcionalidades

- Búsqueda en Ahmia por palabra clave.
- Filtrado de dominios .onion v3 válidos.
- Extracción de títulos y tecnologías básicas.
- Soporte para múltiples palabras clave y archivos.
- Resultados exportados a archivos `.txt`.
- Exportar resultados a JSON con títulos y tecnologías detectadas.
- Mostrar solo dominios .onion v2 o v3 de los resultados.
- Mostrar estadísticas de longitud de títulos.
- Buscar tecnologías específicas en los títulos de los resultados.
- **Resultados en memoria:** las funciones avanzadas operan sobre los últimos resultados obtenidos desde el menú interactivo.

## 📚 Ejemplo de uso

1. Ejecuta el menú principal:
   ```bash
   go run main.go
   ```
2. Usa la opción 1 para buscar y guardar resultados en memoria.
3. Usa las opciones 2–6 para operar sobre los últimos resultados encontrados.

---

**Nota:** Si quieres añadir más motores de búsqueda, edita la función `buscarEnMotores` en `buscadoronion.go`.

# 🛠️ Martillazonion-Go

Versión en Go de las herramientas de martillazonion, con mejoras y nuevas funcionalidades.

## Herramientas portadas y mejoradas

- **Ataque sencillo y multipuerto** (`martillazonion.go`, `martillazonion_multi.go`, `martillazo_alpha_simple.go`, `martillazo_alpha_multi.go`)
  - Soporte para múltiples objetivos desde archivo.
  - Estadísticas detalladas por objetivo y por puerto.
  - Opción de pausa/reanudar ataques en tiempo real (`p`).
  - Estadísticas en tiempo real (`s`).
  - Registro de ataques en archivo de log.
  - Mejor manejo de errores y reconexión automática.
  - Reintentos automáticos en conexiones.
  - Timestamps en cada ataque.
  - Ayuda extendida desde CLI.

- **Rastreador y buscador** (`crawleronion.go`, `browser_onion.go`)
  - Soporte para rastreo/búsqueda de múltiples URLs o palabras clave desde archivo.
  - Exporta todos los enlaces .onion únicos encontrados a TXT y JSON.
  - Estadísticas globales de rastreo/búsqueda.
  - Validación de enlaces antes de exportar.
  - Timestamps en cada operación.

- **Buscador avanzado** (`buscadoronion.go`)
  - Resultados en memoria para análisis posterior.
  - Exportación a JSON, filtrado v2/v3, estadísticas de títulos, búsqueda de tecnologías, etc.

- **Scripts de Tor** (`tor-multi-launcher.go`, `tor-multi-stop.go`)
  - Permite personalizar rango de puertos por argumentos o archivo.
  - Estadísticas de instancias lanzadas/detenidas.

- **Instalador de dependencias** (`instaladordependenciasonion.go`)
  - Imprime e instala dependencias Go automáticamente si el usuario lo desea.

## Ejecución

```bash
go run <tool>.go [opciones]
```

Ejemplo para ataque sencillo:
```bash
go run martillazonion.go -host dominio.onion -p 9050 -t 500 -r 100
go run martillazonion.go -targets objetivos.txt
```

Ejemplo para multipuerto:
```bash
go run martillazonion_multi.go
```

Ejemplo para rastreador:
```bash
go run crawleronion.go
```

Ejemplo para buscador avanzado:
```bash
go run buscadoronion.go
```

## Mejoras destacadas

- **Soporte para múltiples objetivos y palabras clave desde archivo** en la mayoría de herramientas.
- **Estadísticas detalladas** por objetivo, puerto, enlaces encontrados, etc.
- **Pausa y reanuda ataques** en tiempo real escribiendo `p` en consola.
- **Estadísticas en tiempo real** escribiendo `s` en consola.
- **Registro de ataques y errores** en archivos de log.
- **Exportación avanzada** (JSON, TXT, CSV) y análisis de resultados.
- **Menú principal y menú interactivo** para análisis y filtrado de resultados en buscadoronion.go.
- **Instalador de dependencias Go automático**.
- **Reintentos automáticos y timestamps** en ataques y búsquedas.
- **Validación de enlaces antes de exportar**.

## Notas

- Algunas herramientas requieren Tor corriendo en los puertos indicados.
- Los scripts de launcher/stop solo simulan los comandos, no ejecutan sudo.
- El buscadoronion.go es el más avanzado y soporta menú interactivo y exportación.
- Todos los forks respetan banners y funciones originales, solo se añaden mejoras.

---

**Fork realizado respetando banners, funciones y estructura de cada herramienta, añadiendo nuevas funcionalidades y mejoras.**

