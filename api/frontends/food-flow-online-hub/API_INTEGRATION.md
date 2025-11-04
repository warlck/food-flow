# Frontend API Integration

The frontend app has been updated to fetch restaurant details from the backend API.

## How It Works

### Restaurant ID Parameter

The frontend accepts the `restaurant_id` parameter in three ways:

1. **URL Path Parameter** (Recommended):
   ```
   http://localhost:8080/menu/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d
   ```

2. **Query Parameter**:
   ```
   http://localhost:8080/menu?restaurant_id=a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d
   ```

3. **Default** (if no parameter provided):
   ```
   http://localhost:8080/menu
   ```
   Uses the default Donergy restaurant ID: `a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d`

### API Endpoint

The frontend calls:
```
GET /v1/restaurants/{restaurant_id}/details
```

This endpoint returns:
- Restaurant information (name, description, address, contact)
- All categories for the restaurant
- All menu items grouped by category

### Data Transformation

The API response is transformed to match the frontend's expected data structure:

**API Response:**
```json
{
  "id": "...",
  "name": "Donergy",
  "categories": [
    {
      "id": "...",
      "name": "Kebab Roll",
      "mentuItems": [
        {
          "id": "...",
          "name": "Chicken Kebab Roll",
          "price": 11.00,
          "imageUrl": "...",
          "available": true
        }
      ]
    }
  ]
}
```

**Frontend Format:**
```typescript
{
  restaurant: Restaurant,
  menuItems: MenuItem[],
  categories: string[]
}
```

## Configuration

### Build-time Configuration

The API URL is configured during the Docker build:

```dockerfile
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}
```

In the Makefile:
```makefile
--build-arg VITE_API_URL=http://localhost:3000
```

### Runtime Configuration

The frontend uses `import.meta.env.VITE_API_URL` to get the API base URL with a fallback:

```typescript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';
```

## Features

### Loading State
While fetching restaurant data, the frontend displays a loading spinner.

### Error Handling
If the API request fails, the frontend shows:
- An error alert with the error message
- A "Try Again" button to reload the page
- Falls back to mock data if configured

### Caching
Uses React Query for:
- Automatic caching (5-minute stale time)
- Retry on failure (2 retries)
- Optimistic UI updates

## Testing

1. **Access the menu with Donergy restaurant**:
   ```bash
   open http://localhost:8080/menu/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d
   ```

2. **Check the network tab** in browser DevTools to see the API call:
   ```
   GET http://localhost:3000/v1/restaurants/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d/details
   ```

3. **Verify the data** is loaded from the API (check restaurant name, categories, menu items)

## Files Modified

### Frontend Application
- `src/lib/api.ts` - API service layer
- `src/lib/transformers.ts` - Data transformation utilities
- `src/hooks/useRestaurantDetails.ts` - React Query hook for fetching data
- `src/vite-env.d.ts` - TypeScript environment variable definitions
- `src/App.tsx` - Added route for `/menu/:restaurantId`
- `src/pages/Menu.tsx` - Updated to use API data with loading/error states

### Configuration
- `Makefile` - Added `VITE_API_URL` build argument
- `infra/docker/dockerfile.frontend` - Already configured for `VITE_API_URL`

## Deployment

To deploy the updated frontend:

```bash
# Build the frontend with API configuration
make frontend

# Load into Kind cluster
make dev-load

# Restart the deployment
kubectl rollout restart deployment frontend -n sales-system

# Wait for ready
kubectl wait pods --namespace=sales-system --selector app=frontend --timeout=120s --for=condition=Ready
```

## Environment Variables

For different environments, update the `VITE_API_URL` in the Makefile:

- **Local Development**: `http://localhost:3000`
- **Kubernetes (Kind)**: `http://sales-service.sales-system.svc.cluster.local:3000`
- **Production**: `https://api.yourdomain.com`

Note: For Kubernetes internal communication, you may need to update the API URL to use the service DNS name.
