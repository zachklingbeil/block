package value

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"sort"
// 	"strings"
// )

// type Input struct {
// 	Name      string `json:"name"`
// 	LabelName string `json:"labelName"`
// 	Owner     struct {
// 		ID string `json:"id"`
// 	} `json:"owner"`
// }

// type Output struct {
// 	LoopringEns string `json:"loopringEns"`
// 	Address     string `json:"address"`
// }

// func main() {
// 	files, err := filepath.Glob("subdomains_page_*.json")
// 	if err != nil {
// 		log.Fatalf("Error finding files: %v", err)
// 	}
// 	if len(files) == 0 {
// 		log.Fatalf("No subdomains_page_*.json files found")
// 	}

// 	// Sort files by page number
// 	re := regexp.MustCompile(`subdomains_page_(\d+)\.json`)
// 	sort.Slice(files, func(i, j int) bool {
// 		mi := re.FindStringSubmatch(files[i])
// 		mj := re.FindStringSubmatch(files[j])
// 		if len(mi) < 2 || len(mj) < 2 {
// 			return files[i] < files[j]
// 		}
// 		return mi[1] < mj[1]
// 	})

// 	var all []Output

// 	for _, file := range files {
// 		data, err := os.ReadFile(file)
// 		if err != nil {
// 			log.Fatalf("Error reading %s: %v", file, err)
// 		}
// 		var inputs []Input
// 		if err := json.Unmarshal(data, &inputs); err != nil {
// 			log.Fatalf("Error unmarshaling %s: %v", file, err)
// 		}
// 		for _, in := range inputs {
// 			all = append(all, Output{
// 				LoopringEns: strings.ToLower(in.Name),
// 				Address:     strings.ToLower(in.Owner.ID),
// 			})
// 		}
// 	}

// 	outPath := "dotloop.json"
// 	f, err := os.Create(outPath)
// 	if err != nil {
// 		log.Fatalf("Error creating file: %v", err)
// 	}
// 	defer f.Close()

// 	enc := json.NewEncoder(f)
// 	enc.SetIndent("", "  ")
// 	if err := enc.Encode(all); err != nil {
// 		log.Fatalf("Error encoding JSON: %v", err)
// 	}

// 	fmt.Printf("Wrote %d entries to %s\n", len(all), outPath)
// }

// package value

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// )

// const ensSubgraphURL = "https://api.thegraph.com/subgraphs/name/ensdomains/ens"

// type ensDomain struct {
// 	Name      string `json:"name"`
// 	LabelName string `json:"labelName"`
// 	Owner     struct {
// 		ID string `json:"id"`
// 	} `json:"owner"`
// 	CreatedAt string `json:"createdAt,omitempty"` // Optional, if you want to keep createdAt
// }

// type ensDomainsResponse struct {
// 	Data struct {
// 		Domains []ensDomain `json:"domains"`
// 	} `json:"data"`
// }

// // FetchSubdomainsPage fetches a single page of subdomains for the given parent.
// func FetchSubdomainsPage(parent string, pageSize, skip int) ([]ensDomain, error) {
// 	query := fmt.Sprintf(`
//         {
//             domains(
//                 where: { name_ends_with: ".%s" }
//                 first: %d
//                 skip: %d
//                 orderBy: createdAt
//                 orderDirection: asc
//             ) {
//                 name
//                 labelName
//                 owner {
//                   id
//                 }
//             }
//         }`, parent, pageSize, skip)

// 	reqBody := map[string]string{"query": query}
// 	bodyBytes, _ := json.Marshal(reqBody)
// 	resp, err := http.Post(ensSubgraphURL, "application/json", bytes.NewBuffer(bodyBytes))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	respBytes, _ := io.ReadAll(resp.Body)

// 	var result ensDomainsResponse
// 	if err := json.Unmarshal(respBytes, &result); err != nil {
// 		return nil, err
// 	}
// 	// Filter out domains with empty labelName (hashed names)
// 	readable := make([]ensDomain, 0, len(result.Data.Domains))
// 	for _, d := range result.Data.Domains {
// 		if d.LabelName != "" {
// 			readable = append(readable, d)
// 		}
// 	}

// 	return readable, nil
// }

// // ...existing code...
