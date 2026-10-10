-- =====================================================================
-- Seed data: operators (nhà xe khách)
-- Nguồn: trang chính thức/ứng dụng của nhà xe và các trang tổng hợp (tra cứu 10/2026).
--
-- Nguyên tắc:
--   * Chỉ nhập hotline khi có nguồn xác nhận; không có nguồn -> NULL.
--   * email để NULL toàn bộ (chưa có nguồn công khai đáng tin để xác nhận).
--   * name = tên thương hiệu thường dùng; tên pháp nhân chỉ dùng khi có nguồn.
--   * metadata chỉ chứa thông tin đã có nguồn (website, năm thành lập, trụ sở).
-- =====================================================================
BEGIN;

INSERT INTO operators (name, code, brand, hotline, email, status, metadata, created_at, updated_at)
VALUES
-- Phương Trang: hotline 1900 6067 (nhiều nguồn, kể cả trang CSKH); thành lập 2001; website futabus.vn
(jsonb_build_object('vi', 'Công ty Cổ phần Xe khách Phương Trang - FUTA Bus Lines', 'en', 'Phuong Trang - FUTA Bus Lines'),
 'phuong-trang', 'FUTA Bus Lines', '1900 6067', NULL, 'active',
 jsonb_build_object('website', 'futabus.vn', 'founded_year', 2001), NOW(), NOW()),

-- Hoàng Long: hotline 1900 988 920; số 02253 920 920 xuất hiện ở tuyến Hải Phòng
(jsonb_build_object('vi', 'Nhà xe Hoàng Long', 'en', 'Hoang Long Bus'),
 'hoang-long', 'Hoang Long', '1900 988 920', NULL, 'active',
 jsonb_build_object('hotline_secondary', '0225 3920 920'), NOW(), NOW()),

-- Kumho Samco Buslines: tên pháp nhân có nguồn; hotline 1900 6065 (2 nguồn)
(jsonb_build_object('vi', 'Công ty TNHH Vận tải Kumho Samco Buslines', 'en', 'Kumho Samco Buslines'),
 'kumho-samco', 'Kumho Samco Buslines', '1900 6065', NULL, 'active',
 '{}'::jsonb, NOW(), NOW()),

-- Thành Công Buslines: hotline 1900 6962; website thanhcongbus.vn
(jsonb_build_object('vi', 'Nhà xe Thành Công', 'en', 'Thanh Cong Buslines'),
 'thanh-cong', 'Thanh Cong Buslines', '1900 6962', NULL, 'active',
 jsonb_build_object('website', 'thanhcongbus.vn'), NOW(), NOW()),

-- Quang Dũng Limousine: hotline chính 0971 510 510; trụ sở Quảng Trị
(jsonb_build_object('vi', 'Nhà xe Quang Dũng', 'en', 'Quang Dung Limousine'),
 'quang-dung', 'Quang Dung Limousine', '0971 510 510', NULL, 'active',
 jsonb_build_object('headquarters', 'Quảng Trị',
                    'hotline_secondary', jsonb_build_array('0936 375 375', '0971 367 367', '0971 441 441', '0971 442 442')),
 NOW(), NOW()),

-- Sơn Tùng: tổng đài 1900 9696 71 / 1900 5555 23; trụ sở Bình Định (nay thuộc Gia Lai)
(jsonb_build_object('vi', 'Nhà xe Sơn Tùng', 'en', 'Son Tung Bus'),
 'son-tung', NULL, '1900 9696 71', NULL, 'active',
 jsonb_build_object('hotline_secondary', '1900 5555 23', 'headquarters', '234 Tây Sơn, Quy Nhơn'), NOW(), NOW()),

-- Mai Linh Express: thành viên Tập đoàn Mai Linh; chưa xác nhận được hotline riêng
(jsonb_build_object('vi', 'Mai Linh Express', 'en', 'Mai Linh Express'),
 'mai-linh-express', 'Mai Linh Express', NULL, NULL, 'active',
 jsonb_build_object('parent_group', 'Tập đoàn Mai Linh'), NOW(), NOW()),

-- Hà Sơn: thành lập 2008, trụ sở Sơn Tây (Hà Nội); chưa xác nhận hotline riêng
(jsonb_build_object('vi', 'Nhà xe Hà Sơn', 'en', 'Ha Son Bus'),
 'ha-son', NULL, NULL, NULL, 'active',
 jsonb_build_object('founded_year', 2008, 'headquarters', 'Sơn Tây, Hà Nội'), NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

COMMIT;