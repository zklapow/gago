package main

import (
    "fmt"
    "math/rand"
    "math"
    "time"
    "sort"
    "errors"
    "gago/util"
    "log"
)

const (
    kGENERATION_SIZE = 1000
    kMUTATION_RATE = .6
)

var target = []byte("Hello World")
var kCH_LEN = len(target)

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

    t := time.Now()
    gen := initial()

    newgen(gen)

    child := gen
    best := gen[0]
    gc := 0
    for best.score > 0 {
        tmp := newgen(child)

        child = tmp

        // Find the best individual
        sort.Sort(ByScore{child})
        best = child[0]
        gc++

        //time.Sleep(time.Second*2)
    }

    dur := time.Since(t)

    fmt.Printf("Solution found in %v\n", dur)
    fmt.Printf("Found solution in %v generations with score %v: %v\n", gc, best.score, string(best.genes)) 
}

func initial() (gen []*Genome) {
    gen = make([]*Genome, kGENERATION_SIZE)
    for i := range gen {
        gen[i] = &Genome{}
        gen[i].genes = make([]byte, kCH_LEN)
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

func randomf(min, max float32) float32 {
    return (rand.Float32() * (max-min)) + min
}

func score(ind []byte) (dist int) {
    //fmt.Printf("Target: %v Val: %v", len(target), len(ind))
    dist = 0
    for i, val := range ind {
        if i < len(target) && i < len(ind) {
            if val != target[i] {
                dist = dist + int(math.Abs(float64(int(val) - int(target[i]))))
            }
        } else if val != byte(' ') {
            dist ++
        }
    }

    return dist
}

func mutate(ind *Genome) {
    index := rand.Intn(len(ind.genes))

    ind.genes[index] = byte(random(32, 126))
}

func breed(p1 *Genome, p2 *Genome) (child *Genome, err error) {
    if len(p1.genes) != len(p2.genes) {
        return nil, errors.New("Cannot breed parents with different chromosome lengths!")
    }

    // Change a random number of genes and pick random alleles to change
    n := int(randomf(0, 1) * float32(len(p1.genes)))
    change := rand.Perm(len(p1.genes))[:n] 

    // Initialize a new child
    child = &Genome{}
    child.genes = make([]byte, len(p1.genes))

    // Indexes in `change` come from p1 alleles
    // All others come from p2 alleles
    for i := range child.genes {
        if util.In(i, change) {
            child.genes[i] = p1.genes[i]
        } else {
            child.genes[i] = p2.genes[i]
        }
    }

    return child, nil
}

func splice(p1 *Genome, p2 *Genome) (child *Genome, err error) {
    n := rand.Intn(len(p1.genes))

    child = &Genome{}
    child.genes = make([]byte, len(p1.genes))

    for i := range child.genes {
        if i <= n {
            child.genes[i] = p1.genes[i]
        } else {
            child.genes[i] = p2.genes[i]
        }
    }

    return child, nil
}

func newgen(parent []*Genome) (child []*Genome) {
    // Reseed random
    cur := time.Now()
    rand.Seed(cur.Unix())

    sort.Sort(ByScore{parent})

    tpc := int(float32(len(parent)) * .01)
    best := parent[0:tpc]

    child = make([]*Genome, kGENERATION_SIZE)
    for i := 0; i < kGENERATION_SIZE; i++ {
        // Pick two parents from the best
        p1 := rand.Intn(tpc)

        // Make sure we don't breed with the same parent
        p2 := p1
        for p2 != p1 {
            p2 = rand.Intn(tpc)
        }

        // Breed the parents
        tmp, err := splice(best[p1], best[p2])
        if err != nil {
            log.Fatal(err)
        }

        // Mutate the child
        if randomf(0, 1) < kMUTATION_RATE {
            mutate(tmp)
            if err != nil {
                log.Fatal(err)
            }
        }

        // Score the final child genome
        child[i] = tmp
        child[i].score = score(child[i].genes)
    }

    // Keep the best parents in case of shitty offspring
    // child = append(child, best...)

    return child
}
