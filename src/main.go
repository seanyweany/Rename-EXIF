package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
    "strconv"
)

func main() {
    files, err := filepath.Glob("*.*")
    if err != nil {
        log.Fatal(err)
    }

    // Loop through each file
    for _, file := range files {
        ext := filepath.Ext(file)      // Get the original extension

        cmd := exec.Command("exiftool", "-DateTimeOriginal", file)
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Printf("Command failed: %s\n", err)
            fmt.Printf("Output: %s\n", out)
        }

        // Extract the date and time from exiftool output
        re := regexp.MustCompile(`\d{4}:\d{2}:\d{2} \d{2}:\d{2}:\d{2}`)
        match := re.FindString(string(out))

        if match == "" {
            log.Printf("No DateTimeOriginal metadata found for file: %s\n", file)
            continue
        }

        // Parse the date and time
        datetimeParts := strings.Split(match, " ")
        dateParts := strings.Split(datetimeParts[0], ":")
        timeParts := strings.Split(datetimeParts[1], ":")

        year, _ := strconv.Atoi(dateParts[0])
        month, _ := strconv.Atoi(dateParts[1])
        day, _ := strconv.Atoi(dateParts[2])
        hour, _ := strconv.Atoi(timeParts[0])
        minute, _ := strconv.Atoi(timeParts[1])
        second, _ := strconv.Atoi(timeParts[2])

        // PDT = UTC-07:00
        // PST = UTC-08:00
        hour -= 0
        //minute += 30

        // Normalize time
        if minute >= 60 {
            hour += minute / 60
            minute = minute % 60
        }
        if hour >= 24 {
            day += hour / 24
            hour = hour % 24
            // Optionally handle month/day overflow here
        }

        // Generate the new filename
        newFilename := fmt.Sprintf("%d%02d%02d_%02d%02d%02d", year, month, day, hour, minute, second) + ext

        // Rename the file
        err = os.Rename(file, newFilename)
        if err != nil {
            log.Fatal(err)
        }

        // Check if there are duplicate filenames and rename them
        count := 2
        for fileExists(newFilename) {
            newFilename = fmt.Sprintf("%d%02d%02d_%02d%02d%02d_%d", year, month, day, hour, minute, second, count) + ext
            count++
        }
    }
}

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}