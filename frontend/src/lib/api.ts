import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Types
export interface Queue {
  id: number;
  name: string;
  description: string;
  type: string;
  config: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: number;
  queue_id: number;
  queue?: Queue;
  payload: string;
  priority: number;
  status: string;
  scheduled_at?: string;
  processed_at?: string;
  failed_at?: string;
  retry_count: number;
  max_retries: number;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

export interface Worker {
  id: number;
  name: string;
  queue_id: number;
  queue?: Queue;
  status: string;
  last_ping: string;
  processed_count: number;
  failed_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateQueueRequest {
  name: string;
  description: string;
  type: string;
  config: string;
}

export interface CreateMessageRequest {
  queue_id: number;
  payload: string;
  priority?: number;
  scheduled_at?: string;
  max_retries?: number;
}

// API Functions
export const queueApi = {
  getQueues: () => api.get<Queue[]>('/queues'),
  getQueue: (id: number) => api.get<Queue>(`/queues/${id}`),
  createQueue: (data: CreateQueueRequest) => api.post<Queue>('/queues', data),
  deleteQueue: (id: number) => api.delete(`/queues/${id}`),
};

export const messageApi = {
  getMessages: (queueId?: number, limit?: number) => {
    const params = new URLSearchParams();
    if (queueId) params.append('queue_id', queueId.toString());
    if (limit) params.append('limit', limit.toString());
    return api.get<Message[]>(`/messages?${params.toString()}`);
  },
  createMessage: (data: CreateMessageRequest) => api.post<Message>('/messages', data),
};

export const workerApi = {
  getWorkers: (queueId?: number) => {
    const params = new URLSearchParams();
    if (queueId) params.append('queue_id', queueId.toString());
    return api.get<Worker[]>(`/workers?${params.toString()}`);
  },
};
