-- =====================================================================
-- Seed data: route_stops (điểm dừng của từng tuyến)
-- Phụ thuộc (phải chạy trước): 20261004002627_routes.sql
--
--   * Sinh từ routes.route_snapshot->'stops' nên thứ tự, location_id, loại điểm dừng và
--     số phút luôn khớp với snapshot của tuyến.
--   * stop_type: 'pickup' (bến đi), 'rest' (bến dọc đường), 'dropoff' (bến đến).
--   * arrival_offset_minutes   = phút xe tới bến, tính từ giờ xuất bến (bến đi: NULL).
--     departure_offset_minutes = arrival + thời gian dừng (bến đến: NULL).
--   * public_note {vi,en}: hướng dẫn cho hành khách; internal_note: ghi chú cho tài xế/điều hành.
--   * metadata để mặc định '{}'.
--   * Chỉ sinh cho tuyến chưa có route_stops nào -> chạy lại không bị trùng, không đụng dữ liệu
--     đã nhập tay. Muốn sinh lại cho một tuyến: xoá route_stops của tuyến đó rồi chạy lại file.
-- =====================================================================
BEGIN;

-- ---------------------------------------------------------------------
-- 1) Kiểm tra dữ liệu phụ thuộc: chưa có tuyến nào có snapshot stops thì dừng
-- ---------------------------------------------------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM routes
        WHERE NOT deleted AND jsonb_array_length(COALESCE(route_snapshot -> 'stops', '[]'::jsonb)) > 0
    ) THEN
        RAISE EXCEPTION 'Chưa có tuyến nào có route_snapshot->stops, hãy chạy seed routes trước';
    END IF;
END $$;

-- ---------------------------------------------------------------------
-- 2) Insert route_stops từ snapshot
-- ---------------------------------------------------------------------
INSERT INTO route_stops
    (route_id, location_id, stop_order, stop_type,
     arrival_offset_minutes, departure_offset_minutes, public_note, internal_note)
SELECT
    r.id,
    (s ->> 'location_id')::int,
    (s ->> 'order')::int,
    s ->> 'type',
    CASE WHEN s ->> 'type' = 'pickup' THEN NULL
         ELSE (s ->> 'offset_minutes')::int END,
    CASE WHEN s ->> 'type' = 'dropoff' THEN NULL
         ELSE (s ->> 'offset_minutes')::int + (s ->> 'stop_duration_minutes')::int END,
    CASE s ->> 'type'
        WHEN 'pickup' THEN jsonb_build_object(
            'vi', 'Vui lòng có mặt tại bến trước giờ xuất bến 30 phút để làm thủ tục lên xe.',
            'en', 'Please arrive at the station 30 minutes before departure to check in.')
        WHEN 'dropoff' THEN jsonb_build_object(
            'vi', 'Điểm trả khách cuối. Vui lòng kiểm tra hành lý trước khi xuống xe.',
            'en', 'Final drop-off. Please check your belongings before leaving the bus.')
        ELSE jsonb_build_object(
            'vi', format('Xe dừng khoảng %s phút để đón/trả khách và nghỉ ngơi.', s ->> 'stop_duration_minutes'),
            'en', format('The bus stops for about %s minutes for pick-up/drop-off and a short break.',
                         s ->> 'stop_duration_minutes'))
    END,
    CASE s ->> 'type'
        WHEN 'pickup'  THEN 'Đối chiếu danh sách hành khách và vé trước khi xuất bến.'
        WHEN 'dropoff' THEN 'Kiểm tra xe sau chuyến, báo điều hành nếu có hành lý thất lạc.'
        ELSE 'Xác nhận khách lên/xuống tại bến; báo điều hành nếu trễ quá 15 phút so với lịch.'
    END
FROM routes r
CROSS JOIN LATERAL jsonb_array_elements(r.route_snapshot -> 'stops') AS s
WHERE NOT r.deleted
  AND NOT EXISTS (SELECT 1 FROM route_stops x WHERE x.route_id = r.id);

COMMIT;

-- Kiểm tra:
--   -- Số điểm dừng mỗi tuyến phải bằng số phần tử trong snapshot:
--   SELECT r.id, r.name->>'vi', count(rs.id) AS stops,
--          jsonb_array_length(r.route_snapshot->'stops') AS snapshot_stops
--   FROM routes r LEFT JOIN route_stops rs ON rs.route_id = r.id
--   GROUP BY r.id ORDER BY r.id;
--   -- Lịch dừng của một tuyến:
--   SELECT rs.stop_order, l.name->>'vi', rs.stop_type,
--          rs.arrival_offset_minutes, rs.departure_offset_minutes
--   FROM route_stops rs JOIN locations l ON l.id = rs.location_id
--   WHERE rs.route_id = 1 ORDER BY rs.stop_order;
