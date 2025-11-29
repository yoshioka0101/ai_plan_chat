import { apiClient } from './api';
import type {
  InterpretationItem,
  InterpretationItemsResponse,
  InterpretationItemData,
  ApproveItemResponse,
} from '../types/interpretation';

export const interpretationItemService = {
  async getItems(interpretationId: string): Promise<InterpretationItem[]> {
    const response = await apiClient.get<InterpretationItemsResponse>(
      `/interpretations/${interpretationId}/items`
    );
    return response.data.items;
  },

  async updateItem(itemId: string, data: InterpretationItemData): Promise<InterpretationItem> {
    const response = await apiClient.patch<InterpretationItem>(`/interpretation-items/${itemId}`, {
      data,
    });
    return response.data;
  },

  async approveItem(itemId: string): Promise<ApproveItemResponse> {
    const response = await apiClient.post<ApproveItemResponse>(`/interpretation-items/${itemId}/approve`);
    return response.data;
  },
};
