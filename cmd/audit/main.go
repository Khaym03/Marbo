package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Khaym03/Marbo/internal/audit"
	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/embedder"
	"github.com/Khaym03/Marbo/internal/expansion"
	"github.com/Khaym03/Marbo/internal/metrics"
	"github.com/Khaym03/Marbo/internal/planner"
	"github.com/Khaym03/Marbo/internal/validator"
	ort "github.com/yalue/onnxruntime_go"
)

const (
	dataFilePath      = "data.json"
	modelFilePath     = "modelo_e5_onnx/model.onnx"
	tokenizerFilePath = "modelo_e5_onnx/tokenizer.json"
	sharedLib         = "third-party/onnxruntime.dll"
)

func main() {
	ort.SetSharedLibraryPath(sharedLib)
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	defer ort.DestroyEnvironment()

	emb, err := embedder.New(modelFilePath, tokenizerFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer emb.Close()

	data, err := domain.Load(dataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := validator.Validate(data); err != nil {
		log.Fatal(err)
	}

	auditor := audit.TransitionAuditor{Embedder: emb}
	auditReport, err := auditor.Audit(data)
	if err != nil {
		log.Fatal(err)
	}

	printAuditReport(auditReport)

	healthReport, err := metrics.BuildHealthReport(data, auditReport)
	if err != nil {
		log.Fatal(err)
	}

	printHealthReport(healthReport)

	coverageReport, err := planner.AnalyzeCoverage(data)
	if err != nil {
		log.Fatal(err)
	}

	printCoverageReport(coverageReport)

	expansionPack := expansion.GenerateExpansionPack(data, coverageReport)
	printExpansionPack(expansionPack)
}

func printAuditReport(report *audit.AuditReport) {
	fmt.Printf("=================================================\n")
	fmt.Printf("FLOW AUDIT\n")
	fmt.Printf("==========\n\n")

	for _, f := range report.Flows {
		fmt.Printf("Flow:\n%s\n\n", f.FlowID)

		fmt.Printf("Warnings:\n\n")
		if len(f.Warnings) == 0 {
			fmt.Printf("None\n\n")
		} else {
			for _, w := range f.Warnings {
				fmt.Printf("* %s\n", w)
			}
			fmt.Println()
		}

		fmt.Printf("Similar Phrases:\n\n")
		if len(f.SimilarTransitions) == 0 {
			fmt.Printf("None\n\n")
		} else {
			for _, s := range f.SimilarTransitions {
				fmt.Printf("\"%s\"\n\"%s\"\n\nScore: %.4f\n\n", s.PhraseA, s.PhraseB, s.Score)
			}
		}

		if len(f.DuplicatePhrases) > 0 {
			fmt.Printf("Duplicate Phrases:\n\n")
			for _, d := range f.DuplicatePhrases {
				fmt.Printf("* \"%s\" (in transitions %s and %s)\n", d.Phrase, d.TransitionA, d.TransitionB)
			}
			fmt.Println()
		}

		fmt.Printf("---\n\n")
	}

	fmt.Printf("=================================================\n")
}

func printHealthReport(report *metrics.KBHealthReport) {
	fmt.Printf("\n=================================================\n")
	fmt.Printf("KB HEALTH REPORT\n")
	fmt.Printf("================\n\n")

	fmt.Printf("Overall Score:\n%.1f\n\n", report.OverallScore)

	fmt.Printf("---\n\nINTENTS\n\n")
	for _, i := range report.IntentMetrics {
		fmt.Printf("%s\n\nTraining Phrases:\n%d\n\nCoverage:\n%.0f\n\nRisk:\n%.0f\n\n---\n\n", i.IntentID, i.TrainingPhraseCount, i.CoverageScore, i.RiskScore)
	}

	fmt.Printf("FLOWS\n\n")
	for _, f := range report.FlowMetrics {
		fmt.Printf("%s\n\nNodes:\n%d\n\nTransitions:\n%d\n\nMax Depth:\n%d\n\nComplexity:\n%.0f\n\nRisk:\n%.0f\n\n---\n\n", f.FlowID, f.NodeCount, f.TransitionCount, f.MaxDepth, f.ComplexityScore, f.RiskScore)
	}

	fmt.Printf("Highest Risk Intent:\n%s\n\nHighest Risk Flow:\n%s\n\n", report.HighestRiskIntent, report.HighestRiskFlow)
	fmt.Printf("=================================================\n")
}

func printCoverageReport(report *planner.ExpansionReport) {
	fmt.Printf("\n=================================================\n")
	fmt.Printf("TRAINING COVERAGE REPORT\n")
	fmt.Printf("========================\n\n")

	for _, i := range report.Intents {
		fmt.Printf("Intent:\n%s\n\nCurrent Phrases:\n%d\n\nRecommended:\n%d\n\nPriority:\n%s\n\nMissing Coverage:\n", i.IntentID, i.CurrentPhraseCount, i.RecommendedPhraseCount, strings.ToUpper(string(i.Priority)))
		for _, m := range i.MissingCoverage {
			fmt.Printf("* %s\n", m)
		}
		fmt.Printf("\n---\n\n")
	}

	fmt.Printf("=================================================\n")
}

func printExpansionPack(pack *expansion.ExpansionPack) {
	fmt.Printf("\n=================================================\n")
	fmt.Printf("KNOWLEDGE EXPANSION PACK\n")
	fmt.Printf("========================\n\n")

	for _, i := range pack.Intents {
		fmt.Printf("Intent:\n%s\n\nCurrent:\n%d\n\nTarget:\n%d\n\nGenerated:\n", i.IntentID, i.ExistingPhraseCount, i.TargetPhraseCount)
		for _, p := range i.GeneratedPhrases {
			fmt.Printf("* %s\n", p)
		}
		fmt.Printf("\n---\n\n")
	}

	fmt.Printf("=================================================\n")
}
