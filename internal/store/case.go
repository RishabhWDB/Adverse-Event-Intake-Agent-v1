package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Case struct {
	CaseID           string
	PatientAge       int
	PatientSex       string
	SuspectDrug      string
	Dose             string
	EventDescription string
	OnsetDate        string
	Outcome          string
	Reporter         string
	Seriousness      string
	CriteriaMet      string
	RoutingLane      string
	Status           string
	RawNarrative     string
	CreatedAt        time.Time
}

func InsertCase(pool *pgxpool.Pool, c Case) error {
	query := `
		INSERT INTO cases (
			case_id, patient_age, patient_sex, suspect_drug, dose,
			event_description, onset_date, outcome, reporter,
			seriousness, criteria_met, routing_lane, status, raw_narrative
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`
	_, err := pool.Exec(context.Background(), query,
		c.CaseID, c.PatientAge, c.PatientSex, c.SuspectDrug, c.Dose,
		c.EventDescription, c.OnsetDate, c.Outcome, c.Reporter,
		c.Seriousness, c.CriteriaMet, c.RoutingLane, c.Status, c.RawNarrative,
	)
	return err
}

func GetAllCases(pool *pgxpool.Pool) ([]Case, error) {
	query := `
		SELECT case_id, patient_age, patient_sex, suspect_drug, dose,
		       event_description, onset_date, outcome, reporter,
		       seriousness, criteria_met, routing_lane, status, created_at
		FROM cases
		ORDER BY created_at DESC
	`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []Case
	for rows.Next() {
		var c Case
		err := rows.Scan(
			&c.CaseID, &c.PatientAge, &c.PatientSex, &c.SuspectDrug, &c.Dose,
			&c.EventDescription, &c.OnsetDate, &c.Outcome, &c.Reporter,
			&c.Seriousness, &c.CriteriaMet, &c.RoutingLane, &c.Status, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}

func UpdateCaseStatus(pool *pgxpool.Pool, caseID string, newStatus string) error {
	query := `UPDATE cases SET status = $1, updated_at = now() WHERE case_id = $2`
	_, err := pool.Exec(context.Background(), query, newStatus, caseID)
	return err
}

func GetCaseByID(pool *pgxpool.Pool, caseID string) (*Case, error) {
	query := `
		SELECT case_id, patient_age, patient_sex, suspect_drug, dose,
		       event_description, onset_date, outcome, reporter,
		       seriousness, criteria_met, routing_lane, status, raw_narrative, created_at
		FROM cases
		WHERE case_id = $1
	`
	var c Case
	err := pool.QueryRow(context.Background(), query, caseID).Scan(
		&c.CaseID, &c.PatientAge, &c.PatientSex, &c.SuspectDrug, &c.Dose,
		&c.EventDescription, &c.OnsetDate, &c.Outcome, &c.Reporter,
		&c.Seriousness, &c.CriteriaMet, &c.RoutingLane, &c.Status, &c.RawNarrative, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
