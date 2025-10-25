
import React from 'react';
import { Card } from '@/components/ui/card';
import { ChartContainer } from '@/components/ui/chart';
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from 'recharts';

const monthlySalesData = [
  { month: 'Jan', revenue: 3200, orders: 124 },
  { month: 'Feb', revenue: 3500, orders: 145 },
  { month: 'Mar', revenue: 4100, orders: 162 },
  { month: 'Apr', revenue: 4300, orders: 175 },
  { month: 'May', revenue: 5200, orders: 198 },
  { month: 'Jun', revenue: 4800, orders: 186 },
  { month: 'Jul', revenue: 5100, orders: 195 },
  { month: 'Aug', revenue: 5800, orders: 215 },
  { month: 'Sep', revenue: 6200, orders: 236 },
  { month: 'Oct', revenue: 5600, orders: 210 },
  { month: 'Nov', revenue: 6800, orders: 245 },
  { month: 'Dec', revenue: 7400, orders: 268 }
];

const weekdaySalesData = [
  { day: 'Mon', revenue: 820, orders: 32 },
  { day: 'Tue', revenue: 920, orders: 38 },
  { day: 'Wed', revenue: 880, orders: 34 },
  { day: 'Thu', revenue: 1100, orders: 45 },
  { day: 'Fri', revenue: 1700, orders: 74 },
  { day: 'Sat', revenue: 1980, orders: 85 },
  { day: 'Sun', revenue: 1640, orders: 68 }
];

export const SalesMetricsChart: React.FC = () => {
  return (
    <div className="flex flex-col space-y-6 w-full">
      {/* Monthly Sales Chart */}
      <Card className="p-4 w-full">
        <h3 className="font-medium mb-4">Monthly Sales Analytics</h3>
        <div className="h-[350px]">
          <ChartContainer
            config={{
              revenue: { theme: { light: '#22c55e', dark: '#4ade80' } },
              orders: { theme: { light: '#3b82f6', dark: '#60a5fa' } },
            }}
          >
            <LineChart data={monthlySalesData} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" stroke="#888888" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis 
                yAxisId="left"
                stroke="#888888"
                fontSize={12}
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) => `$${value}`}
              />
              <YAxis
                yAxisId="right"
                orientation="right"
                stroke="#888888"
                fontSize={12}
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) => `${value}`}
              />
              <Tooltip />
              <Legend />
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="revenue"
                name="Revenue ($)"
                stroke="var(--color-revenue)"
                strokeWidth={2}
                dot={{ r: 4 }}
                activeDot={{ r: 6 }}
              />
              <Line
                yAxisId="right"
                type="monotone"
                dataKey="orders"
                name="Orders"
                stroke="var(--color-orders)"
                strokeWidth={2}
                dot={{ r: 4 }}
              />
            </LineChart>
          </ChartContainer>
        </div>
      </Card>
      
      {/* Weekly Charts Container */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 w-full">
        <Card className="p-4 w-full">
          <h3 className="font-medium mb-4">Weekly Revenue Trend</h3>
          <div className="h-[200px]">
            <ChartContainer
              config={{
                revenue: { theme: { light: '#22c55e', dark: '#4ade80' } },
              }}
            >
              <AreaChart data={weekdaySalesData} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="day" />
                <YAxis tickFormatter={(value) => `$${value}`} />
                <Tooltip />
                <Area type="monotone" dataKey="revenue" name="Revenue ($)" stroke="var(--color-revenue)" fill="var(--color-revenue)" fillOpacity={0.2} />
              </AreaChart>
            </ChartContainer>
          </div>
        </Card>
        
        <Card className="p-4 w-full">
          <h3 className="font-medium mb-4">Orders by Day of Week</h3>
          <div className="h-[200px]">
            <ChartContainer
              config={{
                orders: { theme: { light: '#3b82f6', dark: '#60a5fa' } },
              }}
            >
              <BarChart data={weekdaySalesData} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="day" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="orders" name="Orders" fill="var(--color-orders)" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ChartContainer>
          </div>
        </Card>
      </div>
    </div>
  );
};
