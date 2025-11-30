# Research: AI Chat-based Task Creation

**Feature**: AI Chat-based Task Creation
**Date**: 2025-10-23
**Phase**: 0 - Research & Outline

## Overview

This document consolidates research findings for implementing a conversational AI interface that enables task creation and management through natural language interaction.

## Technology Stack Research

### 1. AI Service Integration

**Decision**: OpenAI API (GPT-3.5/GPT-4) with function calling

**Rationale**:
- Function calling feature allows structured extraction of task attributes from natural language
- Well-documented Go SDK available
- Reliable uptime and performance
- Support for system prompts to guide task extraction behavior
- Streaming responses for better UX

**Alternatives Considered**:
- **Anthropic Claude**: Excellent at natural language understanding but less established function calling patterns in Go ecosystem
- **Self-hosted LLM**: Higher control but significant operational overhead, latency concerns
- **Azure OpenAI**: Similar capabilities but additional complexity for non-Azure deployments

**Implementation Notes**:
- Use OpenAI Go SDK: `github.com/sashabaranov/go-openai`
- Implement retry logic with exponential backoff for API failures
- Cache API responses when appropriate to reduce costs
- Set timeout to 5 seconds (buffer beyond 3s requirement)

### 2. Real-time Communication Pattern

**Decision**: HTTP polling with optional WebSocket upgrade path

**Rationale**:
- Start with simple HTTP polling (3-5 second intervals) to minimize complexity
- Most chat interactions are request-response, not truly real-time
- WebSocket adds connection management complexity
- Can upgrade to WebSocket later if real-time requirements emerge

**Alternatives Considered**:
- **WebSocket from start**: Over-engineering for initial MVP
- **Server-Sent Events (SSE)**: Good for one-way updates but overkill for chat
- **Long polling**: More complex than simple polling without clear benefits

**Implementation Notes**:
- Frontend polls `/api/conversations/{id}/messages` every 3 seconds
- Use `?since=<timestamp>` query param to fetch only new messages
- Add WebSocket endpoints in Phase 2 if needed

### 3. Natural Language Date Parsing

**Decision**: Use `github.com/olebedev/when` for Go-based date parsing

**Rationale**:
- Handles relative dates ("tomorrow", "next Friday", "in 2 days")
- Pure Go implementation, no external dependencies
- Supports multiple languages (potential for i18n)
- Good test coverage and active maintenance

**Alternatives Considered**:
- **Duckling (Facebook)**: More accurate but requires Haskell runtime
- **Chrono.js**: JavaScript-only, can't use in Go backend
- **Let AI handle it**: Inconsistent, wastes API calls

**Implementation Notes**:
- Parse dates in backend before sending to AI
- Fall back to AI interpretation if `when` library fails
- Always confirm ambiguous dates with user

### 4. Frontend State Management

**Decision**: React Context API + custom hooks

**Rationale**:
- Sufficient for chat state management (messages, conversation history)
- No external dependencies (Redux, MobX)
- Simpler learning curve for team
- Can upgrade to Zustand/Redux later if state complexity grows

**Alternatives Considered**:
- **Redux**: Over-engineering for current scope
- **Zustand**: Lightweight but adds dependency
- **React Query**: Excellent for server state but chat needs local state management

**Implementation Notes**:
- Create `ChatContext` for active conversation state
- Create `ConversationContext` for conversation list
- Use `useChat()` and `useConversation()` custom hooks
- Implement optimistic updates for better UX

## Architecture Patterns

### 1. AI Prompt Engineering

**Decision**: Structured system prompt with JSON schema output

**System Prompt Template**:
```
You are a task management assistant. Extract task information from user messages.

Respond with JSON:
{
  "intent": "create_task" | "update_task" | "delete_task" | "query" | "clarify",
  "confidence": 0.0-1.0,
  "task": {
    "title": "string",
    "description": "string | null",
    "due_date": "ISO 8601 | null",
    "priority": "low" | "medium" | "high" | null
  },
  "clarification_needed": "string | null",
  "referenced_task_id": "uuid | null"
}

If confidence < 0.7, ask for clarification.
```

**Rationale**:
- Structured output ensures consistent parsing
- Confidence score enables intelligent clarification requests
- Intent classification handles multiple conversation types
- JSON parsing is more reliable than free-form text

### 2. Conversation State Machine

**States**:
- `IDLE`: Waiting for user input
- `PROCESSING`: AI analyzing message
- `AWAITING_CONFIRMATION`: AI needs user confirmation
- `CREATING_TASK`: Task being created via CRUD API
- `ERROR`: Recoverable error state

**Transitions**:
```
IDLE → [user sends message] → PROCESSING
PROCESSING → [AI returns high confidence] → CREATING_TASK
PROCESSING → [AI returns low confidence] → AWAITING_CONFIRMATION
AWAITING_CONFIRMATION → [user confirms] → CREATING_TASK
CREATING_TASK → [success] → IDLE
CREATING_TASK → [failure] → ERROR
ERROR → [user retries] → PROCESSING
```

**Rationale**:
- Clear state transitions improve reliability
- Easy to test each state independently
- Supports error recovery and retries
- Enables proper UI feedback at each stage

### 3. Message Threading and Context

**Decision**: Include last N messages as context for AI

**Rationale**:
- Multi-turn conversations require context
- Limit to last 5-10 messages to control token usage
- Summarize older context if conversation extends beyond limit

**Implementation**:
- Store full conversation history in database
- Send last 5 messages to AI for context
- Use token counting to prevent exceeding limits
- Implement context summarization for long conversations

### 4. Error Handling Strategy

**Layers**:
1. **AI Service Layer**: Retry with exponential backoff (3 attempts)
2. **API Layer**: Return user-friendly error messages
3. **Frontend Layer**: Display errors inline, preserve user input
4. **Fallback**: Save draft message locally, allow retry

**Error Categories**:
- **Transient** (AI API timeout): Auto-retry
- **Permanent** (invalid API key): Alert admin, show fallback UI
- **User Error** (empty message): Immediate validation feedback
- **Network** (connection lost): Queue message, retry on reconnect

## Database Schema Considerations

### Conversation & Message Tables

**Conversation**:
```sql
id (UUID, primary key)
user_id (UUID, foreign key) - for future multi-user support
created_at (timestamp)
updated_at (timestamp)
title (varchar) - auto-generated from first message
status (enum: active, archived)
```

**Message**:
```sql
id (UUID, primary key)
conversation_id (UUID, foreign key)
sender (enum: user, ai)
content (text)
created_at (timestamp)
metadata (json) - stores AI confidence, intent, task reference
```

**Indexes**:
- `conversation_id` on messages table
- `created_at` DESC on messages table for chronological queries
- `user_id` on conversations table for user filtering

## Performance Optimization

### 1. Caching Strategy

**Decision**: Redis cache for AI responses

**What to Cache**:
- Similar user queries → AI responses (24 hour TTL)
- User's recent tasks → reduce DB queries
- Conversation summaries → for context window optimization

**Rationale**:
- Reduce AI API costs for common queries
- Improve response times
- Handle AI service outages gracefully

**Implementation**:
- Use cache key: `hash(user_message + context_hash)`
- Invalidate on task updates
- Optional for MVP, implement in Phase 2

### 2. Message Pagination

**Decision**: Cursor-based pagination with 50 messages per page

**Rationale**:
- Offset pagination issues with real-time data
- Cursor (last message timestamp) more reliable
- 50 messages balance UX and performance

**API Pattern**:
```
GET /api/conversations/{id}/messages?limit=50&cursor=<timestamp>
```

## Security Considerations

### 1. Input Sanitization

**Requirements**:
- Validate message length (max 5000 characters)
- Strip HTML/script tags from user input
- Rate limit: 10 messages per minute per user
- Validate AI responses before persisting

### 2. AI Prompt Injection Prevention

**Mitigations**:
- User input clearly marked in system prompt
- Validate AI JSON response structure
- Reject responses that don't match schema
- Log suspicious patterns for review

### 3. Data Privacy

**Considerations**:
- Conversations may contain sensitive task information
- Comply with OpenAI data usage policies
- Provide clear user consent for AI processing
- Option to disable AI features and use forms

## Integration Points

### 1. Existing Task CRUD API

**Integration Approach**:
- Chat service calls existing task endpoints internally
- Reuse validation logic from task handlers
- Maintain audit trail (task created via "chat")
- Use same domain entities

**Benefits**:
- No duplication of business logic
- Consistent validation rules
- Existing tests remain valid
- Clear separation of concerns

### 2. User Authentication

**Assumption**: Existing auth middleware handles user identification

**Requirements**:
- Extract `user_id` from auth context
- Associate conversations with authenticated users
- Enforce authorization on conversation access
- Handle anonymous users (future consideration)

## Testing Strategy

### 1. AI Integration Testing

**Challenges**:
- Non-deterministic AI responses
- API costs for extensive testing
- Rate limits during test runs

**Solutions**:
- Mock AI service for unit tests
- Record/replay AI responses for integration tests
- Use separate test API key with lower rate limits
- Test prompt engineering with small dataset first

### 2. Frontend Testing

**Coverage**:
- Component tests: MessageList, MessageInput, ChatInterface
- Integration tests: Full conversation flow
- E2E tests: Create task via chat, verify in task list

**Tools**:
- Jest + React Testing Library
- Mock Service Worker for API mocking
- Playwright for E2E tests (optional)

### 3. Performance Testing

**Metrics**:
- AI response time (target: < 3s p95)
- Message history load time (target: < 500ms for 100 messages)
- Concurrent conversation handling (target: 50 simultaneous users)

## Open Questions & Future Considerations

### Resolved for MVP

1. **Q**: Which AI service?
   **A**: OpenAI API with function calling

2. **Q**: Real-time or polling?
   **A**: Start with polling, upgrade to WebSocket if needed

3. **Q**: Date parsing approach?
   **A**: Go library (`when`) with AI fallback

4. **Q**: Frontend state management?
   **A**: React Context + custom hooks

### Deferred to Later Phases

1. **Voice input**: Speech-to-text for hands-free task creation
2. **Multi-language support**: i18n for non-English users
3. **Task templates**: AI suggests task templates based on patterns
4. **Smart scheduling**: AI proposes optimal due dates based on workload
5. **Conversation search**: Full-text search across conversation history
6. **Export conversations**: Download chat history as markdown/PDF

## Cost Estimation

### AI API Costs (Monthly)

**Assumptions**:
- 1000 active users
- 10 messages per user per day
- Average 150 tokens per request
- GPT-3.5-turbo pricing: $0.001/1K input tokens, $0.002/1K output tokens

**Calculation**:
```
Input: 1000 users × 10 messages × 30 days × 150 tokens = 45M tokens
Cost: 45M / 1000 × $0.001 = $45/month input

Output: ~100 tokens avg × 300K messages = 30M tokens
Cost: 30M / 1000 × $0.002 = $60/month output

Total: ~$105/month
```

**Optimization**:
- Caching reduces this by ~30-40%: **~$65-75/month**
- Using GPT-4 would increase costs 10-20x

### Storage Costs

**Assumptions**:
- 500 bytes per message average
- 300K messages per month
- MySQL storage

**Calculation**:
```
300K × 500 bytes = 150MB/month
Annual growth: ~1.8GB
```

**Cost**: Negligible with standard hosting plans

## Implementation Phases Summary

### Phase 0: Research (Complete)
✅ AI service selection
✅ Architecture patterns
✅ Database schema design
✅ Testing strategy
✅ Cost estimation

### Phase 1: Design (Next)
- Data model diagrams
- OpenAPI contract for chat endpoints
- Component architecture diagrams
- Quickstart guide for development

### Phase 2: Task Breakdown
- Granular implementation tasks
- Dependency ordering
- Effort estimation
- Milestone planning

## References

- [OpenAI Function Calling Documentation](https://platform.openai.com/docs/guides/function-calling)
- [Go OpenAI SDK](https://github.com/sashabaranov/go-openai)
- [When Date Parser](https://github.com/olebedev/when)
- [React Context API](https://react.dev/reference/react/useContext)
- [BOB Query Builder](https://github.com/stephenafamo/bob)
- [Atlas Migrations](https://atlasgo.io/)
