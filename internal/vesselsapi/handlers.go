package vesselsapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Log   logrus.FieldLogger
	Store Store
}

func (h *Handler) GetVesselByIMO(c *gin.Context) {
	h.Log.Info("Here")
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
