package rest
// sackci
// Copyright (C) 2017 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "time"
    "sync"
    "crypto/rand"
    "encoding/base64"
    "errors"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    // Length of the session token in bytes.
    TOKEN_SIZE = 64
)

var (
    ErrNoCookie = errors.New("cookie not found")
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type SessionStore struct {
    Sessions map[string]session

    mutex sync.Mutex
}

type session struct {
    Creation time.Time
    Timeout time.Duration
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Creates a new session store.
func NewSessionStore() (*SessionStore) {
    return &SessionStore{
        Sessions: make(map[string]session),
    }
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Checks if the given token belongs to a valid session.
func (s *SessionStore) IsValid(token string) (bool) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    // does the session exist?
    session, exists := s.Sessions[token]
    if !exists {
        return false
    }

    // check if the session is expired
    expires := session.Creation.Add(session.Timeout)
    if expires.Before(time.Now()) {
        delete(s.Sessions, token)
        return false
    }

    return true
}

// Creates a new session with the given expiration time.
func (s *SessionStore) Create(timeout time.Duration) (string, error) {
    token, err := s.newSessionToken()
    if err != nil {
        return "", err
    }

    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.Sessions[token] = session{
        Creation: time.Now(),
        Timeout: timeout,
    }

    return token, nil
}

// Deletes a session form the session store.
func (s *SessionStore) Delete(token string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, exists := s.Sessions[token]; exists {
        delete(s.Sessions, token)
    }
}

// Refreshs the given session.
func (s *SessionStore) Refresh(token string) {
    session, exists := s.Sessions[token]
    if exists {
        session.Creation = time.Now()
    }
}


// ----------------------------------------------------------------------------------
//  private members
// ----------------------------------------------------------------------------------

// Generates a new session token.
func (s *SessionStore) newSessionToken() (string, error) {
    // generate a cryptographically secure random token
    b := make([]byte, TOKEN_SIZE)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }

    // TODO: how to make sure this token is unique?

    return base64.StdEncoding.EncodeToString(b), nil
}
