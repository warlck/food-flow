# Frontend API Integration - Summary

## What Was Done

Successfully integrated the frontend application with the backend REST API to fetch restaurant details dynamically.

## Changes Made

### 1. API Service Layer (`src/lib/api.ts`)
- Created `RestaurantApiService` class to handle API communication
- Defined TypeScript interfaces matching the backend API response:
  - `ApiRestaurantDetails`
  - `ApiCategory`
  - `ApiMenuItem`
- Implemented `getRestaurantDetails(restaurantId)` method
- Configured API base URL from environment variable `VITE_API_URL`

### 2. Data Transformers (`src/lib/transformers.ts`)
- Created transformation functions to convert API response to frontend types:
  - `transformApiRestaurant()` - Converts API restaurant to frontend `Restaurant` type
  - `transformApiMenuItems()` - Converts API categories/items to frontend `MenuItem[]`
- Handles data structure differences between backend and frontend

### 3. React Query Hook (`src/hooks/useRestaurantDetails.ts`)
- Created custom hook `useRestaurantDetails(restaurantId)`
- Implements automatic caching (5-minute stale time)
- Provides loading and error states
- Configures retry logic (2 retries on failure)

### 4. Environment Configuration (`src/vite-env.d.ts`)
- Added TypeScript interface for `VITE_API_URL` environment variable
- Ensures type safety for environment variables

### 5. Routing (`src/App.tsx`)
- Added new route: `/menu/:restaurantId`
- Supports dynamic restaurant ID in URL path

### 6. Menu Page (`src/pages/Menu.tsx`)
- Integrated API data fetching with `useRestaurantDetails` hook
- Added three ways to pass `restaurant_id`:
  1. URL path parameter: `/menu/{restaurantId}`
  2. Query parameter: `/menu?restaurant_id={restaurantId}`
  3. Default value: `a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d`
- Implemented loading state with spinner
- Implemented error state with alert and retry button
- Falls back to mock data if API fails (graceful degradation)

### 7. Build Configuration
- Updated `Makefile` to pass `VITE_API_URL` build argument
- Configured Docker build to use `http://localhost:3000` as API URL

## How to Use

### Access the Menu Page

**With Donergy Restaurant (default):**
```bash
open http://localhost:8080/menu/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d
```

**With Query Parameter:**
```bash
open http://localhost:8080/menu?restaurant_id=a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d
```

**Without Parameter (uses default):**
```bash
open http://localhost:8080/menu
```

### API Endpoint Called

```
GET http://localhost:3000/v1/restaurants/{restaurant_id}/details
```

### Expected Response

```json
{
  "id": "a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d",
  "name": "Donergy",
  "description": "Authentic Turkish Kebab & Pide Restaurant",
  "address": "9 Raffles Boulevard, #01-91B, Millenia Walk, Singapore 039596",
  "phone": "+65 6333 0785",
  "email": "info@donergy.sg",
  "imageUrl": "https://www.donergy.sg/Content/images/logo.png",
  "enabled": true,
  "categories": [
    {
      "id": "c1000000-0000-0000-0000-000000000001",
      "name": "Kebab Roll",
      "description": "Delicious kebab wrapped in fresh flatbread",
      "enabled": true,
      "mentuItems": [
        {
          "id": "a1b2c3d4-0001-0000-0000-000000000001",
          "name": "Chicken Kebab Roll",
          "description": "Tender chicken kebab wrapped in fresh flatbread",
          "price": 11,
          "imageUrl": "https://...",
          "available": true
        }
      ]
    }
  ]
}
```

## Deployment

The updated frontend has been:
1. ✅ Built with API integration
2. ✅ Loaded into Kind cluster
3. ✅ Deployed and running

### Verify Deployment

```bash
# Check frontend pod status
kubectl get pods -n sales-system -l app=frontend

# Check frontend logs
kubectl logs -n sales-system -l app=frontend --tail=50

# Test the endpoint
curl http://localhost:8080/health
```

## Technical Details

### State Management
- Uses React Query (`@tanstack/react-query`) for server state management
- Automatic background refetching
- Optimistic UI updates
- Built-in caching strategy

### Error Handling
- Network errors are caught and displayed to user
- Retry mechanism on failure
- Graceful fallback to mock data
- User-friendly error messages

### Loading States
- Shows loading spinner while fetching
- Prevents layout shifts
- Smooth transitions

### Type Safety
- Full TypeScript support
- API response types match backend
- Compile-time type checking

## Next Steps (Optional Enhancements)

### 1. Service Discovery in Kubernetes
Update API URL for internal cluster communication:
```makefile
--build-arg VITE_API_URL=http://sales-service.sales-system.svc.cluster.local:3000
```

### 2. Environment-Specific Builds
```makefile
frontend-dev:
	--build-arg VITE_API_URL=http://localhost:3000

frontend-prod:
	--build-arg VITE_API_URL=https://api.production.com
```

### 3. Add More API Endpoints
- Create order
- User authentication
- Search functionality
- Favorites

### 4. Error Boundary
Add React Error Boundary for better error handling

### 5. Loading Skeleton
Replace spinner with skeleton loading for better UX

### 6. Real-time Updates
Implement WebSocket for real-time menu updates

## Testing

### Manual Testing
1. Open browser DevTools (Network tab)
2. Navigate to: `http://localhost:8080/menu/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d`
3. Verify API call is made to backend
4. Check response data matches menu displayed
5. Test error state (stop backend API)
6. Test loading state (use network throttling)

### Automated Testing (Future)
```typescript
// Example test
describe('Menu Page', () => {
  it('should fetch and display restaurant details', async () => {
    // Test implementation
  });
});
```

## Files Created/Modified

### New Files
- ✅ `src/lib/api.ts`
- ✅ `src/lib/transformers.ts`
- ✅ `src/hooks/useRestaurantDetails.ts`
- ✅ `API_INTEGRATION.md`
- ✅ `INTEGRATION_SUMMARY.md`

### Modified Files
- ✅ `src/vite-env.d.ts`
- ✅ `src/App.tsx`
- ✅ `src/pages/Menu.tsx`
- ✅ `Makefile`

### Unchanged (Already Configured)
- ✅ `infra/docker/dockerfile.frontend`
- ✅ `infra/docker/nginx.conf`
- ✅ `infra/k8s/base/frontend/base-frontend.yaml`

## Success Criteria

✅ Frontend can accept `restaurant_id` parameter  
✅ Frontend calls `/v1/restaurants/{restaurant_id}/details` API  
✅ Frontend displays data from API response  
✅ Loading state implemented  
✅ Error handling implemented  
✅ Deployed to Kubernetes cluster  
✅ Accessible at `http://localhost:8080`  

## Conclusion

The frontend application has been successfully integrated with the backend API. Users can now:
- View dynamic restaurant data from the database
- See all categories and menu items for Donergy restaurant
- Experience proper loading and error states
- Navigate using restaurant_id in the URL

The application is production-ready for the current use case and can be extended with additional features as needed.
