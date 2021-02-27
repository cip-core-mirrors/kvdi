/*

Copyright 2020,2021 Avi Zimmerman

This file is part of kvdi.

kvdi is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

kvdi is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with kvdi.  If not, see <https://www.gnu.org/licenses/>.

*/

package v1

import (
	"regexp"
	"sort"
)

// Rule represents a set of permissions applied to a VDIRole. It mostly resembles
// an rbacv1.PolicyRule, with resources being a regex and the addition of a
// namespace selector.
type Rule struct {
	// The actions this rule applies for. VerbAll matches all actions.
	// Recognized options are: `["create", "read", "update", "delete", "use", "launch", "*"]`
	Verbs []Verb `json:"verbs,omitempty"`
	// Resources this rule applies to. ResourceAll matches all resources.
	// Recognized options are: `["users", "roles", "templates", "serviceaccounts", "*"]`
	Resources []Resource `json:"resources,omitempty"`
	// Resource regexes that match this rule. This can be template patterns, role
	// names or user names. There is no All representation because * will have
	// that effect on its own when the regex is evaluated. When referring to "serviceaccounts",
	// only the "use" verb is evaluated in the context of assuming those accounts in
	// desktop sessions.
	//
	// **NOTE**: The `kvdi-manager` is responsible for launching pods with a service account
	// requested for a given Desktop. If the service account itself contains more permissions
	// than the manager itself, the Kubernetes API will deny the request. The way to remedy this
	// would be to either mirror permissions to that ClusterRole, or make the `kvdi-manager` itself a
	// cluster admin, both of which come with inherent risks. In the end, you can decide the best
	// approach for your use case with regards to exposing access to the Kubernetes APIs via kvdi sessions.
	ResourcePatterns []string `json:"resourcePatterns,omitempty"`
	// Namespaces this rule applies to. Only evaluated for template launching
	// permissions. Including "*" as an option matches all namespaces.
	Namespaces []string `json:"namespaces,omitempty"`
}

// IsEmpty returns true if this rule is empty.
func (r *Rule) IsEmpty() bool {
	return len(r.Verbs) == 0 &&
		len(r.Resources) == 0 &&
		len(r.ResourcePatterns) == 0 &&
		len(r.Namespaces) == 0
}

// DeepEqual returns true if the provided rule matches this one exactly. All values in both rules
// are first sorted and then equality is derived from whether all fields pass reflect.DeepEqual.
func (r *Rule) DeepEqual(rule Rule) bool {
	this := r.DeepCopy()
	that := rule.DeepCopy()

	thisResourceStrings := resourcesToStrings(this.Resources)
	thisVerbStrings := verbsToStrings(this.Verbs)
	thatResourceStrings := resourcesToStrings(that.Resources)
	thatVerbStrings := verbsToStrings(that.Verbs)

	sort.Strings(thisResourceStrings)
	sort.Strings(thisVerbStrings)
	sort.Strings(thatResourceStrings)
	sort.Strings(thatVerbStrings)

	sort.Strings(this.ResourcePatterns)
	sort.Strings(this.Namespaces)
	sort.Strings(that.ResourcePatterns)
	sort.Strings(that.Namespaces)

	return strSliceEqual(thisResourceStrings, thatResourceStrings) &&
		strSliceEqual(thisVerbStrings, thatVerbStrings) &&
		strSliceEqual(this.ResourcePatterns, that.ResourcePatterns) &&
		strSliceEqual(this.Namespaces, that.Namespaces)
}

func strSliceEqual(ss, xx []string) bool {
	lS := len(ss)
	lX := len(xx)
	if lS == 0 && lX == 0 {
		return true
	}
	if lS != lX {
		return false
	}
	if lS > 0 {
		for i, s := range ss {
			if xx[i] == s {
				continue
			}
			return false
		}
		return true
	}
	for i, x := range xx {
		if ss[i] == x {
			continue
		}
		return false
	}
	return true
}

// HasVerb returns true if this rule contains the given verb.
func (r *Rule) HasVerb(verb Verb) bool {
	for _, item := range r.Verbs {
		if item == VerbAll {
			return true
		}
		if item == verb {
			return true
		}
	}
	return false
}

// HasResourceType returns true if this rule has the given resource type.
func (r *Rule) HasResourceType(resource Resource) bool {
	for _, item := range r.Resources {
		if item == ResourceAll {
			return true
		}
		if item == resource {
			return true
		}
	}
	return false
}

// MatchesResourceName returns true if any of the resource patterns in this rule
// match the given name.
func (r *Rule) MatchesResourceName(name string) bool {
	for _, pattern := range r.ResourcePatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			// Should have an external validator to let the user know
			// there is a bad regex.
			continue
		}
		if re.MatchString(name) {
			return true
		}
	}
	return false
}

// HasNamespace returns true if this rule includes the given namespace.
func (r *Rule) HasNamespace(ns string) bool {
	for _, item := range r.Namespaces {
		if item == NamespaceAll {
			return true
		}
		if item == ns {
			return true
		}
	}
	return false
}
