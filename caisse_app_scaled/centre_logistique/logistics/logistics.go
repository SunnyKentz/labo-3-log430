package logistics

import (
	"caisse-app-scaled/caisse_app_scaled/auth"
	"caisse-app-scaled/caisse_app_scaled/centre_logistique/db"
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/models"
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"errors"
	"fmt"
	"net/http"
)

type Commande struct {
	Id        int    `json:"id"`
	Magasin   string `json:"Magsin"`
	ProduitID int    `json:"ProduitID"`
	Message   string `json:"message"`
	Host      string `json:"host"`
}

var commandes []Commande = []Commande{}
var nom = ""

func Nom() (string, error) {
	if nom == "" {
		return "", errors.New("no name")
	}
	return nom, nil
}

func Login(n string, pw string) (string, error) {
	if auth.IsUsernameValid(nom) && auth.IsUserPasswordValid(nom, pw) {
		return auth.CreateJWT(nom)
	}
	return "", errors.New("failed to login")
}

func CheckLogedIn(jwt string) error {
	_, err := auth.ValidateJWT(jwt)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCommands() []Commande {
	return commandes
}
func AjouterUneCommande(produitID int, magasin string, host string) {
	newId := len(commandes) + 1
	prod, _ := db.GetProduitParID(produitID)
	message := magasin + " demande une réaprovisionement de 10 " + prod.Nom
	logger.Info(message)
	commandes = append(commandes, Commande{
		Id:        newId,
		ProduitID: produitID,
		Message:   message,
		Host:      host,
	})
}

func AccepterUneCommande(id int) bool {
	for i, cmd := range commandes {
		if cmd.Id == id {
			err := db.MettreAJourQuantite(cmd.ProduitID, -10)
			Errnotnil(err)
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/api/v1/produit/%d/%d", cmd.Host, cmd.ProduitID, 10), nil)
			if err != nil {
				logger.Error("Erreur lors de la création de la requête: " + err.Error())
				return false
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Error("Erreur lors de l'envoi de reaprovisionement : " + err.Error())
				return false
			}
			defer resp.Body.Close()
			// Remove the command from the slice
			commandes = append(commandes[:i], commandes[i+1:]...)
			return true
		}
	}
	return false
}

func RefuserUneCommande(id int) bool {
	for i, cmd := range commandes {
		if cmd.Id == id {
			// Remove the command from the slice
			commandes = append(commandes[:i], commandes[i+1:]...)
			return true
		}
	}
	return false
}

func TrouverProduit(nomPartiel string) ([]models.Produit, error) {
	return db.GetProduitsParNomWildcard(nomPartiel)
}

func TrouverProduitParID(id int) (*models.Produit, error) {
	return db.GetProduitParID(id)
}

func MiseAJourProduit(produitID int, nom string, prix float64, description string) error {
	produit, err := db.GetProduitParID(produitID)
	if err != nil {
		return errors.New("produit not found")
	}
	produit.Nom = nom
	produit.Prix = prix
	produit.Description = description
	if err := db.MettreAJour(*produit); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
