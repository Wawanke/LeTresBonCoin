package utils

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var templateshowPost *template.Template
var StructComm PostComm

//function that update information in the database

func Update(db *sql.DB, UserSturct Utilisateur, DataStruct Data, PostStruct Post) {

	updateUser, err := db.Prepare("UPDATE Utilisateur SET Nom_Utilisateur=?,Mot_de_passe=?,Adresse_mail=?,Follow=?, Followers=?, Description=?,Nb_Post=? WHERE ID=?  ")
	if err != nil {
		fmt.Println(err)
	}
	defer updateUser.Close()

	_, err = updateUser.Exec(UserSturct.Nom_dutilisateur, UserSturct.Mot_de_passe, UserSturct.Adresse_mail, UserSturct.Follow, UserSturct.Followers, UserSturct.Description, UserSturct.Nb_Post, UserSturct.ID)
	if err != nil {
		fmt.Print(err)
	}

	updateData, err := db.Prepare("UPDATE Data SET UID=?, Pays=?, Prenom=?, Nom=?, Job=?, Naissance=?  WHERE UID=?  ")
	if err != nil {
		fmt.Println(err)
	}
	defer updateData.Close()

	_, err = updateData.Exec(DataStruct.UID, DataStruct.Pays, DataStruct.Prenom, DataStruct.Nom, DataStruct.Job, DataStruct.Naissance, DataStruct.UID)
	if err != nil {
		fmt.Print(err)
	}

	updatePost, err := db.Prepare("UPDATE Post SET ID=?, UID=?, Name=?, Titre=?, Contenu=?, Date_Heure=?, Tag=?,Like=? , Nb_Com=? WHERE UID=?  ")
	if err != nil {
		fmt.Println(err)
	}
	defer updateData.Close()

	_, err = updatePost.Exec(PostStruct.ID, PostStruct.UID, PostStruct.Name, PostStruct.Titre, PostStruct.Contenu, PostStruct.Date_Heure, PostStruct.Tag, PostStruct.Like, PostStruct.Nb_Com, PostStruct.UID)
	if err != nil {
		fmt.Print(err)
	}

}

// function that emplty struct
func Logout(UserSturct Utilisateur, DataStruct Data, PostStruct Post) {

	db, _ := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")

	Update(db, UserSturct, DataStruct, PostStruct)

	UserSturct.ID = 0
	UserSturct.Nom_dutilisateur = ""
	UserSturct.Mot_de_passe = ResetTabBytes(UserSturct.Mot_de_passe)
	UserSturct.Adresse_mail = ""
	UserSturct.Follow = 0
	UserSturct.Followers = 0
	UserSturct.PP = ResetTabBytes(UserSturct.PP)
	UserSturct.Description = ""
	UserSturct.Nb_Post = 0

	fmt.Println(UserSturct)

	DataStruct.UID = 0
	DataStruct.Pays = ""
	DataStruct.Prenom = ""
	DataStruct.Nom = ""
	DataStruct.Job = ""
	DataStruct.Naissance = ""

	fmt.Println(DataStruct)

}

// function check if the argument was empty wen the user changer information
func EditProfile(r *http.Request, keyFormValue string) (string, bool) {
	newEntry := r.FormValue(keyFormValue)
	if newEntry != "" {
		return newEntry, true
	}
	return " ", false
}

// function that empty pp
func ResetTabBytes(slice []byte) []byte {
	for i := range slice {
		slice[i] = 0
	}
	return slice
}

func TestLogin(obtenu_site string, obtenu_bdd string, is_hash bool) bool {
	if !is_hash {
		if obtenu_site == obtenu_bdd {
			return true
		}
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(obtenu_bdd), []byte(obtenu_site))
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func Test(db *sql.DB, Nom_Utilisateur string, Mot_de_passe string, Adresse_mail string) {

	_, err := db.Exec("INSERT INTO Utilisateur (Nom_Utilisateur, Mot_de_passe, Adresse_mail, Follow, Followers )VALUES (?,?,?,?,?)", Nom_Utilisateur, Mot_de_passe, Adresse_mail, 0, 0)
	if err != nil {
		fmt.Println(err, "<-------------ICI ================================")
	}
}
func ActualDateInString() string {
	h, m, s := time.Now().Clock()
	if s == s {
	}
	year, month, day := time.Now().Date()
	DateEtHeure := fmt.Sprint(h) + "h" + fmt.Sprint(m) + " " + fmt.Sprint(year) + "/" + fmt.Sprint(month) + "/" + fmt.Sprint(day)
	return DateEtHeure
}
func fillStructPost(IDpost string) Post { // remplir une struct qui représente un post
	var data Post

	db, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")

	if err != nil {
		fmt.Println("C La Merde")
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT ID, Titre, Contenu, Date_Heure, Tag, Like FROM Post WHERE ID = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(IDpost)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ID int
		var UID int
		var titrePost string
		var Contenu string
		var Date_Heure string
		var Tag string
		err = rows.Scan(&ID, &UID, &titrePost, &Contenu, &Date_Heure, &Tag)
		if err != nil {
			log.Fatal(err)
		}
		data = Post{ID, UID, "", titrePost, Contenu, Date_Heure, Tag, 1, 0, []byte{}}

	}
	return data
}
func FillStructComm() []Commentaire {
	var P []Commentaire
	var TemP Commentaire
	bdd, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	// fmt.Println("Bdd Ouvert ")
	defer bdd.Close()
	stmt, err := bdd.Prepare("SELECT ID, ID, PID, Contenue FROM Commentaire;")
	if err != nil {
		panic(err)
	}
	// fmt.Println(stmt)
	rows, err := stmt.Query() // Faire que le scan se fait sur toute les colonnes des mail/mdp
	defer stmt.Close()
	// fmt.Println("avant------------------", rows)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&TemP.ID, &TemP.UID, &TemP.PID, &TemP.Contenue)
		if err != nil {
			log.Fatal(err)
		}
		P = append(P, TemP)

	}
	return P
}
func TrouveLeNom(r *http.Request) {
	for key, values := range r.Form { // range over map pour trouver quelle name et quelle montant il donne
		for _, value := range values { // range over []string
			fmt.Println(key, value)
		}
	}
}

func DBtableauDesPost() []Post {
	var P []Post
	var TemP Post
	bdd, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	// fmt.Println("Bdd Ouvert ")
	defer bdd.Close()
	stmt, err := bdd.Prepare("SELECT * FROM Post;")
	if err != nil {
		panic(err)
	}
	// fmt.Println(stmt)
	rows, err := stmt.Query() // Faire que le scan se fait sur toute les colonnes des mail/mdp
	defer stmt.Close()
	// fmt.Println("avant------------------", rows)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&TemP.ID, &TemP.UID, &TemP.Name, &TemP.Titre, &TemP.Contenu, &TemP.Date_Heure, &TemP.Tag, &TemP.Like, &TemP.Nb_Com)
		if err != nil {
			log.Fatal(err)
		}
		P = append(P, TemP)

	}
	return P
}
func DBajoutlike(tabPost []Post, id int) {
	bdd, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	if err != nil {
		fmt.Println(err, "BDD FAIled")
	}
	for x, _ := range tabPost {
		fmt.Println(">>>>>>>>>>>>>>", tabPost[x].Like, "<<<<<<<<<<<<<<<<<<")
		if tabPost[x].ID == id {
			templike := tabPost[x].Like + 1
			// cmd := "UPDATE userinfo SET created = ? WHERE uid = ?"
			fmt.Println("YESSSSSSSS" + " Ici la query   " + "UPDATE Post SET Like = '" + strconv.Itoa(templike) + "' WHERE ID =" + strconv.Itoa(id))

			vew, err := bdd.Exec("UPDATE Post SET Like = '" + strconv.Itoa(templike) + "' WHERE ID =" + strconv.Itoa(id))
			if err != nil {
				fmt.Println(err, "Failed du Exec")
				fmt.Println(vew)
			}
		}
	}
	bdd.Close()
}

func DBfindPost(tabPost []Post, id int) Post {
	var err Post
	for x, _ := range tabPost {
		fmt.Println(">>>>>>>>>>>>>>", tabPost[x].ID, "<<<<<<<<<<<<<<<<<<")
		if tabPost[x].ID == id {
			return tabPost[x]
		}
	}
	return err
}
func DBfindPostcomm(tabPost PostComm, id int) []Commentaire {
	var tempCom []Commentaire
	for x, _ := range tabPost.Contenue {
		fmt.Println(">>>>>>Tout les comme lier>>>>>>>>", tabPost.Contenue[x].ID, "<<<<<<<<<<Tout les comme lier<<<<<<<<")
		if tabPost.Contenue[x].PID == id {
			tempCom = append(tempCom, tabPost.Contenue[x])
		}
	}
	return tempCom
}
func DBsearchPost(tabPost []Post, id int) Post {
	var err Post
	for x, _ := range tabPost {
		fmt.Println(">>>>>>>>>>>>>>", tabPost[x].ID, "<<<<<<<<<<<<<<<<<<")
		if tabPost[x].ID == id {
			return tabPost[x]
		}
	}
	return err
}
func showPost(w http.ResponseWriter, r *http.Request) {
	templateshowPost.ExecuteTemplate(w, "showPost.html", nil)

	db, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")

	if err != nil {
		fmt.Println("C La Merde")
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT ID, Titre, Tag FROM Post")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
func Search(r *http.Request) []int {
	db, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	fmt.Println("BDD OUVERT")
	if err != nil {
		fmt.Println("C La Merde")
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT ID, Titre, Tag FROM Post")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	fmt.Println(r.FormValue("search"))
	searched := r.FormValue("search")
	var tabIDpost []int
	for rows.Next() {
		var IDpost int
		var titrePost string
		var tagPost string
		err := rows.Scan(&IDpost, &titrePost, &tagPost)
		if err != nil {
			log.Fatal(err)
		}

		inter_tab := traitement_tabIDpost(searchBarAlgo(searched, IDpost, titrePost, tagPost))
		if inter_tab != nil {
			tabIDpost = append(tabIDpost, inter_tab[0])
		}
	}
	tabIDpostAffiche := traitement_tabIDpost(tabIDpost)
	return tabIDpostAffiche
}

func searchBarAlgo(searched string, IDpost int, titrePost string, tagPost string) []int {
	tabIDpost := []int{}
	searchingtab := map[string][]string{" - Titre": []string{titrePost},
		" - Tag": []string{tagPost}}

	for typeOfData, searchValue := range searchingtab {
		if len(searchValue) == 1 && strings.Contains(strings.ToLower(searchValue[0]), strings.ToLower(string(searched))) {
			//fmt.Println(searchValue[0] + typeOfData) // affichage de l'auto complétion
			if typeOfData == typeOfData {
			}
			tabIDpost = append(tabIDpost, IDpost)
		} else {
			for _, searchValueString := range searchValue {
				if strings.Contains(strings.ToLower(searchValueString), strings.ToLower(string(searched))) {
					//fmt.Println(searchValueString + typeOfData) // affichage de l'auto complétion
					tabIDpost = append(tabIDpost, IDpost)
				}
			}
		}
	}
	return tabIDpost
}

func traitement_tabIDpost(tabIDpost []int) []int {
	var tIpTraitement []int

	for _, id := range tabIDpost {
		var isIn_tIpTraitement bool
		if len(tIpTraitement) == 0 {
			tIpTraitement = append(tIpTraitement, id)
		}
		for _, idTraitement := range tIpTraitement {
			if idTraitement == id {
				isIn_tIpTraitement = true
			}
		}
		fmt.Println(id, isIn_tIpTraitement, tIpTraitement)
		if !isIn_tIpTraitement {
			tIpTraitement = append(tIpTraitement, id)
		}
	}
	return tIpTraitement
}

func PagePost(tab []Post, w http.ResponseWriter, postID int) {

	StructComm.PostC = DBfindPost(tab, postID)
	fmt.Println(StructComm.PostC, "===============regarde bro===========")
	StructComm.Contenue = FillStructComm()
	fmt.Println(StructComm.Contenue, "===============regarde bro===========")
	numeroPostComm := DBfindPostcomm(StructComm, postID)
	StructComm.Contenue = numeroPostComm
	fmt.Println(StructComm.Contenue, "===============regarde bro APRES REcherche===========")
	custTemplate, err := template.ParseFiles("tmpl/inPost.html") //ici sera la page de erwan
	err = custTemplate.Execute(w, StructComm)
	if err != nil {

	}
}
func DBajoutcomm(commStr string, UID int, ID int, PID int) {
	db, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	ID++

	if err != nil {
		fmt.Println("C La Merde")
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Exec("INSERT INTO Commentaire ( UID, PID, Contenue)VALUES (?,?,?)", UID, PID, commStr)
	fmt.Println("BDD OUVERT")
	if err != nil {
		fmt.Println(err)
		fmt.Println(stmt, "<-------------ICI ================================")
	}
	fmt.Println(err, "<-------------ICI LA ================================")
}
func DBajoutPost(UID int, Titre string, Contenu string, Tag string) {
	db, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	if err != nil {
		fmt.Println("C La Merde")
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Exec("INSERT INTO Post ( ID, Titre, Contenu, Date_Heure, Tag, Like)VALUES (?,?,?,?,?,?)", UID, Titre, Contenu, "1200", Tag, 0)
	fmt.Println("BDD OUVERT")
	if err != nil {
		fmt.Println(err)
		fmt.Println(stmt, "<-------------ICI ================================")
	}
	fmt.Println(err, "<-------------ICI LA ================================")
}

func GetUserName(id int) string {

	var UID int
	var Name string

	db, err := sql.Open("sqlite3", "./DbDocker/DockerBack/db.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	getUser, err := db.Prepare("SELECT ID,Nom_Utilisateur FROM Utilisateur")
	if err != nil {
		fmt.Println("Error of get user form database :")
		panic(err)
	}
	defer getUser.Close()

	rows, err := getUser.Query()
	if err != nil {
		fmt.Println("Error of rows of get User :")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&UID, &Name)
		if err != nil {
			fmt.Println("Error of User rows scan :")
			log.Fatal(err)
		}

		if id == UID {
			return Name
		}

	}

	return ""
}
