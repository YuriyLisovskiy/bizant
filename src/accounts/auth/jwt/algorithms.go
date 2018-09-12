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

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
)

//HmacSha256 returns the SingingMethod for HMAC with SHA256
func HmacSha256(key string) JWT {
	return JWT{
		algorithm:   "HS256",
		signingHash: hmac.New(sha256.New, []byte(key)),
	}
}

//HmacSha512 returns the SigningMethod for HMAC with SHA512
func HmacSha512(key string) JWT {
	return JWT{
		algorithm:   "HS512",
		signingHash: hmac.New(sha512.New, []byte(key)),
	}
}

//HmacSha384 returns the SigningMethod for HMAC with SHA384
func HmacSha384(key string) JWT {
	return JWT{
		algorithm:   "HS384",
		signingHash: hmac.New(crypto.SHA384.New, []byte(key)),
	}
}
