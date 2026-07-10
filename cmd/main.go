package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/api"
	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/classifier"
	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/extractor"
	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/intake"
	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/router"
	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func processNewEmails(pool *pgxpool.Pool) {
	emails, err := intake.FetchUnseenEmails()
	if err != nil {
		log.Printf("Email fetch failed: %v", err)
		return
	}

	if len(emails) == 0 {
		return
	}

	log.Printf("Found %d new email(s)", len(emails))

	for _, raw := range emails {
		narrative, err := intake.ExtractPlainText(raw)
		if err != nil {
			log.Printf("Failed to parse email: %v", err)
			continue
		}

		caseData, err := extractor.ExtractCase(narrative)
		if err != nil {
			log.Printf("Extraction failed: %v", err)
			continue
		}

		classResult := classifier.Classify(caseData.EventDescription, caseData.Outcome, narrative)
		routingLane := router.AssignLane(classResult.Seriousness)

		caseID := fmt.Sprintf("AE-%d", time.Now().UnixNano())
		newCase := store.Case{
			CaseID:           caseID,
			PatientAge:       caseData.PatientAge,
			PatientSex:       caseData.PatientSex,
			SuspectDrug:      caseData.SuspectDrug,
			Dose:             caseData.Dose,
			EventDescription: caseData.EventDescription,
			OnsetDate:        caseData.OnsetDate,
			Outcome:          caseData.Outcome,
			Reporter:         caseData.Reporter,
			Seriousness:      classResult.Seriousness,
			CriteriaMet:      classResult.CriteriaMet,
			RoutingLane:      routingLane,
			Status:           "pending",
			RawNarrative:     narrative,
		}

		if err := store.InsertCase(pool, newCase); err != nil {
			log.Printf("Insert failed: %v", err)
			continue
		}

		log.Printf("Processed %s: %s (%s)", caseID, newCase.SuspectDrug, newCase.Seriousness)
	}
}

func startEmailPoller(pool *pgxpool.Pool) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	// Run once immediately on startup
	processNewEmails(pool)

	for range ticker.C {
		processNewEmails(pool)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	pool := store.Connect()
	defer pool.Close()

	// Start email poller in background
	go startEmailPoller(pool)

	// Start API server on main thread
	router := api.SetupRouter(pool)
	log.Println("Server starting on :8080")
	log.Println("Polling for new emails every 20 seconds...")
	router.Run(":8080")
}
