
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

  // Generate Google Maps URL for directions
  const getGoogleMapsUrl = (): string | null => {
    if (restaurant.latitude != null && restaurant.longitude != null) {
      return `https://www.google.com/maps/dir/?api=1&destination=${restaurant.latitude},${restaurant.longitude}`;
    }
    if (restaurant.address && restaurant.address.trim().length > 0) {
      return `https://www.google.com/maps/dir/?api=1&destination=${encodeURIComponent(restaurant.address)}`;
    }
    return null;
  };

  const googleMapsUrl = getGoogleMapsUrl();

  return (
    <div className="w-full">
      {/* Cover Image */}
      <div className="w-full h-64 relative">
        <img
          src={restaurant.coverImage}
          alt={restaurant.name}
          className="w-full h-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent flex items-end">
          <div className="container mx-auto px-4 pb-6 flex items-center">
            <div className="bg-white p-2 rounded-full mr-4 shadow-lg">
              <img
                src={restaurant.logo || "https://via.placeholder.com/80"}
                alt={`${restaurant.name} logo`}
                className="w-20 h-20 rounded-full"
              />
            </div>
            <div className="text-white">
              <h1 className="text-3xl font-bold drop-shadow-lg">{restaurant.name}</h1>
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
                    {isOpen ? (
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
