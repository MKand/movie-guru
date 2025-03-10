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

resource "google_monitoring_dashboard" "chat_quality_dashboard" {
  project    = var.project_id
  dashboard_json = jsonencode({
    "displayName" : "MovieGuru-ChatQualityMetrics-Dashboard",
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
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_outcome_counter_total{Outcome=~\"Engaged\"}[$${__interval}])) / sum(rate(movieguru_chat_outcome_counter_total[$${__interval}]))) * 100, \"legend\", \"Engaged\", \"\", \"\")",
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
            "title" : "Implicit User Engagement",
            "id" : ""
          }
        },
        {
           "xPos" : 24,
          "yPos" : 0,
          "width" : 24,
          "height" : 16,
          "widget" : {
            "xyChart" : {
              "dataSets" : [
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_acceptance_counter_total[$${__interval}])) / sum(rate(movieguru_chat_calls_total[$${__interval}]))) * 100, \"legend\", \"Accepted\", \"\", \"\")",
                    "unitOverride" : "%",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                },
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
            "title" : "Explicit User Acceptance from Genkit",
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
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_sentiment_counter_total{Sentiment=\"Positive\"}[$${__interval}])) / sum(rate(movieguru_chat_sentiment_counter_total[$${__interval}]))) * 100, \"legend\", \"Positive\", \"\", \"\")",
                    "unitOverride" : "%",
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
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_sentiment_counter_total{Sentiment=\"Negative\"}[$${__interval}])) / sum(rate(movieguru_chat_sentiment_counter_total[$${__interval}]))) * 100, \"legend\", \"Negative\", \"\", \"\")",
                    "unitOverride" : "%",
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
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_sentiment_counter_total{Sentiment=\"Neutral\"}[$${__interval}])) / sum(rate(movieguru_chat_sentiment_counter_total[$${__interval}]))) * 100, \"legend\", \"Neutral\", \"\", \"\")",
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
            "title" : "Implicit User Sentiment",
            "id" : ""
          }
        },
        {
           "xPos": 24,
          "yPos" : 16,
          "width" : 24,
          "height" : 16,
          "widget" : {
            "xyChart" : {
              "dataSets" : [
                {
                  "timeSeriesQuery" : {
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_feedback_counter_total{Feedback=~\"positive\"}[$${__interval}])) / sum(rate(movieguru_chat_calls_total[$${__interval}]))) * 100, \"legend\", \"Positive\", \"\", \"\")",
                    "unitOverride" : "%",
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
                    "prometheusQuery" : "label_replace((sum(rate(movieguru_chat_feedback_counter_total{Feedback=~\"negative\"}[$${__interval}])) / sum(rate(movieguru_chat_calls_total[$${__interval}]))) * 100, \"legend\", \"Negative\", \"\", \"\")",
                    "unitOverride" : "%",
                    "outputFullDuration" : false
                  },
                  "plotType" : "LINE",
                  "legendTemplate" : "",
                  "targetAxis" : "Y1",
                  "dimensions" : [],
                  "measures" : [],
                  "breakdowns" : []
                },
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
            "title" : "Explicit User Feedback from Genkit",
            "id" : ""
          }
        },
      ]
    },
    "dashboardFilters" : [],
    "labels" : {}
    })
}