package interfaces

import (
	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source api_handler.go -destination ./mocks/mock_api_handler.go

type ApiHandler interface {
	Index(c *gin.Context)
	// AssignTunesToSet - Assign tunes to a set
	AssignTunesToSet(c *gin.Context)
	// CreateSet - Create a new set
	CreateSet(c *gin.Context)
	// CreateTune - Create a new tune
	CreateTune(c *gin.Context)
	// DeleteSet - Delete a set by ID
	DeleteSet(c *gin.Context)
	// DeleteTune - Delete a tune by ID
	DeleteTune(c *gin.Context)
	// GetSet - Get a set by ID
	GetSet(c *gin.Context)
	// GetTune - Get a tune by ID
	GetTune(c *gin.Context)
	// ListSets - List all sets
	ListSets(c *gin.Context)
	// ListTunes - List all tunes
	ListTunes(c *gin.Context)
	// UpdateSet - Update a set by ID
	UpdateSet(c *gin.Context)
	// UpdateTune - Update a tune by ID
	UpdateTune(c *gin.Context)
	// ImportBww - Import tunes/sets from one or more bww files
	ImportBww(c *gin.Context)
}
