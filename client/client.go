package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"log"
	"bufio"

	st "./structures" // contient la structure Personne
	tr "./travaux" // contient les fonctions de travail sur les Personnes
)


const TAILLE_QUEUE int = 10


/*     A VOIR 




// Dans la fonction main() du serveur
ln, _ := net.Listen("tcp", addr)
for {
  conn, _ := ln.Accept()
  go connectionHandler(conn, canal)
}

// la fonction des handlers
func connectionHandler(conn net.Conn, toMaint chan message) {
  buff := make([]byte, 256)
  n, err := conn.Read(buff)
  if err != nil {
    return
  }
  msg := strings.Split(string(buff[:n]), "|")
  methode := msg[0]
  id, _ := strconv.Atoi(msg[1])
  retour := make(chan string)
  toMaint <- message{id, methode, retour}
  _, _ = conn.Write([]byte(<-retour))
  conn.Close()
}

// Le proxy local
func proxy(fromLocal chan message_proxy, port int) {
  addr := ADRESSE + ":" + fmt.Sprint(port)
  for {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
      continue
    }
    msg := <-fromLocal
    str := msg.methode + "|" + fmt.Sprint(msg.id)
    _, _ = conn.Write([]byte(str))
    reponse := make([]byte, 256)
    n, _ := conn.Read(reponse)
    msg.retour <- string(reponse[:n])
    conn.Close()
  }
}


*/

var ADRESSE string = "localhost"                           // adresse de base pour la Partie 2
var FICHIER_SOURCE string = "./conseillers-municipaux.txt" // fichier dans lequel piocher des personnes
var TAILLE_SOURCE int = 450000                             // inferieure au nombre de lignes du fichier, pour prendre une ligne au hasard
var TAILLE_G int = 5                                       // taille du tampon des gestionnaires
var NB_G int = 2                                           // nombre de gestionnaires
var NB_P int = 2                                           // nombre de producteurs
var NB_O int = 4                                           // nombre d'ouvriers
var NB_PD int = 2                                          // nombre de producteurs distants pour la Partie 2

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"} // une personne vide

type message_lec struct {
	contenu int //string //int  TODO à la base c'était un int
	retour chan string
}

// paquet de personne, sur lequel on peut travailler, implemente l'interface personne_int
type personne_emp struct {
	ligne string
	statut string
	lecteur chan message_lec
	afaire []func(st.Personne) st.Personne
	st.Personne
}

// paquet de personne distante, pour la Partie 2, implemente l'interface personne_int
type personne_dist struct {
	// A FAIRE
}

// interface des personnes manipulees par les ouvriers, les
type personne_int interface {
	initialise()          // appelle sur une personne vide de statut V, remplit les champs de la personne et passe son statut à R
	travaille()           // appelle sur une personne de statut R, travaille une fois sur la personne et passe son statut à C s'il n'y a plus de travail a faire
	vers_string() string  // convertit la personne en string
	donne_statut() string // renvoie V, R ou C
}

// fabrique une personne à partir d'une ligne du fichier des conseillers municipaux
// à changer si un autre fichier est utilisé
func personne_de_ligne(l string) st.Personne {
	separateur := regexp.MustCompile("\u0009") // oui, les donnees sont separees par des tabulations ... merci la Republique Francaise
	separation := separateur.Split(l, -1)
	naiss, _ := time.Parse("2/1/2006", separation[7])
	a1, _, _ := time.Now().Date()
	a2, _, _ := naiss.Date()
	agec := a1 - a2
	return st.Personne{Nom: separation[4], Prenom: separation[5], Sexe: separation[6], Age: agec}
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES ***

func (p *personne_emp) initialise() {
	print("arrivée dans la fonction initialise\n")
	ret := make(chan string)
	//print("a\n")
	//p.lecteur <- message_lec { contenu: p.ligne, retour: ret } // Ils font comme si p.contenu était du string
	//var err error;
	var contenu_temp int64;
	//print("b\n")
	contenu_temp, _ = strconv.ParseInt(p.ligne, 10, 32)
	//print("c\n")

	// Initilisation du lecteur
	//p.lecteur := make(chan message_lec)
	go func (chan message_lec){
		lecteur(p.lecteur)
	}(p.lecteur)

	p.lecteur <- message_lec { contenu: int(contenu_temp), retour: ret } 
	//p.lecteur := <- canal_lecteur 
	
	//print("d\n")
	ligne := <- ret
	//print("e\n")
	p.Personne = personne_de_ligne(ligne)
	//print("f\n")
	for i := 0; i < rand.Intn(6)+1; i++ {
		//print("g\n")
		p.afaire = append(p.afaire, tr.UnTravail())
		//print("h\n")
	}
	//print("i\n")
	p.statut = "R"
	//print("j\n")
}

func (p *personne_emp) travaille() {
	p.Personne = p.afaire[0](p.Personne)
	p.afaire = p.afaire[1:]
	if len(p.afaire) == 0 {
		p.statut = "C"
	}
}

func (p *personne_emp) vers_string() string {
	var add string
	if p.Sexe == "F" {
		add = "Madame "
	} else {
		add = "Monsieur "
	}
	return fmt.Sprint(add, p.Prenom, " ", p.Nom, " : ", p.Age, " ans.")
}

func (p *personne_emp) donne_statut() string {
	return p.statut
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES DISTANTES (PARTIE 2) ***
// ces méthodes doivent appeler le proxy (aucun calcul direct)

func (p personne_dist) initialise() {
	// A FAIRE
}

func (p personne_dist) travaille() {
	// A FAIRE
}

func (p personne_dist) vers_string() string {
	// A FAIRE
	return ""
}

func (p personne_dist) donne_statut() string {
	// A FAIRE
	return ""
}

// *** CODE DES GOROUTINES DU SYSTEME ***

// Partie 2: contacté par les méthodes de personne_dist, le proxy appelle la méthode à travers le réseau et récupère le résultat
// il doit utiliser une connection TCP sur le port donné en ligne de commande
func proxy() {
	// A FAIRE
}

// Partie 1 : contacté par la méthode initialise() de personne_emp, récupère une ligne donnée dans le fichier source
func lecteur(url chan message_lec) {
	for {
		m := <- url
		fmt.Println("Lecteur contacté pour ligne", m.contenu)
		fichier, err := os.Open(FICHIER_SOURCE)
		//print("aa\n")
		if err != nil {
			log.Fatal(err)
		}
		//print("bb\n")

		scanner := bufio.NewScanner(fichier)
		//print("cc\n")
		_ = scanner.Scan()
		//print("dd\n")

		for i := 0; i < 1; i++ { // i < m.contenu mais ça fait bugger et de tt façon ça fait rien
			//print("ee%d\n", i)
			_ = scanner.Scan()
			//print("ff\n")
		}

		resultat := scanner.Scan()
		//print("gg\n")
		if resultat == false {
			//print("hh\n")
			log.Fatal(err)
			//print("ii\n")
		} else {
			//print("jj\n")
			m.retour <- scanner.Text()
			//print("kk\n")
		}
		//print("find de la lecture, fermeture du fichier\n")
		fichier.Close()
		
	}
}

// Partie 1: récupèrent des personne_int depuis les gestionnaires, font une opération dépendant de donne_statut()
// Si le statut est V, ils initialise le paquet de personne puis le repasse aux gestionnaires
// Si le statut est R, ils travaille une fois sur le paquet puis le repasse aux gestionnaires
// Si le statut est C, ils passent le paquet au collecteur
func ouvrier(canal_ouvrier chan personne_int, canal_gestionnaire chan personne_int, canal_collecteur chan personne_int) {
	//fmt.Println("lancement d'ouvrier")
	for{
		//personne := <- canal_ouvrier

		select {
		case personne := <- canal_ouvrier :

				if personne.donne_statut() == "V" {
					print("reception d'un paquet V\n")
					personne.initialise()
					print("renvoie du paquet initialisé au gestionnaire\n")
					canal_gestionnaire <- personne
				} else if personne.donne_statut() == "R" {
					print("reception d'un paquet R\n")
					personne.travaille()
					print("renvoie du paquet travaillé au gestionnaire\n")
					canal_gestionnaire <- personne
				} else if personne.donne_statut() == "C" {
					print("reception d'un paquet C, envoie au collecteur\n")
					canal_collecteur <- personne
				} else {
					print("reception d'un paquet de type inconnu\n")
				}
			
		}
	}
	
}

// Partie 1: les producteurs cree des personne_int implementees par des personne_emp initialement vides,
// de statut V mais contenant un numéro de ligne (pour etre initialisee depuis le fichier texte)
// la personne est passée aux gestionnaires
func producteur(canal_gestionnaire chan personne_int /*, canal_lecteur chan message_lec*/) {
	for {
		np := pers_vide
		nt := make([]func(st.Personne) st.Personne, 0)

		/*
			ligne string
			statut string
			lecteur chan message_lec
			afaire []func(st.Personne) st.Personne
			st.Personne
		*/
		npe := personne_emp{
			statut: "V",
			ligne: strconv.Itoa(rand.Intn(TAILLE_SOURCE)),
			afaire: nt,
			Personne: np,
			//lecteur: canal_lecteur, 
			lecteur:
			make(chan message_lec) }
		
		//print("Producteur crée une ligne\n")

		fmt.Println("Producteur crée une ligne", npe.ligne)
		canal_gestionnaire <- personne_int(&npe) // est ce qu'il ne faut pas qu'il puisse y avoir plusieurs channels de gestionnaires ?
	}
}

// Partie 2: les producteurs distants cree des personne_int implementees par des personne_dist qui contiennent un identifiant unique
// utilisé pour retrouver l'object sur le serveur
// la creation sur le client d'une personne_dist doit declencher la creation sur le serveur d'une "vraie" personne, initialement vide, de statut V
func producteur_distant() {
	// A FAIRE
}

// Partie 1: les gestionnaires recoivent des personne_int des producteurs et des ouvriers et maintiennent chacun une file de personne_int
// ils les passent aux ouvriers quand ils sont disponibles
// ATTENTION: la famine des ouvriers doit être évitée: si les producteurs inondent les gestionnaires de paquets, les ouvrier ne pourront
// plus rendre les paquets sur lesquels ils travaillent pour en prendre des autres
func gestionnaire(canal_gestionnaire chan personne_int, canal_ouvrier chan personne_int, file_personnes []personne_int) {
	for {
		personne := <- canal_gestionnaire
		//var file_personnes []personne_emp

		print("reception d'un paquet dans la fonction gestionnaire\n")
		
		if(len(file_personnes) <= (TAILLE_QUEUE / 2)){
			print("insertion du paquet dans la file\n")
			file_personnes = append(file_personnes, personne)
		} 

		var temp personne_int = file_personnes[0]
		file_personnes = file_personnes[1:]

		print("envoie du paquet dans le canal ouvrier\n")
		canal_ouvrier <- temp
	}
}

// Partie 1: le collecteur recoit des personne_int dont le statut est c, il les collecte dans un journal
// quand il recoit un signal de fin du temps, il imprime son journal.
func collecteur(canal_collecteur chan personne_int, fintemps chan int) {
	for{
		select {
			case <- fintemps:
				print("signal de fin reçu dans collecteur\n")
				fintemps <- 0
				return
			case personne := <- canal_collecteur:
				// Temporaire, juste pour voir si ça marche
				print(personne.vers_string())
				print("\n")

		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // graine pour l'aleatoire
	/*

	// Pour la partie 2

	if len(os.Args) < 3 {
		fmt.Println("Format: client <port> <millisecondes d'attente>")
		return
	}
	*/
	//port, _ := strconv.Atoi(os.Args[1]) // utile pour la partie 2
	//millis, _ := strconv.Atoi(os.Args[2]) // duree du timeout 
	var millis int = 10000
	fintemps := make(chan int)



	// A FAIRE
	// creer les canaux
	// lancer les goroutines (parties 1 et 2): 1 lecteur, 1 collecteur, des producteurs, des gestionnaires, des ouvriers
	// lancer les goroutines (partie 2): des producteurs distants, un proxy
	canal_gestionnaire 	:= make(chan personne_int)
	canal_ouvrier 		:= make(chan personne_int)
	canal_collecteur 	:= make(chan personne_int)
	//canal_lecteur		:= make(chan message_lec)

	var file_personnes []personne_int
	

	// Initialisation du gestionnaire
	go func (chan personne_int, chan personne_int, []personne_int){
		gestionnaire(canal_gestionnaire, canal_ouvrier, file_personnes)
	}(canal_gestionnaire, canal_ouvrier, file_personnes)

	// Initialisation des l'ouvriers
	for i := 0; i < TAILLE_QUEUE; i++ {
		go func (chan personne_int, chan personne_int, chan personne_int){
			ouvrier(canal_ouvrier, canal_gestionnaire, canal_collecteur)
		}(canal_ouvrier, canal_gestionnaire, canal_collecteur)
	}
	


	// Initilisation du collecteur
	go func (chan personne_int, chan int){
		collecteur(canal_collecteur, fintemps)
	}(canal_collecteur, fintemps)

	// Initilisation du producteur
	go func (chan personne_int/*, chan message_lec*/){
		producteur(canal_gestionnaire/*, canal_lecteur*/)
	}(canal_gestionnaire/*, canal_lecteur*/)

	print("fin du lancement de tous les elements\n")

	


		
	time.Sleep(time.Duration(millis) * time.Millisecond)
	fintemps <- 0
	<-fintemps
}
