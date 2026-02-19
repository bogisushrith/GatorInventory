# GatorInventory - Sprint 1

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
