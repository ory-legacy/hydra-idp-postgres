// Package handler
package identity

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/ory-am/hydra/herodot"
)

type Handler struct {
	he     herodot.Herodot
	store  IdentityManager
	warden warden.Warden
}

func permission(id string) string {
	return fmt.Sprintf("rn:idp.hydra:identities:%s", id)
}

func (h *Handler) SetRoutes(r *httprouter.Router) {
	r.GET("/accounts", h.create)
	r.GET("/accounts/:id", h.get)
	r.DELETE("/accounts/:id", h.delete)
	r.PUT("/accounts/:id/password", h.create)
	r.PUT("/accounts/:id/extra", h.updateData)
	r.PUT("/accounts/:id/username", h.updateUsername)
}

func (h *Handler) create(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var ctx = herodot.NewContext()
	var car CreateAccountRequest

	if err := json.NewDecoder(req.Body).Decode(&car); err != nil {
		h.he.WriteErrorCode(ctx, rw, req, err, http.StatusBadRequest)
		return
	}

	// Force ID override
	car.ID = uuid.New()
	user, err := h.store.Create(car)
	if err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	h.he.WriteCreated(ctx, rw, req, user, fmt.Sprintf("/accounts/%s", user.ID))
}

func (h *Handler) updateUsername(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var ctx = herodot.NewContext()
	var uur UpdateUsernameRequest
	var id = p.ByName("id")

	if id == "" {
		h.he.WriteErrorCode(ctx, rw, req, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	if _, err := h.warden.HTTPAuthorized(req, permission(id), "put:username", ladon.NewContext(req, id)); err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&uur); err != nil {
		h.he.WriteErrorCode(ctx, rw, req, err, http.StatusBadRequest)
		return
	}

	user, err := h.store.UpdateUsername(id, uur)
	if err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	h.he.Write(ctx, rw, req, user)
}

func (h *Handler) updatePassword(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var ctx = herodot.NewContext()
	var upr UpdatePasswordRequest
	var id = p.ByName("id")

	if id == "" {
		h.he.WriteErrorCode(ctx, rw, req, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	if _, err := h.warden.HTTPAuthorized(req, permission(id), "put:password", ladon.NewContext(req, id)); err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&upr); err != nil {
		h.he.WriteErrorCode(ctx, rw, req, err, http.StatusBadRequest)
		return
	}

	user, err := h.store.UpdatePassword(id, upr)
	if err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	h.he.Write(ctx, rw, req, user)
}

func (h *Handler) updateData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var ctx = herodot.NewContext()
	var udr UpdateDataRequest
	var id = p.ByName("id")

	if id == "" {
		h.he.WriteErrorCode(ctx, rw, req, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	if _, err := h.warden.HTTPAuthorized(req, permission(id), "put:data", ladon.NewContext(req, id)); err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&udr); err != nil {
		h.he.WriteErrorCode(ctx, rw, req, err, http.StatusBadRequest)
		return
	}

	user, err := h.store.UpdateData(id, udr)
	if err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	h.he.Write(ctx, rw, req, user)
}

func (h *Handler) get(ctx context.Context, rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = p.ByName("id")

	if id == "" {
		h.he.WriteErrorCode(ctx, rw, req, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	if _, err := h.warden.HTTPAuthorized(req, permission(id), "get", ladon.NewContext(req, id)); err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	user, err := h.store.Get(id)
	if err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	h.he.Write(ctx, rw, req, user)
}

func (h *Handler) delete(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var ctx = herodot.NewContext()
	var id = p.ByName("id")

	if id == "" {
		h.he.WriteErrorCode(ctx, rw, req, errors.Errorf("No id given."), http.StatusBadRequest)
		return
	}

	if _, err := h.warden.HTTPAuthorized(req, permission(id), "delete", ladon.NewContext(req, id)); err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	if err := h.store.Delete(id); err != nil {
		h.he.WriteError(ctx, rw, req, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}
