# Sprint 4

## Work Completed in Sprint 4

### New Functionality Implemented

- Completed and stabilized all remaining Sprint 3 issues
- Finalized admin analytics dashboard (revenue trend, order trend, top products, low-stock alerts)
- Finalized user analytics dashboard (summary metrics, recent orders, top products, spending trend, recommendations)
- Polished cart and checkout flows including edge cases (empty cart, stock conflicts, failed submissions)
- Improved order history view with status labels and item detail expansion
- Added product filtering by category, price range, and search term with paginated results
- Implemented role-based access control throughout (admin vs. user views, protected routes)
- Added featured products endpoint for the public landing page
- Completed frontend README with setup and run instructions

### Bug Fixes and Improvements

- Fixed order status field missing from history fixtures (required by Cypress and OrderHistory component)
- Corrected JWT cookie expiry behavior on logout
- Resolved stock conflict handling in cart update and order creation flows
- Stabilized Cypress specs against live UI label changes (`Gator Inventory`, `Open Products`)
- Improved error propagation from repository layer through service and controller layers

---

## Frontend Unit Tests (Vitest / React Testing Library)

### Test Files

| File | What It Covers |
|---|---|
| `Login.test.tsx` | Form render, valid/invalid credentials, network failure, demo credentials block, signup navigation |
| `Add.test.tsx` | Product add form validation, submit behavior, error states |
| `Update.test.tsx` | Product update form pre-population, submit behavior, validation |
| `ProductTable.test.tsx` | Product list render, empty state, filter interaction |
| `FilterPanel.test.tsx` | Filter inputs update UI, category and price range behavior |
| `Pagination.test.tsx` | Page controls render, page size changes, navigation |
| `Dashboard.test.tsx` | Navigation links by role, admin-only action visibility |
| `Cart.test.tsx` | Add item, update quantity, remove item, empty cart state |
| `Checkout.test.tsx` | Successful checkout, empty state, failed submission, order detail render |
| `OrderHistory.test.tsx` | Order cards render, status labels, detail panel expansion |
| `AdminAnalytics.test.tsx` | Summary metrics, chart/table render states |
| `UserAnalytics.test.tsx` | User summary, spending trend, top products, recommendations |

### Run Command

```bash
cd client
npm run test:run
```

---

## Frontend Cypress End-to-End Tests

### Spec Files and Results

| Spec | Tests Passing |
|---|---|
| `app.cy.js` | 17 |
| `cart.cy.js` | 6 |
| `checkout.cy.js` | 5 |
| `history.cy.js` | 5 |
| `login.cy.js` | 5 |
| **Total** | **38 passing, 0 failing** |

### Coverage by Spec

**`app.cy.js`** — Root page load, auth flow, navigation links, product dashboard render, product CRUD (create/update/delete), search and category filters, role-based access control, pagination controls, page-size changes

**`cart.cy.js`** — Add product from dashboard, change quantity, empty cart state, remove items, navigate to checkout, delete last item with minus button

**`checkout.cy.js`** — Successful checkout, empty checkout state, failed order submission, order detail rendering

**`history.cy.js`** — Order history render, empty state, error state, computed totals, item detail expansion

**`login.cy.js`** — Login success, invalid credentials, network failure, demo credential text visibility, navigation to signup

### Run Command

```bash
cd client
npx cypress run --browser electron
```

---

## Backend Unit Tests (Go)

### Test Files

| File | What It Covers |
|---|---|
| `pkg/service/product_service_test.go` | Create product (success + validation), paginated listing, empty results, repo error propagation |
| `pkg/service/user_service_test.go` | Login success, missing user, invalid password, role normalization, update validation, repo failures |
| `pkg/service/user_service_signup_test.go` | Signup with password hashing, missing username/password, username uniqueness |
| `pkg/service/cart_service_test.go` | Get cart, add, update, remove, clear, invalid input, missing product, insufficient stock |
| `pkg/service/order_service_test.go` | Create order from cart, get all orders, get by ID, empty cart, missing product, insufficient stock, transaction rollback |
| `pkg/service/analytics_service_test.go` | Summary metrics, revenue/order trends, top products, low-stock queries, invalid query handling |
| `pkg/controller/product_controller_test.go` | HTTP binding, auth enforcement, status codes, service error mapping |
| `pkg/controller/product_controller_crud_test.go` | Create, update, stock patch, delete handler behavior |
| `pkg/controller/user_controller_test.go` | User list, role update, auth/role restrictions |
| `pkg/controller/order_controller_test.go` | Create order, list orders, get by ID, status codes |
| `pkg/controller/analytics_controller_test.go` | Summary, revenue trend, order trend, top products, low-stock endpoints |
| `pkg/controller/analytics_controller_extended_test.go` | Extended analytics edge cases and query parameter validation |
| `pkg/middleware/auth_test.go` | Missing token, invalid token, expired token, valid token context, authorization checks |

### Run Command

```bash
cd server
go test -v ./...
```

### Result

All backend unit test suites are passing.

---

## Backend API Documentation

### Authentication Model

- JWT stored in the `token` cookie
- `AuthMiddleware` validates the JWT and exposes `user_id` and `role` in request context
- `Authorize(...)` restricts access by role
- Frontend requests go through the `/api` proxy prefix; backend route definitions do not include `/api`

### Common Response Formats

**Error response:**
```json
{ "error_message": "string" }
```

**Login response:**
```json
{ "role": "admin" }
```

**Product list response:**
```json
{
  "data": [
    { "id": 1, "name": "Laptop", "price": 1500, "quantity": 10, "category": "Electronics", "image_url": "https://..." }
  ],
  "pagination": { "page": 1, "limit": 5, "total": 1, "total_pages": 1 }
}
```

---

### Endpoint Reference

#### Public Endpoints

| Method | Path | Description |
|---|---|---|
| POST | `/login` | Authenticate user, set `token` cookie. Returns `{ "role": "admin" }`. **200** success, **400** bad JSON, **422** invalid credentials. |
| POST | `/signup` | Create new account. Body: `{ username, email, password, role }`. **201** created, **400**/**422** errors. |
| POST | `/logout` | Clear `token` cookie. **200** always. |
| GET | `/products/featured` | Featured products for landing page. Query: `limit` (default 6, max 12). Returns product array. **200**. |

---

#### Authenticated Endpoints (admin or user role)

| Method | Path | Description |
|---|---|---|
| GET | `/products` | Paginated product list. Query: `page`, `limit`, `search`, `category`, `min_price`, `max_price`. **200**, **401**, **403**. |
| GET | `/orders` | Orders for current user; admins see all. Query: `search`, `status`, `date_from`, `date_to`. **200**, **400**, **401**, **403**, **500**. |
| GET | `/orders/:id` | Single order by ID with ownership/role check. **200**, **400**, **401**, **403**, **404**, **500**. |

---

#### Authenticated Endpoints (user role only)

| Method | Path | Description |
|---|---|---|
| GET | `/cart` | Current user's cart items. Returns cart item array. **200**, **400**, **401**, **403**, **500**. |
| POST | `/cart/add` | Add product to cart. Body: `{ product_id, quantity }`. **201**, **400**, **401**, **403**, **404**, **409**, **500**. |
| PATCH | `/cart/update` | Update cart item quantity. Body: `{ product_id, quantity }`. **200**, **400**, **401**, **403**, **404**, **409**, **500**. |
| DELETE | `/cart/remove` | Remove product from cart. Body: `{ product_id }`. **200**, **400**, **401**, **403**, **404**, **500**. |
| POST | `/orders` | Create order from submitted items. Body: `{ items: [{ product_id, quantity }] }`. Returns `{ order_id }`. **201**, **400**, **401**, **403**, **404**, **409**, **500**. |
| GET | `/user/analytics/summary` | User totals: orders, spend, pending. **200**, **400**, **401**, **403**, **500**. |
| GET | `/user/analytics/recent-orders` | User's recent orders. **200**, **400**, **401**, **403**, **500**. |
| GET | `/user/analytics/top-products` | Most-purchased products for user. **200**, **400**, **401**, **403**, **500**. |
| GET | `/user/analytics/spending-trend` | Spending trend points for user. **200**, **400**, **401**, **403**, **500**. |
| GET | `/user/analytics/recommendations` | Personalized product recommendations. **200**, **400**, **401**, **403**, **500**. |

---

#### Admin Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/users` | Full user list. Returns array of `{ id, username, email, role }`. **200**, **401**, **403**, **500**. |
| PUT | `/users/:id/role` | Update user role. Body: `{ role }`. **200**, **400**, **401**, **403**, **422**. |
| POST | `/products` | Create product. Body: `{ name, price, quantity, category, image_url }`. **201**, **400**, **401**, **403**, **422**. |
| PUT | `/products/:id` | Update product. Same body as create. **200**, **400**, **401**, **403**, **404**, **422**. |
| PATCH | `/products/:id/stock` | Update stock quantity. Body: `{ quantity }`. Returns updated product object. **200**, **400**, **401**, **403**, **404**, **500**. |
| DELETE | `/products/:id` | Delete product. **200**, **400**, **401**, **403**, **404**, **500**. |
| GET | `/analytics/summary` | Admin summary metrics. Query: `days`. Returns `{ totalRevenue, totalOrders, totalProducts, lowStockCount }`. **200**, **400**, **401**, **403**, **500**. |
| GET | `/analytics/revenue-trend` | Revenue trend points. Query: `days`. Returns `[{ date, value }]`. **200**, **400**, **401**, **403**, **500**. |
| GET | `/analytics/order-trend` | Order count trend points. Query: `days`. Returns `[{ date, value }]`. **200**, **400**, **401**, **403**, **500**. |
| GET | `/analytics/top-products` | Top-selling products. Query: `days`, `top_limit`. Returns `[{ productId, name, totalSold, revenue }]`. **200**, **400**, **401**, **403**, **500**. |
| GET | `/analytics/low-stock` | Low-stock products. Query: `threshold`, `low_stock_limit`. Returns `[{ productId, name, quantity, category }]`. **200**, **400**, **401**, **403**, **500**. |

---

### Notes

- `POST /register` is not implemented — use `POST /signup`
- `PATCH /products/:id/stock` returns the full updated product object
- `POST /logout` expires the `token` cookie regardless of whether it exists
- Order history fixtures and the `OrderHistory` component both require an order `status` field
- Cart endpoints are restricted to `user` role; admin users manage inventory through product endpoints
