# GatorInventory - Sprint 1
# Github Repository Link - https://github.com/bogisushrith/GatorInventory
## Objective
Our first sprint aims to establish the foundational features of the GatorInventory web-based inventory management system. The focus is on setting up secure user authentication, building core inventory CRUD APIs, and developing the essential frontend components that will serve as the base for all future sprints.

---

## User Stories

### 1. User Registration
*Objective:* Allow new users to create an account with their name, email, and password so they can securely access the inventory management system.

*Details:*
•⁠  ⁠*Frontend:*
  - Design and develop the sign-up component with input fields for name, email, and password.
  - Implement client-side validation and error feedback for invalid inputs.
•⁠  ⁠*Backend:*
  - Implement UserRequest DTO to validate and handle incoming registration payloads.
  - Develop User Repository Layer to store new user records in the database.
  - Implement User Service Logic to handle registration business rules.

*Assigned to:* Abhitej Kodakandla (Frontend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 2. User Login
*Objective:* Allow registered users to log in with their email and password so they can access their personalized inventory dashboard.

*Details:*
•⁠  ⁠*Frontend:*
  - Design and develop the login page component with form validation.
  - Display appropriate error messages for invalid credentials.
•⁠  ⁠*Backend:*
  - Implement User Controller endpoints for the login API route.
  - Add LoginResponse DTO to return structured authentication results.
  - Implement User Service Logic to verify credentials against stored records.

*Assigned to:* Abhitej Kodakandla (Frontend), Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 3. Secure Password Storage
*Objective:* Ensure user passwords are stored securely using hashing so that credentials remain protected even if the database is compromised.

*Details:*
•⁠  ⁠*Backend:*
  - Implement bcrypt password hashing in the User Service Logic during registration.
  - Ensure plain-text passwords are never stored or logged at any point.
  - Validate hashed passwords during the login process.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 4. JWT Authentication
*Objective:* Issue a JWT token upon successful login so that users can make authenticated API requests without re-entering their credentials on every action.

*Details:*
•⁠  ⁠*Backend:*
  - Generate a signed JWT token upon successful user login.
  - Add LoginResponse DTO to wrap and return the token to the client.
  - Implement token validation middleware to protect secured API routes.

*Assigned to:* Sushrith Bogi (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 5. Protected Routes
*Objective:* Block unauthenticated access to dashboard pages and redirect users to login so that inventory data is not publicly accessible.

*Details:*
•⁠  ⁠*Frontend:*
  - Implement Root App Component with route guards that check for a valid auth token.
  - Redirect unauthenticated users to the login page automatically.
•⁠  ⁠*Backend:*
  - Apply JWT middleware to all protected API endpoints.

*Assigned to:* Sai Sri Krishna Teja Sanku (Frontend), Sushrith Bogi (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 6. View All Inventory Items
*Objective:* Allow logged-in users to view a list of all inventory items so they can get a quick overview of what is currently in stock.

*Details:*
•⁠  ⁠*Backend:*
  - Implement the GET all products endpoint in ProductController.
  - Add ProductResponse DTO and mapping functions to serialize product data.
  - Implement Product Repository Layer to fetch all items from the database.
  - Implement Product Service Logic to handle retrieval business rules.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 7. Add a New Inventory Item
*Objective:* Allow inventory managers to add a new item with a name, description, quantity, and price so that inventory records stay up to date.

*Details:*
•⁠  ⁠*Backend:*
  - Implement the POST product endpoint in ProductController.
  - Add ProductRequest DTO to validate and handle incoming product creation payloads.
  - Implement Product Repository Layer to insert new items into the database.
  - Implement Product Service Logic to enforce creation business rules.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 8. View Item Details
*Objective:* Allow users to view the full details of a specific inventory item so they can review or verify the information stored for that product.

*Details:*
•⁠  ⁠*Backend:*
  - Implement the GET product by ID endpoint in ProductController.
  - Add ProductResponse DTO to return full item details in a structured format.
  - Implement Product Repository Layer to fetch a single item by its ID.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* Medium  
*Milestone:* Sprint 1

---

### 9. Edit an Existing Inventory Item
*Objective:* Allow inventory managers to update the details of an existing item so that the inventory always reflects the current state of stock.

*Details:*
•⁠  ⁠*Backend:*
  - Implement the PUT product endpoint in ProductController.
  - Add ProductRequest DTO to validate and handle incoming update payloads.
  - Implement Product Repository Layer to apply updates to existing records.
  - Implement Product Service Logic to enforce update business rules.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 10. Delete an Inventory Item
*Objective:* Allow inventory managers to delete an item from the inventory so that discontinued or irrelevant products can be removed from the system.

*Details:*
•⁠  ⁠*Backend:*
  - Implement the DELETE product endpoint in ProductController.
  - Implement Product Repository Layer to remove the item from the database.
  - Implement Product Service Logic to handle deletion rules and validation.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 11. Inventory Dashboard Overview
*Objective:* Provide logged-in users with a dashboard that displays a summary of their inventory so they can assess inventory health at a glance.

*Details:*
•⁠  ⁠*Frontend:*
  - Design and develop the dashboard component as the main post-login landing page.
  - Include sections for displaying inventory summary information.

*Assigned to:* Abhitej Kodakandla (Frontend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 12. Responsive Login Page
*Objective:* Ensure the login page is responsive and cleanly styled so users have a smooth experience regardless of the device or screen size.

*Details:*
•⁠  ⁠*Frontend:*
  - Apply Tailwind CSS utility classes to the login page for responsive layout.
  - Configure Tailwind CSS content paths and base styles to support responsive design.
  - Add Global App Styles and Base CSS Styles to ensure consistent rendering.

*Assigned to:* Abhitej Kodakandla (Frontend), Sai Sri Krishna Teja Sanku (Frontend)  
*Priority:* Medium  
*Milestone:* Sprint 1

---

### 13. Responsive Signup Page
*Objective:* Ensure the signup page is intuitive and clearly laid out so that registration is straightforward and error-free for new users.

*Details:*
•⁠  ⁠*Frontend:*
  - Apply Tailwind CSS to the sign-up component for a clean and responsive layout.
  - Ensure form fields are clearly labeled with accessible validation messages.

*Assigned to:* Abhitej Kodakandla (Frontend), Sai Sri Krishna Teja Sanku (Frontend)  
*Priority:* Medium  
*Milestone:* Sprint 1

---

### 14. Navigation Between Pages
*Objective:* Enable smooth client-side navigation between login, signup, and dashboard pages without full page reloads so that the application feels fast and modern.

*Details:*
•⁠  ⁠*Frontend:*
  - Create a root page that connects all components through a client-side router.
  - Implement Root App Component with defined routes for each page.
  - Set up the application entry point to render the router-wrapped app.

*Assigned to:* Abhitej Kodakandla (Frontend), Sai Sri Krishna Teja Sanku (Frontend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 15. Form Input Validation Feedback
*Objective:* Display clear validation error messages for incorrect form inputs so that users can correct mistakes before submitting.

*Details:*
•⁠  ⁠*Frontend:*
  - Implement client-side validation on login and sign-up forms.
  - Configure Tailwind CSS to style error states for form fields.
  - Show inline error messages for fields like invalid email format or short passwords.

*Assigned to:* Abhitej Kodakandla (Frontend), Sai Sri Krishna Teja Sanku (Frontend)  
*Priority:* Medium  
*Milestone:* Sprint 1

---

### 16. RESTful API for Products
*Objective:* Provide a RESTful API that supports full CRUD operations for products so the frontend can create, read, update, and delete inventory items via HTTP requests.

*Details:*
•⁠  ⁠*Backend:*
  - Implement ProductController with all four CRUD endpoint handlers.
  - Register all product routes through the route registration logic.
  - Define Service DTO Models to standardize data flow across layers.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 17. RESTful API for Users
*Objective:* Ensure user management endpoints follow REST conventions so the system has a consistent and predictable API structure.

*Details:*
•⁠  ⁠*Backend:*
  - Implement User Controller endpoints for register and login routes.
  - Define Service DTO Models for consistent request and response structures.
  - Add missing route registration and server startup logic to expose all user routes.

*Assigned to:* Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 18. CORS Configuration
*Objective:* Configure CORS on the backend so that API requests from the frontend development server are not blocked by the browser.

*Details:*
•⁠  ⁠*Backend:*
  - Update CORS middleware to allow cross-origin requests from the frontend origin.
  - Add required environment variables to configure allowed origins dynamically.

*Assigned to:* Sushrith Bogi (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 19. Database Schema Design
*Objective:* Design a well-structured relational database schema for users and products so that data is stored consistently and entity relationships are properly enforced.

*Details:*
•⁠  ⁠*Backend:*
  - Create the initial SQL schema with users and products tables.
  - Define primary keys, foreign keys, and required field constraints.
  - Initialize and configure project dependencies including the database driver.

*Assigned to:* Sushrith Bogi (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 20. Project Scaffolding and Build Configuration
*Objective:* Set up proper build tooling and configuration files for both the frontend and backend so the development environment can be set up quickly and consistently by all team members.

*Details:*
•⁠  ⁠*Frontend:*
  - Configure Vite Build Setup with the React plugin and development settings.
  - Initialize TypeScript Configuration with strict type checking.
  - Configure Tailwind CSS with content paths and theme settings.
  - Configure ESLint Rules for consistent code style and error detection.
  - Initialize the Application Entry Point to bootstrap the React application.
•⁠  ⁠*Backend:*
  - Initialize and configure Go project dependencies.
  - Maintain the dependency lock file to keep ⁠ go.sum ⁠ consistent.
  - Start the main function and complete the index file for server bootstrapping.

*Assigned to:* Sai Sri Krishna Teja Sanku (Frontend), Sushrith Bogi (Backend), Dheeraj Kodimela (Backend)  
*Priority:* High  
*Milestone:* Sprint 1

---

### 21. Search and Filter Inventory Items
*Objective:* Allow inventory managers to search for items by name and filter them by category or stock level so they can quickly locate specific products without scrolling through the entire inventory list.

*Details:*
•⁠  ⁠*Frontend:*
  - Design a search bar and filter controls on the inventory list page.
  - Dynamically render filtered results based on user input.
•⁠  ⁠*Backend:*
  - Implement a search endpoint that accepts query parameters such as name and category.
  - Update Product Repository Layer to support filtered database queries.

*Assigned to:* Abhitej Kodakandla (Frontend), Dheeraj Kodimela (Backend)  
*Priority:* Medium  
*Milestone:* Sprint 2 (Not completed in Sprint 1 — category field was not finalized in the database schema during this sprint)


## Issues Planned Based on User Stories

The team planned to address 32 GitHub issues during Sprint 1, split across the backend and frontend.

### Backend (19 Issues)

| Issue # | Title |
|---------|-------|
| #2 | Start the main function |
| #4 | Create error message and declare required variables |
| #7 | Update CORS middleware and add more required variables |
| #9 | Add missing route registration and server startup logic |
| #11 | Complete Index file |
| #21 | Initialize and Configure Project Dependencies |
| #22 | Maintain Dependency Lock File |
| #24 | Initial database schema for Inventory Management System |
| #26 | Add the Product request functionality |
| #28 | Add the UserRequest Functionality |
| #30 | Add ProductResponse DTO and mapping functions |
| #32 | Add ErrorResponse and LoginResponse DTOs |
| #34 | Implement ProductController with CRUD operations |
| #36 | Implement User Controller Endpoints |
| #37 | Implement Product Repository Layer |
| #38 | Implement User Repository Layer |
| #39 | Define Service DTO Models |
| #40 | Implement Product Service Logic |
| #41 | Implement User Service Logic |

### Frontend (13 Issues)

| Issue # | Title |
|---------|-------|
| #13 | Create a dashboard component |
| #14 | Create a login page component |
| #16 | Create a sign-up component |
| #19 | Create a root page to connect all these components |
| #48 | Add Global App Styles |
| #49 | Implement Root App Component |
| #50 | Add Base CSS Styles |
| #51 | Initialize Application Entry Point |
| #52 | Configure Vite Build Setup |
| #53 | Initialize TypeScript Configuration |
| #54 | Configure Tailwind CSS |
| #55 | Configure ESLint Rules |
| #56 | Implement Search and Filter for Inventory Items (US-21) |

---

## Issues Completed

31 out of 32 planned issues were successfully closed during Sprint 1.

### Backend — 19/19 Completed ✅

| Issue # | Title | Status |
|---------|-------|--------|
| #2 | Start the main function | ✅ Closed |
| #4 | Create error message and declare required variables | ✅ Closed |
| #7 | Update CORS middleware and add more required variables | ✅ Closed |
| #9 | Add missing route registration and server startup logic | ✅ Closed |
| #11 | Complete Index file | ✅ Closed |
| #21 | Initialize and Configure Project Dependencies | ✅ Closed |
| #22 | Maintain Dependency Lock File | ✅ Closed |
| #24 | Initial database schema for Inventory Management System | ✅ Closed |
| #26 | Add the Product request functionality | ✅ Closed |
| #28 | Add the UserRequest Functionality | ✅ Closed |
| #30 | Add ProductResponse DTO and mapping functions | ✅ Closed |
| #32 | Add ErrorResponse and LoginResponse DTOs | ✅ Closed |
| #34 | Implement ProductController with CRUD operations | ✅ Closed |
| #36 | Implement User Controller Endpoints | ✅ Closed |
| #37 | Implement Product Repository Layer | ✅ Closed |
| #38 | Implement User Repository Layer | ✅ Closed |
| #39 | Define Service DTO Models | ✅ Closed |
| #40 | Implement Product Service Logic | ✅ Closed |
| #41 | Implement User Service Logic | ✅ Closed |

### Frontend — 12/13 Completed ✅

| Issue # | Title | Status |
|---------|-------|--------|
| #13 | Create a dashboard component | ✅ Closed |
| #14 | Create a login page component | ✅ Closed |
| #16 | Create a sign-up component | ✅ Closed |
| #19 | Create a root page to connect all these components | ✅ Closed |
| #48 | Add Global App Styles | ✅ Closed |
| #49 | Implement Root App Component | ✅ Closed |
| #50 | Add Base CSS Styles | ✅ Closed |
| #51 | Initialize Application Entry Point | ✅ Closed |
| #52 | Configure Vite Build Setup | ✅ Closed |
| #53 | Initialize TypeScript Configuration | ✅ Closed |
| #54 | Configure Tailwind CSS | ✅ Closed |
| #55 | Configure ESLint Rules | ✅ Closed |
| #56 | Implement Search and Filter for Inventory Items | ❌ Not Completed |

**Total: 31/32 issues completed. 1 issue incomplete.**

---

## Incomplete Issues

| Issue # | Title | Reason |
|---------|-------|--------|
| #56 | Implement Search and Filter for Inventory Items | ❌ Not Completed — See below |

**Issue #56 — Implement Search and Filter for Inventory Items**

This feature was not completed because implementing search and filter requires the product database schema to include a category field with consistent, predefined values, which was not finalized during Sprint 1. Without a stable category structure in place, building a filter that produces meaningful and accurate results was not possible within the sprint timeline.
Copy this directly into a new file named Sprint1.md and push it to the root of your GitHub repo.Write me prob
