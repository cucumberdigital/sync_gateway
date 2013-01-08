//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package basecouch

import (
	"net/http"
	"time"
)

const kDefaultSessionTTL = 24 * time.Hour

// GET /_session returns info about the current user
func (h *handler) handleSessionGET() error {
	if err := h.checkAuth(); err != nil {
		return err
	}
	userCopy := *h.user
	userCopy.Password = nil
	userCopy.PasswordHash = nil
	h.writeJSON(userCopy)
	return nil
}

// POST /_session creates a login session and sets its cookie
func (h *handler) handleSessionPOST() error {
	var params struct {
		Name string			`json:"name:name"`
		Password string		`json:"name:password"`
	}
	if err := readJSONInto(h.rq.Header, h.rq.Body, &params); err != nil {
		return err
	}
	var user *User
	user,_ = h.context.auth.GetUser(params.Name)
	if !user.Authenticate(params.Password) {
		user = nil
	}
	return h.makeSession(user)
}

func (h *handler) makeSession(user *User) error {
	if user == nil {
		return &HTTPError{http.StatusUnauthorized, "Invalid name/password"}
	}
	auth := h.context.auth
	session := auth.CreateSession(user.Name, kDefaultSessionTTL)
	http.SetCookie(h.response, auth.MakeSessionCookie(session))
	return nil
}
