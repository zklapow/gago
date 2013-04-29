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
    "flag"
)

const (
    kGENERATION_SIZE = 1000
    kSELECT_SIZE = 100
    kMUTATION_RATE = 5
)

var kCH_LEN int 

type Genome struct {
    genes []byte
    score int
}

func (g *Genome) GeneString() (s string) {return string(g.genes)}

type Genomes []*Genome
func (s Genomes) Len() int {return len(s)}
func (s Genomes) Swap(i, j int) {s[i], s[j] = s[j], s[i]}

type ByScore struct { Genomes }
func (s ByScore) Less(i, j int) bool {return s.Genomes[i].score < s.Genomes[j].score}

// Flags
var verbose = flag.Bool("v", false, "Turn on verbose output.")
var target = flag.String("target", "Hello World!", "The target string to search for")

var tb []byte

func main() {
    // Get flage vals
    flag.Parse()

    tb = []byte(*target)
    kCH_LEN = len(tb)

    rand.Seed(util.RandSeed())

    fmt.Printf("Searching for %v\n", string(tb))

    t := time.Now()
    gen := initial()

    if *verbose {
        // Log the entire intial generation
        for i, val := range gen {
            fmt.Printf("Inital %v: %v (%v)\n", i, val.GeneString(), val.score)
        }
    }

    newgen(gen)

    child := gen
    best := gen[0]
    gc := 0
    for best.score > 0 {
        if *verbose {
            fmt.Printf("Running generation %v\n", gc)
        }
        tmp := newgen(child)

        child = tmp

        // Find the best individual
        sort.Sort(ByScore{child})
        best = child[0]
        gc++

        if *verbose {
            fmt.Printf("Current best is: %v\n", best.GeneString())
        }
    }

    dur := time.Since(t)

    fmt.Printf("Solution found in %v\n", dur)
    fmt.Printf("Found solution in %v generations with score %v: %v\n", gc, best.score, string(best.genes)) 
}

func initial() (gen []*Genome) {
    gen = make([]*Genome, kGENERATION_SIZE)
    data := make(chan *Genome, kGENERATION_SIZE)
    for _ = range gen {
        go func() {
            //rand.Seed(util.RandSeed())

            tmp := &Genome{}
            tmp.genes = make([]byte, kCH_LEN)
            for j := range tmp.genes {
                //rand.Seed(util.RandSeed())
                tmp.genes[j] = byte(util.Random(32, 126))
            }
            tmp.score = score(tmp.genes)

            data <- tmp
        }()
    }

    for i := range gen {
        tmp := <-data
        gen[i] = tmp
    }

    return gen
}

func score(ind []byte) (dist int) {
    //fmt.Printf("Target: %v Val: %v", len(target), len(ind))
    dist = 0
    for i, val := range ind {
        if i < len(tb) && i < len(ind) {
            if val != tb[i] {
                dist = dist + int(math.Abs(float64(int(val) - int(tb[i]))))
            }
        } else if val != byte(' ') {
            dist ++
        }
    }

    return dist
}

func mutate(ind *Genome) {
    index := rand.Intn(len(ind.genes))

    ind.genes[index] = byte(util.Random(32, 126))
}

func breed(p1 *Genome, p2 *Genome) (child *Genome, err error) {
    if len(p1.genes) != len(p2.genes) {
        return nil, errors.New("Cannot breed parents with different chromosome lengths!")
    }

    // Change a random number of genes and pick random alleles to change
    n := int(util.Randomf(0, 1) * float32(len(p1.genes)))
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

func splice(p1 *Genome, p2 *Genome) (child *Genome, child2 * Genome, err error) {
    if len(p1.genes) != len(p2.genes) {
        return nil, nil, errors.New("Cannot breed parents with different chromosome lengths!")
    }

    n := rand.Intn(len(p1.genes))

    child = &Genome{}
    child.genes = make([]byte, len(p1.genes))

    child2 = &Genome{}
    child2.genes = make([]byte, len(p1.genes))

    for i := range child.genes {
        if i <= n {
            child.genes[i] = p1.genes[i]
            child2.genes[i] = p2.genes[i]
        } else {
            child.genes[i] = p2.genes[i]
            child2.genes[i] = p1.genes[i]
        }
    }

    return child, child2, nil
}

func tourney(parent Genomes) (res *Genome) {
    a, b := rand.Intn(len(parent)), rand.Intn(len(parent))

    if parent[a].score > parent[b].score {
        return parent[a]
    } else {
        return parent[b]
    }

    return nil
}

func newgen(parent []*Genome) (child []*Genome) {
    // Reseed random
    //rand.Seed(util.RandSeed())

    sort.Sort(ByScore{parent})

    tpc := int(float32(len(parent)) * .01)
    best := parent[0:tpc]

    if *verbose {
        for i, val := range best {
            fmt.Printf("Best Candidate %v is: %v\n", i, val.GeneString())
        }
    }

    for len(best) < kSELECT_SIZE {
        best = append(best, tourney(parent))
    }

    child = make([]*Genome, kGENERATION_SIZE)
    data := make(chan *Genome, kGENERATION_SIZE)
    for i := 0; i < kGENERATION_SIZE/2; i++ {
        go func() {
            //rand.Seed(util.RandSeed())

            // Pick two parents from the best
            p1 := rand.Intn(kSELECT_SIZE)

            // Make sure we don't breed with the same parent
            p2 := p1
            for p2 != p1 {
                //rand.Seed(util.RandSeed())
                p2 = rand.Intn(kSELECT_SIZE)
            }

            // Breed the parents
            tmp1, tmp2, err := splice(best[p1], best[p2])
            if err != nil {
                log.Fatal(err)
            }

            // Score the final child genome
            tmp1.score = score(tmp1.genes)
            tmp2.score = score(tmp2.genes)

            data <- tmp1
            data <- tmp2
        }()
    }

    // Collect all the data
    for i := 0; i < kGENERATION_SIZE; i++ {
        tmp := <- data

        if (i % kMUTATION_RATE) == 0 {
            mutate(tmp)

            // rescore
            tmp.score = score(tmp.genes)
        }
        child[i] = tmp
    }

    // Keep the best parents in case of shitty offspring
    // child = append(child, best...)

    return child
}
