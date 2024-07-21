package interfaces

import "net/http"

type HealthChecker interface {
	GetCheckHandler() (http.Handler, error)
}
