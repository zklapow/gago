package util

func In(val int, arr []int) bool {
    for _, d := range arr {
        if val == d {
            return true
        }
    }

    return false
}
