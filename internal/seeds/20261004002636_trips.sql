-- =====================================================================
-- Seed data: trips (chuyến xe, 14 ngày kể từ ngày mai)
-- Phụ thuộc (phải chạy trước): vehicles, routes (đã có route_snapshot)
--
--   * Mỗi xe chạy con thoi trên 1 cặp tuyến (A->B rồi B->A) của chính nhà xe đó, nên một xe
--     không bao giờ chạy 2 chuyến chồng giờ và luôn xuất bến tại nơi vừa trả khách.
--   * Gán xe cho cặp tuyến: cặp dài nhất được ưu tiên xe giường nằm; xe thứ k của nhà xe nhận
--     cặp (k mod số cặp). Nhà xe có ít xe hơn số cặp tuyến -> các cặp ngắn nhất không có chuyến.
--     Nhiều xe cùng một cặp: xe chẵn xuất phát từ A, xe lẻ từ B, lệch nhau 2 tiếng.
--   * Lịch: giờ xuất bến làm tròn lên mốc 30 phút, chỉ trong khung 05:00 - 22:30 (giờ VN).
--     Chuyến đầu: tuyến >= 8 tiếng xuất bến 19:00 (xe đêm), còn lại 06:00.
--     Nghỉ ở bến cuối trước chuyến về: 60 phút (tuyến <= 4 tiếng), 120 phút (tuyến dài hơn).
--   * arrival_time = departure_time + routes.estimated_duration_minutes.
--   * base_price (VND, làm tròn 10.000): 40.000 + 850đ/km cho 500 km đầu + 600đ/km phần còn lại,
--     nhân hệ số loại xe (limousine 1.3, sleeper 1.0, seat 0.85); xuất bến thứ 6 - CN thêm 10%.
--   * Chưa có booking/giữ chỗ: available_seats_count = vehicles.total_seats, booked/held = 0.
--     status = 'scheduled' (mọi chuyến đều ở tương lai).
--   * searchable_snapshot: dữ liệu hiển thị kết quả tìm kiếm, lấy từ route_snapshot + vehicles:
--       { "operator": {...}, "origin": {...}, "destination": {...},   -- giống route_snapshot
--         "route":   {id, name{vi,en}, distance_km, duration_minutes, stops_count},
--         "vehicle": {id, type, total_seats, amenities} }
--   * Chỉ sinh cho xe chưa có chuyến nào -> chạy lại không bị trùng / chồng lịch.
--     Muốn sinh lại: DELETE FROM trips WHERE vehicle_id = ...; rồi chạy lại file.
-- =====================================================================
BEGIN;

-- ---------------------------------------------------------------------
-- 1) Kiểm tra dữ liệu phụ thuộc
-- ---------------------------------------------------------------------
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM vehicles WHERE NOT deleted AND status = 'active') THEN
        RAISE EXCEPTION 'Chưa có vehicles, hãy chạy seed vehicles trước';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM routes
        WHERE NOT deleted AND jsonb_array_length(COALESCE(route_snapshot -> 'stops', '[]'::jsonb)) > 0
    ) THEN
        RAISE EXCEPTION 'Chưa có tuyến nào có route_snapshot->stops, hãy chạy seed routes trước';
    END IF;
END $$;

-- Mốc giờ xuất bến kế tiếp (giờ VN, không múi giờ): làm tròn lên 30 phút, trong khung 05:00 - 22:30.
CREATE OR REPLACE FUNCTION pg_temp.next_departure_slot(t timestamp) RETURNS timestamp
LANGUAGE sql IMMUTABLE AS $$
    SELECT CASE
               WHEN s::time < '05:00' THEN s::date + time '05:00'
               WHEN s::time > '22:30' THEN s::date + 1 + time '05:00'
               ELSE s
           END
    FROM (SELECT date_trunc('hour', t)
                 + ceil(extract(epoch FROM t - date_trunc('hour', t)) / 1800) * interval '30 minutes' AS s) x
$$;

-- ---------------------------------------------------------------------
-- 2) Ghép 2 chiều của cùng một tuyến thành cặp (A->B, B->A)
-- ---------------------------------------------------------------------
CREATE TEMP TABLE route_pairs ON COMMIT DROP AS
SELECT
    f.operator_id,
    f.id AS fwd_route_id,
    b.id AS back_route_id,
    f.estimated_duration_minutes AS duration_min,
    row_number() OVER (PARTITION BY f.operator_id ORDER BY f.distance_km DESC, f.id) - 1 AS pair_idx,
    count(*)     OVER (PARTITION BY f.operator_id) AS pair_count
FROM routes f
JOIN routes b
  ON b.operator_id = f.operator_id
 AND b.id > f.id
 AND NOT b.deleted
 AND b.route_snapshot -> 'origin' ->> 'location_id' = f.route_snapshot -> 'destination' ->> 'location_id'
 AND b.route_snapshot -> 'destination' ->> 'location_id' = f.route_snapshot -> 'origin' ->> 'location_id'
WHERE NOT f.deleted AND f.status = 'active' AND b.status = 'active'
  AND f.estimated_duration_minutes > 0;

-- ---------------------------------------------------------------------
-- 3) Gán xe (chưa có chuyến nào) vào cặp tuyến
-- ---------------------------------------------------------------------
CREATE TEMP TABLE vehicle_assignments ON COMMIT DROP AS
SELECT
    v.vehicle_id,
    p.fwd_route_id,
    p.back_route_id,
    p.duration_min,
    v.vehicle_idx / p.pair_count AS slot_idx          -- thứ tự của xe trong cùng một cặp tuyến
FROM (
    SELECT
        v.id AS vehicle_id,
        v.operator_id,
        row_number() OVER (PARTITION BY v.operator_id
                           ORDER BY (v.type = 'sleeper') DESC, v.total_seats DESC, v.id) - 1 AS vehicle_idx
    FROM vehicles v
    WHERE NOT v.deleted AND v.status = 'active'
) v
JOIN route_pairs p
  ON p.operator_id = v.operator_id
 AND p.pair_idx = v.vehicle_idx % p.pair_count
WHERE NOT EXISTS (SELECT 1 FROM trips t WHERE t.vehicle_id = v.vehicle_id);

-- ---------------------------------------------------------------------
-- 4) Sinh lịch con thoi cho từng xe trong 14 ngày (giờ VN, không múi giờ)
-- ---------------------------------------------------------------------
CREATE TEMP TABLE vehicle_legs ON COMMIT DROP AS
WITH RECURSIVE legs AS (
    SELECT
        a.vehicle_id,
        a.fwd_route_id,
        a.back_route_id,
        a.duration_min,
        CASE WHEN a.duration_min <= 240 THEN 60 ELSE 120 END AS layover_min,
        a.slot_idx % 2 = 0 AS is_fwd,
        (current_date + 1)
            + CASE WHEN a.duration_min >= 480 THEN time '19:00' ELSE time '06:00' END
            + (a.slot_idx / 2) * interval '2 hours' AS departure_local
    FROM vehicle_assignments a
    UNION ALL
    SELECT
        l.vehicle_id, l.fwd_route_id, l.back_route_id, l.duration_min, l.layover_min,
        NOT l.is_fwd,
        pg_temp.next_departure_slot(
            l.departure_local + (l.duration_min + l.layover_min) * interval '1 minute')
    FROM legs l
    WHERE l.departure_local < (current_date + 15)::timestamp
)
SELECT
    vehicle_id,
    CASE WHEN is_fwd THEN fwd_route_id ELSE back_route_id END AS route_id,
    departure_local
FROM legs
WHERE departure_local < (current_date + 15)::timestamp;

-- ---------------------------------------------------------------------
-- 5) Insert trips
-- ---------------------------------------------------------------------
INSERT INTO trips
    (operator_id, route_id, vehicle_id, origin_id, destination_id, departure_time, arrival_time,
     base_price, status, available_seats_count, booked_seats_count, held_seats_count,
     searchable_snapshot, created_at, updated_at)
SELECT
    r.operator_id,
    r.id,
    v.id,
    r.origin_id,
    r.destination_id,
    l.departure_local AT TIME ZONE 'Asia/Ho_Chi_Minh',
    (l.departure_local + r.estimated_duration_minutes * interval '1 minute') AT TIME ZONE 'Asia/Ho_Chi_Minh',
    round(
        (40000 + 850 * least(r.distance_km, 500) + 600 * greatest(r.distance_km - 500, 0))
        * CASE v.type WHEN 'limousine' THEN 1.3 WHEN 'seat' THEN 0.85 ELSE 1.0 END
        * CASE WHEN extract(isodow FROM l.departure_local) IN (5, 6, 7) THEN 1.1 ELSE 1.0 END
        / 10000
    )::bigint * 10000,
    'scheduled',
    v.total_seats,
    0,
    0,
    jsonb_build_object(
        'operator',    r.route_snapshot -> 'operator',
        'origin',      r.route_snapshot -> 'origin',
        'destination', r.route_snapshot -> 'destination',
        'route', jsonb_build_object(
            'id',               r.id,
            'name',             r.name,
            'distance_km',      r.distance_km,
            'duration_minutes', r.estimated_duration_minutes,
            'stops_count',      jsonb_array_length(r.route_snapshot -> 'stops')
        ),
        'vehicle', jsonb_build_object(
            'id',          v.id,
            'type',        v.type,
            'total_seats', v.total_seats,
            'amenities',   v.amenities
        )
    ),
    NOW(),
    NOW()
FROM vehicle_legs l
JOIN routes r   ON r.id = l.route_id
JOIN vehicles v ON v.id = l.vehicle_id;

COMMIT;

-- Kiểm tra:
--   -- Số chuyến theo tuyến:
--   SELECT r.name->>'vi', count(*) FROM trips t JOIN routes r ON r.id = t.route_id
--   GROUP BY 1 ORDER BY 1;
--   -- Không xe nào chạy chồng giờ (kỳ vọng 0 dòng):
--   SELECT a.id, b.id FROM trips a JOIN trips b
--     ON b.vehicle_id = a.vehicle_id AND b.id <> a.id
--    AND b.departure_time < a.arrival_time AND b.arrival_time > a.departure_time;
--   -- Tìm chuyến (đúng index idx_trips_origin_destination_departure):
--   SELECT id, departure_time AT TIME ZONE 'Asia/Ho_Chi_Minh', base_price, available_seats_count
--   FROM trips
--   WHERE origin_id = ... AND destination_id = ...
--     AND departure_time >= now() AND departure_time < now() + interval '1 day'
--   ORDER BY departure_time;
