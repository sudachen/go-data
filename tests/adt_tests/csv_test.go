package adt_tests

import (
	"fmt"
	"sudachen.xyz/pkg/go-data/adt"
	"sudachen.xyz/pkg/go-data/adt/csv"
	"sudachen.xyz/pkg/go-data/iokit"
	"testing"
)

var irisSource = iokit.Url("https://datahub.io/machine-learning/iris/r/iris.csv", iokit.Cache("dataset/iris.csv"))

func Test_IrisCsv_1(t *testing.T) {
	cls := adt.Enumset{}
	irisCsv := csv.Source(irisSource,
		csv.Float32("sepallength").As("Feature1"),
		csv.Float32("sepalwidth").As("Feature2"),
		csv.Float32("petallength").As("Feature3"),
		csv.Float32("petalwidth").As("Feature4"),
		csv.Meta(cls.Integer(), "class").As("label"))

	var q adt.Table
	irisCsv.MustDrain(q.Sink(),4)

	for i := 0; i < q.Len(); i++ {
		fmt.Println(q.Row(i))
	}
}

func Test_IrisCsv_2(t *testing.T) {
	cls := adt.Enumset{}
	irisCsv := csv.Source(irisSource,
		csv.Float32("sepal*").Group("F1"),
		csv.Float32("petal*").Group("F2"),
		csv.Meta(cls.Integer(), "class").As("label"))

	var q adt.Table
	irisCsv.MustDrain(q.Sink(),4)

	for i := 0; i < q.Len(); i++ {
		fmt.Println(q.Row(i))
	}
}

