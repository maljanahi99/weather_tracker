CREATE TABLE weather_reports
(
    latitude    FLOAT NOT NULL,
    longitude   FLOAT NOT NULL,
    temperature FLOAT NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO weather_reports (latitude, longitude, temperature)
VALUES
    (52.52, 13.41, 34.5);