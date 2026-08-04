-- ============================================================================
-- Users
-- ============================================================================
INSERT INTO users (user_id, name, email, roles, password_hash, department, enabled, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- Restaurants
-- ============================================================================
INSERT INTO restaurants (restaurant_id, name, description, address, phone, email, image_url, enabled, latitude, longitude, max_delivery_distance_km, date_created, date_updated) VALUES
	('a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Donergy', 'Authentic Turkish Kebab & Pide Restaurant', '9 Raffles Boulevard, #01-91B, Millenia Walk, Singapore 039596', '+65 6333 0785', 'info@donergy.sg', 'https://www.donergy.sg/Content/images/logo.png', true, 1.29305, 103.86020, 10, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'A1 Kebab', 'Signature Kebabs & 100% Plant-Based Falafel', '100 Beach Road, #01-12, Singapore 189702', '+65 6789 0123', 'info@a1kebabs.com', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/875772946_10ChickenRoll1.jpg', true, 1.29665, 103.85630, 10, '2025-10-27 00:00:00', '2025-10-27 00:00:00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- Categories
-- ============================================================================
INSERT INTO categories (category_id, name, description, restaurant_id, enabled, date_created, date_updated) VALUES
	-- Donergy Categories
	('c1000000-0000-0000-0000-000000000001', 'Kebab Roll', 'Delicious kebab wrapped in fresh flatbread', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000002', 'Kebab Tombik', 'Thick pita bread filled with kebab and fresh vegetables', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000003', 'Kebab Rice', 'Kebab served over aromatic rice', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000004', 'Kebab Salad', 'Fresh salad topped with grilled kebab', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000005', 'Kebab with Chips', 'Kebab served with crispy french fries', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000006', 'Iskender', 'Traditional Turkish dish with doner on pita bread with tomato sauce and yogurt', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000007', 'Pide and Lahmajun', 'Turkish flatbread pizza with various toppings', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000008', 'Vegetarian Options', 'Delicious meat-free dishes', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000009', 'Side Dishes', 'Appetizers and accompaniments', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000010', 'Dessert', 'Traditional Turkish sweets', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c1000000-0000-0000-0000-000000000011', 'Drinks', 'Refreshing beverages', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	-- A1 Categories
	('c2000000-0000-0000-0000-000000000001', 'Signature Kebabs', 'Fresh ingredients, bold flavours', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c2000000-0000-0000-0000-000000000002', 'Kebab Bowl', 'Choose your base (150g chicken)', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c2000000-0000-0000-0000-000000000003', 'Falafel', '100% Plant-Based goodness', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c2000000-0000-0000-0000-000000000004', 'Sides', 'Perfect accompaniments', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('c2000000-0000-0000-0000-000000000005', 'Drinks', 'Refreshing beverages', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- Menu Items
-- ============================================================================
INSERT INTO menu_items (menu_item_id, name, description, price, category_id, restaurant_id, image_url, available, date_created, date_updated) VALUES
	-- KEBAB ROLL
	('a1b2c3d4-0001-0000-0000-000000000001', 'Chicken Kebab Roll', 'Tender chicken kebab wrapped in fresh flatbread', 11.00, 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/875772946_10ChickenRoll1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000002', 'Beef Kebab Roll', 'Juicy beef kebab wrapped in fresh flatbread', 12.00, 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/28602426_11BeefRoll.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000003', 'Mix Kebab Roll', 'Combination of chicken and beef kebab', 13.50, 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/799344811_12MixRoll1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000004', 'Falafel Chicken Roll', 'Crispy falafel with chicken kebab', 13.50, 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1704566297_30FallafelChicken1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000005', 'Falafel Beef Roll', 'Crispy falafel with beef kebab', 14.50, 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/214602589_31FallafelBeef.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000006', 'Falafel Mix Roll', 'Crispy falafel with mixed kebab', 15.50, 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/752731853_32FallafelMixRoll.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- KEBAB TOMBIK
	('a1b2c3d4-0001-0000-0000-000000000007', 'Chicken Kebab Tombik', 'Chicken kebab in thick pita bread', 11.00, 'c1000000-0000-0000-0000-000000000002', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/548042839_20ChickenTombik1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000008', 'Beef Kebab Tombik', 'Beef kebab in thick pita bread', 12.00, 'c1000000-0000-0000-0000-000000000002', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/559999964_21BeefTombik1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000009', 'Mix Kebab Tombik', 'Mixed kebab in thick pita bread', 13.50, 'c1000000-0000-0000-0000-000000000002', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1994091345_22MixTombik.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- KEBAB RICE
	('a1b2c3d4-0001-0000-0000-000000000010', 'Chicken Kebab Rice', 'Grilled chicken kebab served over rice', 14.50, 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/53393089_40ChickenRice.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000011', 'Beef Kebab Rice', 'Grilled beef kebab served over rice', 16.50, 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1779406103_41BeefRice.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000012', 'Mix Kebab Rice', 'Grilled mixed kebab served over rice', 17.50, 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/299123933_42MixRice.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- KEBAB SALAD
	('a1b2c3d4-0001-0000-0000-000000000013', 'Chicken Kebab Salad', 'Fresh salad topped with grilled chicken kebab', 14.50, 'c1000000-0000-0000-0000-000000000004', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1855254249_60ChickenSalad.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000014', 'Beef Kebab Salad', 'Fresh salad topped with grilled beef kebab', 16.50, 'c1000000-0000-0000-0000-000000000004', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/398837593_61BeefSalad.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000015', 'Mix Kebab Salad', 'Fresh salad topped with grilled mixed kebab', 17.50, 'c1000000-0000-0000-0000-000000000004', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/360612110_62MixSalad.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- KEBAB WITH CHIPS
	('a1b2c3d4-0001-0000-0000-000000000016', 'Chicken Kebab With Fries', 'Chicken kebab served with crispy french fries', 16.50, 'c1000000-0000-0000-0000-000000000005', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1662471778_50ChickenFries1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000017', 'Beef Kebab With Fries', 'Beef kebab served with crispy french fries', 17.50, 'c1000000-0000-0000-0000-000000000005', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1299318997_51BeefFries1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000018', 'Mix Kebab With Fries', 'Mixed kebab served with crispy french fries', 19.00, 'c1000000-0000-0000-0000-000000000005', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/262990838_52MixFries1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- ISKENDER
	('a1b2c3d4-0001-0000-0000-000000000019', 'Chicken Iskender', 'Traditional Turkish dish with chicken doner on pita bread, tomato sauce and yogurt', 17.50, 'c1000000-0000-0000-0000-000000000006', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/779992212_7IskenderChicken.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000020', 'Beef Iskender', 'Traditional Turkish dish with beef doner on pita bread, tomato sauce and yogurt', 19.00, 'c1000000-0000-0000-0000-000000000006', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/2005486030_71IskenderBeef.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000021', 'Mix Iskender', 'Traditional Turkish dish with mixed doner on pita bread, tomato sauce and yogurt', 20.00, 'c1000000-0000-0000-0000-000000000006', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/772495155_72IskenderMix.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- PIDE / LAHMAJUN
	('a1b2c3d4-0001-0000-0000-000000000022', 'Chicken Pide', 'Turkish flatbread topped with chicken', 17.50, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1183812739_115ChickenPide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000023', 'Beef Pide', 'Turkish flatbread topped with beef', 18.50, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1660914797_18BeefPide12.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000024', 'Mix Pide', 'Turkish flatbread topped with mixed meat', 20.00, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/185597017_116MixMeatPide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000025', 'Sujuk Pide', 'Turkish flatbread topped with spicy Turkish sausage', 18.50, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1712998977_114SujukPide.jpg', false, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000026', 'Beef Pepperoni Pide', 'Turkish flatbread topped with beef pepperoni', 18.50, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1708338611_1712998977114SujukPide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000027', 'Minced Lamb Pide', 'Turkish flatbread topped with minced lamb', 19.00, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/202440444_196MincedLambPide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000028', 'Lahmajun', 'Thin Turkish pizza with minced meat', 12.00, 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1382900566_8Lahmajun.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- VEGETARIAN OPTIONS
	('a1b2c3d4-0001-0000-0000-000000000029', 'Cheese Pide', 'Turkish flatbread topped with cheese', 15.50, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1490446861_110CheesePide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000030', 'Mushroom Pide', 'Turkish flatbread topped with mushrooms', 16.50, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1459440086_328109185113MashroomPide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000031', 'Spinach Pide', 'Turkish flatbread topped with spinach', 16.50, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/400124841_1445664858112SpinachPide.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000032', 'Mix Vegetarian Pide', 'Turkish flatbread topped with mixed vegetables', 18.50, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/792560985_2040978936111MixVegetarianPide1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000033', 'Falafel Plate', 'Crispy chickpea fritters served with salad and tahini sauce', 15.50, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/574589317_130FalafelPlate1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000034', 'Falafel Roll', 'Crispy falafel wrapped in flatbread (vegan)', 11.00, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/89806033_120FallafelRollvegan1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000035', 'Falafel Tombik', 'Crispy falafel in thick pita bread', 12.00, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1375458996_falafeltombik.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000036', 'Lentil Soup', 'Traditional Turkish lentil soup', 6.00, 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/866624483_120LentilSoup.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- SIDE DISHES
	('a1b2c3d4-0001-0000-0000-000000000037', 'Babaganoush - Served with 1 Pita Bread', 'Smoky eggplant dip served with pita bread', 11.00, 'c1000000-0000-0000-0000-000000000009', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1914765406_16Babaganush.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000038', 'Hummus - Served with 1 pita bread', 'Creamy chickpea dip served with pita bread', 10.00, 'c1000000-0000-0000-0000-000000000009', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/364237033_140Hummus.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000039', 'Mix Mezze', 'Assorted dips and appetizers', 15.50, 'c1000000-0000-0000-0000-000000000009', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1043193496_mixmazzeh.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000040', 'Pita Bread', 'Fresh Turkish pita bread', 3.50, 'c1000000-0000-0000-0000-000000000009', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1577079933_Pitaaaaa.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000041', 'French Fries', 'Crispy golden french fries', 5.00, 'c1000000-0000-0000-0000-000000000009', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/862897504_150FrenchFries.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000042', 'Sesame Bread', 'Traditional Turkish sesame bread', 6.50, 'c1000000-0000-0000-0000-000000000009', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1655527973_sesamebrad.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- DESSERT
	('a1b2c3d4-0001-0000-0000-000000000043', 'Baklawa', 'Sweet Turkish pastry with nuts and honey', 4.00, 'c1000000-0000-0000-0000-000000000010', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/109642321_180Baklawa1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000044', 'Kunefe', 'Shredded filo pastry with cheese, served warm with syrup', 14.50, 'c1000000-0000-0000-0000-000000000010', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1681489883_200Kunefe1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- DRINKS
	('a1b2c3d4-0001-0000-0000-000000000045', 'Coffee', 'Traditional Turkish coffee', 5.00, 'c1000000-0000-0000-0000-000000000011', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/632926238_170TurkishCoffee1.jpg', false, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000046', 'Ayran Sweet', 'Sweet yogurt drink', 4.50, 'c1000000-0000-0000-0000-000000000011', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/262280452_1006442618200TurkishAyran.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000047', 'Ayran Plain', 'Salted yogurt drink', 4.50, 'c1000000-0000-0000-0000-000000000011', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/904291627_AyranPlain.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000048', 'Tea', 'Traditional Turkish tea', 3.00, 'c1000000-0000-0000-0000-000000000011', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/786264153_220TurkishTea.jpg', false, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000049', 'Mineral Water', 'Bottled mineral water', 1.50, 'c1000000-0000-0000-0000-000000000011', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1484899931_42MineralWater1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a1b2c3d4-0001-0000-0000-000000000050', 'Soft Drinks', 'Assorted carbonated beverages', 3.00, 'c1000000-0000-0000-0000-000000000011', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/272684694_43SoftDrinks2.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- A1 MENU ITEMS
	-- Signature Kebabs
	('d2000000-0001-0000-0000-000000000001', 'Twelve Inch Wrap', '12" wrap (150g) - Signature kebab wrap with grilled chicken', 14.00, 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/875772946_10ChickenRoll1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0001-0000-0000-000000000002', 'Ten Inch Wrap', '10" wrap (100g) - Signature kebab wrap with grilled chicken', 11.00, 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/875772946_10ChickenRoll1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0001-0000-0000-000000000003', 'Fresh Pita', 'Fresh pita (100g) - Signature kebab in fresh pita bread', 12.00, 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/548042839_20ChickenTombik1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- Kebab Bowl
	('d2000000-0002-0000-0000-000000000001', 'Kebab Bowl - Salad Only', 'Kebab bowl with fresh salad base', 14.00, 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1855254249_60ChickenSalad.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0002-0000-0000-000000000002', 'Kebab Bowl - Salad and Rice', 'Kebab bowl with salad and rice base', 15.00, 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/53393089_40ChickenRice.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0002-0000-0000-000000000003', 'Kebab Bowl - Salad and Fries', 'Kebab bowl with salad and french fries base', 16.00, 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1662471778_50ChickenFries1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- Falafel
	('d2000000-0003-0000-0000-000000000001', 'Falafel Twelve Inch Wrap', '100% plant-based falafel in a 12-inch wrap (4 pcs)', 14.00, 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/89806033_120FallafelRollvegan1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0003-0000-0000-000000000002', 'Falafel Ten Inch Wrap', '100% plant-based falafel in a 10-inch wrap (3 pcs)', 11.00, 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/89806033_120FallafelRollvegan1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0003-0000-0000-000000000003', 'Falafel Homemade Pita', 'Homemade pita filled with crispy plant-based falafels (3 pcs)', 12.00, 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/1375458996_falafeltombik.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0003-0000-0000-000000000004', 'Falafel Plate', '100% plant-based falafel plate (6 pcs)', 16.00, 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/574589317_130FalafelPlate1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- Sides
	('d2000000-0004-0000-0000-000000000001', 'Hummus with Pita', 'Creamy hummus dip served with fresh pita bread', 10.00, 'c2000000-0000-0000-0000-000000000004', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/364237033_140Hummus.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0004-0000-0000-000000000002', 'Three Piece Falafel', 'Three pieces of crispy plant-based falafels', 6.00, 'c2000000-0000-0000-0000-000000000004', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/574589317_130FalafelPlate1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0004-0000-0000-000000000003', 'Fries Large', 'Large portion of crispy golden fries', 7.50, 'c2000000-0000-0000-0000-000000000004', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/862897504_150FrenchFries.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0004-0000-0000-000000000004', 'Fries Small', 'Small portion of crispy golden fries', 4.50, 'c2000000-0000-0000-0000-000000000004', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/862897504_150FrenchFries.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- Drinks
	('d2000000-0005-0000-0000-000000000001', 'Soft Drinks', 'Assorted carbonated beverages', 2.50, 'c2000000-0000-0000-0000-000000000005', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/272684694_43SoftDrinks2.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0005-0000-0000-000000000002', 'Mineral Water', 'Bottled mineral water', 2.00, 'c2000000-0000-0000-0000-000000000005', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/1484899931_42MineralWater1.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0005-0000-0000-000000000003', 'Turkish Ayran', 'Salted cold yogurt beverage', 4.00, 'c2000000-0000-0000-0000-000000000005', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/getz-prod/40a27a5c-7901-4397-b244-10dca6109848/262280452_1006442618200TurkishAyran.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('d2000000-0005-0000-0000-000000000004', 'Turkish Lemonade', 'Sweet and refreshing homemade Turkish lemonade', 3.00, 'c2000000-0000-0000-0000-000000000005', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'https://s3-ap-southeast-1.amazonaws.com/smoovpos-live/40a27a5c-7901-4397-b244-10dca6109848/272684694_43SoftDrinks2.jpg', true, '2025-10-27 00:00:00', '2025-10-27 00:00:00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- Addons
-- ============================================================================
INSERT INTO addons (addon_id, category_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated) VALUES
	-- KEBAB ROLL (shared across category)
	('b1000000-0001-0000-0000-000000000001', 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Cheese', 'Additional melted cheese', 2.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0001-0000-0000-000000000002', 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Meat', 'Additional portion of meat', 4.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0001-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Garlic Sauce', 'Extra garlic sauce on the side', 1.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0001-0000-0000-000000000004', 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Spicy Sauce', 'Hot chili sauce', 1.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0002-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000001', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Jalapeños', 'Spicy jalapeño peppers', 1.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- KEBAB RICE (shared across category)
	('b1000000-0010-0000-0000-000000000001', 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Meat', 'Additional portion of meat', 4.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0010-0000-0000-000000000002', 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Hummus', 'Side of creamy hummus', 3.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0010-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Grilled Vegetables', 'Assorted grilled vegetables', 3.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0010-0000-0000-000000000004', 'c1000000-0000-0000-0000-000000000003', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Spicy Sauce', 'Hot chili sauce', 1.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- PIDE AND LAHMAJUN (shared across category)
	('b1000000-0022-0000-0000-000000000001', 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Cheese', 'Additional melted cheese', 2.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0022-0000-0000-000000000002', 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Meat', 'Additional portion of meat', 4.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0022-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000007', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Mushrooms', 'Sautéed mushrooms', 2.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	
	-- VEGETARIAN OPTIONS (shared across category)
	('b1000000-0033-0000-0000-000000000001', 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Falafel', 'Additional falafel pieces (4 pcs)', 3.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0033-0000-0000-000000000002', 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Hummus', 'Additional portion of hummus', 2.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('b1000000-0033-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000008', 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', 'Extra Pita Bread', 'Additional pita bread', 1.50, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- A1 ADDONS
	-- Signature Kebabs Addons
	('a2000000-0001-0000-0000-000000000001', 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Make It A Meal', 'Upgrade your main to a meal (with Fries and Drink)', 5.00, true, 1, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0001-0000-0000-000000000002', 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Chicken Standard', 'Add extra grilled chicken (60g)', 3.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0001-0000-0000-000000000003', 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Chicken Large', 'Add a large portion of chicken (100g)', 5.50, true, 1, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0001-0000-0000-000000000004', 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Salad', 'Add extra fresh salad mix', 3.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0001-0000-0000-000000000005', 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Nacho Cheese Sauce', 'Melted nacho cheese dip', 1.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0001-0000-0000-000000000006', 'c2000000-0000-0000-0000-000000000001', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Hummus Spread', 'Side of creamy hummus dip', 1.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- Kebab Bowl Addons
	('a2000000-0002-0000-0000-000000000002', 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Chicken Standard', 'Add extra grilled chicken (60g)', 3.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0002-0000-0000-000000000003', 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Chicken Large', 'Add a large portion of chicken (100g)', 5.50, true, 1, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0002-0000-0000-000000000004', 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Salad', 'Add extra fresh salad mix', 3.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0002-0000-0000-000000000005', 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Nacho Cheese Sauce', 'Melted nacho cheese dip', 1.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0002-0000-0000-000000000006', 'c2000000-0000-0000-0000-000000000002', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Hummus Spread', 'Side of creamy hummus dip', 1.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),

	-- Falafel Addons
	('a2000000-0003-0000-0000-000000000001', 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Make It A Meal', 'Upgrade your main to a meal (with Fries and Drink)', 5.00, true, 1, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0003-0000-0000-000000000004', 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Extra Salad', 'Add extra fresh salad mix', 3.00, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0003-0000-0000-000000000005', 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Nacho Cheese Sauce', 'Melted nacho cheese dip', 1.00, true, 3, '2025-10-27 00:00:00', '2025-10-27 00:00:00'),
	('a2000000-0003-0000-0000-000000000006', 'c2000000-0000-0000-0000-000000000003', 'a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1', 'Hummus Spread', 'Side of creamy hummus dip', 1.50, true, 2, '2025-10-27 00:00:00', '2025-10-27 00:00:00')
ON CONFLICT DO NOTHING;
