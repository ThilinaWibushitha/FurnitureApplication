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

## Recent Changes
- 2025-12-10: Configured for Replit environment (port 5000, allowedHosts)

## Architecture
- Frontend: SvelteKit with Vite, uses proxy to backend API at /api
- Backend: Rust Actix-web with PostgreSQL (tokio-postgres/deadpool)
- Authentication: JWT-based

## User Preferences
(none recorded yet)
