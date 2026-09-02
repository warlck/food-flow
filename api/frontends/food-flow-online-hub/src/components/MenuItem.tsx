import React, { useState } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { MenuItem as MenuItemType } from '@/types';
import { Plus, Tag, Layers, Puzzle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import MenuItemDialog from './MenuItemDialog';

interface MenuItemProps {
  item: MenuItemType;
  categoryItems?: MenuItemType[];
  onCartUpdate?: () => void;
}

const MenuItem: React.FC<MenuItemProps> = ({ item, categoryItems, onCartUpdate }) => {
  const [imageError, setImageError] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const isOrderable = item.orderable ?? item.available ?? false;

  const handlePrimaryClick = () => {
    setIsDialogOpen(true);
  };

  const handleOpenDialog = () => {
    setIsDialogOpen(true);
  };
  
  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    if (onCartUpdate) {
      onCartUpdate();
    }
  };

  const handleImageError = () => {
    setImageError(true);
  };

  // Generate a consistent fallback image based on the item's category
  const getFallbackImage = () => {
    const categoryImageMap: Record<string, string> = {
      'Appetizers': 'https://images.unsplash.com/photo-1546241072-48010ad2862c?auto=format&fit=crop&q=80',
      'Main Course': 'https://images.unsplash.com/photo-1574484284002-952d92456975?auto=format&fit=crop&q=80',
      'Desserts': 'https://images.unsplash.com/photo-1563729784474-d77dbb933a9e?auto=format&fit=crop&q=80',
      'Beverages': 'https://images.unsplash.com/photo-1544145945-f90425340c7e?auto=format&fit=crop&q=80',
      'Sides': 'https://images.unsplash.com/photo-1573080496219-bb080dd4f877?auto=format&fit=crop&q=80',
      'Pizza': 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?auto=format&fit=crop&q=80',
      'Burgers': 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&q=80',
      'Pasta': 'https://images.unsplash.com/photo-1473093226795-af9932fe5856?auto=format&fit=crop&q=80',
    };
    
    // Default fallback image if category doesn't match
    const defaultImage = 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&q=80';
    
    return categoryImageMap[item.category] || defaultImage;
  };

  const hasModifiers = item.modifierGroups && item.modifierGroups.length > 0;
  const hasAddons = item.addons && item.addons.length > 0;

  return (
    <>
      <Card 
        className={`h-full flex flex-col overflow-hidden transition-all hover:shadow-lg cursor-pointer ${
          !isOrderable ? 'opacity-75' : ''
        }`}
        onClick={handleOpenDialog}
      >
        <div className="relative w-full aspect-[4/3] overflow-hidden bg-muted">
          {!imageError ? (
            <img
              src={item.image}
              alt={item.name}
              className="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
              onError={handleImageError}
            />
          ) : (
            <img 
              src={getFallbackImage()}
              alt={item.name}
              className="w-full h-full object-cover opacity-90"
            />
          )}
          {!isOrderable && (
            <div className="absolute inset-0 bg-black/60 flex items-center justify-center">
              <span className="text-white font-bold text-base tracking-wide uppercase px-3 py-1 bg-red-600/80 rounded-full">
                Out of Stock
              </span>
            </div>
          )}
          <div className="absolute top-2 right-2">
            <Badge className="bg-food-primary text-white">{item.category}</Badge>
          </div>
          <div className="absolute bottom-2 left-2 flex flex-wrap gap-1">
            {hasModifiers && (
              <Badge variant="secondary" className="bg-white/90 text-gray-700 text-[10px] flex items-center gap-1">
                <Layers size={10} />
                Customizable
              </Badge>
            )}
            {hasAddons && (
              <Badge variant="secondary" className="bg-white/90 text-gray-700 text-[10px] flex items-center gap-1">
                <Puzzle size={10} />
                {item.addons!.length} extra{item.addons!.length > 1 ? 's' : ''}
              </Badge>
            )}
          </div>
        </div>
        
        <CardHeader className="pb-2">
          <div className="flex justify-between items-start">
            <CardTitle className="text-base font-bold line-clamp-2 text-gray-900">{item.name}</CardTitle>
            <span className="font-bold text-food-primary shrink-0 ml-2">${item.price.toFixed(2)}</span>
          </div>
          {item.tags && item.tags.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-1">
              {item.tags.map(tag => (
                <Badge key={tag} variant="secondary" className="text-[10px] flex items-center">
                  <Tag size={9} className="mr-1" />
                  {tag}
                </Badge>
              ))}
            </div>
          )}
        </CardHeader>
        
        <CardContent className="pb-2 flex-grow">
          <CardDescription className="text-xs text-gray-600 line-clamp-2">{item.description}</CardDescription>
        </CardContent>
        
        <CardFooter className="pt-2">
          <Button
            onClick={(e) => {
              e.stopPropagation();
              handlePrimaryClick();
            }}
            disabled={!isOrderable}
            size="sm"
            className="w-full bg-food-primary hover:bg-food-accent text-white font-semibold text-xs"
          >
            <Plus className="mr-1 h-3.5 w-3.5" />
            {hasModifiers || hasAddons ? 'Customize' : 'Add to Cart'}
          </Button>
        </CardFooter>
      </Card>

      <MenuItemDialog
        item={item}
        categoryItems={categoryItems}
        isOpen={isDialogOpen}
        onClose={handleCloseDialog}
      />
    </>
  );
};

export default MenuItem;
