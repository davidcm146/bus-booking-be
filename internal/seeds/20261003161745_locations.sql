-- =====================================================================
-- Seed data: locations (bus_station, rest_area)
-- Phụ thuộc: administrative_areas seed (hiệu lực 01/07/2025, 34 tỉnh/thành)
--
-- Cách map administrative_area_id:
--   (province_code, ward_code) -> administrative_areas phường/xã (parent_id = id tỉnh)
--
-- LƯU Ý CHẤT LƯỢNG DỮ LIỆU:
--   * metadata để NULL cho toàn bộ dòng (sẽ migrate sau).
--   * bus_station: tên + địa chỉ là bến xe có thật; toạ độ là xấp xỉ (sai số ~100m-1km),
--     cần đối chiếu bằng Google Maps.
--   * Ward được chọn theo vị trí bến xe sau sáp nhập 01/07/2025; một số bến nằm sát ranh
--     phường (Miền Đông cũ, Chợ Lớn, Vĩnh Long, Ninh Bình...) -> nên rà lại.
--   * rest_area: chỉ nhập các trạm dừng nghỉ cao tốc có thông tin công bố chính thức
--     (Cục Đường bộ VN). Chưa có toạ độ chính xác -> lat/lng NULL, status 'inactive'.
-- =====================================================================
BEGIN;

-- ---------------------------------------------------------------------
-- 1) BUS STATIONS
-- ---------------------------------------------------------------------
INSERT INTO locations
    (administrative_area_id, name, type, address, latitude, longitude, status, created_at, updated_at)
SELECT
    a.id,
    NULL,
    jsonb_build_object('vi', v.name_vi, 'en', v.name_en),
    'bus_station',
    v.address,
    v.latitude,
    v.longitude,
    'active',
    NOW(),
    NOW()
FROM (VALUES
-- ward_code, province_code, name_vi, name_en, address, lat, lng
-- ===== 01 Hà Nội =====
('00592','01','Bến xe Mỹ Đình','My Dinh Bus Station','20 Phạm Hùng, Mỹ Đình, Hà Nội',21.028600,105.778800),
('00322','01','Bến xe Giáp Bát','Giap Bat Bus Station','Số 1 Giải Phóng, Hoàng Mai, Hà Nội',20.980900,105.841300),
('00337','01','Bến xe Nước Ngầm','Nuoc Ngam Bus Station','Km8 Giải Phóng, Hoàng Liệt, Hà Nội',20.962600,105.844700),
('09562','01','Bến xe Yên Nghĩa','Yen Nghia Bus Station','Km9 Quang Trung, Yên Nghĩa, Hà Nội',20.950600,105.743300),
('00118','01','Bến xe Gia Lâm','Gia Lam Bus Station','9 Ngô Gia Khảm, Long Biên, Hà Nội',21.047500,105.878600),
-- ===== 04 Cao Bằng =====
('01273','04','Bến xe Cao Bằng','Cao Bang Bus Station','Phường Thục Phán, Cao Bằng',22.666000,106.258000),
-- ===== 08 Tuyên Quang (gồm Hà Giang cũ) =====
('00694','08','Bến xe Hà Giang','Ha Giang Bus Station','Phường Hà Giang 1, Tuyên Quang',22.824000,104.984000),
('02215','08','Bến xe Tuyên Quang','Tuyen Quang Bus Station','Phường Minh Xuân, Tuyên Quang',21.823000,105.212000),
-- ===== 11 Điện Biên =====
('03127','11','Bến xe Điện Biên Phủ','Dien Bien Phu Bus Station','Phường Điện Biên Phủ, Điện Biên',21.386000,103.016000),
-- ===== 12 Lai Châu =====
('03388','12','Bến xe Lai Châu','Lai Chau Bus Station','Phường Đoàn Kết, Lai Châu',22.396000,103.462000),
-- ===== 14 Sơn La =====
('03646','14','Bến xe Sơn La','Son La Bus Station','Phường Tô Hiệu, Sơn La',21.326000,103.914000),
('03980','14','Bến xe Mộc Châu','Moc Chau Bus Station','Phường Mộc Châu, Sơn La',20.838000,104.639000),
-- ===== 15 Lào Cai (gồm Yên Bái cũ) =====
('02647','15','Bến xe Lào Cai','Lao Cai Bus Station','Phường Lào Cai, Lào Cai',22.478000,103.976000),
('03006','15','Bến xe Sa Pa','Sa Pa Bus Station','Phường Sa Pa, Lào Cai',22.336000,103.844000),
('04252','15','Bến xe Yên Bái','Yen Bai Bus Station','Phường Yên Bái, Lào Cai',21.716000,104.899000),
-- ===== 19 Thái Nguyên (gồm Bắc Kạn cũ) =====
('05443','19','Bến xe Thái Nguyên','Thai Nguyen Bus Station','Phường Phan Đình Phùng, Thái Nguyên',21.588000,105.835000),
('01843','19','Bến xe Bắc Kạn','Bac Kan Bus Station','Phường Bắc Kạn, Thái Nguyên',22.147000,105.834000),
-- ===== 20 Lạng Sơn =====
('05977','20','Bến xe Lạng Sơn','Lang Son Bus Station','Phường Đông Kinh, Lạng Sơn',21.854000,106.758000),
-- ===== 22 Quảng Ninh =====
('06673','22','Bến xe Bãi Cháy','Bai Chay Bus Station','Phường Bãi Cháy, Quảng Ninh',20.956500,107.044000),
('06712','22','Bến xe Móng Cái','Mong Cai Bus Station','Phường Móng Cái 1, Quảng Ninh',21.528000,107.964300),
('06811','22','Bến xe Uông Bí','Uong Bi Bus Station','Phường Uông Bí, Quảng Ninh',21.039000,106.779000),
('06793','22','Bến xe Cẩm Phả','Cam Pha Bus Station','Phường Cẩm Phả, Quảng Ninh',21.010000,107.290000),
('07093','22','Bến xe Đông Triều','Dong Trieu Bus Station','Phường Đông Triều, Quảng Ninh',21.083000,106.512000),
('06994','22','Bến xe Vân Đồn','Van Don Bus Station','Đặc khu Vân Đồn, Quảng Ninh',21.073000,107.420000),
-- ===== 24 Bắc Ninh (gồm Bắc Giang cũ) =====
('09187','24','Bến xe Bắc Ninh','Bac Ninh Bus Station','Phường Kinh Bắc, Bắc Ninh',21.187000,106.072000),
('07210','24','Bến xe Bắc Giang','Bac Giang Bus Station','Phường Bắc Giang, Bắc Ninh',21.279000,106.199000),
-- ===== 25 Phú Thọ (gồm Hòa Bình, Vĩnh Phúc cũ) =====
('07900','25','Bến xe Việt Trì','Viet Tri Bus Station','Phường Việt Trì, Phú Thọ',21.322300,105.401000),
('08707','25','Bến xe Vĩnh Yên','Vinh Yen Bus Station','Phường Vĩnh Yên, Phú Thọ',21.309000,105.604000),
('04795','25','Bến xe Hòa Bình','Hoa Binh Bus Station','Phường Hòa Bình, Phú Thọ',20.813900,105.339400),
-- ===== 33 Hưng Yên (gồm Thái Bình cũ) =====
('11953','33','Bến xe Hưng Yên','Hung Yen Bus Station','Phường Phố Hiến, Hưng Yên',20.646200,106.053500),
('13225','33','Bến xe Thái Bình','Thai Binh Bus Station','Phường Thái Bình, Hưng Yên',20.444000,106.339000),
-- ===== 37 Ninh Bình (gồm Hà Nam, Nam Định cũ) =====
('14359','37','Bến xe Ninh Bình','Ninh Binh Bus Station','Phường Nam Hoa Lư, Ninh Bình',20.251800,105.972200),
('13669','37','Bến xe Nam Định','Nam Dinh Bus Station','Phường Nam Định, Ninh Bình',20.417000,106.170000),
('13285','37','Bến xe Phủ Lý','Phu Ly Bus Station','Phường Phủ Lý, Ninh Bình',20.540000,105.915000),
('14362','37','Bến xe Tam Điệp','Tam Diep Bus Station','Phường Tam Điệp, Ninh Bình',20.156000,105.912000),
-- ===== 38 Thanh Hóa =====
('14797','38','Bến xe phía Bắc Thanh Hóa','Thanh Hoa North Bus Station','Phường Hạc Thành, Thanh Hóa',19.833000,105.771000),
('16531','38','Bến xe Sầm Sơn','Sam Son Bus Station','Phường Sầm Sơn, Thanh Hóa',19.739000,105.904000),
('14812','38','Bến xe Bỉm Sơn','Bim Son Bus Station','Phường Bỉm Sơn, Thanh Hóa',20.085000,105.862000),
-- ===== 40 Nghệ An =====
('16702','40','Bến xe phía Bắc Vinh','Vinh North Bus Station','Phường Vinh Phú, Nghệ An',18.709000,105.695000),
('16708','40','Bến xe phía Nam Vinh','Vinh South Bus Station','Phường Vinh Lộc, Nghệ An',18.624000,105.715000),
('16732','40','Bến xe Cửa Lò','Cua Lo Bus Station','Phường Cửa Lò, Nghệ An',18.809000,105.725000),
('16939','40','Bến xe Thái Hòa','Thai Hoa Bus Station','Phường Thái Hòa, Nghệ An',19.338000,105.444000),
-- ===== 42 Hà Tĩnh =====
('18073','42','Bến xe Hà Tĩnh','Ha Tinh Bus Station','Phường Thành Sen, Hà Tĩnh',18.340000,105.902000),
-- ===== 44 Quảng Trị (gồm Quảng Bình cũ) =====
('18880','44','Bến xe Đồng Hới','Dong Hoi Bus Station','Phường Đồng Hới, Quảng Trị',17.474000,106.603000),
('19333','44','Bến xe Đông Hà','Dong Ha Bus Station','Phường Đông Hà, Quảng Trị',16.816000,107.098000),
('19009','44','Bến xe Ba Đồn','Ba Don Bus Station','Phường Ba Đồn, Quảng Trị',17.756000,106.421000),
-- ===== 46 Huế =====
('19815','46','Bến xe phía Nam Huế','Hue South Bus Station','97 An Dương Vương, Phường An Cựu, Huế',16.433000,107.560000),
-- ===== 48 Đà Nẵng (gồm Quảng Nam cũ) =====
('20200','48','Bến xe Trung tâm Đà Nẵng','Da Nang Central Bus Station','Tôn Đức Thắng, Phường Hòa Khánh, Đà Nẵng',16.096000,108.155000),
('20341','48','Bến xe Tam Kỳ','Tam Ky Bus Station','Phường Tam Kỳ, Đà Nẵng',15.574000,108.475000),
('20410','48','Bến xe Hội An','Hoi An Bus Station','Phường Hội An, Đà Nẵng',15.883000,108.334000),
-- ===== 51 Quảng Ngãi (gồm Kon Tum cũ) =====
('21025','51','Bến xe Quảng Ngãi','Quang Ngai Bus Station','Phường Cẩm Thành, Quảng Ngãi',15.118000,108.796000),
('23293','51','Bến xe Kon Tum','Kon Tum Bus Station','Phường Kon Tum, Quảng Ngãi',14.352000,108.000000),
-- ===== 52 Gia Lai (gồm Bình Định cũ) =====
('21583','52','Bến xe Quy Nhơn','Quy Nhon Bus Station','Phường Quy Nhơn, Gia Lai',13.775000,109.215000),
('23575','52','Bến xe Pleiku','Pleiku Bus Station','Phường Pleiku, Gia Lai',13.974000,108.006000),
('21640','52','Bến xe Bồng Sơn','Bong Son Bus Station','Phường Bồng Sơn, Gia Lai',14.429000,109.007000),
-- ===== 56 Khánh Hòa (gồm Ninh Thuận cũ) =====
('22333','56','Bến xe phía Bắc Nha Trang','Nha Trang North Bus Station','Phường Bắc Nha Trang, Khánh Hòa',12.271000,109.185000),
('22420','56','Bến xe Cam Ranh','Cam Ranh Bus Station','Phường Cam Ranh, Khánh Hòa',11.912000,109.145000),
('22759','56','Bến xe Phan Rang','Phan Rang Bus Station','Phường Phan Rang, Khánh Hòa',11.568000,108.988000),
('22528','56','Bến xe Ninh Hòa','Ninh Hoa Bus Station','Phường Ninh Hòa, Khánh Hòa',12.496000,109.128000),
-- ===== 66 Đắk Lắk (gồm Phú Yên cũ) =====
('24121','66','Bến xe Đắk Lắk','Dak Lak Bus Station','71 Nguyễn Chí Thanh, Buôn Ma Thuột, Đắk Lắk',12.675200,108.061000),
('22015','66','Bến xe Tuy Hòa','Tuy Hoa Bus Station','Phường Tuy Hòa, Đắk Lắk',13.083000,109.310000),
-- ===== 68 Lâm Đồng (gồm Đắk Nông, Bình Thuận cũ) =====
('24781','68','Bến xe Liên tỉnh Đà Lạt','Da Lat Intercity Bus Station','1 Tô Hiến Thành, Đà Lạt, Lâm Đồng',11.918500,108.435800),
('24820','68','Bến xe Bảo Lộc','Bao Loc Bus Station','Phường 2 Bảo Lộc, Lâm Đồng',11.537300,107.808000),
('22945','68','Bến xe Phan Thiết','Phan Thiet Bus Station','Phường Phan Thiết, Lâm Đồng',10.935000,108.108000),
('24611','68','Bến xe Gia Nghĩa','Gia Nghia Bus Station','Phường Bắc Gia Nghĩa, Lâm Đồng',12.004000,107.690000),
-- ===== 75 Đồng Nai (gồm Bình Phước cũ) =====
('26041','75','Bến xe Biên Hòa','Bien Hoa Bus Station','Phường Trấn Biên, Đồng Nai',10.956600,106.844900),
('26080','75','Bến xe Long Khánh','Long Khanh Bus Station','Phường Long Khánh, Đồng Nai',10.936700,107.242000),
('26368','75','Bến xe Long Thành','Long Thanh Bus Station','Xã Long Thành, Đồng Nai',10.793600,107.012700),
('25210','75','Bến xe Đồng Xoài','Dong Xoai Bus Station','Phường Đồng Xoài, Đồng Nai',11.534000,106.908000),
('25217','75','Bến xe Phước Long','Phuoc Long Bus Station','Phường Phước Long, Đồng Nai',11.826000,106.987000),
('25333','75','Bến xe An Lộc','An Loc Bus Station','Phường An Lộc, Đồng Nai',11.649000,106.601000),
-- ===== 79 TP. Hồ Chí Minh (gồm Bình Dương, Bà Rịa - Vũng Tàu cũ) =====
('26833','79','Bến xe Miền Đông mới','Mien Dong New Bus Station','501 Hoàng Hữu Nam, Phường Long Bình, TP. Hồ Chí Minh',10.878600,106.816000),
('26956','79','Bến xe Miền Đông (cũ)','Mien Dong Bus Station (Old)','292 Đinh Bộ Lĩnh, Bình Thạnh, TP. Hồ Chí Minh',10.814900,106.711400),
('27460','79','Bến xe Miền Tây','Mien Tay Bus Station','395 Kinh Dương Vương, Phường An Lạc, TP. Hồ Chí Minh',10.741500,106.619900),
('27592','79','Bến xe An Sương','An Suong Bus Station','Quốc lộ 22, Xã Bà Điểm, TP. Hồ Chí Minh',10.843700,106.612200),
('27367','79','Bến xe Chợ Lớn','Cho Lon Bus Station','279 Lê Quang Sung, Phường Bình Tây, TP. Hồ Chí Minh',10.749700,106.643000),
('27553','79','Bến xe Củ Chi','Cu Chi Bus Station','Xã Củ Chi, TP. Hồ Chí Minh',10.974000,106.493000),
('25747','79','Bến xe Bình Dương','Binh Duong Bus Station','Phường Thủ Dầu Một, TP. Hồ Chí Minh',10.980000,106.660000),
('25813','79','Bến xe Bến Cát','Ben Cat Bus Station','Phường Bến Cát, TP. Hồ Chí Minh',11.145700,106.596000),
('26506','79','Bến xe Vũng Tàu','Vung Tau Bus Station','192 Nam Kỳ Khởi Nghĩa, Phường Vũng Tàu, TP. Hồ Chí Minh',10.360000,107.083000),
('26560','79','Bến xe Bà Rịa','Ba Ria Bus Station','Phường Bà Rịa, TP. Hồ Chí Minh',10.497000,107.169000),
-- ===== 80 Tây Ninh (gồm Long An cũ) =====
('25459','80','Bến xe Tây Ninh','Tay Ninh Bus Station','Phường Tân Ninh, Tây Ninh',11.311200,106.100000),
('27712','80','Bến xe Tân An','Tan An Bus Station','Phường Tân An, Tây Ninh',10.532000,106.414000),
('27787','80','Bến xe Kiến Tường','Kien Tuong Bus Station','Phường Kiến Tường, Tây Ninh',10.779000,105.956000),
('27991','80','Bến xe Bến Lức','Ben Luc Bus Station','Xã Bến Lức, Tây Ninh',10.636000,106.490000),
-- ===== 82 Đồng Tháp (gồm Tiền Giang cũ) =====
('29869','82','Bến xe Cao Lãnh','Cao Lanh Bus Station','Phường Cao Lãnh, Đồng Tháp',10.460000,105.641000),
('28261','82','Bến xe Tiền Giang (Mỹ Tho)','Tien Giang Bus Station (My Tho)','Phường Mỹ Tho, Đồng Tháp',10.356000,106.360000),
('29905','82','Bến xe Sa Đéc','Sa Dec Bus Station','Phường Sa Đéc, Đồng Tháp',10.298000,105.759000),
('28306','82','Bến xe Gò Công','Go Cong Bus Station','Phường Gò Công, Đồng Tháp',10.362000,106.671000),
('29955','82','Bến xe Hồng Ngự','Hong Ngu Bus Station','Phường Hồng Ngự, Đồng Tháp',10.817000,105.335000),
('28439','82','Bến xe Cai Lậy','Cai Lay Bus Station','Phường Cai Lậy, Đồng Tháp',10.409000,106.120000),
-- ===== 86 Vĩnh Long (gồm Bến Tre, Trà Vinh cũ) =====
('29551','86','Bến xe Vĩnh Long','Vinh Long Bus Station','Phường Long Châu, Vĩnh Long',10.249000,105.973000),
('28789','86','Bến xe Bến Tre','Ben Tre Bus Station','Phường Bến Tre, Vĩnh Long',10.242000,106.375000),
('29242','86','Bến xe Trà Vinh','Tra Vinh Bus Station','Phường Trà Vinh, Vĩnh Long',9.935000,106.345000),
-- ===== 91 An Giang (gồm Kiên Giang cũ) =====
('30307','91','Bến xe Long Xuyên','Long Xuyen Bus Station','Phường Long Xuyên, An Giang',10.379000,105.431000),
('30316','91','Bến xe Châu Đốc','Chau Doc Bus Station','Phường Châu Đốc, An Giang',10.709000,105.119000),
('30742','91','Bến xe Rạch Giá','Rach Gia Bus Station','Phường Rạch Giá, An Giang',10.010000,105.085000),
('30769','91','Bến xe Hà Tiên','Ha Tien Bus Station','Phường Hà Tiên, An Giang',10.383000,104.487000),
('31078','91','Bến xe Phú Quốc (Dương Đông)','Phu Quoc Bus Station (Duong Dong)','Đặc khu Phú Quốc, An Giang',10.218000,103.966000),
-- ===== 92 Cần Thơ (gồm Sóc Trăng, Hậu Giang cũ) =====
('31186','92','Bến xe Trung tâm Cần Thơ','Can Tho Central Bus Station','Quốc lộ 1A, Phường Cái Răng, Cần Thơ',10.004400,105.784200),
('31507','92','Bến xe Sóc Trăng','Soc Trang Bus Station','Phường Sóc Trăng, Cần Thơ',9.600000,105.974000),
('31321','92','Bến xe Vị Thanh','Vi Thanh Bus Station','Phường Vị Thanh, Cần Thơ',9.783000,105.470000),
-- ===== 96 Cà Mau (gồm Bạc Liêu cũ) =====
('32002','96','Bến xe Cà Mau','Ca Mau Bus Station','Phường An Xuyên, Cà Mau',9.174000,105.151000),
('31825','96','Bến xe Bạc Liêu','Bac Lieu Bus Station','Phường Bạc Liêu, Cà Mau',9.278000,105.717000),
('31942','96','Bến xe Giá Rai','Gia Rai Bus Station','Phường Giá Rai, Cà Mau',9.238000,105.463000)

-- ===== 31 Hải Phòng (gồm Hải Dương cũ) =====
('11383','31','Bến xe Niệm Nghĩa','Niem Nghia Bus Station','Số 1 Nguyễn Văn Linh, Hải Phòng',20.8400,106.6800),
('10525','31','Bến xe Hải Dương','Hai Duong Bus Station','Phường Hải Dương, Hải Phòng',20.9390,106.3300)
) AS v(ward_code, province_code, name_vi, name_en, address, latitude, longitude)
JOIN administrative_areas p
  ON p.code = v.province_code AND p.parent_id IS NULL
JOIN administrative_areas a
  ON a.code = v.ward_code AND a.parent_id = p.id
WHERE NOT EXISTS (
    SELECT 1 FROM locations l
    WHERE l.type = 'bus_station'
      AND l.name ->> 'vi' = v.name_vi
      AND l.administrative_area_id = a.id
);

-- ---------------------------------------------------------------------
-- 2) REST AREAS (trạm dừng nghỉ cao tốc đang khai thác theo công bố của Cục Đường bộ VN)
--    CHƯA có toạ độ/ward chính xác -> NULL + status 'inactive' để không hiển thị sai.
--    Cập nhật sau khi khảo sát (UPDATE locations SET administrative_area_id=..., latitude=..., longitude=..., status='active').
-- ---------------------------------------------------------------------
INSERT INTO locations
    (administrative_area_id, operator_id, name, type, address, latitude, longitude, status, created_at, updated_at)
SELECT
    NULL,
    NULL,
    jsonb_build_object('vi', v.name_vi, 'en', v.name_en),
    'rest_area',
    v.address,
    NULL,
    NULL,
    'inactive',
    NOW(),
    NOW()
FROM (VALUES
('Trạm dừng nghỉ Cầu Giẽ - Ninh Bình','Cau Gie - Ninh Binh Rest Area','Cao tốc Cầu Giẽ - Ninh Bình'),
('Trạm dừng nghỉ Xuân Khiêm (Cao Bồ - Mai Sơn)','Xuan Khiem Rest Area (Cao Bo - Mai Son)','Cao tốc Cao Bồ - Mai Sơn'),
('Trạm dừng nghỉ Xuân Cương (Cao Bồ - Mai Sơn)','Xuan Cuong Rest Area (Cao Bo - Mai Son)','Cao tốc Cao Bồ - Mai Sơn'),
('Trạm dừng nghỉ La Sơn - Hòa Liên','La Son - Hoa Lien Rest Area','Cao tốc La Sơn - Hòa Liên'),
('Trạm dừng nghỉ TP.HCM - Long Thành - Dầu Giây','HCMC - Long Thanh - Dau Giay Rest Area','Cao tốc TP.HCM - Long Thành - Dầu Giây'),
('Trạm dừng nghỉ TP.HCM - Trung Lương (Km28+200)','HCMC - Trung Luong Rest Area (Km28+200)','Cao tốc TP.HCM - Trung Lương'),
('Trạm dừng nghỉ Vĩnh Hảo - Phan Thiết (Km205)','Vinh Hao - Phan Thiet Rest Area (Km205)','Cao tốc Vĩnh Hảo - Phan Thiết'),
('Trạm dừng tạm Km113 Cam Lâm - Vĩnh Hảo','Cam Lam - Vinh Hao Temporary Stop (Km113)','Cao tốc Cam Lâm - Vĩnh Hảo')
) AS v(name_vi, name_en, address)
WHERE NOT EXISTS (
    SELECT 1 FROM locations l
    WHERE l.type = 'rest_area' AND l.name ->> 'vi' = v.name_vi
);

COMMIT;

-- ---------------------------------------------------------------------
-- Kiểm tra sau khi chạy:
--   SELECT type, count(*) FROM locations GROUP BY type;
--   -- bến xe bị thiếu do sai ward_code/province_code (kỳ vọng 0 dòng):
--   --   so sánh số dòng VALUES (~107) với count bus_station.
-- ---------------------------------------------------------------------