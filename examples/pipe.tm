#!/usr/bin/env tamarin

let array = ["gophers", "are", "burrowing", "rodents"]

let sentence = array | strings.join(" ") | strings.to_upper

print(sentence)
