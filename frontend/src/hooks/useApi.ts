import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { queueApi, messageApi, workerApi, CreateQueueRequest, CreateMessageRequest } from '../lib/api';

// Queue hooks
export const useQueues = () => {
  return useQuery({
    queryKey: ['queues'],
    queryFn: () => queueApi.getQueues().then(res => res.data),
  });
};

export const useQueue = (id: number) => {
  return useQuery({
    queryKey: ['queue', id],
    queryFn: () => queueApi.getQueue(id).then(res => res.data),
    enabled: !!id,
  });
};

export const useCreateQueue = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateQueueRequest) => queueApi.createQueue(data).then(res => res.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['queues'] });
    },
  });
};

export const useDeleteQueue = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: number) => queueApi.deleteQueue(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['queues'] });
    },
  });
};

// Message hooks
export const useMessages = (queueId?: number, limit?: number) => {
  return useQuery({
    queryKey: ['messages', queueId, limit],
    queryFn: () => messageApi.getMessages(queueId, limit).then(res => res.data),
  });
};

export const useCreateMessage = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateMessageRequest) => messageApi.createMessage(data).then(res => res.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['messages'] });
    },
  });
};

// Worker hooks
export const useWorkers = (queueId?: number) => {
  return useQuery({
    queryKey: ['workers', queueId],
    queryFn: () => workerApi.getWorkers(queueId).then(res => res.data),
  });
};
