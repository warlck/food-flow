# Food Flow Agent Rules

These rules apply to all agent interactions within this workspace.

## General Development Rules
- **Commit Granularity**: For each proposed step described in an implementation plan, you MUST create a separate git commit. This ensures commits remain small, atomic, and manageable.
- **Go Standard Practices**: Follow idiomatic Go guidelines. Ensure that `go fmt`, `go vet`, and `go mod tidy` are run when modifying Go code.
- **Microservices & Kubernetes**: This project is built as a set of microservices running on a local Kubernetes (`kind`) cluster. Always consider how changes impact deployment and routing.
- **Makefiles**: Rely on the provided `Makefile` targets (e.g., `make test`, `make dev-run`) for building, testing, and deploying to the cluster.
- **Testing**: Ensure existing tests pass and consider adding tests for new features.

## Git & Commits
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

