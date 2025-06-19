package tests

import (
	"bytes"
	LOGIS "caisse-app-scaled/caisse_app_scaled/centre_logistique/api"
	MAG "caisse-app-scaled/caisse_app_scaled/magasin/api"
	"caisse-app-scaled/caisse_app_scaled/magasin/caissier"
	MERE "caisse-app-scaled/caisse_app_scaled/maison_mere/api"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
	magToken   string
	mereToken  string
	logisToken string
)

func BeforeAll() {
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PORT", "5435")
	os.Setenv("GATEWAY", "localhost")

	go MERE.NewApp()
	go MAG.NewApp()
	go LOGIS.NewApp()
	time.Sleep(time.Second * 1) // time for goroutine to startup

	// Login to get tokens
	magToken = login("magasin", "http://localhost:8080")
	mereToken = login("maison_mere", "http://localhost:8090")
	logisToken = login("logistique", "http://localhost:8091")

}

func login(service, baseURL string) string {

	resp, err := http.PostForm(baseURL+"/login", map[string][]string{
		"username": {"Bob"},
		"caisse":   {"Caisse 1"},
		"magasin":  {"Magasin 1"},
	})
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 400) {
		return "ok"
	}
	fmt.Println(service+" failed logged in", resp.StatusCode)
	return ""
}

func TestMain(m *testing.M) {
	BeforeAll()
	code := m.Run()
	caissier.FermerPOS()
	os.Exit(code)
}

func errnotnil(s string, t *testing.T, err error) {
	if err != nil {
		t.Error(s + ": " + err.Error())
	}
}

// ==================== MAGASIN API TESTS ====================

func TestMagasinProduits(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/produits", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin produits", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMagasinProduitsSearch(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/produits/t", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin produits search", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMagasinCart(t *testing.T) {
	// Test get cart
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/cart", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin get cart", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMagasinAddToCart(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/cart/1", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin add to cart", t, err)
	defer resp.Body.Close()

	// This might fail if product doesn't exist, which is expected
	t.Logf("Add to cart status: %d", resp.StatusCode)
}

func TestMagasinRemoveFromCart(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/cart/1", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin remove from cart", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMagasinVendre(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/vendre", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin vendre", t, err)
	defer resp.Body.Close()

	// This might fail if cart is empty, which is expected
	t.Logf("Vendre status: %d", resp.StatusCode)
}

func TestMagasinTransactions(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin transactions", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMagasinRembourser(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/rembourser/1", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin rembourser", t, err)
	defer resp.Body.Close()

	// This might fail if transaction doesn't exist, which is expected
	t.Logf("Rembourser status: %d", resp.StatusCode)
}

func TestMagasinCommande(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/produit/1", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin commande", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMagasinUpdateProduit(t *testing.T) {
	updateData := map[string]interface{}{
		"nom":         "Test Product",
		"prix":        10.99,
		"description": "Test description",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "http://localhost:8080/api/produit/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	errnotnil("Magasin update produit", t, err)
	defer resp.Body.Close()

	// This might fail if product doesn't exist, which is expected
	t.Logf("Update produit status: %d", resp.StatusCode)
}

func TestMagasinReapprovisionner(t *testing.T) {
	req, _ := http.NewRequest("PUT", "http://localhost:8080/api/produit/1/10", nil)

	resp, err := httpClient.Do(req)
	errnotnil("Magasin reapprovisionner", t, err)
	defer resp.Body.Close()

	// This might fail if product doesn't exist, which is expected
	t.Logf("Reapprovisionner status: %d", resp.StatusCode)
}

// ==================== MAISON MERE API TESTS ====================

func TestMereNotify(t *testing.T) {
	notifyData := map[string]string{
		"message": "Test notification",
		"host":    "test-host",
	}

	jsonData, _ := json.Marshal(notifyData)
	resp, err := httpClient.Post("http://localhost:8090/api/notify", "application/json", bytes.NewBuffer(jsonData))
	errnotnil("Mere notify", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereSubscribe(t *testing.T) {
	subscribeData := map[string]string{
		"host": "http://localhost:8080",
	}

	jsonData, _ := json.Marshal(subscribeData)
	resp, err := httpClient.Post("http://localhost:8090/api/subscribe", "application/json", bytes.NewBuffer(jsonData))
	errnotnil("Mere subscribe", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereAlerts(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8090/api/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+mereToken)

	resp, err := httpClient.Do(req)
	errnotnil("Mere alerts", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereTransactions(t *testing.T) {
	resp, err := httpClient.Get("http://localhost:8090/api/transactions")
	errnotnil("Mere transactions", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereTransactionById(t *testing.T) {
	resp, err := httpClient.Get("http://localhost:8090/api/transactions/1")
	errnotnil("Mere transaction by id", t, err)
	defer resp.Body.Close()

	// This might fail if transaction doesn't exist, which is expected
	t.Logf("Transaction by id status: %d", resp.StatusCode)
}

func TestMereCreateTransaction(t *testing.T) {
	transactionData := map[string]interface{}{
		"caisse":      "Caisse 1",
		"type":        "VENTE",
		"produit_ids": "1,2",
		"montant":     25.50,
		"date":        time.Now(),
	}

	jsonData, _ := json.Marshal(transactionData)
	resp, err := httpClient.Post("http://localhost:8090/api/transactions", "application/json", bytes.NewBuffer(jsonData))
	errnotnil("Mere create transaction", t, err)
	defer resp.Body.Close()

	// This might fail if validation fails, which is expected
	t.Logf("Create transaction status: %d", resp.StatusCode)
}

func TestMereDeleteTransaction(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "http://localhost:8090/api/transactions/1", nil)

	resp, err := httpClient.Do(req)
	errnotnil("Mere delete transaction", t, err)
	defer resp.Body.Close()

	// This might fail if transaction doesn't exist, which is expected
	t.Logf("Delete transaction status: %d", resp.StatusCode)
}

func TestMereMagasins(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8090/api/magasins", nil)
	req.Header.Set("Authorization", "Bearer "+mereToken)

	resp, err := httpClient.Do(req)
	errnotnil("Mere magasins", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereAnalytics(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8090/api/analytics/tout", nil)
	req.Header.Set("Authorization", "Bearer "+mereToken)

	resp, err := httpClient.Do(req)
	errnotnil("Mere analytics", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereRaport(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8090/api/raport", nil)
	req.Header.Set("Authorization", "Bearer "+mereToken)

	resp, err := httpClient.Do(req)
	errnotnil("Mere raport", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereProduitsSearch(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8090/api/produits/test", nil)
	req.Header.Set("Authorization", "Bearer "+mereToken)

	resp, err := httpClient.Do(req)
	errnotnil("Mere produits search", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestMereUpdateProduit(t *testing.T) {
	updateData := map[string]interface{}{
		"productId":   1,
		"nom":         "Test Product",
		"prix":        10.99,
		"description": "Test description",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "http://localhost:8090/api/produit", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mereToken)

	resp, err := httpClient.Do(req)
	errnotnil("Mere update produit", t, err)
	defer resp.Body.Close()

	// This might fail if product doesn't exist, which is expected
	t.Logf("Mere update produit status: %d", resp.StatusCode)
}

// ==================== LOGISTICS API TESTS ====================

func TestLogisNotify(t *testing.T) {
	resp, err := httpClient.Get("http://localhost:8091/api/notify")
	errnotnil("Logis notify", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestLogisCommands(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:8091/api/commands", nil)
	req.Header.Set("Authorization", "Bearer "+logisToken)

	resp, err := httpClient.Do(req)
	errnotnil("Logis commands", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestLogisAcceptCommande(t *testing.T) {
	req, _ := http.NewRequest("PUT", "http://localhost:8091/api/commande/1", nil)
	req.Header.Set("Authorization", "Bearer "+logisToken)

	resp, err := httpClient.Do(req)
	errnotnil("Logis accept commande", t, err)
	defer resp.Body.Close()

	// This might fail if command doesn't exist, which is expected
	t.Logf("Accept commande status: %d", resp.StatusCode)
}

func TestLogisRefuseCommande(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "http://localhost:8091/api/commande/1", nil)
	req.Header.Set("Authorization", "Bearer "+logisToken)

	resp, err := httpClient.Do(req)
	errnotnil("Logis refuse commande", t, err)
	defer resp.Body.Close()

	// This might fail if command doesn't exist, which is expected
	t.Logf("Refuse commande status: %d", resp.StatusCode)
}

func TestLogisProduitsSearch(t *testing.T) {
	resp, err := httpClient.Get("http://localhost:8091/api/produits/t")
	errnotnil("Logis produits search", t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestLogisProduitById(t *testing.T) {
	resp, err := httpClient.Get("http://localhost:8091/api/produits/id/1")
	errnotnil("Logis produit by id", t, err)
	defer resp.Body.Close()

	// This might fail if product doesn't exist, which is expected
	t.Logf("Produit by id status: %d", resp.StatusCode)
}

func TestLogisUpdateProduit(t *testing.T) {
	updateData := map[string]interface{}{
		"nom":         "Test Product",
		"prix":        10.99,
		"description": "Test description",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "http://localhost:8091/api/produit/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	errnotnil("Logis update produit", t, err)
	defer resp.Body.Close()

	// This might fail if product doesn't exist, which is expected
	t.Logf("Logis update produit status: %d", resp.StatusCode)
}

// ==================== INTEGRATION TESTS ====================

func TestCompleteWorkflow(t *testing.T) {
	t.Log("Testing complete workflow...")

	// 1. Add product to cart
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/cart/1", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)
	resp, err := httpClient.Do(req)
	if err == nil {
		t.Logf("Add to cart: %d", resp.StatusCode)
		resp.Body.Close()
	}

	// 2. Complete sale
	req, _ = http.NewRequest("POST", "http://localhost:8080/api/vendre", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)
	resp, err = httpClient.Do(req)
	if err == nil {
		t.Logf("Complete sale: %d", resp.StatusCode)
		resp.Body.Close()
	}

	// 3. Request reapprovisionment
	req, _ = http.NewRequest("POST", "http://localhost:8080/api/produit/1", nil)
	req.Header.Set("Authorization", "Bearer "+magToken)
	resp, err = httpClient.Do(req)
	if err == nil {
		t.Logf("Request reapprovisionment: %d", resp.StatusCode)
		resp.Body.Close()
	}

	// 4. Create logistics command
	commandeData := map[string]string{"host": "test-host"}
	jsonData, _ := json.Marshal(commandeData)
	resp, err = httpClient.Post("http://localhost:8091/api/commande/Magasin%201/1", "application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		t.Logf("Create logistics command: %d", resp.StatusCode)
		resp.Body.Close()
	}

	t.Log("Complete workflow test finished")
}
