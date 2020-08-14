package apn

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/push"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type APNEnv struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "model.apn")

func (env APNEnv) ListMessage(p gorest.Pagination) ([]push.MessageTeaser, error) {
	query := `
    SELECT id,
      page_id AS pageId,
      action_type AS actionType,
      title,
      content_available AS contentAvailable,
      created_by AS createdBy,
      created_utc AS createdAt,
      device_count AS deviceCount,
      invalid_count AS invalidCount,
      time_elapsed AS timeElapsed
    FROM analytic.ios_sent_message
    WHERE prod_server = TRUE
      AND device_group != 0
    ORDER BY created_utc DESC
    LIMIT ? OFFSET ?`

	rows, err := env.DB.Query(query, p.Limit, p.Offset())

	if err != nil {
		logger.WithField("trace", "ListMessage").Error(err)
		return nil, err
	}

	defer rows.Close()

	teasers := make([]push.MessageTeaser, 0)

	for rows.Next() {
		var t push.MessageTeaser
		err := rows.Scan(
			&t.ID,
			&t.PageID,
			&t.Action,
			&t.Title,
			&t.ContentAvailable,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.DeviceCount,
			&t.InvalidCount,
			&t.TimeElapsed)

		if err != nil {
			logger.WithField("trace", "LatestStoryList").Error(err)
			continue
		}

		teasers = append(teasers, t)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListMessage").Error(err)
		return teasers, err
	}

	return teasers, nil
}

func (env APNEnv) TimeZoneDist() ([]push.TimeZone, error) {
	query := `
    SELECT device.timezone AS zoneName,
      zone.utc_offset AS utcOffset,
      count(*) AS deviceCount
    FROM analytic.ios_device_token AS device
      LEFT JOIN analytic.ios_zone_offset AS zone
      ON device.timezone = zone.timezone
    GROUP BY device.timezone;`

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("trace", "TimeZoneDist").Error(err)
		return nil, err
	}

	defer rows.Close()

	zones := make([]push.TimeZone, 0)

	for rows.Next() {
		var z push.TimeZone

		err := rows.Scan(
			&z.ZoneName,
			&z.Offset,
			&z.DeviceCount)

		if err != nil {
			logger.WithField("trace", "TimeZoneDist").Error(err)
			continue
		}

		zones = append(zones, z)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "TimeZoneDist").Error(err)
		return zones, err
	}

	return zones, nil
}

func (env APNEnv) DeviceDist() ([]push.Device, error) {
	query := `
    SELECT device_type AS deviceType,
      count(*) deviceCount
    FROM analytic.ios_device_token
    GROUP BY device_type`

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("trace", "DeviceDist").Error(err)
		return nil, err
	}

	defer rows.Close()

	devices := make([]push.Device, 0)

	for rows.Next() {
		var d push.Device

		err := rows.Scan(
			&d.Name,
			&d.Count)

		if err != nil {
			logger.WithField("trace", "LatestStoryList").Error(err)
			continue
		}

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListMessage").Error(err)
		return devices, err
	}

	return devices, nil
}

func (env APNEnv) InvalidDist() ([]push.InvalidDevice, error) {
	query := `
    SELECT invalid.reason,
      device.device_type AS deviceType,
      count(*) AS deviceCount
    FROM analytic.ios_invalid_token AS invalid
      INNER JOIN analytic.ios_device_token AS device
      ON invalid.device_token = UNHEX(device.device_token)
    GROUP BY reason, deviceType`

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("trace", "InvalidDist").Error(err)
		return nil, err
	}

	defer rows.Close()

	devices := make([]push.InvalidDevice, 0)

	for rows.Next() {
		var d push.InvalidDevice

		err := rows.Scan(
			&d.Reason,
			&d.Name,
			&d.Count)

		if err != nil {
			logger.WithField("trace", "InvalidDist").Error(err)
			continue
		}

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "InvalidDist").Error(err)
		return devices, err
	}

	return devices, nil
}

func (env APNEnv) CreateTestDevice(d push.TestDevice) error {
	query := `
	INSERT INTO analytic.ios_test_device
		SET token = UNHEX(?),
		  description = ?,
		  owned_by = ?,
		  created_utc = UTC_TIMESTAMP()`

	_, err := env.DB.Exec(query,
		d.Token,
		d.Description,
		d.OwnedBy)

	if err != nil {
		logger.WithField("trace", "CreateTestDevice").Error(err)
		return err
	}

	return nil
}

func (env APNEnv) ListTestDevice() ([]push.TestDevice, error) {
	query := `
	SELECT id,
      LOWER(HEX(token)) AS token,
      description AS description,
      owned_by AS ownedBy,
      created_utc AS createdAt
    FROM analytic.ios_test_device
    ORDER BY created_utc DESC`

	rows, err := env.DB.Query(query)

	if err != nil {
		logger.WithField("trace", "ListTestDevice").Error(err)
		return nil, err
	}

	defer rows.Close()

	devices := make([]push.TestDevice, 0)

	for rows.Next() {
		var d push.TestDevice

		err := rows.Scan(
			&d.ID,
			&d.Token,
			&d.Description,
			&d.OwnedBy,
			&d.CreatedAt)

		if err != nil {
			logger.WithField("trace", "ListTestDevice").Error(err)
			continue
		}

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListTestDevice").Error(err)
		return devices, err
	}

	return devices, nil
}

func (env APNEnv) RemoveTestDevice(id int64) error {
	query := `
	DELETE FROM analytic.ios_test_device
    WHERE id = ?
    LIMIT 1`

	_, err := env.DB.Exec(query, id)

	if err != nil {
		logger.WithField("trace", "RemoveTestDevice").Error(err)
		return err
	}

	return nil
}
