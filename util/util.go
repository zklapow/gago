package util

import (
    "math/rand"
    "os"
    "log"
)

func In(val int, arr []int) bool {
    for _, d := range arr {
        if val == d {
            return true
        }
    }

    return false
}

func Random(min, max int) int {
    //cur := time.Now()
    //rand.Seed(RandSeed())
    return rand.Intn(max - min) + min
}

func Randomf(min, max float32) float32 {
    //cur := time.Now()
    //rand.Seed(RandSeed())
    return (rand.Float32() * (max-min)) + min
}

func RandSeed() int64 {
    f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
    // Get 1 bytes at a time
    tmp := make([]byte, 1)
    _, err := f.Read(tmp)
    if err != nil {
        log.Fatal(err)
    }

    return int64(tmp[0])
}
