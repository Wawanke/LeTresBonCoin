package main

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// var
var User = utils.Utilisateur{}
var Data = utils.Data{}
var Post = utils.Post{}
var UserVistit = utils.UserVisit{}
var Commentaire = utils.Commentaire{}
var P AllPost
var postUIDtest string
var giveLike int
var postIDt int
var postIDcID int
var postUID int

type AllPost struct {
	ID   int
	Post []utils.Post
}

//func of handlers

func ProfileVisitor(w http.ResponseWriter, r *http.Request) {

	//get UID of this post

	//open database

	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println(err)
	}

	//prepare query of chart Utilisateur

	getUser, err := db.Prepare("SELECT ID, Nom_Utilisateur, Adresse_mail, Follow, Followers, Description, Nb_Post FROM Utilisateur")
	if err != nil {
		fmt.Println(err)
	}
	defer getUser.Close()

	// Scan all columns selected of Utilisateur

	rows, err := getUser.Query()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	//next rows to scan

	for rows.Next() {

		//var that stock temporaly the data of this rows scan

		var id int
		var username string
		var mail_user string
		var follow int
		var followers int
		var description string
		var nbPost int

		//attribute the values of the scan to the var

		err = rows.Scan(&id, &username, &mail_user, &follow, &followers, &description, &nbPost)
		if err != nil {
			fmt.Println("Error of User rows scan :")
			log.Fatal(err)
		}

		//check if the uid of post correspond to id of user stock in the databese

		if postUID == id {

			// stock all of selected information stock in database of this user in structs

			UserVistit.Nom_dutilisateur = username
			UserVistit.Adresse_mail = mail_user
			UserVistit.Follow = follow
			UserVistit.Followers = followers
			UserVistit.Description = description
			UserVistit.Nb_Post = nbPost

		}

	}

	//parse html

	tmpl, err := template.ParseFiles("tmpl/profileVisitor.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//exectue template with a struct in data

	err = tmpl.Execute(w, UserVistit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Home(w http.ResponseWriter, r *http.Request) {

	var postID int
	custTemplate, err := template.ParseFiles("tmpl/home.html")
	P.Post = utils.DBtableauDesPost()

	if r.Method == "POST" {

		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		utils.TrouveLeNom(r)
		if r.FormValue("user") != "" {
			postUID, _ = strconv.Atoi(r.FormValue("user"))
			ProfileVisitor(w, r)
		} else if r.FormValue("search") != "" {

			temTabPost := utils.Search(r)
			fmt.Println("ici est mon tab ", temTabPost)
			var temTabPostPost AllPost
			for x := 0; x < len(temTabPost); x++ {
				temTabPostPost.Post = append(temTabPostPost.Post, utils.DBfindPost(P.Post, temTabPost[x]))
				fmt.Println("ICI LE COMPTE-----   ", x, "  ICI LE TAB A JOUR  ", temTabPostPost)
			}
			custTemplate, err = template.ParseFiles("tmpl/home.html")
			err = custTemplate.Execute(w, temTabPostPost)
		} else if r.FormValue("like") != "" {
			giveLike, _ = strconv.Atoi(r.FormValue("like"))
			//custTemplate, err = template.ParseFiles("tmpl/home.html") //ici sera la page de erwan
			//err = custTemplate.Execute(w, P)
			utils.TrouveLeNom(r)
			utils.DBajoutlike(P.Post, giveLike)
			fmt.Println("Like + ", giveLike, " a était ajouté")
			custTemplate, err = template.ParseFiles("tmpl/home.html") //ici sera la page de erwan
			err = custTemplate.Execute(w, P)
		} else if r.FormValue("lettre") != "" {
			postID, _ = strconv.Atoi(r.FormValue("lettre"))
			postIDt = postID
			utils.PagePost(P.Post, w, postID)
			utils.TrouveLeNom(r)
		} else if r.FormValue("ajout") != "" {
			utils.TrouveLeNom(r)
			fmt.Println(postIDt, "VERIFIE MOI SI ICI IL Y GARDE LE ID", r.FormValue("ajout"))
			utils.PagePost(P.Post, w, postIDt)
		} else if r.FormValue("likep") != "" {
			giveLike = postIDt
			utils.DBajoutlike(P.Post, giveLike)
			utils.PagePost(P.Post, w, postIDt)
		} else {
			custTemplate, err = template.ParseFiles("tmpl/home.html") //ici sera la page de erwan
			err = custTemplate.Execute(w, P)
		}

	}

	if err != nil {

	}

	if r.Method == "GET" {
		err = custTemplate.Execute(w, P)
	}

	fmt.Println("======UIDtest======== " + postUIDtest + " ===========UIDtest========")
	fmt.Println("========MethUID====== " + r.Method + " =========MethUID==========")

}

func HomeGuest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/homeGuest.html")
}

func Login(w http.ResponseWriter, r *http.Request) {

	//emptying the struct

	utils.Logout(User, Data, Post)

	// var that get input of the user during the signup

	inputUsername := r.FormValue("username")
	inputEmail := r.FormValue("email_signup")
	inputPswd := r.FormValue("pswd_signup")

	pswd_hash, err := bcrypt.GenerateFromPassword([]byte(inputPswd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error of encryption :")
		fmt.Println(err)
	}

	//open database

	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println("Error of open Database :")
		fmt.Println(err)
	}
	defer db.Close()

	//checks if the user has entered getData and if true insert their log in the database

	if inputUsername != "" && inputEmail != "" && inputPswd != "" {

		//prepare and execute User log in the database

		createUser, err := db.Prepare("INSERT INTO Utilisateur (Nom_Utilisateur, Mot_de_passe, Adresse_mail) VALUES (?, ?, ?)")
		if err != nil {
			fmt.Println("Error of create user :")
			panic(err)
		}

		test, err_createUserfunc := createUser.Exec(inputUsername, pswd_hash, inputEmail)
		if err_createUserfunc != nil {
			fmt.Println("Error of open Database :")
			fmt.Println(err_createUserfunc)
		}

		//get the id of the chart Utilisateur that auto increment by the database. To set the link between the cart Utilisateur and data in the column UID of chart Data

		id, _ := test.LastInsertId()

		//prepare and execute User empty Data in the database and link the id user in UID

		createData, err2 := db.Prepare("INSERT INTO Data (UID, Naissance, Pays, Nom, Prenom, Job) VALUES (?, ?, ?, ?, ?, ?)")
		if err2 != nil {
			fmt.Println("Error of create Data :")
			fmt.Println(err2)
		}

		createData.Exec(id, "", "", "", "", "")

		//debug
		fmt.Print(inputUsername, "/ /", inputPswd, "/ /", inputEmail)
	}

	//var for check if the getData input correspond to getData stock in the database

	checkEmail := r.FormValue("email_login")
	checkPassword := r.FormValue("pswd_login")

	//debug

	fmt.Println(checkEmail, checkPassword)

	//prepare query of chart Utilisateur and chart Data

	getUser, err := db.Prepare("SELECT ID,Nom_Utilisateur,Mot_de_passe,Adresse_mail,Follow, Followers, Description, Nb_Post FROM Utilisateur")
	if err != nil {
		fmt.Println("Error of get user form database :")
		panic(err)
	}
	defer getUser.Close()

	getData, err2 := db.Prepare("SELECT UID, Pays, Prenom, Nom, Job, Naissance FROM Data")
	if err2 != nil {
		fmt.Println("Error of get data form database :")
		fmt.Println(err2)
	}
	defer getData.Close()

	// Scan all columns selected of Utilisateur

	rows, err := getUser.Query()
	if err != nil {
		fmt.Println("Error of rows of get User :")
		panic(err)
	}
	defer rows.Close()

	//next rows to scan

	for rows.Next() {

		//var that stock temporaly the data of this rows scan

		var id int
		var username string
		var pswd_user string
		var mail_user string
		var follow int
		var followers int
		var pp []byte
		var description string
		var nbPost int

		//attribute the values of the scan to the var

		err = rows.Scan(&id, &username, &pswd_user, &mail_user, &follow, &followers, &description, &nbPost)
		if err != nil {
			fmt.Println("Error of User rows scan :")
			log.Fatal(err)
		}

		//debug

		fmt.Println("mail : ", mail_user, " pswd : ", pswd_user)

		//check if log option was in the database

		if (checkEmail != "" || checkPassword != "") && (utils.TestLogin(checkPassword, pswd_user, true) && utils.TestLogin(checkEmail, mail_user, false)) {

			//debug

			fmt.Println("testUser : success", pswd_user, checkPassword)

			// script was redirect in the home page if the log option was in the database

			redirectScript := fmt.Sprintf(`<script>window.location.href = "/home";</script>`)
			w.Write([]byte(redirectScript))

			// stock all of selected information stock in database of this user in structs

			User.ID = id
			User.Nom_dutilisateur = username
			User.Mot_de_passe = []byte(pswd_user)
			User.Adresse_mail = mail_user
			User.Follow = follow
			User.Followers = followers
			User.PP = pp
			User.Description = description
			User.Nb_Post = nbPost

			// Scan all columns selected of Data

			row2, err3 := getData.Query()
			if err3 != nil {
				fmt.Println(err3)
			}
			defer row2.Close()

			//next row to scan

			for row2.Next() {

				//var that stock temporaly the data of this rows scan

				var uid int
				var pays string
				var prenom string
				var nom string
				var job string
				var naissance string

				//attribute the values of the scan to the var

				err := row2.Scan(&uid, &pays, &prenom, &nom, &job, &naissance)
				if err != nil {
					fmt.Println(err)
				}

				// stock all of selected information stock in database of this user in structs

				Data.UID = uid
				Data.Pays = pays
				Data.Prenom = prenom
				Data.Nom = nom
				Data.Job = job
				Data.Naissance = naissance

			}

		} else {

			//debug

			fmt.Println("Wrong email/password ", "user : ", pswd_user, "testing : ", checkPassword)
		}
	}

	http.ServeFile(w, r, "tmpl/login.html")

}

func NewPost(w http.ResponseWriter, r *http.Request) {

	newPostTitle := r.FormValue("ttl")
	newPostMsg := r.FormValue("msg")
	newPostTag := r.FormValue("tag")
	DateEtHeure := utils.ActualDateInString()
	UID := User.ID
	name := utils.GetUserName(UID)

	if newPostTitle != "" && newPostMsg != "" {
		fmt.Println("Contenu newPost : ", newPostTitle, newPostMsg, newPostTag, DateEtHeure)
	}

	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println("C La Merde")
	}
	defer db.Close()

	if newPostTitle != "" && newPostMsg != "" {
		statementInfos, err := db.Prepare("INSERT INTO Post (UID,Name, Titre, Contenu, Date_Heure, Tag, Like) VALUES (?,?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}
		statementInfos.Exec(UID, name, newPostTitle, newPostMsg, DateEtHeure, newPostTag, 0)
		User.Nb_Post += 1
		utils.Update(db, User, Data, Post)
	}

	http.ServeFile(w, r, "tmpl/newPost.html")
}

func Profile(w http.ResponseWriter, r *http.Request) {

	//parse html

	tmpl, err := template.ParseFiles("tmpl/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//exectue template with a struct in data

	err = tmpl.Execute(w, User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProfileSettings(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "tmpl/profileSettings.html")

	//open database

	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println(err)
	}

	//get argumente enter in the html an attribute in the good struct

	userName, bool := utils.EditProfile(r, "username")
	if bool == true {
		User.Nom_dutilisateur = userName
		Post.Name = userName
	}

	password, bool := utils.EditProfile(r, "pswd")
	if bool == true {
		User.Mot_de_passe, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	}

	aboutYou, bool := utils.EditProfile(r, "description")
	if bool == true {
		User.Description = aboutYou
	}

	country, bool := utils.EditProfile(r, "country")
	if bool == true {
		Data.Pays = country
	}

	eMail, bool := utils.EditProfile(r, "email")
	if bool == true {
		User.Adresse_mail = eMail
	}

	birth, bool := utils.EditProfile(r, "birth")
	if bool == true {
		Data.Naissance = birth
	}

	surname, bool := utils.EditProfile(r, "surname")
	if bool == true {
		Data.Prenom = surname
	}

	name, bool := utils.EditProfile(r, "name")
	if bool == true {
		Data.Nom = name
	}

	job, bool := utils.EditProfile(r, "job")
	if bool == true {
		Data.Job = job
	}

	//debug

	fmt.Println(Data, User)

	//update the database with the news information

	utils.Update(db, User, Data, Post)

}

func InPost(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/inPost.html")
}

func InPostGuest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/inPost.html")
}

const port = ":8080"

func main() {

	//redirection with her handlers assosiate

	http.HandleFunc("/home", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/post", NewPost)
	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/homeGuest", HomeGuest)
	http.HandleFunc("/profileSettings", ProfileSettings)
	http.HandleFunc("/profileVisitor", ProfileVisitor)
	http.HandleFunc("/inPost", InPost)
	http.HandleFunc("/inPostGuest", InPostGuest)

	//

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	//print url

	fmt.Println("\n(http://localhost"+port+"/login) -  Server started on port", port)

	// start server
	http.ListenAndServe(port, nil)

}
