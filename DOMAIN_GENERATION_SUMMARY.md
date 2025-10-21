# Food Flow Domain Generation Summary - UPDATED

This document describes the domain models generated for the Restaurant, Category, and MenuItem domains.

## ✅ COMPLETED - Generated Files

### Business Layer (business/domain/) - ALL COMPLETE

#### Restaurant Domain (`restaurantbus/`) ✅
- ✅ `model.go` - Restaurant, NewRestaurant, UpdateRestaurant types
- ✅ `restaurantbus.go` - Business logic with CRUD operations
- ✅ `filter.go` - QueryFilter for restaurant queries
- ✅ `order.go` - Ordering constants and defaults

**Restaurant Model Fields:**
- ID (UUID)
- Name (name.Name)
- Description (name.Null)
- Address (string)
- Phone (string)
- Email (string)
- ImageURL (name.Null)
- Enabled (bool)
- DateCreated, DateUpdated (time.Time)

#### Category Domain (`categorybus/`)
- ✅ `model.go` - Category, NewCategory, UpdateCategory types
- ✅ `categorybus.go` - Business logic with CRUD operations
- ✅ `filter.go` - QueryFilter for category queries
- ✅ `order.go` - Ordering constants and defaults

**Category Model Fields:**
- ID (UUID)
- Name (name.Name)
- Description (name.Null)
- RestaurantID (UUID) - Foreign key to Restaurant
- Enabled (bool)
- DateCreated, DateUpdated (time.Time)

#### MenuItem Domain (`menuitembus/`)
- ✅ `model.go` - MenuItem, NewMenuItem, UpdateMenuItem types
- ✅ `menuitembus.go` - Business logic with CRUD operations
- ✅ `filter.go` - QueryFilter for menu item queries
- ✅ `order.go` - Ordering constants and defaults

**MenuItem Model Fields:**
- ID (UUID)
- Name (name.Name)
- Description (name.Null)
- Price (money.Money)
- CategoryID (UUID) - Foreign key to Category
- RestaurantID (UUID) - Foreign key to Restaurant
- ImageURL (name.Null)
- Available (bool)
- DateCreated, DateUpdated (time.Time)

### Application Layer (app/domain/) - MOSTLY COMPLETE

#### Restaurant App (`restaurantapp/`) ✅ COMPLETE
- ✅ `model.go` - API models with JSON tags, conversion functions
- ✅ `restaurantapp.go` - HTTP handlers (create, query, queryByID, update, delete)
- ✅ `route.go` - Route configuration with authentication/authorization
- ✅ `filter.go` - Query parameter parsing
- ✅ `order.go` - Order parameter parsing

**API Endpoints:**
- GET    `/v1/restaurants` - List restaurants (authenticated)
- GET    `/v1/restaurants/{restaurant_id}` - Get by ID (authenticated)
- POST   `/v1/restaurants` - Create (admin only)
- PUT    `/v1/restaurants/{restaurant_id}` - Update (admin only)
- DELETE `/v1/restaurants/{restaurant_id}` - Delete (admin only)

#### Category App (`categoryapp/`) ✅ COMPLETE
- ✅ `model.go` - API models with JSON tags, conversion functions
- ✅ `categoryapp.go` - HTTP handlers (create, query, queryByID, update, delete)
- ✅ `route.go` - Route configuration with authentication/authorization
- ✅ `filter.go` - Query parameter parsing (supports restaurant_id filter)
- ✅ `order.go` - Order parameter parsing

**API Endpoints:**
- GET    `/v1/categories` - List categories (authenticated, can filter by restaurant_id)
- GET    `/v1/categories/{category_id}` - Get by ID (authenticated)
- POST   `/v1/categories` - Create (admin only)
- PUT    `/v1/categories/{category_id}` - Update (admin only)
- DELETE `/v1/categories/{category_id}` - Delete (admin only)

#### MenuItem App (`menuitemapp/`) ⚠️ NEEDS CLEANUP
**Issue:** The folder contains old userapi files that need to be removed/replaced

**Files Needed:**
- `model.go` - API models with JSON tags (needs replacement)
- `menuitemapp.go` - HTTP handlers (needs creation)
- `route.go` - Route configuration (needs replacement)
- `filter.go` - Query parameter parsing (needs replacement)
- `order.go` - Order parameter parsing (needs replacement)

**Remove these incorrect files first:**
- `userapp.go` (wrong file in this directory)
- Any files with `package userapi` declaration

**Then create proper MenuItem app files** (see menuitem_app_template.md below)

## Database Layer (Not Yet Generated)

For full implementation, you'll also need:

### business/domain/*/stores/*db/
- Database store implementation
- SQL queries
- Model conversion (DB ↔ Business)

Example structure (for each domain):
```
business/domain/restaurantbus/stores/restaurantdb/
  - restaurantdb.go (implements Storer interface)
  - model.go (DB model and conversions)
  - filter.go (SQL filter building)
  - order.go (SQL order building)
```

## Next Steps

1. **Complete App Layer** - Create remaining handler and route files for all three domains
2. **Create Database Migrations** - SQL schema for restaurants, categories, menu_items tables
3. **Implement Database Stores** - Implement the Storer interface for each domain
4. **Wire Up Routes** - Add routes to main application
5. **Write Tests** - Create unit and integration tests

## Domain Relationships

```
Restaurant (1) ──< (many) Category (1) ──< (many) MenuItem
```

- A Restaurant has many Categories
- A Category belongs to one Restaurant
- A Category has many MenuItems  
- A MenuItem belongs to one Category and one Restaurant

## API Endpoints (Proposed)

### Restaurant
- GET    /v1/restaurants - List all restaurants
- GET    /v1/restaurants/{id} - Get restaurant by ID
- POST   /v1/restaurants - Create restaurant (Admin only)
- PUT    /v1/restaurants/{id} - Update restaurant (Admin only)
- DELETE /v1/restaurants/{id} - Delete restaurant (Admin only)

### Category
- GET    /v1/categories - List categories (can filter by restaurant_id)
- GET    /v1/categories/{id} - Get category by ID
- POST   /v1/categories - Create category (Admin only)
- PUT    /v1/categories/{id} - Update category (Admin only)
- DELETE /v1/categories/{id} - Delete category (Admin only)

### MenuItem
- GET    /v1/menuitems - List menu items (can filter by category_id, restaurant_id)
- GET    /v1/menuitems/{id} - Get menu item by ID
- POST   /v1/menuitems - Create menu item (Admin only)
- PUT    /v1/menuitems/{id} - Update menu item (Admin only)
- DELETE /v1/menuitems/{id} - Delete menu item (Admin only)
