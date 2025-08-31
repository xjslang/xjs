# XJS Build Commands

Este proyecto usa [Mage](https://magefile.org/) como herramienta de automatización, proporcionando una alternativa moderna y type-safe a los Makefiles tradicionales.

## Instalación

```bash
go install github.com/magefile/mage@latest
```

## Comandos Disponibles

### Testing

```bash
# Ejecutar todos los tests de transpilación (equivalente al script original)
mage test                    # Target por defecto - tests de integración
mage                         # Equivalente a 'mage test'

# Tests por tipo
mage testUnit               # Solo tests unitarios (rápidos)
mage testAll                # Todos los tests (unitarios + integración)

# Tests específicos de integración
mage testTranspilation       # Tests con fixtures del directorio testdata/
mage testInline             # Tests inline de transpilación básica
mage testErrors             # Tests de manejo de errores
mage testMiddleware         # Tests del sistema de middleware
```

## Organización de Tests

El proyecto sigue la **estructura de Kubernetes** para organizar diferentes tipos de tests:

### 📁 **Estructura de Directorios**
```
xjs/
├── parser/
│   └── parser_test.go           # Tests unitarios del parser
├── lexer/
│   └── lexer_test.go           # Tests unitarios del lexer  
├── ast/
│   └── ast_test.go             # Tests unitarios del AST
├── token/
│   └── token_test.go           # Tests unitarios de tokens
└── test/
    ├── integration/            # Tests de integración
    │   ├── transpilation_test.go
    │   ├── middleware_test.go
    │   └── examples_test.go
    └── testdata/               # Fixtures compartidas
        ├── basic.js
        ├── basic.out
        └── ...
```

### 🧪 **Tests por tipo**
```bash
# Tests unitarios (rápidos, sin dependencias externas)
mage testUnit               # Solo tests unitarios, no requiere Node.js
go test ./parser ./lexer    # Equivalente directo

# Tests de integración (requieren Node.js)
mage test                  # Tests de integración (requiere Node.js)
go test ./test/integration  # Equivalente directo
```

### 🎯 **Tests específicos de integración**
```bash
mage testTranspilation       # Tests con fixtures del directorio testdata/
mage testInline             # Tests inline de transpilación básica
mage testErrors             # Tests de manejo de errores
mage testMiddleware         # Tests del sistema de middleware
```

## Ventajas de la Organización Kubernetes

1. **Separación clara**: Tests unitarios cerca del código, integración centralizada
2. **Escalabilidad**: Fácil agregar nuevos tipos de tests (e2e, performance, etc.)
3. **Navegación**: Fácil localizar tests por tipo y propósito
4. **CI eficiente**: Pipelines optimizados por tipo de test
5. **Estándar profesional**: Usado por proyectos enterprise

### Benchmarking

```bash
mage bench                  # Todos los benchmarks
mage benchTranspilation     # Solo benchmarks de transpilación
```

### Build y Release

```bash
mage build                  # Compilar el proyecto
mage clean                  # Limpiar archivos generados
mage install                # Instalar dependencias
mage tidy                   # Limpiar go.mod
mage lint                   # Ejecutar linter (requiere golangci-lint)
mage release                # Pipeline completo de release
mage ci                     # Pipeline de CI
```

### Desarrollo

```bash
mage dev                    # Modo desarrollo con auto-testing (requiere watchexec)
```

## Ventajas sobre Scripts

1. **Type-safe**: Código Go con validación de tipos
2. **Multiplataforma**: Funciona igual en todos los OS
3. **Dependencias**: Manejo automático de dependencias entre tareas
4. **Documentación**: `mage -l` muestra todas las tareas disponibles
5. **IDE Support**: IntelliSense y debugging

## Migración desde el Script

El comando `mage test` es un reemplazo completo de `./test_transpilation.sh`:

### Antes:
```bash
./test_transpilation.sh
```

### Ahora:
```bash
mage test
# o simplemente
mage
```

## Ejemplos de Uso

```bash
# Desarrollo diario
mage test                   # Ejecutar tests rápidos
mage testTranspilation      # Solo tests de transpilación
mage build                  # Compilar

# Antes de commit
mage lint                   # Verificar código
mage testAll               # Todos los tests

# Preparar release
mage release               # Pipeline completo
```

## Personalización

Para añadir nuevos comandos, edita `magefile.go` y agrega nuevas funciones públicas. Mage las detectará automáticamente.

---

**Nota**: El archivo `test_transpilation.sh` original se mantiene como referencia, pero se recomienda usar `mage test` para nuevos desarrollos.
