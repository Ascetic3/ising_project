package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"ising_project/ising"
)

func main() {
	// Основные параметры моделирования.
	params := ising.Params{
		L:      8,
		T1:     5.0,
		T2:     0.1,
		Tcount: 10,
		ASteps: 15000,
		MSteps: 20000,
		Copies: 25,
	}

	Jvalues := []float64{1, -1}
	Hmin, Hmax, Hstep := -5.0, 5.0, 1.0

	// Создаём папку для результатов (по желанию).
	if err := os.MkdirAll("results", 0o755); err != nil {
		log.Fatalf("cannot create results directory: %v", err)
	}

	for _, J := range Jvalues {
		for H := Hmin; H <= Hmax+1e-9; H += Hstep {
			fmt.Printf("Running Ising for J=%.0f, h=%.0f\n", J, H)

			rows, err := ising.RunIsing(J, H, params)
			if err != nil {
				log.Fatalf("RunIsing failed for J=%.0f, h=%.0f: %v", J, H, err)
			}

			filename := fmt.Sprintf("results_J_%d_h_%d.csv", int(J), int(H))
			fullPath := filepath.Join("results", filename)

			f, err := os.Create(fullPath)
			if err != nil {
				log.Fatalf("cannot create file %s: %v", fullPath, err)
			}

			w := csv.NewWriter(f)
			w.Comma = ';'

			for _, r := range rows {
				record := []string{
					fmt.Sprintf("%f", r.T),
					fmt.Sprintf("%f", r.U),
					fmt.Sprintf("%f", r.M),
					fmt.Sprintf("%f", r.C),
					fmt.Sprintf("%f", r.X),
				}
				if err := w.Write(record); err != nil {
					_ = f.Close()
					log.Fatalf("cannot write to file %s: %v", fullPath, err)
				}
			}
			w.Flush()
			if err := w.Error(); err != nil {
				_ = f.Close()
				log.Fatalf("flush error for file %s: %v", fullPath, err)
			}

			if err := f.Close(); err != nil {
				log.Fatalf("cannot close file %s: %v", fullPath, err)
			}

			// Запуск gnuplot для построения графиков по этому файлу.
			cmd := exec.Command("gnuplot", "-e", fmt.Sprintf("file='%s'", fullPath), "plots/ising.plt")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				log.Fatalf("gnuplot failed for %s: %v", fullPath, err)
			}
		}
	}
}
