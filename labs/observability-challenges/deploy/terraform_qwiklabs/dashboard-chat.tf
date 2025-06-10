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

resource "google_monitoring_dashboard" "chat_dashboard" {
  project = var.gcp_project_id
  dashboard_json = jsonencode({
    "displayName" : "MovieGuru-ChatMetrics-Dashboard",
    "mosaicLayout" : {
      "columns" : 48,
      "tiles" : [
        {
          "xPos" : 0,
          "yPos" : 0,
          "width" : 24,
          "height" : 16,
          "widget" : {
            "xyChart" : {
              "dataSets" : [
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_calls_success_total[$${__interval}])) / sum(rate(movieguru_chat_calls_total[$${__interval}]))) * 100, \"legend\", \"Success Rate\", \"\", \"\")",
                    "unitOverride" : "%",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                }
              ],
              "thresholds" : [{
                "value" : 95,
                "targetAxis" : "Y1"
              }],
              "yAxis" : {
                "label" : "Percentage Successful",
                "scale" : "LINEAR"
              },
              "chartOptions" : {
                "mode" : "COLOR",
                "showLegend" : false,
                "displayHorizontal" : false
              }
            },
            "title" : "Chat Success Rate",
            "id" : ""
          }
        },

        {
          "xPos" : 0,
          "yPos" : 16,
          "width" : 24,
          "height" : 16,
          "widget" : {
            "xyChart" : {
              "dataSets" : [
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace(histogram_quantile(0.1, sum(rate(movieguru_chat_latency_bucket[$${__interval}])) by (le)), \"legend\", \"10 Percentile\", \"\", \"\")",
                    "unitOverride" : "ms",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                },
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace(histogram_quantile(0.5, sum(rate(movieguru_chat_latency_bucket[$${__interval}])) by (le)), \"legend\", \"50 Percentile\", \"\", \"\")",
                    "unitOverride" : "ms",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                },
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace(histogram_quantile(0.9, sum(rate(movieguru_chat_latency_bucket[$${__interval}])) by (le)), \"legend\", \"90 Percentile\", \"\", \"\")",
                    "unitOverride" : "ms",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                },
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace(histogram_quantile(0.95, sum(rate(movieguru_chat_latency_bucket[$${__interval}])) by (le)), \"legend\", \"95 Percentile\", \"\", \"\")",
                    "unitOverride" : "ms",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                },
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace(histogram_quantile(0.99, sum(rate(movieguru_chat_latency_bucket[$${__interval}])) by (le)), \"legend\", \"99 Percentile\", \"\", \"\")",
                    "unitOverride" : "ms",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                }
              ],
              "thresholds" : [],
              "yAxis" : {
                "label" : "",
                "scale" : "LINEAR"
              },
              "chartOptions" : {
                "mode" : "COLOR",
                "showLegend" : false,
                "displayHorizontal" : false
              }
            },
            "title" : "Chat Latency",
            "id" : ""
          }
        },
      ]
    },
    "dashboardFilters" : [],
    "labels" : {}
  })
}