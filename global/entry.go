package global

//On définit ici tout ce qui compose une entrée (une porte, deux reader, l'alarme de la porte, et le sensor)

//on crée une entry à chaque fois que l'on reçois la socket (badge lu)
//l'entry va faire la requete api ("/canAccess") avec les paramètres idbadge et idbatiment
//on lance la fonction qui dirige la porte avec le résultat de la requete (on bloque ou débloque la porte) au bout de 30 seconde ferme la porte et écrit porte fermer dans l'espace tuple
//on lance la fonction qui gere les reader (la couleur) elle envoie une socket avec la couleur, et au bout de 5 secondes pour vert et 10 pour rouge une autre pour eteindre
//on lance la function sensor qui va nous dire ecrire dans l'espace tuple. entry écoute, et si + de deux recu, alarme immédiate (si alarme, on ecrit dans espace tuple pour fermer la porte et on ferme l'entry et on enregistre l'incident)
//quand la porte est refermée (recu dans l'espace tuple), l'entry est supprimée, et elle est consignée dans le journal de bord (requete api pour l'enregistrer) faire attention a si la personne est rentrée ou pas

func entry(idPorte string, idBatiment string, idBadge string) {
	print("Je créé une entrée avec :", idPorte, idBatiment, idBadge)
}
