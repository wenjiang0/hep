// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-dump dumps the content of a ROOT file, including the content of
// the Trees (for all entries), if any.
//
// Example:
//
//  $> root-dump ./testdata/small-flat-tree.root
//  >>> file[./testdata/small-flat-tree.root]
//  key[000]: tree;1 "my tree title" (TTree)
//  [000][Int32]: 0
//  [000][Int64]: 0
//  [000][UInt32]: 0
//  [000][UInt64]: 0
//  [000][Float32]: 0
//  [000][Float64]: 0
//  [000][Str]: evt-000
//  [000][ArrayInt32]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayInt64]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayInt32]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayInt64]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayFloat32]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayFloat64]: [0 0 0 0 0 0 0 0 0 0]
//  [000][N]: 0
//  [000][SliceInt32]: []
//  [000][SliceInt64]: []
//  [...]
//
//  $> root-dump -h
//  Usage: root-dump [options] f0.root [f1.root [...]]
//
//  ex:
//   $> root-dump ./testdata/small-flat-tree.root
//   $> root-dump -deep=0 ./testdata/small-flat-tree.root
//
//  options:
//    -deep
//      	enable deep dumping of values (including Trees' entries) (default true)
//    -name string
//      	regex of object names to dump
//
package main // import "go-hep.org/x/hep/groot/cmd/root-dump"

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"golang.org/x/xerrors"
)

var (
	deepFlag = flag.Bool("deep", true, "enable deep dumping of values (including Trees' entries)")
	nameFlag = flag.String("name", "", "regex of object names to dump")
)

func main() {
	log.SetPrefix("root-dump: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: root-dump [options] f0.root [f1.root [...]]

ex:
 $> root-dump ./testdata/small-flat-tree.root
 $> root-dump -deep=0 ./testdata/small-flat-tree.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *nameFlag != "" {
		reName = regexp.MustCompile(*nameFlag)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatalf("need at least one input ROOT file")
	}

	for _, fname := range flag.Args() {
		err := dump(os.Stdout, fname, *deepFlag)
		if err != nil {
			log.Fatalf("error dumping file %q: %+v", fname, err)
		}
	}
}

func dump(w io.Writer, fname string, deep bool) error {
	fmt.Fprintf(w, ">>> file[%s]\n", fname)
	f, err := groot.Open(fname)
	if err != nil {
		return xerrors.Errorf("could not open file with read-access: %w", err)
	}
	defer f.Close()

	return dumpDir(w, f, deep)
}

func dumpDir(w io.Writer, dir riofs.Directory, deep bool) error {
	for i, key := range dir.Keys() {
		fmt.Fprintf(w, "key[%03d]: %s;%d %q (%s)", i, key.Name(), key.Cycle(), key.Title(), key.ClassName())
		if !(deep && match(key.Name())) {
			fmt.Fprint(w, "\n")
			continue
		}
		obj, err := key.Object()
		if err != nil {
			return xerrors.Errorf("could not decode object %q from dir %q: %w", key.Name(), dir.(root.Named).Name(), err)
		}
		err = dumpObj(w, obj, deep)
		if xerrors.Is(err, errIgnoreKey) {
			err = nil
			continue
		}
		if err != nil {
			return xerrors.Errorf("error dumping key %q: %w", key.Name(), err)
		}
	}
	return nil
}

var reName *regexp.Regexp

func match(name string) bool {
	if reName == nil {
		return true
	}
	return reName.MatchString(name)
}

var errIgnoreKey = xerrors.Errorf("root-dump: ignore key")

func dumpObj(w io.Writer, obj root.Object, deep bool) error {
	var err error
	switch obj := obj.(type) {
	case rtree.Tree:
		fmt.Fprintf(w, "\n")
		err = dumpTree(w, obj)
	case riofs.Directory:
		fmt.Fprintf(w, "\n")
		err = dumpDir(w, obj, deep)
	case rhist.H2:
		fmt.Fprintf(w, "\n")
		err = dumpH2(w, obj)
	case rhist.H1: // keep after rhist.H2
		fmt.Fprintf(w, "\n")
		err = dumpH1(w, obj)
	case rhist.Graph:
		fmt.Fprintf(w, "\n")
		err = dumpGraph(w, obj)
	case root.List:
		fmt.Fprintf(w, "\n")
		err = dumpList(w, obj, deep)
	case *rdict.Object:
		fmt.Fprintf(w, " => %v\n", obj)
	case fmt.Stringer:
		fmt.Fprintf(w, " => %q\n", obj.String())
	default:
		fmt.Fprintf(w, " => ignoring key of type %T\n", obj)
		return errIgnoreKey
	}
	return err
}

func dumpList(w io.Writer, lst root.List, deep bool) error {
	for i := 0; i < lst.Len(); i++ {
		fmt.Fprintf(w, "lst[%s][%d]: ", lst.Name(), i)
		err := dumpObj(w, lst.At(i), deep)
		if err != nil && !xerrors.Is(err, errIgnoreKey) {
			return xerrors.Errorf("could not dump list: %w", err)
		}
	}
	return nil
}

func dumpTree(w io.Writer, t rtree.Tree) error {

	vars := rtree.NewScanVars(t)
	sc, err := rtree.NewScannerVars(t, vars...)
	if err != nil {
		return xerrors.Errorf("could not create scanner-vars: %w", err)
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan()
		if err != nil {
			return xerrors.Errorf("error scanning entry %d: %w", sc.Entry(), err)
		}
		for _, v := range vars {
			rv := reflect.Indirect(reflect.ValueOf(v.Value))
			fmt.Fprintf(w, "[%03d][%s]: %v\n", sc.Entry(), v.Name, rv.Interface())
		}
	}
	return nil
}

func dumpH1(w io.Writer, h1 rhist.H1) error {
	h, err := rootcnv.H1D(h1)
	if err != nil {
		return xerrors.Errorf("could not convert TH1x to hbook: %w", err)
	}
	return yodacnv.Write(w, h)
}

func dumpH2(w io.Writer, h2 rhist.H2) error {
	h, err := rootcnv.H2D(h2)
	if err != nil {
		return xerrors.Errorf("could not convert TH2x to hbook: %w", err)
	}
	return yodacnv.Write(w, h)
}

func dumpGraph(w io.Writer, gr rhist.Graph) error {
	g, err := rootcnv.S2D(gr)
	if err != nil {
		return xerrors.Errorf("could not convert TGraph to hbook: %w", err)
	}
	return yodacnv.Write(w, g)
}
