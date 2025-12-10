# Furniture Admin Portal

## Overview
This is a furniture store admin portal built with SvelteKit (frontend) and Rust/Actix-web (backend API). The project manages payments, refunds, clients, items, and user administration for a furniture store.

## Project Structure
- **admin-portal/** - SvelteKit frontend (main application)
- **admin-api/** - Rust/Actix-web backend API
- **Furniture.ClientPortal/** - .NET Blazor client application (not currently active)
- **Furniture.ClientPortal-api/** - Go backend API for client portal (not currently active)

## Current State
- Frontend (admin-portal) is running on port 5000
- Backend API (admin-api) expects to run on port 5200 (requires database setup)

## Running the Application
The admin-portal frontend runs automatically via the configured workflow on port 5000.

## Two-Tier Role-Based Access Control
The system has two admin roles:

### Main Admin (account_type: 'main_admin')
- Full access to all features
- Dashboard with sales reports (daily/monthly)
- Create and manage other admin accounts
- Approve/reject password change requests
- Reject payments and cancel orders
- All client management features

### Normal Admin (account_type: 'admin')
- Manage furniture items (add, edit, delete)
- View client profiles
- Send emails to clients
- View orders and delivery locations
- NO access to: sales reports, admin management, password approvals, payment rejection

### Login Pages
- `/login` - Main Admin login only
- `/admin/login` - Regular Admin login

### Navigation
All nav links are visible to all logged-in users, but unauthorized links appear disabled (grayed out).

## Important Security Note
The frontend implements role-based UI restrictions. For production security, the backend API (admin-api) must implement:
- JWT role validation middleware on all protected endpoints
- Server-side authorization checks for main_admin-only endpoints:
  - POST /api/admins (create admin)
  - GET /api/reports/* (sales reports)
  - POST /api/payments/:id/reject
  - GET/POST /api/password-requests

## Recent Changes
- 2025-12-10: Implemented two-tier admin role system
- 2025-12-10: Added separate login pages for main admin and regular admin
- 2025-12-10: Added Orders page with delivery location view and payment rejection
- 2025-12-10: Added Password Change Requests approval page
- 2025-12-10: Added email functionality to clients and orders pages
- 2025-12-10: Updated navbar with role-based link styling and logout
- 2025-12-10: Configured for Replit environment (port 5000, allowedHosts)

## Architecture
- Frontend: SvelteKit with Vite, uses proxy to backend API at /api
- Backend: Rust Actix-web with PostgreSQL (tokio-postgres/deadpool)
- Authentication: JWT-based with role field in user object

## User Preferences
- All navbar items should be visible but disabled/styled when user lacks access
- Separate login pages for main admin vs regular admin
