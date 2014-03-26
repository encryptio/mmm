package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const CHUNK_SIZE = 128

type Pair struct {
	x, y int
}

func Worker(commonArgs []string, outputPath string, toWorker chan Pair, fromWorker chan Pair, badFromWorker chan Pair) {
	for this := range toWorker {
		fmt.Printf("rendering %d,%d\n", this.x, this.y)

		args := make([]string, len(commonArgs))
		copy(args, commonArgs)
		args = append(args, []string{
			fmt.Sprintf("-x=%v", this.x*CHUNK_SIZE),
			fmt.Sprintf("-z=%v", this.y*CHUNK_SIZE),
			fmt.Sprintf("-o=%v", filepath.Join(outputPath, fmt.Sprintf("map.%d.%d.png", this.x, this.y))),
		}...)

		cmd := exec.Command("rendermmm", args...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			quiet := false

			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitcode := status.ExitStatus()
					if exitcode == 77 {
						// empty image
						quiet = true
					}
				}
			}

			if !quiet {
				log.Printf("Failed to run: %s\n", err.Error())
			}

			badFromWorker <- this
			continue
		}

		fromWorker <- this
	}
}

func main() {
	regionPath := flag.String("region", "hub/region/", "Path to region directory of world")
	outputPath := flag.String("o", "out", "Output directory")
	day := flag.Bool("day", true, "Daytime")
	samples := flag.Int("samples", 100, "Number of samples per pixel")
	cameraType := flag.String("camera", "iso", "Camera type (iso, topdown)")
	flag.Parse()

	commonArgs := []string{
		fmt.Sprintf("-day=%v", *day),
		fmt.Sprintf("-samples=%v", *samples),
		fmt.Sprintf("-camera=%v", *cameraType),
		fmt.Sprintf("-region=%v", *regionPath),
		"-emptyexit=77",
		fmt.Sprintf("-w=%v", CHUNK_SIZE),
		fmt.Sprintf("-h=%v", CHUNK_SIZE),
	}

	toWorker := make(chan Pair)
	fromWorker := make(chan Pair)
	badFromWorker := make(chan Pair)

	go Worker(commonArgs, *outputPath, toWorker, fromWorker, badFromWorker)
	go Worker(commonArgs, *outputPath, toWorker, fromWorker, badFromWorker)
	go Worker(commonArgs, *outputPath, toWorker, fromWorker, badFromWorker)
	go Worker(commonArgs, *outputPath, toWorker, fromWorker, badFromWorker)

	good := make(map[Pair]bool)
	todo := []Pair{{0, 0}}
	done := make(map[Pair]bool)
	done[Pair{0, 0}] = true
	inProgress := 0

	expand := func(this Pair) {
		good[this] = true
		for dx := -1; dx <= 1; dx++ {
			pair := Pair{this.x + dx, this.y}
			exists := done[pair]
			if !exists {
				done[pair] = true
				todo = append(todo, pair)
			}
		}
		for dy := -1; dy <= 1; dy++ {
			pair := Pair{this.x, this.y + dy}
			exists := done[pair]
			if !exists {
				done[pair] = true
				todo = append(todo, pair)
			}
		}
	}

	for len(todo) > 0 || inProgress > 0 {
		if len(todo) > 0 {
			select {
			case toWorker <- todo[0]:
				todo = todo[1:]
				inProgress++

			case this := <-fromWorker:
				expand(this)
				inProgress--

			case <-badFromWorker:
				inProgress--
			}
		} else {
			select {
			case this := <-fromWorker:
				expand(this)
				inProgress--

			case <-badFromWorker:
				inProgress--
			}
		}
	}

	close(toWorker)

	log.Printf("Writing html\n")
	err := writeHtml(good, filepath.Join(*outputPath, "map.html"))
	if err != nil {
		fmt.Println(err)
	}
}

func writeHtml(good map[Pair]bool, path string) error {
	minX := 0
	minY := 0
	maxX := 0
	maxY := 0
	for k := range good {
		if minX > k.x {
			minX = k.x
		}
		if minY > k.y {
			minY = k.y
		}
		if maxX < k.x {
			maxX = k.x
		}
		if maxY < k.y {
			maxY = k.y
		}
	}

	fh, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	fh.Write([]byte("<!DOCTYPE html>\n<html><head><style>table { border-collapse: collapse; } td, tr, body, html, img { padding: 0; margin: 0; } img { display: block; }</style></head><body><table>"))

	for y := minY; y <= maxY; y++ {
		fh.Write([]byte("<tr>"))
		for x := minX; x <= maxX; x++ {
			fh.Write([]byte("<td>"))

			if good[Pair{x, y}] {
				fh.Write([]byte(fmt.Sprintf("<img src=\"map.%d.%d.png\">", x, y)))
			}

			fh.Write([]byte("</td>"))
		}
		fh.Write([]byte("</tr>"))
	}

	fh.Write([]byte("</table></body></html>"))

	return nil
}
