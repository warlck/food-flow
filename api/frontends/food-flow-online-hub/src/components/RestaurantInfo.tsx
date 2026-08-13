
import React from 'react';
import { Restaurant } from '@/types';
import { Clock, MapPin, Phone, Mail, Star, ExternalLink } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';

interface RestaurantInfoProps {
  restaurant: Restaurant;
}

const RestaurantInfo: React.FC<RestaurantInfoProps> = ({ restaurant }) => {
  const today = new Date().toLocaleDateString('en-US', { weekday: 'long' }).toLowerCase();
  const formattedToday = today as keyof typeof restaurant.openingHours;
  const isOpen = restaurant.openingHours[formattedToday] !== undefined;

  // Generate Google Maps URL for directions using restaurant name and address
  const getGoogleMapsUrl = (): string | null => {
    const hasName = restaurant.name && restaurant.name.trim().length > 0;
    const hasAddress = restaurant.address && restaurant.address.trim().length > 0;

    if (hasName && hasAddress) {
      const destinationQuery = `${restaurant.name}, ${restaurant.address}`;
      return `https://www.google.com/maps/dir/?api=1&destination=${encodeURIComponent(destinationQuery)}`;
    }
    if (hasAddress) {
      return `https://www.google.com/maps/dir/?api=1&destination=${encodeURIComponent(restaurant.address)}`;
    }
    if (hasName) {
      return `https://www.google.com/maps/dir/?api=1&destination=${encodeURIComponent(restaurant.name)}`;
    }
    return null;
  };

  const googleMapsUrl = getGoogleMapsUrl();

  return (
    <div className="w-full">
      {/* Paused Alert Banner */}
      {restaurant.enabled === false && (
        <div className="bg-red-50 border-b border-red-200 px-4 py-3 text-sm text-red-900">
          <div className="container mx-auto flex items-center justify-between gap-3">
            <div className="flex items-center gap-2">
              <span className="h-2.5 w-2.5 rounded-full bg-red-600 animate-pulse" />
              <span><strong className="font-bold">Storefront Paused:</strong> {restaurant.name} is currently paused and not taking new online orders at this time.</span>
            </div>
            <span className="shrink-0 font-bold uppercase tracking-wider text-[10px] bg-red-200 text-red-900 px-2.5 py-1 rounded-full">Paused</span>
          </div>
        </div>
      )}

      {/* Cover Image */}
      <div className="w-full h-64 relative">
        <img
          src={restaurant.coverImage}
          alt={restaurant.name}
          className={`w-full h-full object-cover ${restaurant.enabled === false ? 'grayscale-[30%]' : ''}`}
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent flex items-end">
          <div className="container mx-auto px-4 pb-6 flex items-center">
            {restaurant.logo && (
              <div className="bg-white p-2 rounded-full mr-4 shadow-lg">
                <img
                  src={restaurant.logo}
                  alt={`${restaurant.name} logo`}
                  className="w-20 h-20 rounded-full"
                />
              </div>
            )}
            <div className="text-white">
              <div className="flex items-center gap-3">
                <h1 className="text-3xl font-bold drop-shadow-lg">{restaurant.name}</h1>
                {restaurant.enabled === false && (
                  <span className="bg-red-600 text-white text-xs font-bold px-2.5 py-0.5 rounded-full shadow drop-shadow-md uppercase tracking-wider">
                    Paused
                  </span>
                )}
              </div>
              <div className="flex items-center mt-2">
                <Star className="h-5 w-5 text-food-secondary fill-food-secondary mr-1 drop-shadow" />
                <span className="font-medium drop-shadow">{restaurant.rating}</span>
                <span className="mx-2 drop-shadow">•</span>
                <span className="drop-shadow">{restaurant.description}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Restaurant Info Cards */}
      <div className="container mx-auto px-4 -mt-6 mb-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Hours */}
          <Card className="bg-white/90 backdrop-blur-sm">
            <CardContent className="pt-6">
              <div className="flex items-start">
                <Clock className="h-5 w-5 text-food-primary mr-2 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold">Hours</h3>
                  <p className="text-sm text-gray-600 mb-2">
                    {restaurant.enabled === false ? (
                      <span className="text-red-600 font-bold">Paused (Storefront Closed)</span>
                    ) : isOpen ? (
                      <>
                        <span className="text-food-success font-medium">Open Today:</span>{" "}
                        {restaurant.openingHours[formattedToday]?.open} - {restaurant.openingHours[formattedToday]?.close}
                      </>
                    ) : (
                      <span className="text-red-500 font-medium">Closed Today</span>
                    )}
                  </p>
                  <button className="text-food-primary text-sm">See all hours</button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Location */}
          <Card className="bg-white/90 backdrop-blur-sm">
            <CardContent className="pt-6">
              <div className="flex items-start">
                <MapPin className="h-5 w-5 text-food-primary mr-2 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold">Location</h3>
                  <p className="text-sm text-gray-600">
                    {restaurant.address || "Address unavailable"}
                  </p>
                  {googleMapsUrl ? (
                    <a
                      href={googleMapsUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center text-food-primary text-sm font-medium mt-2 hover:underline focus:outline-none focus:ring-2 focus:ring-food-primary rounded"
                    >
                      Get directions
                      <ExternalLink className="h-3.5 w-3.5 ml-1" />
                    </a>
                  ) : (
                    <span className="text-gray-400 text-sm mt-2 inline-block">
                      Directions unavailable
                    </span>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Contact */}
          <Card className="bg-white/90 backdrop-blur-sm">
            <CardContent className="pt-6">
              <div className="flex items-start">
                <Phone className="h-5 w-5 text-food-primary mr-2 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold">Contact</h3>
                  <p className="text-sm text-gray-600">{restaurant.phone}</p>
                  <div className="flex items-center mt-1">
                    <Mail className="h-4 w-4 text-gray-500 mr-1" />
                    <span className="text-sm text-gray-600">{restaurant.email}</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

    </div>
  );
};

export default RestaurantInfo;
