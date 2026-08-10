-- +goose Up

-- ─── Roles ──────────────────────────────────────────────────────────────
INSERT INTO roles (id, name) VALUES
  (1, 'passenger'),
  (2, 'operator'),
  (3, 'admin')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- ─── Operators ───────────────────────────────────────────────────────────
INSERT INTO operators (id, name, code, hotline, email, status, metadata) VALUES
  (1, '{"vi":"Phương Trang","en":"Phuong Trang"}'::jsonb, 'PT', '+84900000001', 'contact@phuongtrang.vn', 'active', '{"logoColor":"#059669","rating":4.6,"totalReviews":12480}'::jsonb),
  (2, '{"vi":"Hoàng Yến","en":"Hoang Yen"}'::jsonb,     'HY', '+84900000002', 'contact@hoangyen.vn',   'active', '{"logoColor":"#2563eb","rating":4.3,"totalReviews":28910}'::jsonb),
  (3, '{"vi":"Mai Linh","en":"Mai Linh"}'::jsonb,       'ML', '+84900000003', 'contact@mailinh.vn',     'active', '{"logoColor":"#ea580c","rating":4.1,"totalReviews":8730}'::jsonb),
  (4, '{"vi":"Kumho Express","en":"Kumho Express"}'::jsonb, 'KE', '+84900000004', 'contact@kumho.vn',  'active', '{"logoColor":"#7c3aed","rating":4.7,"totalReviews":15620}'::jsonb),
  (5, '{"vi":"Sapa Express","en":"Sapa Express"}'::jsonb,  'SE', '+84900000005', 'contact@sapaexpress.vn', 'active', '{"logoColor":"#dc2626","rating":3.9,"totalReviews":5410}'::jsonb)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, code = EXCLUDED.code;

-- ─── Administrative areas (cities) ──────────────────────────────────────
INSERT INTO administrative_areas (id, name, code, level) VALUES
  (1,  '{"vi":"Hà Nội","en":"Hanoi"}'::jsonb,        'ha-noi',       'city'),
  (2,  '{"vi":"Hồ Chí Minh","en":"Ho Chi Minh"}'::jsonb, 'ho-chi-minh', 'city'),
  (3,  '{"vi":"Đà Nẵng","en":"Da Nang"}'::jsonb,     'da-nang',      'city'),
  (4,  '{"vi":"Hải Phòng","en":"Hai Phong"}'::jsonb, 'hai-phong',    'city'),
  (5,  '{"vi":"Nha Trang","en":"Nha Trang"}'::jsonb, 'nha-trang',    'city'),
  (6,  '{"vi":"Đà Lạt","en":"Da Lat"}'::jsonb,       'da-lat',       'city'),
  (7,  '{"vi":"Huế","en":"Hue"}'::jsonb,             'hue',          'city'),
  (8,  '{"vi":"Vũng Tàu","en":"Vung Tau"}'::jsonb,   'vung-tau',     'city'),
  (9,  '{"vi":"Cần Thơ","en":"Can Tho"}'::jsonb,     'can-tho',      'city'),
  (10, '{"vi":"Quy Nhơn","en":"Quy Nhon"}'::jsonb,   'quy-nhon',     'city'),
  (11, '{"vi":"Phú Quốc","en":"Phu Quoc"}'::jsonb,   'phu-quoc',     'city'),
  (12, '{"vi":"Hội An","en":"Hoi An"}'::jsonb,       'hoi-an',       'city'),
  (13, '{"vi":"Ninh Bình","en":"Ninh Binh"}'::jsonb, 'ninh-binh',    'city'),
  (14, '{"vi":"Sapa","en":"Sapa"}'::jsonb,           'sapa',         'city'),
  (15, '{"vi":"Long Xuyên","en":"Long Xuyen"}'::jsonb, 'long-xuyen', 'city')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, code = EXCLUDED.code;

-- ─── Locations (boarding points for each city) ──────────────────────────
INSERT INTO locations (id, administrative_area_id, name, type, address, latitude, longitude, status)
SELECT aa.id, aa.id, aa.name, 'bus_station', '', NULL, NULL, 'active'
FROM administrative_areas aa
WHERE aa.level = 'city'
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- ─── Vehicles ──────────────────────────────────────────────────────────
INSERT INTO vehicles (id, operator_id, plate_number, vehicle_type, total_seats, status, seat_layout, amenities) VALUES
  (1,  1, '51B-1234', 'standard',      40, 'active', '{"layoutId":"standard-2-2","rows":10,"cols":[2,2],"aisleAfter":2}'::jsonb,     '["wifi","ac","power"]'::jsonb),
  (2,  1, '51B-2234', 'premium',       32, 'active', '{"layoutId":"premium-2-2","rows":8,"cols":[2,2],"aisleAfter":2}'::jsonb,      '["wifi","ac","power","usb","snack","legroom"]'::jsonb),
  (3,  2, '51B-3234', 'standard',      40, 'active', '{"layoutId":"standard-2-2","rows":10,"cols":[2,2],"aisleAfter":2}'::jsonb,     '["wifi","ac","power"]'::jsonb),
  (4,  2, '51B-4234', 'double-decker', 64, 'active', '{"layoutId":"double-decker-2-2","rows":16,"cols":[2,2],"aisleAfter":2}'::jsonb, '["wifi","ac","toilet","power","entertainment"]'::jsonb),
  (5,  3, '51B-5234', 'standard',      40, 'active', '{"layoutId":"standard-2-2","rows":10,"cols":[2,2],"aisleAfter":2}'::jsonb,     '["wifi","ac","power"]'::jsonb),
  (6,  3, '51B-6234', 'premium',       32, 'active', '{"layoutId":"premium-2-2","rows":8,"cols":[2,2],"aisleAfter":2}'::jsonb,      '["wifi","ac","power","usb","snack","legroom"]'::jsonb),
  (7,  4, '51B-7234', 'double-decker', 64, 'active', '{"layoutId":"double-decker-2-2","rows":16,"cols":[2,2],"aisleAfter":2}'::jsonb, '["wifi","ac","toilet","power","entertainment"]'::jsonb),
  (8,  4, '51B-8234', 'standard',      40, 'active', '{"layoutId":"standard-2-2","rows":10,"cols":[2,2],"aisleAfter":2}'::jsonb,     '["wifi","ac","power"]'::jsonb),
  (9,  5, '51B-9234', 'standard',      40, 'active', '{"layoutId":"standard-2-2","rows":10,"cols":[2,2],"aisleAfter":2}'::jsonb,     '["wifi","ac","power"]'::jsonb),
  (10, 5, '51B-0234', 'double-decker', 64, 'active', '{"layoutId":"double-decker-2-2","rows":16,"cols":[2,2],"aisleAfter":2}'::jsonb, '["wifi","ac","toilet","power","entertainment"]'::jsonb)
ON CONFLICT (id) DO UPDATE SET plate_number = EXCLUDED.plate_number;

-- ─── Routes ──────────────────────────────────────────────────────────────
INSERT INTO routes (id, operator_id, origin_id, destination_id, name, estimated_duration_minutes, status) VALUES
  (1,  1, 1,  2,  '{"vi":"Hà Nội - Hồ Chí Minh","en":"Hanoi - Ho Chi Minh"}'::jsonb,  2100, 'active'),
  (2,  2, 2,  1,  '{"vi":"Hồ Chí Minh - Hà Nội","en":"Ho Chi Minh - Hanoi"}'::jsonb,  2100, 'active'),
  (3,  1, 2,  6,  '{"vi":"Hồ Chí Minh - Đà Lạt","en":"Ho Chi Minh - Da Lat"}'::jsonb, 420, 'active'),
  (4,  2, 6,  2,  '{"vi":"Đà Lạt - Hồ Chí Minh","en":"Da Lat - Ho Chi Minh"}'::jsonb, 420, 'active'),
  (5,  3, 2,  8,  '{"vi":"Hồ Chí Minh - Vũng Tàu","en":"Ho Chi Minh - Vung Tau"}'::jsonb, 120, 'active'),
  (6,  3, 8,  2,  '{"vi":"Vũng Tàu - Hồ Chí Minh","en":"Vung Tau - Ho Chi Minh"}'::jsonb, 120, 'active'),
  (7,  4, 2,  9,  '{"vi":"Hồ Chí Minh - Cần Thơ","en":"Ho Chi Minh - Can Tho"}'::jsonb, 240, 'active'),
  (8,  4, 9,  2,  '{"vi":"Cần Thơ - Hồ Chí Minh","en":"Can Tho - Ho Chi Minh"}'::jsonb, 240, 'active'),
  (9,  5, 2,  5,  '{"vi":"Hồ Chí Minh - Nha Trang","en":"Ho Chi Minh - Nha Trang"}'::jsonb, 540, 'active'),
  (10, 5, 5,  2,  '{"vi":"Nha Trang - Hồ Chí Minh","en":"Nha Trang - Ho Chi Minh"}'::jsonb, 540, 'active'),
  (11, 1, 1,  4,  '{"vi":"Hà Nội - Hải Phòng","en":"Hanoi - Hai Phong"}'::jsonb,      120, 'active'),
  (12, 1, 4,  1,  '{"vi":"Hải Phòng - Hà Nội","en":"Hai Phong - Hanoi"}'::jsonb,      120, 'active'),
  (13, 2, 1,  13, '{"vi":"Hà Nội - Ninh Bình","en":"Hanoi - Ninh Binh"}'::jsonb,       90, 'active'),
  (14, 2, 13, 1, '{"vi":"Ninh Bình - Hà Nội","en":"Ninh Binh - Hanoi"}'::jsonb,       90, 'active'),
  (15, 3, 1,  14, '{"vi":"Hà Nội - Sapa","en":"Hanoi - Sapa"}'::jsonb,                360, 'active'),
  (16, 3, 14, 1, '{"vi":"Sapa - Hà Nội","en":"Sapa - Hanoi"}'::jsonb,                360, 'active'),
  (17, 4, 3,  7,  '{"vi":"Đà Nẵng - Huế","en":"Da Nang - Hue"}'::jsonb,                90, 'active'),
  (18, 4, 7,  3,  '{"vi":"Huế - Đà Nẵng","en":"Hue - Da Nang"}'::jsonb,                90, 'active'),
  (19, 5, 3,  12, '{"vi":"Đà Nẵng - Hội An","en":"Da Nang - Hoi An"}'::jsonb,           45, 'active'),
  (20, 5, 12, 3, '{"vi":"Hội An - Đà Nẵng","en":"Hoi An - Da Nang"}'::jsonb,           45, 'active'),
  (21, 1, 3,  10, '{"vi":"Đà Nẵng - Quy Nhơn","en":"Da Nang - Quy Nhon"}'::jsonb,     240, 'active'),
  (22, 1, 10, 3, '{"vi":"Quy Nhơn - Đà Nẵng","en":"Quy Nhon - Da Nang"}'::jsonb,     240, 'active'),
  (23, 2, 5,  6,  '{"vi":"Nha Trang - Đà Lạt","en":"Nha Trang - Da Lat"}'::jsonb,     240, 'active'),
  (24, 2, 6,  5,  '{"vi":"Đà Lạt - Nha Trang","en":"Da Lat - Nha Trang"}'::jsonb,     240, 'active'),
  (25, 3, 2,  15, '{"vi":"Hồ Chí Minh - Long Xuyên","en":"Ho Chi Minh - Long Xuyen"}'::jsonb, 240, 'active'),
  (26, 3, 15, 2, '{"vi":"Long Xuyên - Hồ Chí Minh","en":"Long Xuyen - Ho Chi Minh"}'::jsonb, 240, 'active'),
  (27, 4, 2,  11, '{"vi":"Hồ Chí Minh - Phú Quốc","en":"Ho Chi Minh - Phu Quoc"}'::jsonb, 360, 'active'),
  (28, 4, 11, 2, '{"vi":"Phú Quốc - Hồ Chí Minh","en":"Phu Quoc - Ho Chi Minh"}'::jsonb, 360, 'active'),
  (29, 5, 1,  3,  '{"vi":"Hà Nội - Đà Nẵng","en":"Hanoi - Da Nang"}'::jsonb,           780, 'active'),
  (30, 5, 3,  1,  '{"vi":"Đà Nẵng - Hà Nội","en":"Da Nang - Hanoi"}'::jsonb,           780, 'active')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- ─── Route stops ──────────────────────────────────────────────────────────
INSERT INTO route_stops (route_id, location_id, stop_order, stop_type, arrival_offset_minutes, departure_offset_minutes) VALUES
  (1, 1,  1, 'origin',       0,    0),
  (1, 3,  2, 'intermediate',  960,  960),
  (1, 5,  3, 'intermediate',  1380, 1380),
  (1, 2,  4, 'destination',   2100, 2100),
  (2, 2,  1, 'origin',       0,    0),
  (2, 5,  2, 'intermediate',  720,  720),
  (2, 3,  3, 'intermediate',  1140, 1140),
  (2, 1,  4, 'destination',   2100, 2100),
  (3, 2,  1, 'origin',       0,    0),
  (3, 8,  2, 'intermediate',  120,  120),
  (3, 6,  3, 'destination',   420,  420),
  (4, 6,  1, 'origin',       0,    0),
  (4, 8,  2, 'intermediate',  300,  300),
  (4, 2,  3, 'destination',   420,  420),
  (9, 2,  1, 'origin',       0,    0),
  (9, 10, 2, 'intermediate',  360,  360),
  (9, 5,  3, 'destination',   540,  540),
  (10, 5,  1, 'origin',       0,    0),
  (10, 10, 2, 'intermediate',  180,  180),
  (10, 2,  3, 'destination',   540,  540),
  (15, 1,  1, 'origin',       0,    0),
  (15, 13, 2, 'intermediate',  90,   90),
  (15, 14, 3, 'destination',  360,  360),
  (16, 14, 1, 'origin',       0,    0),
  (16, 13, 2, 'intermediate', 270,  270),
  (16, 1,  3, 'destination',  360,  360),
  (29, 1,  1, 'origin',       0,    0),
  (29, 13, 2, 'intermediate',  90,   90),
  (29, 3,  3, 'destination',  780,  780),
  (30, 3,  1, 'origin',       0,    0),
  (30, 13, 2, 'intermediate', 690,  690),
  (30, 1,  3, 'destination',  780,  780)
ON CONFLICT DO NOTHING;

-- ─── Trips ────────────────────────────────────────────────────────────────
-- +goose StatementBegin
DO $$
DECLARE
  trip_id_seq INTEGER := 1;
  trip_date DATE;
  r RECORD;
  num_trips INTEGER;
  dep_hour INTEGER;
  dep_min INTEGER;
  dep_time TIMESTAMPTZ;
  arr_time TIMESTAMPTZ;
  bus_idx INTEGER;
  vehicle RECORD;
  operator_id INTEGER;
  base_price BIGINT;
  price BIGINT;
  total_seats INTEGER;
  available_seats INTEGER;
  day_offset INTEGER;
  i INTEGER;
  bus_type TEXT;
  amenities JSONB;
BEGIN
  DELETE FROM seat_inventories WHERE trip_id <= 10000;
  DELETE FROM trips WHERE id <= 10000;

  FOR day_offset IN 0..6 LOOP
    trip_date := CURRENT_DATE + day_offset;

    FOR r IN SELECT * FROM routes WHERE id <= 30 LOOP
      num_trips := 6 + ((day_offset + r.id) % 5);

      FOR i IN 0..num_trips-1 LOOP
        dep_hour := CASE (i % 4)
          WHEN 0 THEN 6 + ((day_offset + r.id + i) % 3)
          WHEN 1 THEN 11 + ((day_offset + r.id + i) % 3)
          WHEN 2 THEN 14 + ((day_offset + r.id + i) % 4)
          WHEN 3 THEN 18 + ((day_offset + r.id + i) % 5)
        END;
        dep_min := ((day_offset + r.id + i) % 12) * 5;
        dep_time := trip_date + (dep_hour || ' hours')::interval + (dep_min || ' minutes')::interval;
        arr_time := dep_time + (r.estimated_duration_minutes || ' minutes')::interval;

        bus_idx := (day_offset + r.id + i) % 10;
        SELECT * INTO vehicle FROM vehicles WHERE id = bus_idx + 1;
        operator_id := (bus_idx % 5) + 1;

        bus_type := vehicle.vehicle_type;
        total_seats := CASE WHEN bus_type = 'double-decker' THEN 64 WHEN bus_type = 'premium' THEN 32 ELSE 40 END;
        available_seats := (day_offset + r.id + i + 7) % (total_seats + 1);

        base_price := (80000 + r.estimated_duration_minutes * 500);
        base_price := CASE bus_type
          WHEN 'premium' THEN ROUND(base_price * 1.4)
          WHEN 'double-decker' THEN ROUND(base_price * 0.85)
          ELSE base_price
        END;
        price := ROUND((base_price + ((day_offset + r.id + i) % 5) * 10000) / 10000) * 10000;

        amenities := CASE bus_type
          WHEN 'premium' THEN '["wifi","ac","power","usb","snack","legroom"]'::jsonb
          WHEN 'double-decker' THEN '["wifi","ac","toilet","power","entertainment"]'::jsonb
          ELSE '["wifi","ac","power"]'::jsonb
        END;

        INSERT INTO trips (
          id, operator_id, route_id, vehicle_id,
          origin_id, destination_id,
          departure_time, arrival_time,
          base_price, status,
          available_seats_count, booked_seats_count, held_seats_count,
          searchable_snapshot
        ) VALUES (
          trip_id_seq, operator_id, r.id, vehicle.id,
          r.origin_id, r.destination_id,
          dep_time, arr_time,
          price, 'scheduled',
          available_seats, total_seats - available_seats, 0,
          jsonb_build_object(
            'busType', bus_type,
            'amenities', amenities,
            'totalSeats', total_seats,
            'durationMinutes', r.estimated_duration_minutes
          )
        );

        INSERT INTO seat_inventories (trip_id, seat_no, floor, row_no, column_no, seat_type, price, status)
        SELECT
          trip_id_seq,
          'A' || row_num || '-' || col_num,
          CASE WHEN bus_type = 'double-decker' AND row_num > 8 THEN 2 ELSE 1 END,
          row_num,
          col_num,
          CASE WHEN bus_type = 'premium' THEN 'premium' ELSE 'standard' END,
          price,
          CASE WHEN (row_num * 10 + col_num + trip_id_seq) % (total_seats + 1) < available_seats
               THEN 'available' ELSE 'occupied' END
        FROM generate_series(1, total_seats / CASE WHEN bus_type = 'double-decker' THEN 4 ELSE 2 END) AS row_num
        CROSS JOIN generate_series(1, 2) AS col_num
        WHERE (row_num - 1) * 2 + col_num <= total_seats;

        trip_id_seq := trip_id_seq + 1;
      END LOOP;
    END LOOP;
  END LOOP;
END $$;
-- +goose StatementEnd

-- ─── Users ────────────────────────────────────────────────────────────────
-- Demo users: password is "password123" (bcrypt hash below). Not a real secret.
INSERT INTO users (id, full_name, email, phone, role_id, operator_id, role, password, language, status) VALUES
  (1, 'Phương Trang Operations', 'operator@busgo.app', '+84900000010', 2, 1, 'operator',
   '$2a$10$2SKxrw75fx1aiHZKIybmF.GlXPp6STVIX6.4ltyOVRigKURMMfDXa', 'vi', 'active')
ON CONFLICT (id) DO UPDATE SET full_name = EXCLUDED.full_name, email = EXCLUDED.email, role = EXCLUDED.role;

INSERT INTO users (id, full_name, email, phone, role_id, role, password, language, status) VALUES
  (2, 'Nguyen Minh Anh', 'passenger@busgo.app', '+84900000020', 1, 'passenger',
   '$2a$10$2SKxrw75fx1aiHZKIybmF.GlXPp6STVIX6.4ltyOVRigKURMMfDXa', 'vi', 'active')
ON CONFLICT (id) DO UPDATE SET full_name = EXCLUDED.full_name, email = EXCLUDED.email, role = EXCLUDED.role;

-- +goose Down
DELETE FROM seat_inventories WHERE trip_id <= 10000;
DELETE FROM trips WHERE id <= 10000;
DELETE FROM route_stops WHERE route_id <= 30;
DELETE FROM routes WHERE id <= 30;
DELETE FROM vehicles WHERE id <= 10;
DELETE FROM locations WHERE id <= 15;
DELETE FROM administrative_areas WHERE id <= 15;
DELETE FROM operators WHERE id <= 5;
DELETE FROM users WHERE id <= 2;
