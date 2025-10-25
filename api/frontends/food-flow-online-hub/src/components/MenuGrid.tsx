import React, { useState } from 'react';
import { SidebarProvider, Sidebar, SidebarContent, SidebarInset } from '@/components/ui/sidebar';
import MenuItem from './MenuItem';
import CategorySidebar from './CategorySidebar';
import { MenuItem as MenuItemType } from '@/types';
import { Input } from '@/components/ui/input';
import { Search } from 'lucide-react';

interface MenuGridProps {
  items: MenuItemType[];
  categories: string[];
  onCartUpdate?: () => void;
}

const MenuGrid: React.FC<MenuGridProps> = ({ items, categories, onCartUpdate }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');

  const filteredItems = items.filter((item) => {
    const matchesSearch = item.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          item.description.toLowerCase().includes(searchQuery.toLowerCase());
    
    const matchesCategory = selectedCategory === 'All' || item.category === selectedCategory;
    
    return matchesSearch && matchesCategory;
  });

  return (
    <SidebarProvider>
      <div className="flex min-h-screen">
        <Sidebar variant="floating" collapsible="icon" className="sticky top-16 h-[calc(100vh-4rem)]">
          <CategorySidebar 
            categories={categories} 
            selectedCategory={selectedCategory} 
            onCategorySelect={setSelectedCategory} 
          />
        </Sidebar>
        <SidebarInset className="flex-1">
          <div className="w-full px-0 sm:px-2 py-4">
            <div className="mb-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
                <Input
                  type="search"
                  placeholder="Search menu items..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>

            {filteredItems.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
                {filteredItems.map((item) => (
                  <MenuItem key={item.id} item={item} onCartUpdate={onCartUpdate} />
                ))}
              </div>
            ) : (
              <div className="text-center py-10">
                <p className="text-gray-500 mb-2">No menu items found.</p>
                <button
                  className="text-food-primary hover:underline"
                  onClick={() => {
                    setSearchQuery('');
                    setSelectedCategory('All');
                  }}
                >
                  Clear filters
                </button>
              </div>
            )}
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
};

export default MenuGrid;
