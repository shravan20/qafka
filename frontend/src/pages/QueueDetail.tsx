import React from 'react';
import { useParams } from 'react-router-dom';
import { ArrowLeft, Send, Users, BarChart3 } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { useQueue, useMessages, useWorkers } from '../hooks/useApi';

const QueueDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const queueId = parseInt(id || '0', 10);
  
  const { data: queue, isLoading: queueLoading } = useQueue(queueId);
  const { data: messages, isLoading: messagesLoading } = useMessages(queueId, 50);
  const { data: workers, isLoading: workersLoading } = useWorkers(queueId);

  if (queueLoading) {
    return <div className="text-center py-8">Loading queue details...</div>;
  }

  if (!queue) {
    return <div className="text-center py-8">Queue not found</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center space-x-4">
        <Button variant="outline" onClick={() => window.history.back()}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back
        </Button>
        <div>
          <h1 className="text-3xl font-bold">{queue.name}</h1>
          <p className="text-muted-foreground">{queue.description}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Messages</CardTitle>
            <Send className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{messages?.length || 0}</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Workers</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {workers?.filter(w => w.status === 'busy' || w.status === 'idle').length || 0}
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Queue Type</CardTitle>
            <BarChart3 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold capitalize">{queue.type}</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Recent Messages</CardTitle>
            <CardDescription>Latest messages in the queue</CardDescription>
          </CardHeader>
          <CardContent>
            {messagesLoading ? (
              <div>Loading messages...</div>
            ) : messages && messages.length > 0 ? (
              <div className="space-y-2">
                {messages.slice(0, 10).map((message) => (
                  <div key={message.id} className="p-3 border rounded-lg">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <p className="text-sm font-medium">ID: {message.id}</p>
                        <p className="text-xs text-muted-foreground">
                          {message.payload.substring(0, 100)}
                          {message.payload.length > 100 ? '...' : ''}
                        </p>
                      </div>
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        message.status === 'completed' ? 'bg-green-100 text-green-800' :
                        message.status === 'failed' ? 'bg-red-100 text-red-800' :
                        message.status === 'processing' ? 'bg-blue-100 text-blue-800' :
                        'bg-yellow-100 text-yellow-800'
                      }`}>
                        {message.status}
                      </span>
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      Created: {new Date(message.created_at).toLocaleString()}
                    </p>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center text-muted-foreground py-8">
                No messages in this queue
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Workers</CardTitle>
            <CardDescription>Active workers processing this queue</CardDescription>
          </CardHeader>
          <CardContent>
            {workersLoading ? (
              <div>Loading workers...</div>
            ) : workers && workers.length > 0 ? (
              <div className="space-y-2">
                {workers.map((worker) => (
                  <div key={worker.id} className="p-3 border rounded-lg">
                    <div className="flex justify-between items-center">
                      <div>
                        <p className="text-sm font-medium">{worker.name}</p>
                        <p className="text-xs text-muted-foreground">
                          Processed: {worker.processed_count} | Failed: {worker.failed_count}
                        </p>
                      </div>
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        worker.status === 'busy' ? 'bg-blue-100 text-blue-800' :
                        worker.status === 'idle' ? 'bg-green-100 text-green-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {worker.status}
                      </span>
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">
                      Last ping: {new Date(worker.last_ping).toLocaleString()}
                    </p>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center text-muted-foreground py-8">
                No workers assigned to this queue
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default QueueDetail;
