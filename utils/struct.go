package utils

// struct of Utilistaeur
type Utilisateur struct {
	ID               int
	Nom_dutilisateur string
	Mot_de_passe     []byte
	Adresse_mail     string
	Follow           int
	Followers        int
	PP               []byte
	Description      string
	Nb_Post          int
}

type UserVisit struct {
	Nom_dutilisateur string
	Adresse_mail     string
	Follow           int
	Followers        int
	Description      string
	Nb_Post          int
}

// struct of Post
type Post struct {
	ID         int
	UID        int
	Name       string
	Titre      string
	Contenu    string
	Date_Heure string
	Tag        string
	Like       int
	Nb_Com     int
	UPP        []byte
}

// struct of data of user
type Data struct {
	UID       int
	Pays      string
	Prenom    string
	Nom       string
	Job       string
	Naissance string
}

type Commentaire struct {
	ID       int
	UID      int
	PID      int
	Contenue string
}

type PostComm struct {
	PostC    Post
	Contenue []Commentaire
}
