package crud

import (
	"code.vikunja.io/api/models"
	"github.com/labstack/echo"
	"net/http"
)

// DeleteWeb is the web handler to delete something
func (c *WebHandler) DeleteWeb(ctx echo.Context) error {
	// Bind params to struct
	if err := ParamBinder(c.CObject, ctx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL param.")
	}

	// Check if the user has the right to delete
	user, err := models.GetCurrentUser(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !c.CObject.CanDelete(&user) {
		models.Log.Noticef("%s [ID: %d] tried to delete while not having the rights for it", user.Username, user.ID)
		return echo.NewHTTPError(http.StatusForbidden)
	}

	err = c.CObject.Delete()
	if err != nil {
		return HandleHTTPError(err)
	}

	return ctx.JSON(http.StatusOK, models.Message{"Successfully deleted."})
}
