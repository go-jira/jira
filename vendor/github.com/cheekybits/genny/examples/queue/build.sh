#!/bin/bash
cat ./queue_generic.go | ../../genny gen "Generic=string,int" > queue_generic_gen.go
