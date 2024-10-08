package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type ChatMeters struct {
	CCounter            metric.Int64Counter
	CSuccessCounter     metric.Int64Counter
	CSentimentCounter   metric.Int64Counter
	COutcomeCounter     metric.Int64Counter
	CSafetyIssueCounter metric.Int64Counter
	CLatencyHistogram   metric.Int64Histogram
}

func NewChatMeters() *ChatMeters {
	meter := otel.Meter("chat-handler")

	cCounter, err := meter.Int64Counter("movieguru_chat_calls_total", metric.WithDescription("Total number of chat calls"))
	if err != nil {
		log.Printf("Error creating chat calls counter: %v", err)
	}
	cSuccessCounter, err := meter.Int64Counter("movieguru_chat_calls_success_total", metric.WithDescription("Total number of chat calls that are successful"))
	if err != nil {
		log.Printf("Error creating chat calls success counter: %v", err)
	}
	cSentimentCounter, err := meter.Int64Counter("movieguru_chat_sentiment_counter", metric.WithDescription("Bucketed Sentiment counter"))
	if err != nil {
		log.Printf("Error creating bucketed sentiment counter: %v", err)
	}

	cOutcomeCounter, err := meter.Int64Counter("movieguru_chat_outcome_counter", metric.WithDescription("Bucketed Outcome counter"))
	if err != nil {
		log.Printf("Error creating bucketed outcome counter: %v", err)
	}

	cSafetyIssueCounter, err := meter.Int64Counter("movieguru_chat_safetyissue_counter", metric.WithDescription("Safety issue counter"))
	if err != nil {
		log.Printf("Error creating safety issue counter: %v", err)
	}
	cLatencyHistogram, err := meter.Int64Histogram("movieguru_chat_latency", metric.WithDescription("Histogram of chat request latency"))
	if err != nil {
		log.Printf("Error creating login latency histogram: %v", err)
	}
	return &ChatMeters{
		CCounter:            cCounter,
		CLatencyHistogram:   cLatencyHistogram,
		CSuccessCounter:     cSuccessCounter,
		CSafetyIssueCounter: cSafetyIssueCounter,
		CSentimentCounter:   cSentimentCounter,
		COutcomeCounter:     cOutcomeCounter,
	}
}
