# Food Flow Agent Rules

These rules apply to all agent interactions within this workspace.

## General Development Rules
- **Commit Granularity**: For each proposed step described in an implementation plan, you MUST create a separate git commit. This ensures commits remain small, atomic, and manageable.
- **Go Standard Practices**: Follow idiomatic Go guidelines. Ensure that `go fmt`, `go vet`, and `go mod tidy` are run when modifying Go code.
- **Microservices & Kubernetes**: This project is built as a set of microservices running on a local Kubernetes (`kind`) cluster. Always consider how changes impact deployment and routing.
- **Makefiles**: Rely on the provided `Makefile` targets (e.g., `make test`, `make dev-run`) for building, testing, and deploying to the cluster.
- **Testing**: Ensure existing tests pass and consider adding tests for new features.

## Git & Commits
- **Branch Creation**: Always create and switch to a separate, appropriately named git branch (e.g., `feat/...`, `fix/...`) for any new task, feature, or bug fix before making code edits or commits. Never commit directly to `master`, and NEVER merge branches into `master` (merging to `master` must only be done by the user).
- **Conventional Commits**: Use the Conventional Commits specification (e.g., `feat:`, `fix:`, `refactor:`, `chore:`, `docs:`) for all commit messages.
- **Atomic Commits**: Group edits logically. Do not mix refactoring, feature work, and configuration updates in a single commit.

## Database & Seeding Guidelines
- **Strict Name Validation**: The business layer strictly validates names for categories, menu items, and addons (e.g., using `^[\p{L}\p{N}' -]{3,100}$`). Do NOT use special characters like double quotes (`"`), parentheses (`()`), or ampersands (`&`) in `name` fields. Place weights, sizing, and details into the `description` fields instead.
- **Cascade Deletes**: Foreign keys are configured with `ON DELETE CASCADE`. To reset database seeds cleanly, delete the top-level parent record (e.g., `DELETE FROM restaurants WHERE restaurant_id = '...'`) instead of running manual deletions across all child tables.

## Network & Port Mappings
- Refer to the following local access configurations in development:
  * **Frontend UI**: `http://localhost:8080` (mapped via Kind NodePort `30080` to Nginx container port `8080`).
  * **Sales API**: `http://localhost:3000` (mapped to container port `3000`).
  * **Auth API**: `http://localhost:6000` (mapped to container port `6000`).
- **Nginx Proxying**: The frontend Nginx container proxies all requests matching `/v1/` to the backend `sales-service:3000` in Kubernetes. Ensure frontend builds configure `VITE_API_URL` as empty (`""`) to utilize relative URLs and routing via this proxy.

## Validation & Exhaustive Testing Guidelines
- **End-to-End Field Mapping Checklist**: When adding or modifying any domain field or entity property, you MUST systematically trace and verify every layer in the data pipeline:
  1. **Database Schema**: Migration SQL (`migrate.sql`) and Seed SQL (`seed.sql`).
  2. **Database Store Layer**: Struct definitions (`*db/model.go`), SQL query strings (SELECT, INSERT, UPDATE in `*db.go`), and converters (`toDB*`, `toBus*`).
  3. **Business Domain Models**: Entity, `New*`, and `Update*` structs in `*bus/model.go`.
  4. **Business Domain Methods**: Explicitly audit both `Create()` and `Update()` functions in `*bus.go` to ensure new fields are assigned when constructing or updating entity structs.
  5. **Application/API DTOs**: Request/response structs (`*app/model.go`), HTTP handlers, and DTO converters (`toBusUpdate*`, `ToApp*`).
  6. **Frontend API & UI Components**: API Service interfaces, transformers, state management, forms, and render conditionals.
- **Exhaustive Unit & Integration Testing**:
  - Never declare a feature or field addition complete without adding unit and integration tests that specifically target and assert the newly added logic.
  - Every field in an `Update` model MUST have an explicit unit test verifying that updating the field via the business layer correctly mutates and persists the value in the database.
  - Include tests for edge cases (e.g. `0` / zero-values, `nil` optionals, boundary limits) as well as non-zero / active values.


