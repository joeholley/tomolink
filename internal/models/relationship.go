// Copyright 2019 Google LLC, with excerpts 2019 Matt Silverlock as noted
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

package models

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

//Relationship ...
type Relationship struct {
	Direction    string `json:"direction"`
	Relationship string `json:"relationship"`
	Delta        int    `json:"delta"`
	UUIDSource   string `json:"uuidsource"`
	UUIDTarget   string `json:"uuidtarget"`
}

//Validate ...
func (rel *Relationship) Validate() error {
	if rel.IsSingleDirection() == true ||
		rel.IsMultipleDirection() == true {
		return nil
	}
	return fmt.Errorf("invalid relationship direction '%v'", rel.Direction)
}

// Merge uses reflection to iterate over the values in both input relationships
// and attempt to merge them.  Note that it doesn't try to handle conflicts
// where both input relationship structs have defined different values for the
// same field - that just produces an error.
func (rel *Relationship) Merge(rel2 *Relationship) (*Relationship, error) {
	mergedRel := Relationship{}
	// Use reflection to get iterate over the values of both relationships
	//v := reflect.ValueOf(&rel).Elem()
	//w := reflect.ValueOf(&rel2).Elem()
	//m := reflect.ValueOf(&mergedRel).Elem()
	v := reflect.ValueOf(rel).Elem()
	w := reflect.ValueOf(rel2).Elem()
	//ve := reflect.ValueOf(&rel).Elem()
	//we := reflect.ValueOf(&rel2).Elem()
	m := reflect.ValueOf(&mergedRel).Elem()

	for i := 0; i < v.NumField(); i++ {
		sourceField := v.Field(i)
		otherField := w.Field(i)
		if sourceField.Interface() != otherField.Interface() {
			// The two structs have different values in this field; (at least)
			// one of them must be non-empty.
			// Initially the sourceField value is set with the assumption
			// that the first struct has the value for this field, and the
			// second struct doesn't. In that case, we want to populate the
			// value from the first struct to the merged struct, which
			// we're already set up to do.
			switch sourceField.Kind() {
			case reflect.String:
				m.Field(i).SetString(sourceField.String())
				if sourceField.String() == "" {
					// This field in the first struct is empty, so copy the
					// value from the second struct into the merged struct
					m.Field(i).SetString(otherField.String())

					//sourceField.SetString(otherField.String())
				} else if otherField.String() != "" {
					// The two structs both have non-empty values for this
					// field. Just fail so this can get sent back to the client
					// as an error (the client should re-submit this call
					// without conflicting values)
					// TODO: log the values
					return &mergedRel, errors.New("relationship field conflict - the same field has two different, non-empty values")
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				m.Field(i).SetInt(sourceField.Int())
				if sourceField.Int() == 0 {
					// This field in the first struct is empty, so copy the
					// value from the second struct into the merged stuct
					m.Field(i).SetInt(otherField.Int())
					//sourceField.SetInt(otherField.Int())
				} else if otherField.Int() != 0 {
					// The two structs both have non-empty values for this
					// field. Just fail so this can get sent back to the client
					// as an error (the client should re-submit this call
					// without conflicting values)
					// TODO: log the values
					return &mergedRel, errors.New("relationship field conflict - the same field has two different, non-empty values")
				}

			default:
				// This should never happen unless someone adds a new field to
				// the Relationships struct without understanding the code!!
				return &mergedRel, errors.New("Unhandled relationship field")
			}

		}

		// After the above logic, the sourceField should contain the value from
		// the two input relationship structs that we want to populate to the
		// final merged struct.  Now just figure out what kind of value it is
		// so it can be correctly copied.
		/*switch a := sourceField; a.Kind() {
		case reflect.String:
			m.Field(i).SetString(sourceField.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m.Field(i).SetInt(sourceField.Int())
		default:
			// This should never happen unless someone adds a new field to
			// the Relationships struct without understanding the code!!
			return &mergedRel, errors.New("Unhandled relationship field")
		}
		*/
	}
	// Assign the results back to the original relationship.
	return &mergedRel, nil
}

//VerboseLogger ...
func (rel *Relationship) VerboseLogger() *logrus.Entry {
	// Populate fields in the logrus structured logging
	logger := logrus.WithFields(logrus.Fields{
		"relationship": rel.Relationship,
		"delta":        rel.Delta,
		"direction":    rel.Direction,
		"uuidsource":   rel.UUIDSource,
		"uuidtarget":   rel.UUIDTarget,
	})

	return logger
}

//IsSingle ...
func (rel *Relationship) RelationshipInArray(relArray []string) bool {

	// Check if the relationship type in rel is contained in the input
	// array relArray.
	valid := false
	for _, r := range relArray {
		if rel.Relationship == r {
			valid = true
		}
	}

	return valid
}

//IsSingle ...
func (rel *Relationship) IsSingleDirection() bool {
	// Not specifying a direction defaults to source->target UUID relationship update
	if rel.Direction == "" ||
		rel.Direction == "unary" ||
		rel.Direction == "discrete" ||
		rel.Direction == "one" ||
		rel.Direction == "uni" ||
		rel.Direction == "single" {
		return true
	}
	return false
}

//IsMultiple ...
func (rel *Relationship) IsMultipleDirection() bool {
	if rel.Direction == "multiple" ||
		rel.Direction == "reciprocal" ||
		rel.Direction == "bi" ||
		rel.Direction == "mutual" {
		return true
	}
	return false
}

// IsValid ...
func (rel *Relationship) IsValid() bool {

	return false
}
