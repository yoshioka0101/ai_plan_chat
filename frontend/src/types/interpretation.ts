import type { TaskStatus } from './task';

// AI Interpretation type definitions

export type InterpretationType = 'todo' | 'reminder' | 'question' | 'other';
export type PriorityType = 'low' | 'medium' | 'high';

export interface InterpretationMetadata {
  deadline?: string;
  priority?: PriorityType;
  tags?: string[];
}

export interface StructuredResult {
  title?: string;
  description?: string;
  type?: InterpretationType;
  metadata?: InterpretationMetadata;
}

export interface AIInterpretation {
  id: string;
  user_id: string;
  input_text: string;
  structured_result: StructuredResult;
  ai_model: string;
  ai_prompt_tokens?: number;
  ai_completion_tokens?: number;
  created_at: string;
}

export interface CreateInterpretationRequest {
  input_text: string;
}

export interface InterpretationResponse {
  type: InterpretationType;
  interpretation: AIInterpretation;
  message?: string;
}

export interface InterpretationsListResponse {
  interpretations: AIInterpretation[];
  total: number;
  limit: number;
  offset: number;
}

export type InterpretationItemStatus = 'pending' | 'created';

export type ResourceType = 'task' | 'event' | 'wallet';

export interface InterpretationItemData {
  title?: string;
  description?: string;
  due_at?: string;
  status?: TaskStatus;
  tags?: string[];
  [key: string]: unknown;
}

export interface InterpretationItem {
  id: string;
  interpretation_id: string;
  item_index: number;
  resource_type: ResourceType;
  resource_id?: string;
  status: InterpretationItemStatus;
  data: InterpretationItemData;
  original_data: Record<string, unknown>;
  reviewed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface InterpretationItemsResponse {
  items: InterpretationItem[];
}

export interface ApproveItemResponse {
  resource_id: string;
}
