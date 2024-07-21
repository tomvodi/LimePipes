/*
 * Set and Tune API
 *
 * API for managing sets and tunes
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apimodel

import (
	"time"
)

type Health200ResponseDetailsDatabase struct {

	Status HealthStatus `json:"status,omitempty"`

	Timestamp time.Time `json:"timestamp,omitempty"`
}
