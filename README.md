# Flame CRM

A simple CRM system built with Go (Backend) and React + TypeScript + Shadcn UI (Frontend).

## Prerequisites

- Go 1.23+
- Node.js 18+
- PostgreSQL

## Setup

2. **Database**:
   You can use the built-in CLI tool to create the database and seed it.
   ```bash
   cd backend
   # Create the database (flame_crm)
   go run cmd/manage/main.go createdb
   
   # Seed the database with an initial Admin user
   go run cmd/manage/main.go seed
   ```
   *Default Admin User:*
   - Email: `admin@example.com`
   - Password: `password123`

3. **Backend**:
   ```bash
   cd backend
   # Create a .env file from .env.example and configure your database and JWT secret
   cp .env.example .env
   
   # Edit .env to set DB_HOST, DB_USER, DB_PASSWORD, DB_NAME etc.
   
   go run cmd/server/main.go
   ```
   The server runs on `http://localhost:8080`.

4. **Frontend**:

## Features

- **Authentication**: Sign up (first user becomes Admin) and Sign in.
- **Users**: Manage system users (Admin, Sales, Head of Sales).
- **Companies**: Manage client companies.
- **Customers**: Manage customers associated with companies and funnels.
- **Dashboard**: Overview of key metrics.

## Tech Stack

- **Backend**: Go, Gin, GORM, PostgreSQL.
- **Frontend**: React, TypeScript, Vite, Tailwind CSS, Shadcn UI.
