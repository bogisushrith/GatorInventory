# Sprint 2 Documentation

## Inventory Management System
## Github Repository link: https://github.com/bogisushrith/GatorInventory
---

# Work Completed in Sprint 2

During Sprint 2, we focused on integrating the frontend and backend systems, enhancing core functionality, and ensuring reliability through structured testing.

---

### Full Stack Integration

* Successfully integrated the React frontend with Go backend APIs
* Ensured smooth communication between UI and backend services
* Verified end-to-end functionality across all major features

---

### Inventory Management Features

* Implemented product creation and retrieval
* Added pagination (limit = 5) for efficient data handling
* Built filtering functionality:

  * Search by product name
  * Filter by category
  * Filter by price range

---

### Role-Based Access Control (RBAC)

* Implemented two roles: Admin and User
* Integrated JWT-based authentication
* Enforced authorization using middleware:

  * Admin has full access
  * User has read-only access

---

### User Management System

* Developed admin-only user management functionality
* Admin can:

  * View all users
  * Update user roles dynamically

---

### UI/UX Enhancements

* Improved dashboard layout using Tailwind CSS
* Added:

  * Sidebar navigation
  * Left-aligned filter panel
  * Responsive product table
  * Pagination controls
* Improved layout consistency and usability

---

### Testing Implementation

* Implemented frontend unit tests using Vitest
* Added Cypress for end-to-end testing
* Created backend unit tests using Go’s testing framework

---

# Frontend Testing

Frontend testing was implemented using both unit testing and end-to-end testing to ensure coverage at both component and user interaction levels.

---

## Cypress (End-to-End Testing)

Cypress was used to simulate real user behavior and validate the overall application flow.

### Tests Implemented:

App Load Test

* Verifies that the application loads successfully
* Ensures core UI elements are rendered

Dashboard Navigation Test

* Simulates navigation to the dashboard
* Verifies routing functionality

Products Page Test

* Confirms that the product table is displayed
* Ensures the filter panel is visible

Filter Interaction Test

* Simulates typing into the search input
* Applies filters using the button
* Verifies that user interaction works correctly

These tests validate UI rendering, navigation, and user interaction flows.

---

## Unit Tests (React)

We used Vitest and React Testing Library to test components individually.

### Components Tested:

ProductTable

* Verifies correct rendering of product data
* Ensures props are displayed accurately

FilterPanel

* Tests user input handling
* Verifies state updates when typing

Pagination / UI Behavior

* Tests interaction with pagination controls
* Ensures UI updates correctly based on user actions

These tests ensure component reliability, proper state handling, and correct rendering logic.

---

# Backend Testing

Backend unit testing was implemented using Go’s built-in testing package, focusing on the service layer, which contains the core business logic of the application.

Mock repositories were used to isolate logic and avoid dependency on a real database.

---

## Product Service Tests

Create Product – Success Case

* Valid product input is accepted
* Service maps input correctly and calls the repository

Create Product – Invalid Input

* Tests invalid input such as empty product name
* Ensures validation prevents repository call

Create Product – Repository Failure

* Simulates database failure
* Ensures error is returned properly

Get Products – Success with Pagination

* Verifies product retrieval
* Validates pagination logic including total pages calculation
* Ensures filters are forwarded correctly

Get Products – Empty Dataset

* Handles scenario where no products exist
* Returns empty array and zero counts safely

Get Products – Repository Failure

* Simulates repository error
* Ensures graceful fallback behavior

---

## User / Authentication Service Tests

Login – Success Case

* Valid username and password authentication
* JWT token generation
* Role normalization

Login – User Not Found

* Handles missing user scenario
* Returns appropriate error

Login – Invalid Password

* Simulates incorrect password
* Ensures authentication fails correctly

Update User Role – Success Case

* Updates role with normalization (e.g., ADMIN → admin)
* Ensures repository is called correctly

Update User Role – Invalid Role

* Rejects invalid roles such as "manager"
* Prevents repository call

Update User Role – Repository Failure

* Simulates failure during role update
* Ensures error propagation

---

## Mocking Strategy

* Mock repositories were implemented using structs
* They return predefined responses and simulate errors
* They capture method inputs for verification

This ensures no real database dependency, proper isolation of business logic, and deterministic test behavior.

---

# Backend API Documentation

---

## Authentication APIs

POST /login

Description:
Authenticates a user and returns a JWT token

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "token": "JWT_TOKEN",
  "role": "admin | user"
}
```

---

POST /register 

Description:
Registers a new user

Request Body:

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

Response:

```json
{
  "message": "User registered successfully"
}
```

---

## Product APIs

GET /products

Description:
Fetch products with pagination and filtering

Query Parameters:

* page (integer)
* limit (integer)
* search (string)
* category (string)
* min_price (number)
* max_price (number)

Response:

```json
{
  "data": [
    {
      "id": 1,
      "name": "string",
      "price": 100,
      "quantity": 10,
      "category": "string"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 5,
    "total": 100,
    "total_pages": 20
  }
}
```

---

POST /products

Description:
Create a new product
Access: Admin only

Request Body:

```json
{
  "name": "string",
  "price": 100,
  "quantity": 10,
  "category": "string"
}
```

Response:

```json
{
  "message": "Product created successfully"
}
```

---

PUT /products/:id

Description:
Update product
Access: Admin only

Path Parameter:

* id (product ID)

Request Body:

```json
{
  "name": "string",
  "price": 100,
  "quantity": 10,
  "category": "string"
}
```

Response:

```json
{
  "message": "Product updated successfully"
}
```

---

DELETE /products/:id

Description:
Delete product
Access: Admin only

Path Parameter:

* id (product ID)

Response:

```json
{
  "message": "Product deleted successfully"
}
```

---

## User APIs (Admin Only)

GET /users

Description:
Fetch all users

Headers:

* Authorization: Bearer token

Response:

```json
[
  {
    "id": 1,
    "username": "string",
    "email": "string",
    "role": "user"
  }
]
```

---

PUT /users/:id/role

Description:
Update user role

Path Parameter:

* id (user ID)

Request Body:

```json
{
  "role": "admin"
}
```

Response:

```json
{
  "message": "User role updated successfully"
}
```

---

# Authentication Notes

* JWT token must be included in request headers:
  Authorization: Bearer <token>

* Role-based access is enforced:

  * Admin has full access
  * User has restricted access

---

# Summary

In Sprint 2, we successfully transformed the project into a fully integrated, role-based, and tested full-stack application.

We ensured system reliability through frontend and backend testing while improving usability through thoughtful UI design.

This sprint establishes a strong foundation for future enhancements and scalability.

---
