schema "ai_plan_chat" {}

table "users" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "nickname" {
    type = text
    null = false
  }
  column "avatar" {
    type = text
    null = true
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  column "updated_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
}

table "user_auths" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "user_id" {
    type = varchar(36)
    null = false
  }
  column "identity_type" {
    type = varchar(50)
    null = false
  }
  column "identifier" {
    type = varchar(255)
    null = false
  }
  column "credential" {
    type = text
    null = true
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  column "updated_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_auths_user_id_fkey" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "user_auths_user_id_identity_type_identifier_key" {
    columns = [column.user_id, column.identity_type, column.identifier]
    unique = true
  }
  index "idx_user_auths_user_id" {
    columns = [column.user_id]
  }
}

table "interpretations" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "user_id" {
    type = varchar(36)
    null = false
  }
  column "name" {
    type = varchar(255)
    null = true
  }
  column "description" {
    type = text
    null = true
  }
  column "context" {
    type = text
    null = true
  }
  column "input_text" {
    type = text
    null = false
  }
  column "structured_result" {
    type = json
    null = true
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "interpretations_user_id_fkey" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "idx_interpretations_user_id" {
    columns = [column.user_id]
  }
}

table "tasks" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "user_id" {
    type = varchar(36)
    null = false
  }
  column "title" {
    type = varchar(500)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "due_at" {
    type = timestamp
    null = true
  }
  column "status" {
    type = varchar(20)
    null = false
    default = "pending"
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  column "updated_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tasks_user_id_fkey" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "idx_tasks_user_id" {
    columns = [column.user_id]
  }
  index "idx_tasks_status" {
    columns = [column.status]
  }
  index "idx_tasks_due_at" {
    columns = [column.due_at]
  }
}

table "events" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "user_id" {
    type = varchar(36)
    null = false
  }
  column "title" {
    type = varchar(500)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "starts_at" {
    type = timestamp
    null = false
  }
  column "ends_at" {
    type = timestamp
    null = false
  }
  column "location" {
    type = varchar(500)
    null = true
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  column "updated_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "events_user_id_fkey" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "idx_events_user_id" {
    columns = [column.user_id]
  }
  index "idx_events_starts_at" {
    columns = [column.starts_at]
  }
}

table "tags" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "user_id" {
    type = varchar(36)
    null = false
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "color" {
    type = varchar(20)
    null = true
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tags_user_id_fkey" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "idx_tags_user_id" {
    columns = [column.user_id]
  }
  index "tags_user_id_name_key" {
    columns = [column.user_id, column.name]
    unique = true
  }
}

table "tag_tasks" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "tag_id" {
    type = varchar(36)
    null = false
  }
  column "task_id" {
    type = varchar(36)
    null = false
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tag_tasks_tag_id_fkey" {
    columns = [column.tag_id]
    ref_columns = [table.tags.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  foreign_key "tag_tasks_task_id_fkey" {
    columns = [column.task_id]
    ref_columns = [table.tasks.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "tag_tasks_tag_id_task_id_key" {
    columns = [column.tag_id, column.task_id]
    unique = true
  }
  index "idx_tag_tasks_tag_id" {
    columns = [column.tag_id]
  }
  index "idx_tag_tasks_task_id" {
    columns = [column.task_id]
  }
}

table "tag_events" {
  schema = schema.ai_plan_chat
  column "id" {
    type = varchar(36)
    null = false
  }
  column "tag_id" {
    type = varchar(36)
    null = false
  }
  column "event_id" {
    type = varchar(36)
    null = false
  }
  column "created_at" {
    type = timestamp
    default = sql("CURRENT_TIMESTAMP")
    null = false
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tag_events_tag_id_fkey" {
    columns = [column.tag_id]
    ref_columns = [table.tags.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  foreign_key "tag_events_event_id_fkey" {
    columns = [column.event_id]
    ref_columns = [table.events.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
  index "tag_events_tag_id_event_id_key" {
    columns = [column.tag_id, column.event_id]
    unique = true
  }
  index "idx_tag_events_tag_id" {
    columns = [column.tag_id]
  }
  index "idx_tag_events_event_id" {
    columns = [column.event_id]
  }
}
