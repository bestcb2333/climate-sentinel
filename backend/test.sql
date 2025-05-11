-- 插入测试数据
DELIMITER $$

CREATE PROCEDURE generate_test_data()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE j INT DEFAULT 1;
    DECLARE current_date DATE;
    DECLARE base_date DATE;

    -- 设定基准日期（5月8日的前一周）
    SET base_date = CURDATE() - INTERVAL 7 DAY;

    -- 循环生成91条数据，13个区域，每个区域生成7天的数据
    WHILE i <= 13 DO
        SET current_date = base_date;
        
        -- 每个区域生成7条数据
        WHILE j <= 7 DO
            INSERT INTO histories (
                created_at, 
                region_id, 
                type, 
                time, 
                max_temperature, 
                min_temperature, 
                avg_temperature, 
                wind_speed, 
                visibility, 
                rain_fall, 
                severity, 
                source
            )
            VALUES (
                NOW(),                               -- 当前时间
                i,                                   -- 区域ID
                ROUND(RAND() * 10, 2),               -- 随机天气类型 (0~10)
                current_date,                        -- 对应日期
                ROUND(RAND() * 35 + 15, 2),          -- 随机最高气温 (15~50)
                ROUND(RAND() * 15, 2),               -- 随机最低气温 (0~15)
                ROUND(RAND() * 25 + 15, 2),          -- 随机平均气温 (15~40)
                ROUND(RAND() * 10, 2),               -- 随机风速 (0~10)
                ROUND(RAND() * 10, 2),               -- 随机能见度 (0~10)
                ROUND(RAND() * 50, 2),               -- 随机降水量 (0~50)
                ROUND(RAND() * 5, 2),                -- 随机严重性 (0~5)
                'Test Source'                        -- 数据来源
            );
            
            -- 增加日期
            SET current_date = current_date + INTERVAL 1 DAY;
            SET j = j + 1;
        END WHILE;

        -- 重置j并继续下一个区域
        SET j = 1;
        SET i = i + 1;
    END WHILE;
    
END$$

DELIMITER ;

-- 调用存储过程生成测试数据
CALL generate_test_data();

