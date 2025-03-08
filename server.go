package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

const DefaultHost = "localhost"
const DefaultPort = 8001

type RoutesMap map[string][]string

type Server struct {
	Host     string
	Port     int
	Handlers []*Handler
	// To create something like that: "GET": {"/books", "/libros", "/llibres"},
	Routes RoutesMap
	server *gin.Engine
}

func NewServer() *Server {
	s := &Server{}
	s.Reset()
	debugPrint("Reset!")
	return s
}

// Start
func (s *Server) Start() error {
	return s.server.Run(s.Host + ":" + fmt.Sprintf("%d", s.Port))
}

func (s *Server) Reset() {
	s.Host = DefaultHost
	s.Port = DefaultPort
	s.Handlers = nil //Best way to clear a slice
	s.Routes = make(RoutesMap)
	gin.SetMode(gin.ReleaseMode) // Avoid debug messages
	s.server = gin.Default()
}

// Implement `print` trait (for debug)
func (s *Server) String() string {
	// Convert to indented JSON
	jsonData, _ := json.MarshalIndent(s, "", "  ") // 2 spaces as TAB
	return string(jsonData)
}

// AddHandler Adds a handler or "virtual server"
func (s *Server) AddHandler(h *Handler) {
	s.Handlers = append(s.Handlers, h)
	// Set active methods only
	for method, active := range h.Methods {
		if active {
			s.Routes[method] = h.Paths
		}
	}
	s.addPathsToFunction(s.Routes)
}

// addPathsToFunction: Used by AddHandle to add entries to the Request Handlers
// all handleXXX functions are generics.
func (s *Server) addPathsToFunction(rm RoutesMap) {
	for method, paths := range rm {
		for _, path := range paths {
			debugPrint("addPathsToFunction -> %s - %s", method, path)
			switch method {
			case "GET":
				s.server.GET(path, s.handleGET)
			case "POST":
				s.server.POST(path, s.handlePOST)
			case "PUT":
				s.server.PUT(path, s.handlePUT)
			case "DELETE":
				s.server.DELETE(path, s.handleDELETE)
			case "PATCH":
				s.server.PATCH(path, s.handlePATCH)
			case "HEAD":
				s.server.HEAD(path, s.handleGET)
			}
		}
	}
}

// Internal Handlers for all methods request
//
func (s *Server) handleGET(c *gin.Context) {
	handler := s.getHandlerFromPath(c.FullPath())
	// If not ID, returns all JSON fields
	if c.Param("id") == "" {
		c.JSON(http.StatusOK, handler.jsonData.GetItems())
		return
	}

	// If ID is not a number returns a Bad Request error.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// item will be the JSON structure
	var item map[string]any
	item, err = handler.jsonData.GetItemByID(id)
	// Check if ID exists
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
		return
	}
	c.JSON(http.StatusOK, item) // ie: c.JSON(200, gin.H{"id": id, "title": "Go best book!"})
}

func (s *Server) handlePOST(c *gin.Context) {
	debugPrint("FullPath(): %s. Request.URL.Path: %s", c.FullPath(), c.Request.URL.Path)
	handler := s.getHandlerFromPath(c.FullPath())
	if handler.ReadOnlyMode {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed in readonly mode"})
		return
	}

	// Get Raw Data from the Body
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read BODY"})
		return
	}

	// Add new item and check error.
	createdItem, err := handler.jsonData.AddItem(rawData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdItem)
}

func (s *Server) handlePUT(c *gin.Context) {
	debugPrint("FullPath(): %s. Request.URL.Path: %s", c.FullPath(), c.Request.URL.Path)
	handler := s.getHandlerFromPath(c.FullPath())
	if handler.ReadOnlyMode {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed in readonly mode"})
		return
	}
	// If ID is not a number returns a Bad Request error.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get Raw Data from the Body
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read BODY"})
		return
	}

	// Modify the item and get error
	updatedItem, err := handler.jsonData.UpdateItem(id, rawData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedItem)
}

func (s *Server) handleDELETE(c *gin.Context) {
	handler := s.getHandlerFromPath(c.FullPath())
	if handler.ReadOnlyMode {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed in readonly mode"})
		return
	}
	// If ID is not a number returns a Bad Request error.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = handler.jsonData.DeleteItem(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (s *Server) handlePATCH(c *gin.Context) {
	handler := s.getHandlerFromPath(c.FullPath())
	if handler.ReadOnlyMode {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed in readonly mode"})
		return
	}

	// If ID is not a number returns a Bad Request error.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get Raw Data from the Body
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read BODY"})
		return
	}

	// Patch the item and get error
	patchedItem, err := handler.jsonData.PatchItem(id, rawData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patchedItem)
}

// getHandlerFromPath Returns the Handler associated to the path
func (s *Server) getHandlerFromPath(path string) *Handler {
	var result *Handler = nil
	var exists bool
	for _, h := range s.Handlers {
		exists = slices.Contains(h.Paths, path)
		if exists {
			result = h
			break
		}
	}
	return result
}
