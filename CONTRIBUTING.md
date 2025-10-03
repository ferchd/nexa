# Guía de Contribución

¡Gracias por considerar contribuir a Nexa!

## 📋 Tabla de Contenidos

- [Código de Conducta](#código-de-conducta)
- [¿Cómo Puedo Contribuir?](#cómo-puedo-contribuir)
- [Guías de Estilo](#guías-de-estilo)
- [Proceso de Desarrollo](#proceso-de-desarrollo)
- [Reportar Bugs](#reportar-bugs)
- [Solicitar Features](#solicitar-features)

## Código de Conducta

Este proyecto se adhiere al [Código de Conducta](CODE_OF_CONDUCT.md). Al participar, se espera que mantengas este código.

## ¿Cómo Puedo Contribuir?

### Reportando Bugs

Antes de crear un reporte de bug, por favor:

1. **Verifica** que no exista un issue similar
2. **Recopila** información detallada sobre el problema
3. **Incluye** pasos para reproducir

Usa el [template de bug report](.github/ISSUE_TEMPLATE/bug_report.md).

### Solicitando Features

Antes de solicitar un feature:

1. **Verifica** que no exista una solicitud similar
2. **Describe** claramente el caso de uso
3. **Explica** por qué este feature sería útil

Usa el [template de feature request](.github/ISSUE_TEMPLATE/feature_request.md).

### Pull Requests

1. Fork el repositorio
2. Crea una rama desde `develop`
3. Realiza tus cambios
4. Asegúrate de que los tests pasen
5. Actualiza la documentación si es necesario
6. Envía un pull request

## Guías de Estilo

### Git Commit Messages

Usamos [Conventional Commits](https://www.conventionalcommits.org/):

<type>[optional scope]: <description>

[optional body]

[optional footer(s)]

Types:

- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Solo cambios en documentación
- `style`: Cambios que no afectan el significado del código
- `refactor`: Código que no corrige bugs ni agrega features
- `perf`: Mejoras de rendimiento
- `test`: Agregar o corregir tests
- `chore`: Cambios en el proceso de build o herramientas auxiliares

### Ejemplos:

```
feat(checker): add UDP connectivity check
fix(tcp): correct port conversion bug
docs(readme): add installation instructions
test(checker): improve test coverage for edge cases
```

### Go Code Style

Seguimos las Go Code Review Comments:

```go
// ✅ Good
func CheckTCP(host string, port int, timeout time.Duration) bool {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return false
    }
    defer conn.Close()
    return true
}

// ❌ Bad
func check_tcp(h string, p int) bool { // snake_case, abreviaciones
    // missing timeout
    conn, _ := net.Dial("tcp", h + ":" + string(p)) // ignora error, conversión incorrecta
    return conn != nil
}
```

### Convenciones de Código

- **Nombres**: Usa `camelCase` para variables locales, `PascalCase` para exportados
- **Errores**: Siempre verifica y maneja errores
- **Comentarios**: Documenta funciones exportadas
- **Tests**: Escribe tests para toda nueva funcionalidad
- **Formateo**: Usa `gofmt` o `goimports`

### Estructura de Tests

```go
func TestFunctionName(t *testing.T) {
    testCases := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
        },
        {
            name:     "invalid input",
            input:    invalidInput,
            expected: expectedError,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := FunctionName(tc.input)
            if result != tc.expected {
                t.Errorf("Expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```

## Proceso de Desarrollo

### Setup Local

```bash
# 1. Fork y clonar
git clone https://github.com/TU_USUARIO/nexa.git
cd nexa

# 2. Agregar upstream
git remote add upstream https://github.com/ferchd/nexa.git

# 3. Crear rama de desarrollo
git checkout -b feature/mi-feature develop

# 4. Instalar dependencias
go mod download

# 5. Instalar herramientas de desarrollo
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### Workflow de Desarrollo

```bash
# 1. Mantener tu fork actualizado
git fetch upstream
git checkout develop
git merge upstream/develop

# 2. Crear branch para tu feature
git checkout -b feature/nueva-funcionalidad develop

# 3. Hacer cambios y commits
git add .
git commit -m "feat(scope): descripción del cambio"

# 4. Ejecutar tests
make test
make lint
make security

# 5. Push a tu fork
git push origin feature/nueva-funcionalidad

# 6. Crear Pull Request en GitHub
```

### Antes de Enviar PR

```bash
# Ejecutar suite completa de checks
make test          # Tests unitarios
make lint          # Linter
make security      # Security scan
make coverage      # Coverage report

# Verificar que compile en todas las plataformas
GOOS=linux make build
GOOS=windows make build
GOOS=darwin make build
```

## Reportar Bugs

### Información Necesaria

**Versión de Nexa**: `nexa --version`
**Sistema Operativo**: OS y versión
**Versión de Go**: `go version`
**Configuración**: Archivo config.yaml (sin datos sensibles)
**Logs**: Logs relevantes
**Pasos para reproducir**: Detallados

### Ejemplo de Bug Report

```markdown
**Versión**: v1.0.0
**OS**: Ubuntu 22.04 LTS
**Go**: 1.21.0

**Descripción**:
TCP check falla con timeout aunque el puerto esté abierto.

**Pasos para reproducir**:
1. Configurar host: `8.8.8.8:53`
2. Ejecutar: `nexa --external 8.8.8.8:53`
3. Observar timeout

**Comportamiento esperado**:
Check exitoso

**Logs**:
[ERROR] TCP check failed: i/o timeout

**Configuración**:
```yaml
tcp_timeout: 2s
attempts: 2
```

## Solicitar Features

### Template de Feature Request

1. **Problema que resuelve**: Describe el problema actual
2. **Solución propuesta**: Cómo debería funcionar el feature
3. **Alternativas**: Otras soluciones consideradas
4. **Contexto adicional**: Screenshots, diagramas, etc.

### Ejemplo
```markdown
**Problema**:
No hay forma de verificar conectividad UDP (ej: DNS sobre UDP)

**Solución propuesta**:
Agregar `CheckUDP()` similar a `CheckTCP()`

**Alternativas**:
- Usar herramienta externa como `nc -u`
- Solo verificar DNS con resolución

**Contexto**:
Muchos servicios usan UDP (DNS, NTP, SNMP)
```

## Tests

### Ejecutar Tests

```bash
# Todos los tests
go test ./...

# Test específico
go test -v -run TestCheckTCP ./internal/checker/

# Con coverage
go test -cover ./...

# Con race detection
go test -race ./...

# Generar reporte HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Escribir Tests

- Usa mocks para dependencias externas
- Tests deben ser deterministas
- Usa table-driven tests para múltiples casos
- Nombra tests descriptivamente: `TestFunction_Scenario_ExpectedResult`

## Documentación

### Actualizar Documentación

Al agregar features, actualiza:

[] README.md
[] ARCHITECTURE.md (si aplica)
[] Comentarios en código
[] Examples/ (si aplica)
[] Wiki (para guías extensas)

### Formato de Documentación

- Usa Markdown
- Incluye ejemplos de código
- Agrega diagramas cuando sea útil
- Mantén lenguaje claro y conciso

## Preguntas Frecuentes

### ¿Cómo inicio con mi primera contribución?
Busca issues etiquetados como `good first issue` o `help wanted`.

### ¿Cuánto tiempo toma revisar un PR?
Normalmente 2-5 días. PRs más grandes pueden tomar más tiempo.

### ¿Puedo trabajar en múltiples features simultáneamente?
Sí, pero usa branches separadas para cada feature.

### Mi PR fue rechazado, ¿qué hago?
Revisa los comentarios, haz los cambios solicitados, y actualiza el PR.

## Contacto

- Issues: [GitHub Issues](https://github.com/ferchd/nexa/issues)
- Discussions: [GitHub Discussions](https://github.com/ferchd/nexa/discussions)
- Email: ferchd@outlook.com

---

¡Gracias por contribuir a Nexa! 🚀