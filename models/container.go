/*
 * Copyright 2020 Roman Miro
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Container struct {
	ID             string            `json:"id" db:"id"`
	Name           string            `json:"name" db:"name"`
	OSTemplate     string            `json:"ostemplate" db:"os_template"`
	Parameters     map[string]string `json:"parameters" db:"-"`
	ParametersJSON sql.NullString    `db:"parameters"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
}

func (c *Container) UnmarshalParametersDB() error {
	if !c.ParametersJSON.Valid {
		return nil
	}

	return json.Unmarshal([]byte(c.ParametersJSON.String), &c.Parameters)
}
