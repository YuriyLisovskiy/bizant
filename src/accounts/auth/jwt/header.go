// Copyright (c) 2018 Yuriy Lisovskiy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package jwt

// Represents JWT Header.
// Contains important information for encrypting/decrypting.
type Header struct {
	// Token type.
	Typ string `json:"typ"`

	// Message authentication code algorithm - the issuer can freely set an algorithm
	// to verify the signature on the token. However, some asymmetrical algorithms
	// pose security concerns.
	Alg string `json:"alg"`

	// Content type - this always is JWT.
	Cty string `json:"cty"`
}
