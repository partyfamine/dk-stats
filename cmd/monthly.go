package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
)

type Line struct {
	Reporting         string  `csv:"Reporting"`
	Date              string  `csv:"Date"`
	SaleMonth         string  `csv:"Sale Month"`
	Store             string  `csv:"Store"`
	Artist            string  `csv:"Artist"`
	Title             string  `csv:"Title"`
	ISRC              string  `csv:"ISRC"`
	UPC               string  `csv:"UPC"`
	Quantity          int     `csv:"Quantity"`
	TeamPct           string  `csv:"Team Percentage"`
	SongAndAlbum      string  `csv:"Song/Album"`
	CustomerPaid      string  `csv:"Customer Paid"`
	CustomerCurrency  string  `csv:"Customer Currency"`
	CountryOfSale     string  `csv:"Country of Sale"`
	Songwriter        string  `csv:"Songwriter"`
	RoyaltiesWithheld string  `csv:"Royalties Withheld"`
	Earnings          float64 `csv:"Earnings (USD)"`
}

type Plays struct {
	Plays    int
	Earnings float64
}

var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "monthly stats",
	Long:  "monthly stats",
	Run:   monthlyStats,
}

func monthlyStats(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		panic("must provide single argument for tsv file path")
	}

	file, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'

	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("headers: %v", dec.Header())

	var lines []Line
	for {
		var line Line
		if err := dec.Decode(&line); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, line)
	}

	fmt.Printf("lines: %d\n", len(lines))

	artistPlays := make(map[string]map[string]Plays)
	for _, line := range lines {
		if _, ok := artistPlays[line.Artist]; !ok {
			artistPlays[line.Artist] = make(map[string]Plays)
		}
		if _, ok := artistPlays[line.Artist]; !ok {
			artistPlays[line.Artist] = make(map[string]Plays)
		}

		plays := artistPlays[line.Artist][line.SaleMonth]
		artistPlays[line.Artist][line.SaleMonth] = Plays{Plays: plays.Plays + line.Quantity, Earnings: plays.Earnings + line.Earnings}
	}

	for artist, monthlyPlays := range artistPlays {
		fmt.Println("artist: " + artist)
		totalPlays := 0
		totalEarnings := 0.0
		for month, plays := range monthlyPlays {
			fmt.Printf("monthly plays for %s: %d, total earnings: %.2f\n", month, plays.Plays, plays.Earnings)
			totalPlays += plays.Plays
			totalEarnings += plays.Earnings
		}

		fmt.Printf("total plays: %d, total earnings: %.2f\n", totalPlays, totalEarnings)
	}
}
