
import React from 'react';
import { Button } from '@/components/ui/button';
import { SidebarGroup, SidebarGroupContent, SidebarGroupLabel } from '@/components/ui/sidebar';
import { List } from 'lucide-react';

interface CategorySidebarProps {
  categories: string[];
  selectedCategory: string;
  onCategorySelect: (category: string) => void;
}

const CategorySidebar: React.FC<CategorySidebarProps> = ({ 
  categories, 
  selectedCategory, 
  onCategorySelect 
}) => {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>
        <List className="mr-2" />
        Categories
      </SidebarGroupLabel>
      <SidebarGroupContent>
        <div className="space-y-2">
          {categories.map((category) => (
            <Button
              key={category}
              variant={selectedCategory === category ? "default" : "outline"}
              onClick={() => onCategorySelect(category)}
              className={`
                w-full justify-start
                ${selectedCategory === category 
                  ? 'bg-food-primary hover:bg-food-accent text-white' 
                  : 'text-gray-700 hover:bg-gray-100'
                } 
              `}
              size="sm"
            >
              {category}
            </Button>
          ))}
        </div>
      </SidebarGroupContent>
    </SidebarGroup>
  );
};

export default CategorySidebar;
