import React, { useState } from 'react';
import { Plus, Activity, Clock, AlertCircle } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { useQueues, useCreateQueue } from '../hooks/useApi';

const Dashboard: React.FC = () => {
  const { data: queues, isLoading, error } = useQueues();
  const createQueueMutation = useCreateQueue();
  const [isCreating, setIsCreating] = useState(false);
  const [newQueue, setNewQueue] = useState({
    name: '',
    description: '',
    type: 'fifo',
    config: '{}',
  });

  const handleCreateQueue = async () => {
    try {
      await createQueueMutation.mutateAsync(newQueue);
      setIsCreating(false);
      setNewQueue({ name: '', description: '', type: 'fifo', config: '{}' });
    } catch (error) {
      console.error('Failed to create queue:', error);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg">Loading queues...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-lg text-destructive">Failed to load queues</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold">Queue Dashboard</h2>
        <Button onClick={() => setIsCreating(true)}>
          <Plus className="w-4 h-4 mr-2" />
          Create Queue
        </Button>
      </div>

      {isCreating && (
        <Card>
          <CardHeader>
            <CardTitle>Create New Queue</CardTitle>
            <CardDescription>
              Set up a new message queue for your application
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Input
              placeholder="Queue name"
              value={newQueue.name}
              onChange={(e) => setNewQueue({ ...newQueue, name: e.target.value })}
            />
            <Input
              placeholder="Description"
              value={newQueue.description}
              onChange={(e) => setNewQueue({ ...newQueue, description: e.target.value })}
            />
            <select
              className="w-full p-2 border rounded-md"
              value={newQueue.type}
              onChange={(e) => setNewQueue({ ...newQueue, type: e.target.value })}
            >
              <option value="fifo">FIFO</option>
              <option value="priority">Priority</option>
              <option value="delay">Delay</option>
            </select>
            <div className="flex space-x-2">
              <Button onClick={handleCreateQueue} disabled={!newQueue.name}>
                Create
              </Button>
              <Button variant="outline" onClick={() => setIsCreating(false)}>
                Cancel
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {queues?.map((queue) => (
          <Card key={queue.id} className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <div className="flex justify-between items-start">
                <div>
                  <CardTitle className="text-lg">{queue.name}</CardTitle>
                  <CardDescription>{queue.description}</CardDescription>
                </div>
                <div className={`px-2 py-1 rounded-full text-xs ${
                  queue.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                }`}>
                  {queue.is_active ? 'Active' : 'Inactive'}
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex items-center space-x-2">
                  <Activity className="w-4 h-4 text-blue-500" />
                  <span className="text-sm">Type: {queue.type}</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Clock className="w-4 h-4 text-green-500" />
                  <span className="text-sm">
                    Created: {new Date(queue.created_at).toLocaleDateString()}
                  </span>
                </div>
                <Button 
                  variant="outline" 
                  className="w-full"
                  onClick={() => window.location.href = `/queue/${queue.id}`}
                >
                  View Details
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {queues?.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <AlertCircle className="w-12 h-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">No queues found</h3>
            <p className="text-muted-foreground text-center mb-4">
              Get started by creating your first message queue
            </p>
            <Button onClick={() => setIsCreating(true)}>
              <Plus className="w-4 h-4 mr-2" />
              Create Your First Queue
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default Dashboard;
