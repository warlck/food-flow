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
  pedestrian?: string;
  footway?: string;
  building?: string;
  residential?: string;
  amenity?: string;
  suburb?: string;
  neighbourhood?: string;
  quarter?: string;
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
  
  // If the query is a 6-digit Singapore postal code, format or query appropriately
  const trimmedQuery = query.trim();
  const isPostalCode = /^\d{6}$/.test(trimmedQuery);

  url.searchParams.set("q", isPostalCode ? `${trimmedQuery}, Singapore` : trimmedQuery);
  url.searchParams.set("format", "jsonv2");
  url.searchParams.set("addressdetails", "1");
  url.searchParams.set("countrycodes", "sg");
  url.searchParams.set("accept-language", "en");
  url.searchParams.set("limit", "5");

  const response = await fetch(url.toString(), {
    headers: { 
      Accept: "application/json",
      "Accept-Language": "en",
    },
  });

  if (!response.ok) {
    throw new Error("Address search failed");
  }

  const data: NominatimResult[] = await response.json();

  return data.map((item) => {
    const addr = item.address;
    const building = addr?.building || addr?.residential || addr?.amenity || "";
    const houseNum = addr?.house_number || "";
    const road = addr?.road || addr?.pedestrian || addr?.footway || "";
    const suburb = addr?.suburb || addr?.neighbourhood || addr?.quarter || "";

    const parts: string[] = [];
    if (building) parts.push(building);
    if (houseNum || road) parts.push([houseNum, road].filter(Boolean).join(" "));
    if (suburb && !parts.some((p) => p.includes(suburb))) parts.push(suburb);

    let street = parts.join(", ");
    if (!street) {
      street = item.display_name.split(",")[0] || item.display_name;
    }

    const postalCode = addr?.postcode || (isPostalCode ? trimmedQuery : "");
    const city = addr?.city || addr?.town || addr?.village || "Singapore";
    const state = addr?.state || (addr as unknown as Record<string, string> | undefined)?.region || (addr as unknown as Record<string, string> | undefined)?.county || "SG";

    return {
      displayName: item.display_name,
      latitude: parseFloat(item.lat),
      longitude: parseFloat(item.lon),
      address: {
        street,
        city,
        state,
        postalCode,
      },
    };
  });
}
