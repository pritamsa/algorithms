package gohealthcheck

//go:generate counterfeiter . Healthcheckable
type Healthcheckable interface {
	Name() string
	IsHealthy() (bool, error)
}
