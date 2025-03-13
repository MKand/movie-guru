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
  "displayName": "MovieGuru-ChatQualityMetrics-Dashboard",
  "dashboardFilters": [
    {
      "filterType": "",
      "labelKey": "undefined",
      "stringValue": "",
      "valueType": "STRING"
    }
  ],
  "mosaicLayout": {
    "columns": 48,
    "tiles": [
      {
        "height": 16,
        "width": 24,
        "widget": {
          "title": "Implicit User Sentiment: Positive",
          "scorecard": {
            "gaugeView": {
              "lowerBound": 0,
              "upperBound": 100
            },
            "thresholds": [
              {
                "color": "YELLOW",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 60
              },
              {
                "color": "RED",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 30
              }
            ],
            "timeSeriesQuery": {
              "outputFullDuration": true,
              "prometheusQuery": "label_replace((sum(rate(movieguru_chat_sentiment_counter_total{Sentiment=\"Positive\"}[${__interval}])) / sum(rate(movieguru_chat_sentiment_counter_total[${__interval}]))) * 100, \"legend\", \"Positive\", \"\", \"\")",
              "unitOverride": "%"
            }
          }
        }
      },
      {
        "xPos": 24,
        "height": 16,
        "width": 24,
        "widget": {
          "title": "Implicit User Sentiment: Negative",
          "scorecard": {
            "gaugeView": {
              "lowerBound": 0,
              "upperBound": 100
            },
            "thresholds": [
              {
                "color": "YELLOW",
                "direction": "ABOVE",
                "targetAxis": "Y1",
                "value": 20
              },
              {
                "color": "RED",
                "direction": "ABOVE",
                "targetAxis": "Y1",
                "value": 30
              }
            ],
            "timeSeriesQuery": {
              "outputFullDuration": true,
              "prometheusQuery": "label_replace((sum(rate(movieguru_chat_sentiment_counter_total{Sentiment=\"Negative\"}[${__interval}])) / sum(rate(movieguru_chat_sentiment_counter_total[${__interval}]))) * 100, \"legend\", \"Negative\", \"\", \"\")",
              "unitOverride": "%"
            }
          }
        }
      },
      {
        "yPos": 16,
        "height": 16,
        "width": 24,
        "widget": {
          "title": "Explicit Positive User Feedback from Genkit",
          "scorecard": {
            "gaugeView": {
              "lowerBound": 0,
              "upperBound": 100
            },
            "thresholds": [
              {
                "color": "YELLOW",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 30
              },
              {
                "color": "RED",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 15
              }
            ],
            "timeSeriesQuery": {
              "outputFullDuration": true,
              "prometheusQuery": "label_replace((sum(rate(movieguru_chat_feedback_counter_total{Feedback=~\"positive\"}[${__interval}])) / sum(rate(movieguru_chat_calls_total[${__interval}]))) * 100, \"legend\", \"Positive\", \"\", \"\")",
              "unitOverride": "%"
            }
          }
        }
      },
      {
        "yPos": 16,
        "xPos": 24,
        "height": 16,
        "width": 24,
        "widget": {
          "title": "Explicit User Negative Feedback from Genkit",
          "scorecard": {
            "gaugeView": {
              "lowerBound": 0,
              "upperBound": 100
            },
            "thresholds": [
              {
                "color": "YELLOW",
                "direction": "ABOVE",
                "targetAxis": "Y1",
                "value": 20
              },
              {
                "color": "RED",
                "direction": "ABOVE",
                "targetAxis": "Y1",
                "value": 40
              }
            ],
            "timeSeriesQuery": {
              "outputFullDuration": true,
              "prometheusQuery": "label_replace((sum(rate(movieguru_chat_feedback_counter_total{Feedback=~\"negative\"}[${__interval}])) / sum(rate(movieguru_chat_calls_total[${__interval}]))) * 100, \"legend\", \"Negative\", \"\", \"\")",
              "unitOverride": "%"
            }
          }
        }
      },
      {
        "yPos": 32,
        "height": 16,
        "width": 24,
        "widget": {
          "title": "Implicit User Engagement",
          "scorecard": {
            "gaugeView": {
              "lowerBound": 0,
              "upperBound": 100
            },
            "thresholds": [
              {
                "color": "YELLOW",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 60
              },
              {
                "color": "RED",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 30
              }
            ],
            "timeSeriesQuery": {
              "outputFullDuration": true,
              "prometheusQuery": "label_replace((sum(rate(movieguru_chat_outcome_counter_total{Outcome=~\"Engaged\"}[${__interval}])) / sum(rate(movieguru_chat_outcome_counter_total[${__interval}]))) * 100, \"legend\", \"Engaged\", \"\", \"\")",
              "unitOverride": "%"
            }
          }
        }
      },
      {
        "yPos": 32,
        "xPos": 24,
        "height": 16,
        "width": 24,
        "widget": {
          "title": "Explicit User Acceptance from Genkit",
          "scorecard": {
            "gaugeView": {
              "lowerBound": 0,
              "upperBound": 100
            },
            "thresholds": [
              {
                "color": "RED",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 10
              },
              {
                "color": "YELLOW",
                "direction": "BELOW",
                "targetAxis": "Y1",
                "value": 30
              }
            ],
            "timeSeriesQuery": {
              "outputFullDuration": true,
              "prometheusQuery": "label_replace((sum(rate(movieguru_chat_acceptance_counter_total[${__interval}])) / sum(rate(movieguru_chat_calls_total[${__interval}]))) * 100, \"legend\", \"Accepted\", \"\", \"\")",
              "unitOverride": "%"
            }
          }
        }
      }
    ]
  }
})
}