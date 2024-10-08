package metrics

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type ChatMeters struct {
	CCounter                      metric.Int64Counter
	CSuccessCounter               metric.Int64Counter
	CSentimentPositiveCounter     metric.Int64Counter
	CSentimentNegativeCounter     metric.Int64Counter
	CSentimentNeutralCounter      metric.Int64Counter
	CSentimentUnclassifiedCounter metric.Int64Counter
	COutcomeAcknowledgedCounter   metric.Int64Counter
	COutcomeEngagedCounter        metric.Int64Counter
	COutcomeIrrelevantCounter     metric.Int64Counter
	COutcomeRejectedCounter       metric.Int64Counter
	COutcomeUnclassifiedCounter   metric.Int64Counter
	COutcomeCounter               metric.Int64Counter
	CSafetyIssueCounter           metric.Int64Counter
	CLatencyHistogram             metric.Int64Histogram
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
	cSentimentCounterPositive, err := meter.Int64Counter("movieguru_chat_sentimentpositive_counter", metric.WithDescription("Positive Sentiment counter"))
	if err != nil {
		log.Printf("Error creating bucketed sentiment counter: %v", err)
	}
	cSentimentCounterNegative, err := meter.Int64Counter("movieguru_chat_sentimentnegative_counter", metric.WithDescription("Negative Sentiment counter"))
	if err != nil {
		log.Printf("Error creating bucketed sentiment counter: %v", err)
	}
	cSentimentCounterNeutral, err := meter.Int64Counter("movieguru_chat_sentimentneutral_counter", metric.WithDescription("Neutral Sentiment counter"))
	if err != nil {
		log.Printf("Error creating bucketed sentiment counter: %v", err)
	}
	cSentimentCounterUnclassified, err := meter.Int64Counter("movieguru_chat_sentimentunclassified_counter", metric.WithDescription("Unclassified Sentiment counter"))
	if err != nil {
		log.Printf("Error creating bucketed sentiment counter: %v", err)
	}

	cOutcomeAcknowledgedCounter, err := meter.Int64Counter("movieguru_chat_outcomeAck_counter", metric.WithDescription("Acknowledged Outcome counter"))
	if err != nil {
		log.Printf("Error creating bucketed outcome counter: %v", err)
	}
	cOutcomeEngagedCounter, err := meter.Int64Counter("movieguru_chat_outcomeEngaged_counter", metric.WithDescription("Engaged Outcome counter"))
	if err != nil {
		log.Printf("Error creating bucketed outcome counter: %v", err)
	}
	cOutcomeIrrelevantCounter, err := meter.Int64Counter("movieguru_chat_outcomeIrrelevant_counter", metric.WithDescription("Irrelevant Outcome counter"))
	if err != nil {
		log.Printf("Error creating bucketed outcome counter: %v", err)
	}
	cOutcomeRejectedCounter, err := meter.Int64Counter("movieguru_chat_outcomeRejected_counter", metric.WithDescription("Rejected Outcome counter"))
	if err != nil {
		log.Printf("Error creating bucketed outcome counter: %v", err)
	}
	cOutcomeUnclassfiedCounter, err := meter.Int64Counter("movieguru_chat_outcomeUnclassified_counter", metric.WithDescription("Unclassified Outcome counter"))
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
		CCounter:                      cCounter,
		CLatencyHistogram:             cLatencyHistogram,
		CSuccessCounter:               cSuccessCounter,
		CSafetyIssueCounter:           cSafetyIssueCounter,
		CSentimentPositiveCounter:     cSentimentCounterPositive,
		CSentimentNegativeCounter:     cSentimentCounterNegative,
		CSentimentNeutralCounter:      cSentimentCounterNeutral,
		CSentimentUnclassifiedCounter: cSentimentCounterUnclassified,
		COutcomeAcknowledgedCounter:   cOutcomeAcknowledgedCounter,
		COutcomeEngagedCounter:        cOutcomeEngagedCounter,
		COutcomeIrrelevantCounter:     cOutcomeIrrelevantCounter,
		COutcomeRejectedCounter:       cOutcomeRejectedCounter,
		COutcomeUnclassifiedCounter:   cOutcomeUnclassfiedCounter,
	}
}
