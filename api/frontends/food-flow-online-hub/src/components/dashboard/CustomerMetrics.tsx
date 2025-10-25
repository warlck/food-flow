
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, TrendingUp } from 'lucide-react';

interface Metric {
  title: string;
  value: string | number;
  change: number;
  icon: React.ReactNode;
}

export const CustomerMetrics: React.FC = () => {
  const metrics: Metric[] = [
    {
      title: 'New Customers',
      value: '32',
      change: 12.3,
      icon: <Users className="h-5 w-5 text-blue-500" />,
    },
    {
      title: 'Retention Rate',
      value: '68%',
      change: 5.2,
      icon: <TrendingUp className="h-5 w-5 text-green-500" />,
    }
  ];

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-lg font-medium">Customer Growth</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {metrics.map((metric) => (
            <div key={metric.title} className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <div className="rounded-full bg-muted p-1">
                  {metric.icon}
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">{metric.title}</p>
                  <p className="text-xl font-bold">{metric.value}</p>
                </div>
              </div>
              <div className={`flex items-center ${metric.change > 0 ? 'text-green-500' : 'text-red-500'}`}>
                <TrendingUp className={`h-4 w-4 ${metric.change > 0 ? '' : 'rotate-180'} mr-1`} />
                <span className="text-sm font-medium">{metric.change}%</span>
              </div>
            </div>
          ))}
        </div>

        <div className="mt-4 pt-4 border-t">
          <div className="text-sm text-muted-foreground">Customer Segments</div>
          <div className="mt-2 space-y-1">
            <div className="flex justify-between items-center text-sm">
              <span>New customers</span>
              <span>32%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div className="bg-blue-500 h-2 rounded-full" style={{ width: '32%' }}></div>
            </div>

            <div className="flex justify-between items-center text-sm mt-2">
              <span>Returning customers</span>
              <span>48%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div className="bg-green-500 h-2 rounded-full" style={{ width: '48%' }}></div>
            </div>

            <div className="flex justify-between items-center text-sm mt-2">
              <span>Loyal customers</span>
              <span>20%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div className="bg-purple-500 h-2 rounded-full" style={{ width: '20%' }}></div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
