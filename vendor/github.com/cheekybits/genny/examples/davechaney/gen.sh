#!/bin/bash
cat ./generic_max.go | ../../genny gen "NumberType=NUMBERS" > numbers_max_get.go
cat ./func_thing.go | ../../genny gen "ThisNumberType=NUMBERS" > numbers_func_thing.go
