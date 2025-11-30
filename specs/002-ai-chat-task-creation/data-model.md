# Data Model: AI Chat-based Task Creation

**Feature**: AI Chat-based Task Creation
**Date**: 2025-10-23
**Phase**: 1 - Design

## Entity Relationship Diagram

```
┌─────────────────┐
│   User          │
│  (existing)     │
└────────┬────────┘
         │ 1
         │
         │ N
┌────────▼────────────────┐         ┌──────────────────────┐
│   Conversation          │ 1     N │   Message            │
├─────────────────────────┤◄────────┤──────────────────────┤
│ id: UUID (PK)           │         │ id: UUID (PK)        │
│ user_id: UUID (FK)      │         │ conversation_id (FK) │
│ title: VARCHAR(255)     │         │ sender: ENUM         │
│ status: ENUM            │         │ content: TEXT        │
│ created_at: TIMESTAMP   │         │ metadata: JSON       │
│ updated_at: TIMESTAMP   │         │ created_at: TIMESTAMP│
└─────────┬───────────────┘         └──────────┬───────────┘
          │                                    │
          │ N                                  │ 0..1
          │                                    │
          │                         ┌──────────▼───────────┐
          │                         │   Task               │
          │                         │  (existing)          │
          └─────────────────────────┤──────────────────────┤
                                    │ id: UUID (PK)        │
                       (referenced) │ title: VARCHAR(255)  │
                                    │ description: TEXT    │
                                    │ due_date: TIMESTAMP  │
                                    │ status: ENUM         │
                                    │ priority: ENUM       │
                                    │ created_via: VARCHAR │
                                    └──────────────────────┘
```

## Entity Definitions

### Conversation

Represents a chat session between a user and the AI assistant.

**Attributes**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier for the conversation |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | Reference to the user who owns this conversation |
| title | VARCHAR(255) | NOT NULL, DEFAULT "New Conversation" | Auto-generated from first user message (first 50 chars) |
| status | ENUM | NOT NULL, DEFAULT 'active' | Values: 'active', 'archived' |
| created_at | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | When conversation was created |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE | Last message timestamp |

**Indexes**:
- PRIMARY KEY: `id`
- INDEX: `idx_user_conversations` on (`user_id`, `updated_at` DESC) - for user's conversation list
- INDEX: `idx_status` on (`status`) - for filtering archived conversations

**Validation Rules**:
- `title` must be 1-255 characters
- `status` must be one of: ['active', 'archived']
- `user_id` must reference existing user

**State Transitions**:
```
active → archived (user archives conversation)
archived → active (user restores conversation)
```

**Business Rules**:
- Each user can have unlimited conversations
- Archiving a conversation does not delete messages
- Deleted conversations should soft-delete (add `deleted_at` field)

---

### Message

Represents a single message within a conversation (from user or AI).

**Attributes**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier for the message |
| conversation_id | UUID | FOREIGN KEY (conversations.id), NOT NULL | Parent conversation |
| sender | ENUM | NOT NULL | Values: 'user', 'ai' |
| content | TEXT | NOT NULL | The message content (max 10,000 chars enforced in app) |
| metadata | JSON | NULL | Additional data (see metadata schema below) |
| created_at | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | When message was sent |

**Indexes**:
- PRIMARY KEY: `id`
- INDEX: `idx_conversation_messages` on (`conversation_id`, `created_at` ASC) - for chronological message retrieval
- INDEX: `idx_created_at` on (`created_at` DESC) - for recent messages queries

**Metadata Schema** (JSON field):

```json
{
  "ai_model": "gpt-3.5-turbo",          // Which AI model was used
  "confidence": 0.85,                    // AI confidence score (0.0-1.0)
  "intent": "create_task",               // AI detected intent
  "task_id": "uuid",                     // Reference to created/modified task
  "tokens_used": 150,                    // API token consumption
  "processing_time_ms": 1200,            // How long AI took to respond
  "error": "string | null"               // Error message if AI failed
}
```

**Validation Rules**:
- `content` must be 1-10,000 characters
- `sender` must be one of: ['user', 'ai']
- `conversation_id` must reference existing conversation
- `metadata` must be valid JSON if provided

**Business Rules**:
- Messages are immutable (no updates allowed)
- Messages are ordered by `created_at` ASC
- User messages must be followed by AI response (enforced in app layer)
- Deleting a conversation cascades to delete messages

---

### Task (Existing Entity - Extensions Only)

Represents a task in the system. This entity already exists; we only add a field to track creation source.

**New/Modified Attributes**:

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| created_via | VARCHAR(50) | NULL, DEFAULT 'manual' | Values: 'manual', 'chat', 'import', 'api' |
| source_message_id | UUID | FOREIGN KEY (messages.id), NULL | Reference to chat message that created this task |

**Index Additions**:
- INDEX: `idx_created_via` on (`created_via`) - for analytics

**Validation Rules**:
- `created_via` must be one of: ['manual', 'chat', 'import', 'api']
- `source_message_id` must reference existing message if set

**Business Rules**:
- Tasks created via chat have `created_via = 'chat'`
- `source_message_id` links back to the AI message that created the task
- Deleting a message does NOT delete the task (cascade prevention)

---

## Relationship Details

### Conversation ↔ Message (1:N)

- **Cardinality**: One conversation has many messages
- **Cascade**: DELETE conversation → DELETE messages
- **Referential Integrity**: Every message must belong to a conversation
- **Ordering**: Messages ordered by `created_at` ASC within a conversation

### User ↔ Conversation (1:N)

- **Cardinality**: One user has many conversations
- **Cascade**: DELETE user → DELETE conversations (which cascades to messages)
- **Referential Integrity**: Every conversation must belong to a user
- **Ordering**: Conversations ordered by `updated_at` DESC for user's conversation list

### Message → Task (N:1, optional)

- **Cardinality**: Many messages can reference one task (create, update, delete actions)
- **Cascade**: DELETE task → SET NULL on message.metadata.task_id (no cascade)
- **Referential Integrity**: Weak reference via JSON metadata
- **Note**: This is not a foreign key constraint, just a JSON field for traceability

## Database Migration Schema (Atlas HCL)

### File: `schemas/conversation.hcl`

```hcl
table "conversations" {
  schema = schema.ai_plan_chat

  column "id" {
    type = binary(16)
    comment = "UUID stored as BINARY(16)"
  }

  column "user_id" {
    type = binary(16)
    null = false
    comment = "FK to users table"
  }

  column "title" {
    type = varchar(255)
    null = false
    default = "New Conversation"
  }

  column "status" {
    type = enum("active", "archived")
    null = false
    default = "active"
  }

  column "created_at" {
    type = timestamp
    null = false
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    type = timestamp
    null = false
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_conversations_user" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = CASCADE
  }

  index "idx_user_conversations" {
    columns = [column.user_id, column.updated_at]
  }

  index "idx_status" {
    columns = [column.status]
  }
}

table "messages" {
  schema = schema.ai_plan_chat

  column "id" {
    type = binary(16)
    comment = "UUID stored as BINARY(16)"
  }

  column "conversation_id" {
    type = binary(16)
    null = false
    comment = "FK to conversations table"
  }

  column "sender" {
    type = enum("user", "ai")
    null = false
  }

  column "content" {
    type = text
    null = false
  }

  column "metadata" {
    type = json
    null = true
    comment = "Stores AI confidence, intent, task references, etc."
  }

  column "created_at" {
    type = timestamp
    null = false
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_messages_conversation" {
    columns = [column.conversation_id]
    ref_columns = [table.conversations.column.id]
    on_delete = CASCADE
    on_update = CASCADE
  }

  index "idx_conversation_messages" {
    columns = [column.conversation_id, column.created_at]
  }

  index "idx_created_at" {
    columns = [column.created_at]
  }
}
```

### Migration for Task Extension

```hcl
table "tasks" {
  schema = schema.ai_plan_chat

  # ... existing columns ...

  column "created_via" {
    type = varchar(50)
    null = true
    default = "manual"
  }

  column "source_message_id" {
    type = binary(16)
    null = true
    comment = "Optional FK to messages table"
  }

  index "idx_created_via" {
    columns = [column.created_via]
  }

  # Note: We intentionally do NOT create a foreign key for source_message_id
  # to avoid deletion cascades. This is a weak reference for audit purposes only.
}
```

## Data Access Patterns

### Common Queries

**1. Get user's active conversations (ordered by most recent)**
```sql
SELECT * FROM conversations
WHERE user_id = ? AND status = 'active'
ORDER BY updated_at DESC
LIMIT 20;
```

**2. Get messages for a conversation (chronological)**
```sql
SELECT * FROM messages
WHERE conversation_id = ?
ORDER BY created_at ASC;
```

**3. Get recent messages (since timestamp) for polling**
```sql
SELECT * FROM messages
WHERE conversation_id = ? AND created_at > ?
ORDER BY created_at ASC;
```

**4. Get messages with pagination (cursor-based)**
```sql
SELECT * FROM messages
WHERE conversation_id = ? AND created_at > ?
ORDER BY created_at ASC
LIMIT 50;
```

**5. Get tasks created via chat**
```sql
SELECT * FROM tasks
WHERE created_via = 'chat' AND user_id = ?
ORDER BY created_at DESC;
```

**6. Get conversation with message count**
```sql
SELECT c.*, COUNT(m.id) as message_count
FROM conversations c
LEFT JOIN messages m ON m.conversation_id = c.id
WHERE c.user_id = ?
GROUP BY c.id
ORDER BY c.updated_at DESC;
```

### Performance Considerations

**Index Usage**:
- `idx_user_conversations`: Optimizes user's conversation list query
- `idx_conversation_messages`: Optimizes message retrieval within a conversation
- `idx_created_at`: Enables efficient "since timestamp" polling queries

**Query Optimization**:
- Use LIMIT on conversation and message queries to prevent large result sets
- Implement cursor-based pagination for messages (better than OFFSET)
- Consider materialized view for conversation message counts if query is slow

**Storage Optimization**:
- Archive old conversations periodically (move to `status = 'archived'`)
- Consider partitioning messages table by `created_at` if volume exceeds 10M rows
- Compress JSON metadata field if MySQL version supports it

## Validation Summary

### Conversation Validations

- ✅ Title length: 1-255 characters
- ✅ Status enum: 'active' or 'archived'
- ✅ User exists: Foreign key constraint
- ✅ Created/Updated timestamps: Auto-managed by database

### Message Validations

- ✅ Content length: 1-10,000 characters (enforced at application layer)
- ✅ Sender enum: 'user' or 'ai'
- ✅ Conversation exists: Foreign key constraint
- ✅ Metadata JSON: Valid JSON structure
- ✅ Created timestamp: Auto-managed by database

### Task Validations (Extensions)

- ✅ Created via enum: 'manual', 'chat', 'import', 'api'
- ✅ Source message reference: Optional, weak reference (no FK)

## Future Considerations

### Potential Schema Enhancements (Not in MVP)

1. **Conversation Tags**: Add M:N relationship for categorizing conversations
2. **Message Reactions**: Allow users to react to AI responses (helpful/not helpful)
3. **Conversation Templates**: Pre-defined conversation starters
4. **Message Edits**: Track edit history for user messages
5. **Shared Conversations**: Multi-user collaboration on conversations
6. **Message Attachments**: Support file uploads in chat
7. **Conversation Analytics**: Aggregate metrics (response times, success rates)

### Scalability Considerations

**When to Shard**:
- Shard by `user_id` when user base exceeds 1M active users
- Partition messages by `created_at` when table exceeds 100M rows

**Caching Strategy**:
- Cache recent conversations (last 24 hours) in Redis
- Cache conversation message counts
- Invalidate cache on new message creation

**Archival Strategy**:
- Move conversations older than 90 days to cold storage
- Implement soft delete with `deleted_at` timestamp
- Periodic cleanup jobs for truly deleted data
