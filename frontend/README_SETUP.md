# Frontend Setup Guide

## Prerequisites

- Node.js 18+ and npm
- Backend server running at `http://localhost:8080`

## Installation

1. Install dependencies:
```bash
npm install
```

2. Create `.env` file from example:
```bash
cp .env.example .env
```

3. Configure environment variables in `.env`:
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_GOOGLE_AUTH_URL=http://localhost:8080/auth/google/callback
VITE_ENV=development
```

## Running the Application

### Development Mode
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Build for Production
```bash
npm run build
```

### Preview Production Build
```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/         # Reusable components
│   │   ├── ProtectedRoute.tsx  # Route guard for authentication
│   │   ├── Modal/          # Modal component
│   │   └── Tasks/          # Task-related components
│   ├── contexts/           # React contexts
│   │   └── AuthContext.tsx # Authentication context
│   ├── pages/              # Page components
│   │   ├── LoginPage.tsx   # Google OAuth login page
│   │   ├── AuthCallbackPage.tsx  # OAuth callback handler
│   │   └── DashboardPage.tsx     # Main dashboard
│   ├── services/           # API services
│   │   ├── api.ts          # Axios client configuration
│   │   ├── authService.ts  # Authentication service
│   │   └── taskService.ts  # Task management service
│   ├── types/              # TypeScript type definitions
│   │   ├── auth.ts         # Auth-related types
│   │   └── task.ts         # Task-related types
│   ├── App.tsx             # Main app component with routing
│   └── main.tsx            # Application entry point
├── public/                 # Static assets
├── .env.example            # Environment variables template
└── package.json            # Project dependencies
```

## Authentication Flow

1. User clicks "Sign in with Google" on the login page
2. Browser redirects to Google OAuth consent screen
3. After user approves, Google redirects back to backend at `/auth/google/callback`
4. Backend processes the OAuth callback and redirects to frontend at `/auth/callback?token=xxx&user=xxx`
5. Frontend AuthCallbackPage extracts token and user data from URL
6. Frontend stores credentials and redirects to dashboard

## Features

### Implemented
- Google OAuth authentication
- Protected routes with authentication guard
- Task management UI (List, Kanban, Calendar views)
- Responsive dashboard layout
- Token-based API authentication

### Coming Soon
- AI Interpretation feature
- Task priority management
- User profile settings

## API Integration

The frontend communicates with the backend API using Axios. All requests automatically include the JWT token in the Authorization header when the user is authenticated.

### Available Endpoints
- `GET /api/v1/tasks` - Get all tasks
- `POST /api/v1/tasks` - Create new task
- `GET /api/v1/tasks/:id` - Get single task
- `PUT /api/v1/tasks/:id` - Update task (full)
- `PATCH /api/v1/tasks/:id` - Edit task (partial)
- `DELETE /api/v1/tasks/:id` - Delete task

## Troubleshooting

### CORS Errors
Make sure the backend CORS configuration includes `http://localhost:5173` in allowed origins.

### Authentication Issues
1. Verify Google OAuth credentials are configured in backend
2. Check that VITE_GOOGLE_AUTH_URL points to the correct backend endpoint
3. Ensure backend is running and accessible

### Build Errors
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Development Tips

- Use React DevTools for debugging component state
- Check browser console for API errors
- Use Network tab to inspect API requests/responses
- Token is stored in localStorage as 'token'
- User data is stored in localStorage as 'user'
