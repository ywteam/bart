package tests.zerotrust

import rego.v1

default allow := false

allow if {
	input.method == "PUT"
	some petid
	input.path = ["pets", petid]
	input.user == input.owner
}