import unittest

from great_salt_lake import parse_usgs_payload, summarize_observations


def series(site, parameter, values, *, name="station", lat=40.0, lon=-112.0):
    return {
        "sourceInfo": {
            "siteCode": [{"value": site}],
            "siteName": name,
            "geoLocation": {"geogLocation": {"latitude": lat, "longitude": lon}},
        },
        "variable": {
            "variableCode": [{"value": parameter}],
            "unit": {"unitCode": "ft" if parameter in {"62614", "00065"} else "ft3/s"},
        },
        "values": [{"value": values}],
    }


class GreatSaltLakeProfileTests(unittest.TestCase):
    def test_gage_height_is_not_absolute_lake_elevation(self):
        payload = {
            "value": {
                "timeSeries": [
                    series(
                        "10010000",
                        "00065",
                        [{"value": "4.25", "dateTime": "2026-09-24T18:00:00Z"}],
                    )
                ]
            }
        }
        rows = parse_usgs_payload(payload)
        self.assertEqual(rows[0].field, "gage_height_ft")
        summary = summarize_observations(rows, fetched_at=1790276400)
        self.assertIsNone(summary["context"]["lake_surface_elevation_ft"])

    def test_missing_values_remain_missing(self):
        payload = {
            "value": {
                "timeSeries": [
                    series(
                        "10126000",
                        "00060",
                        [{"value": "", "dateTime": "2026-09-24T18:00:00Z"}],
                    )
                ]
            }
        }
        rows = parse_usgs_payload(payload)
        self.assertIsNone(rows[0].value)
        summary = summarize_observations(rows, fetched_at=1790276400)
        self.assertIsNone(summary["context"]["combined_inflow_cfs"])

    def test_elevation_and_inflows_are_kept_separate(self):
        payload = {
            "value": {
                "timeSeries": [
                    series(
                        "10010000",
                        "62614",
                        [
                            {"value": "4191.20", "dateTime": "2026-09-24T17:45:00Z"},
                            {"value": "4191.25", "dateTime": "2026-09-24T18:00:00Z"},
                        ],
                    ),
                    series(
                        "10126000",
                        "00060",
                        [
                            {"value": "100", "dateTime": "2026-09-24T17:45:00Z"},
                            {"value": "110", "dateTime": "2026-09-24T18:00:00Z"},
                        ],
                    ),
                    series(
                        "10141000",
                        "00060",
                        [
                            {"value": "40", "dateTime": "2026-09-24T17:45:00Z"},
                            {"value": "42", "dateTime": "2026-09-24T18:00:00Z"},
                        ],
                    ),
                ]
            }
        }
        rows = parse_usgs_payload(payload)
        summary = summarize_observations(
            rows, fetched_at=1790276400, baseline_elevation_ft=4191.5
        )
        context = summary["context"]
        self.assertEqual(context["lake_surface_elevation_ft"], 4191.25)
        self.assertAlmostEqual(context["lake_surface_elevation_delta_ft"], 0.05)
        self.assertEqual(context["combined_inflow_cfs"], 152.0)
        self.assertEqual(context["combined_inflow_delta_cfs"], 12.0)
        self.assertAlmostEqual(context["elevation_vs_reference_ft"], -0.25)
        self.assertTrue(summary["provenance"]["output_digest"].startswith("sha256:"))


if __name__ == "__main__":
    unittest.main()
