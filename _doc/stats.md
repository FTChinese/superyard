## Daily Sigup

    GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD

If both `start` and `end` are missing from query parameters, the time range defaults to the past 7 days.

If `start` is missing, it defaults to 7 days earlier before `end`.
If `end` is missing, it defaults to 7 days later after `start`.
UTC+08:00 is used rather than UTC time.

* `200 OK` with body:
```json
[
    {
		"count": 123,
		"date": ""
	}
]
```