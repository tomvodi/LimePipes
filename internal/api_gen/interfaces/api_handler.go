/*
 * Set and Tune API
 *
 * API for managing sets and tunes
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package interfaces

import (
	"github.com/gin-gonic/gin"
)

type ApiHandler interface {


    // AssignTunesToSet Put /sets/:setId/tunes
    // Assign tunes to a set 
     AssignTunesToSet(c *gin.Context)

    // CreateSet Post /sets
    // Create a new set 
     CreateSet(c *gin.Context)

    // CreateTune Post /tunes
    // Create a new tune 
     CreateTune(c *gin.Context)

    // DeleteSet Delete /sets/:setId
    // Delete a set by ID 
     DeleteSet(c *gin.Context)

    // DeleteTune Delete /tunes/:tuneId
    // Delete a tune by ID 
     DeleteTune(c *gin.Context)

    // GetSet Get /sets/:setId
    // Get a set by ID 
     GetSet(c *gin.Context)

    // GetTune Get /tunes/:tuneId
    // Get a tune by ID 
     GetTune(c *gin.Context)

    // Health Get /health
    // Check the health of the service 
     Health(c *gin.Context)

    // ImportBww Post /imports/bww
    // Import tunes/sets from one or more bww files 
     ImportBww(c *gin.Context)

    // ListSets Get /sets
    // List all sets 
     ListSets(c *gin.Context)

    // ListTunes Get /tunes
    // List all tunes 
     ListTunes(c *gin.Context)

    // UpdateSet Put /sets/:setId
    // Update a set by ID 
     UpdateSet(c *gin.Context)

    // UpdateTune Put /tunes/:tuneId
    // Update a tune by ID 
     UpdateTune(c *gin.Context)

}