package tests.public

import rego.v1

default allow := true

deny if {
	input.method == "PUT"
	some petid
	input.path = ["pets", petid]
	input.user == input.owner
}