# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

resource "google_monitoring_dashboard" "quote-rate-dashboard" {
  project        = var.project_id
  dashboard_json = <<EOF
  {
  "displayName": "Model API Rate",
  "mosaicLayout": {
    "columns": 48,
    "tiles": [
      {
        "width": 24,
        "height": 16,
        "widget": {
          "xyChart": {
            "dataSets": [
              {
                "timeSeriesQuery": {
                  "timeSeriesFilter": {
                    "filter": "metric.type=\"aiplatform.googleapis.com/quota/generate_content_requests_per_minute_per_project_per_base_model/usage\" resource.type=\"aiplatform.googleapis.com/Location\" metric.label.\"base_model\"=\"gemini-1.5-flash\" resource.label.\"location\"=\"us-west1\"",
                    "aggregation": {
                      "alignmentPeriod": "60s",
                      "perSeriesAligner": "ALIGN_RATE",
                      "crossSeriesReducer": "REDUCE_SUM",
                      "groupByFields": []
                    }
                  },
                  "unitOverride": "",
                  "outputFullDuration": false
                },
                "plotType": "LINE",
                "legendTemplate": "",
                "minAlignmentPeriod": "60s",
                "targetAxis": "Y1",
                "dimensions": [],
                "measures": [],
                "breakdowns": []
              }
            ],
            "thresholds": [
              {
                "label": "",
                "value": 3.33,
                "color": "COLOR_UNSPECIFIED",
                "direction": "DIRECTION_UNSPECIFIED",
                "targetAxis": "Y1"
              }
            ],
            "yAxis": {
              "label": "",
              "scale": "LINEAR"
            },
            "chartOptions": {
              "mode": "COLOR",
              "showLegend": false,
              "displayHorizontal": false
            }
          },
          "title": "RateLimit",
          "id": ""
        }
      }
    ]
  },
  "dashboardFilters": [],
  "labels": {}
}
  EOF
}