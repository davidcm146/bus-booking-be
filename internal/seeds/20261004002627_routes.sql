-- =====================================================================
-- Seed data: routes (tuyến xe, sinh cả chiều đi và chiều về)
-- Phụ thuộc (phải chạy trước): locations_seed.sql, operators_seed.sql
--
-- !!! DỮ LIỆU MẪU - HỢP LÝ NHƯNG KHÔNG PHẢI LỊCH CHẠY THẬT !!!
--   * Cặp thành phố đầu/cuối lấy theo thông tin đã tra về từng nhà xe (Phương Trang,
--     Hoàng Long, Kumho Samco, Quang Dũng, Hà Sơn). Tuyến của Thành Công, Sơn Tùng,
--     Mai Linh Express là tuyến giả định cho hợp lý.
--   * Điểm dừng, khoảng cách (km) và thời gian (phút) là ước lượng theo quốc lộ/cao tốc
--     thực tế, không phải số liệu công bố của nhà xe.
--   * Điểm dừng 'rest' = các bến xe dọc đường (bảng locations). Trạm dừng nghỉ cao tốc
--     chưa có toạ độ nên chưa đưa vào.
--   * Mỗi tuyến tự sinh thêm chiều ngược lại (đảo origin/destination và thứ tự stops).
--   * metadata = {"seed":"routes_v1","sample":true} để dọn dữ liệu:
--       DELETE FROM routes WHERE metadata->>'seed' = 'routes_v1';
--   * route_snapshot được dựng từ chính bảng operators/locations/administrative_areas nên
--     location_id luôn khớp. Định dạng:
--       { "operator":    {id, code, name{vi,en}, brand, hotline, logo_url},
--         "origin":      {location_id, name{vi,en}, city{id, name{vi,en}}, address, lat, lng},
--         "destination": {... giống origin},
--         "stops": [{order, location_id, name{vi,en}, type, address, lat, lng,
--                    offset_minutes, distance_from_origin_km, stop_duration_minutes}] }
--     - city = tỉnh/thành (parent_id của phường/xã nơi đặt bến).
--     - hotline bỏ khoảng trắng; logo_url lấy từ operators.metadata->>'logo_url' (chưa có -> null).
--     - distance_from_origin_km: cộng dồn khoảng cách đường chim bay giữa các bến liên tiếp,
--       co giãn theo tỉ lệ để bến cuối = distance_km của tuyến.
--     - Bến dọc đường dừng 10 phút; offset_minutes = phút xe tới bến tính từ giờ xuất bến,
--       thời gian chạy (duration - tổng phút dừng) chia theo tỉ lệ quãng đường; bến cuối = duration.
--   * origin_id / destination_id = id của administrative_areas (phường/xã nơi đặt bến đi / bến đến),
--     lấy từ locations.administrative_area_id của bến đầu/cuối. Muốn trỏ về cấp tỉnh/thành thì
--     thay bằng parent_id của phường/xã đó (xem ghi chú cuối file).
-- =====================================================================
BEGIN;

-- ---------------------------------------------------------------------
-- 0) Định nghĩa tuyến (chiều đi)
-- ---------------------------------------------------------------------
CREATE TEMP TABLE route_defs (
    op_code        text,
    origin         text,      -- tên vi của bến xe đi
    dest           text,      -- tên vi của bến xe đến
    origin_city    text,
    origin_city_en text,
    dest_city      text,
    dest_city_en   text,
    via            text[],    -- các bến dừng dọc đường (theo thứ tự)
    distance_km    numeric(10,2),
    duration_min   integer
) ON COMMIT DROP;

INSERT INTO route_defs VALUES
-- ===== Phương Trang =====
('phuong-trang','Bến xe Miền Đông mới','Bến xe Liên tỉnh Đà Lạt','Sài Gòn','Ho Chi Minh City','Đà Lạt','Da Lat',
  ARRAY['Bến xe Long Khánh','Bến xe Bảo Lộc'], 300, 420),
('phuong-trang','Bến xe Miền Tây','Bến xe Cà Mau','Sài Gòn','Ho Chi Minh City','Cà Mau','Ca Mau',
  ARRAY['Bến xe Tiền Giang (Mỹ Tho)','Bến xe Vĩnh Long','Bến xe Trung tâm Cần Thơ','Bến xe Sóc Trăng','Bến xe Bạc Liêu'], 340, 510),
('phuong-trang','Bến xe Miền Tây','Bến xe Trung tâm Cần Thơ','Sài Gòn','Ho Chi Minh City','Cần Thơ','Can Tho',
  ARRAY['Bến xe Tiền Giang (Mỹ Tho)','Bến xe Vĩnh Long'], 170, 210),
('phuong-trang','Bến xe Miền Đông mới','Bến xe phía Bắc Nha Trang','Sài Gòn','Ho Chi Minh City','Nha Trang','Nha Trang',
  ARRAY['Bến xe Phan Thiết','Bến xe Phan Rang'], 430, 600),
-- ===== Hoàng Long =====
('hoang-long','Bến xe Giáp Bát','Bến xe Niệm Nghĩa','Hà Nội','Ha Noi','Hải Phòng','Hai Phong',
  ARRAY['Bến xe Hải Dương'], 120, 150),
('hoang-long','Bến xe Mỹ Đình','Bến xe phía Bắc Thanh Hóa','Hà Nội','Ha Noi','Thanh Hóa','Thanh Hoa',
  ARRAY['Bến xe Phủ Lý','Bến xe Ninh Bình'], 155, 200),
('hoang-long','Bến xe Giáp Bát','Bến xe phía Bắc Nha Trang','Hà Nội','Ha Noi','Nha Trang','Nha Trang',
  ARRAY['Bến xe Ninh Bình','Bến xe phía Bắc Thanh Hóa','Bến xe phía Bắc Vinh','Bến xe Hà Tĩnh','Bến xe Đồng Hới',
        'Bến xe Đông Hà','Bến xe phía Nam Huế','Bến xe Trung tâm Đà Nẵng','Bến xe Quảng Ngãi','Bến xe Quy Nhơn','Bến xe Tuy Hòa'], 1300, 1500),
-- ===== Kumho Samco =====
('kumho-samco','Bến xe Miền Đông mới','Bến xe Phan Thiết','Sài Gòn','Ho Chi Minh City','Phan Thiết','Phan Thiet',
  ARRAY[]::text[], 200, 270),
('kumho-samco','Bến xe Miền Đông mới','Bến xe Vũng Tàu','Sài Gòn','Ho Chi Minh City','Vũng Tàu','Vung Tau',
  ARRAY['Bến xe Long Thành','Bến xe Bà Rịa'], 100, 150),
('kumho-samco','Bến xe Miền Đông mới','Bến xe Đắk Lắk','Sài Gòn','Ho Chi Minh City','Buôn Ma Thuột','Buon Ma Thuot',
  ARRAY['Bến xe Đồng Xoài','Bến xe Gia Nghĩa'], 350, 480),
('kumho-samco','Bến xe Miền Tây','Bến xe Rạch Giá','Sài Gòn','Ho Chi Minh City','Rạch Giá','Rach Gia',
  ARRAY['Bến xe Tân An','Bến xe Vĩnh Long','Bến xe Trung tâm Cần Thơ'], 290, 380),
('kumho-samco','Bến xe Rạch Giá','Bến xe Hà Tiên','Rạch Giá','Rach Gia','Hà Tiên','Ha Tien',
  ARRAY[]::text[], 100, 150),
-- ===== Quang Dũng =====
('quang-dung','Bến xe Đông Hà','Bến xe Trung tâm Đà Nẵng','Đông Hà','Dong Ha','Đà Nẵng','Da Nang',
  ARRAY['Bến xe phía Nam Huế'], 175, 210),
('quang-dung','Bến xe phía Nam Huế','Bến xe Hà Tĩnh','Huế','Hue','Hà Tĩnh','Ha Tinh',
  ARRAY['Bến xe Đông Hà','Bến xe Đồng Hới'], 330, 390),
('quang-dung','Bến xe Mỹ Đình','Bến xe phía Nam Huế','Hà Nội','Ha Noi','Huế','Hue',
  ARRAY['Bến xe Ninh Bình','Bến xe phía Bắc Thanh Hóa','Bến xe phía Bắc Vinh','Bến xe Hà Tĩnh','Bến xe Đồng Hới','Bến xe Đông Hà'], 660, 780),
-- ===== Sơn Tùng (giả định) =====
('son-tung','Bến xe Quy Nhơn','Bến xe Miền Đông mới','Quy Nhơn','Quy Nhon','Sài Gòn','Ho Chi Minh City',
  ARRAY['Bến xe Tuy Hòa','Bến xe phía Bắc Nha Trang','Bến xe Phan Thiết'], 650, 720),
('son-tung','Bến xe Quy Nhơn','Bến xe Pleiku','Quy Nhơn','Quy Nhon','Pleiku','Pleiku',
  ARRAY[]::text[], 190, 270),
-- ===== Thành Công (giả định) =====
('thanh-cong','Bến xe Mỹ Đình','Bến xe Lào Cai','Hà Nội','Ha Noi','Lào Cai','Lao Cai',
  ARRAY['Bến xe Việt Trì'], 300, 330),
('thanh-cong','Bến xe Mỹ Đình','Bến xe Bãi Cháy','Hà Nội','Ha Noi','Hạ Long','Ha Long',
  ARRAY[]::text[], 160, 180),
-- ===== Mai Linh Express (giả định) =====
('mai-linh-express','Bến xe Trung tâm Đà Nẵng','Bến xe phía Nam Huế','Đà Nẵng','Da Nang','Huế','Hue',
  ARRAY[]::text[], 100, 120),
('mai-linh-express','Bến xe Trung tâm Đà Nẵng','Bến xe Quảng Ngãi','Đà Nẵng','Da Nang','Quảng Ngãi','Quang Ngai',
  ARRAY['Bến xe Tam Kỳ'], 130, 180),
-- ===== Hà Sơn =====
('ha-son','Bến xe Mỹ Đình','Bến xe Sa Pa','Hà Nội','Ha Noi','Sa Pa','Sa Pa',
  ARRAY['Bến xe Việt Trì','Bến xe Lào Cai'], 320, 360),
('ha-son','Bến xe phía Bắc Vinh','Bến xe Sa Pa','Vinh','Vinh','Sa Pa','Sa Pa',
  ARRAY['Bến xe Ninh Bình','Bến xe Việt Trì','Bến xe Lào Cai'], 700, 840);

-- ---------------------------------------------------------------------
-- 1) Kiểm tra dữ liệu phụ thuộc: thiếu operator/location thì dừng, không insert dở
-- ---------------------------------------------------------------------
DO $$
DECLARE
    missing text;
BEGIN
    SELECT string_agg(DISTINCT s.n, ', ') INTO missing
    FROM (SELECT unnest(ARRAY[origin] || via || ARRAY[dest]) AS n FROM route_defs) s
    WHERE NOT EXISTS (
        SELECT 1 FROM locations l
        WHERE l.type = 'bus_station' AND l.name ->> 'vi' = s.n AND NOT l.deleted
    );
    IF missing IS NOT NULL THEN
        RAISE EXCEPTION 'Thiếu bến xe trong locations: %', missing;
    END IF;

    SELECT string_agg(DISTINCT l.name ->> 'vi', ', ') INTO missing
    FROM locations l
    WHERE l.type = 'bus_station' AND NOT l.deleted
      AND (l.latitude IS NULL OR l.longitude IS NULL)
      AND l.name ->> 'vi' IN (SELECT unnest(ARRAY[origin] || via || ARRAY[dest]) FROM route_defs);
    IF missing IS NOT NULL THEN
        RAISE EXCEPTION 'Bến xe chưa có toạ độ (cần để ước lượng km/phút): %', missing;
    END IF;

    SELECT string_agg(DISTINCT l.name ->> 'vi', ', ') INTO missing
    FROM locations l
    WHERE l.type = 'bus_station' AND NOT l.deleted
      AND l.administrative_area_id IS NULL
      AND l.name ->> 'vi' IN (SELECT origin FROM route_defs UNION SELECT dest FROM route_defs);
    IF missing IS NOT NULL THEN
        RAISE EXCEPTION 'Bến đầu/cuối chưa có administrative_area_id: %', missing;
    END IF;

    SELECT string_agg(DISTINCT d.op_code, ', ') INTO missing
    FROM route_defs d
    WHERE NOT EXISTS (SELECT 1 FROM operators o WHERE o.code = d.op_code AND NOT o.deleted);
    IF missing IS NOT NULL THEN
        RAISE EXCEPTION 'Thiếu operator: %', missing;
    END IF;
END $$;

-- ---------------------------------------------------------------------
-- 2) Mở rộng thành 2 chiều (chiều về = đảo thứ tự điểm dừng)
-- ---------------------------------------------------------------------
CREATE TEMP TABLE route_rows ON COMMIT DROP AS
SELECT
    row_number() OVER (ORDER BY s.op_code, s.origin, s.dest, s.dir) AS rid,
    s.op_code,
    s.distance_km,
    s.duration_min,
    CASE WHEN s.dir = 1 THEN s.origin_city    ELSE s.dest_city    END AS o_city,
    CASE WHEN s.dir = 1 THEN s.origin_city_en ELSE s.dest_city_en END AS o_city_en,
    CASE WHEN s.dir = 1 THEN s.dest_city      ELSE s.origin_city  END AS d_city,
    CASE WHEN s.dir = 1 THEN s.dest_city_en   ELSE s.origin_city_en END AS d_city_en,
    s.stop_names
FROM (
    SELECT
        d.*,
        g.dir,
        CASE
            WHEN g.dir = 1 THEN ARRAY[d.origin] || d.via || ARRAY[d.dest]
            ELSE (SELECT array_agg(u.x ORDER BY u.ord DESC)
                    FROM unnest(ARRAY[d.origin] || d.via || ARRAY[d.dest]) WITH ORDINALITY AS u(x, ord))
        END AS stop_names
    FROM route_defs d
    CROSS JOIN generate_series(1, 2) AS g(dir)
) s;

-- ---------------------------------------------------------------------
-- 3) Điểm dừng của từng tuyến + ước lượng km / phút tính từ bến đi
-- ---------------------------------------------------------------------
CREATE TEMP TABLE route_stop_rows ON COMMIT DROP AS
WITH s AS (
    SELECT r.rid, u.ord, cardinality(r.stop_names) AS n,
           l.id AS location_id, l.name, l.address, l.latitude, l.longitude,
           CASE WHEN u.ord IN (1, cardinality(r.stop_names)) THEN 0 ELSE 10 END AS stop_min
    FROM route_rows r
    CROSS JOIN LATERAL unnest(r.stop_names) WITH ORDINALITY AS u(nm, ord)
    JOIN locations l ON l.type = 'bus_station' AND NOT l.deleted AND l.name ->> 'vi' = u.nm
), seg AS (
    -- Khoảng cách đường chim bay (haversine, km) từ bến trước tới bến này.
    SELECT s.*,
           COALESCE(2 * 6371 * asin(sqrt(
               power(sin(radians(s.latitude - lag(s.latitude) OVER w) / 2), 2)
             + cos(radians(lag(s.latitude) OVER w)) * cos(radians(s.latitude))
             * power(sin(radians(s.longitude - lag(s.longitude) OVER w) / 2), 2)
           )), 0) AS seg_km
    FROM s
    WINDOW w AS (PARTITION BY s.rid ORDER BY s.ord)
)
SELECT seg.*,
       sum(seg_km) OVER (PARTITION BY rid ORDER BY ord)
         / sum(seg_km) OVER (PARTITION BY rid)                     AS frac,         -- tỉ lệ quãng đường đã đi
       sum(stop_min) OVER (PARTITION BY rid)                       AS total_stop_min,
       COALESCE(sum(stop_min) OVER (PARTITION BY rid ORDER BY ord
                ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING), 0) AS stop_min_before
FROM seg;

-- Snapshot bến đầu/cuối (city = tỉnh/thành chứa phường/xã của bến).
CREATE TEMP TABLE endpoint_snapshots ON COMMIT DROP AS
SELECT l.id AS location_id,
       jsonb_build_object(
           'location_id', l.id,
           'name',        l.name,
           'city',        jsonb_build_object('id', c.id, 'name', c.name),
           'address',     l.address,
           'lat',         l.latitude,
           'lng',         l.longitude
       ) AS snapshot
FROM locations l
JOIN administrative_areas w ON w.id = l.administrative_area_id
JOIN administrative_areas c ON c.id = COALESCE(w.parent_id, w.id)
WHERE l.type = 'bus_station' AND NOT l.deleted;

-- ---------------------------------------------------------------------
-- 4) Insert routes + route_snapshot
-- ---------------------------------------------------------------------
INSERT INTO routes
    (operator_id, origin_id, destination_id, name, distance_km, estimated_duration_minutes,
     status, route_snapshot, created_at, updated_at)
SELECT
    o.id,
    lo.administrative_area_id,
    ld.administrative_area_id,
    jsonb_build_object('vi', r.o_city || ' - ' || r.d_city,
                       'en', r.o_city_en || ' - ' || r.d_city_en),
    r.distance_km,
    r.duration_min,
    'active',
    jsonb_build_object(
        'operator', jsonb_build_object(
            'id',       o.id,
            'code',     o.code,
            'name',     o.name,
            'brand',    o.brand,
            'hotline',  regexp_replace(o.hotline, '\s', '', 'g'),
            'logo_url', o.metadata ->> 'logo_url'
        ),
        'origin',      eo.snapshot,
        'destination', ed.snapshot,
        'stops', (
            SELECT jsonb_agg(
                       jsonb_build_object(
                           'order',       st.ord,
                           'location_id', st.location_id,
                           'name',        st.name,
                           'type',        CASE
                                              WHEN st.ord = 1 THEN 'pickup'
                                              WHEN st.ord = st.n THEN 'dropoff'
                                              ELSE 'rest'
                                          END,
                           'address',     st.address,
                           'lat',         st.latitude,
                           'lng',         st.longitude,
                           'offset_minutes',
                               round((r.duration_min - st.total_stop_min) * st.frac)::int + st.stop_min_before,
                           'distance_from_origin_km', round(r.distance_km * st.frac::numeric, 1),
                           'stop_duration_minutes',   st.stop_min
                       ) ORDER BY st.ord)
            FROM route_stop_rows st
            WHERE st.rid = r.rid
        )
    ),
    NOW(),
    NOW()
FROM route_rows r
JOIN operators o        ON o.code = r.op_code
JOIN locations lo       ON lo.type = 'bus_station' AND NOT lo.deleted
                       AND lo.name ->> 'vi' = r.stop_names[1]
JOIN locations ld       ON ld.type = 'bus_station' AND NOT ld.deleted
                       AND ld.name ->> 'vi' = r.stop_names[cardinality(r.stop_names)]
JOIN endpoint_snapshots eo ON eo.location_id = lo.id
JOIN endpoint_snapshots ed ON ed.location_id = ld.id
WHERE NOT EXISTS (
    SELECT 1 FROM routes x
    WHERE x.operator_id = o.id
      AND x.origin_id = lo.administrative_area_id
      AND x.destination_id = ld.administrative_area_id
      AND x.name ->> 'vi' = r.o_city || ' - ' || r.d_city
);

COMMIT;

-- Kiểm tra (kỳ vọng 46 tuyến = 23 tuyến x 2 chiều):
--   SELECT count(*) FROM routes WHERE metadata->>'seed' = 'routes_v1';
--   SELECT id, name->>'vi', distance_km, estimated_duration_minutes,
--          jsonb_array_length(route_snapshot->'stops') AS stops
--   FROM routes ORDER BY id LIMIT 10;
--   -- Xem tên phường/xã đi - đến của từng tuyến:
--   SELECT r.id, r.name->>'vi', ao.name->>'vi' AS origin_area, ad.name->>'vi' AS destination_area
--   FROM routes r
--   JOIN administrative_areas ao ON ao.id = r.origin_id
--   JOIN administrative_areas ad ON ad.id = r.destination_id
--   WHERE r.metadata->>'seed' = 'routes_v1';
--   -- Nếu muốn origin_id/destination_id là cấp tỉnh/thành: dùng ao.parent_id / ad.parent_id
--   -- (trong INSERT: thay lo.administrative_area_id bằng
--   --  (SELECT parent_id FROM administrative_areas WHERE id = lo.administrative_area_id)).
--   -- Tìm tuyến đi qua một bến mà không cần join:
--   SELECT id, name->>'vi' FROM routes
--   WHERE route_snapshot->'stops' @> '[{"location_id": 1}]';