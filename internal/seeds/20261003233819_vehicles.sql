BEGIN;

INSERT INTO vehicles
    (operator_id, plate_number, type, total_seats, status, seat_layout, amenities, created_at, updated_at)
SELECT
    o.id,
    '00X-' || lpad(x.rn::text, 5, '0'),
    x.type,
    x.total_seats,
    'active',
    -- seat_layout: {"floors":[{"floor":1,"seats":[{"seat_no","row","col"}...]}...]}}
    -- sleeper = 2 tầng (A/B), còn lại = 1 tầng (A)
    jsonb_build_object(
        'floors', (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'floor', fl,
                    'seats', (
                        SELECT jsonb_agg(
                            jsonb_build_object(
                                'seat_no', chr(64 + fl) || g,
                                'row', CASE WHEN g <= x.spf - 5
                                            THEN (g - 1) / x.per_row + 1
                                            ELSE (x.spf - 5) / x.per_row + 1 END,
                                'col', CASE WHEN g <= x.spf - 5
                                            THEN x.col_map[((g - 1) % x.per_row) + 1]
                                            ELSE g - (x.spf - 5) END
                            ) ORDER BY g)
                        FROM generate_series(1, x.spf) AS g
                    )
                ) ORDER BY fl)
            FROM generate_series(1, x.floors) AS fl
        )
    ),
    -- amenities theo loại xe
    CASE x.type
        WHEN 'sleeper' THEN jsonb_build_object(
            'wifi', true, 'usb_charger', true, 'air_conditioner', true,
            'blanket', true, 'water', true, 'tv', false, 'toilet', false)
        WHEN 'limousine' THEN jsonb_build_object(
            'wifi', true, 'usb_charger', true, 'air_conditioner', true,
            'blanket', true, 'water', true, 'tv', true, 'toilet', false)
        ELSE jsonb_build_object(
            'wifi', false, 'usb_charger', false, 'air_conditioner', true,
            'blanket', false, 'water', true, 'tv', false, 'toilet', false)
    END,
    NOW(),
    NOW()
FROM (
    SELECT
        f.operator_code,
        f.type,
        f.total_seats,
        f.floors,
        f.total_seats / f.floors AS spf,          -- số ghế mỗi tầng
        f.per_row,
        f.col_map,
        row_number() OVER (ORDER BY f.operator_code, f.type, f.total_seats, s.n) AS rn
    FROM (VALUES
        -- operator_code,    type,        total_seats, qty, floors, per_row, col_map
        ('phuong-trang',     'sleeper',   40, 3, 2, 3, ARRAY[1,3,5]),
        ('phuong-trang',     'limousine', 33, 2, 1, 2, ARRAY[1,5]),
        ('hoang-long',       'sleeper',   40, 3, 2, 3, ARRAY[1,3,5]),
        ('kumho-samco',      'limousine', 33, 3, 1, 2, ARRAY[1,5]),
        ('kumho-samco',      'sleeper',   40, 1, 2, 3, ARRAY[1,3,5]),
        ('thanh-cong',       'sleeper',   40, 2, 2, 3, ARRAY[1,3,5]),
        ('thanh-cong',       'seat',      45, 1, 1, 4, ARRAY[1,2,4,5]),
        ('quang-dung',       'limousine', 33, 2, 1, 2, ARRAY[1,5]),
        ('son-tung',         'sleeper',   34, 2, 2, 2, ARRAY[1,5]),
        ('mai-linh-express', 'seat',      45, 2, 1, 4, ARRAY[1,2,4,5]),
        ('mai-linh-express', 'limousine', 33, 1, 1, 2, ARRAY[1,5]),
        ('ha-son',           'sleeper',   40, 2, 2, 3, ARRAY[1,3,5])
    ) AS f(operator_code, type, total_seats, qty, floors, per_row, col_map)
    CROSS JOIN LATERAL generate_series(1, f.qty) AS s(n)
) x
JOIN operators o ON o.code = x.operator_code
ON CONFLICT (plate_number) DO NOTHING;

COMMIT;

-- Kiểm tra:
--   SELECT plate_number, jsonb_array_length(seat_layout->'floors') AS floors FROM vehicles LIMIT 5;
--   -- tổng ghế trong layout phải bằng total_seats (kỳ vọng 0 dòng):
--   SELECT v.plate_number FROM vehicles v
--   WHERE v.plate_number LIKE '00X-%'
--     AND (SELECT sum(jsonb_array_length(f->'seats')) FROM jsonb_array_elements(v.seat_layout->'floors') f) <> v.total_seats;