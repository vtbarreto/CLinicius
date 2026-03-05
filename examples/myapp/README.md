# myapp — CLinicius example project

This is a intentionally broken Go project used to demonstrate CLinicius in action.

## Violations

| Package | Imports | Rule broken |
|---|---|---|
| `internal/domain` | `internal/infra` | domain cannot depend on infra |
| `internal/domain` | `internal/repository` | domain cannot depend on repository |
| `internal/handler` | `internal/repository` | handler cannot depend on repository |

## Running

From this directory, with CLinicius installed:

```bash
clinicius check ./...
```

Expected output:

```
❌ Architectural Violation
  Layer:   domain
  Package: myapp/internal/domain
  Imports: myapp/internal/infra
  Rule:    layer-boundary
  Detail:  domain layer cannot depend on internal/infra

❌ Architectural Violation
  Layer:   domain
  Package: myapp/internal/domain
  Imports: myapp/internal/repository
  Rule:    layer-boundary
  Detail:  domain layer cannot depend on internal/repository

❌ Architectural Violation
  Layer:   handler
  Package: myapp/internal/handler
  Imports: myapp/internal/repository
  Rule:    layer-boundary
  Detail:  handler layer cannot depend on internal/repository

3 violation(s) found.
```
