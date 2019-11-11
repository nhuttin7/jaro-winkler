// Research-paper: http://www.semantic-web-journal.net/system/files/swj1128.pdf
// The Jaro-Winkler similarity is an algorithm helping problem solvers calculate 2 strings and produce the percentage similarity
// Develop sets the threshold = 0.7 for the result and intergate it into searching engine. Ranking score for searching data.

package main

import ( //"encoding/csv"
	//"io"
	//"log"

	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	//"os"
	//"time"
)

func main() {

	//fmt.Println(JaroWinkle("the cat ate the mouse", "the mouse ate the cat food"))
	//fmt.Println(JaroWinkle("1", "12"))
	//fmt.Println(JaroWinkle("1", "17"))
	// fmt.Println(time.Now())
	//fmt.Println(JaroWinkle("martha", "marhta"))
	//fmt.Println(JaroWinkle("JONES", "JOHNSON"))
	// fmt.Println(time.Now())

	var arrT []float64
	var arrS []string
	// Open the file
	csvfile, err := os.Open("room_type.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	var arr1 []string
	var arr2 []string
	// Iterate through the records
	count := 0
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if count != 0 {
			arr1 = append(arr1, record[0])
			arr2 = append(arr2, record[1])

		}
		count = 1
	}

	t1 := time.Now()
	for i, _ := range arr1 {
		t2 := time.Now()
		JaroWinkler(arr1[i], arr2[i], 0.7)
		arrT = append(arrT, time.Now().Sub(t2).Seconds())
		arrS = append(arrS, arr1[i]+arr2[i])

	}

	fmt.Println("Total:", time.Now().Sub(t1).Seconds())
	var smallest float64 = arrT[0]
	var biggest float64
	for _, i := range arrT {
		if i <= smallest {
			smallest = i
		}
		if i > biggest {
			biggest = i
		}
	}

	fmt.Println("Smallest", smallest)
	fmt.Println("Biggest", biggest)
	fmt.Println("Average", (biggest-smallest)/2)

	var ssl string = arrS[0]
	var sbig string
	for _, i := range arrS {
		if len(ssl) >= len(i) {
			ssl = i
		}
		if len(sbig) < len(i) {
			sbig = i
		}

	}
	fmt.Println(len(ssl))
	fmt.Println(len(sbig))
}
