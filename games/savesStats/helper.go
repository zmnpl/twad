package savesStats

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"math/bits"
	"os"
	"regexp"
)

// BinaryStartPosition returns the position after the search series has been found
func binaryStartPosition(binaryData []byte, startAfterSeries []byte, startAt int) int {
	seriesLength := len(startAfterSeries)
	readFrom := -1
	for index, _ := range binaryData {
		if index >= startAt {
			if bytes.Equal(startAfterSeries, binaryData[index:index+seriesLength]) {
				readFrom = index + seriesLength
				break
			}
		}
	}
	return readFrom
}

// readSize gets the size of coming string in variable length encoding
// Here: looks like reversed order; lowest bits are first bytes
//
// https://en.wikipedia.org/wiki/Variable-length_quantity
func readSize(reader io.ReadSeeker) int {
	b := make([]byte, 1)
	count := 0
	ofset := 0

	for {
		reader.Read(b)
		count = count | (int(b[0])&0x7f)<<ofset
		ofset = 7

		// Checks if the MSB is 0
		if (int(b[0]) & 0x80) == 0 {
			break
		}
	}

	return count
}

func reverseBitsIfNeeded(i uint32) uint32 {
	if i > 0x0000FFFF {
		i = bits.Reverse32(i)
	}
	return i
}

func getFileContentFromZip(src string, fileName string) ([]byte, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == fileName {
			content := make([]byte, 0, f.FileInfo().Size())
			contentBuffer := bytes.NewBuffer(content)

			rc, err := f.Open()
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(contentBuffer, rc)

			rc.Close()

			if err != nil {
				return nil, err
			}
			return contentBuffer.Bytes(), nil
		}
	}
	return nil, err
}

func fileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}
