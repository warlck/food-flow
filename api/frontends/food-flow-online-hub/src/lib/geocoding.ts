// Geocoding helper backed by Nominatim (OpenStreetMap). Free to use without
// an API key; see https://operations.osmfoundation.org/policies/nominatim/.

export interface GeocodedAddress {
  street: string;
  city: string;
  state: string;
  postalCode: string;
}

export interface GeocodingResult {
  displayName: string;
  latitude: number;
  longitude: number;
  address: GeocodedAddress;
}

interface NominatimAddress {
  house_number?: string;
  road?: string;
  city?: string;
  town?: string;
  village?: string;
  state?: string;
  postcode?: string;
}

interface NominatimResult {
  display_name: string;
  lat: string;
  lon: string;
  address?: NominatimAddress;
}

export async function searchAddress(query: string): Promise<GeocodingResult[]> {
  const url = new URL("https://nominatim.openstreetmap.org/search");
  url.searchParams.set("q", query);
  url.searchParams.set("format", "jsonv2");
  url.searchParams.set("addressdetails", "1");
  url.searchParams.set("limit", "5");

  const response = await fetch(url.toString(), {
    headers: { Accept: "application/json" },
  });

  if (!response.ok) {
    throw new Error("Address search failed");
  }

  const data: NominatimResult[] = await response.json();

  return data.map((item) => ({
    displayName: item.display_name,
    latitude: parseFloat(item.lat),
    longitude: parseFloat(item.lon),
    address: {
      street: [item.address?.house_number, item.address?.road].filter(Boolean).join(" "),
      city: item.address?.city ?? item.address?.town ?? item.address?.village ?? "",
      state: item.address?.state ?? "",
      postalCode: item.address?.postcode ?? "",
    },
  }));
}
