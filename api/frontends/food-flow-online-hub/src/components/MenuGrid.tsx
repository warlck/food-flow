import React, { useMemo, useState } from 'react';
import { SidebarProvider, Sidebar, SidebarInset } from '@/components/ui/sidebar';
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

  const categoriesWithAll = useMemo(() => ['All', ...categories], [categories]);

  // Group items by category.
  // IMPORTANT: backend returns items within each category rank-ordered (ranked first, then by price).
  const itemsByCategory = useMemo(() => {
    const map = new Map<string, MenuItemType[]>();

    items.forEach((mi) => {
      const existing = map.get(mi.category);
      if (existing) {
        existing.push(mi);
      } else {
        map.set(mi.category, [mi]);
      }
    });

    return map;
  }, [items]);

  const displayedItems = useMemo(() => {
    const query = searchQuery.trim().toLowerCase();

    const matchesQuery = (mi: MenuItemType) => {
      if (!query) return true;
      return (
        mi.name.toLowerCase().includes(query) ||
        mi.description.toLowerCase().includes(query)
      );
    };

    // Spec §10.2: every menu item is an individual card. `items` arrives in
    // backend order (categories in order, then rank/price within a category),
    // so filtering preserves it.
    if (selectedCategory !== 'All') {
      return (itemsByCategory.get(selectedCategory) ?? []).filter(matchesQuery);
    }

    return items.filter(matchesQuery);
  }, [items, itemsByCategory, searchQuery, selectedCategory]);

  return (
    <SidebarProvider>
      <div className="flex min-h-screen">
        <Sidebar variant="floating" collapsible="icon" className="sticky top-20 h-[calc(100vh-5rem)]">
          <CategorySidebar
            categories={categoriesWithAll}
            selectedCategory={selectedCategory}
            onCategorySelect={setSelectedCategory}
          />
        </Sidebar>
        <SidebarInset className="flex-1">
          <div className="w-full px-0 sm:px-2 pt-2 pb-4">
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

            {displayedItems.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
                {displayedItems.map((item) => (
                  <MenuItem
                    key={item.id}
                    item={item}
                    onCartUpdate={onCartUpdate}
                  />
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
