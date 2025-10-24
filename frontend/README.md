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
│   └── Tasks/             # Task-related components
├── hooks/                 # Custom React hooks
├── pages/                 # Page components
├── services/              # API services
└── types/                 # TypeScript type definitions

tests/
├── components/            # Component tests
└── services/              # Service tests
```

## Features

### Phase 1 (MVP)
- ✅ Chat interface for task creation
- ✅ Natural language task input
- ✅ Real-time AI responses
- ✅ Task preview in chat

### Phase 2
- 🔲 Conversation history
- 🔲 Task modification via chat
- 🔲 Multi-turn conversations

## Development Notes

- The app currently uses mock services for AI responses
- Real backend integration will replace mock services
- All components are typed with TypeScript strict mode
