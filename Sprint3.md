# Sprint 3 Report — GatorInventory
---
Video Presentation - https://drive.google.com/file/d/1Bs0DRnE4MB7Rw1Fm07y56JwNPmmyPMxQ/view?usp=sharing
---

## Work Completed in Sprint 3

Sprint 3 extended our inventory management system from a basic product-management tool into a fully functional inventory commerce platform. The major additions this sprint were a cart system, a complete order processing pipeline, role-based access control, backend unit testing, and several frontend UX improvements.

### Cart System

We designed and implemented a cart module from scratch, as this feature was not present in Sprint 2. Users can now add products to their cart, adjust quantities using increment/decrement controls, remove individual items, or clear the entire cart. Cart state persists across page navigations and updates in real time.

We added the following API endpoints to support this:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/cart/add` | Add a product to the cart |
| `GET` | `/cart` | Retrieve the current user's cart |
| `PATCH` | `/cart/update` | Update item quantity (quantity `0` removes the item) |
| `DELETE` | `/cart/remove` | Remove a specific item from the cart |

---

### Order and Transaction System

We implemented a full order processing system that ties together the cart and inventory layers. When a user checks out, an order is created from their cart, stock is automatically deducted from the relevant products, and the cart is cleared on success. If the cart is empty, the order can also be placed directly using an items list in the request body.

To support orders, we added two new database tables — `orders` and `order_items` — and the following API endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/orders` | Create a new order |
| `GET` | `/orders` | Fetch orders (role-dependent, see below) |
| `GET` | `/orders/:id` | Fetch a specific order by ID |

---

### Role-Based Access Control

We extended the authorization layer to distinguish between admin and user roles across the order system. Regular users can only view their own orders. Admins can view all orders placed across all users. This access logic is enforced at the service layer and backed by JWT-based role claims set at login.

---

### Inventory Enhancements

We added quantity tracking to the product model and wired stock deduction into the order creation flow so that placing an order automatically reduces available inventory. We also added a dedicated stock update endpoint (`PATCH /products/:id/stock`) to allow admins to manually adjust stock levels without needing to update the full product record.

---

### Backend Unit Testing

We added comprehensive unit tests covering the service layer, which is where all business logic, validation, and authorization rules are enforced. Tests use mock repositories so they run independently of the real database and execute quickly. We organized tests around four service areas — Product, User, Cart, and Order — and covered success paths, edge cases, and failure/rollback scenarios for each.

---

### Frontend Enhancements

We made significant improvements to the frontend user experience this sprint:

- Integrated the cart into the main UI with quantity controls and a running total
- Built a checkout flow that connects to the order creation API
- Added an order history page where users can review past purchases
- Implemented role-based UI rendering so admins and regular users see different options
- Added a profile dropdown (user info + logout) that appears only when logged in
- Redesigned the navigation bar layout for better usability
- Added protected routes so unauthenticated users cannot access internal pages

---

## Frontend Unit Tests

We used **Cypress** for end-to-end tests and **Vitest with React Testing Library** for component-level unit tests.

### Cypress End-to-End Tests

| Test ID | Description |
|---------|-------------|
| APP-001 | Homepage load — verify navigation elements render correctly |
| AUTH-001 | Login flow — enter credentials, submit, verify redirect to dashboard |
| AUTH-002 | Login flow — invalid credentials, verify error message is displayed |
| AUTH-003 | Dashboard verification after successful login |
| DASH-001 | Dashboard page load and element rendering |
| DASH-002 | Product search — type query, apply filter, verify results |
| DASH-003 | Search reset — type query, clear, verify state is reset |
| CRUD-001 | Login and navigate to Add Product page |
| CRUD-002 | Admin dashboard — verify admin-specific controls are visible |
| CRUD-003 | Admin dashboard — verify product table renders correctly |
| FILTER-001 | Category filter interaction and result verification |
| FILTER-002 | Multi-filter application and result verification |
| ACCESS-001 | User role — verify correct UI elements are shown |
| ACCESS-002 | Admin role — verify correct UI elements are shown |
| CART-001 | Add a product to the cart and verify cart updates |
| CART-002 | Update cart item quantity using controls |
| CART-003 | Complete checkout flow from cart to order confirmation |
| ORDER-001 | Order history page — verify orders are listed correctly |

**Results: 17 tests — 17 passing, 0 failing**

---

### Component Unit Tests (Vitest / React Testing Library)

| Component | Tests Covered |
|-----------|---------------|
| `Login.tsx` | Form rendering, input handling, error display on failed auth |
| `Dashboard.tsx` | Product table rendering, category filter interaction |
| `Cart.tsx` | Add item, update quantity, remove item behavior |
| `Checkout.tsx` | Form rendering, order submission, success/error feedback |

---

## Backend Unit Tests

All backend unit tests are written using Go's built-in `testing` package with mock repositories. No real database connection is used in any unit test.

### Product Service (`pkg/service/product_service_test.go`)

| Test | Scenario |
|------|----------|
| `TestAddProduct_Success` | Valid product data creates a product successfully |
| `TestAddProduct_InvalidInput` | Missing or malformed fields return a validation error |
| `TestGetProducts_Success` | Products are returned with correct pagination metadata |
| `TestGetProducts_EmptyList` | Empty product list is handled without error |
| `TestGetProducts_FilterByCategory` | Category filter is applied and returns the correct subset |
| `TestUpdateStock_Success` | Stock quantity is updated correctly by admin |
| `TestUpdateStock_ProductNotFound` | Returns not-found error when product does not exist |
| `TestUpdateStock_InvalidQuantity` | Returns error for negative or zero quantity values |
| `TestDeleteProduct_Success` | Product is removed from the system successfully |
| `TestDeleteProduct_NotFound` | Returns error when product ID does not exist |

---

### User Service (`pkg/service/user_service_test.go`)

| Test | Scenario |
|------|----------|
| `TestLogin_Success` | Valid credentials return the correct role in the response |
| `TestLogin_InvalidCredentials` | Wrong password returns an authentication error |
| `TestLogin_UserNotFound` | Non-existent username returns a not-found error |
| `TestUpdateRole_Success` | Role is updated correctly for a valid user |
| `TestUpdateRole_InvalidRole` | Invalid role string returns a validation error |
| `TestUpdateRole_UserNotFound` | Returns error when the target user does not exist |

---

### Cart Service

| Test | Scenario |
|------|----------|
| `TestAddToCart_Success` | Item is added to cart with correct quantity |
| `TestAddToCart_InsufficientStock` | Returns conflict error when requested quantity exceeds stock |
| `TestAddToCart_ProductNotFound` | Returns not-found error when product does not exist |
| `TestUpdateCartItem_Success` | Cart item quantity is updated correctly |
| `TestUpdateCartItem_SetZeroRemoves` | Setting quantity to 0 deletes the item from the cart |
| `TestRemoveCartItem_Success` | Item is removed from cart successfully |
| `TestRemoveCartItem_NotFound` | Returns error when item is not in the user's cart |

---

### Order Service

| Test | Scenario |
|------|----------|
| `TestCreateOrder_Success` | Order is created from cart, stock deducted, cart cleared |
| `TestCreateOrder_EmptyCart` | Falls back to request body items when cart is empty |
| `TestCreateOrder_InsufficientStock` | Returns conflict error and rolls back the transaction |
| `TestCreateOrder_ProductNotFound` | Returns not-found error and rolls back the transaction |
| `TestGetOrders_User` | Returns only the authenticated user's own orders |
| `TestGetOrders_Admin` | Returns all orders across all users for admin role |
| `TestGetOrderByID_Success` | Returns the correct order details for a valid ID |
| `TestGetOrderByID_NotFound` | Returns not-found error when order ID does not exist |

---

### Controller Tests

We also added controller-level tests to validate API request parsing, response formatting, and HTTP status code correctness across all four modules.

| Controller | Coverage |
|------------|----------|
| Product Controller | Endpoint validation, request binding, response shape |
| User Controller | Auth endpoint validation, request binding, response shape |
| Cart Controller | Cart endpoint validation, request binding, response shape |
| Order Controller | Order endpoint validation, request binding, response shape |

---

## Updated Backend API Documentation

All routes are registered on the Echo server without the `/api` prefix. The frontend uses a Vite proxy to prepend `/api` before forwarding requests to the backend at `http://localhost:8080`.

Authenticated routes require a valid JWT stored in the `token` cookie. The `AuthMiddleware` validates the cookie and injects `user_id` and `role` into the request context. The `Authorize(...)` helper enforces role-based access and returns `403 Forbidden` when the role is not permitted.

**Total APIs documented: 16**

---

### Authentication / User APIs

#### `POST /login`
Authenticates a user and returns their role.

- **Auth required:** No
- **Request body:**
```json
{ "username": "string", "password": "string" }
```
- **Response:**
```json
{ "role": "admin | user" }
```
- **Status codes:** `200 OK` · `400 Bad Request` · `422 Unprocessable Entity`

---

#### `POST /signup`
Registers a new user account.

- **Auth required:** No
- **Request body:**
```json
{ "username": "string", "email": "string", "password": "string", "role": "string" }
```
- **Status codes:** `201 Created` · `400 Bad Request` · `422 Unprocessable Entity`

---

#### `POST /logout`
Clears the `token` cookie and ends the session.

- **Auth required:** No
- **Status codes:** `200 OK`

---

#### `GET /users`
Returns a list of all users. Admin only.

- **Auth required:** Yes — admin
- **Response:**
```json
[{ "id": 1, "username": "string", "email": "string", "role": "string" }]
```
- **Status codes:** `200 OK` · `401 Unauthorized` · `403 Forbidden` · `500 Internal Server Error`

---

#### `PUT /users/:id/role`
Updates the role of a user by ID. Admin only.

- **Auth required:** Yes — admin
- **Request body:**
```json
{ "role": "admin | user" }
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `422 Unprocessable Entity`

---

### Product APIs

#### `GET /products`
Returns paginated products with optional filters.

- **Auth required:** Yes — admin or user
- **Query parameters:** `page`, `limit`, `search`, `category`, `min_price`, `max_price`
- **Response:**
```json
{
  "data": [{ "id": 1, "name": "string", "price": 100, "quantity": 10, "category": "string" }],
  "pagination": { "page": 1, "limit": 5, "total": 100, "total_pages": 20 }
}
```
- **Status codes:** `200 OK` · `401 Unauthorized` · `403 Forbidden`

---

#### `POST /products`
Creates a new product. Admin only.

- **Auth required:** Yes — admin
- **Request body:**
```json
{ "name": "string", "price": 100, "quantity": 10, "category": "string" }
```
- **Status codes:** `201 Created` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `422 Unprocessable Entity`

---

#### `PUT /products/:id`
Updates a product by ID. Admin only.

- **Auth required:** Yes — admin
- **Request body:**
```json
{ "name": "string", "price": 100, "quantity": 10, "category": "string" }
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `422 Unprocessable Entity`

---

#### `PATCH /products/:id/stock`
Updates only the stock quantity for a product. Admin only.

- **Auth required:** Yes — admin
- **Request body:**
```json
{ "quantity": 10 }
```
- **Response:**
```json
{ "id": 1, "name": "string", "price": 100, "quantity": 10, "category": "string" }
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `500 Internal Server Error`

---

#### `DELETE /products/:id`
Deletes a product by ID. Admin only.

- **Auth required:** Yes — admin
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `500 Internal Server Error`

---

### Cart APIs

#### `GET /cart`
Returns the authenticated user's current cart contents.

- **Auth required:** Yes — user
- **Response:**
```json
[{
  "id": 1,
  "user_id": 10,
  "product_id": 1,
  "quantity": 2,
  "product_name": "string",
  "product_price": 100,
  "product_stock": 10
}]
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `500 Internal Server Error`

---

#### `POST /cart/add`
Adds a product to the cart. If the product already exists in the cart, the quantity is increased.

- **Auth required:** Yes — user
- **Request body:**
```json
{ "product_id": 1, "quantity": 2 }
```
- **Status codes:** `201 Created` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `409 Conflict` · `500 Internal Server Error`

---

#### `PATCH /cart/update`
Updates the quantity of a cart item. Sending a quantity of `0` removes the item.

- **Auth required:** Yes — user
- **Request body:**
```json
{ "product_id": 1, "quantity": 3 }
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `409 Conflict` · `500 Internal Server Error`

---

#### `DELETE /cart/remove`
Removes a specific product from the authenticated user's cart.

- **Auth required:** Yes — user
- **Request body:**
```json
{ "product_id": 1 }
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `500 Internal Server Error`

---

### Order APIs

#### `POST /orders`
Creates an order for the authenticated user. The order is built from the cart if items exist; otherwise, the request body `items` list is used.

- **Auth required:** Yes — user
- **Request body:**
```json
{ "items": [{ "product_id": 1, "quantity": 2 }] }
```
- **Response:**
```json
{ "order_id": 101 }
```
- **Status codes:** `201 Created` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `409 Conflict` · `500 Internal Server Error`

---

#### `GET /orders`
Returns orders for the authenticated user. Admins receive all orders; regular users receive only their own.

- **Auth required:** Yes — admin or user
- **Response:**
```json
[{
  "id": 101,
  "created_at": "2026-04-12T10:00:00Z",
  "items": [{ "id": 1, "order_id": 101, "product_id": 1, "quantity": 2 }]
}]
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `500 Internal Server Error`

---

#### `GET /orders/:id`
Returns a single order by ID for the authenticated user.

- **Auth required:** Yes — admin or user
- **Response:**
```json
{
  "id": 101,
  "created_at": "2026-04-12T10:00:00Z",
  "items": [{ "id": 1, "order_id": 101, "product_id": 1, "quantity": 2 }]
}
```
- **Status codes:** `200 OK` · `400 Bad Request` · `401 Unauthorized` · `403 Forbidden` · `404 Not Found` · `500 Internal Server Error`

---

> **Note:** `POST /register` is not a valid endpoint in this system. The correct registration route is `POST /signup`.
