package main


import "fmt" 
import "log" 
import "net/http" 
import "io/ioutil"
import "github.com/tidwall/gjson"
import "strconv"
import "os"


func main() {

	if len(os.Args) == 1 {
		showUsage(os.Args[0])
		return
	}

	site := "http://localhost:9000"
	command := `get-coverage`
	expected_coverage := 100

	if len(os.Args) > 1 {
		command, site, expected_coverage = processArgs(os.Args[1:])
	}
	
	dispatch(command, site, expected_coverage)

	
}

func processArgs(args []string) (command string, site string, expected_coverage int) {
	command = args[0]
	
	i := 1

	for i < len(args) {
		if args[i] == "--site" {
			site = args[i+1]
			i++
		} else if args[i] == "--expected-coverage" {
			expected_coverage, _ = strconv.Atoi(args[i+1])
			i++
		}
		i++
	}
	return command, site, expected_coverage
}

func showUsage(pgmName string) {
	fmt.Println(fmt.Sprintf("%s command [options]", pgmName))
	fmt.Println("\ncommands:") 
	fmt.Println(`get-coverage: retrieve code coverage in percentage`)
	fmt.Println("\noptions:")
	fmt.Println("--site sitename: i.e. --site http://localhost:9000")
	fmt.Println("--expected-coverage value: i.e. --expected-coverage 80")
}

func dispatch(command string, site string, expected_coverage int) {
	if "get-coverage" == command {
		coverage := getCoverage(site)
		fmt.Printf("%d\n", coverage)
	} else if "get-complexity" == command {
		complexity := getComplexity(site)
		fmt.Printf("%d\n", complexity)
	} else if "check-coverage" == command {
		coverage := getCoverage(site)
		if coverage >= expected_coverage {
			fmt.Println(fmt.Sprintf("passed: expected %d and got %d", expected_coverage, coverage))
		} else {
			fmt.Println(fmt.Sprintf("failed: expected at least %d but got %d", expected_coverage, coverage))
		}

	}
}

func getCoverage(site string) int {
	data := callSonarQubeServer(site)

	coverageStr := gjson.Get(string(data),`component.measures.#[metric="line_coverage"].value`).String()
	f, _ := strconv.ParseFloat(coverageStr, 64)
	return int(f)
}

func getComplexity(site string) int {
	data := callSonarQubeServer(site)

	coverageStr := gjson.Get(string(data),`component.measures.#[metric="complexity"].value`).String()
	i, _ := strconv.Atoi(coverageStr)
	return i
}


func callSonarQubeServer(site string) string {
	path := "/api/measures/component"
	project := "org.springframework.samples:spring-petclinic"
	metrics := "ncloc,line_coverage,code_smells,complexity"
	additionalFields := "metrics,periods"
	
	url := fmt.Sprintf("%s%s?componentKey=%s&metricKeys=%s&additionalFields=%s", site, path, project, metrics, additionalFields)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err);
		return ""
	}

	client := &http.Client{}
	resp, err := client.Do(req);
	if err != nil {
		log.Fatal("Do: ", err)
		return ""
	}

	data, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println("data: " + string(data))

	return string(data)
}

