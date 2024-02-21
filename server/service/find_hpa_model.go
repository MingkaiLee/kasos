package service

type HpaModel struct {
	Name        *string `json:"name"`
	TrainScript *string `json:"train_script"`
	TestScript  *string `json:"test_script"`
}
