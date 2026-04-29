# GatorInventory
## Problem Statement

Managing inventory manually or across disconnected tools leads to inefficiencies, data errors, and poor visibility into stock levels. Small teams and businesses often struggle to track what items they have, update quantities in real time, and quickly locate specific products - resulting in overstocking, stockouts, and wasted time. GatorInventory solves this by providing a centralized, easy-to-use platform where users can add, update, delete, and search inventory items through a clean web interface, backed by a secure and reliable REST API.

## Team Members

- Sushrith Bogi (bogisushrith) - Backend
- Sai Sri Krishna Teja Sanku (krishnatejasai) - Frontend
- Dheeraj Kodimela (Dheeraj2125) - Backend
- Abhitej Kodakandla (Abhitej23) - Frontend

## What the app does

This full-stack inventory management application built with React, TypeScript, Vite, Go, Echo, and PostgreSQL.

The application supports two roles:

- Admin users can manage products, users, and sales analytics.
- Regular users can browse products, manage a cart, place orders, and review order history and personal analytics.

The current frontend entry point uses the branding "Gator Inventory" and the primary user flows start from the home page with the Open Products, View Cart, and See all products actions.

## Requirements

Before running the project, install:

- Node.js 18 or newer
- Go 1.23 or newer
- PostgreSQL 12 or newer
- A browser for Cypress runs if you want to execute the end-to-end tests locally

## Environment variables

The backend reads these environment variables from `.env`:

- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `JWT_KEY`
- `PORT` is optional and defaults to `8080`

## Setup

### Backend

```bash
cd server
go mod download
```

Make sure PostgreSQL is running and the backend database credentials are set in `.env`, then start the API:

```bash
go run ./cmd/imsapi
```

The backend listens on `http://localhost:8080` by default.

### Frontend

```bash
cd client
npm install
npm run dev
```

The frontend runs on `http://localhost:5173`.

## Using the application

1. Open the frontend in your browser.
2. Sign up or log in with a valid account.
3. Browse the products page.
4. Add items to the cart.
5. Go to checkout and place an order.
6. Review order history and analytics from the dashboard.

Admin users can also create, edit, delete, and restock products, manage users, and review sales analytics.

## Test commands

From the repository root:

```bash
cd server
go test -v ./...
```

From the frontend folder:

```bash
cd client
npm run test:run
npx cypress run --browser electron
```
