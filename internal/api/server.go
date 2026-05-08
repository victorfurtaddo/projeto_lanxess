package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"projeto-lanxess/internal/inference"
	"projeto-lanxess/internal/model"
	"projeto-lanxess/internal/safety"
	"projeto-lanxess/internal/simulator"
)

type Server struct {
	mux http.Handler
}

func NewServer() *Server {
	s := &Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /api/state", s.handleState)
	mux.HandleFunc("GET /api/crane", s.handleCrane)
	mux.HandleFunc("GET /api/reactors", s.handleReactors)
	mux.HandleFunc("GET /api/reactors/{id}", s.handleReactor)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/alerts", s.handleAlerts)
	mux.HandleFunc("GET /api/demo", s.handleDemo)
	mux.HandleFunc("GET /api/simulate", s.handleSimulate)
	mux.HandleFunc("POST /api/infer", s.handleInfer)
	mux.HandleFunc("POST /api/simulation/start", s.handleSimulationCommand)
	mux.HandleFunc("POST /api/simulation/stop", s.handleSimulationCommand)
	mux.HandleFunc("POST /api/simulation/step", s.handleSimulationCommand)
	mux.HandleFunc("POST /api/simulation/reset", s.handleSimulationCommand)
	mux.Handle("/", http.FileServer(http.Dir("web")))
	s.mux = mux
	return s
}

func (s *Server) Routes() http.Handler {
	return s.mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	demo := buildDemo(configFromRequest(r))
	people := []model.PersonDetection{}
	if len(demo.Samples) > 0 {
		people = demo.Samples[len(demo.Samples)-1].People
	}

	respondJSON(w, model.DigitalTwinState{
		Crane:    demo.Crane,
		Reactors: demo.Reactors,
		People:   people,
		Events:   demo.Inferred,
		Alerts:   demo.Alerts,
	})
}

func (s *Server) handleCrane(w http.ResponseWriter, r *http.Request) {
	demo := buildDemo(configFromRequest(r))
	respondJSON(w, demo.Crane)
}

func (s *Server) handleReactors(w http.ResponseWriter, r *http.Request) {
	demo := buildDemo(configFromRequest(r))
	respondJSON(w, demo.Reactors)
}

func (s *Server) handleReactor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	demo := buildDemo(configFromRequest(r))
	for _, reactor := range demo.Reactors {
		if reactor.ID == id {
			respondJSON(w, reactor)
			return
		}
	}

	http.Error(w, "reator nao encontrado", http.StatusNotFound)
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	demo := buildDemo(configFromRequest(r))
	respondJSON(w, demo.Inferred)
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	demo := buildDemo(configFromRequest(r))
	respondJSON(w, demo.Alerts)
}

func (s *Server) handleDemo(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, buildDemo(configFromRequest(r)))
}

func (s *Server) handleSimulate(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, simulator.Generate(configFromRequest(r)))
}

func (s *Server) handleInfer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload struct {
		Reactors []model.Reactor     `json:"reatores"`
		Samples  []model.CraneSample `json:"amostras"`
		Config   inference.Config    `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON invalido: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(payload.Reactors) == 0 {
		payload.Reactors = simulator.BuildReactors(28, 8, 107)
	}

	events := inference.DetectEvents(payload.Samples, payload.Reactors, payload.Config)
	respondJSON(w, events)
}

func (s *Server) handleSimulationCommand(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status":  "accepted",
		"message": "MVP usa simulacao deterministica por requisicao; estado persistente em memoria fica para a proxima etapa.",
	})
}

func configFromRequest(r *http.Request) simulator.Config {
	cfg := simulator.DefaultConfig()
	cfg.Events = intQuery(r, "events", cfg.Events)
	cfg.Seed = int64(intQuery(r, "seed", int(cfg.Seed)))
	return cfg
}

func buildDemo(cfg simulator.Config) model.DemoResult {
	sim := simulator.Generate(cfg)
	events := inference.DetectEvents(sim.Samples, sim.Reactors, inference.DefaultConfig())
	reactors := simulator.ApplyLoads(sim.Reactors, events)
	alerts := collectAlerts(sim.Samples, sim.Reactors, events)

	total := 0.0
	for _, event := range events {
		total += event.AmountKg
	}

	return model.DemoResult{
		Reactors:        reactors,
		Crane:           craneStateFromSamples(sim.Samples),
		Samples:         sim.Samples,
		Truth:           sim.Truth,
		Inferred:        events,
		Alerts:          alerts,
		TotalInferredKg: total,
		Accuracy:        inference.EvaluateAccuracy(sim.Truth, events),
	}
}

func collectAlerts(samples []model.CraneSample, reactors []model.Reactor, events []model.InferredEvent) []model.Alert {
	alerts := make([]model.Alert, 0)
	for _, sample := range samples {
		alerts = append(alerts, safety.EvaluateSample(sample, reactors)...)
	}
	for _, event := range events {
		alerts = append(alerts, model.Alert{
			ID:        "info_loaded_" + strconv.Itoa(event.ReactorID) + "_" + event.Timestamp.Format("150405"),
			Level:     "INFO",
			Message:   "Reator recebeu novo ciclo de carregamento.",
			Timestamp: event.Timestamp,
			Metadata: map[string]any{
				"reactor_id": event.ReactorID,
				"amount_kg":  event.AmountKg,
				"confidence": event.Confidence,
			},
		})
	}
	return alerts
}

func craneStateFromSamples(samples []model.CraneSample) model.CraneState {
	if len(samples) == 0 {
		return model.CraneState{}
	}
	last := samples[len(samples)-1]
	return model.CraneState{
		AngleDeg:           last.AngleDeg,
		RadiusM:            last.RadiusM,
		HeightM:            last.HeightM,
		CurrentLoadKg:      last.WeightKg,
		MaxAllowedLoadKg:   last.OperationalLimitKg,
		CraneCapacityKg:    last.CraneCapacityKg,
		OperationalLimitKg: last.OperationalLimitKg,
		State:              last.State,
		TargetReactorID:    last.NearestReactorID,
	}
}

func intQuery(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func respondJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(payload)
}
