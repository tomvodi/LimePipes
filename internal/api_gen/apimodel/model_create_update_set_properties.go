/*
 * Set and Tune API
 *
 * API for managing sets and tunes
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apimodel

type CreateUpdateSetProperties struct {

	// The name of the Set
	Title string `json:"title"`

	// A description of the Set
	Description string `json:"description,omitempty"`

	// The name of the creator of the set
	Creator string `json:"creator,omitempty"`

	Tunes []int64 `json:"tunes,omitempty"`
}
