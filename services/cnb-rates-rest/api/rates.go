// Copyright (c) 2016-2020, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"net/http"

	"github.com/jancajthaml-openbank/cnb-rates-rest/persistence"
	"github.com/jancajthaml-openbank/cnb-rates-rest/utils"

	"github.com/gorilla/mux"
)

// RatesPartial returns http handler for single currency
func RatesPartial(server *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		currency := vars["currency"]

		if currency == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(emptyJSONObject)
			return
		}

		switch r.Method {

		case "GET":
			GetRates(server, currency, w, r)
			return

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(emptyJSONObject)
			return

		}
	}
}

// GetRates returns rates for given currency
func GetRates(server *Server, currency string, w http.ResponseWriter, r *http.Request) {
	tokens, err := persistence.LoadRates(server.Storage, currency)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONArray)
		return
	}

	resp, err := utils.JSON.Marshal(tokens)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONArray)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}
