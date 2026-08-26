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
- **No Agent Attribution**: Never add `Co-authored-by` trailers or any other agent/AI attribution to commit messages.
- **Atomic Commits**: Group edits logically. Do not mix refactoring, feature work, and configuration updates in a single commit.
- **Self-Contained & Compilable Commits**: Every individual commit MUST be self-contained, fully compilable, and pass all relevant automated tests in isolation. Never leave test helpers, struct definitions, or dependent logic missing from the commit that introduces code requiring them.

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

## Adversarial Testing & Quality Assurance Mandates
To eliminate blindspots, denominator dilution, and untested fallback branches, all agents MUST adhere to these testing mandates:
1. **Adversarial Multi-State Test Fixtures**:
   - Never rely exclusively on homogeneous "happy path" seed data where all records share the same state (e.g. all `completed` orders).
   - Every aggregation, query, and analytics test MUST execute against a fixture containing conflicting and mixed entity states on the same grouping key (e.g., 1 completed, 1 cancelled, and 1 in-flight order on the same date bucket) to verify that denominators and filters isolate the exact intended states without dilution.
2. **Branch & Boundary Condition Testing**:
   - Every alternative parsing branch, format fallback (e.g., date-only `YYYY-MM-DD` alongside `RFC3339`), and optional query parameter MUST have a dedicated subtest asserting exact boundary behavior.
   - For date-only filters, tests must explicitly assert that events occurring up to the final microsecond of the day (`23:59:59.999999 UTC`) are included.
3. **Exact Mathematical Value Assertions**:
   - Loose assertions such as `assert value > 0` or `assert len(slice) > 0` are STRICTLY BANNED for financial calculations, rates, percentages, and metrics.
   - Tests MUST assert exact, hand-calculated mathematical values and explicit decimal precision rules (e.g., rounded to 2 decimal places).
4. **Triangular Consistency Protocol**:
   - Every metric or domain calculation MUST be strictly synchronized across all three tiers:
     * **Tier 1 (Store Layer)**: SQL query arithmetic and filtering invariants.
     * **Tier 2 (Application/API Layer)**: DTO transformation, scaling, and precision rounding.
     * **Tier 3 (Frontend/UI Layer)**: Display formatting, badge labels, and user-facing formula legends.

## Code Review Flow
Finding problems and writing the report are two different jobs. Always do the finding first and the writing last. When asked to review a branch (and always before producing a Feature Review Report), follow these steps in order:

1. **Freeze the Diff**: Record the merge-base and head SHA once (`git merge-base master HEAD`, `git rev-parse HEAD`). Every anchor, stat, and claim in the review refers to this exact range.
2. **Findings Pass**: Fan out area-scoped passes over the diff (schema/seeds, store layer, business logic, API + tests, frontends), using parallel subagents when available. Collect *candidates only*: `file:line` plus the concrete trigger that makes it a bug. No severities, no prose, no report sections at this stage.
3. **Validation Gate**: Prove every candidate against the actual code before believing it: open the code, confirm the trigger fires, and check provenance via git history (e.g. `git log -S`) before blaming the branch. Deduplicate mirrored-domain twins (e.g. a menuitem bug and its identical addon twin collapse into one finding). Kill anything you cannot demonstrate. Keep a list of the risks you checked and cleared — it becomes the Verified Non-Issues section for free.
4. **Spec Compliance Pass** (spec-based features): Read the spec and map every requirement and design decision to the commit(s) implementing it. Flag anything unimplemented or implemented differently; intent-vs-implementation gaps are invisible to pure diff review.
5. **Human Triage**: Present the validated findings as a compact ledger (one line per finding, grouped by severity) and wait for fix / skip / defer decisions. Never start fix work before triage — declined findings must cost zero engineering time.
6. **Fix Commits**: One atomic conventional commit per approved finding, with its tests riding along, verified green (build + relevant tests) before committing. Reference the finding number in the commit body so the report's Resolved Findings mapping is trivial.
7. **Final Verification**: Re-run the full gate: build, `go vet`, `gofmt`, affected test suites, and frontend typecheck/lint. The exact commands and outcomes feed the Verification Log.
8. **Render the Report Last**: Generate the report from the validated ledger, never freehand. Split findings into Resolved (finding → fix commit SHA) and Open (accepted, deferred, declined). Compute line anchors once, against the final post-fix tree. `suggestion` blocks are only required for open findings — for resolved ones the fix diff is the suggestion.

### Hard Rules
- **No report prose before the validation gate.** The report is generated from verified evidence, never authored on speculation.
- **No fixes before human triage.** Otherwise the report degenerates into a changelog of decisions the user was never asked about.
- **Review only the branch diff.** Pre-existing issues are labeled as pre-existing and normally do not block the branch, but security-relevant ones must be called out in the Pre-Merge Checklist.
- **Scale depth to the diff.** Small branches may use a lighter findings pass, but the validation gate and human triage are never skipped.

## Feature Review Reports & Walkthrough Documentation
- **Mandatory Review Report Generation**: Whenever asked to create a code review report, review a branch, or finish implementing a feature spec, you MUST generate a CodeRabbit-style review report covering the full branch diff and save it directly in `docs/reviews/<feature-slug>-review.md` (e.g. `docs/reviews/sales-insights-dashboard-review.md`).
- **Review Assets**: Save all review screenshots and visual artifacts in `docs/reviews/images/`.
- **Local Walkthroughs**: Save walkthroughs in `docs/local/walkthrough.md` with visual artifacts in `docs/local/images/`.
- **Gitignored Location**: Both `docs/reviews/` and `docs/local/` are local-only gitignored directories and are never committed.
- **Required Report Sections**:
  1. **Header metadata**: base branch and merge-base SHA, head SHA, commit count, files changed, and diffstat (`git diff --stat $(git merge-base master HEAD)..HEAD`).
  2. **Walkthrough**: a per-area table (schema, store, business, API, frontends, tests, docs) summarizing what changed and why.
  3. **Spec Compliance**: map each spec requirement and design decision to the commit(s) implementing it; explicitly flag any spec item left unimplemented or implemented differently than specified.
  4. **Findings**: actionable comments tagged by severity (🔴 Critical, 🟠 Major, 🔵 Minor, 🟡 Trivial), each anchored to `file:line` in the final branch state and including fenced `suggestion` blocks where a concrete fix exists. Every finding MUST be validated against the actual code before being reported; never report speculative or unverified issues. Distinguish branch-introduced issues from pre-existing ones (check provenance via git history before blaming the branch).
  5. **Resolved Findings**: when review-driven fixes already landed mid-branch, a table mapping each finding to its fix commit for auditability.
  6. **Verified Non-Issues**: notable risks that were checked and found safe, each with the reason it is safe.
  7. **Verification Log**: the exact build, vet, format, lint, typecheck, and test commands run, with their outcomes.
  8. **Pre-Merge Checklist**: remaining gates (declined suggestions, follow-up branches, reseed/deploy steps).

## Frontend UX & UI Design Invariants
- **Lean Error States (No Leaky Parsing)**:
  * When handling failed entity lookups or client error states (e.g., order tracking, missing menus, 404s), avoid overengineered client-side error parsing cascades, string sanitizers, or regex filters.
  * Render clean, human-friendly static domain copy directly (e.g., *"Order Not Found"*, *"We could not find the order you were looking for."*). Never leak raw backend validation JSON arrays (e.g. `[{"field":"...", "error":"..."}]`) or internal server errors to user interfaces.
- **Dark Surface UI Component Theming**:
  * Default library component variants (e.g., Shadcn `variant="outline"`) bind light `bg-background` (`#ffffff`) by default.
  * On dark surfaces (`bg-gray-800`, `bg-ink-950`), never rely on unthemed outline variants that produce solid white boxes. Explicitly style buttons and cards with dark glass tokens (e.g. `border border-white/20 bg-white/10 hover:bg-white/20 text-white rounded-xl` or brand primary styles).
