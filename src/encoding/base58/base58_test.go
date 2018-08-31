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

package base58

import "testing"

var Base58Decode_Data = []struct {
	input    []byte
	expected []byte
}{
	{
		input:    []byte("14QDwy7VLCMuvcaSAXcKrqAd3XtVzQWdfiVh14RGLbYepbFQpc8TeFy8Ctrex6tB8Kctj94KrTTY84tpzwRFeujdH7mYJV4qK"),
		expected: []byte{0, 77, 97, 108, 101, 115, 117, 97, 100, 97, 32, 110, 105, 115, 108, 32, 117, 108, 116, 114, 105, 99, 105, 101, 115, 32, 99, 117, 114, 97, 101, 59, 46, 32, 76, 101, 99, 116, 117, 115, 32, 104, 97, 98, 105, 116, 97, 110, 116, 32, 100, 105, 115, 32, 112, 108, 97, 116, 101, 97, 46, 32, 65, 116, 32, 101, 116, 105, 97, 109, 46},
	},
	{
		input:    []byte("17vFeFEjfUUx1PU73y3pucVqAyEMnzLtsXjuCSWxqNd2os6sDWbY8GavEucb2v53D3GUhhHXvKgykcKySDPeY2fR1o"),
		expected: []byte{0, 78, 101, 113, 117, 101, 32, 109, 97, 115, 115, 97, 32, 118, 101, 108, 32, 118, 105, 118, 97, 109, 117, 115, 32, 84, 101, 108, 108, 117, 115, 32, 118, 97, 114, 105, 117, 115, 46, 32, 67, 108, 97, 115, 115, 32, 100, 105, 97, 109, 32, 102, 117, 115, 99, 101, 32, 105, 110, 116, 101, 114, 100, 117, 109, 46},
	},
	{
		input:    []byte("16X2XbHhikXCytvGxC5EFDz3so3Uhnr6LYFRroCCZdTn18xu8hcmXJ2ALPqAmoQmHnMaC7919EgBhwD2PK2aaLSyHg9Qw8FMtCZuxpariZDsq7WnyntV"),
		expected: []byte{0, 77, 97, 108, 101, 115, 117, 97, 100, 97, 32, 112, 111, 116, 101, 110, 116, 105, 32, 97, 108, 105, 113, 117, 101, 116, 32, 83, 111, 108, 108, 105, 99, 105, 116, 117, 100, 105, 110, 32, 114, 105, 100, 105, 99, 117, 108, 117, 115, 46, 32, 86, 101, 108, 105, 116, 32, 98, 105, 98, 101, 110, 100, 117, 109, 32, 109, 111, 110, 116, 101, 115, 32, 117, 114, 110, 97, 32, 116, 101, 109, 112, 111, 114, 46},
	},
	{
		input:    []byte("14EA8HigkVXU4FGjuTgivtskr5mmG12rtDqZf4cZ5dkd3No5FvQPeUjzdY7sVdgeF2iZCG5Lg1dtA51tdMqMS8XWobLEDdFky"),
		expected: []byte{0, 73, 110, 99, 101, 112, 116, 111, 115, 44, 32, 110, 97, 109, 32, 108, 101, 99, 116, 117, 115, 32, 117, 108, 116, 114, 105, 99, 101, 115, 32, 97, 117, 99, 116, 111, 114, 32, 113, 117, 105, 115, 113, 117, 101, 32, 105, 110, 32, 116, 111, 114, 116, 111, 114, 46, 32, 67, 111, 110, 117, 98, 105, 97, 32, 112, 111, 114, 116, 97, 46},
	},
	{
		input:    []byte("1zSF5ery8S7KV4gSraVeEKHXneb7DZZF8TUn9VQf5oe8ezgqMLiN3ecR75FokGtpKdRsB6LA8V2G4CpeUzP7"),
		expected: []byte{0, 73, 110, 46, 32, 86, 97, 114, 105, 117, 115, 32, 113, 117, 97, 109, 32, 110, 111, 115, 116, 114, 97, 32, 108, 97, 99, 105, 110, 105, 97, 32, 115, 101, 100, 32, 110, 105, 115, 105, 32, 110, 105, 115, 108, 32, 109, 111, 110, 116, 101, 115, 32, 115, 111, 99, 105, 111, 115, 113, 117, 46},
	},
	{
		input:    []byte("132h8WYAJtCqi6rSj3pbJ1YmV2o1kRTon6BzQA7as1iJLC6gMJUxTDdotoTjUTvF9dvsGze8GVKB1UYeihzQAM7Ck2BRX"),
		expected: []byte{0, 68, 117, 105, 32, 112, 108, 97, 99, 101, 114, 97, 116, 32, 101, 110, 105, 109, 32, 105, 100, 32, 99, 111, 110, 118, 97, 108, 108, 105, 115, 32, 115, 101, 110, 101, 99, 116, 117, 115, 46, 32, 80, 104, 97, 115, 101, 108, 108, 117, 115, 32, 105, 110, 32, 112, 104, 97, 114, 101, 116, 114, 97, 32, 110, 97, 109, 46},
	},
	{
		input:    []byte("13FdZRynTwFbNueGLxkRZKrsFb7H72p2QuP1L1tXrQKWVwBUcmeyvySq4nUFycjSbeGpCjRJ5zvrTRYSjdA4PcEaxdgcuNMG3KkdC3ogWb8R"),
		expected: []byte{0, 69, 108, 101, 105, 102, 101, 110, 100, 32, 118, 101, 110, 101, 110, 97, 116, 105, 115, 32, 109, 111, 108, 101, 115, 116, 105, 101, 32, 118, 101, 108, 32, 109, 97, 115, 115, 97, 32, 113, 117, 105, 115, 113, 117, 101, 32, 115, 111, 99, 105, 111, 115, 113, 117, 32, 99, 117, 109, 32, 99, 111, 110, 100, 105, 109, 101, 110, 116, 117, 109, 32, 112, 108, 97, 116, 101, 97, 46},
	},
	{
		input:    []byte("14aiJhcFoVZs1bvVFPtysL4s1ep79vqFrGw67no6f8k7zAKGzpRiKNDF8UzehS95EFgBPkuDUG5WwoNgDL7W4ZEcCTqfLn1VTcdqBMVHdCLqTfEu"),
		expected: []byte{0, 74, 117, 115, 116, 111, 32, 108, 101, 99, 116, 117, 115, 32, 108, 101, 99, 116, 117, 115, 32, 102, 97, 117, 99, 105, 98, 117, 115, 32, 108, 97, 111, 114, 101, 101, 116, 32, 112, 101, 108, 108, 101, 110, 116, 101, 115, 113, 117, 101, 32, 112, 108, 97, 99, 101, 114, 97, 116, 32, 108, 97, 99, 117, 115, 32, 100, 111, 108, 111, 114, 32, 99, 111, 110, 118, 97, 108, 108, 105, 115, 46},
	},
	{
		input:    []byte("13cPufgKmy5NUYKGFoYzrHRpi3WJmF4qz5Xay1napuxwLAK3A1unfDyzfbruSthKEprdxTMwgrcYPcU7ZLsaXxFqBvoD6uJBJryypYUMMhmP"),
		expected: []byte{0, 80, 117, 114, 117, 115, 32, 109, 101, 116, 117, 115, 32, 117, 116, 46, 32, 67, 111, 110, 100, 105, 109, 101, 110, 116, 117, 109, 46, 32, 83, 117, 115, 99, 105, 112, 105, 116, 44, 32, 102, 101, 114, 109, 101, 110, 116, 117, 109, 32, 105, 110, 116, 101, 103, 101, 114, 32, 118, 105, 118, 97, 109, 117, 115, 32, 112, 111, 114, 116, 97, 32, 116, 117, 114, 112, 105, 115, 46},
	},
	{
		input:    []byte("17Hqb1ZA3mRWydXgtGDK6vGQ6ydLTjcu4to5oUCvW2Pqrzd51BEauyeASX3dThHQesBBKkr1mzkJwivPA2Rypqo6DuqkpT8jGUvXNoSS1"),
		expected: []byte{0, 65, 32, 97, 108, 105, 113, 117, 101, 116, 32, 99, 111, 110, 100, 105, 109, 101, 110, 116, 117, 109, 32, 97, 99, 99, 117, 109, 115, 97, 110, 32, 115, 101, 110, 101, 99, 116, 117, 115, 32, 105, 110, 116, 101, 103, 101, 114, 46, 32, 76, 117, 99, 116, 117, 115, 32, 99, 117, 114, 115, 117, 115, 32, 101, 108, 105, 116, 32, 105, 110, 116, 101, 103, 101, 114, 46},
	},
	{
		input:    []byte("1XYBMqXpnYZa9ATfe1G1vaY8mcN9kx287hwHg1Gf5t1ZXqcRFsmWGNh4VS5hYuEZpQxgao7kJxhztJvj6Y6w4adigJo9oBeiJ7SW51KAhSYVMYh91Z55HvcRB"),
		expected: []byte{0, 65, 114, 99, 117, 44, 32, 115, 117, 115, 112, 101, 110, 100, 105, 115, 115, 101, 44, 32, 116, 105, 110, 99, 105, 100, 117, 110, 116, 32, 115, 117, 115, 112, 101, 110, 100, 105, 115, 115, 101, 32, 102, 97, 117, 99, 105, 98, 117, 115, 32, 100, 105, 97, 109, 32, 99, 111, 109, 109, 111, 100, 111, 32, 100, 105, 99, 116, 117, 109, 115, 116, 44, 32, 116, 111, 114, 113, 117, 101, 110, 116, 46, 32, 69, 114, 97, 116, 46},
	},
	{
		input:    []byte("1EH8nWtjCrkYSDs7jXHpkJ2DFMxKhuxZWLk7sFqYx7o1V3Q6xp92wKGcFx5nCvbEoLdWhRssLqAPhzG2bkRnRkfZPpanLNMsvM"),
		expected: []byte{0, 68, 117, 105, 32, 99, 117, 98, 105, 108, 105, 97, 32, 102, 97, 99, 105, 108, 105, 115, 105, 32, 110, 101, 116, 117, 115, 32, 105, 110, 116, 101, 103, 101, 114, 32, 102, 97, 99, 105, 108, 105, 115, 105, 32, 115, 117, 115, 99, 105, 112, 105, 116, 32, 101, 116, 32, 99, 111, 110, 100, 105, 109, 101, 110, 116, 117, 109, 32, 65, 116, 46},
	},
	{
		input:    []byte("14d38gbgsEi2sRdEdX5cuVKRKqbvzdKfEnrkMXroKvwx47CRtcuMfgWa1ZZ3cYr62LqoBuG3te4XiqhXCDqtQoWPHLWKnBxYD"),
		expected: []byte{0, 82, 104, 111, 110, 99, 117, 115, 32, 113, 117, 97, 109, 32, 101, 114, 111, 115, 32, 102, 101, 108, 105, 115, 32, 80, 104, 97, 114, 101, 116, 114, 97, 32, 112, 111, 114, 116, 116, 105, 116, 111, 114, 32, 110, 105, 98, 104, 32, 115, 101, 109, 32, 116, 105, 110, 99, 105, 100, 117, 110, 116, 32, 97, 108, 105, 113, 117, 97, 109, 46},
	},
}

func TestBase58Decode(test *testing.T) {
	for j, data := range Base58Decode_Data {
		actual := Decode(data.input)
			if len(actual) != len(data.expected) {
				test.Errorf("base58.TestBase58Decode[%d], size: %d != %d", j, len(actual), len(data.expected))
			}
			for i, b := range actual {
				if b != data.expected[i] {
					test.Errorf("base58.TestBase58Decode[%d], data: %x !=\n\t\t\t\t\t%x", j, actual, data.expected)
				}
			}
	}
}

var Base58Encode_Data = []struct {
	input    []byte
	expected []byte
}{
	{
		input:    []byte("Malesuada nisl ultricies curae;. Lectus habitant dis platea. At etiam."),
		expected: []byte("14QDwy7VLCMuvcaSAXcKrqAd3XtVzQWdfiVh14RGLbYepbFQpc8TeFy8Ctrex6tB8Kctj94KrTTY84tpzwRFeujdH7mYJV4qK"),
	},
	{
		input:    []byte("Neque massa vel vivamus Tellus varius. Class diam fusce interdum."),
		expected: []byte("17vFeFEjfUUx1PU73y3pucVqAyEMnzLtsXjuCSWxqNd2os6sDWbY8GavEucb2v53D3GUhhHXvKgykcKySDPeY2fR1o"),
	},
	{
		input:    []byte("Malesuada potenti aliquet Sollicitudin ridiculus. Velit bibendum montes urna tempor."),
		expected: []byte("16X2XbHhikXCytvGxC5EFDz3so3Uhnr6LYFRroCCZdTn18xu8hcmXJ2ALPqAmoQmHnMaC7919EgBhwD2PK2aaLSyHg9Qw8FMtCZuxpariZDsq7WnyntV"),
	},
	{
		input:    []byte("Inceptos, nam lectus ultrices auctor quisque in tortor. Conubia porta."),
		expected: []byte("14EA8HigkVXU4FGjuTgivtskr5mmG12rtDqZf4cZ5dkd3No5FvQPeUjzdY7sVdgeF2iZCG5Lg1dtA51tdMqMS8XWobLEDdFky"),
	},
	{
		input:    []byte("In. Varius quam nostra lacinia sed nisi nisl montes sociosqu."),
		expected: []byte("1zSF5ery8S7KV4gSraVeEKHXneb7DZZF8TUn9VQf5oe8ezgqMLiN3ecR75FokGtpKdRsB6LA8V2G4CpeUzP7"),
	},
	{
		input:    []byte("Dui placerat enim id convallis senectus. Phasellus in pharetra nam."),
		expected: []byte("132h8WYAJtCqi6rSj3pbJ1YmV2o1kRTon6BzQA7as1iJLC6gMJUxTDdotoTjUTvF9dvsGze8GVKB1UYeihzQAM7Ck2BRX"),
	},
	{
		input:    []byte("Eleifend venenatis molestie vel massa quisque sociosqu cum condimentum platea."),
		expected: []byte("13FdZRynTwFbNueGLxkRZKrsFb7H72p2QuP1L1tXrQKWVwBUcmeyvySq4nUFycjSbeGpCjRJ5zvrTRYSjdA4PcEaxdgcuNMG3KkdC3ogWb8R"),
	},
	{
		input:    []byte("Justo lectus lectus faucibus laoreet pellentesque placerat lacus dolor convallis."),
		expected: []byte("14aiJhcFoVZs1bvVFPtysL4s1ep79vqFrGw67no6f8k7zAKGzpRiKNDF8UzehS95EFgBPkuDUG5WwoNgDL7W4ZEcCTqfLn1VTcdqBMVHdCLqTfEu"),
	},
	{
		input:    []byte("Purus metus ut. Condimentum. Suscipit, fermentum integer vivamus porta turpis."),
		expected: []byte("13cPufgKmy5NUYKGFoYzrHRpi3WJmF4qz5Xay1napuxwLAK3A1unfDyzfbruSthKEprdxTMwgrcYPcU7ZLsaXxFqBvoD6uJBJryypYUMMhmP"),
	},
	{
		input:    []byte("A aliquet condimentum accumsan senectus integer. Luctus cursus elit integer."),
		expected: []byte("17Hqb1ZA3mRWydXgtGDK6vGQ6ydLTjcu4to5oUCvW2Pqrzd51BEauyeASX3dThHQesBBKkr1mzkJwivPA2Rypqo6DuqkpT8jGUvXNoSS1"),
	},
	{
		input:    []byte("Arcu, suspendisse, tincidunt suspendisse faucibus diam commodo dictumst, torquent. Erat."),
		expected: []byte("1XYBMqXpnYZa9ATfe1G1vaY8mcN9kx287hwHg1Gf5t1ZXqcRFsmWGNh4VS5hYuEZpQxgao7kJxhztJvj6Y6w4adigJo9oBeiJ7SW51KAhSYVMYh91Z55HvcRB"),
	},
	{
		input:    []byte("Dui cubilia facilisi netus integer facilisi suscipit et condimentum At."),
		expected: []byte("1EH8nWtjCrkYSDs7jXHpkJ2DFMxKhuxZWLk7sFqYx7o1V3Q6xp92wKGcFx5nCvbEoLdWhRssLqAPhzG2bkRnRkfZPpanLNMsvM"),
	},
	{
		input:    []byte("Rhoncus quam eros felis Pharetra porttitor nibh sem tincidunt aliquam."),
		expected: []byte("14d38gbgsEi2sRdEdX5cuVKRKqbvzdKfEnrkMXroKvwx47CRtcuMfgWa1ZZ3cYr62LqoBuG3te4XiqhXCDqtQoWPHLWKnBxYD"),
	},
}

func TestBase58Encode(test *testing.T) {
	for j, data := range Base58Encode_Data {
		actual := Encode(data.input)
		if len(actual) != len(data.expected) {
			test.Errorf("base58.TestBase58Encode[%d], size: %d != %d", j, len(actual), len(data.expected))
		}
		for i, b := range actual {
			if b != data.expected[i] {
				test.Errorf("base58.TestBase58Encode[%d], data: %s != %s", j, actual, data.expected)
			}
		}
	}
}

var ReverseBytes_Data = []struct{
	input []byte
	expected []byte
} {
	{
		input: []byte{0, 0, 0, 0, 0, 0, 0, 10},
		expected: []byte{10, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 54, 157, 86, 18, 54},
		expected: []byte{54, 18, 86, 157, 54, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 0, 0, 1, 35, 66},
		expected: []byte{66, 35, 1, 0, 0, 0, 0, 0},
	},
	{
		input: []byte{0, 0, 0, 0, 0, 2, 233, 45},
		expected: []byte{45, 233, 2, 0, 0, 0, 0, 0},
	},
}

func TestReverseBytes(test *testing.T) {
	for j, data := range ReverseBytes_Data {
		actual := data.input
		reverseBytes(actual)
		if len(actual) != len(data.expected) {
			test.Errorf("base58.TestReverseBytes[%d], size: %d != %d", j, len(actual), len(data.expected))
		}
		for i, b := range actual {
			if b != data.expected[i] {
				test.Errorf("base58.TestReverseBytes[%d], data: %s != %s", j, actual, data.expected)
			}
		}
	}
}
