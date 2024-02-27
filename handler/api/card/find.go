// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

// +build !oss

package card

import (
	"net/http"
	"strconv"

	"github.com/drone/drone/core"
	"github.com/drone/drone/handler/api/render"

	"github.com/go-chi/chi"
)

// HandleFind returns an http.HandlerFunc that writes a json-encoded
func HandleFind(
	buildStore core.BuildStore,
	cardStore core.CardStore,
	stageStore core.StageStore,
	stepStore core.StepStore,
	repoStore core.RepositoryStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			namespace = chi.URLParam(r, "owner")
			name      = chi.URLParam(r, "name")
		)

		buildNumber, err := strconv.ParseInt(chi.URLParam(r, "build"), 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		stageNumber, err := strconv.Atoi(chi.URLParam(r, "stage"))
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		stepNumber, err := strconv.Atoi(chi.URLParam(r, "step"))
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		repo, err := repoStore.FindName(r.Context(), namespace, name)
		if err != nil {
			render.NotFound(w, err)
			return
		}
		build, err := buildStore.FindNumber(r.Context(), repo.ID, buildNumber)
		if err != nil {
			render.NotFound(w, err)
			return
		}
		stage, err := stageStore.FindNumber(r.Context(), build.ID, stageNumber)
		if err != nil {
			render.NotFound(w, err)
			return
		}
		step, err := stepStore.FindNumber(r.Context(), stage.ID, stepNumber)
		if err != nil {
			render.NotFound(w, err)
			return
		}

		card, err := cardStore.Find(r.Context(), step.ID)
		if err != nil {
			render.NotFound(w, err)
			return
		}
		render.JSON(w, card, 200)
	}
}
