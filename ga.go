package main

import (
    "fmt"
    "math/rand"
    "time"
    "sort"
//    "errors"
)

var target = []byte("Hi my name is Ze'ev")

type Genome struct {
    genes []byte
    score int
}
type Genomes []*Genome

func (s Genomes) Len() int {return len(s)}
func (s Genomes) Swap(i, j int) {s[i], s[j] = s[j], s[i]}

type ByScore struct { Genomes }

func (s ByScore) Less(i, j int) bool {return s.Genomes[i].score < s.Genomes[j].score}

func main() {
    cur := time.Now()
    rand.Seed(cur.Unix())

    //for i := 0; i < 20; i++ {
    //    base[i] = byte(rand.Intn(128)) 
    //}

    gen := initial()
    for i := range gen {
        fmt.Printf("Inital %v: %v\n", gen[i].score, string(gen[i].genes))
    }

    newgen(gen)

    //for bs > 0 {
    //    newgen(gen)
    //}
}

func twodslice(x, y int) (arr [][]byte) {
    arr = make([][]byte, x)
    for i := range arr {
        arr[i] = make([]byte, y)
    }

    return arr
}

func initial() (gen []*Genome) {
    gen = make([]*Genome, 10)
    for i := range gen {
        gen[i] = &Genome{}
        gen[i].genes = make([]byte, 20)
        for j := range gen[i].genes {
            gen[i].genes[j] = byte(random(32, 126))
        }
        gen[i].score = score(gen[i].genes)
    }

    return gen
}

func random(min, max int) int {
    return rand.Intn(max - min) + min
}

func score(ind []byte) (dist int) {
    //fmt.Printf("Target: %v Val: %v", len(target), len(ind))
    dist = 0
    for i, val := range ind {
        if i < len(target) && i < len(ind) {
            if val != target[i] {
                dist++
            }
        } else if val != byte(' ') {
            dist ++
        }
    }

    return dist
}

func mutate(ind []byte, rate float32) (err error) {
    if rate >= 1 {
        err := fmt.Errorf("Error rate must be less than one! You had %v", rate)
        return err
    }

    change := rand.Perm(len(ind))

    for _, val := range change {
        ind[val] = byte(random(32, 126))
    }

    return nil
}

func newgen(parent []*Genome) (child []*Genome) {
    sort.Sort(ByScore{parent})

    best := parent[0:4]

    for i := range best {
        fmt.Printf("Sorted %v: %v\n", parent[i].score, string(parent[i].genes))
    }

    return parent
}
