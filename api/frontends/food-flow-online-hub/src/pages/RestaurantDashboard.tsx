
import React, { useState } from 'react';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { mockMenuItems, mockOrders, mockCustomers } from '@/data/mockData';
import { Plus, Trash, Edit, Search, Clock, ArrowUpDown, Filter, Users, CreditCard, Receipt, TrendingUp } from 'lucide-react';
import { toast } from '@/components/ui/use-toast';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { SalesMetricsChart } from '@/components/dashboard/SalesMetricsChart';
import { CustomerMetrics } from '@/components/dashboard/CustomerMetrics';
import { RevenueStats } from '@/components/dashboard/RevenueStats';
import { CustomerTable } from '@/components/dashboard/CustomerTable';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { SidebarProvider, Sidebar, SidebarContent, SidebarMenu, SidebarMenuItem, SidebarMenuButton, SidebarInset, SidebarGroup, SidebarGroupContent } from '@/components/ui/sidebar';

const RestaurantDashboard: React.FC = () => {
  const [menuItems, setMenuItems] = useState(mockMenuItems);
  const [orders, setOrders] = useState(mockOrders);
  const [customers, setCustomers] = useState(mockCustomers);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [customerSearchQuery, setCustomerSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState('analytics');

  const handleToggleAvailability = (id: string) => {
    setMenuItems(items => 
      items.map(item => 
        item.id === id ? { ...item, available: !item.available } : item
      )
    );
    
    const item = menuItems.find(item => item.id === id);
    if (item) {
      toast({
        description: `${item.name} is now ${!item.available ? 'available' : 'unavailable'}`,
      });
    }
  };

  const handleUpdateStatus = (orderId: string, newStatus: string) => {
    setOrders(prevOrders => 
      prevOrders.map(order => 
        order.id === orderId 
          ? { ...order, status: newStatus as any } 
          : order
      )
    );
    
    toast({
      description: `Order #${orderId.slice(-3)} status updated to ${newStatus}`,
    });
  };

  const filteredOrders = orders.filter(order => {
    if (statusFilter !== 'all') {
      return order.status === statusFilter;
    }
    return true;
  });

  const filteredMenuItems = menuItems.filter(item => 
    item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const filteredCustomers = customers.filter(customer =>
    customer.name.toLowerCase().includes(customerSearchQuery.toLowerCase()) ||
    customer.email.toLowerCase().includes(customerSearchQuery.toLowerCase())
  );

  const renderContent = () => {
    switch (activeTab) {
      case 'analytics':
        return (
          <>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              <RevenueStats />
              
              <CustomerMetrics />
              
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg font-medium">
                    <div className="flex items-center">
                      <Receipt className="mr-2 h-5 w-5 text-muted-foreground" />
                      Recent Orders
                    </div>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {orders.slice(0, 5).map(order => (
                      <div key={order.id} className="flex items-center justify-between border-b pb-2">
                        <div>
                          <div className="font-medium">Order #{order.id.slice(-3)}</div>
                          <div className="text-sm text-muted-foreground">
                            {order.createdAt.toLocaleDateString()}
                          </div>
                        </div>
                        <div className="font-semibold">${order.totalAmount.toFixed(2)}</div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
            
            <div className="grid gap-6 mt-6">
              <Card className="col-span-1">
                <CardHeader>
                  <CardTitle>
                    <div className="flex items-center">
                      <TrendingUp className="mr-2 h-5 w-5 text-muted-foreground" />
                      Sales Analytics
                    </div>
                  </CardTitle>
                  <CardDescription>Monthly revenue and order trends</CardDescription>
                </CardHeader>
                <CardContent>
                  <SalesMetricsChart />
                </CardContent>
              </Card>
            </div>

            <div className="mt-6">
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle>
                        <div className="flex items-center">
                          <Users className="mr-2 h-5 w-5 text-muted-foreground" />
                          Customer Information
                        </div>
                      </CardTitle>
                      <CardDescription>Details of customers for marketing purposes</CardDescription>
                    </div>
                    <div className="relative w-full max-w-sm">
                      <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                      <Input
                        type="search"
                        placeholder="Search customers..."
                        value={customerSearchQuery}
                        onChange={(e) => setCustomerSearchQuery(e.target.value)}
                        className="pl-10"
                      />
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <CustomerTable customers={filteredCustomers} />
                </CardContent>
              </Card>
            </div>
          </>
        );
      case 'orders':
        return (
          <div className="grid gap-6 mt-6">
            <Card className="col-span-1">
              <CardHeader>
                <CardTitle>
                  <div className="flex items-center">
                    <Receipt className="mr-2 h-5 w-5 text-muted-foreground" />
                    Order Management
                  </div>
                </CardTitle>
                <CardDescription>View and manage customer orders</CardDescription>
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mt-4">
                  <div className="relative w-full max-w-sm">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                    <Input
                      type="search"
                      placeholder="Search orders..."
                      className="pl-10"
                    />
                  </div>
                  <div className="flex gap-2">
                    <Select value={statusFilter} onValueChange={setStatusFilter}>
                      <SelectTrigger className="w-[180px]">
                        <SelectValue placeholder="Filter by status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="all">All Orders</SelectItem>
                        <SelectItem value="pending">Pending</SelectItem>
                        <SelectItem value="preparing">Preparing</SelectItem>
                        <SelectItem value="ready">Ready for Pickup</SelectItem>
                        <SelectItem value="delivered">Delivered</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="overflow-x-auto">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Order ID</TableHead>
                        <TableHead>Customer</TableHead>
                        <TableHead>Items</TableHead>
                        <TableHead>Total</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredOrders.map(order => {
                        const customer = mockCustomers.find(c => c.id === order.customerId);
                        return (
                          <TableRow key={order.id}>
                            <TableCell>#{order.id.slice(-3)}</TableCell>
                            <TableCell>{customer?.name || 'Unknown Customer'}</TableCell>
                            <TableCell>{order.items.length} items</TableCell>
                            <TableCell>${order.totalAmount.toFixed(2)}</TableCell>
                            <TableCell>
                              <Select
                                value={order.status}
                                onValueChange={(value) => handleUpdateStatus(order.id, value)}
                              >
                                <SelectTrigger className="w-[140px]">
                                  <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="pending">Pending</SelectItem>
                                  <SelectItem value="preparing">Preparing</SelectItem>
                                  <SelectItem value="ready">Ready</SelectItem>
                                  <SelectItem value="delivered">Delivered</SelectItem>
                                </SelectContent>
                              </Select>
                            </TableCell>
                            <TableCell>
                              <Button variant="outline" size="sm">
                                View Details
                              </Button>
                            </TableCell>
                          </TableRow>
                        );
                      })}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          </div>
        );
      case 'menu':
        return (
          <div className="grid gap-6 mt-6">
            <Card className="col-span-1">
              <CardHeader>
                <CardTitle>
                  <div className="flex items-center">
                    <CreditCard className="mr-2 h-5 w-5 text-muted-foreground" />
                    Menu Management
                  </div>
                </CardTitle>
                <CardDescription>Add, edit, or remove menu items</CardDescription>
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mt-4">
                  <div className="relative w-full max-w-sm">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                    <Input
                      type="search"
                      placeholder="Search menu items..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-10"
                    />
                  </div>
                  <Button className="bg-food-primary hover:bg-food-accent">
                    <Plus className="mr-2 h-4 w-4" /> Add New Item
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="overflow-x-auto">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Image</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Category</TableHead>
                        <TableHead>Price</TableHead>
                        <TableHead>Available</TableHead>
                        <TableHead>Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredMenuItems.map(item => (
                        <TableRow key={item.id} className="border-b hover:bg-gray-50">
                          <TableCell>
                            <img 
                              src={item.image} 
                              alt={item.name} 
                              className="w-12 h-12 object-cover rounded"
                            />
                          </TableCell>
                          <TableCell className="font-medium">{item.name}</TableCell>
                          <TableCell>{item.category}</TableCell>
                          <TableCell>${item.price.toFixed(2)}</TableCell>
                          <TableCell className="text-center">
                            <Switch
                              checked={item.available}
                              onCheckedChange={() => handleToggleAvailability(item.id)}
                            />
                          </TableCell>
                          <TableCell className="text-right">
                            <div className="flex justify-end gap-2">
                              <Button variant="outline" size="sm">
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button variant="outline" size="sm" className="text-red-500 hover:text-red-700 hover:bg-red-50">
                                <Trash className="h-4 w-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <Layout>
      <div className="container mx-auto py-8">
        <h1 className="text-3xl font-bold mb-6 px-4">Restaurant Dashboard</h1>

        <SidebarProvider>
          <div className="flex min-h-screen w-full">
            <Sidebar variant="floating" collapsible="icon" className="sticky top-16 h-[calc(100vh-4rem)] z-10">
              <SidebarContent>
                <SidebarGroup>
                  <SidebarGroupContent>
                    <SidebarMenu>
                      <SidebarMenuItem>
                        <SidebarMenuButton 
                          isActive={activeTab === 'analytics'} 
                          onClick={() => setActiveTab('analytics')}
                          tooltip="Analytics"
                          size="lg"
                          className="text-base font-medium hover:bg-food-primary/10 transition-colors hover:scale-105 transform duration-200 shadow-sm hover:shadow-md"
                        >
                          <TrendingUp className="mr-2 h-5 w-5 text-food-primary" />
                          <span>Analytics</span>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                      <SidebarMenuItem>
                        <SidebarMenuButton 
                          isActive={activeTab === 'orders'} 
                          onClick={() => setActiveTab('orders')}
                          tooltip="Orders"
                          size="lg"
                          className="text-base font-medium hover:bg-food-primary/10 transition-colors hover:scale-105 transform duration-200 shadow-sm hover:shadow-md"
                        >
                          <Receipt className="mr-2 h-5 w-5 text-food-primary" />
                          <span>Orders</span>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                      <SidebarMenuItem>
                        <SidebarMenuButton 
                          isActive={activeTab === 'menu'} 
                          onClick={() => setActiveTab('menu')}
                          tooltip="Menu"
                          size="lg" 
                          className="text-base font-medium hover:bg-food-primary/10 transition-colors hover:scale-105 transform duration-200 shadow-sm hover:shadow-md"
                        >
                          <CreditCard className="mr-2 h-5 w-5 text-food-primary" />
                          <span>Menu</span>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    </SidebarMenu>
                  </SidebarGroupContent>
                </SidebarGroup>
              </SidebarContent>
            </Sidebar>
            <SidebarInset className="flex-1 px-4 max-w-[1200px] overflow-y-auto">
              {renderContent()}
            </SidebarInset>
          </div>
        </SidebarProvider>
      </div>
    </Layout>
  );
};

export default RestaurantDashboard;
