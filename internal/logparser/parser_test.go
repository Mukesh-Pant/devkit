package logparser

import (
	"strings"
	"testing"
)

// Sample log data for testing
const sampleCombinedLog = `192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
192.168.1.2 - - [01/Jan/2024:12:00:01 +0000] "POST /api/login HTTP/1.1" 401 567 "-" "curl/7.68.0"
192.168.1.1 - - [01/Jan/2024:12:00:02 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
192.168.1.3 - - [01/Jan/2024:12:00:03 +0000] "GET /api/products HTTP/1.1" 200 5678 "-" "Mozilla/5.0"
192.168.1.1 - - [01/Jan/2024:12:00:04 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
192.168.1.4 - - [01/Jan/2024:12:00:05 +0000] "GET /api/orders HTTP/1.1" 404 0 "-" "curl/7.68.0"
192.168.1.2 - - [01/Jan/2024:12:00:06 +0000] "POST /api/login HTTP/1.1" 200 890 "-" "curl/7.68.0"
192.168.1.5 - - [01/Jan/2024:12:00:07 +0000] "GET /api/users HTTP/1.1" 500 0 "-" "Mozilla/5.0"
192.168.1.1 - - [01/Jan/2024:12:00:08 +0000] "GET /api/products HTTP/1.1" 200 5678 "-" "Mozilla/5.0"
192.168.1.3 - - [01/Jan/2024:12:00:09 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`

const sampleJSONLog = `{"ip":"192.168.1.1","timestamp":"2024-01-01T12:00:00Z","method":"GET","path":"/api/users","status_code":200,"size":1234,"user_agent":"Mozilla/5.0","referrer":"-"}
{"ip":"192.168.1.2","timestamp":"2024-01-01T12:00:01Z","method":"POST","path":"/api/login","status_code":401,"size":567,"user_agent":"curl/7.68.0","referrer":"-"}
{"ip":"192.168.1.1","timestamp":"2024-01-01T12:00:02Z","method":"GET","path":"/api/users","status_code":200,"size":1234,"user_agent":"Mozilla/5.0","referrer":"-"}
{"ip":"192.168.1.3","timestamp":"2024-01-01T12:00:03Z","method":"GET","path":"/api/products","status_code":200,"size":5678,"user_agent":"Mozilla/5.0","referrer":"-"}
{"ip":"192.168.1.1","timestamp":"2024-01-01T12:00:04Z","method":"GET","path":"/api/users","status_code":200,"size":1234,"user_agent":"Mozilla/5.0","referrer":"-"}`

const samplePlainLog = `2024-01-01 12:00:00 INFO 192.168.1.1 GET /api/users 200
2024-01-01 12:00:01 ERROR 192.168.1.2 POST /api/login 401
2024-01-01 12:00:02 INFO 192.168.1.1 GET /api/users 200
2024-01-01 12:00:03 INFO 192.168.1.3 GET /api/products 200
2024-01-01 12:00:04 ERROR 192.168.1.4 GET /api/orders 404`

func TestParser_ParseReader_Combined(t *testing.T) {
	parser := NewParser(Options{
		Format: FormatCombined,
		TopN:   10,
	})

	reader := strings.NewReader(sampleCombinedLog)
	stats, err := parser.ParseReader(reader)

	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	// Verify total lines
	if stats.TotalLines != 10 {
		t.Errorf("TotalLines = %d, want 10", stats.TotalLines)
	}

	// Verify status code counts
	if stats.StatusCounts[200] != 7 {
		t.Errorf("StatusCounts[200] = %d, want 7", stats.StatusCounts[200])
	}
	if stats.StatusCounts[401] != 1 {
		t.Errorf("StatusCounts[401] = %d, want 1", stats.StatusCounts[401])
	}
	if stats.StatusCounts[404] != 1 {
		t.Errorf("StatusCounts[404] = %d, want 1", stats.StatusCounts[404])
	}
	if stats.StatusCounts[500] != 1 {
		t.Errorf("StatusCounts[500] = %d, want 1", stats.StatusCounts[500])
	}

	// Verify error rate (3 errors out of 10 = 30%)
	expectedErrorRate := 30.0
	if stats.ErrorRate != expectedErrorRate {
		t.Errorf("ErrorRate = %.2f, want %.2f", stats.ErrorRate, expectedErrorRate)
	}

	// Verify top IP
	if len(stats.TopIPs) == 0 {
		t.Fatal("TopIPs is empty")
	}
	if stats.TopIPs[0].Value != "192.168.1.1" {
		t.Errorf("TopIPs[0].Value = %s, want 192.168.1.1", stats.TopIPs[0].Value)
	}
	if stats.TopIPs[0].Count != 4 {
		t.Errorf("TopIPs[0].Count = %d, want 4", stats.TopIPs[0].Count)
	}

	// Verify top path
	if len(stats.TopPaths) == 0 {
		t.Fatal("TopPaths is empty")
	}
	if stats.TopPaths[0].Value != "/api/users" {
		t.Errorf("TopPaths[0].Value = %s, want /api/users", stats.TopPaths[0].Value)
	}
	if stats.TopPaths[0].Count != 5 {
		t.Errorf("TopPaths[0].Count = %d, want 5", stats.TopPaths[0].Count)
	}

	// Verify top status code
	if len(stats.TopStatusCodes) == 0 {
		t.Fatal("TopStatusCodes is empty")
	}
	if stats.TopStatusCodes[0].Value != "200" {
		t.Errorf("TopStatusCodes[0].Value = %s, want 200", stats.TopStatusCodes[0].Value)
	}
	if stats.TopStatusCodes[0].Count != 7 {
		t.Errorf("TopStatusCodes[0].Count = %d, want 7", stats.TopStatusCodes[0].Count)
	}
}

func TestParser_ParseReader_JSON(t *testing.T) {
	parser := NewParser(Options{
		Format: FormatJSON,
		TopN:   10,
	})

	reader := strings.NewReader(sampleJSONLog)
	stats, err := parser.ParseReader(reader)

	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	// Verify total lines
	if stats.TotalLines != 5 {
		t.Errorf("TotalLines = %d, want 5", stats.TotalLines)
	}

	// Verify status code counts
	if stats.StatusCounts[200] != 4 {
		t.Errorf("StatusCounts[200] = %d, want 4", stats.StatusCounts[200])
	}
	if stats.StatusCounts[401] != 1 {
		t.Errorf("StatusCounts[401] = %d, want 1", stats.StatusCounts[401])
	}

	// Verify error rate (1 error out of 5 = 20%)
	expectedErrorRate := 20.0
	if stats.ErrorRate != expectedErrorRate {
		t.Errorf("ErrorRate = %.2f, want %.2f", stats.ErrorRate, expectedErrorRate)
	}

	// Verify top IP
	if len(stats.TopIPs) == 0 {
		t.Fatal("TopIPs is empty")
	}
	if stats.TopIPs[0].Value != "192.168.1.1" {
		t.Errorf("TopIPs[0].Value = %s, want 192.168.1.1", stats.TopIPs[0].Value)
	}
	if stats.TopIPs[0].Count != 3 {
		t.Errorf("TopIPs[0].Count = %d, want 3", stats.TopIPs[0].Count)
	}
}

func TestParser_ParseReader_Plain(t *testing.T) {
	parser := NewParser(Options{
		Format: FormatPlain,
		TopN:   10,
	})

	reader := strings.NewReader(samplePlainLog)
	stats, err := parser.ParseReader(reader)

	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	// Verify total lines
	if stats.TotalLines != 5 {
		t.Errorf("TotalLines = %d, want 5", stats.TotalLines)
	}

	// Verify status code counts
	if stats.StatusCounts[200] != 3 {
		t.Errorf("StatusCounts[200] = %d, want 3", stats.StatusCounts[200])
	}
	if stats.StatusCounts[401] != 1 {
		t.Errorf("StatusCounts[401] = %d, want 1", stats.StatusCounts[401])
	}
	if stats.StatusCounts[404] != 1 {
		t.Errorf("StatusCounts[404] = %d, want 1", stats.StatusCounts[404])
	}

	// Verify error rate (2 errors out of 5 = 40%)
	expectedErrorRate := 40.0
	if stats.ErrorRate != expectedErrorRate {
		t.Errorf("ErrorRate = %.2f, want %.2f", stats.ErrorRate, expectedErrorRate)
	}

	// Verify IP extraction
	if len(stats.TopIPs) == 0 {
		t.Fatal("TopIPs is empty")
	}
	if stats.TopIPs[0].Value != "192.168.1.1" {
		t.Errorf("TopIPs[0].Value = %s, want 192.168.1.1", stats.TopIPs[0].Value)
	}
	if stats.TopIPs[0].Count != 2 {
		t.Errorf("TopIPs[0].Count = %d, want 2", stats.TopIPs[0].Count)
	}
}

func TestParser_ParseReader_TopN(t *testing.T) {
	parser := NewParser(Options{
		Format: FormatCombined,
		TopN:   3,
	})

	reader := strings.NewReader(sampleCombinedLog)
	stats, err := parser.ParseReader(reader)

	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	// Verify top N limit is applied
	if len(stats.TopIPs) > 3 {
		t.Errorf("len(TopIPs) = %d, want <= 3", len(stats.TopIPs))
	}
	if len(stats.TopPaths) > 3 {
		t.Errorf("len(TopPaths) = %d, want <= 3", len(stats.TopPaths))
	}
	if len(stats.TopStatusCodes) > 3 {
		t.Errorf("len(TopStatusCodes) = %d, want <= 3", len(stats.TopStatusCodes))
	}
}

func TestParser_ParseReader_EmptyFile(t *testing.T) {
	parser := NewParser(Options{
		Format: FormatCombined,
		TopN:   10,
	})

	reader := strings.NewReader("")
	stats, err := parser.ParseReader(reader)

	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	// Verify empty file handling
	if stats.TotalLines != 0 {
		t.Errorf("TotalLines = %d, want 0", stats.TotalLines)
	}
	if stats.ErrorRate != 0 {
		t.Errorf("ErrorRate = %.2f, want 0.00", stats.ErrorRate)
	}
}

func TestParser_ParseReader_MalformedLines(t *testing.T) {
	parser := NewParser(Options{
		Format: FormatCombined,
		TopN:   10,
	})

	// Mix of valid and malformed lines
	malformedLog := `192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
this is not a valid log line
192.168.1.2 - - [01/Jan/2024:12:00:01 +0000] "POST /api/login HTTP/1.1" 401 567 "-" "curl/7.68.0"
another invalid line
192.168.1.1 - - [01/Jan/2024:12:00:02 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`

	reader := strings.NewReader(malformedLog)
	stats, err := parser.ParseReader(reader)

	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}

	// Verify total lines includes malformed lines
	if stats.TotalLines != 5 {
		t.Errorf("TotalLines = %d, want 5", stats.TotalLines)
	}

	// Verify only valid lines are counted in status codes
	totalStatusCounts := 0
	for _, count := range stats.StatusCounts {
		totalStatusCounts += count
	}
	if totalStatusCounts != 3 {
		t.Errorf("Total status counts = %d, want 3", totalStatusCounts)
	}
}

func TestParseCombined_ValidLine(t *testing.T) {
	line := `192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`

	entry, err := parseCombined(line)

	if err != nil {
		t.Fatalf("parseCombined() error = %v", err)
	}

	if entry.IP != "192.168.1.1" {
		t.Errorf("IP = %s, want 192.168.1.1", entry.IP)
	}
	if entry.Method != "GET" {
		t.Errorf("Method = %s, want GET", entry.Method)
	}
	if entry.Path != "/api/users" {
		t.Errorf("Path = %s, want /api/users", entry.Path)
	}
	if entry.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", entry.StatusCode)
	}
	if entry.Size != 1234 {
		t.Errorf("Size = %d, want 1234", entry.Size)
	}
}

func TestParseCombined_InvalidLine(t *testing.T) {
	line := `this is not a valid log line`

	_, err := parseCombined(line)

	if err == nil {
		t.Error("parseCombined() should return error for invalid line")
	}
}

func TestParseJSON_ValidLine(t *testing.T) {
	line := `{"ip":"192.168.1.1","timestamp":"2024-01-01T12:00:00Z","method":"GET","path":"/api/users","status_code":200,"size":1234,"user_agent":"Mozilla/5.0","referrer":"-"}`

	entry, err := parseJSON(line)

	if err != nil {
		t.Fatalf("parseJSON() error = %v", err)
	}

	if entry.IP != "192.168.1.1" {
		t.Errorf("IP = %s, want 192.168.1.1", entry.IP)
	}
	if entry.Method != "GET" {
		t.Errorf("Method = %s, want GET", entry.Method)
	}
	if entry.Path != "/api/users" {
		t.Errorf("Path = %s, want /api/users", entry.Path)
	}
	if entry.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", entry.StatusCode)
	}
}

func TestParseJSON_InvalidLine(t *testing.T) {
	line := `this is not valid JSON`

	_, err := parseJSON(line)

	if err == nil {
		t.Error("parseJSON() should return error for invalid JSON")
	}
}

func TestParsePlain_ValidLine(t *testing.T) {
	line := `2024-01-01 12:00:00 INFO 192.168.1.1 GET /api/users 200`

	entry, err := parsePlain(line)

	if err != nil {
		t.Fatalf("parsePlain() error = %v", err)
	}

	if entry.IP != "192.168.1.1" {
		t.Errorf("IP = %s, want 192.168.1.1", entry.IP)
	}
	if entry.Method != "GET" {
		t.Errorf("Method = %s, want GET", entry.Method)
	}
	if entry.Path != "/api/users" {
		t.Errorf("Path = %s, want /api/users", entry.Path)
	}
	if entry.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", entry.StatusCode)
	}
}
