package vesselsapi

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/vesssels-api/pkg/vessels"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Log   logrus.FieldLogger
	Store Store
}

func (h *Handler) UpdateVessel(c *gin.Context) {
	var request vessels.UpdateVessel
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid json"})
		h.Log.WithError(err).Debug("AddVessel could not bind json")
		return
	}
	imoRaw := c.Param("imo")
	if imoRaw == "" {
		h.Log.Errorf("No imo provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "No imo provided"})
		return
	}
	imo, err := strconv.Atoi(imoRaw)
	if err != nil {
		h.Log.WithError(err).Errorf("Unable to convert imo to int: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "No valid imo provided"})
		return
	}

	err = h.Store.UpdateVessel(imo, &request)
	if err != nil {
		h.Log.WithError(err).Errorf("Unable to add vessel: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to add vessel"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Vessel added"})
}

func (h *Handler) GetVessels(c *gin.Context) {
	vessels, err := h.Store.GetVessels()
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get vessels: %s", err)
		h.Log.WithError(err).Error(errorMessage)
		return
	}
	c.JSON(http.StatusOK, vessels)
}

func (h *Handler) GetVesselByIMO(c *gin.Context) {
	imo, err := strconv.Atoi(c.Param("imo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No valid imo provided"})
		return
	}

	vessel, err := h.Store.GetVesselByIMO(imo)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to get vessel by imo: %s", err)
		h.Log.WithError(err).Error(errorMessage)
		return
	}
	if vessel == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Vessel not found"})
		return
	}
	c.JSON(http.StatusOK, vessel)
}

func (h *Handler) DeleteVessel(c *gin.Context) {
	imo, err := strconv.Atoi(c.Param("imo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No valid imo provided"})
		return
	}

	err = h.Store.DeleteVessel(imo)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to delete vessel: %s", err)
		h.Log.WithError(err).Error(errorMessage)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vessel deleted"})
}
