# AI Chat Task Creation - Frontend

React + TypeScript frontend for the AI Chat Task Creation application.

## Tech Stack

- **React 18** with TypeScript
- **Vite** for build tooling
- **Axios** for API calls
- **React Router** for routing
- **date-fns** for date handling

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create `.env` file from template:
   ```bash
   cp .env.example .env
   ```

3. Update `.env` with your API endpoint:
   ```
   VITE_API_BASE_URL=http://localhost:8080/api/v1
   ```

### Development

Start the development server:

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Build

Build for production:

```bash
npm run build
```

### Lint

Run linter:

```bash
npm run lint
```

Format code with Prettier:

```bash
npx prettier --write src/
```

## Project Structure

```
src/
├── components/
│   ├── Chat/              # Chat-related components
│   │   ├── AIChat.tsx     # AI chat interface
│   │   └── InterpretationHistory.tsx  # AI interpretation history
│   ├── Modal/             # Modal components
│   ├── Sidebar/           # Sidebar components
│   └── Tasks/             # Task-related components
├── contexts/              # React contexts
├── hooks/                 # Custom React hooks
├── pages/                 # Page components
│   ├── AuthCallbackPage.tsx
│   ├── DashboardPage.tsx
│   └── LoginPage.tsx
├── services/              # API services
│   ├── api.ts             # Base API client
│   ├── authService.ts     # Authentication service
│   ├── interpretationService.ts  # AI interpretation service
│   └── taskService.ts     # Task service
└── types/                 # TypeScript type definitions
    ├── auth.ts
    ├── interpretation.ts  # AI interpretation types
    └── task.ts

tests/
├── components/            # Component tests
└── services/              # Service tests
```

## Features

### Implemented Features
- ✅ Google OAuth authentication
- ✅ Task management (CRUD operations)
- ✅ AI-powered interpretation of natural language
- ✅ AI Chat interface with real-time responses
- ✅ AI interpretation history view
- ✅ Real backend integration with Gemini API

### AI Features
The application now includes full AI capabilities:

1. **AI Chat Interface**:
   - Natural language input for task creation
   - Real-time AI responses using Gemini API
   - Interactive chat experience
   - Support for various input types (todos, reminders, questions)

2. **AI Interpretation History**:
   - View all past AI interpretations
   - Detailed view of interpretation results
   - Metadata display (priority, deadline, tags)
   - Search and filter capabilities

3. **Structured AI Output**:
   - Automatic extraction of task titles
   - Description generation
   - Priority detection
   - Deadline parsing
   - Tag suggestions

## How to Use AI Features

1. **Login**: Use Google OAuth to authenticate
2. **Navigate to AI Chat**: Click on the "AI Chat" tab in the dashboard
3. **Start Chatting**: Type natural language requests like:
   - "I need to finish the project report by Friday"
   - "Remind me to call John tomorrow at 3pm"
   - "Create a high priority task to review the code"
4. **View History**: Click on "AI History" to see all past interpretations

## Development Notes

- The app integrates with a Go backend using Gemini API for AI interpretations
- All API communication is handled through Axios with automatic auth token injection
- TypeScript strict mode is enabled for type safety
- Components follow React best practices with hooks
- CSS modules are used for component styling
