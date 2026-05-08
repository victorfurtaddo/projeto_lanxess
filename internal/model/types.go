package model

import "time"

type Reactor struct {
	ID             int        `json:"id"`
	Name           string     `json:"nome"`
	AngleDeg       float64    `json:"angulo"`
	RadiusM        float64    `json:"raio"`
	CapacityKg     float64    `json:"capacidade_kg"`
	LoadKg         float64    `json:"carga_kg"`
	CycleCount     int        `json:"ciclos"`
	Status         string     `json:"status"`
	LidStatus      string     `json:"status_tampa"`
	LastLoadedAt   *time.Time `json:"ultimo_carregamento,omitempty"`
	LastLoadMassKg float64    `json:"ultima_carga_kg"`
}

type CraneSample struct {
	Timestamp          time.Time         `json:"timestamp"`
	AngleDeg           float64           `json:"angulo"`
	RadiusM            float64           `json:"raio"`
	HeightM            float64           `json:"altura"`
	Speed              float64           `json:"velocidade"`
	WeightKg           float64           `json:"peso_kg"`
	VolumeM3           float64           `json:"volume_m3"`
	CraneCapacityKg    float64           `json:"capacidade_grua_kg"`
	OperationalLimitKg float64           `json:"limite_operacional_kg"`
	SuspendedLoad      bool              `json:"carga_suspensa"`
	State              string            `json:"estado"`
	NearestReactorID   int               `json:"reator_mais_proximo_id"`
	People             []PersonDetection `json:"pessoas"`
}

type Candidate struct {
	ReactorID          int     `json:"reator_id"`
	AngularDistanceDeg float64 `json:"distancia_angular"`
	RadialDistanceM    float64 `json:"distancia_radial"`
	Score              float64 `json:"score"`
	Probability        float64 `json:"probabilidade"`
}

type InferredEvent struct {
	Timestamp  time.Time   `json:"timestamp"`
	ReactorID  int         `json:"reator_id"`
	AmountKg   float64     `json:"quantidade_kg"`
	Confidence float64     `json:"confianca"`
	AngleDeg   float64     `json:"angulo"`
	RadiusM    float64     `json:"raio"`
	Reason     string      `json:"motivo"`
	Candidates []Candidate `json:"candidatos"`
}

type SimulationResult struct {
	Reactors []Reactor     `json:"reatores"`
	Samples  []CraneSample `json:"amostras"`
	Truth    []TruthEvent  `json:"eventos_reais"`
}

type TruthEvent struct {
	Timestamp       time.Time `json:"timestamp"`
	ReactorID       int       `json:"reator_id"`
	VolumeM3        float64   `json:"volume_m3"`
	EstimatedMassKg float64   `json:"massa_estimada_kg"`
	CraneLimitKg    float64   `json:"limite_grua_kg"`
	NetworkLimitKg  float64   `json:"limite_rede_kg"`
	FinalMassKg     float64   `json:"massa_final_kg"`
	AmountKg        float64   `json:"quantidade_kg"`
}

type DemoResult struct {
	Reactors        []Reactor       `json:"reatores"`
	Crane           CraneState      `json:"grua"`
	Samples         []CraneSample   `json:"amostras"`
	Truth           []TruthEvent    `json:"eventos_reais"`
	Inferred        []InferredEvent `json:"eventos_inferidos"`
	Alerts          []Alert         `json:"alertas"`
	TotalInferredKg float64         `json:"total_inferido_kg"`
	Accuracy        float64         `json:"acuracia"`
}

type CraneState struct {
	AngleDeg           float64 `json:"angulo"`
	RadiusM            float64 `json:"raio"`
	HeightM            float64 `json:"altura"`
	CurrentLoadKg      float64 `json:"carga_atual_kg"`
	MaxAllowedLoadKg   float64 `json:"carga_maxima_permitida_kg"`
	CraneCapacityKg    float64 `json:"capacidade_grua_kg"`
	OperationalLimitKg float64 `json:"limite_operacional_kg"`
	State              string  `json:"estado"`
	TargetReactorID    int     `json:"reator_alvo_provavel"`
}

type PersonDetection struct {
	ID         string    `json:"id"`
	X          float64   `json:"x"`
	Y          float64   `json:"y"`
	SpeedMps   float64   `json:"velocidade_mps"`
	State      string    `json:"estado"`
	Confidence float64   `json:"confianca"`
	InRiskZone bool      `json:"em_zona_risco"`
	Timestamp  time.Time `json:"timestamp"`
}

type Alert struct {
	ID        string         `json:"id"`
	Level     string         `json:"nivel"`
	Message   string         `json:"mensagem"`
	Timestamp time.Time      `json:"timestamp"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type MassEstimate struct {
	VolumeM3           float64 `json:"volume_m3"`
	DensityKgM3        float64 `json:"densidade_kg_m3"`
	EstimatedMassKg    float64 `json:"massa_estimada_kg"`
	NetworkLimitKg     float64 `json:"limite_rede_kg"`
	CraneLimitKg       float64 `json:"limite_grua_kg"`
	OperationalLimitKg float64 `json:"limite_operacional_kg"`
	FinalMassKg        float64 `json:"massa_final_kg"`
	LimitedBy          string  `json:"limitado_por"`
}

type DigitalTwinState struct {
	Crane    CraneState        `json:"crane"`
	Reactors []Reactor         `json:"reactors"`
	People   []PersonDetection `json:"people"`
	Events   []InferredEvent   `json:"events"`
	Alerts   []Alert           `json:"alerts"`
}
