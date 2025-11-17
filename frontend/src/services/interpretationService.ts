import { apiClient } from './api';
import type {
  CreateInterpretationRequest,
  InterpretationResponse,
  InterpretationsListResponse,
  AIInterpretation,
} from '../types/interpretation';

export const interpretationService = {
  /**
   * Create a new AI interpretation from natural language input
   */
  async createInterpretation(
    request: CreateInterpretationRequest
  ): Promise<InterpretationResponse> {
    const response = await apiClient.post<InterpretationResponse>(
      '/interpretations',
      request
    );
    return response.data;
  },

  /**
   * Get list of AI interpretations
   */
  async listInterpretations(
    limit: number = 20,
    offset: number = 0
  ): Promise<InterpretationsListResponse> {
    const response = await apiClient.get<InterpretationsListResponse>(
      '/interpretations',
      {
        params: { limit, offset },
      }
    );
    return response.data;
  },

  /**
   * Get a specific AI interpretation by ID
   */
  async getInterpretation(id: string): Promise<AIInterpretation> {
    const response = await apiClient.get<AIInterpretation>(`/interpretations/${id}`);
    return response.data;
  },
};
